package concurrency

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// errSentinel is a stable test-only error for comparison.
var errSentinel = errors.New("sentinel")

func TestGroup_NoTasksWaitReturnsNil(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	if err := g.Wait(); err != nil {
		t.Fatalf("Wait on empty group = %v, want nil", err)
	}
}

func TestGroup_AllSucceed(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	for i := 0; i < 5; i++ {
		g.Go(func(ctx context.Context) error { return nil })
	}
	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
}

// --- Strict mode ---------------------------------------------------------

func TestGroup_Strict_FirstErrorReturnedAndCancels(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	var secondStarted, secondDone atomic.Bool

	g.Go(func(ctx context.Context) error {
		// First task fails immediately.
		return errSentinel
	})

	g.Go(func(ctx context.Context) error {
		secondStarted.Store(true)
		<-ctx.Done() // wait for cancellation
		secondDone.Store(true)
		return ctx.Err()
	})

	err := g.Wait()
	if !errors.Is(err, errSentinel) {
		t.Fatalf("Wait err = %v, want sentinel", err)
	}
	if !secondStarted.Load() {
		t.Fatalf("second task never started")
	}
	if !secondDone.Load() {
		t.Fatalf("second task did not observe cancellation")
	}
}

func TestGroup_Strict_OnlyFirstErrorReturned(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	g.Go(func(ctx context.Context) error { return errors.New("first") })
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return errors.New("second")
	})
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return errors.New("third")
	})

	err := g.Wait()
	if err == nil || err.Error() != "first" {
		t.Fatalf("Wait err = %v, want \"first\"", err)
	}
}

func TestGroup_Strict_GoAfterErrorStillRuns(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	g.Go(func(ctx context.Context) error { return errSentinel })

	// Calling Go after the first failure should still spawn the
	// goroutine; it just shouldn't change the Wait adt.
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})

	if err := g.Wait(); !errors.Is(err, errSentinel) {
		t.Fatalf("Wait err = %v, want sentinel", err)
	}
}

func TestGroup_Strict_ParentCtxCancelReturnsCtxErr(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	g := NewGroup(ctx)

	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	cancel()

	err := g.Wait()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait err = %v, want context.Canceled", err)
	}
}

// --- BestEffort mode ----------------------------------------------------

func TestGroup_BestEffort_AllSucceedReturnsNil(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))
	for i := 0; i < 5; i++ {
		g.Go(func(ctx context.Context) error { return nil })
	}
	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
}

func TestGroup_BestEffort_OneSuccessStillReturnsAggregatedError(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))

	g.Go(func(ctx context.Context) error { return errors.New("a failed") })
	g.Go(func(ctx context.Context) error { return nil })
	g.Go(func(ctx context.Context) error { return errors.New("b failed") })

	err := g.Wait()
	if err == nil {
		t.Fatal("Wait = nil, want *BestEffortError (any failure aggregates)")
	}
	var be *BestEffortError
	if !errors.As(err, &be) {
		t.Fatalf("Wait err = %v, want *BestEffortError", err)
	}
	if len(be.Errors) != 2 {
		t.Fatalf("BestEffortError.Errors has %d entries, want 2", len(be.Errors))
	}
}

func TestGroup_BestEffort_AllFailedReturnsAggregatedError(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))

	g.Go(func(ctx context.Context) error { return errors.New("a") })
	g.Go(func(ctx context.Context) error { return errors.New("b") })
	g.Go(func(ctx context.Context) error { return errors.New("c") })

	err := g.Wait()
	if err == nil {
		t.Fatal("Wait = nil, want aggregated error")
	}

	var be *BestEffortError
	if !errors.As(err, &be) {
		t.Fatalf("Wait err = %v, want *BestEffortError", err)
	}
	if len(be.Errors) != 3 {
		t.Fatalf("BestEffortError.Errors has %d entries, want 3", len(be.Errors))
	}

	// Each underlying error should be reachable via errors.Is.
	for _, want := range []string{"a", "b", "c"} {
		found := false
		for _, e := range be.Errors {
			if strings.Contains(e.Error(), want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BestEffortError missing %q in %v", want, be.Errors)
		}
	}
}

func TestGroup_BestEffort_SingleFailureReturnsErrorDirectly(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))
	g.Go(func(ctx context.Context) error { return errSentinel })

	err := g.Wait()
	if !errors.Is(err, errSentinel) {
		t.Fatalf("Wait err = %v, want sentinel (single error passes through)", err)
	}
	if _, ok := err.(*BestEffortError); ok {
		t.Fatalf("single-error Wait returned *BestEffortError, want bare error")
	}
}

