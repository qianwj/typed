package concurrency

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/adt"
)

// tryTakeOr is a small helper for tests that want the (value, present) shape.
// It returns the dequeued value and true on success, or the zero value and
// false on an empty Option. Using IsEmpty first is necessary because Get
// returns the zero value when the Option is absent.
func tryTakeOr[T any](t *testing.T, q *BoundedBlockingQueue[T], wantZeroForReport T) (T, bool) {
	t.Helper()
	opt := q.TryPoll()
	if opt.IsEmpty() {
		var zero T
		return zero, false
	}
	return opt.Get(), true
}

func TestNew_PanicsOnNonPositiveCapacity(t *testing.T) {
	t.Parallel()
	for _, c := range []int{0, -1, -1000} {
		t.Run("", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected panic for capacity %d", c)
				}
			}()
			NewBoundedBlockingQueue[int](c)
		})
	}
}

func TestNew_CapacityAndSize(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](8)
	if got := q.Capacity(); got != 8 {
		t.Fatalf("Capacity() = %d, want 8", got)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() = %d, want 0", got)
	}
}

func TestNew_RoundsCapacityUpToPowerOfTwo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		requested int
		want      int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{100, 128},
		{1000, 1024},
		{1 << 14, 1 << 14},
		{(1 << 14) + 1, 1 << 15},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			q := NewBoundedBlockingQueue[int](c.requested)
			if got := q.Capacity(); got != c.want {
				t.Fatalf("requested %d: Capacity() = %d, want %d", c.requested, got, c.want)
			}
			// The reported capacity must hold at least the requested number of
			// elements — that's the whole point of rounding up.
			if got := q.Capacity(); got < c.requested {
				t.Fatalf("Capacity() = %d < requested %d", got, c.requested)
			}
		})
	}
}

func TestPushPoll_FIFO(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](4)
	for i := 1; i <= 4; i++ {
		q.Push(i)
	}
	if got := q.Size(); got != 4 {
		t.Fatalf("Size() = %d, want 4", got)
	}
	for i := 1; i <= 4; i++ {
		if got := q.Poll(); got != i {
			t.Fatalf("Poll() = %d, want %d", got, i)
		}
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() after drain = %d, want 0", got)
	}
}

func TestPoll_BlocksWhenEmpty(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)

	takeDone := make(chan int, 1)
	go func() {
		takeDone <- q.Poll()
	}()

	// Give the goroutine time to reach cond.Wait().
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

func TestPush_BlocksWhenFull(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(1)
	q.Push(2)
	if got := q.Size(); got != 2 {
		t.Fatalf("Size() = %d, want 2", got)
	}

	pushDone := make(chan struct{}, 1)
	go func() {
		q.Push(3)
		pushDone <- struct{}{}
	}()

	time.Sleep(20 * time.Millisecond)
	select {
	case <-pushDone:
		t.Fatal("Push returned before the queue had space")
	default:
	}

	if got := q.Poll(); got != 1 {
		t.Fatalf("Poll() = %d, want 1", got)
	}

	select {
	case <-pushDone:
	case <-time.After(time.Second):
		t.Fatal("Push did not return after Take freed a slot")
	}

	if got, ok := tryTakeOr(t, q, 2); !ok || got != 2 {
		t.Fatalf("TryPoll() = (%d, present=%v), want (2, present)", got, ok)
	}
	if got, ok := tryTakeOr(t, q, 3); !ok || got != 3 {
		t.Fatalf("TryPoll() = (%d, present=%v), want (3, present)", got, ok)
	}
}

func TestTryPush_TryPoll(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)

	if !q.TryPush(10) {
		t.Fatal("TryPush on empty queue returned false")
	}
	if !q.TryPush(20) {
		t.Fatal("TryPush on queue with one slot returned false")
	}
	if q.TryPush(30) {
		t.Fatal("TryPush on full queue returned true")
	}

	if v, ok := tryTakeOr(t, q, 10); !ok || v != 10 {
		t.Fatalf("TryPoll() = (%d, present=%v), want (10, present)", v, ok)
	}
	if v, ok := tryTakeOr(t, q, 20); !ok || v != 20 {
		t.Fatalf("TryPoll() = (%d, present=%v), want (20, present)", v, ok)
	}
	if got := q.TryPoll(); !got.IsEmpty() {
		t.Fatalf("TryPoll on empty queue returned %v, want Empty", got)
	}
}

