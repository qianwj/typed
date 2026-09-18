package queues

import (
	"cmp"
	"testing"
	"time"

	"github.com/qianwj/typed/adt/option"
)

// --- Empty / peek / pop --------------------------------------------------

func TestPriorityQueue_NewIsEmpty(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	if got := q.Size(); got != 0 {
		t.Fatalf("Len on fresh queue = %d, want 0", got)
	}
	if opt := q.Peek(); !opt.IsEmpty() {
		t.Errorf("Peek on empty queue = present (%v), want empty", opt.Get())
	}
	if opt := q.Pop(); !opt.IsEmpty() {
		t.Errorf("Pop on empty queue = present (%v), want empty", opt.Get())
	}
}

// --- Ordering ------------------------------------------------------------

// TestPriorityQueue_LessDrivesOrder confirms the comparator
// controls the dequeue order. Pushing 2, 0, 1 with cmp.Less[int]
// (smaller first) yields 0, 1, 2 on Pop.
func TestPriorityQueue_LessDrivesOrder(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	q.Push(2)
	q.Push(0)
	q.Push(1)

	for _, want := range []int{0, 1, 2} {
		opt := q.Pop()
		if opt.IsEmpty() {
			t.Fatalf("Pop returned empty; want %d", want)
		}
		if got := opt.Get(); got != want {
			t.Errorf("Pop = %d, want %d", got, want)
		}
	}
	if !q.Pop().IsEmpty() {
		t.Error("Pop after draining returned present; want empty")
	}
}

// TestPriorityQueue_DescendingComparator confirms that swapping
// the comparator's sense reverses the dequeue order. The same
// pushes (2, 0, 1) come out as 2, 1, 0 under "greater first".
func TestPriorityQueue_DescendingComparator(t *testing.T) {
	t.Parallel()
	greaterFirst := func(a, b int) bool { return a > b }
	q := NewPriorityQueue(0, greaterFirst)
	q.Push(2)
	q.Push(0)
	q.Push(1)

	for _, want := range []int{2, 1, 0} {
		if got := q.Pop().Get(); got != want {
			t.Errorf("Pop = %d, want %d", got, want)
		}
	}
}

// TestPriorityQueue_CustomComparatorByLength orders strings by
// length, shortest first. Demonstrates that the comparator can
// encode any total order over T, not just numeric.
func TestPriorityQueue_CustomComparatorByLength(t *testing.T) {
	t.Parallel()
	byLength := func(a, b string) bool { return len(a) < len(b) }
	q := NewPriorityQueue(0, byLength)
	q.Push("hi")
	q.Push("hello")
	q.Push("a")
	q.Push("world")

	if got := q.Pop().Get(); got != "a" {
		t.Errorf("Pop[0] = %q, want \"a\"", got)
	}
	if got := q.Pop().Get(); got != "hi" {
		t.Errorf("Pop[1] = %q, want \"hi\"", got)
	}
	rest := []string{q.Pop().Get(), q.Pop().Get()}
	if !contains(rest, "hello") || !contains(rest, "world") {
		t.Errorf("Pop[2..3] = %v, want both \"hello\" and \"world\"", rest)
	}
}

// TestPriorityQueue_TiesPreservedAsMultiset: with a strict-less
// comparator, equal-priority elements have no guaranteed relative
// order. The heap must not drop or duplicate them.
func TestPriorityQueue_TiesPreservedAsMultiset(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[string])
	q.Push("first")
	q.Push("second")
	q.Push("third")

	seen := make(map[string]bool, 3)
	for q.Size() > 0 {
		seen[q.Pop().Get()] = true
	}
	for _, want := range []string{"first", "second", "third"} {
		if !seen[want] {
			t.Errorf("Pop sequence missing %q", want)
		}
	}
	if len(seen) != 3 {
		t.Errorf("Pop yielded %d distinct values, want 3", len(seen))
	}
}

// --- Peek doesn't mutate -------------------------------------------------

func TestPriorityQueue_PeekDoesNotMutate(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	q.Push(42)

	if got := q.Peek().Get(); got != 42 {
		t.Errorf("first Peek = %d, want 42", got)
	}
	if got := q.Size(); got != 1 {
		t.Errorf("Len after Peek = %d, want 1 (Peek must not mutate)", got)
	}
	if got := q.Peek().Get(); got != 42 {
		t.Errorf("second Peek = %d, want 42", got)
	}
}

