package concurrency

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// tryTakeUnboundedOr is a small helper for tests that want the
// (value, present) shape. Returns the dequeued value and true on
// success, or the zero value and false on an empty Optional.
func tryTakeUnboundedOr[T any](t *testing.T, q *UnboundedBlockingQueue[T], wantZeroForReport T) (T, bool) {
	t.Helper()
	opt := q.TryTake()
	if opt.IsEmpty() {
		var zero T
		return zero, false
	}
	return opt.Get(), true
}

func TestUnboundedNewAndSize(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	if got := q.Size(); got != 0 {
		t.Fatalf("fresh queue Size() = %d, want 0", got)
	}
}

func TestUnboundedPushTakeFIFO(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	for i := 1; i <= 5; i++ {
		if got := q.Take(); got != i {
			t.Fatalf("Take #%d = %d, want %d", i, got, i)
		}
	}
}

func TestUnboundedTryPushAlwaysSucceeds(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	for i := 0; i < 1000; i++ {
		if !q.TryPush(i) {
			t.Fatalf("TryPush returned false at i=%d (queue is unbounded)", i)
		}
	}
	if got := q.Size(); got != 1000 {
		t.Fatalf("Size() = %d, want 1000", got)
	}
}

func TestUnboundedTryTake(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	q.Push(7)
	if v, ok := tryTakeUnboundedOr(t, q, -1); !ok || v != 7 {
		t.Fatalf("TryTake() = (%d, present=%v), want (7, present)", v, ok)
	}
	if got := q.TryTake(); !got.IsEmpty() {
		t.Fatalf("TryTake on empty queue returned %v, want Empty", got)
	}
}

func TestUnboundedTakeBlocksWhenEmpty(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()

	takeDone := make(chan int, 1)
	go func() {
		takeDone <- q.Take()
	}()

	time.Sleep(20 * time.Millisecond)
	select {
	case v := <-takeDone:
		t.Fatalf("Take returned %d before any Push", v)
	default:
	}

	q.Push(42)

	select {
	case v := <-takeDone:
		if v != 42 {
			t.Fatalf("Take returned %d, want 42", v)
		}
	case <-time.After(time.Second):
		t.Fatal("Take did not return after Push")
	}
}

func TestUnboundedPushCtxAlwaysSucceedsWhenCtxAlive(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	for i := 0; i < 10; i++ {
		if err := q.PushCtx(context.Background(), i); err != nil {
			t.Fatalf("PushCtx returned %v, want nil", err)
		}
	}
	if got := q.Size(); got != 10 {
		t.Fatalf("Size() = %d, want 10", got)
	}
}

func TestUnboundedPushCtxRespectsCanceledCtx(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before push
	if err := q.PushCtx(ctx, 99); !errors.Is(err, context.Canceled) {
		t.Fatalf("PushCtx err = %v, want context.Canceled", err)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() after canceled PushCtx = %d, want 0", got)
	}
}

func TestUnboundedTakeCtxDeadline(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := q.TakeCtx(ctx)
	elapsed := time.Since(start)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("TakeCtx err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed < 20*time.Millisecond {
		t.Fatalf("TakeCtx returned in %v, want >= ~30ms", elapsed)
	}
}

func TestUnboundedTakeCtxDeadlineReturnsAfterItem(t *testing.T) {
	t.Parallel()
	q := NewUnboundedBlockingQueue[int]()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go func() {
		time.Sleep(10 * time.Millisecond)
		q.Push(1)
	}()
	v, err := q.TakeCtx(ctx)
	if err != nil {
		t.Fatalf("TakeCtx err = %v, want nil", err)
	}
	if v != 1 {
		t.Fatalf("TakeCtx value = %d, want 1", v)
	}
}

func TestUnboundedRingGrowth(t *testing.T) {
	t.Parallel()
	// Push far more than the initial capacity (16) to exercise the
	// grow path several times. With power-of-2 doubling, going from
	// 16 to 5_000 takes 9 grows.
	q := NewUnboundedBlockingQueue[int]()
	const n = 5000
	for i := 0; i < n; i++ {
		q.Push(i)
	}
	for i := 0; i < n; i++ {
		if got := q.Take(); got != i {
			t.Fatalf("Take #%d = %d, want %d (ring growth broken)", i, got, i)
		}
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size after drain = %d, want 0", got)
	}
}

