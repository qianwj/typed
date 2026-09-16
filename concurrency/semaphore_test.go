package concurrency

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// --- Construction --------------------------------------------------------

func TestSemaphore_NewPanicsOnNonPositive(t *testing.T) {
	t.Parallel()
	for _, n := range []int{0, -1, -100, -1 << 30} {
		n := n
		t.Run("", func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("NewSemaphore(%d) did not panic", n)
				}
			}()
			NewSemaphore(n)
		})
	}
}

func TestSemaphore_FreshSemaphoreHasFullCapacity(t *testing.T) {
	t.Parallel()
	for _, n := range []int{1, 2, 8, 1024} {
		n := n
		t.Run("", func(t *testing.T) {
			t.Parallel()
			s := NewSemaphore(n)
			if got := s.Available(); got != n {
				t.Errorf("Available on fresh Semaphore(%d) = %d, want %d", n, got, n)
			}
		})
	}
}

// --- Acquire / Release ---------------------------------------------------

func TestSemaphore_AcquireRelease(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(2)

	s.Acquire()
	if got := s.Available(); got != 1 {
		t.Errorf("Available after one Acquire = %d, want 1", got)
	}

	s.Acquire()
	if got := s.Available(); got != 0 {
		t.Errorf("Available after two Acquires = %d, want 0", got)
	}

	s.Release()
	if got := s.Available(); got != 1 {
		t.Errorf("Available after one Release = %d, want 1", got)
	}

	s.Release()
	if got := s.Available(); got != 2 {
		t.Errorf("Available after two Releases = %d, want 2", got)
	}
}

func TestSemaphore_AcquireBlocksUntilRelease(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	s.Acquire() // full

	acquired := make(chan struct{})
	go func() {
		s.Acquire() // should block until Release below
		close(acquired)
	}()

	// Give the goroutine a chance to enter the blocked receive.
	time.Sleep(20 * time.Millisecond)

	select {
	case <-acquired:
		t.Fatal("second Acquire returned before Release")
	default:
	}

	s.Release()

	select {
	case <-acquired:
	case <-time.After(time.Second):
		t.Fatal("second Acquire did not unblock after Release")
	}
	s.Release()
}

func TestSemaphore_OverReleaseBlocks(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	// Semaphore starts at full capacity. An extra Release has no
	// matching Acquire; per the documented contract it blocks
	// indefinitely (the underlying channel send cannot complete
	// against a full buffered channel). We verify the block by
	// running the call in a goroutine and asserting that it does
	// not return within a short window.
	done := make(chan struct{})
	go func() {
		s.Release() // no matching Acquire — blocks
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Release returned despite no matching Acquire")
	case <-time.After(50 * time.Millisecond):
		// Expected: the goroutine is still blocked. Test exits; the
		// blocked goroutine is reclaimed when the test binary
		// exits.
	}
}

// --- TryAcquire ----------------------------------------------------------

func TestSemaphore_TryAcquire(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)

	if !s.TryAcquire() {
		t.Fatal("first TryAcquire = false, want true")
	}
	if s.TryAcquire() {
		t.Fatal("second TryAcquire = true, want false (semaphore full)")
	}
	s.Release()
	if !s.TryAcquire() {
		t.Fatal("TryAcquire after Release = false, want true")
	}
	s.Release()
}

// --- AcquireWithContext --------------------------------------------------

func TestSemaphore_AcquireWithContext_Success(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.AcquireWithContext(ctx); err != nil {
		t.Fatalf("AcquireWithContext err = %v, want nil", err)
	}
	if got := s.Available(); got != 0 {
		t.Errorf("Available after AcquireWithContext = %d, want 0", got)
	}
	s.Release()
}

func TestSemaphore_AcquireWithContext_CtxExpiresFirst(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	s.Acquire() // full

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := s.AcquireWithContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("AcquireWithContext err = %v, want context.DeadlineExceeded", err)
	}
	// Critical: ctx firing must NOT have taken a slot. If it had, the
	// semaphore would now have Available() == -1 from the user's
	// perspective, and the next Release would push Available() to 0,
	// silently under-counting the slots.
	if got := s.Available(); got != 0 {
		t.Errorf("Available after failed Acquire = %d, want 0 (slot must not be taken on ctx fire)", got)
	}
	s.Release()
}

