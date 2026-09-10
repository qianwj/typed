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
	if value, ok := set.MinBy(func(a, b int) int { return a - b }); !ok || value != 1 {
		t.Fatalf("MinBy: got (%d, %v), want (1, true)", value, ok)
	}
	if value, ok := set.MaxBy(func(a, b int) int { return a - b }); !ok || value != 3 {
		t.Fatalf("MaxBy: got (%d, %v), want (3, true)", value, ok)
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