func TestGroup_BestEffort_FailureDoesNotCancelSiblings(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))

	var secondRan, thirdRan atomic.Bool

	g.Go(func(ctx context.Context) error {
		// Fail immediately. Siblings should still get to run.
		return errors.New("first failed")
	})

	g.Go(func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond) // ensure ordering
		secondRan.Store(true)
		return nil
	})

	g.Go(func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		thirdRan.Store(true)
		return errors.New("third failed")
	})

	// Two failures → *BestEffortError; we only assert siblings ran.
	if err := g.Wait(); err == nil {
		t.Fatalf("Wait = nil, want *BestEffortError (two tasks failed)")
	}
	if !secondRan.Load() {
		t.Fatalf("second task did not run (was canceled by sibling error)")
	}
	if !thirdRan.Load() {
		t.Fatalf("third task did not run (was canceled by sibling error)")
	}
}

func TestGroup_BestEffort_ParentCtxCancelCollectsCtxErrs(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	g := NewGroup(ctx, WithMode(BestEffort))

	for i := 0; i < 3; i++ {
		g.Go(func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		})
	}

	cancel()

	err := g.Wait()
	// Every task returned ctx.Err() → *BestEffortError wrapping them.
	// errors.Is must still match via the standard Unwrap() []error
	// contract.
	var be *BestEffortError
	if !errors.As(err, &be) {
		t.Fatalf("Wait err = %v, want *BestEffortError", err)
	}
	if len(be.Errors) != 3 {
		t.Fatalf("BestEffortError.Errors has %d entries, want 3", len(be.Errors))
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("errors.Is(err, context.Canceled) = false, want true (via Unwrap)")
	}
}

// --- WithLimit -----------------------------------------------------------

func TestGroup_WithLimit_BlocksExcessGo(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithLimit(2))

	var (
		running    atomic.Int32
		maxRunning atomic.Int32
	)

	observe := func() func() {
		n := running.Add(1)
		for {
			cur := maxRunning.Load()
			if n <= cur || maxRunning.CompareAndSwap(cur, n) {
				break
			}
		}
		return func() { running.Add(-1) }
	}

	for i := 0; i < 10; i++ {
		g.Go(func(ctx context.Context) error {
			defer observe()()
			time.Sleep(20 * time.Millisecond)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if m := maxRunning.Load(); m > 2 {
		t.Errorf("max concurrent running tasks = %d, want <= 2", m)
	}
}

func TestGroup_WithLimit_OneTaskRuns(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithLimit(1))

	var ran atomic.Bool
	g.Go(func(ctx context.Context) error {
		ran.Store(true)
		return nil
	})

	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if !ran.Load() {
		t.Fatal("task never ran")
	}
}

// --- WaitWithContext -----------------------------------------------------

func TestGroup_WaitWithContext_NoTasksReturnsNil(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	if err := g.WaitWithContext(context.Background()); err != nil {
		t.Fatalf("WaitWithContext on empty group = %v, want nil", err)
	}
}

func TestGroup_WaitWithContext_AllSucceedReturnsNil(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	for i := 0; i < 5; i++ {
		g.Go(func(ctx context.Context) error { return nil })
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := g.WaitWithContext(ctx); err != nil {
		t.Fatalf("WaitWithContext = %v, want nil", err)
	}
}

func TestGroup_WaitWithContext_Strict_TaskErrorBeforeCtx(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return errSentinel })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := g.WaitWithContext(ctx); !errors.Is(err, errSentinel) {
		t.Fatalf("WaitWithContext err = %v, want sentinel", err)
	}
}

func TestGroup_WaitWithContext_BestEffort_TaskErrorBeforeCtx(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))
	g.Go(func(ctx context.Context) error { return errSentinel })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := g.WaitWithContext(ctx); !errors.Is(err, errSentinel) {
		t.Fatalf("WaitWithContext err = %v, want sentinel", err)
	}
}

// makeBlockingTask spawns a single task that blocks on `release` and
// returns nil, and returns a channel that is closed when the task has
// entered the blocking receive. Tests use this to make "task still
// running at WaitWithContext entry" deterministic.
func makeBlockingTask(t *testing.T, g *Group, release chan struct{}) <-chan struct{} {
	t.Helper()
	started := make(chan struct{})
	g.Go(func(ctx context.Context) error {
		close(started)
		<-release
		return nil
	})
	return started
}

func TestGroup_WaitWithContext_Strict_CtxExpiresFirstReturnsCtxErr(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	g := NewGroup(context.Background())
	started := makeBlockingTask(t, g, release)
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := g.WaitWithContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitWithContext err = %v, want context.DeadlineExceeded", err)
	}

	close(release) // let the orphan task complete so test cleanup is clean
}

