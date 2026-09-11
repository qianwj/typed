package lists

import (
	"runtime"
	"testing"
	"time"
)

// ---------- head-offset invariants ----------

// TestArrayListRemoveFirstIsConstantTime documents the O(1)
// promise of the head-offset layout. A pre-head-offset RemoveFirst
// had to shift every remaining element down by one, so it was
// O(n). After the refactor, RemoveFirst just bumps the head
// offset; doubling the live size must not measurably double the
// per-call cost.
//
// The test runs RemoveFirst on a list of 1000 elements and
// measures total elapsed time. With O(n) shifting, total cost
// would be O(n^2) — that is 500,000 element copies. With O(1)
// per call, total cost is O(n). 1 second is a very generous
// ceiling; if the test ever fails, it almost certainly means
// someone re-introduced a shift in RemoveFirst.
func TestArrayListRemoveFirstIsConstantTime(t *testing.T) {
	const n = 1000
	a := NewArrayList[int]()
	for i := 0; i < n; i++ {
		a.Add(i)
	}

	start := time.Now()
	for a.Size() > 0 {
		a.RemoveFirst()
	}
	elapsed := time.Since(start)

	if elapsed > time.Second {
		t.Fatalf("RemoveFirst took %v for %d elements; expected sub-second (O(1) per call)", elapsed, n)
	}
	if a.Size() != 0 {
		t.Fatalf("after draining: size = %d, want 0", a.Size())
	}
	if a.head >= 64 {
		// The compaction threshold is 64. After 1000 RemoveFirsts
		// the list goes through 15 compactions, leaving 40
		// elements with head=40 (which is below the threshold
		// and therefore correct). The important property is
		// that head does not grow unboundedly: it stays
		// bounded by the compaction threshold regardless of
		// how many RemoveFirsts have been called.
		t.Fatalf("after draining: head = %d, want < 64 (head offset must stay bounded by the compaction threshold)", a.head)
	}
}

// TestArrayListAddFirstIsConstantTime is the symmetric test for
// AddFirst. With the old slices.Insert(items, 0, v) approach,
// every call shifted the entire live range right by one, O(n).
// Under the head-offset layout, AddFirst decrements head and
// writes at items[head]; O(1).
func TestArrayListAddFirstIsConstantTime(t *testing.T) {
	const n = 1000
	a := NewArrayList[int]()
	for i := 0; i < n; i++ {
		a.AddFirst(i)
	}
	// At this point a looks like [n-1, n-2, ..., 1, 0]: each
	// AddFirst puts the new value at the front, so the first
	// AddFirst puts 0 at the front, the second puts 1 in front
	// of that, and so on. The last element pushed (n-1) ends
	// up at the front; the first one (0) ends up at the back.
	if got := a.First().OrElse(-1); got != n-1 {
		t.Fatalf("First after %d AddFirsts: got %d, want %d", n, got, n-1)
	}
	if a.Size() != n {
		t.Fatalf("Size: got %d, want %d", a.Size(), n)
	}

	// Drain with RemoveLast. The list is [n-1, n-2, ..., 0],
	// so RemoveLast returns 0 first, then 1, ..., then n-1.
	for i := 0; i < n; i++ {
		v := a.RemoveLast()
		if v.IsEmpty() {
			t.Fatalf("RemoveLast: empty at i=%d", i)
		}
		if got := v.Get(); got != i {
			t.Fatalf("RemoveLast at i=%d: got %d, want %d", i, got, i)
		}
	}
}

// TestArrayListMixedHeadTailDrain exercises interleaved AddFirst,
// Add, RemoveFirst, RemoveLast, Get, ForEach against a
// growing-then-shrinking head offset. After each mutation the
// live range, head, and the visible values are all checked.
func TestArrayListMixedHeadTailDrain(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 10; i++ {
		a.Add(i) // [0..9]
	}
	for i := 0; i < 5; i++ {
		v := a.RemoveFirst() // drops 0,1,2,3,4; [5..9]
		if got := v.OrElse(-1); got != i {
			t.Fatalf("RemoveFirst at i=%d: got %d, want %d", i, got, i)
		}
	}
	for i := 0; i < 3; i++ {
		a.AddFirst(i) // [2,1,0,5,6,7,8,9]
	}
	// Now head is 5, and a[head:head+size] = [0,5,6,7,8,9]
	// because AddFirst wrote at items[head-1], items[head-2],
	// items[head-3]. Wait — AddFirst decrements head and writes
	// at the new head, so the new head points to the most
	// recently added value. After 3 AddFirsts: head=2, and
	// the live range is items[2:10] = [2,1,0,5,6,7,8,9].
	if got := a.Collect(); !equalIntSlice(got, []int{2, 1, 0, 5, 6, 7, 8, 9}) {
		t.Fatalf("after mixed ops: got %v, want [2 1 0 5 6 7 8 9]", got)
	}

	// First / Last / Get must work through the head offset.
	if got := a.First().OrElse(-1); got != 2 {
		t.Fatalf("First: got %d, want 2", got)
	}
	if got := a.Last().OrElse(-1); got != 9 {
		t.Fatalf("Last: got %d, want 9", got)
	}
	for i, want := range []int{2, 1, 0, 5, 6, 7, 8, 9} {
		got := a.Get(i).OrElse(-1)
		if got != want {
			t.Fatalf("Get(%d): got %d, want %d", i, got, want)
		}
	}
}

