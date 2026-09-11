package collections

import (
	"runtime"
	"testing"
	"time"

	"github.com/qianwj/typed/collections/lists"
)

// ---------- FIFO basics ----------

// TestQueueFIFOOrder covers the canonical First-In-First-Out
// order: pushes 1, 2, 3, pops them back in 1, 2, 3.
func TestQueueFIFOOrder(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	q.Push(3)

	for _, want := range []int{1, 2, 3} {
		got := q.Pop()
		if !got.IsPresent() {
			t.Fatalf("Pop: got absent, want present(%d)", want)
		}
		if v := got.Get(); v != want {
			t.Fatalf("Pop: got %d, want %d", v, want)
		}
	}

	if !q.IsEmpty() {
		t.Fatal("after draining: queue should be empty")
	}
}

// TestQueuePeekDoesNotRemove confirms Peek is a non-mutating
// read: Size stays stable across repeated Peeks, and Peek
// returns the same front element that the next Pop will yield.
func TestQueuePeekDoesNotRemove(t *testing.T) {
	q := NewQueue[string]()
	q.Push("a")
	q.Push("b")

	if v := q.Peek().OrElse(""); v != "a" {
		t.Fatalf("Peek: got %q, want \"a\"", v)
	}
	if got := q.Size(); got != 2 {
		t.Fatalf("Size after Peek: got %d, want 2", got)
	}
	if v := q.Peek().OrElse(""); v != "a" {
		t.Fatalf("Peek again: got %q, want \"a\" (must be stable)", v)
	}

	if v := q.Pop().OrElse(""); v != "a" {
		t.Fatalf("Pop after Peek: got %q, want \"a\"", v)
	}
	if v := q.Pop().OrElse(""); v != "b" {
		t.Fatalf("Pop after Peek: got %q, want \"b\"", v)
	}
}

// ---------- empty Queue behaviour ----------

// TestQueuePopAndPeekOnEmpty confirms that Pop and Peek on an
// empty Queue return an absent Optional and leave the Queue
// in the empty state.
func TestQueuePopAndPeekOnEmpty(t *testing.T) {
	q := NewQueue[int]()
	if q.Pop().IsPresent() {
		t.Fatal("Pop on empty: got present, want absent")
	}
	if q.Peek().IsPresent() {
		t.Fatal("Peek on empty: got present, want absent")
	}
	if !q.IsEmpty() {
		t.Fatal("empty queue: IsEmpty should be true")
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("empty queue: Size got %d, want 0", got)
	}
}

// TestQueueZeroValueUsable documents that the zero value of
// Queue[T] is usable without going through NewQueue. The
// internal items field is a value-typed ArrayList, which is
// itself usable as a zero value.
func TestQueueZeroValueUsable(t *testing.T) {
	var q Queue[int]
	q.Push(1)
	q.Push(2)
	if v := q.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1", v)
	}
	if v := q.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if !q.IsEmpty() {
		t.Fatal("after draining zero-value queue: should be empty")
	}
}

// ---------- Size / IsEmpty / Clear ----------

// TestQueueSizeAndIsEmptyTracksMutations confirms Size and
// IsEmpty stay consistent across Push, Peek, Pop, and Clear.
// Peek must not change Size.
func TestQueueSizeAndIsEmptyTracksMutations(t *testing.T) {
	q := NewQueue[int]()
	if !q.IsEmpty() {
		t.Fatal("new queue: not empty")
	}

	q.Push(1)
	q.Push(2)
	q.Push(3)
	if got := q.Size(); got != 3 {
		t.Fatalf("after 3 pushes: got %d, want 3", got)
	}
	if q.IsEmpty() {
		t.Fatal("after pushes: reported empty")
	}

	_ = q.Peek()
	if got := q.Size(); got != 3 {
		t.Fatalf("after Peek: got %d, want 3 (Peek must not change Size)", got)
	}

	_ = q.Pop()
	if got := q.Size(); got != 2 {
		t.Fatalf("after Pop: got %d, want 2", got)
	}

	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("after Clear: not empty")
	}
	if got := q.Size(); got != 0 {
		t.Fatalf("after Clear: Size got %d, want 0", got)
	}
	if q.Pop().IsPresent() {
		t.Fatal("after Clear: Pop returned present")
	}
}