func TestUnboundedPointerTGCReleased(t *testing.T) {
	t.Parallel()
	// Smoke test: pointer-typed T must round-trip correctly, including
	// through slot zeroing on Take. -race plus the slot-zeroing code
	// in Take is the real verification.
	type box struct{ payload *int }
	q := NewUnboundedBlockingQueue[*box]()
	const n = 100
	for i := 0; i < n; i++ {
		q.Push(&box{payload: new(int)})
	}
	for i := 0; i < n; i++ {
		v := q.Take()
		if v == nil {
			t.Fatalf("Take #%d returned nil", i)
		}
	}
	// After draining, the queue is empty; run a GC pass and make sure
	// nothing panics.
	runtime.GC()
}

// TestUnboundedConcurrent_ManyProducersManyConsumers is a stress test that
// exercises Push / Take under heavy contention. It must:
//   - not deadlock
//   - not lose or duplicate elements
//   - be race-free (verified separately with `go test -race`)
//
// The consumer uses TryTake in a tight loop (so it never blocks), which
// forces the producer side to exercise the blocking path when the
// queue is full — except this queue is unbounded, so producers never
// block. We still keep many producers to put pressure on the cond
// Signal path and the lock.
func TestUnboundedConcurrent_ManyProducersManyConsumers(t *testing.T) {
	t.Parallel()
	const (
		producers      = 8
		perProducer    = 5_000
		expectedTotal  = producers * perProducer
		overallTimeout = 30 * time.Second
	)

	q := NewUnboundedBlockingQueue[int]()

	var produced atomic.Int64
	var consumed atomic.Int64
	var wgProd sync.WaitGroup
	wgProd.Add(producers)
	for p := 0; p < producers; p++ {
		base := p * perProducer
		go func() {
			defer wgProd.Done()
			for i := 0; i < perProducer; i++ {
				q.Push(base + i)
				produced.Add(1)
			}
		}()
	}

	doneProducing := make(chan struct{})
	go func() {
		wgProd.Wait()
		close(doneProducing)
	}()

loop:
	for {
		if opt := q.TryTake(); !opt.IsEmpty() {
			_ = opt.Get()
			consumed.Add(1)
			continue
		}
		select {
		case <-doneProducing:
			for {
				opt := q.TryTake()
				if opt.IsEmpty() {
					break loop
				}
				_ = opt.Get()
				consumed.Add(1)
			}
		default:
			runtime.Gosched()
		}
	}

	select {
	case <-doneProducing:
	case <-time.After(overallTimeout):
		t.Fatalf("producer timeout: produced=%d", produced.Load())
	}

	if got := produced.Load(); got != int64(expectedTotal) {
		t.Fatalf("produced = %d, want %d", got, expectedTotal)
	}
	if got := consumed.Load(); got != int64(expectedTotal) {
		t.Fatalf("consumed = %d, want %d", got, expectedTotal)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() after drain = %d, want 0", got)
	}
}

// TestUnboundedBurstWakesAllWaiters is the regression test for the
// "lost wakeup" pitfall: if many pushes happen before any taker
// wakes, the queue must still hand out all items, not just the first
// one. With sync.Cond.Signal, each Push pops one waiter, so a burst of
// N pushes wakes N takers. (A naive cap-1 channel design would fail
// this test by leaving items stranded.)
func TestUnboundedBurstWakesAllWaiters(t *testing.T) {
	t.Parallel()
	const (
		waiters = 16
		bursts  = 64
	)
	q := NewUnboundedBlockingQueue[int]()

	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(waiters)
	for w := 0; w < waiters; w++ {
		go func() {
			defer wg.Done()
			// Each taker takes one item. Loop until we get our share.
			for i := 0; i < bursts/waiters; i++ {
				q.Take()
				received.Add(1)
			}
		}()
	}

	// Push bursts/waiters items per taker. Because pushes happen
	// before takers start draining, all waiters are already in the
	// cond's wait list when the first Signal fires.
	for i := 0; i < bursts; i++ {
		q.Push(i)
	}

	wg.Wait()

	if got := received.Load(); got != int64(bursts) {
		t.Fatalf("received = %d, want %d (lost wakeup?)", got, int64(bursts))
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() after drain = %d, want 0", got)
	}
}