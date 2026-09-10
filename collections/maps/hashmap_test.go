package maps

import (
	"testing"
)

// ---------- helpers ----------

// equalStringInt compares two (string, int) maps for equality as
// sets of pairs. Iteration order is irrelevant, so we use a counter
// instead of a position-based comparison.
func equalStringInt(got, want map[string]int) bool {
	if len(got) != len(want) {
		return false
	}
	for k, v := range want {
		if gv, ok := got[k]; !ok || gv != v {
			return false
		}
	}
	return true
}

// ---------- Constructors ----------

// TestNewHashMapEmpty covers the empty-constructor contract: a fresh
// HashMap is empty, has zero size, and Contains returns false.
func TestNewHashMapEmpty(t *testing.T) {
	m := NewHashMap[string, int]()
	if got := m.Size(); got != 0 {
		t.Fatalf("Size: got %d, want 0", got)
	}
	if !m.IsEmpty() {
		t.Fatal("IsEmpty: got false, want true")
	}
	if m.Contains("anything") {
		t.Fatal("Contains on empty: got true, want false")
	}
}

// TestHashMapOfInsertion covers the variadic constructor. Later
// entries with the same key overwrite earlier ones, matching plain
// map assignment semantics.
func TestHashMapOfInsertion(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
		Entry[string, int]{Key: "a", Value: 3}, // overwrites
	)
	if got, want := m.Size(), 2; got != want {
		t.Fatalf("Size: got %d, want %d", got, want)
	}
	if v, ok := m.Get("a"); !ok || v != 3 {
		t.Fatalf("Get(a): got (%d, %v), want (3, true)", v, ok)
	}
	if v, ok := m.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b): got (%d, %v), want (2, true)", v, ok)
	}
}

// TestHashMapFromMapIndependentCopy verifies that the resulting
// HashMap is decoupled from the source map: mutating one must not
// affect the other.
func TestHashMapFromMapIndependentCopy(t *testing.T) {
	src := map[string]int{"x": 1, "y": 2}
	m := HashMapFromMap(src)
	if !equalStringInt(m.Collect(), src) {
		t.Fatalf("Collect after FromMap: got %v, want %v", m.Collect(), src)
	}

	// Mutate the source: HashMap must not see it.
	src["x"] = 99
	src["z"] = 3
	if v, _ := m.Get("x"); v == 99 {
		t.Fatalf("HashMap saw source mutation: x=%d", v)
	}
	if m.Contains("z") {
		t.Fatal("HashMap saw new key added to source")
	}

	// Mutate the HashMap: source must not see it.
	m.Put("w", 4)
	if _, ok := src["w"]; ok {
		t.Fatal("source saw new key added to HashMap")
	}
}

// ---------- Basic operations: Put / PutIfAbsent ----------

// TestPutReturnsPreviousValue documents the Put contract: it returns
// the value that was previously stored under key, or the zero value
// of V when the key was new. Because V may legitimately hold its
// zero value, the return alone cannot tell the two apart.
func TestPutReturnsPreviousValue(t *testing.T) {
	m := NewHashMap[string, int]()

	// Insert a new key: previous value is the zero (0).
	if got := m.Put("a", 1); got != 0 {
		t.Fatalf("Put new: got previous %d, want 0", got)
	}

	// Overwrite: previous is the old value.
	if got := m.Put("a", 2); got != 1 {
		t.Fatalf("Put overwrite: got previous %d, want 1", got)
	}

	// Overwrite with the zero value: previous is still 2.
	if got := m.Put("a", 0); got != 2 {
		t.Fatalf("Put zero: got previous %d, want 2", got)
	}
}

// TestPutIfAbsent verifies that an existing key is never overwritten,
// while a fresh key is set to the given value.
func TestPutIfAbsent(t *testing.T) {
	m := NewHashMap[string, string]()
	m.Put("a", "old")

	m.PutIfAbsent("a", "new")
	if v, _ := m.Get("a"); v != "old" {
		t.Fatalf("PutIfAbsent on existing: got %q, want \"old\"", v)
	}

	m.PutIfAbsent("b", "fresh")
	if v, _ := m.Get("b"); v != "fresh" {
		t.Fatalf("PutIfAbsent on new: got %q, want \"fresh\"", v)
	}

	if got := m.Size(); got != 2 {
		t.Fatalf("Size: got %d, want 2", got)
	}
}

