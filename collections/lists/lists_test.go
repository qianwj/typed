package lists

import (
	"testing"

	"github.com/qianwj/typed/utils/option"
)

func TestArrayListAndLinkedListCommonOperations(t *testing.T) {
	array := NewArrayList[int]()
	linked := NewLinkedList[int]()

	for _, list := range []interface {
		Add(int)
		AddFirst(int)
		RemoveFirst() option.Optional[int]
		RemoveLast() option.Optional[int]
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
		if v := list.RemoveFirst().OrElse(0); v != 1 {
			t.Fatalf("RemoveFirst: got %d, want 1", v)
		}
		if v := list.RemoveLast().OrElse(0); v != 3 {
			t.Fatalf("RemoveLast: got %d, want 3", v)
		}
		list.Clear()
		if !list.IsEmpty() || list.Size() != 0 {
			t.Fatalf("Clear did not empty the collection")
		}
	}
}

// TestArrayListOptionalReturnsOnEmpty exercises the absent-Optional
// branches of every ArrayList method that returns option.Optional[T].
// These are defensive branches: callers should not normally call Get
// on an empty list, but the contract is "absent Optional, not panic",
// and the branches must be covered.
func TestArrayListOptionalReturnsOnEmpty(t *testing.T) {
	a := NewArrayList[int]()
	if v := a.Get(0); v.IsPresent() {
		t.Fatalf("Get(0) on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.First(); v.IsPresent() {
		t.Fatalf("First on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.Last(); v.IsPresent() {
		t.Fatalf("Last on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.RemoveFirst(); v.IsPresent() {
		t.Fatalf("RemoveFirst on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.RemoveLast(); v.IsPresent() {
		t.Fatalf("RemoveLast on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.MinBy(func(x, y int) int { return x - y }); v.IsPresent() {
		t.Fatalf("MinBy on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := a.MaxBy(func(x, y int) int { return x - y }); v.IsPresent() {
		t.Fatalf("MaxBy on empty: got present %d, want absent", v.OrElse(0))
	}
}

// TestLinkedListOptionalReturnsOnEmpty is the LinkedList-side parallel
// of TestArrayListOptionalReturnsOnEmpty. Get, First, Last,
// RemoveFirst, RemoveLast, MinBy and MaxBy all return an absent
// Optional when the list is empty, and the absent branches are
// covered here.
func TestLinkedListOptionalReturnsOnEmpty(t *testing.T) {
	l := NewLinkedList[int]()
	if v := l.Get(0); v.IsPresent() {
		t.Fatalf("Get(0) on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.Get(-1); v.IsPresent() {
		t.Fatalf("Get(-1) on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.First(); v.IsPresent() {
		t.Fatalf("First on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.Last(); v.IsPresent() {
		t.Fatalf("Last on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.RemoveFirst(); v.IsPresent() {
		t.Fatalf("RemoveFirst on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.RemoveLast(); v.IsPresent() {
		t.Fatalf("RemoveLast on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.MinBy(func(x, y int) int { return x - y }); v.IsPresent() {
		t.Fatalf("MinBy on empty: got present %d, want absent", v.OrElse(0))
	}
	if v := l.MaxBy(func(x, y int) int { return x - y }); v.IsPresent() {
		t.Fatalf("MaxBy on empty: got present %d, want absent", v.OrElse(0))
	}
}

// TestLinkedListGetRemoveFirstLastBounds exercises the present
// branches of LinkedList.Get / RemoveFirst / RemoveLast on a
// non-empty list. Combined with TestLinkedListOptionalReturnsOnEmpty
// these bring the present and absent branches to coverage.
func TestLinkedListGetRemoveFirstLastBounds(t *testing.T) {
	l := LinkedListOf(10, 20, 30)

	if v := l.First().OrElse(0); v != 10 {
		t.Fatalf("First: got %d, want 10", v)
	}
	if v := l.Last().OrElse(0); v != 30 {
		t.Fatalf("Last: got %d, want 30", v)
	}
	if v := l.Get(0).OrElse(0); v != 10 {
		t.Fatalf("Get(0): got %d, want 10", v)
	}
	if v := l.Get(2).OrElse(0); v != 30 {
		t.Fatalf("Get(2): got %d, want 30", v)
	}
	if v := l.Get(3).OrElse(-1); v != -1 {
		t.Fatalf("Get(3) out-of-range: got %d, want absent (-1 sentinel)", v)
	}

	if v := l.RemoveFirst().OrElse(0); v != 10 {
		t.Fatalf("RemoveFirst: got %d, want 10", v)
	}
	if v := l.RemoveLast().OrElse(0); v != 30 {
		t.Fatalf("RemoveLast: got %d, want 30", v)
	}
	if got := l.Collect(); !equalValues(got, []int{20}) {
		t.Fatalf("after RemoveFirst+RemoveLast: got %v, want [20]", got)
	}

	// MinBy / MaxBy on a single-element list return that element.
	if v := l.MinBy(func(x, y int) int { return x - y }).OrElse(0); v != 20 {
		t.Fatalf("MinBy: got %d, want 20", v)
	}
	if v := l.MaxBy(func(x, y int) int { return x - y }).OrElse(0); v != 20 {
		t.Fatalf("MaxBy: got %d, want 20", v)
	}
}

// TestArrayListMinMaxByNonEmpty exercises the present branches of
// ArrayList.MinBy / MaxBy (the empty branches are covered in
// TestArrayListOptionalReturnsOnEmpty).
func TestArrayListMinMaxByNonEmpty(t *testing.T) {
	a := ArrayListOf(3, 1, 4, 1, 5, 9, 2, 6)
	if v := a.MinBy(func(x, y int) int { return x - y }).OrElse(0); v != 1 {
		t.Fatalf("MinBy: got %d, want 1", v)
	}
	if v := a.MaxBy(func(x, y int) int { return x - y }).OrElse(0); v != 9 {
		t.Fatalf("MaxBy: got %d, want 9", v)
	}
}

// TestLinkedListMinMaxByMultiElement exercises the inner-loop
// branches of LinkedList.MinBy / MaxBy, which only fire on lists
// of two or more elements. With a single-element list the loop
// body is skipped and only the initial "best" is returned.
func TestLinkedListMinMaxByMultiElement(t *testing.T) {
	l := LinkedListOf(3, 1, 4, 1, 5, 9, 2, 6)
	if v := l.MinBy(func(x, y int) int { return x - y }).OrElse(0); v != 1 {
		t.Fatalf("MinBy: got %d, want 1", v)
	}
	if v := l.MaxBy(func(x, y int) int { return x - y }).OrElse(0); v != 9 {
		t.Fatalf("MaxBy: got %d, want 9", v)
	}
}

// TestLinkedListRemoveFirstLastSingleElement exercises the
// "removed the only element" branch of LinkedList.RemoveFirst /
// RemoveLast: after the removal both head and tail must be nil,
// and Size must be 0.
func TestLinkedListRemoveFirstLastSingleElement(t *testing.T) {
	l := LinkedListOf(42)
	if v := l.RemoveFirst().OrElse(0); v != 42 {
		t.Fatalf("RemoveFirst single: got %d, want 42", v)
	}
	if l.head != nil || l.tail != nil {
		t.Fatalf("after RemoveFirst single: head=%v tail=%v, want both nil", l.head, l.tail)
	}

	l2 := LinkedListOf(42)
	if v := l2.RemoveLast().OrElse(0); v != 42 {
		t.Fatalf("RemoveLast single: got %d, want 42", v)
	}
	if l2.head != nil || l2.tail != nil {
		t.Fatalf("after RemoveLast single: head=%v tail=%v, want both nil", l2.head, l2.tail)
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