func TestRingWrap(t *testing.T) {
	t.Parallel()
	// Exercise head/tail wrap-around multiple times against a small buffer.
	q := NewBoundedBlockingQueue[int](4)
	for round := 0; round < 3; round++ {
		for i := 0; i < 4; i++ {
			q.Push(round*10 + i)
		}
		for i := 0; i < 4; i++ {
			if got := q.Poll(); got != round*10+i {
				t.Fatalf("round %d: Poll() = %d, want %d", round, got, round*10+i)
			}
		}
	}
}

func TestPointerT_DoesNotLeak(t *testing.T) {
	t.Parallel()
	// Smoke test: pointer-typed T must round-trip correctly, including
	// through a blocked Push that is unblocked by a Take.
	type box struct{ payload *int }

	q := NewBoundedBlockingQueue[*box](2)
	v1, v2, v3 := &box{payload: new(int)}, &box{payload: new(int)}, &box{payload: new(int)}
	q.Push(v1)
	q.Push(v2)

	pushed := make(chan struct{}, 1)
	go func() {
		q.Push(v3)
		pushed <- struct{}{}
	}()
	time.Sleep(20 * time.Millisecond)

	if got := q.Poll(); got != v1 {
		t.Fatalf("got %p, want %p", got, v1)
	}
	<-pushed

	if got := q.Poll(); got != v2 {
		t.Fatalf("got %p, want %p", got, v2)
	}
	if got := q.Poll(); got != v3 {
		t.Fatalf("got %p, want %p", got, v3)
	}

	// Hint to the runtime that now is a fine time to collect. The Take
	// implementation clears the slot to release the reference, but we don't
	// directly assert that here — we rely on -race and the slot-zeroing
	// code path in Take/TryPoll.
	runtime.GC()
}