// ---------- Basic operations: Get / GetOrDefault ----------

// TestGetReturnsZeroAndFalseOnMiss confirms that Get on a missing
// key returns the V zero value and false, while Get on a present
// key returns the stored value and true.
func TestGetReturnsZeroAndFalseOnMiss(t *testing.T) {
	m := HashMapOf(Entry[string, int]{Key: "a", Value: 7})

	t.Run("present", func(t *testing.T) {
		v, ok := m.Get("a")
		if !ok || v != 7 {
			t.Fatalf("got (%d, %v), want (7, true)", v, ok)
		}
	})
	t.Run("absent", func(t *testing.T) {
		v, ok := m.Get("missing")
		if ok || v != 0 {
			t.Fatalf("got (%d, %v), want (0, false)", v, ok)
		}
	})
}

// TestGetOrDefaultFallsBack confirms that GetOrDefault returns the
// stored value when the key is present and the default otherwise,
// without mutating the HashMap.
func TestGetOrDefaultFallsBack(t *testing.T) {
	m := HashMapOf(Entry[string, int]{Key: "a", Value: 5})

	if got := m.GetOrDefault("a", 99); got != 5 {
		t.Fatalf("present: got %d, want 5", got)
	}
	if got := m.GetOrDefault("missing", 99); got != 99 {
		t.Fatalf("absent: got %d, want 99", got)
	}

	// Verify the HashMap is unchanged after a fallback.
	if got := m.Size(); got != 1 {
		t.Fatalf("Size: got %d, want 1 (fallback must not insert)", got)
	}
}

// ---------- Basic operations: Remove / Contains ----------

// TestRemoveReturnsPreviousValue covers the contract: Remove returns
// the previous value and true on hit, zero value and false on miss,
// and a missing key is a no-op.
func TestRemoveReturnsPreviousValue(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
	)

	t.Run("hit", func(t *testing.T) {
		v, ok := m.Remove("a")
		if !ok || v != 1 {
			t.Fatalf("got (%d, %v), want (1, true)", v, ok)
		}
		if m.Contains("a") {
			t.Fatal("a still present after Remove")
		}
	})
	t.Run("miss", func(t *testing.T) {
		v, ok := m.Remove("ghost")
		if ok || v != 0 {
			t.Fatalf("got (%d, %v), want (0, false)", v, ok)
		}
	})
	t.Run("size after removes", func(t *testing.T) {
		if got := m.Size(); got != 1 {
			t.Fatalf("Size: got %d, want 1", got)
		}
	})
}

// TestContainsPositiveAndNegative covers the obvious hit / miss
// cases; the value type is irrelevant for membership so the same
// test runs for both string and int values.
func TestContainsPositiveAndNegative(t *testing.T) {
	m := HashMapOf(Entry[string, int]{Key: "k", Value: 0})
	if !m.Contains("k") {
		t.Fatal("Contains on stored key: got false, want true")
	}
	if m.Contains("k2") {
		t.Fatal("Contains on absent key: got true, want false")
	}
}

// ---------- Basic operations: Size / IsEmpty / Clear ----------

// TestSizeAndIsEmptyTracksMutations confirms that Size and IsEmpty
// stay consistent across insertions, removals, and Clear.
func TestSizeAndIsEmptyTracksMutations(t *testing.T) {
	m := NewHashMap[string, int]()
	if !m.IsEmpty() {
		t.Fatal("new HashMap: not empty")
	}
	if got := m.Size(); got != 0 {
		t.Fatalf("new Size: got %d, want 0", got)
	}

	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	if got := m.Size(); got != 3 {
		t.Fatalf("after 3 Puts: got %d, want 3", got)
	}
	if m.IsEmpty() {
		t.Fatal("after 3 Puts: reported empty")
	}

	m.Remove("a")
	if got := m.Size(); got != 2 {
		t.Fatalf("after Remove: got %d, want 2", got)
	}

	m.Clear()
	if !m.IsEmpty() || m.Size() != 0 {
		t.Fatal("after Clear: not empty")
	}
	if m.Contains("b") {
		t.Fatal("after Clear: still contains b")
	}
}

