package stream_test

import (
	"strconv"
	"testing"

	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/stream"
)

// TestArrayListMap exercises the Go 1.27 generic method Map[R any] on the
// concrete ArrayList type. The test is intentionally end-to-end: it chains
// ArrayList.Filter then ArrayList.Map then converts to Stream and chains
// Stream.Map, verifying that method-level type parameters are accepted on
// concrete types.
func TestArrayListMap(t *testing.T) {
	users := []User{
		{Name: "Alice", Age: 17},
		{Name: "Bob", Age: 21},
		{Name: "Carol", Age: 30},
	}

	// Eager pipeline. The compiler must instantiate Map's R parameter as
	// string here, even though the type ArrayList[User] itself uses T any.
	usersList := lists.ArrayListOf(users...)
	adults := usersList.Filter(func(u User) bool { return u.Age >= 18 })
	names := adults.Map(func(u User) string { return u.Name })

	got := names.Collect()
	want := []string{"Bob", "Carol"}
	if !equal(got, want) {
		t.Fatalf("ArrayList.Map: got %v, want %v", got, want)
	}
}

// TestStreamMapAndFlatMap exercises the same generic method feature on
// Stream, plus FlatMap[R any] (which returns Stream[R] from Stream[T]).
func TestStreamMapAndFlatMap(t *testing.T) {
	nums := []int{1, 2, 3, 4}

	// Lazy pipeline. Map's R parameter is inferred as string.
	src := lists.ArrayListOf(nums...)
	words := src.Stream().
		Filter(func(n int) bool { return n%2 == 0 }).
		Map(strconv.Itoa)

	if got, want := words.Collect(), []string{"2", "4"}; !equal(got, want) {
		t.Fatalf("Stream.Map: got %v, want %v", got, want)
	}

	// FlatMap changes the element type as well: Stream[int] -> Stream[string].
	chars := src.Stream().
		FlatMap(func(n int) stream.Stream[string] {
			return stream.Of(strconv.Itoa(n), strconv.Itoa(n*10))
		})

	if got, want := chars.Collect(), []string{"1", "10", "2", "20", "3", "30", "4", "40"}; !equal(got, want) {
		t.Fatalf("Stream.FlatMap: got %v, want %v", got, want)
	}
}

// TestStreamChainedMapTypeChange verifies the chained call form documented
// in the README:
//
//	arr.Stream().Map(...).Map(...).Collect()
//
// Each Map introduces a fresh R parameter; the chain must compile, meaning
// every method-level type parameter is accepted on the concrete receiver.
func TestStreamChainedMapTypeChange(t *testing.T) {
	nums := lists.ArrayListOf(0, 2, 4)
	got := nums.Stream().
		Map(func(n int) int { return n + 2 }). // int -> int
		Map(func(n int) int { return n / 2 }). // int -> int
		Collect()

	if want := []int{1, 2, 3}; !equal(got, want) {
		t.Fatalf("chained Map: got %v, want %v", got, want)
	}

	// Now exercise a real type change in the chain: int -> string -> int.
	strs := lists.ArrayListOf("a", "bb", "ccc")
	lengths := strs.Stream().
		Map(func(s string) int { return len(s) }). // string -> int
		Filter(func(n int) bool { return n > 1 }).
		Collect()

	if want := []int{2, 3}; !equal(lengths, want) {
		t.Fatalf("string->int chain: got %v, want %v", lengths, want)
	}
}