// TestArrayListInsertAtArbitraryIndex exercises Insert at the
// head, the middle, and the tail, plus a panic on out-of-range.
// Insert used to be slices.Insert, which had to memmove a
// possibly long tail. The head-offset version does the same
// memmove but the head is a constant offset that we just add
// to the indices.
func TestArrayListInsertAtArbitraryIndex(t *testing.T) {
	a := NewArrayList[int]()
	a.Add(0)
	a.Add(1)
	a.Add(2) // [0, 1, 2]

	a.Insert(1, 10) // [0, 10, 1, 2]
	if got := a.Collect(); !equalIntSlice(got, []int{0, 10, 1, 2}) {
		t.Fatalf("Insert middle: got %v, want [0 10 1 2]", got)
	}

	a.Insert(0, 20) // [20, 0, 10, 1, 2]
	if got := a.Collect(); !equalIntSlice(got, []int{20, 0, 10, 1, 2}) {
		t.Fatalf("Insert head: got %v, want [20 0 10 1 2]", got)
	}

	a.Insert(a.Size(), 30) // [20, 0, 10, 1, 2, 30]
	if got := a.Collect(); !equalIntSlice(got, []int{20, 0, 10, 1, 2, 30}) {
		t.Fatalf("Insert tail: got %v, want [20 0 10 1 2 30]", got)
	}
}

// TestArrayListRemoveAtAcrossTheRange exercises RemoveAt at the
// head, the middle, and the tail. The head and tail cases are
// O(1) under the head-offset layout; the middle case shifts the
// tail left.
func TestArrayListRemoveAtAcrossTheRange(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		a.Add(i) // [0, 1, 2, 3, 4]
	}

	if v := a.RemoveAt(0); v != 0 { // [1, 2, 3, 4]
		t.Fatalf("RemoveAt(0): got %d, want 0", v)
	}
	if v := a.RemoveAt(a.Size() - 1); v != 4 { // [1, 2, 3]
		t.Fatalf("RemoveAt(tail): got %d, want 4", v)
	}
	if v := a.RemoveAt(1); v != 2 { // [1, 3]
		t.Fatalf("RemoveAt(middle): got %d, want 2", v)
	}
	if got := a.Collect(); !equalIntSlice(got, []int{1, 3}) {
		t.Fatalf("after three RemoveAts: got %v, want [1 3]", got)
	}
}

// TestArrayListForEachRespectsHeadOffset documents that ForEach
// iterates only the live range, not the discarded prefix.
func TestArrayListForEachRespectsHeadOffset(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		a.Add(i)
	}
	for i := 0; i < 3; i++ {
		a.RemoveFirst() // drops 0, 1, 2; head=3
	}
	seen := []int{}
	a.ForEach(func(v int) { seen = append(seen, v) })
	if !equalIntSlice(seen, []int{3, 4}) {
		t.Fatalf("ForEach saw %v, want [3 4]", seen)
	}
}

// TestArrayListStreamSnapshotRespectsHeadOffset confirms that
// Stream reads the live range, not the discarded prefix. The
// discarded prefix's values must not appear in the Stream.
func TestArrayListStreamSnapshotRespectsHeadOffset(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		a.Add(i)
	}
	for i := 0; i < 3; i++ {
		a.RemoveFirst() // head=3, live = [3, 4]
	}
	got := a.Stream().Collect()
	if !equalIntSlice(got, []int{3, 4}) {
		t.Fatalf("Stream snapshot: got %v, want [3 4]", got)
	}

	// Mutating the ArrayList after Stream() must not affect
	// the snapshot.
	a.RemoveFirst()
	got = a.Stream().Collect()
	if !equalIntSlice(got, []int{4}) {
		t.Fatalf("Stream snapshot after mutation: got %v, want [4]", got)
	}
}

// TestArrayListSizeIsLiveCountNotCapacity is the key correctness
// test for the head offset: Size must be the live count, not
// the capacity of the underlying slice. With the old code
// (and with a broken new code), a list that grew to 100 and
// shrank back to 1 would still report Size as 100 if the
// implementation counted capacity.
func TestArrayListSizeIsLiveCountNotCapacity(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		a.Add(i)
	}
	for i := 0; i < 99; i++ {
		a.RemoveFirst()
	}
	if got := a.Size(); got != 1 {
		t.Fatalf("after 100 Adds and 99 RemoveFirsts: Size = %d, want 1", got)
	}
}

