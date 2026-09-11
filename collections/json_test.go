package collections

import (
	"encoding/json"
	"reflect"
	"testing"
)

// ---------- Stack ----------

// TestStackMarshalBasicRoundTrip confirms that a Stack
// round-trips through JSON, preserving push order. The top
// of the stack is the last element of the JSON array.
func TestStackMarshalBasicRoundTrip(t *testing.T) {
	src := NewStack[int]()
	src.Push(1)
	src.Push(2)
	src.Push(3)

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[1,2,3]" {
		t.Fatalf("Marshal: got %s, want [1,2,3]", data)
	}

	var dst Stack[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	// Top of stack must be 3 (the last pushed).
	if v := dst.Pop().OrElse(0); v != 3 {
		t.Fatalf("Pop: got %d, want 3 (top of stack)", v)
	}
	if v := dst.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if v := dst.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1 (bottom of stack)", v)
	}
}

// TestStackMarshalEmpty confirms the empty round-trip.
func TestStackMarshalEmpty(t *testing.T) {
	src := NewStack[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst Stack[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal empty: should be empty")
	}
}

// TestStackMarshalNullInput confirms that JSON null gives
// an empty Stack.
func TestStackMarshalNullInput(t *testing.T) {
	var dst Stack[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal null: should be empty")
	}
}

// TestStackMarshalInvalidJSON confirms the error path.
func TestStackMarshalInvalidJSON(t *testing.T) {
	var dst Stack[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestStackMarshalReplacesContents confirms that Unmarshal
// replaces the previous items.
func TestStackMarshalReplacesContents(t *testing.T) {
	dst := NewStack[int]()
	dst.Push(99)
	dst.Push(98)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	// After replace, the top of the stack is 3.
	if v := dst.Pop().OrElse(0); v != 3 {
		t.Fatalf("Pop after replace: got %d, want 3", v)
	}
}

// ---------- Queue ----------

// TestQueueMarshalBasicRoundTrip confirms the FIFO round-
// trip. The first element of the JSON array is the head of
// the queue.
func TestQueueMarshalBasicRoundTrip(t *testing.T) {
	src := NewQueue[int]()
	src.Push(1)
	src.Push(2)
	src.Push(3)

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[1,2,3]" {
		t.Fatalf("Marshal: got %s, want [1,2,3]", data)
	}

	var dst Queue[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if v := dst.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop: got %d, want 1 (FIFO head)", v)
	}
	if v := dst.Pop().OrElse(0); v != 2 {
		t.Fatalf("Pop: got %d, want 2", v)
	}
	if v := dst.Pop().OrElse(0); v != 3 {
		t.Fatalf("Pop: got %d, want 3 (FIFO tail)", v)
	}
}

// TestQueueMarshalEmpty confirms the empty round-trip.
func TestQueueMarshalEmpty(t *testing.T) {
	src := NewQueue[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst Queue[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal empty: should be empty")
	}
}

// TestQueueMarshalNullInput confirms the null round-trip.
func TestQueueMarshalNullInput(t *testing.T) {
	var dst Queue[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal null: should be empty")
	}
}

// TestQueueMarshalInvalidJSON confirms the error path.
func TestQueueMarshalInvalidJSON(t *testing.T) {
	var dst Queue[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestQueueMarshalDrainedRoundTrip confirms that a Queue
// that has been drained from the front still marshals to
// just the live range. A Queue that was Pushed 100 times
// and Popped 100 times must marshal to "[]" and round-trip
// back to an empty Queue.
func TestQueueMarshalDrainedRoundTrip(t *testing.T) {
	src := NewQueue[int]()
	for i := 0; i < 100; i++ {
		src.Push(i)
	}
	for i := 0; i < 100; i++ {
		_ = src.Pop()
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal drained: got %s, want []", data)
	}

	var dst Queue[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal drained: should be empty")
	}
}

// TestQueueMarshalReplacesContents confirms that Unmarshal
// on an existing Queue replaces the items. After replace,
// the head must be the new first element, not a leftover.
func TestQueueMarshalReplacesContents(t *testing.T) {
	dst := NewQueue[int]()
	dst.Push(99)
	dst.Push(98)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if v := dst.Pop().OrElse(0); v != 1 {
		t.Fatalf("Pop after replace: got %d, want 1 (FIFO head)", v)
	}
}

// TestQueueMarshalDelegatesToArrayList confirms that Queue's
// MarshalJSON delegates correctly to the embedded ArrayList.
// This is a regression test for a bug where passing
// q.items (a value) instead of &q.items (a pointer) caused
// json to marshal the struct fields instead of invoking
// ArrayList's MarshalJSON method, producing "{}".
func TestQueueMarshalDelegatesToArrayList(t *testing.T) {
	src := NewQueue[int]()
	src.Push(1)
	src.Push(2)
	src.Push(3)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) == "{}" {
		t.Fatalf("Marshal produced %q: ArrayList.MarshalJSON was not invoked (delegation bug)", data)
	}
	if string(data) != "[1,2,3]" {
		t.Fatalf("Marshal: got %s, want [1,2,3]", data)
	}
}

// ---------- Deque ----------

// TestDequeMarshalBasicRoundTrip confirms the two-ended
// round-trip. The first element of the JSON array is the
// front of the Deque.
func TestDequeMarshalBasicRoundTrip(t *testing.T) {
	src := NewDeque[int]()
	src.PushBack(1)
	src.PushFront(0)
	src.PushBack(2)

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// Front is 0, then 1, then 2 (back).
	if string(data) != "[0,1,2]" {
		t.Fatalf("Marshal: got %s, want [0,1,2]", data)
	}

	var dst Deque[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if v := dst.PopFront().OrElse(-1); v != 0 {
		t.Fatalf("PopFront: got %d, want 0 (front)", v)
	}
	if v := dst.PopFront().OrElse(-1); v != 1 {
		t.Fatalf("PopFront: got %d, want 1", v)
	}
	if v := dst.PopBack().OrElse(-1); v != 2 {
		t.Fatalf("PopBack: got %d, want 2 (back)", v)
	}
}

// TestDequeMarshalEmpty confirms the empty round-trip.
func TestDequeMarshalEmpty(t *testing.T) {
	src := NewDeque[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst Deque[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal empty: should be empty")
	}
}

// TestDequeMarshalNullInput confirms the null round-trip.
func TestDequeMarshalNullInput(t *testing.T) {
	var dst Deque[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if !dst.IsEmpty() {
		t.Fatal("Unmarshal null: should be empty")
	}
}

// TestDequeMarshalInvalidJSON confirms the error path.
func TestDequeMarshalInvalidJSON(t *testing.T) {
	var dst Deque[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestDequeMarshalDelegatesToLinkedList confirms that Deque's
// MarshalJSON delegates correctly to the embedded
// LinkedList. Same regression concern as the Queue test.
func TestDequeMarshalDelegatesToLinkedList(t *testing.T) {
	src := NewDeque[int]()
	src.PushBack(1)
	src.PushFront(0)
	src.PushBack(2)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) == "{}" {
		t.Fatalf("Marshal produced %q: LinkedList.MarshalJSON was not invoked (delegation bug)", data)
	}
	if string(data) != "[0,1,2]" {
		t.Fatalf("Marshal: got %s, want [0,1,2]", data)
	}
}

// TestDequeMarshalReplacesContents confirms that Unmarshal
// replaces the previous node chain.
func TestDequeMarshalReplacesContents(t *testing.T) {
	dst := NewDeque[int]()
	dst.PushBack(99)
	dst.PushBack(98)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	// After replace, front is 1, back is 3.
	if v := dst.Front().OrElse(0); v != 1 {
		t.Fatalf("Front after replace: got %d, want 1", v)
	}
	if v := dst.Back().OrElse(0); v != 3 {
		t.Fatalf("Back after replace: got %d, want 3", v)
	}
	// Drain the Deque to confirm the underlying node chain
	// is in the expected front-to-back order.
	got := []int{}
	for !dst.IsEmpty() {
		got = append(got, dst.PopFront().OrElse(0))
	}
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Drain: got %v, want [1 2 3]", got)
	}
}