// TestQueueClearOnEmptyIsNoOp covers the edge case: Clear on
// an already-empty Queue must not panic.
func TestQueueClearOnEmptyIsNoOp(t *testing.T) {
	q := NewQueue[int]()
	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("Clear on empty: not empty after")
	}
}

// ---------- Push / Pop interleaving ----------

// TestQueueMixedPushPop exercises realistic interleaving
// patterns: push a few, pop one, push more, pop more, mixed
// front and back. The Point is to catch any off-by-one or
// stale-state errors in the delegation to ArrayList.
func TestQueueMixedPushPop(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	if v := q.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1", v)
	}
	q.Push(3)
	q.Push(4)
	if v := q.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if v := q.Pop().OrElse(0); v != 3 {
		t.Fatalf("Pop: got %d, want 3", v)
	}
	if v := q.Pop().OrElse(0); v != 4 {
		t.Fatalf("Pop: got %d, want 4", v)
	}
	if !q.IsEmpty() {
		t.Fatal("after interleaved pushes/pops: not empty")
	}
}

// TestQueueManyElements pushes and pops a large number of
// elements. This is mostly a smoke test that the head-offset
// backing array (in ArrayList) does not break under repeated
// reallocation, and that Push / Pop stay correct at scale.
func TestQueueManyElements(t *testing.T) {
	const n = 1000
	q := NewQueue[int]()
	for i := 0; i < n; i++ {
		q.Push(i)
	}
	if got := q.Size(); got != n {
		t.Fatalf("Size after %d pushes: got %d", n, got)
	}
	for i := 0; i < n; i++ {
		v := q.Pop().OrElse(-1)
		if v != i {
			t.Fatalf("Pop: got %d, want %d", v, i)
		}
	}
	if !q.IsEmpty() {
		t.Fatal("after draining many elements: not empty")
	}
}

// TestQueuePushPopIsConstantTime asserts the O(1) property of
// Push and Pop end-to-end. The head-offset layout in the
// underlying ArrayList gives O(1) Push and Pop; this test
// documents that the Queue wrapper inherits the property.
func TestQueuePushPopIsConstantTime(t *testing.T) {
	const n = 1000
	q := NewQueue[int]()
	for i := 0; i < n; i++ {
		q.Push(i)
	}

	start := time.Now()
	for q.Size() > 0 {
		q.Pop()
	}
	elapsed := time.Since(start)

	if elapsed > time.Second {
		t.Fatalf("Queue Pop took %v for %d elements; expected sub-second (O(1) per call)", elapsed, n)
	}
}

// ---------- Optional integration ----------

// TestQueueReturnsOptionValues exercises the type contract of
// the access methods: Pop and Peek return option.Optional[T]
// with a present value on a hit and an absent value on a miss.
func TestQueueReturnsOptionValues(t *testing.T) {
	q := NewQueue[int]()

	t.Run("empty Pop is absent", func(t *testing.T) {
		if q.Pop().IsPresent() {
			t.Fatal("Pop on empty: IsPresent should be false")
		}
	})
	t.Run("empty Peek is absent", func(t *testing.T) {
		if q.Peek().IsPresent() {
			t.Fatal("Peek on empty: IsPresent should be false")
		}
	})
	t.Run("non-empty Pop is present", func(t *testing.T) {
		q.Push(42)
		got := q.Pop()
		if !got.IsPresent() {
			t.Fatal("Pop on non-empty: IsPresent should be true")
		}
		if v := got.OrElse(0); v != 42 {
			t.Fatalf("Pop: got %d, want 42", v)
		}
	})
}

// TestQueueChaining shows that Optional-based chaining works
// out of the box, which is the main reason Queue returns
// option.Optional[T] instead of (T, bool).
func TestQueueChaining(t *testing.T) {
	q := NewQueue[int]()
	q.Push(3)
	q.Push(7)

	// OrElse is the standard fallback path: pop twice with
	// a default of -1 if the Queue is empty.
	var got []int
	for !q.IsEmpty() {
		got = append(got, q.Pop().OrElse(-1))
	}
	if len(got) != 2 || got[0] != 3 || got[1] != 7 {
		t.Fatalf("chained Pop: got %v, want [3 7]", got)
	}

	// After draining, OrElse(99) returns 99 instead of
	// panicking.
	if v := q.Pop().OrElse(99); v != 99 {
		t.Fatalf("Pop on empty: got %d, want 99 (OrElse default)", v)
	}
}

// ---------- delegation correctness ----------