// TestStreamEarlyTermination makes sure that Any / First / Take stop walking
// the source as soon as the answer is known. This is one of the reasons a
// lazy Stream is preferred over an eager ArrayList for large or infinite
// sources.
func TestStreamEarlyTermination(t *testing.T) {
	visited := 0
	counter := func(yield func(int) bool) {
		for i := 0; ; i++ {
			visited++
			if !yield(i) {
				return
			}
		}
	}

	s := stream.From(counter)

	// Any must stop after the first match.
	if !s.Any(func(n int) bool { return n == 5 }) {
		t.Fatal("Any should have returned true")
	}
	if visited != 6 { // 0..5 are checked
		t.Fatalf("Any visited %d values, want 6", visited)
	}

	// First must stop after yielding one value.
	visited = 0
	v := stream.From(counter).First().OrElse(-1)
	if v != 0 {
		t.Fatalf("First: got %d, want 0", v)
	}
	if visited != 1 {
		t.Fatalf("First visited %d values, want 1", visited)
	}
}

// TestStreamSortByMinByMaxBy verifies the comparator-provided sort and
// aggregation variants. A no-arg Sort/Min/Max would require T cmp.Ordered,
// which a single Stream[T any] type cannot express on its method set.
func TestStreamSortByMinByMaxBy(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6}

	src := lists.ArrayListOf(in...)
	sorted := src.Stream().SortBy(func(a, b int) int {
		return a - b
	}).Collect()

	if want := []int{1, 1, 2, 3, 4, 5, 6, 9}; !equal(sorted, want) {
		t.Fatalf("SortBy: got %v, want %v", sorted, want)
	}

	min, ok := src.Stream().MinBy(func(a, b int) int { return a - b })
	if !ok || min != 1 {
		t.Fatalf("MinBy: got (%v, %v), want (1, true)", min, ok)
	}

	max, ok := src.Stream().MaxBy(func(a, b int) int { return a - b })
	if !ok || max != 9 {
		t.Fatalf("MaxBy: got (%v, %v), want (9, true)", max, ok)
	}
}

// TestArrayListIsConcreteType documents (and asserts at compile time) that
// ArrayList is a concrete generic type, not an interface. The blog post on
// Go 1.27 generic methods is explicit: this is the only way to support
// type-changing fluent methods such as Map[R].
func TestArrayListIsConcreteType(t *testing.T) {
	var _ *lists.ArrayList[int] = lists.NewArrayList[int]()
	var _ *lists.LinkedList[int] = lists.NewLinkedList[int]()
	arr := lists.ArrayListOf(1, 2)
	var _ stream.Stream[int] = arr.Stream()
}

// TestLinkedListMapReduce exercises the Map / Reduce / Any / All / None /
// Find methods on LinkedList.
func TestLinkedListMapReduce(t *testing.T) {
	ll := lists.LinkedListOf(1, 2, 3, 4, 5)

	// Map changes the element type — generic method on a concrete
	// receiver, allowed by Go 1.27.
	strings := ll.Map(func(n int) string { return strconv.Itoa(n * 10) })
	if got, want := strings.Size(), 5; got != want {
		t.Fatalf("Map: size got %d, want %d", got, want)
	}
	if v, ok := strings.Get(0); !ok || v != "10" {
		t.Fatalf("Map: Get(0) got (%q, %v), want (\"10\", true)", v, ok)
	}

	// Reduce folds left-to-right (U = T case).
	sum := ll.Reduce(0, func(a, v int) int { return a + v })
	if got, want := sum, 15; got != want {
		t.Fatalf("Reduce (U=T): got %d, want %d", got, want)
	}

	// Reduce with a different accumulator type (U = string). This is the
	// reason the method takes its own type parameter; the simple T→T
	// version that used to live on the interface cannot do this.
	joined := ll.Reduce("", func(acc string, v int) string {
		if acc == "" {
			return strconv.Itoa(v)
		}
		return acc + "," + strconv.Itoa(v)
	})
	if got, want := joined, "1,2,3,4,5"; got != want {
		t.Fatalf("Reduce (U=string): got %q, want %q", got, want)
	}

	// Any / All / None short-circuit.
	if !ll.Any(func(n int) bool { return n == 3 }) {
		t.Fatalf("Any: expected true for n == 3")
	}
	if ll.Any(func(n int) bool { return n == 99 }) {
		t.Fatalf("Any: expected false for n == 99")
	}
	if !ll.All(func(n int) bool { return n > 0 }) {
		t.Fatalf("All: expected true for n > 0")
	}
	if ll.All(func(n int) bool { return n > 3 }) {
		t.Fatalf("All: expected false for n > 3")
	}
	if !ll.None(func(n int) bool { return n > 100 }) {
		t.Fatalf("None: expected true for n > 100")
	}
	if ll.None(func(n int) bool { return n == 3 }) {
		t.Fatalf("None: expected false for n == 3")
	}

	// Find returns the first match.
	if v, ok := ll.Find(func(n int) bool { return n > 2 }); !ok || v != 3 {
		t.Fatalf("Find: got (%d, %v), want (3, true)", v, ok)
	}
	if v, ok := ll.Find(func(n int) bool { return n > 100 }); ok || v != 0 {
		t.Fatalf("Find (no match): got (%d, %v), want (0, false)", v, ok)
	}
}

