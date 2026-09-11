package collections

import (
	"runtime"
	"testing"
	"time"
)

// ---------- basic LIFO behaviour ----------

// TestStackLIFO covers the canonical Last-In-First-Out order:
// pushes 1, 2, 3 in order, pops them back in 3, 2, 1.
func TestStackLIFO(t *testing.T) {
	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for _, want := range []int{3, 2, 1} {
		got := s.Pop()
		if !got.IsPresent() {
			t.Fatalf("Pop: got absent, want present(%d)", want)
		}
		if v := got.Get(); v != want {
			t.Fatalf("Pop: got %d, want %d", v, want)
		}
	}

	if !s.IsEmpty() {
		t.Fatal("after draining: stack should be empty")
	}
}

// TestStackPeekDoesNotRemove confirms Peek is a non-mutating
// read: Size stays the same across repeated Peek calls, and the
// returned value is stable.
func TestStackPeekDoesNotRemove(t *testing.T) {
	s := NewStack[string]()
	s.Push("a")
	s.Push("b")

	if v := s.Peek().OrElse(""); v != "b" {
		t.Fatalf("Peek: got %q, want \"b\"", v)
	}
	if got := s.Size(); got != 2 {
		t.Fatalf("Size after Peek: got %d, want 2", got)
	}
	if v := s.Peek().OrElse(""); v != "b" {
		t.Fatalf("Peek again: got %q, want \"b\" (must be stable)", v)
	}

	// Peek should also not interfere with a subsequent Pop.
	if v := s.Pop().OrElse(""); v != "b" {
		t.Fatalf("Pop after Peek: got %q, want \"b\"", v)
	}
	if v := s.Pop().OrElse(""); v != "a" {
		t.Fatalf("Pop after Peek: got %q, want \"a\"", v)
	}
}

// ---------- empty Stack behaviour ----------

// TestStackPopAndPeekOnEmpty confirms that Pop and Peek on an
// empty Stack return an absent Optional, and do not panic or
// corrupt the Stack.
func TestStackPopAndPeekOnEmpty(t *testing.T) {
	s := NewStack[int]()
	if s.Pop().IsPresent() {
		t.Fatal("Pop on empty: got present, want absent")
	}
	if s.Peek().IsPresent() {
		t.Fatal("Peek on empty: got present, want absent")
	}
	if !s.IsEmpty() {
		t.Fatal("empty stack: IsEmpty should be true")
	}
	if got := s.Size(); got != 0 {
		t.Fatalf("empty stack: Size got %d, want 0", got)
	}
}

// TestStackZeroValueUsable documents that the zero value of
// Stack[T] is usable without going through NewStack. The fields
// (just a slice) start as nil, but every operation handles that
// the same way as an explicitly constructed empty Stack.
func TestStackZeroValueUsable(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	if v := s.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if v := s.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1", v)
	}
	if !s.IsEmpty() {
		t.Fatal("after draining zero-value stack: should be empty")
	}
}

// ---------- Size / IsEmpty / Clear ----------

// TestStackSizeAndIsEmptyTracksMutations verifies that Size and
// IsEmpty stay consistent across pushes, pops, peeks, and Clear.
func TestStackSizeAndIsEmptyTracksMutations(t *testing.T) {
	s := NewStack[int]()
	if !s.IsEmpty() {
		t.Fatal("new stack: not empty")
	}

	s.Push(1)
	s.Push(2)
	s.Push(3)
	if got := s.Size(); got != 3 {
		t.Fatalf("after 3 pushes: got %d, want 3", got)
	}
	if s.IsEmpty() {
		t.Fatal("after pushes: reported empty")
	}

	_ = s.Peek() // Peek must not change Size
	if got := s.Size(); got != 3 {
		t.Fatalf("after Peek: got %d, want 3 (Peek must not change Size)", got)
	}

	_ = s.Pop()
	if got := s.Size(); got != 2 {
		t.Fatalf("after Pop: got %d, want 2", got)
	}

	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("after Clear: not empty")
	}
	if got := s.Size(); got != 0 {
		t.Fatalf("after Clear: Size got %d, want 0", got)
	}
	if s.Pop().IsPresent() {
		t.Fatal("after Clear: Pop returned present")
	}
}

// TestStackClearOnEmptyIsNoOp covers the edge case: Clear on an
// already-empty Stack must not panic.
func TestStackClearOnEmptyIsNoOp(t *testing.T) {
	s := NewStack[int]()
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear on empty: not empty after")
	}
}

// ---------- Push / Pop interleaving ----------

// TestStackPopAlternating exercises the realistic push / pop
// pattern where pushes and pops interleave, to catch any
// off-by-one or stale-state errors.
func TestStackPopAlternating(t *testing.T) {
	s := NewStack[int]()
	s.Push(1)
	if v := s.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1", v)
	}
	s.Push(2)
	s.Push(3)
	if v := s.Pop().OrElse(0); v != 3 {
		t.Fatalf("Pop: got %d, want 3", v)
	}
	s.Push(4)
	if v := s.Pop().OrElse(0); v != 4 {
		t.Fatalf("Pop: got %d, want 4", v)
	}
	if v := s.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if !s.IsEmpty() {
		t.Fatal("after draining interleaved pushes/pops: not empty")
	}
}

