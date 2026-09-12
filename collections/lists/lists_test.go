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

// ---------- gap-fill coverage for ArrayList and LinkedList ----------
//
// The tests in this block cover the function paths that the original
// lists_test.go suite did not exercise: every 0% function in the
// package coverage report, plus the empty-source and "n <= 0" branches
// of the partial-coverage functions. They are written in the same
// style as the rest of the file (table-free, focused assertions).

// TestArrayListFilter covers keep-all, drop-all, keep-some, and the
// empty source.
func TestArrayListFilter(t *testing.T) {
	t.Run("keep some", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3, 4, 5, 6).Filter(func(n int) bool { return n%2 == 0 })
		if want := []int{2, 4, 6}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("keep all", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Filter(func(int) bool { return true })
		if want := []int{1, 2, 3}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("drop all", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Filter(func(int) bool { return false })
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewArrayList[int]().Filter(func(int) bool { return true })
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("result is independent copy", func(t *testing.T) {
		src := ArrayListOf(1, 2, 3, 4)
		filtered := src.Filter(func(n int) bool { return n%2 == 0 })
		// Mutating src after Filter must not change the result.
		src.Add(99)
		if want := []int{2, 4}; !equalValues(filtered.Collect(), want) {
			t.Fatalf("filtered: got %v, want %v (mutated src leaked)", filtered.Collect(), want)
		}
	})
}

// TestArrayListMap covers type-changing Map, same-type Map, empty
// source, and verifying that the result is an independent copy.
func TestArrayListMap(t *testing.T) {
	t.Run("int -> string", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Map(func(n int) string { return string(rune('a' + n - 1)) })
		if want := []string{"a", "b", "c"}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("int -> int", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Map(func(n int) int { return n * 10 })
		if want := []int{10, 20, 30}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("string -> []int", func(t *testing.T) {
		got := ArrayListOf("a", "bb", "ccc").Map(func(s string) int { return len(s) })
		if want := []int{1, 2, 3}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewArrayList[int]().Map(func(n int) int { return n + 1 })
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("chained Map", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).
			Map(func(n int) int { return n + 1 }).
			Map(func(n int) int { return n * n })
		if want := []int{4, 9, 16}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
}

// TestArrayListFlatMap covers the cases the source has to consider:
// multi-element inner, nil inner (treated as empty), empty source, and
// single-element inner.
func TestArrayListFlatMap(t *testing.T) {
	t.Run("multi-element inner", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).FlatMap(func(n int) *ArrayList[int] {
			return ArrayListOf(n, n*10)
		})
		if want := []int{1, 10, 2, 20, 3, 30}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("nil inner is empty", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).FlatMap(func(n int) *ArrayList[int] {
			if n%2 == 0 {
				return nil
			}
			return ArrayListOf(n)
		})
		if want := []int{1, 3}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("empty inner", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).FlatMap(func(int) *ArrayList[int] {
			return NewArrayList[int]()
		})
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewArrayList[int]().FlatMap(func(n int) *ArrayList[int] {
			return ArrayListOf(n)
		})
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("type-changing FlatMap", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).FlatMap(func(n int) *ArrayList[string] {
			return ArrayListOf(string(rune('a' + n - 1)))
		})
		if want := []string{"a", "b", "c"}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
}

// TestArrayListTake covers n <= 0 (empty), n > 0 but smaller than
// the source, n equal to the source, and n larger than the source.
func TestArrayListTake(t *testing.T) {
	src := ArrayListOf(1, 2, 3, 4, 5)

	t.Run("n smaller", func(t *testing.T) {
		got := src.Take(3).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n equal", func(t *testing.T) {
		got := src.Take(5).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n larger", func(t *testing.T) {
		got := src.Take(100).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n zero", func(t *testing.T) {
		if got := src.Take(0).Size(); got != 0 {
			t.Fatalf("Take(0): got %d, want 0", got)
		}
	})
	t.Run("n negative", func(t *testing.T) {
		if got := src.Take(-3).Size(); got != 0 {
			t.Fatalf("Take(-3): got %d, want 0", got)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if got := NewArrayList[int]().Take(5).Size(); got != 0 {
			t.Fatalf("Take on empty: got %d, want 0", got)
		}
	})
}

// TestArrayListDrop covers the n <= 0 fast path, n < size, n == size,
// n > size, and the empty source.
func TestArrayListDrop(t *testing.T) {
	src := ArrayListOf(1, 2, 3, 4, 5)

	t.Run("n smaller", func(t *testing.T) {
		got := src.Drop(2).Collect()
		if want := []int{3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n equal", func(t *testing.T) {
		if got := src.Drop(5).Size(); got != 0 {
			t.Fatalf("Drop(5): got %d, want 0", got)
		}
	})
	t.Run("n larger", func(t *testing.T) {
		if got := src.Drop(100).Size(); got != 0 {
			t.Fatalf("Drop(100): got %d, want 0", got)
		}
	})
	t.Run("n zero", func(t *testing.T) {
		got := src.Drop(0).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n negative", func(t *testing.T) {
		got := src.Drop(-3).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if got := NewArrayList[int]().Drop(5).Size(); got != 0 {
			t.Fatalf("Drop on empty: got %d, want 0", got)
		}
	})
}

// TestArrayListDistinct covers the dedup path with the default
// comparable comparison: keep-all, drop-all-duplicates, mixed, and
// a single-element source.
func TestArrayListDistinct(t *testing.T) {
	t.Run("all distinct", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3, 4).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{1, 2, 3, 4}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("all duplicates", func(t *testing.T) {
		got := ArrayListOf(7, 7, 7, 7).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{7}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("first occurrence wins", func(t *testing.T) {
		got := ArrayListOf(1, 2, 1, 3, 2, 1).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if got := NewArrayList[int]().Distinct(func(a, b int) bool { return a == b }).Size(); got != 0 {
			t.Fatalf("got %d, want 0", got)
		}
	})
	t.Run("custom eq on structs", func(t *testing.T) {
		type pt struct{ x, y int }
		eq := func(a, b pt) bool { return a.x == b.x && a.y == b.y }
		got := ArrayListOf(pt{1, 1}, pt{1, 2}, pt{2, 1}, pt{1, 1}).
			Distinct(eq).Collect()
		if len(got) != 3 {
			t.Fatalf("got %v, want 3 elements", got)
		}
		if got[0] != (pt{1, 1}) || got[1] != (pt{1, 2}) || got[2] != (pt{2, 1}) {
			t.Fatalf("got %v, want [{1,1} {1,2} {2,1}]", got)
		}
	})
}

// TestArrayListConcat covers both empty, one empty, both populated,
// and verifying that mutating the source after Concat does not affect
// the result.
func TestArrayListConcat(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		if got := NewArrayList[int]().Concat(NewArrayList[int]()).Size(); got != 0 {
			t.Fatalf("got %d, want 0", got)
		}
	})
	t.Run("empty + populated", func(t *testing.T) {
		got := NewArrayList[int]().Concat(ArrayListOf(1, 2, 3)).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + empty", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Concat(NewArrayList[int]()).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + populated", func(t *testing.T) {
		got := ArrayListOf(1, 2).Concat(ArrayListOf(3, 4, 5)).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("result is independent of sources", func(t *testing.T) {
		left := ArrayListOf(1, 2)
		right := ArrayListOf(3, 4)
		combined := left.Concat(right)
		left.Add(99)
		right.Add(100)
		if want := []int{1, 2, 3, 4}; !equalValues(combined.Collect(), want) {
			t.Fatalf("got %v, want %v (mutations leaked)", combined.Collect(), want)
		}
	})
}

// TestArrayListPeek verifies that Peek observes every element in
// order, does not modify the returned list, and visit is not called
// when the source is empty.
func TestArrayListPeek(t *testing.T) {
	t.Run("observes in order", func(t *testing.T) {
		var observed []int
		src := ArrayListOf(1, 2, 3, 4).Peek(func(n int) { observed = append(observed, n*10) })
		if want := []int{10, 20, 30, 40}; !equalValues(observed, want) {
			t.Fatalf("Peek observed: got %v, want %v", observed, want)
		}
		if want := []int{1, 2, 3, 4}; !equalValues(src.Collect(), want) {
			t.Fatalf("Peek returned: got %v, want %v", src.Collect(), want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		called := 0
		src := NewArrayList[int]().Peek(func(int) { called++ })
		if called != 0 {
			t.Fatalf("Peek on empty: visit called %d times, want 0", called)
		}
		if src.Size() != 0 {
			t.Fatalf("Peek on empty: size %d, want 0", src.Size())
		}
	})
	t.Run("result is independent copy", func(t *testing.T) {
		src := ArrayListOf(1, 2, 3)
		peeked := src.Peek(func(int) {})
		src.Add(99)
		if want := []int{1, 2, 3}; !equalValues(peeked.Collect(), want) {
			t.Fatalf("got %v, want %v (mutated src leaked)", peeked.Collect(), want)
		}
	})
}

// TestArrayListAnyAllNone covers Any, All, None for the truthy and
// falsy branches plus the empty source. All three are walked through
// the same for-loop body in the implementation, so a single test
// function keeps the symmetry visible.
func TestArrayListAnyAllNone(t *testing.T) {
	t.Run("Any", func(t *testing.T) {
		if !ArrayListOf(1, 2, 3).Any(func(n int) bool { return n == 2 }) {
			t.Fatal("Any: expected true for n == 2")
		}
		if ArrayListOf(1, 2, 3).Any(func(n int) bool { return n == 99 }) {
			t.Fatal("Any: expected false for n == 99")
		}
		if NewArrayList[int]().Any(func(int) bool { return true }) {
			t.Fatal("Any on empty: expected false")
		}
	})
	t.Run("All", func(t *testing.T) {
		if !ArrayListOf(2, 4, 6).All(func(n int) bool { return n%2 == 0 }) {
			t.Fatal("All: expected true for n%2==0")
		}
		if ArrayListOf(2, 3, 6).All(func(n int) bool { return n%2 == 0 }) {
			t.Fatal("All: expected false for n%2==0")
		}
		// All on an empty list is vacuously true.
		if !NewArrayList[int]().All(func(int) bool { return false }) {
			t.Fatal("All on empty: expected vacuously true")
		}
	})
	t.Run("None", func(t *testing.T) {
		if !ArrayListOf(1, 2, 3).None(func(n int) bool { return n > 100 }) {
			t.Fatal("None: expected true when nothing matches")
		}
		if ArrayListOf(1, 2, 3).None(func(n int) bool { return n == 2 }) {
			t.Fatal("None: expected false when something matches")
		}
		if !NewArrayList[int]().None(func(int) bool { return true }) {
			t.Fatal("None on empty: expected true")
		}
	})
}

// TestArrayListFind covers the found, not-found, and empty branches.
func TestArrayListFind(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		got := ArrayListOf(10, 20, 30, 40).Find(func(n int) bool { return n >= 30 }).OrElse(0)
		if got != 30 {
			t.Fatalf("got %d, want 30", got)
		}
	})
	t.Run("not found", func(t *testing.T) {
		if ArrayListOf(1, 2, 3).Find(func(n int) bool { return n > 100 }).IsPresent() {
			t.Fatal("Find: expected absent")
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if NewArrayList[int]().Find(func(int) bool { return true }).IsPresent() {
			t.Fatal("Empty.Find: expected absent")
		}
	})
}

// TestArrayListReduce covers same-type acc, type-changing acc, empty
// source, and a single-element source. Reduce is the only method on
// ArrayList that supports a different accumulator type via its own
// type parameter [R any].
func TestArrayListReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3, 4, 5).Reduce(0, func(acc, v int) int { return acc + v })
		if got != 15 {
			t.Fatalf("got %d, want 15", got)
		}
	})
	t.Run("type-changing (int -> string)", func(t *testing.T) {
		got := ArrayListOf(1, 2, 3).Reduce("", func(acc string, v int) string {
			if acc == "" {
				return string(rune('0' + v))
			}
			return acc + "," + string(rune('0'+v))
		})
		if got != "1,2,3" {
			t.Fatalf("got %q, want %q", got, "1,2,3")
		}
	})
	t.Run("type-changing (int -> []int)", func(t *testing.T) {
		// Build a reversed copy via Reduce with a different acc type.
		got := ArrayListOf(1, 2, 3, 4).Reduce([]int(nil), func(acc []int, v int) []int {
			return append([]int{v}, acc...)
		})
		if want := []int{4, 3, 2, 1}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewArrayList[int]().Reduce(99, func(acc, v int) int { return acc + v })
		if got != 99 {
			t.Fatalf("got %d, want 99 (init)", got)
		}
	})
	t.Run("single element", func(t *testing.T) {
		got := ArrayListOf(42).Reduce(0, func(acc, v int) int { return acc + v })
		if got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})
}

// TestArrayListSortBy covers empty source, single element, multi
// element with a comparator, and verifying that the result is an
// independent copy.
func TestArrayListSortBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if got := NewArrayList[int]().SortBy(func(a, b int) int { return a - b }).Size(); got != 0 {
			t.Fatalf("got %d, want 0", got)
		}
	})
	t.Run("single", func(t *testing.T) {
		got := ArrayListOf(42).SortBy(func(a, b int) int { return a - b }).Collect()
		if want := []int{42}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("multi", func(t *testing.T) {
		got := ArrayListOf(3, 1, 4, 1, 5, 9, 2, 6).
			SortBy(func(a, b int) int { return a - b }).Collect()
		if want := []int{1, 1, 2, 3, 4, 5, 6, 9}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("result is independent", func(t *testing.T) {
		src := ArrayListOf(3, 1, 2)
		sorted := src.SortBy(func(a, b int) int { return a - b })
		src.Add(99)
		if want := []int{1, 2, 3}; !equalValues(sorted.Collect(), want) {
			t.Fatalf("got %v, want %v (mutated src leaked)", sorted.Collect(), want)
		}
	})
}

// TestLinkedListInsert covers Insert at the start, middle, end, the
// panic paths (negative index, index > size), and the post-insert
// state of the doubly-linked structure.
func TestLinkedListInsert(t *testing.T) {
	t.Run("middle", func(t *testing.T) {
		l := LinkedListOf(1, 3)
		l.Insert(1, 2)
		if got := l.Collect(); !equalValues(got, []int{1, 2, 3}) {
			t.Fatalf("got %v, want [1 2 3]", got)
		}
		// Internal pointers after a middle insert: head=1, head.next=2,
		// middle.prev=1, middle.next=3, tail=3.
		if l.size != 3 {
			t.Fatalf("size: got %d, want 3", l.size)
		}
		if l.head.value != 1 || l.tail.value != 3 {
			t.Fatalf("head=%d tail=%d, want 1 / 3", l.head.value, l.tail.value)
		}
		if l.head.next.value != 2 {
			t.Fatalf("head.next: got %d, want 2", l.head.next.value)
		}
		if l.head.next.next.value != 3 {
			t.Fatalf("head.next.next: got %d, want 3", l.head.next.next.value)
		}
	})
	t.Run("start (delegates to AddFirst)", func(t *testing.T) {
		l := LinkedListOf(2, 3)
		l.Insert(0, 1)
		if got := l.Collect(); !equalValues(got, []int{1, 2, 3}) {
			t.Fatalf("got %v, want [1 2 3]", got)
		}
	})
	t.Run("end is append (delegates to Add)", func(t *testing.T) {
		l := LinkedListOf(1, 2)
		l.Insert(2, 3) // index == size
		if got := l.Collect(); !equalValues(got, []int{1, 2, 3}) {
			t.Fatalf("got %v, want [1 2 3]", got)
		}
	})
	t.Run("negative index panics", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for negative index")
			}
		}()
		l.Insert(-1, 99)
	})
	t.Run("index > size panics", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index > size")
			}
		}()
		l.Insert(4, 99)
	})
}

// TestLinkedListRemoveAt covers RemoveAt at the start (delegates to
// RemoveFirst), middle (unlinks a node), end (delegates to RemoveLast),
// and the panic paths.
func TestLinkedListRemoveAt(t *testing.T) {
	t.Run("middle unlinks a node", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3, 4)
		if v := l.RemoveAt(1); v != 2 {
			t.Fatalf("got %d, want 2", v)
		}
		if got := l.Collect(); !equalValues(got, []int{1, 3, 4}) {
			t.Fatalf("got %v, want [1 3 4]", got)
		}
		// After unlink: head.next must skip the removed node and
		// point to 3 directly.
		if l.head.next.value != 3 {
			t.Fatalf("head.next: got %d, want 3 (node not properly unlinked)", l.head.next.value)
		}
		if l.head.next.prev != l.head {
			t.Fatalf("head.next.prev != head: head.next.prev=%v, head=%v", l.head.next.prev, l.head)
		}
	})
	t.Run("first delegates to RemoveFirst", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		if v := l.RemoveAt(0); v != 1 {
			t.Fatalf("got %d, want 1", v)
		}
		if got := l.Collect(); !equalValues(got, []int{2, 3}) {
			t.Fatalf("got %v, want [2 3]", got)
		}
	})
	t.Run("last delegates to RemoveLast", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		if v := l.RemoveAt(2); v != 3 {
			t.Fatalf("got %d, want 3", v)
		}
		if got := l.Collect(); !equalValues(got, []int{1, 2}) {
			t.Fatalf("got %v, want [1 2]", got)
		}
	})
	t.Run("negative index panics", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for negative index")
			}
		}()
		l.RemoveAt(-1)
	})
	t.Run("index >= size panics", func(t *testing.T) {
		l := LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index >= size")
			}
		}()
		l.RemoveAt(3)
	})
}

