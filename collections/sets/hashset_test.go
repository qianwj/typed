package sets

import (
	"testing"

	"github.com/qianwj/typed/collections/lists"
)

func TestHashSetBasicOperations(t *testing.T) {
	set := HashSetOf(1, 2, 2, 3)
	if got, want := set.Size(), 3; got != want {
		t.Fatalf("Size: got %d, want %d", got, want)
	}
	if !set.Contains(2) || set.Contains(4) {
		t.Fatalf("Contains returned an incorrect result")
	}

	set.Add(3)
	set.Add(4)
	set.Remove(2)
	if got, want := set.Size(), 3; got != want {
		t.Fatalf("mutations: got size %d, want %d", got, want)
	}
	if set.Contains(2) || !set.Contains(4) {
		t.Fatalf("mutations did not update membership")
	}

	set.Clear()
	if !set.IsEmpty() || set.Size() != 0 {
		t.Fatalf("Clear did not empty the set")
	}
}

func TestHashSetZeroValueAndCollectCopy(t *testing.T) {
	var set HashSet[string]
	set.Add("a")
	set.Add("b")

	values := set.Collect()
	if !sameMembers(values, []string{"a", "b"}) {
		t.Fatalf("Collect: got %v", values)
	}
	values[0] = "changed"
	if !set.Contains("a") || !set.Contains("b") {
		t.Fatalf("Collect returned a slice that aliases the set")
	}
}

func TestHashSetStreamIsSnapshot(t *testing.T) {
	set := HashSetOf(1, 2)
	stream := set.Stream()

	set.Add(3)
	set.Remove(1)

	if got := stream.Collect(); !sameMembers(got, []int{1, 2}) {
		t.Fatalf("Stream snapshot: got %v, want members {1, 2}", got)
	}
}

func TestHashSetTransformationsAndSetOperations(t *testing.T) {
	left := HashSetOf(1, 2, 3)
	right := HashSetOf(3, 4)

	if got := left.Filter(func(value int) bool { return value%2 == 1 }).Collect(); !sameMembers(got, []int{1, 3}) {
		t.Fatalf("Filter: got %v", got)
	}
	if got := left.MapSet(func(value int) int { return value % 2 }).Collect(); !sameMembers(got, []int{0, 1}) {
		t.Fatalf("MapSet: got %v", got)
	}
	if got := left.FlatMapSet(func(value int) *HashSet[int] {
		return HashSetOf(value, value%2)
	}).Collect(); !sameMembers(got, []int{0, 1, 2, 3}) {
		t.Fatalf("FlatMapSet: got %v", got)
	}
	if got := left.Union(right).Collect(); !sameMembers(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Union: got %v", got)
	}
	if got := left.Intersect(right).Collect(); !sameMembers(got, []int{3}) {
		t.Fatalf("Intersect: got %v", got)
	}
	if got := left.Difference(right).Collect(); !sameMembers(got, []int{1, 2}) {
		t.Fatalf("Difference: got %v", got)
	}
	if got := left.SymmetricDifference(right).Collect(); !sameMembers(got, []int{1, 2, 4}) {
		t.Fatalf("SymmetricDifference: got %v", got)
	}

	if !HashSetOf(1, 2).IsSubsetOf(left) || !left.IsSupersetOf(HashSetOf(1, 2)) {
		t.Fatalf("subset/superset check returned false for valid sets")
	}
	if left.IsSubsetOf(HashSetOf(1, 2)) {
		t.Fatalf("subset check returned true for an invalid subset")
	}

	if left.Size() != 3 || right.Size() != 2 {
		t.Fatalf("set operations mutated their inputs")
	}
}

func TestHashSetMapAllowsNonComparableResults(t *testing.T) {
	set := HashSetOf(1, 2, 3)
	values := set.Map(func(value int) []int {
		return []int{value, value * 10}
	}).Collect()

	if len(values) != 3 {
		t.Fatalf("Map: got %d values, want 3", len(values))
	}
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		if len(value) != 2 || value[1] != value[0]*10 {
			t.Fatalf("Map produced an invalid result: %v", value)
		}
		seen[value[0]] = true
	}
	for _, value := range []int{1, 2, 3} {
		if !seen[value] {
			t.Fatalf("Map did not preserve source value %d", value)
		}
	}
}