// --- Len tracks size -----------------------------------------------------

func TestPriorityQueue_LenShrinksAfterPop(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	for range 10 {
		q.Push(0)
	}
	if got := q.Size(); got != 10 {
		t.Fatalf("Len after 10 Pushes = %d, want 10", got)
	}
	for range 10 {
		q.Pop()
	}
	if got := q.Size(); got != 0 {
		t.Errorf("Len after 10 Pops = %d, want 0", got)
	}
}

// --- Heap invariant under repeated drain-and-refill --------------------

// TestPriorityQueue_DrainAndRefillPreservesHeapInvariant drains
// the queue, refills with mixed values, drains again, and checks
// the second drain is monotonically non-decreasing under less.
// Catches "stale heap state" bugs.
func TestPriorityQueue_DrainAndRefillPreservesHeapInvariant(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])

	for i := range 100 {
		q.Push(i)
	}
	first := make([]int, 0, 100)
	for q.Size() > 0 {
		first = append(first, q.Pop().Get())
	}
	if !monotonicNonDecreasing(first) {
		t.Fatalf("first drain not sorted: %v", first)
	}

	for i := range 100 {
		q.Push(99 - i)
	}
	second := make([]int, 0, 100)
	for q.Size() > 0 {
		second = append(second, q.Pop().Get())
	}
	if !monotonicNonDecreasing(second) {
		t.Fatalf("second drain not sorted: %v", second)
	}
}