// TestStackManyElements pushes and pops a large number of
// elements. This is mostly a smoke test that the backing-array
// growth and shrink logic does not break under repeated
// reallocation.
func TestStackManyElements(t *testing.T) {
	const n = 1000
	s := NewStack[int]()
	for i := 0; i < n; i++ {
		s.Push(i)
	}
	if got := s.Size(); got != n {
		t.Fatalf("Size after %d pushes: got %d", n, got)
	}
	for i := n - 1; i >= 0; i-- {
		v := s.Pop().OrElse(-1)
		if v != i {
			t.Fatalf("Pop: got %d, want %d", v, i)
		}
	}
	if !s.IsEmpty() {
		t.Fatal("after draining many elements: not empty")
	}
}

// ---------- memory-leak fix verification ----------

// finalizableBox is a struct that registers a finalizer on
// construction. The test uses the finalizer as a witness that
// the object was actually garbage-collected. This is the
// standard idiom for "this object is unreachable" tests, used
// for example in the standard library's container tests.
type finalizableBox struct {
	id int
}

// finalizableBoxTracker counts how many finalizableBox
// instances have been finalised. Read/Write from the test only.
var finalizableBoxTracker struct {
	mu        chan struct{}
	finalised int
}

func init() {
	finalizableBoxTracker.mu = make(chan struct{}, 1)
}

// newFinalizableBox allocates a finalizableBox and registers a
// finalizer that increments the tracker when the box is
// collected.
func newFinalizableBox(id int) *finalizableBox {
	b := &finalizableBox{id: id}
	runtime.SetFinalizer(b, func(fb *finalizableBox) {
		finalizableBoxTracker.mu <- struct{}{}
		finalizableBoxTracker.finalised++
		<-finalizableBoxTracker.mu
	})
	return b
}

// waitForFinalisation spins until the finalizer tracker reaches
// at least target, or the timeout elapses. Returns the number
// actually finalised.
//
// The timeout must be generous: Go's runtime may delay
// finalization arbitrarily, especially under -race, and the
// garbage collector is non-deterministic. Empirically 1 second
// is enough for a handful of objects on a quiet system, but we
// use 5 seconds to be safe.
func waitForFinalisation(target int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		runtime.GC()
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		finalizableBoxTracker.mu <- struct{}{}
		n := finalizableBoxTracker.finalised
		<-finalizableBoxTracker.mu
		if n >= target {
			return n
		}
	}
	finalizableBoxTracker.mu <- struct{}{}
	n := finalizableBoxTracker.finalised
	<-finalizableBoxTracker.mu
	return n
}

// TestStackPopReleasesPointerElements is the regression test for
// the memory-leak fix. It pushes a sequence of *finalizableBox
// onto a Stack, then pops them all. After the pops, the popped
// boxes must become unreachable so their finalizers can run.
//
// Before the fix, the backing array kept the popped boxes alive
// (through the now-out-of-range slots), and the test would time
// out waiting for the finalizers. After the fix, the slots are
// zeroed on each Pop, releasing the references.
func TestStackPopReleasesPointerElements(t *testing.T) {
	const n = 5
	baseline := waitForFinalisation(0, 0)

	s := NewStack[*finalizableBox]()
	boxes := make([]*finalizableBox, n)
	for i := 0; i < n; i++ {
		boxes[i] = newFinalizableBox(i)
		s.Push(boxes[i])
	}

	// Pop everything. After the fix, the popped boxes are
	// released and their finalizers are eligible to run.
	for i := 0; i < n; i++ {
		if !s.Pop().IsPresent() {
			t.Fatalf("Pop: got absent at i=%d", i)
		}
	}
	if !s.IsEmpty() {
		t.Fatal("after popping all: stack should be empty")
	}

	// Drop our strong references to the boxes too, so the only
	// thing keeping them alive would be the Stack's backing
	// array.
	for i := range boxes {
		boxes[i] = nil
	}

	// Force a few GC cycles and wait for finalizers.
	got := waitForFinalisation(baseline+n, 5*time.Second)
	if got < baseline+n {
		t.Fatalf("finalisers ran %d times, want at least %d (popped boxes still reachable through backing array?)", got, baseline+n)
	}
}

// TestStackClearReleasesPointerElements is the parallel test for
// Clear. Clear re-creates the backing array, so all references
// to old elements are released and the finalizers should run.
func TestStackClearReleasesPointerElements(t *testing.T) {
	const n = 5
	baseline := waitForFinalisation(0, 0)

	s := NewStack[*finalizableBox]()
	for i := 0; i < n; i++ {
		s.Push(newFinalizableBox(i))
	}

	s.Clear()

	got := waitForFinalisation(baseline+n, 5*time.Second)
	if got < baseline+n {
		t.Fatalf("after Clear: finalisers ran %d times, want at least %d", got, baseline+n)
	}
}

// ---------- Optional integration ----------

// TestStackPopReturnsOptionValues verifies the type contract:
// Pop returns option.Optional[T] with a present value when the
// stack is non-empty and an absent value when it is empty.
func TestStackPopReturnsOptionValues(t *testing.T) {
	s := NewStack[int]()

	empty := s.Pop()
	if !empty.IsEmpty() {
		t.Fatal("Pop on empty: IsEmpty should be true")
	}
	if empty.IsPresent() {
		t.Fatal("Pop on empty: IsPresent should be false")
	}

	s.Push(42)
	present := s.Pop()
	if !present.IsPresent() {
		t.Fatal("Pop on non-empty: IsPresent should be true")
	}
	if v := present.OrElse(0); v != 42 {
		t.Fatalf("Pop: got %d, want 42", v)
	}
}