func TestGroup_WaitWithContext_BestEffort_CtxExpiresFirstReturnsCtxErr(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	g := NewGroup(context.Background(), WithMode(BestEffort))
	started := makeBlockingTask(t, g, release)
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := g.WaitWithContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitWithContext err = %v, want context.DeadlineExceeded", err)
	}

	close(release)
}

func TestGroup_WaitWithContext_CtxAlreadyDone(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	g := NewGroup(context.Background())
	started := makeBlockingTask(t, g, release)
	<-started

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already done before the call

	err := g.WaitWithContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("WaitWithContext err = %v, want context.Canceled", err)
	}

	close(release)
}

func TestGroup_WaitWithContext_NilCtxTreatedAsBackground(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return nil })
	if err := g.WaitWithContext(nil); err != nil {
		t.Fatalf("WaitWithContext(nil) = %v, want nil", err)
	}
}

// TestGroup_WaitWithContext_AfterTimeoutWaitReturnsRealError confirms
// the "abandon, don't kill" policy: when WaitWithContext times out,
// the goroutines keep running; calling Wait after they finish returns
// the real outcome (not ctx.Err()).
func TestGroup_WaitWithContext_AfterTimeoutWaitReturnsRealError(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	var finished atomic.Bool
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error {
		defer finished.Store(true)
		<-release
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	if err := g.WaitWithContext(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitWithContext err = %v, want context.DeadlineExceeded", err)
	}
	if finished.Load() {
		t.Fatalf("task finished before WaitWithContext returned (it should still be blocked)")
	}

	// Release the task; Wait should now return the real outcome.
	close(release)

	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if !finished.Load() {
		t.Fatalf("task never finished after release")
	}
}

// TestGroup_WaitWithContext_TaskFinishWinsOverExpiredCtx confirms the
// "ctx only governs the wait" semantics: when all tasks have already
// finished by the time WaitWithContext is called with an expired ctx,
// the outcome-based error is returned, not ctx.Err().
func TestGroup_WaitWithContext_TaskFinishWinsOverExpiredCtx(t *testing.T) {
	t.Parallel()

	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return errSentinel })

	// Let the task finish before we even reach WaitWithContext.
	if err := g.Wait(); err != errSentinel {
		t.Fatalf("priming Wait err = %v, want sentinel", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already done

	if err := g.WaitWithContext(ctx); !errors.Is(err, errSentinel) {
		t.Fatalf("WaitWithContext err = %v, want sentinel (tasks done before ctx expired)", err)
	}
}

// TestGroup_WaitWithContext_ConcurrentWithWait exercises Wait and
// WaitWithContext called concurrently against the same Group. Both
// must observe the same outcome; the race detector will catch any
// stale-channel or unsynchronised-state-read issues.
func TestGroup_WaitWithContext_ConcurrentWithWait(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	var waitErr, waitCtxErr error
	go func() {
		defer wg.Done()
		waitErr = g.Wait()
	}()
	go func() {
		defer wg.Done()
		waitCtxErr = g.WaitWithContext(ctx)
	}()
	wg.Wait()

	if waitErr != nil {
		t.Errorf("Wait err = %v, want nil", waitErr)
	}
	if waitCtxErr != nil {
		t.Errorf("WaitWithContext err = %v, want nil", waitCtxErr)
	}
}

// TestGroup_WaitWithContext_Stress spawns many concurrent
// WaitWithContext callers against a single task; every caller must
// observe the same outcome and the race detector must stay clean.
func TestGroup_WaitWithContext_Stress(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	const N = 100
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := g.WaitWithContext(ctx); err != nil {
				t.Errorf("WaitWithContext err = %v, want nil", err)
			}
		}()
	}

	g.Go(func(ctx context.Context) error { return nil })
	wg.Wait()
}

// --- Cleanup: Wait / WaitWithContext must cancel the Group ctx ----------
// to release any propagateCancel watcher the parent might have spawned
// and to free the cancelCtx from the parent's children map. Without
// this, BestEffort Groups (and Strict Groups whose tasks all succeed)
// leak resources tied to the cancelCtx.

// TestGroup_Wait_CancelsGroupContext: the cancel func is fired on the
// Wait return path, even when no task errored.
func TestGroup_Wait_CancelsGroupContext(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return nil })
	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if err := g.ctx.Err(); err == nil {
		t.Fatalf("g.ctx not canceled after Wait; err = nil")
	}
}

// TestGroup_WaitWithContext_CancelsGroupContextOnCompletion: even
// when the Wait returns because tasks finished (not because ctx
// fired), the Group ctx is still canceled.
func TestGroup_WaitWithContext_CancelsGroupContextOnCompletion(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return nil })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := g.WaitWithContext(ctx); err != nil {
		t.Fatalf("WaitWithContext = %v, want nil", err)
	}
	if err := g.ctx.Err(); err == nil {
		t.Fatalf("g.ctx not canceled after WaitWithContext; err = nil")
	}
}

