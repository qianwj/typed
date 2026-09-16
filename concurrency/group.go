package concurrency

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/qianwj/typed/adt"
)

// Mode selects the failure policy of a [Group].
type Mode int

const (
	// Strict mode (default) is the standard "errgroup" semantics: when any
	// task returns a non-nil error, the group's ctx is canceled, the
	// in-flight siblings see the cancellation and abort, and [Group.Wait]
	// returns that error. Subsequent calls to [Group.Go] after the first
	// error still spawn goroutines but their results are ignored by
	// [Group.Wait].
	Strict Mode = iota

	// BestEffort mode runs every task to its natural end (or until the
	// parent ctx is canceled), regardless of how many siblings fail. The
	// group's ctx is NOT canceled on a task error. [Group.Wait] blocks
	// until all spawned tasks have returned, then returns:
	//   - nil if every task succeeded (including the trivial case of
	//     zero tasks spawned);
	//   - the bare single task error if exactly one task failed;
	//   - a [*BestEffortError] wrapping every failed task otherwise.
	//
	// Parent-ctx cancellation does not change the shape of the returned
	// error: any task that was still running when the parent ctx was
	// canceled simply returns ctx.Err() as its error, and those errors
	// are collected just like any other failure.
	BestEffort
)

// String returns the mode's name ("Strict" or "BestEffort") for
// diagnostics and error messages.
func (m Mode) String() string {
	switch m {
	case Strict:
		return "Strict"
	case BestEffort:
		return "BestEffort"
	default:
		return fmt.Sprintf("Mode(%d)", int(m))
	}
}

// Option configures a [Group] at construction time. Use [WithMode] or
// [WithLimit].
type Option func(*groupConfig)

type groupConfig struct {
	mode  Mode
	limit int
}

// WithMode selects the group's failure policy. The default is [Strict].
// See [Mode] for the semantics of each.
func WithMode(m Mode) Option {
	return func(c *groupConfig) { c.mode = m }
}

// WithLimit caps the number of concurrently-running tasks spawned via
// [Group.Go] at n. Calls to [Group.Go] beyond the cap block until a slot
// is freed. A non-positive n means no limit.
func WithLimit(n int) Option {
	return func(c *groupConfig) { c.limit = n }
}

// Group is a typed wrapper around "spawn N goroutines, wait for them all,
// return one error" with two failure policies selectable via [WithMode].
//
// Use [NewGroup] to construct one. The zero value is not usable.
type Group struct {
	cfg groupConfig

	wg sync.WaitGroup

	// sem is nil if no limit was configured, otherwise it has capacity
	// cfg.limit and slots are sent on at the start of each Go call and
	// received from at the end. Sending on a full chan blocks; that's the
	// concurrency cap.
	sem chan struct{}

	ctx    context.Context // derived from parent in NewGroup, canceled by cancel()
	cancel context.CancelFunc

	mu    sync.Mutex
	state groupState

	// doneOnce / doneCh back [Group.Wait] and [Group.WaitWithContext].
	// doneCh is closed when every spawned task has returned; it is
	// lazily allocated on the first call to [Group.doneChan], so the
	// watcher goroutine is spawned at most once per Group regardless of
	// how many times Wait / WaitWithContext is invoked.
	doneOnce sync.Once
	doneCh   chan struct{}
}

// groupState is the result-tracking half of Group, protected by mu.
//
// In Strict mode, only firstErr matters; errs is always nil.
// In BestEffort mode, firstErr is always nil and errs accumulates every
// failed task in completion order. len(errs)==0 at Wait time means every
// task succeeded (or no task was spawned) → Wait returns nil.
type groupState struct {
	firstErr error
	errs     []error
}

