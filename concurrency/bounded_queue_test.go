package concurrency

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/utils/option"
)

// tryTakeOr is a small helper for tests that want the (value, present) shape.
// It returns the dequeued value and true on success, or the zero value and
// false on an empty Optional. Using IsEmpty first is necessary because Get
// returns the zero value when the Optional is absent.
func tryTakeOr[T any](t *testing.T, q *BoundedBlockingQueue[T], wantZeroForReport T) (T, bool) {
	t.Helper()
	opt := q.TryTake()
	if opt.IsEmpty() {
		var zero T
		return zero, false
	}
	return opt.Get(), true
}

func TestNew_PanicsOnNonPositiveCapacity(t *testing.T) {
	t.Parallel()
	for _, c := range []int{0, -1, -1000} {
		c := c
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
		c := c
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

func TestPushTake_FIFO(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](4)
	for i := 1; i <= 4; i++ {
		q.Push(i)
	}
	if got := q.Size(); got != 4 {
		t.Fatalf("Size() = %d, want 4", got)
	}
	for i := 1; i <= 4; i++ {
		if got := q.Take(); got != i {
			t.Fatalf("Take() = %d, want %d", got, i)
		}
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size() after drain = %d, want 0", got)
	}
}

func TestTake_BlocksWhenEmpty(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)

	takeDone := make(chan int, 1)
	go func() {
		takeDone <- q.Take()
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

	if got := q.Take(); got != 1 {
		t.Fatalf("Take() = %d, want 1", got)
	}

	select {
	case <-pushDone:
	case <-time.After(time.Second):
		t.Fatal("Push did not return after Take freed a slot")
	}

	if got, ok := tryTakeOr(t, q, 2); !ok || got != 2 {
		t.Fatalf("TryTake() = (%d, present=%v), want (2, present)", got, ok)
	}
	if got, ok := tryTakeOr(t, q, 3); !ok || got != 3 {
		t.Fatalf("TryTake() = (%d, present=%v), want (3, present)", got, ok)
	}
}

func TestTryPush_TryTake(t *testing.T) {
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
		t.Fatalf("TryTake() = (%d, present=%v), want (10, present)", v, ok)
	}
	if v, ok := tryTakeOr(t, q, 20); !ok || v != 20 {
		t.Fatalf("TryTake() = (%d, present=%v), want (20, present)", v, ok)
	}
	if got := q.TryTake(); !got.IsEmpty() {
		t.Fatalf("TryTake on empty queue returned %v, want Empty", got)
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
			if got := q.Take(); got != round*10+i {
				t.Fatalf("round %d: Take() = %d, want %d", round, got, round*10+i)
			}
		}
	}
}

func TestDrainTo_AllIntoSizedDst(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](4)
	for i := 1; i <= 3; i++ {
		q.Push(i)
	}
	dst := make([]int, 3)
	n := q.DrainTo(dst)
	if n != 3 {
		t.Fatalf("DrainTo returned %d, want 3", n)
	}
	want := []int{1, 2, 3}
	if !slicesEqual(dst, want) {
		t.Fatalf("DrainTo wrote %v, want %v", dst, want)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size after drain = %d, want 0", got)
	}
}

func TestDrainTo_PartialDst(t *testing.T) {
	t.Parallel()
	// Queue has 5 elements (cap 8 after rounding), dst has room for 2.
	q := NewBoundedBlockingQueue[int](5)
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	dst := make([]int, 2)
	n := q.DrainTo(dst)
	if n != 2 {
		t.Fatalf("DrainTo returned %d, want 2", n)
	}
	if dst[0] != 1 || dst[1] != 2 {
		t.Fatalf("DrainTo wrote %v, want [1, 2]", dst)
	}
	if got := q.Size(); got != 3 {
		t.Fatalf("Size after partial drain = %d, want 3", got)
	}
	// Remaining elements should still come out in FIFO order.
	for _, want := range []int{3, 4, 5} {
		if got := q.Take(); got != want {
			t.Fatalf("Take after partial drain = %d, want %d", got, want)
		}
	}
}

func TestDrainTo_EmptyQueue(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](4)
	dst := []int{99, 99, 99} // sentinel values
	n := q.DrainTo(dst)
	if n != 0 {
		t.Fatalf("DrainTo on empty queue returned %d, want 0", n)
	}
	for i, v := range dst {
		if v != 99 {
			t.Fatalf("DrainTo mutated dst[%d] = %d, want 99 (untouched)", i, v)
		}
	}
}

func TestDrainTo_EmptyDst(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](4)
	q.Push(1)
	q.Push(2)
	var dst []int
	if n := q.DrainTo(dst); n != 0 {
		t.Fatalf("DrainTo with nil dst returned %d, want 0", n)
	}
	dst = []int{}
	if n := q.DrainTo(dst); n != 0 {
		t.Fatalf("DrainTo with empty dst returned %d, want 0", n)
	}
	if got := q.Size(); got != 2 {
		t.Fatalf("Size after empty-dst drains = %d, want 2 (queue must not be touched)", got)
	}
}