// TestLinkedListStream covers the snapshot semantics: a Stream taken
// from a LinkedList must not be affected by mutations to the list
// after Stream() is called.
func TestLinkedListStream(t *testing.T) {
	l := LinkedListOf(1, 2, 3, 4, 5)
	s := l.Stream()
	// Walk the stream and mutate the list; the result must reflect
	// the original [1, 2, 3, 4, 5].
	l.Add(99)
	l.RemoveFirst()
	l.AddFirst(0)
	if want := []int{1, 2, 3, 4, 5}; !equalValues(s.Collect(), want) {
		t.Fatalf("got %v, want %v (mutations leaked into stream)", s.Collect(), want)
	}

	// Empty list Stream.
	empty := NewLinkedList[int]().Stream()
	if empty.Count() != 0 {
		t.Fatalf("Empty.Stream: got %d, want 0", empty.Count())
	}
}

// TestLinkedListPeek verifies that Peek observes every element in
// order, does not modify the returned list, and visit is not called
// on an empty list.
func TestLinkedListPeek(t *testing.T) {
	t.Run("observes in order", func(t *testing.T) {
		var observed []int
		src := LinkedListOf(1, 2, 3, 4).Peek(func(n int) { observed = append(observed, n*10) })
		if want := []int{10, 20, 30, 40}; !equalValues(observed, want) {
			t.Fatalf("Peek observed: got %v, want %v", observed, want)
		}
		if want := []int{1, 2, 3, 4}; !equalValues(src.Collect(), want) {
			t.Fatalf("Peek returned: got %v, want %v", src.Collect(), want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		called := 0
		src := NewLinkedList[int]().Peek(func(int) { called++ })
		if called != 0 {
			t.Fatalf("Peek on empty: visit called %d times, want 0", called)
		}
		if src.Size() != 0 {
			t.Fatalf("Peek on empty: size %d, want 0", src.Size())
		}
	})
	t.Run("result is independent", func(t *testing.T) {
		src := LinkedListOf(1, 2, 3)
		peeked := src.Peek(func(int) {})
		src.Add(99)
		if want := []int{1, 2, 3}; !equalValues(peeked.Collect(), want) {
			t.Fatalf("got %v, want %v (mutated src leaked)", peeked.Collect(), want)
		}
	})
}

// TestLinkedListAnyAllNone is the LinkedList-side parallel of
// TestArrayListAnyAllNone.
func TestLinkedListAnyAllNone(t *testing.T) {
	t.Run("Any", func(t *testing.T) {
		if !LinkedListOf(1, 2, 3).Any(func(n int) bool { return n == 2 }) {
			t.Fatal("Any: expected true for n == 2")
		}
		if LinkedListOf(1, 2, 3).Any(func(n int) bool { return n == 99 }) {
			t.Fatal("Any: expected false for n == 99")
		}
		if NewLinkedList[int]().Any(func(int) bool { return true }) {
			t.Fatal("Any on empty: expected false")
		}
	})
	t.Run("All", func(t *testing.T) {
		if !LinkedListOf(2, 4, 6).All(func(n int) bool { return n%2 == 0 }) {
			t.Fatal("All: expected true for n%2==0")
		}
		if LinkedListOf(2, 3, 6).All(func(n int) bool { return n%2 == 0 }) {
			t.Fatal("All: expected false for n%2==0")
		}
		if !NewLinkedList[int]().All(func(int) bool { return false }) {
			t.Fatal("All on empty: expected vacuously true")
		}
	})
	t.Run("None", func(t *testing.T) {
		if !LinkedListOf(1, 2, 3).None(func(n int) bool { return n > 100 }) {
			t.Fatal("None: expected true when nothing matches")
		}
		if LinkedListOf(1, 2, 3).None(func(n int) bool { return n == 2 }) {
			t.Fatal("None: expected false when something matches")
		}
		if !NewLinkedList[int]().None(func(int) bool { return true }) {
			t.Fatal("None on empty: expected true")
		}
	})
}

// TestLinkedListFind covers the found, not-found, and empty branches.
func TestLinkedListFind(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		got := LinkedListOf(10, 20, 30, 40).Find(func(n int) bool { return n >= 30 }).OrElse(0)
		if got != 30 {
			t.Fatalf("got %d, want 30", got)
		}
	})
	t.Run("not found", func(t *testing.T) {
		if LinkedListOf(1, 2, 3).Find(func(n int) bool { return n > 100 }).IsPresent() {
			t.Fatal("Find: expected absent")
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if NewLinkedList[int]().Find(func(int) bool { return true }).IsPresent() {
			t.Fatal("Empty.Find: expected absent")
		}
	})
}

// TestLinkedListReduce covers same-type acc, type-changing acc, empty
// source, and a single-element source. The accumulator type U is
// independent of T.
func TestLinkedListReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3, 4, 5).Reduce(0, func(acc, v int) int { return acc + v })
		if got != 15 {
			t.Fatalf("got %d, want 15", got)
		}
	})
	t.Run("type-changing (int -> string)", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Reduce("", func(acc string, v int) string {
			if acc == "" {
				return string(rune('0' + v))
			}
			return acc + "," + string(rune('0'+v))
		})
		if got != "1,2,3" {
			t.Fatalf("got %q, want %q", got, "1,2,3")
		}
	})
	t.Run("type-changing (int -> []int)", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3, 4).Reduce([]int(nil), func(acc []int, v int) []int {
			return append([]int{v}, acc...)
		})
		if want := []int{4, 3, 2, 1}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewLinkedList[int]().Reduce(99, func(acc, v int) int { return acc + v })
		if got != 99 {
			t.Fatalf("got %d, want 99 (init)", got)
		}
	})
	t.Run("single element", func(t *testing.T) {
		got := LinkedListOf(42).Reduce(0, func(acc, v int) int { return acc + v })
		if got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})
}

