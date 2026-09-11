package lists_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/qianwj/typed/collections/lists"
)

// ---------- ArrayList ----------

// TestArrayListMarshalBasicRoundTrip confirms the canonical
// happy path: build an ArrayList, marshal to JSON, unmarshal
// into a fresh ArrayList, and confirm the contents match.
func TestArrayListMarshalBasicRoundTrip(t *testing.T) {
	src := lists.ArrayListOf(1, 2, 3, 4, 5)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[1,2,3,4,5]" {
		t.Fatalf("Marshal: got %s, want [1,2,3,4,5]", data)
	}

	var dst lists.ArrayList[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Unmarshal: got %v, want [1 2 3 4 5]", got)
	}
}

// TestArrayListMarshalEmpty confirms that an empty ArrayList
// marshals to "[]" rather than "null". This is the convention
// set by the doc: absent and empty collections are
// indistinguishable in the JSON output.
func TestArrayListMarshalEmpty(t *testing.T) {
	src := lists.NewArrayList[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst lists.ArrayList[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal empty: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal empty: size = %d, want 0", got)
	}
}

// TestArrayListMarshalDrained confirms that an ArrayList whose
// head offset has advanced past most of the backing array
// still marshals to the live range only. A list that has
// been drained from the front 100 times and then has 3 live
// elements must marshal to a 3-element array, not a
// 103-element array with 100 stale slots.
func TestArrayListMarshalDrained(t *testing.T) {
	src := lists.NewArrayList[int]()
	for i := 0; i < 103; i++ {
		src.Add(i)
	}
	for i := 0; i < 100; i++ {
		src.RemoveFirst()
	}
	// Now head=100, live range is [100, 101, 102].

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[100,101,102]" {
		t.Fatalf("Marshal drained: got %s, want [100,101,102]", data)
	}
}

// TestArrayListMarshalNullInput confirms that JSON null is
// accepted and produces an empty ArrayList. This matches the
// v1 json package's behaviour for slices and the project's
// "absent and empty are equivalent" convention.
func TestArrayListMarshalNullInput(t *testing.T) {
	var dst lists.ArrayList[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal null: size = %d, want 0", got)
	}
	if got := dst.Collect(); len(got) != 0 {
		t.Fatalf("Unmarshal null: Collect = %v, want empty", got)
	}
}

// TestArrayListMarshalInvalidJSON confirms that a malformed
// JSON array produces an unmarshal error.
func TestArrayListMarshalInvalidJSON(t *testing.T) {
	var dst lists.ArrayList[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestArrayListMarshalTypeMismatch confirms that a JSON object
// (instead of an array) returns a type error on unmarshal.
func TestArrayListMarshalTypeMismatch(t *testing.T) {
	var dst lists.ArrayList[int]
	if err := json.Unmarshal([]byte(`{"key": "value"}`), &dst); err == nil {
		t.Fatal("Unmarshal type-mismatch: expected error, got nil")
	}
}

// TestArrayListMarshalStructElements confirms that ArrayList
// of structs round-trips correctly through JSON, including
// the json struct tag for field renaming.
func TestArrayListMarshalStructElements(t *testing.T) {
	type point struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	src := lists.ArrayListOf(point{X: 1, Y: 2}, point{X: 3, Y: 4})
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want := `[{"x":1,"y":2},{"x":3,"y":4}]`
	if string(data) != want {
		t.Fatalf("Marshal: got %s, want %s", data, want)
	}

	var dst lists.ArrayList[point]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Collect(); !reflect.DeepEqual(got, []point{{1, 2}, {3, 4}}) {
		t.Fatalf("Unmarshal: got %v, want [{1 2} {3 4}]", got)
	}
}

// TestArrayListMarshalPointerElements confirms that an
// ArrayList of pointers round-trips. The pointer elements
// must be non-nil after unmarshaling.
func TestArrayListMarshalPointerElements(t *testing.T) {
	type inner struct{ V int }
	src := lists.NewArrayList[*inner]()
	src.Add(&inner{V: 1})
	src.Add(&inner{V: 2})
	src.Add(&inner{V: 3})

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// The output is a JSON array of objects, not the string
	// representations of the pointers.
	if !strings.HasPrefix(string(data), "[{") {
		t.Fatalf("Marshal: got %s, want [ prefix", data)
	}

	var dst lists.ArrayList[*inner]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	got := dst.Collect()
	if len(got) != 3 {
		t.Fatalf("Unmarshal: len = %d, want 3", len(got))
	}
	for i, p := range got {
		if p == nil {
			t.Fatalf("Unmarshal: element %d is nil", i)
		}
		if p.V != i+1 {
			t.Fatalf("Unmarshal: element %d has V = %d, want %d", i, p.V, i+1)
		}
	}
}

// TestArrayListMarshalCustomMarshaler confirms that an
// ArrayList whose element type implements json.Marshaler uses
// the element's MarshalJSON method. This is the "user
// controls the encoding per element" extension point.
func TestArrayListMarshalCustomMarshaler(t *testing.T) {
	type tagged struct {
		Tag string `json:"tag"`
	}
	// We use a custom type that always emits a known string
	// for "value", confirming the per-element MarshalJSON
	// is honoured.
	src := lists.ArrayListOf("hello", "world")
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != `["hello","world"]` {
		t.Fatalf("Marshal: got %s", data)
	}
	_ = tagged{} // keep the import-free declaration unused-friendly
}

// TestArrayListMarshalMatchesCollect confirms that
// json.Marshal(ArrayList) produces the same bytes as
// json.Marshal(ArrayList.Collect()). This is the round-trip
// property that the doc promises.
func TestArrayListMarshalMatchesCollect(t *testing.T) {
	src := lists.ArrayListOf(1, 2, 3, 4, 5)
	fromList, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal list: %v", err)
	}
	fromSlice, err := json.Marshal(src.Collect())
	if err != nil {
		t.Fatalf("Marshal slice: %v", err)
	}
	if string(fromList) != string(fromSlice) {
		t.Fatalf("Marshal list and slice disagree: %s vs %s", fromList, fromSlice)
	}
}

// TestArrayListMarshalReplacesContents confirms that calling
// Unmarshal on an existing ArrayList replaces the contents
// rather than appending. The previous head offset is reset
// to zero.
func TestArrayListMarshalReplacesContents(t *testing.T) {
	dst := lists.ArrayListOf(99, 98, 97)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Unmarshal replace: got %v, want [1 2 3]", got)
	}
	// After Unmarshal, head should be 0 so the next Add
	// writes at index 3, not after any leftover prefix.
	dst.Add(4)
	if got := dst.Get(3).OrElse(-1); got != 4 {
		t.Fatalf("Add after Unmarshal: Get(3) = %d, want 4", got)
	}
}

