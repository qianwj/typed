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