// TestArrayListReduceGeneral exercises the general Reduce[U any] form on
// ArrayList: building a string from a list of ints.
func TestArrayListReduceGeneral(t *testing.T) {
	arr := lists.ArrayListOf(1, 2, 3, 4, 5)

	// U = int (same as T): a sum.
	if got, want := arr.Reduce(0, func(a, v int) int { return a + v }), 15; got != want {
		t.Fatalf("Reduce int: got %d, want %d", got, want)
	}

	// U = string: a comma-joined representation. This is the case that
	// the simple T→T Reduce cannot express.
	joined := arr.Reduce("", func(acc string, v int) string {
		if acc == "" {
			return strconv.Itoa(v)
		}
		return acc + "," + strconv.Itoa(v)
	})
	if got, want := joined, "1,2,3,4,5"; got != want {
		t.Fatalf("Reduce string: got %q, want %q", got, want)
	}

	// U = []int: build a reversed slice. Different type entirely.
	reversed := arr.Reduce([]int(nil), func(acc []int, v int) []int {
		return append([]int{v}, acc...)
	})
	if want := []int{5, 4, 3, 2, 1}; !equal(reversed, want) {
		t.Fatalf("Reduce []int: got %v, want %v", reversed, want)
	}
}

// TestArrayListStreamIsSnapshot verifies that ArrayList.Stream() copies the
// internal slice. With ArrayList as a struct that owns a private []T, the
// only way to mutate after Stream() is via the API (Insert / RemoveAt);
// the snapshot is unaffected regardless of which method is used.
func TestArrayListStreamIsSnapshot(t *testing.T) {
	arr := lists.ArrayListOf(1, 2, 3)
	s := arr.Stream()

	// Mutate the ArrayList after Stream() has been built, using the
	// mutating API. The snapshot taken by Stream() must still see [1, 2, 3].
	arr.RemoveAt(0)   // arr is now [2, 3]
	arr.RemoveAt(0)   // arr is now [3]
	arr.Insert(0, 99) // arr is now [99, 3]
	arr.Insert(1, 99) // arr is now [99, 99, 3]
	arr.RemoveAt(0)   // arr is now [99, 3]
	arr.RemoveAt(0)   // arr is now [3]

	if got, want := s.Collect(), []int{1, 2, 3}; !equal(got, want) {
		t.Fatalf("snapshot: got %v, want %v (mutations leaked into stream)", got, want)
	}
}