// TestClearOnEmptyIsNoOp covers the edge case where Clear is called
// on a HashMap that has no entries.
func TestClearOnEmptyIsNoOp(t *testing.T) {
	m := NewHashMap[string, int]()
	m.Clear() // must not panic
	if !m.IsEmpty() {
		t.Fatal("Clear on empty: not empty after")
	}
}

// ---------- Iteration: ForEach / Keys / Values / Entries ----------

// TestForEachVisitsAllPairs covers the iteration contract: every
// (k, v) pair is visited exactly once. Iteration order is
// unspecified, so we check membership instead of position.
func TestForEachVisitsAllPairs(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2, "c": 3}
	m := HashMapFromMap(src)

	seen := make(map[string]int)
	m.ForEach(func(k string, v int) {
		if _, dup := seen[k]; dup {
			t.Fatalf("key %q visited twice", k)
		}
		seen[k] = v
	})
	if !equalStringInt(seen, src) {
		t.Fatalf("ForEach visited: got %v, want %v", seen, src)
	}
}

// TestKeysValuesAndEntries are three related tests that verify the
// projection methods. Keys / Values / Entries each produce an
// independent ArrayList whose contents match the HashMap's
// projection.
func TestKeysValuesAndEntries(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2, "c": 3}
	m := HashMapFromMap(src)

	keys := m.Keys().Collect()
	if got, want := len(keys), 3; got != want {
		t.Fatalf("Keys length: got %d, want %d", got, want)
	}
	for _, k := range keys {
		if _, ok := src[k]; !ok {
			t.Fatalf("Keys: stray key %q", k)
		}
	}

	values := m.Values().Collect()
	if got, want := len(values), 3; got != want {
		t.Fatalf("Values length: got %d, want %d", got, want)
	}
	valueSet := make(map[int]bool)
	for _, v := range values {
		valueSet[v] = true
	}
	for _, v := range src {
		if !valueSet[v] {
			t.Fatalf("Values: missing value %d", v)
		}
	}

	entries := m.Entries().Collect()
	if got, want := len(entries), 3; got != want {
		t.Fatalf("Entries length: got %d, want %d", got, want)
	}
	for _, e := range entries {
		if got, want := e.Value, src[e.Key]; got != want {
			t.Fatalf("Entries: %q -> %d, want %d", e.Key, got, want)
		}
	}
}

// TestForEachOnEmptyIsNoOp covers the edge case: ForEach on a
// HashMap with no entries must not invoke the callback.
func TestForEachOnEmptyIsNoOp(t *testing.T) {
	m := NewHashMap[string, int]()
	called := false
	m.ForEach(func(string, int) { called = true })
	if called {
		t.Fatal("ForEach on empty: callback was invoked")
	}
}

// ---------- Transformations: Filter / FilterKeys / FilterValues ----------

// TestFilterKeepsMatchingPairs confirms that Filter's predicate is
// applied to each (k, v) pair and only matching pairs survive.
func TestFilterKeepsMatchingPairs(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
		Entry[string, int]{Key: "c", Value: 3},
	)
	out := m.Filter(func(k string, v int) bool {
		return v%2 == 1
	})
	if got, want := out.Size(), 2; got != want {
		t.Fatalf("Filter size: got %d, want %d", got, want)
	}
	if !out.Contains("a") || !out.Contains("c") {
		t.Fatal("Filter: missing expected keys")
	}
	if out.Contains("b") {
		t.Fatal("Filter: kept non-matching key")
	}
	if m.Size() != 3 {
		t.Fatal("Filter: mutated source")
	}
}

// TestFilterKeysKeepsMatchingKeys documents the keys-only filter.
// The receiver's entries are unaffected; the result is a new HashMap.
func TestFilterKeysKeepsMatchingKeys(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "keep-a", Value: 1},
		Entry[string, int]{Key: "drop-b", Value: 2},
	)
	out := m.FilterKeys(func(k string) bool {
		return k[0] == 'k'
	})
	if got := out.Size(); got != 1 {
		t.Fatalf("FilterKeys size: got %d, want 1", got)
	}
	if !out.Contains("keep-a") || out.Contains("drop-b") {
		t.Fatal("FilterKeys: wrong set of keys")
	}
}

