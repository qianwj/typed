package sets_test

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/qianwj/typed/collections/sets"
)

// TestHashSetMarshalBasicRoundTrip confirms the canonical
// happy path: build a HashSet, marshal, unmarshal into a
// fresh HashSet, and confirm the set's contents match.
//
// Because map iteration order is non-deterministic, the test
// asserts set equality (sorted slice comparison) rather than
// raw byte equality.
func TestHashSetMarshalBasicRoundTrip(t *testing.T) {
	src := sets.NewHashSet[int]()
	src.Add(1)
	src.Add(2)
	src.Add(3)
	src.Add(4)
	src.Add(5)

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst sets.HashSet[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 5 {
		t.Fatalf("Unmarshal: size = %d, want 5", got)
	}
	got := dst.Collect()
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Unmarshal: got %v, want %v", got, want)
		}
	}
}

// TestHashSetMarshalEmpty confirms that an empty HashSet
// marshals to "[]" and round-trips back to an empty HashSet.
func TestHashSetMarshalEmpty(t *testing.T) {
	src := sets.NewHashSet[int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("Marshal empty: got %s, want []", data)
	}

	var dst sets.HashSet[int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal empty: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal empty: size = %d, want 0", got)
	}
}

// TestHashSetMarshalNullInput confirms that JSON null
// produces an empty HashSet.
func TestHashSetMarshalNullInput(t *testing.T) {
	var dst sets.HashSet[int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal null: size = %d, want 0", got)
	}
}

// TestHashSetMarshalInvalidJSON confirms the error path.
func TestHashSetMarshalInvalidJSON(t *testing.T) {
	var dst sets.HashSet[int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestHashSetMarshalTypeMismatch confirms that a JSON object
// (instead of an array) returns a type error on unmarshal.
func TestHashSetMarshalTypeMismatch(t *testing.T) {
	var dst sets.HashSet[int]
	if err := json.Unmarshal([]byte(`{"key": "value"}`), &dst); err == nil {
		t.Fatal("Unmarshal type-mismatch: expected error, got nil")
	}
}

// TestHashSetMarshalDuplicatesCollapsed confirms that a JSON
// array with duplicate values produces a set with one entry
// per distinct value. This is the "set" semantic; the
// underlying collection type does not preserve multiplicity.
func TestHashSetMarshalDuplicatesCollapsed(t *testing.T) {
	var dst sets.HashSet[int]
	if err := json.Unmarshal([]byte("[1,1,2,2,2,3,3,3,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 3 {
		t.Fatalf("Unmarshal duplicates: size = %d, want 3", got)
	}
	got := dst.Collect()
	sort.Ints(got)
	want := []int{1, 2, 3}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Unmarshal: got %v, want %v", got, want)
		}
	}
}

// TestHashSetMarshalStructElements confirms that struct
// elements round-trip.
func TestHashSetMarshalStructElements(t *testing.T) {
	type point struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	src := sets.NewHashSet[point]()
	src.Add(point{X: 1, Y: 2})
	src.Add(point{X: 3, Y: 4})

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst sets.HashSet[point]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 2 {
		t.Fatalf("Unmarshal: size = %d, want 2", got)
	}
	if !dst.Contains(point{X: 1, Y: 2}) {
		t.Fatal("Unmarshal: missing {1 2}")
	}
	if !dst.Contains(point{X: 3, Y: 4}) {
		t.Fatal("Unmarshal: missing {3 4}")
	}
}

// TestHashSetMarshalReplacesContents confirms that Unmarshal
// on an existing HashSet replaces the contents.
func TestHashSetMarshalReplacesContents(t *testing.T) {
	dst := sets.NewHashSet[int]()
	dst.Add(99)
	dst.Add(98)
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 3 {
		t.Fatalf("Unmarshal: size = %d, want 3 (99 and 98 should be gone)", got)
	}
	if dst.Contains(99) {
		t.Fatal("Unmarshal: 99 still present after replace")
	}
	if !dst.Contains(1) || !dst.Contains(2) || !dst.Contains(3) {
		t.Fatalf("Unmarshal: missing some of 1/2/3: %v", dst.Collect())
	}
}

// TestHashSetMarshalMatchesCollect confirms that the JSON
// output, after a round-trip, produces a set with the same
// elements. Byte equality cannot be asserted directly because
// the underlying map's iteration order is non-deterministic,
// but set equality (as a sorted slice) is stable.
func TestHashSetMarshalMatchesCollect(t *testing.T) {
	src := sets.NewHashSet[int]()
	src.Add(1)
	src.Add(2)
	src.Add(3)

	fromSet, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal set: %v", err)
	}
	fromCollect, err := json.Marshal(src.Collect())
	if err != nil {
		t.Fatalf("Marshal slice: %v", err)
	}

	// Decode both back into a HashSet and compare.
	var gotFromSet, gotFromCollect sets.HashSet[int]
	if err := json.Unmarshal(fromSet, &gotFromSet); err != nil {
		t.Fatalf("Unmarshal set: %v", err)
	}
	if err := json.Unmarshal(fromCollect, &gotFromCollect); err != nil {
		t.Fatalf("Unmarshal slice: %v", err)
	}
	if gotFromSet.Size() != gotFromCollect.Size() {
		t.Fatalf("size mismatch: set=%d, slice=%d", gotFromSet.Size(), gotFromCollect.Size())
	}
	gotFromSet.ForEach(func(v int) {
		if !gotFromCollect.Contains(v) {
			t.Fatalf("set has %d but slice does not", v)
		}
	})
}