// TestLinkedListFilter covers keep-all, drop-all, keep-some, and
// the empty source. Filter returns a new LinkedList.
func TestLinkedListFilter(t *testing.T) {
	t.Run("keep some", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3, 4, 5, 6).Filter(func(n int) bool { return n%2 == 0 })
		if want := []int{2, 4, 6}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("keep all", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Filter(func(int) bool { return true })
		if want := []int{1, 2, 3}; !equalValues(got.Collect(), want) {
			t.Fatalf("got %v, want %v", got.Collect(), want)
		}
	})
	t.Run("drop all", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Filter(func(int) bool { return false })
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
	t.Run("empty source", func(t *testing.T) {
		got := NewLinkedList[int]().Filter(func(int) bool { return true })
		if got.Size() != 0 {
			t.Fatalf("got %v, want empty", got.Collect())
		}
	})
}

// TestLinkedListDistinct covers the dedup path.
func TestLinkedListDistinct(t *testing.T) {
	t.Run("all distinct", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3, 4).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{1, 2, 3, 4}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("all duplicates", func(t *testing.T) {
		got := LinkedListOf(7, 7, 7, 7).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{7}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("first occurrence wins", func(t *testing.T) {
		got := LinkedListOf(1, 2, 1, 3, 2, 1).Distinct(func(a, b int) bool { return a == b }).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if got := NewLinkedList[int]().Distinct(func(a, b int) bool { return a == b }).Size(); got != 0 {
			t.Fatalf("got %d, want 0", got)
		}
	})
}

// TestLinkedListConcat covers both empty, one empty, both populated,
// the nil-other defensive branch, and verifying that mutating the
// sources after Concat does not affect the result.
func TestLinkedListConcat(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		if got := NewLinkedList[int]().Concat(NewLinkedList[int]()).Size(); got != 0 {
			t.Fatalf("got %d, want 0", got)
		}
	})
	t.Run("empty + populated", func(t *testing.T) {
		got := NewLinkedList[int]().Concat(LinkedListOf(1, 2, 3)).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + empty", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Concat(NewLinkedList[int]()).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + populated", func(t *testing.T) {
		got := LinkedListOf(1, 2).Concat(LinkedListOf(3, 4, 5)).Collect()
		if want := []int{1, 2, 3, 4, 5}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("nil other", func(t *testing.T) {
		// Concat has a nil guard: other == nil must be treated as empty.
		got := LinkedListOf(1, 2, 3).Concat(nil).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("result is independent", func(t *testing.T) {
		left := LinkedListOf(1, 2)
		right := LinkedListOf(3, 4)
		combined := left.Concat(right)
		left.Add(99)
		right.Add(100)
		if want := []int{1, 2, 3, 4}; !equalValues(combined.Collect(), want) {
			t.Fatalf("got %v, want %v (mutations leaked)", combined.Collect(), want)
		}
	})
}

// TestLinkedListEmptySourceTransformations covers the
// empty-source branch of every transformation method in one go, so
// the empty branches that the per-function tests don't already hit
// are still exercised.
func TestLinkedListEmptySourceTransformations(t *testing.T) {
	empty := NewLinkedList[int]()
	if got := empty.Take(5).Size(); got != 0 {
		t.Fatalf("Take on empty: got %d, want 0", got)
	}
	if got := empty.Drop(5).Size(); got != 0 {
		t.Fatalf("Drop on empty: got %d, want 0", got)
	}
	if got := empty.SortBy(func(a, b int) int { return a - b }).Size(); got != 0 {
		t.Fatalf("SortBy on empty: got %d, want 0", got)
	}
}

// TestLinkedListAddFirstTailBranch covers the LinkedList AddFirst
// fast path (head > 0, prepend without grow) and the grow path
// (head == 0, prepend requires growing the slice-like node chain).
// The implementation uses preallocated node memory, so AddFirst is
// always O(1) but the "head == 0" branch must be exercised
// separately.
func TestLinkedListAddFirstBranches(t *testing.T) {
	t.Run("prepend once (single growth)", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Add(2)
		l.Add(3)
		l.AddFirst(1)
		if got := l.Collect(); !equalValues(got, []int{1, 2, 3}) {
			t.Fatalf("got %v, want [1 2 3]", got)
		}
	})
	t.Run("prepend many (exercises head==0 grow path)", func(t *testing.T) {
		l := NewLinkedList[int]()
		// Prepending to an empty list: head grows from 0.
		for i := 100; i >= 1; i-- {
			l.AddFirst(i)
		}
		want := make([]int, 100)
		for i := range want {
			want[i] = i + 1
		}
		if got := l.Collect(); !equalValues(got, want) {
			t.Fatalf("got %v, want [1..100]", got)
		}
	})
}

// TestArrayListCompactInternal exercises ArrayList.compact's main
// path: after many head-side removals, head grows large enough that
// the periodic compaction kicks in. We drive a 1 → 0 → 1 sequence
// through AddFirst/RemoveFirst/Add to force the compact branch.
func TestArrayListCompactInternal(t *testing.T) {
	a := NewArrayList[int]()

	// Build a non-empty list, then drain it from the head many
	// times to push head past 64 (the documented compaction
	// threshold).
	for i := 0; i < 200; i++ {
		a.Add(i)
	}
	for i := 0; i < 150; i++ {
		a.RemoveFirst()
	}
	// After draining, the list has 50 elements but head == 150
	// (beyond the 64 threshold). AddFirst + Add should still work
	// correctly: the live range is items[150:200], but compact
	// should have folded the prefix back.
	a.AddFirst(-1)
	if a.First().OrElse(0) != -1 {
		t.Fatalf("First after compact+AddFirst: got %d, want -1", a.First().OrElse(0))
	}
	// All the original elements are still in the list in order.
	if a.Size() != 51 {
		t.Fatalf("Size: got %d, want 51", a.Size())
	}
	if a.Last().OrElse(0) != 199 {
		t.Fatalf("Last: got %d, want 199", a.Last().OrElse(0))
	}
}

// TestArrayListCompactHeadZeroIsNoop exercises the defensive
// early-exit in ArrayList.compact when called on a list whose head
// is already 0. compact() is normally only invoked via
// compactIfNeeded (which requires head >= 64), so the head==0 path
// is effectively dead in regular usage; the test calls compact()
// directly through internal access to document the contract.
func TestArrayListCompactHeadZeroIsNoop(t *testing.T) {
	a := ArrayListOf(1, 2, 3)
	if a.head != 0 {
		t.Fatalf("precondition: head=%d, want 0", a.head)
	}
	// compact() with head == 0 must be a no-op.
	a.compact()
	if got := a.Collect(); !equalValues(got, []int{1, 2, 3}) {
		t.Fatalf("compact on head==0: got %v, want [1 2 3]", got)
	}
}

// TestArrayListInsertPanic covers the panic branches of Insert:
// negative index and index > size.
func TestArrayListInsertPanic(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		a := ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for negative index")
			}
		}()
		a.Insert(-1, 99)
	})
	t.Run("index > size", func(t *testing.T) {
		a := ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index > size")
			}
		}()
		a.Insert(4, 99)
	})
}