// TestLinkedListBasics covers the doubly-linked list core operations: add
// from both ends, remove from both ends, random access, and length.
func TestLinkedListBasics(t *testing.T) {
	l := lists.NewLinkedList[int]()

	l.Add(2)
	l.Add(3)
	l.AddFirst(1)
	l.Add(4) // tail

	if got := l.Size(); got != 4 {
		t.Fatalf("Size: got %d, want 4", got)
	}
	if v := l.First().OrElse(0); v != 1 {
		t.Fatalf("First: got %d, want 1", v)
	}
	if v := l.Last().OrElse(0); v != 4 {
		t.Fatalf("Last: got %d, want 4", v)
	}
	if v, _ := l.Get(2); v != 3 {
		t.Fatalf("Get(2): got %d, want 3", v)
	}

	// RemoveFirst returns 1.
	if v, ok := l.RemoveFirst(); !ok || v != 1 {
		t.Fatalf("RemoveFirst: got (%d, %v), want (1, true)", v, ok)
	}
	// RemoveLast returns 4.
	if v, ok := l.RemoveLast(); !ok || v != 4 {
		t.Fatalf("RemoveLast: got (%d, %v), want (4, true)", v, ok)
	}
	if got := l.Size(); got != 2 {
		t.Fatalf("Size after pops: got %d, want 2", got)
	}
}

// TestLinkedListMap verifies the Go 1.27 generic method Map[R any] also
// works on the concrete LinkedList receiver, returning an ArrayList of the
// new type with the original list order.
func TestLinkedListMap(t *testing.T) {
	src := lists.LinkedListOf(1, 2, 3, 4)
	ll := src.Map(func(n int) string { return strconv.Itoa(n * 10) })

	if got, want := ll.Size(), 4; got != want {
		t.Fatalf("Map: size got %d, want %d", got, want)
	}
	for i, want := range []string{"10", "20", "30", "40"} {
		v, ok := ll.Get(i)
		if !ok || v != want {
			t.Fatalf("Map: Get(%d) got (%q, %v), want (%q, true)", i, v, ok, want)
		}
	}
}

// TestLinkedListToStream makes sure the LinkedList.Stream() adapter walks
// head-to-tail and the resulting Stream participates in the same fluent
// pipeline as ArrayList.Stream().
func TestLinkedListToStream(t *testing.T) {
	ll := lists.LinkedListOf(1, 2, 3, 4, 5)

	got := ll.Stream().
		Filter(func(n int) bool { return n%2 == 1 }).
		Map(strconv.Itoa).
		Collect()

	if want := []string{"1", "3", "5"}; !equal(got, want) {
		t.Fatalf("LinkedList.Stream: got %v, want %v", got, want)
	}
}

// TestLinkedListStreamIsSnapshot documents and verifies the snapshot
// contract on LinkedList.Stream(): the Stream must not be affected by
// mutations to the list after Stream() is called. Without the snapshot
// materialization, the Stream's underlying closure still walks the live
// node links and would see Add/Remove operations performed afterwards.
func TestLinkedListStreamIsSnapshot(t *testing.T) {
	ll := lists.LinkedListOf(1, 2, 3)
	s := ll.Stream()

	// Mutate the source list after Stream() has been built.
	ll.Add(4)
	ll.AddFirst(0)
	ll.RemoveFirst() // removes the 0 we just added
	ll.RemoveLast()  // removes the 4 we just added

	// The stream should still hold the original [1, 2, 3], not whatever
	// state the list is in after the mutations.
	if got, want := s.Collect(), []int{1, 2, 3}; !equal(got, want) {
		t.Fatalf("snapshot: got %v, want %v (mutations leaked into stream)", got, want)
	}
}

// TestLinkedListStreamSnapshotIsIndependent makes sure two independent
// Streams taken at different points in time capture different snapshots
// of the same list, and that neither is affected by subsequent mutations.
func TestLinkedListStreamSnapshotIsIndependent(t *testing.T) {
	ll := lists.LinkedListOf(1, 2, 3)

	s1 := ll.Stream()
	ll.Add(4)
	s2 := ll.Stream()
	ll.Add(5)

	if got, want := s1.Collect(), []int{1, 2, 3}; !equal(got, want) {
		t.Fatalf("s1 (pre-Add(4)): got %v, want %v", got, want)
	}
	if got, want := s2.Collect(), []int{1, 2, 3, 4}; !equal(got, want) {
		t.Fatalf("s2 (post-Add(4), pre-Add(5)): got %v, want %v", got, want)
	}
}