// NewGroup returns a fresh Group whose derived ctx is canceled either by
// the parent ctx, by the first task error (in [Strict] mode), or by
// [Group.Wait]'s completion signaling. Pass [WithMode] / [WithLimit] to
// configure the failure policy and concurrency cap respectively.
//
// A nil parent ctx is treated as [context.Background].
func NewGroup(parent context.Context, opts ...Option) *Group {
	if parent == nil {
		parent = context.Background()
	}
	cfg := groupConfig{mode: Strict}
	for _, opt := range opts {
		opt(&cfg)
	}

	ctx, cancel := context.WithCancel(parent)

	g := &Group{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
	if cfg.limit > 0 {
		g.sem = make(chan struct{}, cfg.limit)
	}
	return g
}

// Go spawns a new task. fn is called in a fresh goroutine with the
// group's derived ctx; the task should respect that ctx's cancellation
// and return promptly.
//
// In [Strict] mode, when fn returns a non-nil error, the group's ctx is
// canceled (siblings see it), and [Group.Wait] will return that error.
//
// In [BestEffort] mode, the ctx is NOT canceled. Every error is
// recorded; [Group.Wait] returns nil only if every task succeeded, and
// a bare error or [*BestEffortError] otherwise.
//
// Go is safe to call from multiple goroutines concurrently. Calls to Go
// after the first error are still accepted (the goroutines run); they
// just have no effect on [Group.Wait]'s return value.
func (g *Group) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)

	if g.sem != nil {
		g.sem <- struct{}{}
	}

	go func() {
		defer func() {
			if g.sem != nil {
				<-g.sem
			}
			g.wg.Done()
		}()

		err := fn(g.ctx)
		g.handleResult(err)
	}()
}

// handleResult is called once per task on its return path. It records the
// task's outcome under mu and (in Strict mode) cancels the group ctx on
// the first error.
func (g *Group) handleResult(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err == nil {
		// Success: nothing to record in either mode. In BestEffort
		// mode, len(errs)==0 at Wait time means "all succeeded",
		// which collapses to nil.
		return
	}

	if g.cfg.mode == Strict {
		// First error wins; cancel siblings. Subsequent errors are
		// ignored (the first one is already stable in firstErr and
		// ctx is already canceled).
		if g.state.firstErr != nil {
			return
		}
		g.state.firstErr = err
		g.cancel()
		return
	}

	// BestEffort: every failure is recorded.
	g.state.errs = append(g.state.errs, err)
}

// Wait blocks until every task spawned by [Group.Go] has returned, then
// returns the error that best describes the outcome:
//
//   - In [Strict] mode: the first task error, or nil if every task
//     succeeded (including the trivial case of zero tasks).
//   - In [BestEffort] mode: nil if every task succeeded; the bare
//     single task error if exactly one task failed; a [*BestEffortError]
//     wrapping every failed task if more than one task failed.
//
// Wait also cancels the Group's derived ctx on return, which releases
// the cancelCtx from any parent's children map and lets any
// propagateCancel watcher goroutine exit. The cancellation is
// idempotent and harmless for [Strict] mode where the first task
// error has already canceled the ctx.
//
// Wait is safe to call multiple times; subsequent calls return the
// same value once all tasks have finished.
func (g *Group) Wait() error {
	defer g.cancel()

	<-g.doneChan()

	g.mu.Lock()
	defer g.mu.Unlock()
	return g.snapshotError()
}

