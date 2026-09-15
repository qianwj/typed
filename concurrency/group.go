// Package concurrency provides synchronization primitives that complement
// the standard library.
package concurrency

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// Mode selects the failure policy of a [Group].
type Mode int

const (
	// Strict mode (default) is the standard "errgroup" semantics: when any
	// task returns a non-nil error, the group's ctx is canceled, siblings
	// see the cancellation and abort, and [Group.Wait] returns that error.
	// Subsequent calls to [Group.Go] after the first error still spawn
	// goroutines but their results are ignored by [Group.Wait].
	Strict Mode = iota

	// BestEffort mode runs all tasks to completion regardless of
	// siblings' failures. The group's ctx is NOT canceled on a task
	// error, so all goroutines run to their natural end (or until the
	// parent ctx is canceled). [Group.Wait] returns nil if at least one
	// task succeeded, the parent ctx error if the parent was canceled
	// while no task had succeeded yet, or a [*BestEffortError] wrapping
	// every failed task otherwise.
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
// is freed. A non-positive n means no limit. [WithLimit] must be supplied
// at construction; [Group.SetLimit] covers the equivalent runtime API
// for callers that need to defer the decision.
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

	parentCtx context.Context
	ctx       context.Context // derived from parentCtx, canceled by cancel()
	cancel    context.CancelFunc

	mu    sync.Mutex
	state groupState

	// started becomes true the first time Go is called. It guards
	// SetLimit, which is only valid before any goroutine has launched
	// (matches golang.org/x/sync/errgroup semantics).
	started atomic.Bool
}

// groupState is the result-tracking half of Group. protected by mu.
//
// In Strict mode, only firstErr matters; errs is always nil.
// In BestEffort mode, firstErr is always nil and errs accumulates every
// failed task. successCount counts how many tasks returned nil.
type groupState struct {
	firstErr     error
	errs         []error
	successCount int
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
		cfg:       cfg,
		parentCtx: parent,
		ctx:       ctx,
		cancel:    cancel,
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
// In [BestEffort] mode, the ctx is NOT canceled. The error is recorded
// and [Group.Wait] will return a [*BestEffortError] only if ALL tasks
// failed.
//
// Go is safe to call from multiple goroutines concurrently. Calls to Go
// after the first error are still accepted (the goroutines run); they
// just have no effect on [Group.Wait]'s return value.
func (g *Group) Go(fn func(ctx context.Context) error) {
	g.started.Store(true)
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

	// In Strict mode, the first error wins. Subsequent errors don't
	// change firstErr (it's already the "first") and we don't bother
	// canceling again (it's a no-op). Skip the record to keep the
	// first error stable and avoid extra work.
	if g.cfg.mode == Strict && g.state.firstErr != nil {
		return
	}

	if err == nil {
		g.state.successCount++
		return
	}

	if g.cfg.mode == Strict {
		g.state.firstErr = err
		g.cancel()
		return
	}

	// BestEffort
	g.state.errs = append(g.state.errs, err)
}

// Wait blocks until every task spawned by [Group.Go] has returned, then
// returns the error that best describes the outcome:
//
//   - In [Strict] mode: the first task error, or nil if every task
//     succeeded (including the trivial case of zero tasks).
//   - In [BestEffort] mode: nil if at least one task succeeded; the
//     parent ctx's error if the parent was canceled before any task
//     succeeded; or a [*BestEffortError] wrapping every failed task
//     when all tasks failed and the parent ctx is still live.
//
// Wait is safe to call multiple times; subsequent calls return the
// same value once all tasks have finished.
func (g *Group) Wait() error {
	g.wg.Wait()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cfg.mode == Strict {
		if g.state.firstErr != nil {
			return g.state.firstErr
		}
		return nil
	}

	// BestEffort
	if g.state.successCount > 0 {
		return nil
	}
	if err := g.parentCtx.Err(); err != nil {
		return err
	}
	if len(g.state.errs) == 0 {
		// No tasks at all. Treat as success.
		return nil
	}
	if len(g.state.errs) == 1 {
		return g.state.errs[0]
	}
	return &BestEffortError{Errors: g.state.errs}
}

// SetLimit caps the number of tasks that may run concurrently. Calls to
// [Group.Go] beyond the cap block until a slot is freed. n <= 0 removes
// any previously-set limit.
//
// SetLimit must be called before any [Group.Go] (matches
// golang.org/x/sync/errgroup semantics); calling it after a task has been
// spawned panics, because the limit at construction time and the
// concurrency cap are not designed to be mutated once goroutines are in
// flight.
func (g *Group) SetLimit(n int) {
	if g.started.Load() {
		panic("concurrency: Group.SetLimit called after Go")
	}
	if n <= 0 {
		g.sem = nil
		return
	}
	g.sem = make(chan struct{}, n)
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
	fmt.Fprintf(&sb, "concurrency: %d tasks failed: ", len(e.Errors))
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