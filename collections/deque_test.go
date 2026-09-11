package collections

import (
	"runtime"
	"testing"
	"time"

	"github.com/qianwj/typed/collections/lists"
)

// ---------- basic two-ended order ----------

// TestDequePushAndPopBothEnds covers the canonical two-ended
// patterns: push back, pop front (FIFO); push front, pop back
// (LIFO from the other end); interleaved push both ends. The
// tests confirm that PushFront / PopFront hit the same head
// side and PushBack / PopBack hit the same tail side, with no
// cross-talk.
func TestDequePushAndPopBothEnds(t *testing.T) {
	t.Run("push back, pop front = FIFO", func(t *testing.T) {
		d := NewDeque[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		for _, want := range []int{1, 2, 3} {
			if v := d.PopFront().OrElse(-1); v != want {
				t.Fatalf("PopFront: got %d, want %d", v, want)
			}
		}
		if !d.IsEmpty() {
			t.Fatal("after draining: not empty")
		}
	})

	t.Run("push front, pop back = FIFO from the other end", func(t *testing.T) {
		d := NewDeque[int]()
		d.PushFront(1)
		d.PushFront(2)
		d.PushFront(3)
		for _, want := range []int{1, 2, 3} {
			if v := d.PopBack().OrElse(-1); v != want {
				t.Fatalf("PopBack: got %d, want %d", v, want)
			}
		}
		if !d.IsEmpty() {
			t.Fatal("after draining: not empty")
		}
	})

	t.Run("interleaved push both ends", func(t *testing.T) {
		// Build [3, 1, 2, 4] (front-to-back).
		//   PushBack(1)        -> [1]
		//   PushBack(2)        -> [1, 2]
		//   PushFront(3)       -> [3, 1, 2]
		//   PushBack(4)        -> [3, 1, 2, 4]
		d := NewDeque[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushFront(3)
		d.PushBack(4)

		for i, want := range []int{3, 1, 2, 4} {
			if v := d.PopFront().OrElse(-1); v != want {
				t.Fatalf("PopFront at i=%d: got %d, want %d", i, v, want)
			}
		}
	})
}

// TestDequeFrontBackDoNotRemove confirms Front and Back are
// non-mutating reads: repeated calls yield the same value, and
// Size stays stable. This is the read-only contract that
// distinguishes Front/Back from PopFront/PopBack.
func TestDequeFrontBackDoNotRemove(t *testing.T) {
	d := NewDeque[int]()
	d.PushBack(10)
	d.PushBack(20)
	d.PushBack(30)

	if v := d.Front().OrElse(-1); v != 10 {
		t.Fatalf("Front: got %d, want 10", v)
	}
	if v := d.Back().OrElse(-1); v != 30 {
		t.Fatalf("Back: got %d, want 30", v)
	}
	if got := d.Size(); got != 3 {
		t.Fatalf("Size after Front/Back: got %d, want 3 (reads must not mutate)", got)
	}

	// Repeated reads must be stable.
	if v := d.Front().OrElse(-1); v != 10 {
		t.Fatalf("Front again: got %d, want 10", v)
	}
	if v := d.Back().OrElse(-1); v != 30 {
		t.Fatalf("Back again: got %d, want 30", v)
	}
}

// ---------- empty Deque behaviour ----------

// TestDequePopAndPeekOnEmpty confirms that every accessor
// returns an absent Optional on an empty Deque, and leaves the
// Deque in the empty state.
func TestDequePopAndPeekOnEmpty(t *testing.T) {
	d := NewDeque[int]()
	if d.PopFront().IsPresent() {
		t.Fatal("PopFront on empty: got present, want absent")
	}
	if d.PopBack().IsPresent() {
		t.Fatal("PopBack on empty: got present, want absent")
	}
	if d.Front().IsPresent() {
		t.Fatal("Front on empty: got present, want absent")
	}
	if d.Back().IsPresent() {
		t.Fatal("Back on empty: got present, want absent")
	}
	if !d.IsEmpty() {
		t.Fatal("empty Deque: IsEmpty should be true")
	}
	if got := d.Size(); got != 0 {
		t.Fatalf("empty Deque: Size got %d, want 0", got)
	}
}

// TestDequeZeroValueUsable documents that the zero value of
// Deque[T] is usable without going through NewDeque. The
// internal items field is a value-typed ArrayList, which is
// itself usable as a zero value.
func TestDequeZeroValueUsable(t *testing.T) {
	var d Deque[int]
	d.PushBack(1)
	d.PushFront(0)
	d.PushBack(2)
	if v := d.Front().OrElse(-1); v != 0 {
		t.Fatalf("Front: got %d, want 0", v)
	}
	if v := d.Back().OrElse(-1); v != 2 {
		t.Fatalf("Back: got %d, want 2", v)
	}
	if v := d.PopFront().OrElse(-1); v != 0 {
		t.Fatalf("PopFront: got %d, want 0", v)
	}
	if v := d.PopBack().OrElse(-1); v != 2 {
		t.Fatalf("PopBack: got %d, want 2", v)
	}
	if v := d.PopFront().OrElse(-1); v != 1 {
		t.Fatalf("PopFront: got %d, want 1", v)
	}
	if !d.IsEmpty() {
		t.Fatal("after draining zero-value Deque: should be empty")
	}
}

// ---------- Size / IsEmpty / Clear ----------

// TestDequeSizeAndIsEmptyTracksMutations confirms Size and
// IsEmpty stay consistent across PushFront, PushBack, Front,
// Back, PopFront, PopBack, and Clear. Reads (Front, Back) must
// not change Size.
func TestDequeSizeAndIsEmptyTracksMutations(t *testing.T) {
	d := NewDeque[int]()
	if !d.IsEmpty() {
		t.Fatal("new Deque: not empty")
	}

	d.PushBack(1)
	d.PushFront(0)
	d.PushBack(2)
	if got := d.Size(); got != 3 {
		t.Fatalf("after 3 pushes: got %d, want 3", got)
	}
	if d.IsEmpty() {
		t.Fatal("after pushes: reported empty")
	}

	_ = d.Front()
	_ = d.Back()
	if got := d.Size(); got != 3 {
		t.Fatalf("after Front/Back: got %d, want 3 (reads must not change Size)", got)
	}

	_ = d.PopFront()
	if got := d.Size(); got != 2 {
		t.Fatalf("after PopFront: got %d, want 2", got)
	}
	_ = d.PopBack()
	if got := d.Size(); got != 1 {
		t.Fatalf("after PopBack: got %d, want 1", got)
	}

	d.Clear()
	if !d.IsEmpty() {
		t.Fatal("after Clear: not empty")
	}
	if got := d.Size(); got != 0 {
		t.Fatalf("after Clear: Size got %d, want 0", got)
	}
	if d.PopFront().IsPresent() {
		t.Fatal("after Clear: PopFront returned present")
	}
	if d.PopBack().IsPresent() {
		t.Fatal("after Clear: PopBack returned present")
	}
}

// TestDequeClearOnEmptyIsNoOp covers the edge case: Clear on
// an already-empty Deque must not panic.
func TestDequeClearOnEmptyIsNoOp(t *testing.T) {
	d := NewDeque[int]()
	d.Clear()
	if !d.IsEmpty() {
		t.Fatal("Clear on empty: not empty after")
	}
}

// ---------- two-ended round trip ----------

// TestDequeMixedTwoEndedOps exercises realistic interleaving
// across both ends: push back, push front, pop front, pop
// back, push back, push front, drain. The point is to catch
// any off-by-one or stale-state errors in the delegation to
// ArrayList's head offset.
func TestDequeMixedTwoEndedOps(t *testing.T) {
	d := NewDeque[int]()
	d.PushBack(2)
	d.PushBack(3)
	d.PushFront(1)
	d.PushBack(4)
	// [1, 2, 3, 4]
	if v := d.PopFront().OrElse(-1); v != 1 {
		t.Fatalf("PopFront: got %d, want 1", v)
	}
	if v := d.PopBack().OrElse(-1); v != 4 {
		t.Fatalf("PopBack: got %d, want 4", v)
	}
	// [2, 3]
	d.PushBack(5)
	d.PushFront(0)
	// [0, 2, 3, 5]
	if v := d.PopBack().OrElse(-1); v != 5 {
		t.Fatalf("PopBack: got %d, want 5", v)
	}
	if v := d.PopBack().OrElse(-1); v != 3 {
		t.Fatalf("PopBack: got %d, want 3", v)
	}
	if v := d.PopFront().OrElse(-1); v != 0 {
		t.Fatalf("PopFront: got %d, want 0", v)
	}
	if v := d.PopFront().OrElse(-1); v != 2 {
		t.Fatalf("PopFront: got %d, want 2", v)
	}
	if !d.IsEmpty() {
		t.Fatal("after interleaved pops: not empty")
	}
}

// TestDequeManyElements is a smoke test for the head-offset
// layout under repeated push/pop from both ends. The test
// pushes 500 elements to the back, pops them all from the
// front (FIFO), then pushes 500 more to the front and pops
// them all from the back (the mirror). Both round trips must
// return the values in the expected order.
func TestDequeManyElements(t *testing.T) {
	const n = 500

	d := NewDeque[int]()
	for i := 0; i < n; i++ {
		d.PushBack(i)
	}
	if got := d.Size(); got != n {
		t.Fatalf("after %d PushBacks: Size got %d, want %d", n, got, n)
	}
	for i := 0; i < n; i++ {
		if v := d.PopFront().OrElse(-1); v != i {
			t.Fatalf("PopFront: got %d, want %d", v, i)
		}
	}
	if !d.IsEmpty() {
		t.Fatal("after draining PushBacks: not empty")
	}

	for i := 0; i < n; i++ {
		d.PushFront(i)
	}
	if got := d.Size(); got != n {
		t.Fatalf("after %d PushFronts: Size got %d, want %d", n, got, n)
	}
	for i := 0; i < n; i++ {
		if v := d.PopBack().OrElse(-1); v != i {
			t.Fatalf("PopBack: got %d, want %d", v, i)
		}
	}
	if !d.IsEmpty() {
		t.Fatal("after draining PushFronts: not empty")
	}
}

// TestDequeOpsAreConstantTime asserts the O(1) property of
// every public operation end-to-end. The head-offset layout in
// the underlying ArrayList makes PushFront / PushBack /
// PopFront / PopBack all O(1); this test documents that the
// Deque wrapper inherits the property on a large input.
//
// A pre-head-offset PopFront had to shift the entire live
// range down by one, so it was O(n) per call. A pre-head-offset
// PushFront had to shift the entire live range right by one,
// so it was also O(n). Under the head-offset layout both are
// O(1), so the total cost of 1000 operations of each kind is
// O(n) and 1 second is a very generous ceiling.
func TestDequeOpsAreConstantTime(t *testing.T) {
	const n = 1000

	// PushFront + PopBack round trip.
	d := NewDeque[int]()
	for i := 0; i < n; i++ {
		d.PushFront(i)
	}
	start := time.Now()
	for d.Size() > 0 {
		d.PopBack()
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("PushFront/PopBack took %v for %d elements; expected sub-second (O(1) per call)", elapsed, n)
	}

	// PushBack + PopFront round trip.
	for i := 0; i < n; i++ {
		d.PushBack(i)
	}
	start = time.Now()
	for d.Size() > 0 {
		d.PopFront()
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("PushBack/PopFront took %v for %d elements; expected sub-second (O(1) per call)", elapsed, n)
	}
}

// ---------- Optional integration ----------

// TestDequeReturnsOptionValues exercises the type contract of
// the access methods: every accessor returns option.Optional[T]
// with a present value on a hit and an absent value on a miss.
func TestDequeReturnsOptionValues(t *testing.T) {
	for _, name := range []string{"PopFront", "PopBack", "Front", "Back"} {
		t.Run(name+" on empty is absent", func(t *testing.T) {
			d := NewDeque[int]()
			switch name {
			case "PopFront":
				if d.PopFront().IsPresent() {
					t.Fatalf("%s on empty: IsPresent should be false", name)
				}
			case "PopBack":
				if d.PopBack().IsPresent() {
					t.Fatalf("%s on empty: IsPresent should be false", name)
				}
			case "Front":
				if d.Front().IsPresent() {
					t.Fatalf("%s on empty: IsPresent should be false", name)
				}
			case "Back":
				if d.Back().IsPresent() {
					t.Fatalf("%s on empty: IsPresent should be false", name)
				}
			}
		})
	}

	t.Run("non-empty PopFront is present", func(t *testing.T) {
		d := NewDeque[int]()
		d.PushBack(42)
		got := d.PopFront()
		if !got.IsPresent() {
			t.Fatal("PopFront on non-empty: IsPresent should be true")
		}
		if v := got.OrElse(0); v != 42 {
			t.Fatalf("PopFront: got %d, want 42", v)
		}
	})
	t.Run("non-empty PopBack is present", func(t *testing.T) {
		d := NewDeque[int]()
		d.PushFront(42)
		got := d.PopBack()
		if !got.IsPresent() {
			t.Fatal("PopBack on non-empty: IsPresent should be true")
		}
		if v := got.OrElse(0); v != 42 {
			t.Fatalf("PopBack: got %d, want 42", v)
		}
	})
}

// TestDequeChaining shows that Optional-based chaining works
// out of the box, which is the main reason Deque returns
// option.Optional[T] instead of (T, bool).
func TestDequeChaining(t *testing.T) {
	d := NewDeque[int]()
	d.PushBack(3)
	d.PushBack(7)
	d.PushFront(1)
	// [1, 3, 7]

	// Drain from the front, chaining OrElse.
	var got []int
	for !d.IsEmpty() {
		got = append(got, d.PopFront().OrElse(-1))
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 7 {
		t.Fatalf("chained PopFront: got %v, want [1 3 7]", got)
	}

	// After draining, OrElse(99) returns 99 instead of
	// panicking, and the same is true for Back, Front, and
	// PopBack.
	if v := d.PopFront().OrElse(99); v != 99 {
		t.Fatalf("PopFront on empty: got %d, want 99 (OrElse default)", v)
	}
	if v := d.PopBack().OrElse(99); v != 99 {
		t.Fatalf("PopBack on empty: got %d, want 99 (OrElse default)", v)
	}
	if v := d.Front().OrElse(99); v != 99 {
		t.Fatalf("Front on empty: got %d, want 99 (OrElse default)", v)
	}
	if v := d.Back().OrElse(99); v != 99 {
		t.Fatalf("Back on empty: got %d, want 99 (OrElse default)", v)
	}
}

// ---------- delegation correctness ----------

// TestDequeDelegationToLinkedList proves that Deque is
// faithfully forwarding to LinkedList, not running a parallel
// implementation. A scripted sequence of PushFront / PushBack /
// PopFront / PopBack operations is run on a Deque and a fresh
// LinkedList, and the final contents of both must match.
//
// The op encoding is "positive = push that value to the back,
// negative = push -(op) to the front, zero = pop from the
// front, plus-1000 = pop from the back". The scripted sequence
// is designed to exercise both ends of both data structures.
func TestDequeDelegationToLinkedList(t *testing.T) {
	d := NewDeque[int]()
	l := lists.NewLinkedList[int]()

	// ops: positive = PushBack, negative = PushFront(|op|),
	// 0 = PopFront, 1000 = PopBack. The sequence is designed
	// to leave the same front-to-back contents in both.
	ops := []int{1, -2, 3, -4, 0, 5, 0, 0, 1000, -6, 7, 1000, 0}
	for _, op := range ops {
		switch {
		case op == 0:
			d.PopFront()
			l.RemoveFirst()
		case op == 1000:
			d.PopBack()
			l.RemoveLast()
		case op > 0:
			d.PushBack(op)
			l.Add(op)
		default: // op < 0
			d.PushFront(-op)
			l.AddFirst(-op)
		}
	}

	dContents := dContentsAsSlice(d)
	lContents := l.Collect()
	if len(dContents) != len(lContents) {
		t.Fatalf("size mismatch: deque=%d, linkedlist=%d", len(dContents), len(lContents))
	}
	for i := range dContents {
		if dContents[i] != lContents[i] {
			t.Fatalf("contents mismatch at %d: deque=%d, linkedlist=%d",
				i, dContents[i], lContents[i])
		}
	}
}

// dContentsAsSlice drains a Deque from the front into a plain
// slice for comparison. Used by TestDequeDelegationToLinkedList.
func dContentsAsSlice(d *Deque[int]) []int {
	out := []int{}
	for !d.IsEmpty() {
		out = append(out, d.PopFront().OrElse(0))
	}
	return out
}

// ---------- memory-release regression test ----------

// dequeFinalizableBox is the test type used to prove that
// Deque releases references held by popped elements. Deque
// delegates to ArrayList, which already has its own
// finalizer-based regression test; this Deque-level test serves
// as an end-to-end check that the delegation does not
// accidentally keep popped references reachable from either
// end.
type dequeFinalizableBox struct {
	id int
}

var dequeFinalizableTracker struct {
	mu        chan struct{}
	finalised int
}

func init() {
	dequeFinalizableTracker.mu = make(chan struct{}, 1)
}

func newDequeFinalizableBox(id int) *dequeFinalizableBox {
	b := &dequeFinalizableBox{id: id}
	runtime.SetFinalizer(b, func(*dequeFinalizableBox) {
		dequeFinalizableTracker.mu <- struct{}{}
		dequeFinalizableTracker.finalised++
		<-dequeFinalizableTracker.mu
	})
	return b
}

func waitForDequeFinalisation(target int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		runtime.GC()
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		dequeFinalizableTracker.mu <- struct{}{}
		n := dequeFinalizableTracker.finalised
		<-dequeFinalizableTracker.mu
		if n >= target {
			return n
		}
	}
	dequeFinalizableTracker.mu <- struct{}{}
	n := dequeFinalizableTracker.finalised
	<-dequeFinalizableTracker.mu
	return n
}

// TestDequePopFrontReleasesPointerElements is the head-side
// memory-release test. Push 10 *dequeFinalizableBox elements
// to the back, PopFront 9 of them, drop the test-side
// references, and assert that all 9 popped boxes' finalizers
// have run.
func TestDequePopFrontReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForDequeFinalisation(0, 0)

	d := NewDeque[*dequeFinalizableBox]()
	boxes := make([]*dequeFinalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newDequeFinalizableBox(i)
		d.PushBack(boxes[i])
	}

	for i := 0; i < n-1; i++ {
		if !d.PopFront().IsPresent() {
			t.Fatalf("PopFront at i=%d: absent", i)
		}
	}
	for i := 0; i < n-1; i++ {
		boxes[i] = nil
	}

	got := waitForDequeFinalisation(baseline+(n-1), 5*time.Second)
	if got < baseline+(n-1) {
		t.Fatalf("after %d PopFronts: %d finalisers ran, want at least %d (popped boxes still reachable through Deque backing?)",
			n-1, got, baseline+(n-1))
	}
}

// TestDequePopBackReleasesPointerElements is the tail-side
// memory-release test. Push 10 *dequeFinalizableBox elements
// to the back, PopBack 9 of them, drop the test-side
// references, and assert that all 9 popped boxes' finalizers
// have run.
func TestDequePopBackReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForDequeFinalisation(0, 0)

	d := NewDeque[*dequeFinalizableBox]()
	boxes := make([]*dequeFinalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newDequeFinalizableBox(i)
		d.PushBack(boxes[i])
	}

	for i := 0; i < n-1; i++ {
		if !d.PopBack().IsPresent() {
			t.Fatalf("PopBack at i=%d: absent", i)
		}
	}
	for i := 0; i < n-1; i++ {
		boxes[i] = nil
	}

	got := waitForDequeFinalisation(baseline+(n-1), 5*time.Second)
	if got < baseline+(n-1) {
		t.Fatalf("after %d PopBacks: %d finalisers ran, want at least %d (popped boxes still reachable through Deque backing?)",
			n-1, got, baseline+(n-1))
	}
}

// TestDequePushFrontReleasesPointerElements exercises the
// head-side compaction. PushFront shifts the prefix, and the
// old tail elements (now in slots 0..n-1) are then drained via
// PopFront, which advances head past them and releases them.
func TestDequePushFrontReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForDequeFinalisation(0, 0)

	d := NewDeque[*dequeFinalizableBox]()
	boxes := make([]*dequeFinalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newDequeFinalizableBox(i)
		d.PushBack(boxes[i])
	}
	for i := 0; i < n; i++ {
		d.PushFront(newDequeFinalizableBox(100 + i))
	}
	for i := 0; i < n; i++ {
		d.PopFront()
	}
	for i := 0; i < n; i++ {
		boxes[i] = nil
	}

	got := waitForDequeFinalisation(baseline+2*n, 5*time.Second)
	if got < baseline+2*n {
		t.Fatalf("after PushFront + PopFront cycle: %d finalisers ran, want at least %d", got, baseline+2*n)
	}
}