// WaitWithContext is the ctx-aware variant of [Group.Wait]. It blocks
// until every spawned task has returned OR ctx is canceled or its
// deadline expires, whichever comes first.
//
// When ctx fires first, WaitWithContext returns ctx.Err() and does NOT
// kill the in-flight goroutines — Go has no safe kill primitive, and
// the toolkit's policy is "abandon, don't kill". Tasks that respect
// the Group's ctx will see it canceled (see "Cleanup" below) and
// exit promptly; tasks that ignore the ctx will continue running, and
// callers who want the eventual outcome can call [Group.Wait] after
// this method returns.
//
// When all spawned tasks have finished first, the return value is the
// same outcome-based error that [Group.Wait] would return — the ctx
// only governs the wait, not the result. Once all tasks are done,
// every subsequent call to WaitWithContext (or [Group.Wait]) returns
// the same outcome-based error, regardless of whether ctx is still
// alive. This is implemented with a non-blocking "tasks already done"
// probe before the blocking ctx race, so the deterministic "tasks
// done" case wins over a concurrently-canceled ctx instead of being
// resolved by Go's randomly-tied select.
//
// A nil ctx is treated as [context.Background].
//
// Cleanup: like [Group.Wait], this method always cancels the Group's
// derived ctx on return (idempotent with any Strict-mode error-path
// cancel already done), so the cancelCtx is released from any
// parent's children map and any propagateCancel watcher goroutine
// started by [context.WithCancel] against a non-cancelCtx parent can
// exit. This is the standard library's own cleanup pattern (see
// golang.org/x/sync/errgroup); without it the Group would leak one
// propagateCancel goroutine per instantiation when the parent is a
// custom non-cancelCtx Context.
func (g *Group) WaitWithContext(ctx context.Context) error {
	defer g.cancel()

	if ctx == nil {
		ctx = context.Background()
	}
	// Fast path: if every task is already done, skip the ctx race so
	// that an already-canceled ctx doesn't steal the outcome from a
	// finished Group (Go's select picks randomly when both cases are
	// ready; the contract here is "done beats canceled").
	select {
	case <-g.doneChan():
		// All tasks finished; fall through to the shared state-read
		// path used by [Group.Wait].
	default:
		select {
		case <-g.doneChan():
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.snapshotError()
}

// doneChan returns a channel that is closed when every task spawned by
// [Group.Go] has returned. It is safe to call concurrently and
// repeatedly; a single watcher goroutine is spawned per Group on the
// first call, and that goroutine terminates naturally when the
// underlying [sync.WaitGroup] reaches zero.
func (g *Group) doneChan() <-chan struct{} {
	g.doneOnce.Do(func() {
		g.doneCh = make(chan struct{})
		go func() {
			g.wg.Wait()
			close(g.doneCh)
		}()
	})
	return g.doneCh
}

// snapshotError reads the group's result state and returns the error
// that best describes it. The caller must hold g.mu.
func (g *Group) snapshotError() error {
	if g.cfg.mode == Strict {
		return g.state.firstErr
	}

	// BestEffort.
	switch len(g.state.errs) {
	case 0:
		// Every task succeeded (or no task was spawned).
		return nil
	case 1:
		// Common case: exactly one failure. Return the bare error
		// so errors.Is / errors.As match without unwrapping.
		return g.state.errs[0]
	default:
		return &BestEffortError{Errors: g.state.errs}
	}
}

// BestEffortError aggregates the errors from a [Group] running in
// [BestEffort] mode whose tasks all failed. It implements [errors.Is]
// and [errors.As] via [errors.Unwrap]: callers can match any underlying
// error with errors.Is(err, targetErr).
type BestEffortError struct {
	// Errors holds the per-task errors in completion order. Always
	// non-empty when Wait returns a *BestEffortError.
	Errors []error
}

// Error returns a single-line summary listing every task error in order.
// For a single error, it returns that error's message verbatim (so
// errors.Is on the wrapped error just works). For multiple errors, the
// message is prefixed with the count.
func (e *BestEffortError) Error() string {
	if len(e.Errors) == 0 {
		return "concurrency: BestEffortError with no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "concurrency: %d tasks failed: ", len(e.Errors))
	for i, err := range e.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Unwrap returns the slice of underlying errors so [errors.Is] and
// [errors.As] can walk through them. The first non-nil error is what
// [errors.Is] compares against; a target matching any of the Errors
// causes Is to return true.
func (e *BestEffortError) Unwrap() []error {
	return e.Errors
}

// Get returns the i-th task error in completion order as an
// [adt.Option]. The result is present when 0 <= i < len(e.Errors)
// and empty otherwise (no panic, no out-of-range signal — same
// "absent value is observable via IsEmpty" rule as
// [BoundedBlockingQueue.TryPoll]).
//
// Most callers won't need this: [errors.Is] / [errors.As] walk every
// element via [BestEffortError.Unwrap] without indexing. Get exists
// for the rare case where the caller wants positional access — for
// example, formatting "task #3 failed: ..." while leaving the rest
// for errors.Is.
func (e *BestEffortError) Get(i int) adt.Option[error] {
	if i >= 0 && i < len(e.Errors) {
		return adt.Of(e.Errors[i])
	}
	return adt.Empty[error]()
}
