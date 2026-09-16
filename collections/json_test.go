package collections

import (
	"encoding/json"
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