// TestGroup_WaitWithContext_CancelsGroupContextOnTimeout: even when
// the Wait returns because the caller's ctx fired, the Group ctx is
// still canceled — so the propagateCancel watcher (if any) and the
// parent's children map both get cleaned up.
func TestGroup_WaitWithContext_CancelsGroupContextOnTimeout(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	g := NewGroup(context.Background())
	started := makeBlockingTask(t, g, release)
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	if err := g.WaitWithContext(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitWithContext err = %v, want context.DeadlineExceeded", err)
	}
	if err := g.ctx.Err(); err == nil {
		t.Fatalf("g.ctx not canceled after WaitWithContext timeout; err = nil")
	}

	close(release) // let the orphan task complete so test cleanup is clean
}

// --- Stress / concurrent safety ------------------------------------------

// TestGroup_Stress_ConcurrentGoAndWait is a stress test that exercises
// many goroutines calling Go concurrently while another goroutine
// repeatedly polls Wait. Must be race-free under `go test -race`.
func TestGroup_Stress_ConcurrentGoAndWait(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	const N = 1000
	done := make(chan struct{})

	// Producer: spawn N tasks as fast as possible.
	go func() {
		for i := 0; i < N; i++ {
			g.Go(func(ctx context.Context) error { return nil })
		}
		close(done)
	}()

	// Consumer: poll Wait repeatedly. The first Wait that observes all
	// goroutines returned wins.
	<-done
	for {
		if err := g.Wait(); err == nil {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func TestGroup_Stress_StrictWithManyFailures(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	const N = 100
	for i := 0; i < N; i++ {
		g.Go(func(ctx context.Context) error { return errSentinel })
	}

	err := g.Wait()
	if !errors.Is(err, errSentinel) {
		t.Fatalf("Wait err = %v, want sentinel", err)
	}
}

func TestGroup_Stress_BestEffortMixedSuccessAndFailure(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background(), WithMode(BestEffort))

	const N = 100
	var successCount atomic.Int32
	for i := 0; i < N; i++ {
		i := i
		g.Go(func(ctx context.Context) error {
			if i%2 == 0 {
				successCount.Add(1)
				return nil
			}
			return errors.New("fail")
		})
	}

	err := g.Wait()
	if err == nil {
		t.Fatal("Wait = nil, want *BestEffortError (half tasks failed)")
	}
	var be *BestEffortError
	if !errors.As(err, &be) {
		t.Fatalf("Wait err = %v, want *BestEffortError", err)
	}
	if len(be.Errors) != N/2 {
		t.Fatalf("BestEffortError.Errors has %d entries, want %d", len(be.Errors), N/2)
	}
	if got := successCount.Load(); got != N/2 {
		t.Errorf("successCount = %d, want %d", got, N/2)
	}
}

// TestGroup_WaitIsIdempotent: after all tasks finish, subsequent Wait
// calls return the same result without blocking.
func TestGroup_WaitIsIdempotent(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())
	g.Go(func(ctx context.Context) error { return nil })

	for i := 0; i < 3; i++ {
		if err := g.Wait(); err != nil {
			t.Fatalf("Wait #%d = %v, want nil", i, err)
		}
	}
}

func TestGroup_ModeString(t *testing.T) {
	t.Parallel()
	if got := Strict.String(); got != "Strict" {
		t.Errorf("Strict.String() = %q, want \"Strict\"", got)
	}
	if got := BestEffort.String(); got != "BestEffort" {
		t.Errorf("BestEffort.String() = %q, want \"BestEffort\"", got)
	}
}

// TestGroup_BestEffortError_UnwrapMatches confirms that errors.Is can
// find an underlying error in a *BestEffortError via the standard
// errors.Unwrap contract.
func TestGroup_BestEffortError_UnwrapMatches(t *testing.T) {
	t.Parallel()
	e := &BestEffortError{Errors: []error{errSentinel, errors.New("other")}}
	if !errors.Is(e, errSentinel) {
		t.Fatalf("errors.Is(e, sentinel) = false, want true")
	}
}

// --- existing concurrent infrastructure smoke test ---------------------

// TestGroup_ConcurrentGoCallers confirms Go is safe to call from
// multiple goroutines concurrently.
func TestGroup_ConcurrentGoCallers(t *testing.T) {
	t.Parallel()
	g := NewGroup(context.Background())

	const callers = 8
	const perCaller = 100

	var wg sync.WaitGroup
	wg.Add(callers)
	for c := 0; c < callers; c++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perCaller; i++ {
				g.Go(func(ctx context.Context) error { return nil })
			}
		}()
	}
	wg.Wait()

	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
}