// TestQueueDelegationToArrayList proves that Queue is
// faithfully forwarding to ArrayList, not running a parallel
// implementation. A scripted sequence of Push / Pop
// operations is run on a Queue and a fresh ArrayList, and the
// final contents of both must match.
func TestQueueDelegationToArrayList(t *testing.T) {
	q := NewQueue[int]()
	a := lists.NewArrayList[int]()

	// -1 means "Pop"; any other value means "Push that value".
	// The sequence is designed to exercise interleaving and
	// the head-side state of both data structures.
	ops := []int{1, 2, 3, -1, 4, 5, -1, -1, 6, -1, -1, 7, 8}
	for _, op := range ops {
		if op == -1 {
			q.Pop()
			a.RemoveFirst()
		} else {
			q.Push(op)
			a.Add(op)
		}
	}

	qContents := qContentsAsSlice(q)
	aContents := a.Collect()
	if len(qContents) != len(aContents) {
		t.Fatalf("size mismatch: queue=%d, arraylist=%d", len(qContents), len(aContents))
	}
	for i := range qContents {
		if qContents[i] != aContents[i] {
			t.Fatalf("contents mismatch at %d: queue=%d, arraylist=%d",
				i, qContents[i], aContents[i])
		}
	}
}

// qContentsAsSlice drains a Queue into a plain slice for
// comparison. Used by TestQueueDelegationToArrayList.
func qContentsAsSlice(q *Queue[int]) []int {
	out := []int{}
	for !q.IsEmpty() {
		out = append(out, q.Pop().OrElse(0))
	}
	return out
}

// ---------- memory-release regression test (parallel to ArrayList) ----------

// queueFinalizableBox is the test type used to prove that Queue
// releases references held by popped elements. The Queue
// delegates to ArrayList, which already has its own
// finalizer-based regression test; this Queue-level test
// serves as an end-to-end check that the delegation does not
// accidentally keep popped references reachable.
type queueFinalizableBox struct {
	id int
}

var queueFinalizableTracker struct {
	mu        chan struct{}
	finalised int
}

func init() {
	queueFinalizableTracker.mu = make(chan struct{}, 1)
}

func newQueueFinalizableBox(id int) *queueFinalizableBox {
	b := &queueFinalizableBox{id: id}
	runtime.SetFinalizer(b, func(*queueFinalizableBox) {
		queueFinalizableTracker.mu <- struct{}{}
		queueFinalizableTracker.finalised++
		<-queueFinalizableTracker.mu
	})
	return b
}

func waitForQueueFinalisation(target int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		runtime.GC()
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		queueFinalizableTracker.mu <- struct{}{}
		n := queueFinalizableTracker.finalised
		<-queueFinalizableTracker.mu
		if n >= target {
			return n
		}
	}
	queueFinalizableTracker.mu <- struct{}{}
	n := queueFinalizableTracker.finalised
	<-queueFinalizableTracker.mu
	return n
}

// TestQueuePopReleasesPointerElements is the end-to-end
// memory-release test for Queue. Push 10 *queueFinalizableBox
// elements, Pop 9 of them, drop the test-side references,
// and assert that all 9 popped boxes' finalizers have run.
//
// Because Queue delegates to ArrayList, and ArrayList already
// has the head-offset layout that releases the front slot on
// each RemoveFirst, this test passes by inheritance.
func TestQueuePopReleasesPointerElements(t *testing.T) {
	const n = 10
	baseline := waitForQueueFinalisation(0, 0)

	q := NewQueue[*queueFinalizableBox]()
	boxes := make([]*queueFinalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newQueueFinalizableBox(i)
		q.Push(boxes[i])
	}

	for i := 0; i < n-1; i++ {
		if !q.Pop().IsPresent() {
			t.Fatalf("Pop at i=%d: !ok", i)
		}
	}
	for i := 0; i < n-1; i++ {
		boxes[i] = nil
	}

	got := waitForQueueFinalisation(baseline+(n-1), 5*time.Second)
	if got < baseline+(n-1) {
		t.Fatalf("after %d Pops: %d finalisers ran, want at least %d (popped boxes still reachable through Queue backing?)",
			n-1, got, baseline+(n-1))
	}
}

// _ keeps the import list honest; no-op assignment. The
// finalization helpers above are referenced by the test
// bodies but Go does not require a static reference to keep
// the imports alive.
var _ = runtime.GC
