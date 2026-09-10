package lists

import "testing"

func TestArrayListAndLinkedListCommonOperations(t *testing.T) {
	array := NewArrayList[int]()
	linked := NewLinkedList[int]()

	for _, list := range []interface {
		Add(int)
		AddFirst(int)
		RemoveFirst() (int, bool)
		RemoveLast() (int, bool)
		Size() int
		IsEmpty() bool
		Clear()
	}{array, linked} {
		list.Add(2)
		list.AddFirst(1)
		list.Add(3)
		if list.Size() != 3 || list.IsEmpty() {
			t.Fatalf("common size API is inconsistent: size=%d empty=%v",
				list.Size(), list.IsEmpty())
		}
		if value, ok := list.RemoveFirst(); !ok || value != 1 {
			t.Fatalf("RemoveFirst: got (%d, %v), want (1, true)", value, ok)
		}
		if value, ok := list.RemoveLast(); !ok || value != 3 {
			t.Fatalf("RemoveLast: got (%d, %v), want (3, true)", value, ok)
		}
		list.Clear()
		if !list.IsEmpty() || list.Size() != 0 {
			t.Fatalf("Clear did not empty the collection")
		}
	}
}

func TestLinkedListTransformationsMatchArrayListShape(t *testing.T) {
	linked := LinkedListOf(1, 2, 3)

	mapped := linked.Map(func(value int) string {
		return string(rune('0' + value))
	})
	if got, want := mapped.Collect(), []string{"1", "2", "3"}; !equalValues(got, want) {
		t.Fatalf("Map: got %v, want %v", got, want)
	}

	flat := linked.FlatMap(func(value int) *ArrayList[int] {
		return ArrayListOf(value, value*10)
	})
	if got, want := flat.Collect(), []int{1, 10, 2, 20, 3, 30}; !equalValues(got, want) {
		t.Fatalf("FlatMap: got %v, want %v", got, want)
	}

	if got, want := linked.Take(2).Collect(), []int{1, 2}; !equalValues(got, want) {
		t.Fatalf("Take: got %v, want %v", got, want)
	}
	if got, want := linked.Drop(1).Collect(), []int{2, 3}; !equalValues(got, want) {
		t.Fatalf("Drop: got %v, want %v", got, want)
	}
	if got, want := linked.SortBy(func(a, b int) int { return a - b }).Collect(), []int{1, 2, 3}; !equalValues(got, want) {
		t.Fatalf("SortBy: got %v, want %v", got, want)
	}
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