// TestPriorityQueue_LargeRandomSequenceIsMonotonicallyNonDecreasing
// pushes a deterministic pseudo-random sequence of values, drains,
// and checks the dequeue order is monotonically non-decreasing AND
// the multiset matches the input.
func TestPriorityQueue_LargeRandomSequenceIsMonotonicallyNonDecreasing(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])

	const N = 200
	want := make([]int, N)
	for i := range N {
		v := (i * 7) % 113
		want[i] = v
		q.Push(v)
	}
	insertionSort(want)

	got := make([]int, 0, N)
	for q.Size() > 0 {
		got = append(got, q.Pop().Get())
	}

	if !monotonicNonDecreasing(got) {
		t.Fatalf("dequeue sequence not sorted: %v", got)
	}
	if len(got) != N {
		t.Fatalf("dequeue length = %d, want %d", len(got), N)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("multiset mismatch at index %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

// TestPriorityQueue_StableOrderViaCompositeKey demonstrates the
// canonical "encode the tiebreaker in the comparator" pattern.
// Two elements with the same primary key come out in seq order,
// even though the heap itself doesn't store seq separately.
func TestPriorityQueue_StableOrderViaCompositeKey(t *testing.T) {
	t.Parallel()

	type keyed struct {
		primary int
		seq     int
	}
	less := func(a, b keyed) bool {
		if a.primary != b.primary {
			return a.primary < b.primary
		}
		return a.seq < b.seq
	}
	q := NewPriorityQueue(0, less)

	q.Push(keyed{primary: 1, seq: 10})
	q.Push(keyed{primary: 0, seq: 1})
	q.Push(keyed{primary: 1, seq: 11})
	q.Push(keyed{primary: 0, seq: 2})
	q.Push(keyed{primary: 1, seq: 12})
	q.Push(keyed{primary: 0, seq: 3})

	want := []keyed{
		{primary: 0, seq: 1},
		{primary: 0, seq: 2},
		{primary: 0, seq: 3},
		{primary: 1, seq: 10},
		{primary: 1, seq: 11},
		{primary: 1, seq: 12},
	}
	for i, w := range want {
		if got := q.Pop().Get(); got != w {
			t.Errorf("Pop[%d] = %+v, want %+v", i, got, w)
		}
	}
}

// --- Empty Option type assertion -----------------------------------------

func TestPriorityQueue_PopOnEmptyRepeatedlyReturnsEmpty(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	for range 5 {
		if opt := q.Pop(); opt.IsPresent() {
			t.Fatalf("Pop on empty returned present: %v", opt.Get())
		}
	}
	if got := q.Size(); got != 0 {
		t.Errorf("Len after repeated empty Pops = %d, want 0", got)
	}
}

func TestPriorityQueue_EmptyReturnsOptionNotNil(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[string])
	var opt option.Option[string] = q.Pop()
	if opt.IsPresent() {
		t.Errorf("Pop on empty queue returned present (%v)", opt.Get())
	}
}

// --- Capacity / top-K ----------------------------------------------------

// TestPriorityQueue_CapacityUnboundedReportsZero confirms the
// "0 = unbounded" reading of Capacity: a queue built with
// capacity 0 reports Capacity() == 0.
func TestPriorityQueue_CapacityUnboundedReportsZero(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	if got := q.Capacity(); got != 0 {
		t.Errorf("Capacity on unbounded queue = %d, want 0", got)
	}
}

// TestPriorityQueue_CapacityBoundedReportsN confirms Capacity
// returns the value passed at construction.
func TestPriorityQueue_CapacityBoundedReportsN(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(42, cmp.Less[int])
	if got := q.Capacity(); got != 42 {
		t.Errorf("Capacity = %d, want 42", got)
	}
}

// TestPriorityQueue_PushPanicsOnNegativeCapacity confirms the
// constructor rejects a misconfigured capacity loudly.
func TestPriorityQueue_PushPanicsOnNegativeCapacity(t *testing.T) {
	t.Parallel()

	done := make(chan any, 1)
	go func() {
		defer func() {
			done <- recover()
		}()
		NewPriorityQueue(-1, cmp.Less[int])
	}()

	select {
	case r := <-done:
		if r == nil {
			t.Fatal("NewPriorityQueue(-1, ...) did not panic")
		}
	case <-time.After(time.Second):
		t.Fatal("NewPriorityQueue(-1, ...) hung instead of panicking")
	}
}

// TestPriorityQueue_UnboundedPushAlwaysReturnsTrue: in an
// unbounded queue (capacity 0), every Push succeeds regardless of
// how full the queue is.
func TestPriorityQueue_UnboundedPushAlwaysReturnsTrue(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(0, cmp.Less[int])
	for i := range 1000 {
		if !q.Push(i) {
			t.Errorf("Push(%d) returned false on unbounded queue", i)
		}
	}
	if got := q.Size(); got != 1000 {
		t.Errorf("Len after 1000 Pushes = %d, want 1000", got)
	}
}

// TestPriorityQueue_BoundedPushBelowCapacityAcceptsAll: until the
// queue reaches capacity, every Push returns true and the queue
// grows.
func TestPriorityQueue_BoundedPushBelowCapacityAcceptsAll(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(5, cmp.Less[int])
	for i := range 5 {
		if !q.Push(i) {
			t.Errorf("Push(%d) returned false below capacity", i)
		}
	}
	if got := q.Size(); got != 5 {
		t.Errorf("Len after 5 Pushes = %d, want 5", got)
	}
}

// TestPriorityQueue_BoundedPushAcceptsHigherPriorityReplacesBoundary
// is the core top-K test: once the queue is full, pushing a new
// element that out-prioritises the boundary (the lowest-priority
// element in the heap, which is the eviction candidate) must
// replace the boundary, and Push must return true. The current
// highest-priority element (the root) is preserved.
func TestPriorityQueue_BoundedPushAcceptsHigherPriorityReplacesBoundary(t *testing.T) {
	t.Parallel()
	// Min-heap: smaller = higher priority. Push 1..5 — queue
	// ends up holding [1,2,3,4,5] with root=1 (the smallest = highest
	// priority) and boundary=5 (the largest = lowest priority).
	q := NewPriorityQueue(5, cmp.Less[int])
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	// Verify the initial top.
	if got := q.Peek().Get(); got != 1 {
		t.Fatalf("Peek before replacement = %d, want 1 (root of min-heap is smallest)", got)
	}

	// Push 0 — out-prioritises the boundary (5). Boundary should be
	// replaced; the queue ends up holding {0, 1, 2, 3, 4} with root=0.
	if !q.Push(0) {
		t.Fatal("Push(0) returned false; want true (0 < boundary=5)")
	}
	if got := q.Peek().Get(); got != 0 {
		t.Errorf("Peek after replacement = %d, want 0 (0 should be at root)", got)
	}
	if got := q.Size(); got != 5 {
		t.Errorf("Len after replacement = %d, want 5 (capacity unchanged)", got)
	}

	// Drain and verify the queue holds {0, 1, 2, 3, 4}.
	drained := make([]int, 0, 5)
	for q.Size() > 0 {
		drained = append(drained, q.Pop().Get())
	}
	wantDrained := []int{0, 1, 2, 3, 4}
	if !equalSlices(drained, wantDrained) {
		t.Errorf("drained = %v, want %v", drained, wantDrained)
	}
}

// TestPriorityQueue_BoundedPushDropsLowerPriorityThanBoundary is
// the complement: once the queue is full, pushing a new element
// that is NOT higher priority than the boundary must be dropped,
// and Push must return false. The root is preserved.
func TestPriorityQueue_BoundedPushDropsLowerPriorityThanBoundary(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(5, cmp.Less[int])
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	// Boundary is 5 (the largest). Push 10 — does NOT out-prioritise
	// boundary. Drop.
	if q.Push(10) {
		t.Error("Push(10) returned true; want false (10 not < boundary=5)")
	}
	if got := q.Size(); got != 5 {
		t.Errorf("Len after drop = %d, want 5 (capacity unchanged)", got)
	}
	if got := q.Peek().Get(); got != 1 {
		t.Errorf("Peek after drop = %d, want 1 (root unchanged)", got)
	}
}

// TestPriorityQueue_BoundedPushEqualToBoundaryDrops: ties against
// the boundary are NOT considered "higher priority" — strict-less
// comparator means less(boundary, boundary) is false, so a tied
// push is dropped. This preserves the multiset invariant when
// the comparator is strict-less: a tied push would replace an
// equally-bad element with another equally-bad element, swapping
// one for the other — no information gain.
func TestPriorityQueue_BoundedPushEqualToBoundaryDrops(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(3, cmp.Less[int])
	q.Push(5)
	q.Push(3)
	q.Push(7)
	// Heap: 3 at root, 5 and 7 below. Boundary is 7.
	if got := q.Peek().Get(); got != 3 {
		t.Fatalf("Peek = %d, want 3", got)
	}
	// Push 7 — tied with boundary. Strict-less: less(7, 7) is
	// false, so drop.
	if q.Push(7) {
		t.Error("Push(7) returned true; want false (tied with boundary)")
	}
	if got := q.Size(); got != 3 {
		t.Errorf("Len after tied push = %d, want 3", got)
	}
}

// TestPriorityQueue_BoundedHoldsTopK is the integration test:
// push N=100 elements with priorities 0..99 into a top-10 queue
// and verify the queue ends up holding {0..9}. This catches the
// whole top-K pipeline: root comparison, replacement, sift down.
func TestPriorityQueue_BoundedHoldsTopK(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(10, cmp.Less[int])
	for i := range 100 {
		q.Push(i)
	}
	if got := q.Size(); got != 10 {
		t.Fatalf("Len after 100 pushes into top-10 = %d, want 10", got)
	}
	drained := make([]int, 0, 10)
	for q.Size() > 0 {
		drained = append(drained, q.Pop().Get())
	}
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !equalSlices(drained, want) {
		t.Errorf("drained = %v, want %v (the top-10 priorities)", drained, want)
	}
}

// TestPriorityQueue_BoundedReverseOrder: pushing in reverse order
// into a top-K queue still ends up holding the top K. Catches
// "depends on insertion order" bugs.
func TestPriorityQueue_BoundedReverseOrder(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(5, cmp.Less[int])
	for i := 99; i >= 0; i-- {
		q.Push(i)
	}
	if got := q.Size(); got != 5 {
		t.Fatalf("Len = %d, want 5", got)
	}
	drained := make([]int, 0, 5)
	for q.Size() > 0 {
		drained = append(drained, q.Pop().Get())
	}
	want := []int{0, 1, 2, 3, 4}
	if !equalSlices(drained, want) {
		t.Errorf("drained = %v, want %v", drained, want)
	}
}

// --- Helpers -------------------------------------------------------------

func monotonicNonDecreasing(s []int) bool {
	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false
		}
	}
	return true
}

func contains(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func equalSlices(a, b []int) bool {
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

// insertionSort is a tiny in-place stable sort. Used here to
// compute the expected multiset order without pulling in sort.
func insertionSort(s []int) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i - 1
		for j >= 0 && s[j] > v {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = v
	}
}