func TestHashSetFlatMapAllowsNonComparableResults(t *testing.T) {
	set := HashSetOf(1, 2, 3)
	values := set.FlatMap(func(value int) *lists.ArrayList[[]int] {
		return lists.ArrayListOf([]int{value}, []int{value * 10})
	}).Collect()

	if len(values) != 6 {
		t.Fatalf("FlatMap: got %d values, want 6", len(values))
	}
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		if len(value) != 1 {
			t.Fatalf("FlatMap produced an invalid result: %v", value)
		}
		seen[value[0]] = true
	}
	for _, value := range []int{1, 2, 3, 10, 20, 30} {
		if !seen[value] {
			t.Fatalf("FlatMap did not preserve flattened value %d", value)
		}
	}
}

func TestHashSetCommonCollectionOperations(t *testing.T) {
	set := HashSetOf(3, 1, 2)

	visited := make(map[int]bool, set.Size())
	peeked := set.Peek(func(value int) {
		visited[value] = true
	})
	if !sameMembers(peeked.Collect(), []int{1, 2, 3}) || !sameMembers(keys(visited), []int{1, 2, 3}) {
		t.Fatalf("Peek: got set %v, visited %v", peeked.Collect(), keys(visited))
	}

	sum := set.Reduce(0, func(acc, value int) int { return acc + value })
	if sum != 6 {
		t.Fatalf("Reduce: got %d, want 6", sum)
	}
	if got := set.SortBy(func(a, b int) int { return a - b }).Collect(); !equalValues(got, []int{1, 2, 3}) {
		t.Fatalf("SortBy: got %v", got)
	}
	if value := set.MinBy(func(a, b int) int { return a - b }).OrElse(0); value != 1 {
		t.Fatalf("MinBy: got %d, want 1", value)
	}
	if value := set.MaxBy(func(a, b int) int { return a - b }).OrElse(0); value != 3 {
		t.Fatalf("MaxBy: got %d, want 3", value)
	}
}

func TestHashSetPredicatesShortCircuit(t *testing.T) {
	set := HashSetOf(1, 2, 3)

	anyCalls := 0
	if !set.Any(func(int) bool {
		anyCalls++
		return true
	}) || anyCalls != 1 {
		t.Fatalf("Any: calls %d, want 1", anyCalls)
	}

	allCalls := 0
	if set.All(func(int) bool {
		allCalls++
		return false
	}) || allCalls != 1 {
		t.Fatalf("All: calls %d, want 1", allCalls)
	}
}

func sameMembers[T comparable](got, want []T) bool {
	if len(got) != len(want) {
		return false
	}
	seen := make(map[T]struct{}, len(got))
	for _, value := range got {
		seen[value] = struct{}{}
	}
	for _, value := range want {
		if _, ok := seen[value]; !ok {
			return false
		}
	}
	return true
}