// TestArrayListRemoveAtPanic covers the panic branches of RemoveAt.
func TestArrayListRemoveAtPanic(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		a := ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for negative index")
			}
		}()
		a.RemoveAt(-1)
	})
	t.Run("index >= size", func(t *testing.T) {
		a := ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index >= size")
			}
		}()
		a.RemoveAt(3)
	})
}

// TestLinkedListInsertMiddleLoopBody covers the prev = prev.next
// assignment inside Insert's middle-case for-range loop. The earlier
// Insert test inserted at index 1 into a 2-element list, which
// walks the loop zero times. Here we insert in a longer list to
// force the loop body to run.
func TestLinkedListInsertMiddleLoopBody(t *testing.T) {
	l := LinkedListOf(1, 2, 3, 4, 5)
	l.Insert(3, 99) // index 3 into a 5-element list: walks prev 2 hops
	if got := l.Collect(); !equalValues(got, []int{1, 2, 3, 99, 4, 5}) {
		t.Fatalf("got %v, want [1 2 3 99 4 5]", got)
	}
	// Verify the doubly-linked pointers around the inserted node.
	inserted := l.head.next.next.next // 1→2→3→99
	if inserted.value != 99 {
		t.Fatalf("inserted node: got %d, want 99", inserted.value)
	}
	if inserted.prev.value != 3 {
		t.Fatalf("inserted.prev: got %d, want 3", inserted.prev.value)
	}
	if inserted.next.value != 4 {
		t.Fatalf("inserted.next: got %d, want 4", inserted.next.value)
	}
}