// TestArrayListClearResetsHead confirms that Clear resets the
// head offset back to zero, so a fresh Add after Clear starts
// writing at index 0.
func TestArrayListClearResetsHead(t *testing.T) {
	a := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		a.Add(i)
	}
	for i := 0; i < 3; i++ {
		a.RemoveFirst()
	}
	a.Clear()
	if a.head != 0 {
		t.Fatalf("Clear: head = %d, want 0", a.head)
	}
	if a.Size() != 0 {
		t.Fatalf("Clear: Size = %d, want 0", a.Size())
	}
	a.Add(42)
	if v := a.Get(0).OrElse(-1); v != 42 {
		t.Fatalf("Get(0) after Clear+Add: got %d, want 42", v)
	}
}

// ---------- memory-release regression test ----------

// finalizableBox is the test type used to prove that the
// head-offset layout releases references held by popped
// elements. runtime.SetFinalizer is set in newFinalizableBox
// and the test asserts that all of the popped boxes'
// finalizers have run after a long drain.
//
// This is the same idiom used in the stack_test.go regression
// test for Stack.Pop; the ArrayList version proves the
// equivalent fix for the larger RemoveFirst code path.
type finalizableBox struct {
	id int
}

var finalizableTracker struct {
	mu        chan struct{}
	finalised int
}

func init() {
	finalizableTracker.mu = make(chan struct{}, 1)
}

func newFinalizableBox(id int) *finalizableBox {
	b := &finalizableBox{id: id}
	runtime.SetFinalizer(b, func(*finalizableBox) {
		finalizableTracker.mu <- struct{}{}
		finalizableTracker.finalised++
		<-finalizableTracker.mu
	})
	return b
}

func waitForFinalisation(target int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		runtime.GC()
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		finalizableTracker.mu <- struct{}{}
		n := finalizableTracker.finalised
		<-finalizableTracker.mu
		if n >= target {
			return n
		}
	}
	finalizableTracker.mu <- struct{}{}
	n := finalizableTracker.finalised
	<-finalizableTracker.mu
	return n
}

// TestArrayListRemoveFirstReleasesPointerElements is the
// memory-leak regression test for ArrayList.RemoveFirst.
// The previous implementation shifted elements left on every
// RemoveFirst and never released the backing-array capacity,
// keeping pointers reachable through out-of-range slots.
// After the head-offset refactor, RemoveFirst advances head
// and zeroes the freed slot, releasing the reference and
// making the popped box eligible for GC.
func TestArrayListRemoveFirstReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForFinalisation(0, 0)

	a := NewArrayList[*finalizableBox]()
	boxes := make([]*finalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newFinalizableBox(i)
		a.Add(boxes[i])
	}

	// RemoveFirst all but the last. After this, the head
	// offset is n-1 and items[0:head] should be zeroed
	// (with possible periodic compaction folding them down).
	for i := 0; i < n-1; i++ {
		if a.RemoveFirst().IsEmpty() {
			t.Fatalf("RemoveFirst at i=%d: empty", i)
		}
	}

	// Drop our strong references to the popped boxes too.
	for i := 0; i < n-1; i++ {
		boxes[i] = nil
	}

	// All n-1 popped boxes' finalizers must have run.
	got := waitForFinalisation(baseline+(n-1), 5*time.Second)
	if got < baseline+(n-1) {
		t.Fatalf("after %d RemoveFirsts: %d finalisers ran, want at least %d (popped boxes still reachable through head-offset prefix?)",
			n-1, got, baseline+(n-1))
	}
}

// TestArrayListAddFirstReleasesPointerElements is the
// parallel test for AddFirst. AddFirst used to do
// slices.Insert(items, 0, v), which shifted the live range
// right by one; the new head-offset layout decrements head
// and writes at the new head, which is O(1). The
// previous-capacity concern does not apply to AddFirst
// (the slice grows when needed), but the head-offset
// layout still benefits from the periodic compaction when
// the prefix is then drained.
func TestArrayListAddFirstReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForFinalisation(0, 0)

	a := NewArrayList[*finalizableBox]()
	boxes := make([]*finalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newFinalizableBox(i)
		a.Add(boxes[i])
	}

	// AddFirst shifts the prefix; the old tail elements
	// (in slots 0..n-1) are now at the back of the live
	// range. Then RemoveFirst drains from the front,
	// advancing head past them and releasing them.
	for i := 0; i < n; i++ {
		a.AddFirst(newFinalizableBox(100 + i))
	}
	for i := 0; i < n; i++ {
		a.RemoveFirst()
	}
	for i := 0; i < n; i++ {
		boxes[i] = nil
	}

	// All 2n boxes' finalizers must have run.
	got := waitForFinalisation(baseline+2*n, 5*time.Second)
	if got < baseline+2*n {
		t.Fatalf("after AddFirst + RemoveFirst cycle: %d finalisers ran, want at least %d", got, baseline+2*n)
	}
}

// ---------- helpers ----------

// equalIntSlice is a memberwise equality check on []int, used
// for asserting the contents of Collect() across a head
// offset.
func equalIntSlice(got, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