func keys[T comparable](values map[T]bool) []T {
	out := make([]T, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	return out
}

func equalValues[T comparable](got, want []T) bool {
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

// TestNone covers the previously 0%-covered HashSet.None: true on no match,
// false when at least one value matches, true on an empty set.
func TestNone(t *testing.T) {
	t.Parallel()

	empty := NewHashSet[int]()
	if !empty.None(func(int) bool { return true }) {
		t.Fatal("empty set: None should return true")
	}

	s := NewHashSet[int]()
	s.Add(1); s.Add(2); s.Add(3); s.Add(4); s.Add(5)
	if !s.None(func(v int) bool { return v > 10 }) {
		t.Fatal("no values > 10: None should return true")
	}
	if s.None(func(v int) bool { return v%2 == 0 }) {
		t.Fatal("even values exist: None should return false")
	}
}

// TestFind covers the previously 0%-covered HashSet.Find: present Optional
// on a hit, absent Optional on a miss.
func TestFind(t *testing.T) {
	t.Parallel()

	empty := NewHashSet[int]()
	if got := empty.Find(func(int) bool { return true }); !got.IsEmpty() {
		t.Fatalf("empty set: Find should return absent; got %v", got)
	}

	s := NewHashSet[int]()
	s.Add(1); s.Add(2); s.Add(3); s.Add(4); s.Add(5)
	if got := s.Find(func(v int) bool { return v == 3 }); got.IsEmpty() || got.Get() != 3 {
		t.Fatalf("hit: got %v, want Optional{3}", got)
	}
	if got := s.Find(func(v int) bool { return v > 100 }); !got.IsEmpty() {
		t.Fatalf("miss: got %v, want absent", got)
	}
}

// TestSimpleMethodBranches exercises the previously-low-covered basic
// accessors (Add, Remove, Contains, Size, Clear, ForEach) with both the
// "found" and "not found" branches.
func TestSimpleMethodBranches(t *testing.T) {
	t.Parallel()

	s := NewHashSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(1) // duplicate, no return value to inspect
	s.Add(3)
	if s.Size() != 3 {
		t.Fatalf("Size = %d, want 3", s.Size())
	}
	if !s.Contains(2) {
		t.Fatal("Contains(2) should be true")
	}
	if s.Contains(99) {
		t.Fatal("Contains(99) should be false")
	}
	s.Remove(2)
	if s.Contains(2) {
		t.Fatal("Contains(2) after Remove should be false")
	}
	s.Remove(99) // not present, must not panic

	count := 0
	s.ForEach(func(int) { count++ })
	if count != 2 {
		t.Fatalf("ForEach visited %d, want 2", count)
	}

	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear should leave set empty")
	}
	if s.Size() != 0 {
		t.Fatalf("Size after Clear = %d, want 0", s.Size())
	}
}

// TestSetAlgebra covers the IsSubsetOf / IsSupersetOf branches that were
// partially uncovered: self-comparison, equal sets, and the partial
// overlap that produces a strict (non-) subset.
func TestSetAlgebra(t *testing.T) {
	t.Parallel()

	mkSet := func(vs ...int) *HashSet[int] {
		s := NewHashSet[int]()
		for _, v := range vs {
			s.Add(v)
		}
		return s
	}
	a := mkSet(1, 2, 3)
	b := mkSet(1, 2, 3, 4, 5)

	if !a.IsSubsetOf(b) {
		t.Fatal("a should be a subset of b")
	}
	if b.IsSubsetOf(a) {
		t.Fatal("b should not be a subset of a")
	}
	if !a.IsSubsetOf(a) {
		t.Fatal("a should be a subset of itself")
	}
	if !b.IsSupersetOf(a) {
		t.Fatal("b should be a superset of a")
	}
	if a.IsSupersetOf(b) {
		t.Fatal("a should not be a superset of b")
	}

	// Symmetric difference: items in one but not both.
	diff := a.SymmetricDifference(b)
	if diff.Size() != 2 {
		t.Fatalf("SymmetricDifference size = %d, want 2", diff.Size())
	}

	// Intersect: items in both.
	both := a.Intersect(mkSet(2, 3, 4))
	if both.Size() != 2 {
		t.Fatalf("Intersect size = %d, want 2", both.Size())
	}
}

// TestMinMaxBy covers the MinBy / MaxBy branches that compare values.
func TestMinMaxBy(t *testing.T) {
	t.Parallel()

	empty := NewHashSet[int]()
	if !empty.MinBy(func(a, b int) int { return a - b }).IsEmpty() {
		t.Fatal("MinBy on empty set should return absent")
	}
	if !empty.MaxBy(func(a, b int) int { return a - b }).IsEmpty() {
		t.Fatal("MaxBy on empty set should return absent")
	}

	s := NewHashSet[int]()
	for _, v := range []int{3, 1, 4, 1, 5, 9, 2, 6} { // duplicates collapse
		s.Add(v)
	}
	if got := s.MinBy(func(a, b int) int { return a - b }).Get(); got != 1 {
		t.Fatalf("MinBy = %d, want 1", got)
	}
	if got := s.MaxBy(func(a, b int) int { return a - b }).Get(); got != 9 {
		t.Fatalf("MaxBy = %d, want 9", got)
	}
}
