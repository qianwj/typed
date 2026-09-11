package maps_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qianwj/typed/collections/maps"
)

// TestHashMapMarshalBasicRoundTrip confirms the canonical
// happy path: build a HashMap, marshal, unmarshal into a
// fresh HashMap, and confirm the entries match.
func TestHashMapMarshalBasicRoundTrip(t *testing.T) {
	src := maps.NewHashMap[string, int]()
	src.Put("a", 1)
	src.Put("b", 2)
	src.Put("c", 3)

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst maps.HashMap[string, int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 3 {
		t.Fatalf("Unmarshal: size = %d, want 3", got)
	}
	if v, ok := dst.Get("a"); !ok || v != 1 {
		t.Fatalf("Unmarshal: Get(a) = (%d, %v), want (1, true)", v, ok)
	}
	if v, ok := dst.Get("b"); !ok || v != 2 {
		t.Fatalf("Unmarshal: Get(b) = (%d, %v), want (2, true)", v, ok)
	}
	if v, ok := dst.Get("c"); !ok || v != 3 {
		t.Fatalf("Unmarshal: Get(c) = (%d, %v), want (3, true)", v, ok)
	}
}

// TestHashMapMarshalEmpty confirms that an empty HashMap
// marshals to "{}" and round-trips back to an empty HashMap.
func TestHashMapMarshalEmpty(t *testing.T) {
	src := maps.NewHashMap[string, int]()
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(data) != "{}" {
		t.Fatalf("Marshal empty: got %s, want {}", data)
	}

	var dst maps.HashMap[string, int]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal empty: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal empty: size = %d, want 0", got)
	}
}

// TestHashMapMarshalNullInput confirms that JSON null
// produces an empty HashMap.
func TestHashMapMarshalNullInput(t *testing.T) {
	var dst maps.HashMap[string, int]
	if err := json.Unmarshal([]byte("null"), &dst); err != nil {
		t.Fatalf("Unmarshal null: %v", err)
	}
	if got := dst.Size(); got != 0 {
		t.Fatalf("Unmarshal null: size = %d, want 0", got)
	}
}

// TestHashMapMarshalInvalidJSON confirms the error path.
func TestHashMapMarshalInvalidJSON(t *testing.T) {
	var dst maps.HashMap[string, int]
	if err := json.Unmarshal([]byte("{not json"), &dst); err == nil {
		t.Fatal("Unmarshal invalid: expected error, got nil")
	}
}

// TestHashMapMarshalTypeMismatch confirms that a JSON array
// (instead of an object) returns a type error on unmarshal.
func TestHashMapMarshalTypeMismatch(t *testing.T) {
	var dst maps.HashMap[string, int]
	if err := json.Unmarshal([]byte("[1,2,3]"), &dst); err == nil {
		t.Fatal("Unmarshal type-mismatch: expected error, got nil")
	}
}

// TestHashMapMarshalStructValues confirms that struct
// values round-trip correctly.
func TestHashMapMarshalStructValues(t *testing.T) {
	type point struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	src := maps.NewHashMap[string, point]()
	src.Put("origin", point{X: 0, Y: 0})
	src.Put("top-right", point{X: 10, Y: 20})

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var dst maps.HashMap[string, point]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 2 {
		t.Fatalf("Unmarshal: size = %d, want 2", got)
	}
	if v, ok := dst.Get("origin"); !ok || v.X != 0 || v.Y != 0 {
		t.Fatalf("Unmarshal: Get(origin) = (%+v, %v), want ({0 0}, true)", v, ok)
	}
	if v, ok := dst.Get("top-right"); !ok || v.X != 10 || v.Y != 20 {
		t.Fatalf("Unmarshal: Get(top-right) = (%+v, %v), want ({10 20}, true)", v, ok)
	}
}

// TestHashMapMarshalIntKey confirms that non-string key
// types also round-trip, as long as they are json-marshalable.
// The Marshal assertion is set-based because map iteration
// order is non-deterministic; the test parses the output back
// and checks the resulting map.
func TestHashMapMarshalIntKey(t *testing.T) {
	src := maps.NewHashMap[int, string]()
	src.Put(1, "one")
	src.Put(2, "two")
	src.Put(3, "three")

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	// Parse the output and check the entries as a set,
	// since the underlying Go map iteration order is
	// non-deterministic. The previous version of this test
	// asserted byte-equality with sorted keys, which only
	// happened to pass by chance on a particular run.
	var rawMap map[string]string
	if err := json.Unmarshal(data, &rawMap); err != nil {
		t.Fatalf("parse Marshal output: %v", err)
	}
	wantRaw := map[string]string{"1": "one", "2": "two", "3": "three"}
	if !reflect.DeepEqual(rawMap, wantRaw) {
		t.Fatalf("Marshal: got %v, want %v", rawMap, wantRaw)
	}

	var dst maps.HashMap[int, string]
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 3 {
		t.Fatalf("Unmarshal: size = %d, want 3", got)
	}
	if v, ok := dst.Get(2); !ok || v != "two" {
		t.Fatalf("Unmarshal: Get(2) = (%s, %v), want (two, true)", v, ok)
	}
}

// TestHashMapMarshalReplacesContents confirms that Unmarshal
// on an existing HashMap replaces the entries.
func TestHashMapMarshalReplacesContents(t *testing.T) {
	dst := maps.NewHashMap[string, int]()
	dst.Put("x", 99)
	dst.Put("y", 98)
	if err := json.Unmarshal([]byte(`{"a":1,"b":2}`), &dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got := dst.Size(); got != 2 {
		t.Fatalf("Unmarshal: size = %d, want 2 (x and y should be gone)", got)
	}
	if _, ok := dst.Get("x"); ok {
		t.Fatal("Unmarshal: Get(x) still present after replace")
	}
	if v, ok := dst.Get("a"); !ok || v != 1 {
		t.Fatalf("Unmarshal: Get(a) = (%d, %v), want (1, true)", v, ok)
	}
}

// TestHashMapMarshalMatchesCollect confirms that
// json.Marshal(HashMap) produces the same byte content as
// json.Marshal(HashMap.Collect()). Byte-exact equality is not
// asserted because Go map iteration order is non-deterministic;
// the test parses both outputs back into maps and compares
// them as a set.
func TestHashMapMarshalMatchesCollect(t *testing.T) {
	src := maps.NewHashMap[string, int]()
	src.Put("a", 1)
	src.Put("b", 2)

	fromMap, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal map: %v", err)
	}
	fromCollect, err := json.Marshal(src.Collect())
	if err != nil {
		t.Fatalf("Marshal slice: %v", err)
	}

	var fromMapParsed, fromCollectParsed map[string]int
	if err := json.Unmarshal(fromMap, &fromMapParsed); err != nil {
		t.Fatalf("parse fromMap: %v", err)
	}
	if err := json.Unmarshal(fromCollect, &fromCollectParsed); err != nil {
		t.Fatalf("parse fromCollect: %v", err)
	}
	if !reflect.DeepEqual(fromMapParsed, fromCollectParsed) {
		t.Fatalf("Marshal map and slice disagree: %v vs %v", fromMapParsed, fromCollectParsed)
	}
}