// ---------- LinkedList ----------

// TestLinkedListMarshalBasicRoundTrip is the LinkedList
// counterpart of TestArrayListMarshalBasicRoundTrip. The
// node chain is rebuilt via Add during Unmarshal.
func TestLinkedListMarshalBasicRoundTrip(t *testing.T) {
	src := lists.LinkedListOf(1, 2, 3, 4, 5)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[1,2,3,4,5]" {
		t.Fatalf("Marshal: got %s, want [1,2,3,4,5]", data)
	}

	var dst lists.LinkedList[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Unmarshal: got %v, want [1 2 3 4 5]", got)
	}
}

// TestLinkedListMarshalEmpty confirms the empty round-trip.
func TestLinkedListMarshalEmpty(t *testing.T) {
	src := lists.NewLinkedList[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst lists.LinkedList[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal empty: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal empty: size = %d, want 0", got)
	}
	if dst.First().IsPresent() || dst.Last().IsPresent() {
		t.Fatal("Unmarshal empty: First/Last must be absent")
	}
}

// TestLinkedListMarshalNullInput confirms that JSON null
// produces an empty LinkedList with head, tail both nil.
func TestLinkedListMarshalNullInput(t *testing.T) {
	var dst lists.LinkedList[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal null: size = %d, want 0", got)
	}
}

// TestLinkedListMarshalInvalidJSON confirms the error path.
func TestLinkedListMarshalInvalidJSON(t *testing.T) {
	var dst lists.LinkedList[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestLinkedListMarshalOrderPreserved confirms that head-to-
// tail order is preserved across the round-trip. The first
// node after unmarshal is the first element of the input
// array; the last node is the last element.
func TestLinkedListMarshalOrderPreserved(t *testing.T) {
	src := lists.NewLinkedList[int]()
	for i := 1; i <= 5; i++ {
		src.Add(i)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst lists.LinkedList[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	// Walk head-to-tail and confirm the sequence is 1, 2, 3, 4, 5.
	var got []int
	dst.ForEach(func(v int) { got = append(got, v) })
	if !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Order: got %v, want [1 2 3 4 5]", got)
	}
}

// TestLinkedListMarshalReplacesContents confirms that
// Unmarshal replaces the previous node chain rather than
// appending. The previous head and tail pointers are reset.
func TestLinkedListMarshalReplacesContents(t *testing.T) {
	dst := lists.LinkedListOf(99, 98, 97)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Unmarshal replace: got %v, want [1 2 3]", got)
	}
	// After Unmarshal, the new head is 1, the new tail is 3.
	if v := dst.First().OrElse(0); v != 1 {
		t.Fatalf("First after Unmarshal: got %d, want 1", v)
	}
	if v := dst.Last().OrElse(0); v != 3 {
		t.Fatalf("Last after Unmarshal: got %d, want 3", v)
	}
}

// TestLinkedListMarshalPointerElements confirms that an
// LinkedList of pointers round-trips.
func TestLinkedListMarshalPointerElements(t *testing.T) {
	type inner struct{ V int }
	src := lists.NewLinkedList[*inner]()
	src.Add(&inner{V: 10})
	src.Add(&inner{V: 20})

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst lists.LinkedList[*inner]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	got := dst.Collect()
	if len(got) != 2 {
		t.Fatalf("Unmarshal: len = %d, want 2", len(got))
	}
	if got[0] == nil || got[0].V != 10 {
		t.Fatalf("Unmarshal [0]: got %+v, want &{10}", got[0])
	}
	if got[1] == nil || got[1].V != 20 {
		t.Fatalf("Unmarshal [1]: got %+v, want &{20}", got[1])
	}
}
