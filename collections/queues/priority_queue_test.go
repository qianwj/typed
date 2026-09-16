package queues

import (
	"cmp"
	"testing"

	"github.com/qianwj/typed/adt"
)

// --- Empty / peek / pop --------------------------------------------------

func TestPriorityQueue_NewIsEmpty(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(cmp.Less[int])
	if got := q.Len(); got != 0 {
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
	q := NewPriorityQueue(cmp.Less[int])
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
	q := NewPriorityQueue(greaterFirst)
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
	q := NewPriorityQueue(byLength)
	q.Push("hi")
	q.Push("hello")
	q.Push("a")
	q.Push("world")

	// Expected dequeue order: "a" (1), "hi" (2), "hello" / "world" (5, ties).
	if got := q.Pop().Get(); got != "a" {
		t.Errorf("Pop[0] = %q, want \"a\"", got)
	}
	if got := q.Pop().Get(); got != "hi" {
		t.Errorf("Pop[1] = %q, want \"hi\"", got)
	}
	// Two 5-char strings remain in unspecified order.
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
	q := NewPriorityQueue(cmp.Less[string])
	q.Push("first")
	q.Push("second")
	q.Push("third")

	seen := make(map[string]bool, 3)
	for q.Len() > 0 {
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
	q := NewPriorityQueue(cmp.Less[int])
	q.Push(42)

	if got := q.Peek().Get(); got != 42 {
		t.Errorf("first Peek = %d, want 42", got)
	}
	if got := q.Len(); got != 1 {
		t.Errorf("Len after Peek = %d, want 1 (Peek must not mutate)", got)
	}
	if got := q.Peek().Get(); got != 42 {
		t.Errorf("second Peek = %d, want 42", got)
	}
}

// --- Len tracks size -----------------------------------------------------

func TestPriorityQueue_LenShrinksAfterPop(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(cmp.Less[int])
	for range 10 {
		q.Push(0)
	}
	if got := q.Len(); got != 10 {
		t.Fatalf("Len after 10 Pushes = %d, want 10", got)
	}
	for range 10 {
		q.Pop()
	}
	if got := q.Len(); got != 0 {
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
	q := NewPriorityQueue(cmp.Less[int])

	// First batch: 0..99, with the comparator ordering ints, so
	// the natural drain order is 0, 1, ..., 99.
	for i := range 100 {
		q.Push(i)
	}
	first := make([]int, 0, 100)
	for q.Len() > 0 {
		first = append(first, q.Pop().Get())
	}
	if !monotonicNonDecreasing(first) {
		t.Fatalf("first drain not sorted: %v", first)
	}

	// Second batch: 99..0, reversed — the natural drain order is
	// still 0, 1, ..., 99 (the comparator doesn't care about
	// insertion order).
	for i := range 100 {
		q.Push(99 - i)
	}
	second := make([]int, 0, 100)
	for q.Len() > 0 {
		second = append(second, q.Pop().Get())
	}
	if !monotonicNonDecreasing(second) {
		t.Fatalf("second drain not sorted: %v", second)
	}
}

// TestPriorityQueue_LargeRandomSequenceIsMonotonicallyNonDecreasing
// pushes a deterministic pseudo-random sequence of values, drains,
// and checks the dequeue order is monotonically non-decreasing AND
// the multiset matches the input. This is the strongest
// single-thread invariant test: it exercises every internal code
// path and verifies both ordering and completeness.
func TestPriorityQueue_LargeRandomSequenceIsMonotonicallyNonDecreasing(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(cmp.Less[int])

	// Deterministic: same sequence every run. Modulus 113 with N=200
	// gives some duplicates, exercising tie-handling.
	const N = 200
	want := make([]int, N)
	for i := range N {
		v := (i * 7) % 113
		want[i] = v
		q.Push(v)
	}
	// Sort want so the comparison is multiset, not order.
	insertionSort(want)

	got := make([]int, 0, N)
	for q.Len() > 0 {
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
		primary int // tie group
		seq     int // FIFO within tie group
	}
	less := func(a, b keyed) bool {
		if a.primary != b.primary {
			return a.primary < b.primary
		}
		return a.seq < b.seq
	}
	q := NewPriorityQueue(less)

	// Push three elements in primary group 1, then three in group
	// 0 — interleaved, to make sure the heap can't rely on
	// insertion order.
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
	q := NewPriorityQueue(cmp.Less[int])
	for range 5 {
		if opt := q.Pop(); opt.IsPresent() {
			t.Fatalf("Pop on empty returned present: %v", opt.Get())
		}
	}
	if got := q.Len(); got != 0 {
		t.Errorf("Len after repeated empty Pops = %d, want 0", got)
	}
}

func TestPriorityQueue_EmptyReturnsOptionNotNil(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue(cmp.Less[string])
	var opt adt.Option[string] = q.Pop()
	if opt.IsPresent() {
		t.Errorf("Pop on empty queue returned present (%v)", opt.Get())
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