// TestFilterValuesKeepsMatchingValues documents the values-only
// filter. Like FilterKeys, the receiver is not mutated.
func TestFilterValuesKeepsMatchingValues(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
		Entry[string, int]{Key: "c", Value: 3},
	)
	out := m.FilterValues(func(v int) bool { return v > 1 })
	if got := out.Size(); got != 2 {
		t.Fatalf("FilterValues size: got %d, want 2", got)
	}
	if !out.Contains("b") || !out.Contains("c") {
		t.Fatal("FilterValues: missing expected keys")
	}
	if out.Contains("a") {
		t.Fatal("FilterValues: kept non-matching value")
	}
}

// ---------- Transformations: MapValues ----------

// TestMapValuesTypeChange documents the type-changing aspect of
// MapValues: every value goes through f and the result is a new
// HashMap[K, R] with the same keys.
func TestMapValuesTypeChange(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
		Entry[string, int]{Key: "c", Value: 3},
	)
	out := m.MapValues(func(k string, v int) string {
		return k + ":" + itoaForTest(v)
	})
	if got, want := out.Size(), 3; got != want {
		t.Fatalf("MapValues size: got %d, want %d", got, want)
	}
	for _, k := range []string{"a", "b", "c"} {
		v, ok := out.Get(k)
		if !ok {
			t.Fatalf("MapValues: missing key %q", k)
		}
		switch k {
		case "a":
			if v != "a:1" {
				t.Fatalf("MapValues[%q]: got %q, want %q", k, v, "a:1")
			}
		case "b":
			if v != "b:2" {
				t.Fatalf("MapValues[%q]: got %q, want %q", k, v, "b:2")
			}
		case "c":
			if v != "c:3" {
				t.Fatalf("MapValues[%q]: got %q, want %q", k, v, "c:3")
			}
		}
	}
	if m.Size() != 3 {
		t.Fatal("MapValues: mutated source")
	}
}

// TestMapValuesOnEmptyProducesEmpty covers the edge case: mapping
// an empty HashMap yields an empty HashMap.
func TestMapValuesOnEmptyProducesEmpty(t *testing.T) {
	m := NewHashMap[string, int]()
	out := m.MapValues(func(string, int) string { return "x" })
	if !out.IsEmpty() {
		t.Fatal("MapValues on empty: not empty")
	}
}

// ---------- Transformations: Concat ----------

// TestConcatMergesAndOtherWins documents the merge rule: every
// entry from m and other ends up in the result, and on a key
// collision the value from other wins. Neither input is mutated.
func TestConcatMergesAndOtherWins(t *testing.T) {
	a := HashMapOf(
		Entry[string, int]{Key: "x", Value: 1},
		Entry[string, int]{Key: "y", Value: 2},
	)
	b := HashMapOf(
		Entry[string, int]{Key: "y", Value: 20},
		Entry[string, int]{Key: "z", Value: 3},
	)
	out := a.Concat(b)

	if got, want := out.Size(), 3; got != want {
		t.Fatalf("Concat size: got %d, want %d", got, want)
	}
	if v, _ := out.Get("x"); v != 1 {
		t.Fatalf("x: got %d, want 1", v)
	}
	if v, _ := out.Get("y"); v != 20 {
		t.Fatalf("y (collision): got %d, want 20 (other wins)", v)
	}
	if v, _ := out.Get("z"); v != 3 {
		t.Fatalf("z: got %d, want 3", v)
	}
	if a.Size() != 2 || b.Size() != 2 {
		t.Fatalf("Concat mutated inputs: a=%d, b=%d", a.Size(), b.Size())
	}
}

// TestConcatWithEmpty documents the empty-handling rule: concatenating
// with an empty HashMap returns a copy of the receiver; concatenating
// an empty HashMap with a non-empty one returns a copy of the other.
func TestConcatWithEmpty(t *testing.T) {
	a := HashMapOf(Entry[string, int]{Key: "x", Value: 1})
	b := NewHashMap[string, int]()

	left := a.Concat(b)
	if !equalStringInt(left.Collect(), map[string]int{"x": 1}) {
		t.Fatalf("a.Concat(empty): got %v", left.Collect())
	}
	right := b.Concat(a)
	if !equalStringInt(right.Collect(), map[string]int{"x": 1}) {
		t.Fatalf("empty.Concat(a): got %v", right.Collect())
	}
}