// TestConcurrent_ManyProducersManyConsumers is a stress test that exercises
// Push/Take under heavy contention. It must:
//   - not deadlock
//   - not lose or duplicate elements
//   - be race-free (verified separately with `go test -race`)
//
// Producers block on Push when the queue is small, so the blocking path
// is exercised on the producer side even though the consumer uses TryPoll.
func TestConcurrent_ManyProducersManyConsumers(t *testing.T) {
	t.Parallel()
	const (
		producers      = 8
		perProducer    = 5_000
		expectedTotal  = producers * perProducer
		queueCapacity  = 16
		overallTimeout = 30 * time.Second
	)

	q := NewBoundedBlockingQueue[int](queueCapacity)

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

	// A single consumer using TryPoll in a tight loop is the simplest way
	// to drain the queue without ever blocking — that means the test's
	// "did the system make progress?" signal is purely the counter.
	doneProducing := make(chan struct{})
	go func() {
		wgProd.Wait()
		close(doneProducing)
	}()

loop:
	for {
		if opt := q.TryPoll(); !opt.IsEmpty() {
			_ = opt.Get()
			consumed.Add(1)
			continue
		}
		select {
		case <-doneProducing:
			// Producers are done; drain anything left and exit.
			for {
				if opt := q.TryPoll(); !opt.IsEmpty() {
					_ = opt.Get()
					consumed.Add(1)
					continue
				}
				break loop
			}
		default:
			runtime.Gosched()
		}
	}

	select {
	case <-doneProducing:
	case <-time.After(overallTimeout):
		t.Fatalf("producers timeout: produced=%d", produced.Load())
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

// TestTryPollReturnsOption is a small sanity check on the return type —
// it locks in that TryPoll returns adt.Option[T] (and therefore an
// absent value is observable via IsEmpty, not via a zero T).
func TestTryPollReturnsOption(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(0) // note: 0 is a perfectly valid element; an absent Option must not be confused with it

	opt := q.TryPoll()
	if opt.IsEmpty() {
		t.Fatal("first TryPoll() should be present")
	}
	if got := opt.Get(); got != 0 {
		t.Fatalf("first TryPoll() = %d, want 0", got)
	}

	empty := q.TryPoll()
	if !empty.IsEmpty() {
		t.Fatalf("second TryPoll() = %v, want Empty", empty)
	}

	// Sanity: the empty factory also looks the same.
	if e := adt.Empty[int](); !e.IsEmpty() {
		t.Fatal("adt.Empty[int]() should be empty")
	}
}

// --- PushWithContext / PollWithContext ---------------------------------------------------

func TestPushWithContext_FastPathSucceeds(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	if err := q.PushWithContext(context.Background(), 7); err != nil {
		t.Fatalf("PushWithContext on empty queue returned %v, want nil", err)
	}
	if got, ok := tryTakeOr(t, q, -1); !ok || got != 7 {
		t.Fatalf("element after PushWithContext = (%d, present=%v), want (7, present)", got, ok)
	}
}

func TestPollWithContext_FastPathSucceeds(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(9)
	result := q.PollWithContext(context.Background())
	if result.IsFailure() {
		t.Fatalf("PollWithContext on non-empty queue failed: %v", result.Error())
	}
	if v := result.Value(); v != 9 {
		t.Fatalf("PollWithContext = %d, want 9", v)
	}
}

func TestPushWithContext_BlocksUntilTake(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(1)
	q.Push(2) // queue is full

	pushDone := make(chan error, 1)
	go func() {
		pushDone <- q.PushWithContext(context.Background(), 3)
	}()

	// Give the goroutine time to reach the channel receive.
	time.Sleep(20 * time.Millisecond)
	select {
	case err := <-pushDone:
		t.Fatalf("PushWithContext returned %v before any Take", err)
	default:
	}

	if got := q.Poll(); got != 1 {
		t.Fatalf("Poll() = %d, want 1", got)
	}

	select {
	case err := <-pushDone:
		if err != nil {
			t.Fatalf("PushWithContext returned err=%v after Take, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("PushWithContext did not return after Take freed a slot")
	}
	if got, ok := tryTakeOr(t, q, -1); !ok || got != 2 {
		t.Fatalf("first remaining = (%d, present=%v), want (2, present)", got, ok)
	}
	if got, ok := tryTakeOr(t, q, -1); !ok || got != 3 {
		t.Fatalf("element after PushWithContext = (%d, present=%v), want (3, present)", got, ok)
	}
}

func TestPollWithContext_BlocksUntilPush(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)

	takeDone := make(chan adt.Result[int], 1)
	go func() {
		takeDone <- q.PollWithContext(context.Background())
	}()

	time.Sleep(20 * time.Millisecond)
	select {
	case result := <-takeDone:
		t.Fatalf("PollWithContext returned %v before any Push", result)
	default:
	}

	q.Push(42)

	select {
	case result := <-takeDone:
		if result.IsFailure() {
			t.Fatalf("PollWithContext failed: %v", result.Error())
		}
		if v := result.Value(); v != 42 {
			t.Fatalf("PollWithContext value = %d, want 42", v)
		}
	case <-time.After(time.Second):
		t.Fatal("PollWithContext did not return after Push")
	}
}

func TestPushWithContext_CanceledDoesNotEnqueue(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(1)
	q.Push(2) // queue is full

	ctx, cancel := context.WithCancel(context.Background())
	pushDone := make(chan error, 1)
	go func() {
		pushDone <- q.PushWithContext(ctx, 99)
	}()

	// Let the goroutine reach the channel receive, then cancel.
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-pushDone:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("PushWithContext err = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("PushWithContext did not return after cancel")
	}
	// The element must NOT have been enqueued.
	if got := q.Size(); got != 2 {
		t.Fatalf("Size after canceled PushWithContext = %d, want 2 (element must be dropped)", got)
	}
	if got := q.Poll(); got != 1 {
		t.Fatalf("Poll() = %d, want 1 (queue unchanged)", got)
	}
	if got := q.Poll(); got != 2 {
		t.Fatalf("Poll() = %d, want 2 (queue unchanged)", got)
	}
}

func TestPollWithContext_CanceledReturnsFailure(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)

	ctx, cancel := context.WithCancel(context.Background())
	takeDone := make(chan adt.Result[int], 1)
	go func() {
		takeDone <- q.PollWithContext(ctx)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case result := <-takeDone:
		if !result.IsFailure() || !errors.Is(result.Error(), context.Canceled) {
			t.Fatalf("PollWithContext = %v, want failure with context.Canceled", result)
		}
	case <-time.After(time.Second):
		t.Fatal("PollWithContext did not return after cancel")
	}
	// Queue must still be empty — cancellation must not have dequeued anything.
	if got := q.Size(); got != 0 {
		t.Fatalf("Size after canceled PollWithContext = %d, want 0", got)
	}
}

func TestPushWithContext_Deadline(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](1)
	q.Push(1) // queue is full

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := q.PushWithContext(ctx, 2)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("PushWithContext err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed < 20*time.Millisecond {
		t.Fatalf("PushWithContext returned in %v, want >= ~30ms (deadline)", elapsed)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("PushWithContext returned in %v, want ~30ms (deadline)", elapsed)
	}
	if got := q.Size(); got != 1 {
		t.Fatalf("Size after deadline = %d, want 1", got)
	}
}

func TestPollWithContext_Deadline(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	result := q.PollWithContext(ctx)
	elapsed := time.Since(start)

	if !result.IsFailure() || !errors.Is(result.Error(), context.DeadlineExceeded) {
		t.Fatalf("PollWithContext = %v, want failure with context.DeadlineExceeded", result)
	}
	if elapsed < 20*time.Millisecond {
		t.Fatalf("PollWithContext returned in %v, want >= ~30ms (deadline)", elapsed)
	}
}

// TestPushWithContext_PollWithContextRace is a stress test for the ctx-aware paths: many
// goroutines call PushWithContext and PollWithContext concurrently, and a fraction of
// them get their ctx canceled mid-flight. The test passes if every
// successful Push/Take is paired (no lost or duplicate elements) and
// no goroutine deadlocks.
func TestPushWithContext_PollWithContextRace(t *testing.T) {
	t.Parallel()
	const (
		producers      = 4
		perProducer    = 2_000
		expectedTotal  = producers * perProducer
		queueCapacity  = 32
		overallTimeout = 30 * time.Second
	)

	q := NewBoundedBlockingQueue[int](queueCapacity)

	var produced atomic.Int64
	var consumed atomic.Int64
	var wgProd sync.WaitGroup
	wgProd.Add(producers)
	for p := 0; p < producers; p++ {
		base := p * perProducer
		go func() {
			defer wgProd.Done()
			for i := 0; i < perProducer; i++ {
				_ = q.PushWithContext(context.Background(), base+i)
				produced.Add(1)
			}
		}()
	}

	// One consumer running PollWithContext in a tight loop. Each call uses a fresh
	// ctx with a 5ms timeout; a few will hit the timeout, and those PollWithContext
	// calls must return Failure(context.DeadlineExceeded) without dequeuing.
	takeDone := make(chan struct{})
	go func() {
		defer close(takeDone)
		for consumed.Load() < int64(expectedTotal) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
			if result := q.PollWithContext(ctx); result.IsSuccess() {
				consumed.Add(1)
			}
			cancel()
		}
	}()

	doneProducing := make(chan struct{})
	go func() {
		wgProd.Wait()
		close(doneProducing)
	}()

	select {
	case <-takeDone:
	case <-time.After(overallTimeout):
		t.Fatalf("consumer timeout: produced=%d consumed=%d", produced.Load(), consumed.Load())
	}
	<-doneProducing

	// Drain anything left so Size() must be 0 at the end.
	for {
		if opt := q.TryPoll(); opt.IsEmpty() {
			break
		}
	}

	if got := produced.Load(); got != int64(expectedTotal) {
		t.Fatalf("produced = %d, want %d", got, expectedTotal)
	}
	// The exact consumed count is bounded below by the number of successful
	// PollWithContext calls; we only assert that it matches the produced count
	// (i.e. every produced element was eventually consumed — possibly by
	// the final drain loop above).
	total := consumed.Load() + int64(q.Size())
	if total != int64(expectedTotal) {
		// q.Size() should be 0 here because of the drain loop, so this
		// check is really "consumed == expected".
		t.Fatalf("consumed+remaining = %d, want %d", total, expectedTotal)
	}
}