// TestArrayListInsert covers Insert at the start, middle, and end (which
// behaves as an append). Insert is a mutating operation, so the receiver
// must be addressable — these tests assign to a variable first.
func TestArrayListInsert(t *testing.T) {
	t.Run("middle", func(t *testing.T) {
		a := lists.ArrayListOf(1, 3)
		a.Insert(1, 2)
		if got, want := a.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("start", func(t *testing.T) {
		a := lists.ArrayListOf(2, 3)
		a.Insert(0, 1)
		if got, want := a.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("end is append", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2)
		a.Insert(2, 3) // index == len
		if got, want := a.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("out of range panics", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index > len")
			}
		}()
		a.Insert(4, 99) // 4 > 3
	})
}

// TestArrayListRemoveAt covers RemoveAt at the start, middle, last, and
// the panic path. RemoveAt returns the removed value, like the
// RemoveFirst / RemoveLast pair.
func TestArrayListRemoveAt(t *testing.T) {
	t.Run("middle returns value", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2, 3, 4)
		v := a.RemoveAt(1)
		if v != 2 {
			t.Fatalf("got %d, want 2", v)
		}
		if got, want := a.Collect(), []int{1, 3, 4}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("first", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2, 3)
		if v := a.RemoveAt(0); v != 1 {
			t.Fatalf("got %d, want 1", v)
		}
		if got, want := a.Collect(), []int{2, 3}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("last", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2, 3)
		if v := a.RemoveAt(2); v != 3 {
			t.Fatalf("got %d, want 3", v)
		}
		if got, want := a.Collect(), []int{1, 2}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("out of range panics", func(t *testing.T) {
		a := lists.ArrayListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index >= len")
			}
		}()
		a.RemoveAt(3) // 3 == len, invalid
	})
}

// TestLinkedListInsert covers Insert at the start, middle, and end. The
// start and end cases delegate to AddFirst and Add; the middle case walks
// the list once to find the predecessor.
func TestLinkedListInsert(t *testing.T) {
	t.Run("middle", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 3)
		ll.Insert(1, 2)
		if got, want := ll.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("start", func(t *testing.T) {
		ll := lists.LinkedListOf(2, 3)
		ll.Insert(0, 1)
		if got, want := ll.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("end is append", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2)
		ll.Insert(2, 3) // index == len
		if got, want := ll.Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("out of range panics", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index > len")
			}
		}()
		ll.Insert(4, 99)
	})
}

// TestLinkedListRemoveAt covers RemoveAt at the start (delegates to
// RemoveFirst), end (delegates to RemoveLast), middle (unlinks a node),
// and the panic path.
func TestLinkedListRemoveAt(t *testing.T) {
	t.Run("middle returns value", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2, 3, 4)
		if v := ll.RemoveAt(1); v != 2 {
			t.Fatalf("got %d, want 2", v)
		}
		if got, want := ll.Collect(), []int{1, 3, 4}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("first delegates to RemoveFirst", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2, 3)
		if v := ll.RemoveAt(0); v != 1 {
			t.Fatalf("got %d, want 1", v)
		}
		if got, want := ll.Collect(), []int{2, 3}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("last delegates to RemoveLast", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2, 3)
		if v := ll.RemoveAt(2); v != 3 {
			t.Fatalf("got %d, want 3", v)
		}
		if got, want := ll.Collect(), []int{1, 2}; !equal(got, want) {
			t.Fatalf("after: got %v, want %v", got, want)
		}
	})
	t.Run("out of range panics", func(t *testing.T) {
		ll := lists.LinkedListOf(1, 2, 3)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for index >= len")
			}
		}()
		ll.RemoveAt(3)
	})
}

// ---------- helpers ----------

type User struct {
	Name string
	Age  int
}

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