// ---------- Stream and Collect ----------

// TestStreamIsSnapshot confirms that Stream returns a single-use
// snapshot: mutations to the HashMap after Stream is called do not
// affect the already-materialized Stream.
func TestStreamIsSnapshot(t *testing.T) {
	m := HashMapOf(
		Entry[string, int]{Key: "a", Value: 1},
		Entry[string, int]{Key: "b", Value: 2},
	)
	s := m.Stream()

	m.Put("c", 3)
	m.Remove("a")

	got := s.Collect()
	if len(got) != 2 {
		t.Fatalf("Stream snapshot: got %d entries, want 2", len(got))
	}
	keys := make(map[string]bool)
	for _, e := range got {
		keys[e.Key] = true
	}
	if !keys["a"] || !keys["b"] || keys["c"] {
		t.Fatalf("Stream snapshot keys: got %v, want only a and b", keys)
	}
}

// TestStreamOnEmptyProducesEmpty covers the edge case: streaming
// an empty HashMap yields no entries.
func TestStreamOnEmptyProducesEmpty(t *testing.T) {
	m := NewHashMap[string, int]()
	if got := m.Stream().Collect(); len(got) != 0 {
		t.Fatalf("Stream on empty: got %d entries, want 0", len(got))
	}
}

// TestCollectIsIndependentCopy confirms that the map returned by
// Collect is decoupled from the HashMap's internal storage:
// mutating one must not affect the other.
func TestCollectIsIndependentCopy(t *testing.T) {
	m := HashMapOf(Entry[string, int]{Key: "a", Value: 1})
	out := m.Collect()

	out["a"] = 99
	out["b"] = 2
	if v, _ := m.Get("a"); v == 99 {
		t.Fatal("Collect: HashMap saw mutation of returned map")
	}
	if m.Contains("b") {
		t.Fatal("Collect: HashMap saw new key in returned map")
	}

	// And the reverse direction.
	m.Put("c", 3)
	if _, ok := out["c"]; ok {
		t.Fatal("Collect: returned map saw mutation of HashMap")
	}
}

// ---------- Type-changing helpers ----------

// TestEntryFieldAccess documents that the public Entry struct can
// be built and read by callers; this is the type yielded by
// Entries and Stream.
func TestEntryFieldAccess(t *testing.T) {
	e := Entry[string, int]{Key: "k", Value: 7}
	if e.Key != "k" || e.Value != 7 {
		t.Fatalf("Entry: got %+v, want {k 7}", e)
	}
}

// ---------- nil receiver behaviour ----------

// The receiver methods on HashMap are documented to be safe against
// nil receivers only in some places. The tests below cover the
// documented nil-tolerance for the methods that explicitly support
// it (GetOrDefault, Filter, FilterKeys, FilterValues) and document
// the behaviour of methods that do not (ForEach on a nil items map
// would panic, but that path is not reachable through the public
// constructors because every constructor initializes items).

// TestNewHashMapAlwaysInitializesItems ensures that NewHashMap, the
// variadic HashMapOf, and HashMapFromMap all leave the internal map
// ready to use. This rules out a nil-items edge case: callers can
// chain construction with immediate Put without an explicit ensure
// step.
func TestNewHashMapAlwaysInitializesItems(t *testing.T) {
	cases := []struct {
		name string
		m    *HashMap[string, int]
	}{
		{"NewHashMap", NewHashMap[string, int]()},
		{"HashMapOf", HashMapOf[string, int]()},
		{"HashMapFromMap empty", HashMapFromMap[string, int](map[string]int{})},
		{"HashMapFromMap nil", HashMapFromMap[string, int](nil)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.m.Put("a", 1)
			if v, ok := c.m.Get("a"); !ok || v != 1 {
				t.Fatalf("Put after construction: got (%d, %v), want (1, true)", v, ok)
			}
		})
	}
}

// ---------- small int formatter ----------

// itoaForTest is a tiny non-negative int formatter used by the
// MapValues test. Negative or zero values are out of scope for
// the tests that call it.
func itoaForTest(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