// TestLinkedListUnmarshalJSONError covers the error branch of
// LinkedList.UnmarshalJSON, which is reachable when the input is
// not a JSON array or the elements fail to decode.
func TestLinkedListUnmarshalJSONError(t *testing.T) {
	t.Run("not an array", func(t *testing.T) {
		l := NewLinkedList[int]()
		if err := l.UnmarshalJSON([]byte(`{"a":1}`)); err == nil {
			t.Fatal("expected error for non-array JSON")
		}
	})
	t.Run("malformed JSON", func(t *testing.T) {
		l := NewLinkedList[int]()
		if err := l.UnmarshalJSON([]byte(`not json`)); err == nil {
			t.Fatal("expected error for malformed JSON")
		}
	})
}

// TestLinkedListFlatMapNilInner covers the nil-inner branch of
// LinkedList.FlatMap. A nil mapped list must be treated as empty
// (continue without appending), not as a panic.
func TestLinkedListFlatMapNilInner(t *testing.T) {
	got := LinkedListOf(1, 2, 3, 4).FlatMap(func(n int) *ArrayList[int] {
		if n%2 == 0 {
			return nil
		}
		return ArrayListOf(n)
	}).Collect()
	if want := []int{1, 3}; !equalValues(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// TestLinkedListTakeNonPositive covers the n <= 0 branch of
// LinkedList.Take, which returns an empty list without iterating.
func TestLinkedListTakeNonPositive(t *testing.T) {
	t.Run("n zero", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Take(0).Collect()
		if len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("n negative", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Take(-5).Collect()
		if len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
}

// TestLinkedListDropNonPositive covers the n < 0 branch of
// LinkedList.Drop, which clamps n to 0 (returns the full list).
func TestLinkedListDropNonPositive(t *testing.T) {
	t.Run("n zero", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Drop(0).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n negative clamps to zero", func(t *testing.T) {
		got := LinkedListOf(1, 2, 3).Drop(-5).Collect()
		if want := []int{1, 2, 3}; !equalValues(got, want) {
			t.Fatalf("got %v, want %v (negative n should clamp to 0)", got, want)
		}
	})
}