func TestSemaphore_AcquireWithContext_CtxAlreadyCanceled(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	s.Acquire() // full

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.AcquireWithContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("AcquireWithContext err = %v, want context.Canceled", err)
	}
	if got := s.Available(); got != 0 {
		t.Errorf("Available after failed Acquire = %d, want 0", got)
	}
	s.Release()
}

func TestSemaphore_AcquireWithContext_NilCtx(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	if err := s.AcquireWithContext(nil); err != nil {
		t.Fatalf("AcquireWithContext(nil) err = %v, want nil", err)
	}
	if got := s.Available(); got != 0 {
		t.Errorf("Available after AcquireWithContext(nil) = %d, want 0", got)
	}
	s.Release()
}

func TestSemaphore_AcquireWithContext_RetryAfterCtxFire(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(1)
	s.Acquire() // full

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := s.AcquireWithContext(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("first AcquireWithContext err = %v, want DeadlineExceeded", err)
	}

	// Now release the slot and retry with a fresh, longer ctx. The
	// first failure must not have leaked a slot; the second call
	// should succeed promptly.
	s.Release()
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	if err := s.AcquireWithContext(ctx2); err != nil {
		t.Fatalf("second AcquireWithContext err = %v, want nil", err)
	}
	s.Release()
}

// --- Concurrency / stress ------------------------------------------------

// TestSemaphore_AvailableUnderConcurrentReads is a smoke test that
// Available does not deadlock or panic under concurrent calls. The
// window between a length read and a subsequent op is not
// synchronised — that's the documented contract — but the read itself
// must be safe.
func TestSemaphore_AvailableUnderConcurrentReads(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(10)

	const N = 100
	var wg sync.WaitGroup
	wg.Add(N)
	for range N {
		go func() {
			defer wg.Done()
			_ = s.Available()
		}()
	}
	wg.Wait()
}

func TestSemaphore_ConcurrentAcquireRelease(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(5)

	const N = 200
	var (
		wg             sync.WaitGroup
		maxConcurrent  atomic.Int32
		currentRunning atomic.Int32
	)

	for range N {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire()
			defer s.Release()

			cur := currentRunning.Add(1)
			// Track the high-water mark of concurrent holders.
			for {
				m := maxConcurrent.Load()
				if cur <= m || maxConcurrent.CompareAndSwap(m, cur) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			currentRunning.Add(-1)
		}()
	}
	wg.Wait()

	if m := maxConcurrent.Load(); m > 5 {
		t.Errorf("max concurrent holders = %d, want <= 5", m)
	}
	if m := maxConcurrent.Load(); m < 1 {
		t.Errorf("max concurrent holders = %d, want >= 1 (test should actually contend)", m)
	}
	if got := s.Available(); got != 5 {
		t.Errorf("Available after all Releases = %d, want 5", got)
	}
}

func TestSemaphore_AcquireWithContext_ConcurrentUnblocksOnRelease(t *testing.T) {
	t.Parallel()
	s := NewSemaphore(2)

	// Spawn N goroutines that all try to acquire. With only 2 slots,
	// the first 2 wake immediately and the remaining N-2 block
	// inside the channel-receive arm of the select. None of them
	// release after acquiring, so the final Available() is 0
	// (every slot acquired is still held).
	const N = 10
	errs := make([]error, N)
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		i := i
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			errs[i] = s.AcquireWithContext(ctx)
		}()
	}

	// Release N-2 more slots, one at a time with a small gap so the
	// runtime can wake one goroutine per send. Each Release is paired
	// with exactly one wake-up; no goroutine should miss its turn.
	for range N - 2 {
		s.Release()
		time.Sleep(2 * time.Millisecond)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("errs[%d] = %v, want nil", i, err)
		}
	}
	// All 10 acquisitions are still held by the goroutines (none of
	// them release), so Available is 0 — the semaphore is drained.
	if got := s.Available(); got != 0 {
		t.Errorf("Available after all goroutines done = %d, want 0 (all slots still held)", got)
	}
}
