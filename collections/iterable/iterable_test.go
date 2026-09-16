package iterable_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/qianwj/typed/collections/iterable"
	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/sets"
	"github.com/qianwj/typed/collections/stream"
)

var (
	_ iterable.Iterable[int] = (*lists.ArrayList[int])(nil)
	_ iterable.Iterable[int] = (*lists.LinkedList[int])(nil)
	_ iterable.Iterable[int] = (*sets.HashSet[int])(nil)
	_ iterable.Iterable[int] = stream.Stream[int]{}
)

func TestCollectionConversions(t *testing.T) {
	array := lists.ArrayListOf(99, 3, 1, 3)
	array.RemoveFirst() // Conversion must respect the live range.
	for _, tc := range []struct {
		name    string
		source  iterable.Iterable[int]
		want    []int
		ordered bool
	}{
		{"array list", array, []int{3, 1, 3}, true},
		{"linked list", lists.LinkedListOf(3, 1, 3), []int{3, 1, 3}, true},
		{"hash set", sets.HashSetOf(3, 1, 3), []int{1, 3}, false},
		{"empty array list", lists.NewArrayList[int](), nil, true},
		{"empty linked list", lists.NewLinkedList[int](), nil, true},
		{"empty hash set", sets.NewHashSet[int](), nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			array := lists.ArrayListFrom(tc.source)
			linked := lists.LinkedListFrom(tc.source)
			for _, got := range [][]int{array.Collect(), linked.Collect()} {
				if !tc.ordered {
					slices.Sort(got)
				}
				if !slices.Equal(got, tc.want) {
					t.Fatalf("list conversion = %v, want %v", got, tc.want)
				}
			}
			set := sets.HashSetFrom(tc.source)
			wantUnique := make(map[int]bool)
			for _, v := range tc.want {
				wantUnique[v] = true
				if !set.Contains(v) {
					t.Fatalf("set is missing %d", v)
				}
			}
			if set.Size() != len(wantUnique) {
				t.Fatalf("set size = %d, want %d", set.Size(), len(wantUnique))
			}
		})
	}
}

func TestConversionsHaveIndependentStorage(t *testing.T) {
	source := lists.ArrayListOf(1, 2)
	array := lists.ArrayListFrom(source)
	linked := lists.LinkedListFrom(source)
	set := sets.HashSetFrom(source)
	source.RemoveFirst()
	source.Add(3)
	if !slices.Equal(array.Collect(), []int{1, 2}) || !slices.Equal(linked.Collect(), []int{1, 2}) || set.Contains(3) || !set.Contains(1) {
		t.Fatal("source mutation affected converted collections")
	}
	array.Add(4)
	linked.Add(5)
	if !slices.Equal(source.Collect(), []int{2, 3}) || set.Size() != 2 {
		t.Fatal("destination mutation affected another collection")
	}
	fromSet := lists.ArrayListFrom(set)
	linkedFromSet := lists.LinkedListFrom(set)
	set.Clear()
	if fromSet.Size() != 2 || linkedFromSet.Size() != 2 {
		t.Fatal("set mutation affected converted lists")
	}
}

type traversal[T any] func(func(T))

func (f traversal[T]) ForEach(visit func(T)) { f(visit) }

func TestFromCustomIterableAndStream(t *testing.T) {
	calls := 0
	source := traversal[[]int](func(visit func([]int)) {
		calls++
		visit([]int{1, 2})
	})
	array := lists.ArrayListFrom(source)
	if calls != 1 || !slices.Equal(array.First().Get(), []int{1, 2}) {
		t.Fatal("ArrayListFrom must consume the source once and support non-comparable values")
	}
	linked := lists.LinkedListFrom(source)
	if calls != 2 || !slices.Equal(linked.First().Get(), []int{1, 2}) {
		t.Fatal("LinkedListFrom must consume the source once and support non-comparable values")
	}
	set := sets.HashSetFrom(stream.Of(1, 1, 2, 3).Filter(func(v int) bool { return v < 3 }))
	if set.Size() != 2 || !set.Contains(1) || !set.Contains(2) {
		t.Fatalf("stream conversion = %v, want set {1, 2}", set.Collect())
	}
}

func ExampleIterable() {
	list := lists.ArrayListOf(3, 1, 3, 2)
	set := sets.HashSetFrom(list)
	array := lists.ArrayListFrom(set).SortBy(func(a, b int) int { return a - b })
	linked := lists.LinkedListFrom(array)
	fmt.Println(set.Size())
	fmt.Println(linked.Collect())
	// Output:
	// 3
	// [1 2 3]
}