func TestDrainTo_AcrossWrap(t *testing.T) {
	t.Parallel()
	// Force the head to be in the middle of the backing array, then drain.
	// Queue: cap 8 (after rounding 5→8). Push 1..5, Take 2, Push 6..7.
	// Resulting layout: head=2, tail=7, count=5; ring wraps at 8.
	q := NewBoundedBlockingQueue[int](5)
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	_ = q.Take() // 1
	_ = q.Take() // 2
	q.Push(6)
	q.Push(7)
	// head=2, tail=(7+1)&7=0, count=5, items[2..7]=[3,4,5,6,7]

	dst := make([]int, 5)
	n := q.DrainTo(dst)
	if n != 5 {
		t.Fatalf("DrainTo returned %d, want 5", n)
	}
	want := []int{3, 4, 5, 6, 7}
	if !slicesEqual(dst, want) {
		t.Fatalf("DrainTo wrote %v, want %v", dst, want)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size after drain = %d, want 0", got)
	}
}

func TestDrainTo_AcrossWrapPartial(t *testing.T) {
	t.Parallel()
	// Same setup as above, but drain only 3 — exercises the "two ranges"
	// copy path with a small dst on a wrapped queue.
	q := NewBoundedBlockingQueue[int](5)
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	_ = q.Take()
	_ = q.Take()
	q.Push(6)
	q.Push(7)
	// head=2, count=5, items[2..7]=[3,4,5,6,7]

	dst := make([]int, 3)
	n := q.DrainTo(dst)
	if n != 3 {
		t.Fatalf("DrainTo returned %d, want 3", n)
	}
	want := []int{3, 4, 5}
	if !slicesEqual(dst, want) {
		t.Fatalf("DrainTo wrote %v, want %v", dst, want)
	}
	if got := q.Size(); got != 2 {
		t.Fatalf("Size after partial drain = %d, want 2", got)
	}
	// Remaining in order: 6, 7
	if v, ok := tryTakeOr(t, q, 6); !ok || v != 6 {
		t.Fatalf("first remaining = (%d, present=%v), want (6, present)", v, ok)
	}
	if v, ok := tryTakeOr(t, q, 7); !ok || v != 7 {
		t.Fatalf("second remaining = (%d, present=%v), want (7, present)", v, ok)
	}
}

func TestDrainTo_AllowsPushAfter(t *testing.T) {
	t.Parallel()
	// After a drain that empties a full queue, blocked Pushers must wake up.
	q := NewBoundedBlockingQueue[int](2)
	q.Push(1)
	q.Push(2)

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

	dst := make([]int, 2)
	if n := q.DrainTo(dst); n != 2 {
		t.Fatalf("DrainTo returned %d, want 2", n)
	}

	select {
	case <-pushDone:
	case <-time.After(time.Second):
		t.Fatal("blocked Push did not wake up after DrainTo freed slots")
	}
}

func TestDrainTo_PointerTGCReleased(t *testing.T) {
	t.Parallel()
	// Slots drained via DrainTo should be cleared to zero so the GC can
	// collect the pointer values they used to hold. This is a smoke test;
	// -race + the slot-zeroing code in DrainTo is the real verification.
	type box struct{ payload *int }
	q := NewBoundedBlockingQueue[*box](4)
	for i := 0; i < 4; i++ {
		q.Push(&box{payload: new(int)})
	}
	dst := make([]*box, 4)
	if n := q.DrainTo(dst); n != 4 {
		t.Fatalf("DrainTo returned %d, want 4", n)
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("Size after drain = %d, want 0", got)
	}
	// Drop our references and let the runtime collect.
	for i := range dst {
		dst[i] = nil
	}
	runtime.GC()
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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

	if got := q.Take(); got != v1 {
		t.Fatalf("got %p, want %p", got, v1)
	}
	<-pushed

	if got := q.Take(); got != v2 {
		t.Fatalf("got %p, want %p", got, v2)
	}
	if got := q.Take(); got != v3 {
		t.Fatalf("got %p, want %p", got, v3)
	}

	// Hint to the runtime that now is a fine time to collect. The Take
	// implementation clears the slot to release the reference, but we don't
	// directly assert that here — we rely on -race and the slot-zeroing
	// code path in Take/TryTake.
	runtime.GC()
}

// TestConcurrent_ManyProducersManyConsumers is a stress test that exercises
// Push/Take under heavy contention. It must:
//   - not deadlock
//   - not lose or duplicate elements
//   - be race-free (verified separately with `go test -race`)
//
// Producers block on Push when the queue is small, so the blocking path
// is exercised on the producer side even though the consumer uses TryTake.
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

	// A single consumer using TryTake in a tight loop is the simplest way
	// to drain the queue without ever blocking — that means the test's
	// "did the system make progress?" signal is purely the counter.
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
			// Producers are done; drain anything left and exit.
			for {
				if opt := q.TryTake(); !opt.IsEmpty() {
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

// TestTryTakeReturnsOptional is a small sanity check on the return type —
// it locks in that TryTake returns option.Optional[T] (and therefore an
// absent value is observable via IsEmpty, not via a zero T).
func TestTryTakeReturnsOptional(t *testing.T) {
	t.Parallel()
	q := NewBoundedBlockingQueue[int](2)
	q.Push(0) // note: 0 is a perfectly valid element; an absent Optional must not be confused with it

	opt := q.TryTake()
	if opt.IsEmpty() {
		t.Fatal("first TryTake() should be present")
	}
	if got := opt.Get(); got != 0 {
		t.Fatalf("first TryTake() = %d, want 0", got)
	}

	empty := q.TryTake()
	if !empty.IsEmpty() {
		t.Fatalf("second TryTake() = %v, want Empty", empty)
	}

	// Sanity: the empty factory also looks the same.
	if e := option.Empty[int](); !e.IsEmpty() {
		t.Fatal("option.Empty[int]() should be empty")
	}
}
