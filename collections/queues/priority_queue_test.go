package queues

import (
	"testing"

	"github.com/qianwj/typed/adt"
)

// --- Empty / peek / pop --------------------------------------------------

func TestPriorityQueue_NewIsEmpty(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()
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

// --- Push / Pop order ----------------------------------------------------

// TestPriorityQueue_MinPriorityDequeuedFirst confirms the
// min-heap invariant: among (p=2, p=0, p=1) pushed in that order,
// Pop returns 0, then 1, then 2.
func TestPriorityQueue_MinPriorityDequeuedFirst(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()
	q.PushWithPriority(2, 2)
	q.PushWithPriority(0, 0)
	q.PushWithPriority(1, 1)

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

// TestPriorityQueue_TiesPreservedAsMultiset: among equal
// priorities, Pop returns each element exactly once — order is
// not guaranteed (a binary heap with strict-less comparison has
// no natural tiebreaker; this test only checks that the heap
// doesn't drop or duplicate items).
func TestPriorityQueue_TiesPreservedAsMultiset(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[string]()
	q.PushWithPriority("first", 5)
	q.PushWithPriority("second", 5)
	q.PushWithPriority("third", 5)

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

// --- Push uses default priority 0 ----------------------------------------

// TestPriorityQueue_PushUsesZeroPriority confirms that bare Push
// is equivalent to PushWithPriority(_, 0). Items pushed with Push
// sort among items pushed with PushWithPriority at priority 0.
func TestPriorityQueue_PushUsesZeroPriority(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[string]()
	q.PushWithPriority("high", -1)
	q.Push("middle") // priority 0
	q.PushWithPriority("low", 1)

	for _, want := range []string{"high", "middle", "low"} {
		if got := q.Pop().Get(); got != want {
			t.Errorf("Pop = %q, want %q", got, want)
		}
	}
}

// --- Peek doesn't mutate -------------------------------------------------

func TestPriorityQueue_PeekDoesNotMutate(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()
	q.PushWithPriority(42, 0)

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
	q := NewPriorityQueue[int]()
	for range 10 {
		q.PushWithPriority(0, 0)
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
// the queue, refills with mixed priorities, drains again, and
// checks the second drain is monotonically non-decreasing by
// priority. Catches "stale heap state" bugs.
func TestPriorityQueue_DrainAndRefillPreservesHeapInvariant(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()

	// First batch: priorities (i*7)%11 for i in 0..100.
	for i := range 100 {
		q.PushWithPriority(i, (i*7)%11)
	}
	priorities := make([]int, 0, 100)
	for q.Len() > 0 {
		v := q.Pop().Get()
		priorities = append(priorities, (v*7)%11)
	}
	if !monotonicNonDecreasing(priorities) {
		t.Fatalf("first drain not sorted: %v", priorities)
	}

	// Second batch: priorities reversed.
	for i := range 100 {
		q.PushWithPriority(i, 100-i)
	}
	priorities = priorities[:0]
	for q.Len() > 0 {
		v := q.Pop().Get()
		priorities = append(priorities, 100-v)
	}
	if !monotonicNonDecreasing(priorities) {
		t.Fatalf("second drain not sorted: %v", priorities)
	}
}

// TestPriorityQueue_LargeRandomOrderIsMonotonicallyNonDecreasing
// is the strongest invariant test: push priorities from a
// deterministic pseudo-random sequence with duplicates, drain,
// and check the dequeue order is monotonically non-decreasing
// (each item's priority >= previous item's priority). Ties may
// come out in any order — only the multiset ordering matters.
func TestPriorityQueue_LargeRandomOrderIsMonotonicallyNonDecreasing(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()

	const N = 200
	priorities := make([]int, N)
	for i := range N {
		p := (i * 7) % 11 // 0..10, with duplicates
		priorities[i] = p
		q.PushWithPriority(i, p)
	}

	got := make([]int, 0, N)
	for q.Len() > 0 {
		v := q.Pop().Get()
		got = append(got, (v*7)%11)
	}

	if len(got) != len(priorities) {
		t.Fatalf("dequeue length = %d, want %d", len(got), len(priorities))
	}
	if !monotonicNonDecreasing(got) {
		t.Fatalf("dequeue sequence not sorted: %v", got)
	}
}

// TestPriorityQueue_PopOnEmptyRepeatedlyReturnsEmpty: repeated
// Pops on an empty queue must keep returning absent Options,
// not panic or alternate between empty / present.
func TestPriorityQueue_PopOnEmptyRepeatedlyReturnsEmpty(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[int]()
	for range 5 {
		if opt := q.Pop(); opt.IsPresent() {
			t.Fatalf("Pop on empty returned present: %v", opt.Get())
		}
	}
	if got := q.Len(); got != 0 {
		t.Errorf("Len after repeated empty Pops = %d, want 0", got)
	}
}

// --- Empty Option type assertion -----------------------------------------

func TestPriorityQueue_EmptyReturnsOptionNotNil(t *testing.T) {
	t.Parallel()
	q := NewPriorityQueue[string]()
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

// stableSort is a tiny insertion sort. Stable, so equal-priority
// items keep their input order — which is what the heap's
// FIFO tiebreaker promises.
func stableSort(s []int) {
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
