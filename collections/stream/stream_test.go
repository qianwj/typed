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

	if min := src.Stream().MinBy(func(a, b int) int { return a - b }).OrElse(0); min != 1 {
		t.Fatalf("MinBy: got %d, want 1", min)
	}

	if max := src.Stream().MaxBy(func(a, b int) int { return a - b }).OrElse(0); max != 9 {
		t.Fatalf("MaxBy: got %d, want 9", max)
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
	if v := strings.Get(0).OrElse(""); v != "10" {
		t.Fatalf("Map: Get(0) got %q, want \"10\"", v)
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
	if v := ll.Find(func(n int) bool { return n > 2 }).OrElse(0); v != 3 {
		t.Fatalf("Find: got %d, want 3", v)
	}
	if got := ll.Find(func(n int) bool { return n > 100 }); got.IsPresent() {
		t.Fatalf("Find (no match): got present %v, want absent", got)
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
	if v := l.Get(2).OrElse(0); v != 3 {
		t.Fatalf("Get(2): got %d, want 3", v)
	}

	// RemoveFirst returns 1.
	if v := l.RemoveFirst().OrElse(0); v != 1 {
		t.Fatalf("RemoveFirst: got %d, want 1", v)
	}
	// RemoveLast returns 4.
	if v := l.RemoveLast().OrElse(0); v != 4 {
		t.Fatalf("RemoveLast: got %d, want 4", v)
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
		v := ll.Get(i).OrElse("")
		if v != want {
			t.Fatalf("Map: Get(%d) got %q, want %q", i, v, want)
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

// TestStreamAssociate exercises the list -> map collector. The mapping
// function provides both the key and the value for each element, and a
// later element with the same key overwrites an earlier one.
func TestStreamAssociate(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}
	users := []user{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Carol"},
	}

	got := lists.ArrayListOf(users...).
		Stream().
		Associate(func(u user) (int, string) { return u.ID, u.Name })

	want := map[int]string{1: "Alice", 2: "Bob", 3: "Carol"}
	if len(got) != len(want) {
		t.Fatalf("len: got %d, want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("got[%d] = %q, want %q", k, got[k], v)
		}
	}

	// Same key from multiple elements: later one wins.
	type pair struct {
		k string
		v int
	}
	pairs := []pair{
		{"a", 1},
		{"b", 2},
		{"a", 99}, // overwrites the first "a"
	}
	merged := stream.FromSlice(pairs).
		Associate(func(p pair) (string, int) { return p.k, p.v })
	if merged["a"] != 99 {
		t.Fatalf("collision: got[\"a\"] = %d, want 99", merged["a"])
	}
	if merged["b"] != 2 {
		t.Fatalf("got[\"b\"] = %d, want 2", merged["b"])
	}

	// Empty stream -> empty map (not nil).
	empty := lists.NewArrayList[int]().Stream().Associate(func(n int) (int, int) { return n, n })
	if empty == nil {
		t.Fatal("Associate on empty stream should return non-nil empty map")
	}
	if len(empty) != 0 {
		t.Fatalf("empty map len: got %d, want 0", len(empty))
	}
}

// ---------- gap-fill coverage for stream package ----------
//
// The tests in this block target the function paths that the original
// stream_test.go suite did not exercise: every 0% function in the
// package's coverage report, plus the empty-stream and downstream
// early-termination branches of the partial-coverage functions. They
// follow the same style as the rest of the file (table-free, focused
// assertions, no test frameworks) so the file stays uniform.

// TestStreamEmptyConstructor covers stream.Empty[T]() and confirms the
// empty stream is well-behaved across every terminal operation.
func TestStreamEmptyConstructor(t *testing.T) {
	empty := stream.Empty[int]()

	if got := empty.Count(); got != 0 {
		t.Fatalf("Empty.Count: got %d, want 0", got)
	}
	if got := empty.Collect(); len(got) != 0 {
		t.Fatalf("Empty.Collect: got %v, want []", got)
	}
	if empty.First().IsPresent() {
		t.Fatal("Empty.First: expected absent")
	}
	if empty.Last().IsPresent() {
		t.Fatal("Empty.Last: expected absent")
	}
	if empty.Find(func(int) bool { return true }).IsPresent() {
		t.Fatal("Empty.Find: expected absent")
	}
	if got := empty.Reduce(42, func(a, v int) int { return a + v }); got != 42 {
		t.Fatalf("Empty.Reduce: got %d, want 42 (init)", got)
	}
	// ForEach on empty must not call visit at all.
	called := 0
	empty.ForEach(func(int) { called++ })
	if called != 0 {
		t.Fatalf("Empty.ForEach visited %d times, want 0", called)
	}
	// Any / All / None: All and None are vacuously true on empty input.
	if empty.Any(func(int) bool { return true }) {
		t.Fatal("Empty.Any: expected false")
	}
	if !empty.All(func(int) bool { return false }) {
		t.Fatal("Empty.All: expected vacuously true")
	}
	if !empty.None(func(int) bool { return true }) {
		t.Fatal("Empty.None: expected true")
	}
}

// TestStreamEmptyAfterChainedOps confirms that downstream operations
// preserve the empty contract: an empty source must produce an empty
// result regardless of which intermediates are wired up.
func TestStreamEmptyAfterChainedOps(t *testing.T) {
	empty := stream.Empty[int]()

	if got := empty.Filter(func(int) bool { return true }).Count(); got != 0 {
		t.Fatalf("Empty.Filter.Count: got %d, want 0", got)
	}
	if got := empty.Map(func(n int) int { return n + 1 }).Collect(); len(got) != 0 {
		t.Fatalf("Empty.Map.Collect: got %v, want []", got)
	}
	if got := empty.Take(5).Count(); got != 0 {
		t.Fatalf("Empty.Take.Count: got %d, want 0", got)
	}
	if got := empty.Drop(5).Count(); got != 0 {
		t.Fatalf("Empty.Drop.Count: got %d, want 0", got)
	}
	if got := empty.Distinct(func(a, b int) bool { return a == b }).Count(); got != 0 {
		t.Fatalf("Empty.Distinct.Count: got %d, want 0", got)
	}
	if got := empty.Peek(func(int) {}).Count(); got != 0 {
		t.Fatalf("Empty.Peek.Count: got %d, want 0", got)
	}
	if got := empty.Concat(stream.Of(1, 2, 3)).Count(); got != 3 {
		t.Fatalf("Empty.Concat(Of).Count: got %d, want 3", got)
	}
	if got := stream.Of(1, 2, 3).Concat(empty).Count(); got != 3 {
		t.Fatalf("Of.Concat(Empty).Count: got %d, want 3", got)
	}
	if got := empty.Concat(empty).Count(); got != 0 {
		t.Fatalf("Empty.Concat(Empty).Count: got %d, want 0", got)
	}
}

// TestStreamOfAndFrom cover the remaining constructors that were
// exercised only indirectly. Of(values...) is the variadic sibling of
// FromSlice; From(slices.Values) builds from a raw iter.Seq.
func TestStreamOfAndFrom(t *testing.T) {
	t.Run("Of", func(t *testing.T) {
		got := stream.Of(10, 20, 30).Collect()
		if want := []int{10, 20, 30}; !equal(got, want) {
			t.Fatalf("Of.Collect: got %v, want %v", got, want)
		}
	})

	t.Run("From raw iter.Seq", func(t *testing.T) {
		seq := func(yield func(int) bool) {
			for i := 1; i <= 3; i++ {
				if !yield(i * 100) {
					return
				}
			}
		}
		got := stream.From(seq).Collect()
		if want := []int{100, 200, 300}; !equal(got, want) {
			t.Fatalf("From.Collect: got %v, want %v", got, want)
		}
	})
}

// TestStreamPeek verifies that Peek observes every element in order,
// does not modify the yielded value, and that the visit function is
// not called when the source yields nothing.
func TestStreamPeek(t *testing.T) {
	var observed []int
	src := stream.Of(1, 2, 3, 4).Peek(func(n int) { observed = append(observed, n*10) })

	got := src.Collect()
	if want := []int{1, 2, 3, 4}; !equal(got, want) {
		t.Fatalf("Peek must not transform the values: got %v, want %v", got, want)
	}
	if want := []int{10, 20, 30, 40}; !equal(observed, want) {
		t.Fatalf("Peek observed: got %v, want %v", observed, want)
	}

	// Peek is an intermediate; further ops should still see the original
	// values, not the side-effect arguments.
	sum := stream.Of(1, 2, 3, 4).
		Peek(func(int) {}).
		Reduce(0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Fatalf("Peek+Reduce sum: got %d, want 10", sum)
	}
}

// TestStreamTake covers n>0, n<0, n==0, n smaller than the source,
// n equal to the source, and n larger than the source.
func TestStreamTake(t *testing.T) {
	src := stream.Of(1, 2, 3, 4, 5)

	t.Run("n smaller", func(t *testing.T) {
		if got, want := src.Take(3).Collect(), []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n equal", func(t *testing.T) {
		if got, want := src.Take(5).Collect(), []int{1, 2, 3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n larger", func(t *testing.T) {
		if got, want := src.Take(100).Collect(), []int{1, 2, 3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n zero", func(t *testing.T) {
		if got := src.Take(0).Count(); got != 0 {
			t.Fatalf("Take(0).Count: got %d, want 0", got)
		}
	})
	t.Run("n negative", func(t *testing.T) {
		if got := src.Take(-7).Count(); got != 0 {
			t.Fatalf("Take(-7).Count: got %d, want 0", got)
		}
	})
}

// TestStreamTakeStopsSourceEarly uses a custom iter.Seq that records
// how many values it produced. Take must stop the source as soon as
// it has enough elements, not walk the whole sequence.
func TestStreamTakeStopsSourceEarly(t *testing.T) {
	produced := 0
	infinite := func(yield func(int) bool) {
		for i := 0; ; i++ {
			produced++
			if !yield(i) {
				return
			}
		}
	}

	if got := stream.From(infinite).Take(3).Collect(); !equal(got, []int{0, 1, 2}) {
		t.Fatalf("Take(3) on infinite: got %v, want [0 1 2]", got)
	}
	if produced != 3 {
		t.Fatalf("Take(3) produced %d values from the source, want 3", produced)
	}
}

// TestStreamDrop covers n smaller than, equal to, and larger than the
// source, plus the n<=0 fast path (which returns s itself).
func TestStreamDrop(t *testing.T) {
	src := stream.Of(1, 2, 3, 4, 5)

	t.Run("n smaller", func(t *testing.T) {
		if got, want := src.Drop(2).Collect(), []int{3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n equal", func(t *testing.T) {
		if got := src.Drop(5).Count(); got != 0 {
			t.Fatalf("Drop(5).Count: got %d, want 0", got)
		}
	})
	t.Run("n larger", func(t *testing.T) {
		if got := src.Drop(100).Count(); got != 0 {
			t.Fatalf("Drop(100).Count: got %d, want 0", got)
		}
	})
	t.Run("n zero returns full stream", func(t *testing.T) {
		if got, want := src.Drop(0).Collect(), []int{1, 2, 3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("n negative returns full stream", func(t *testing.T) {
		if got, want := src.Drop(-3).Collect(), []int{1, 2, 3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

// TestStreamDistinct covers the dedup path with a custom eq function.
// It exercises all the relevant shapes: no duplicates, all duplicates,
// mixed, and stream with a single element.
func TestStreamDistinct(t *testing.T) {
	eq := func(a, b int) bool { return a == b }

	t.Run("no duplicates", func(t *testing.T) {
		got := stream.Of(1, 2, 3, 4).Distinct(eq).Collect()
		if want := []int{1, 2, 3, 4}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("all duplicates", func(t *testing.T) {
		got := stream.Of(7, 7, 7, 7).Distinct(eq).Collect()
		if want := []int{7}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("first occurrence wins", func(t *testing.T) {
		got := stream.Of(1, 2, 1, 3, 2, 1).Distinct(eq).Collect()
		if want := []int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("custom eq on structs", func(t *testing.T) {
		type pt struct{ x, y int }
		eqPt := func(a, b pt) bool { return a.x == b.x && a.y == b.y }
		in := []pt{{1, 1}, {1, 2}, {2, 1}, {1, 1}}
		want := []pt{{1, 1}, {1, 2}, {2, 1}}
		got := stream.FromSlice(in).Distinct(eqPt).Collect()
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})
}

// TestStreamConcat covers the four matrix cells of (empty, populated)
// and also verifies that a downstream early-termination (Take) stops
// the right side of the concat from being walked.
func TestStreamConcat(t *testing.T) {
	t.Run("empty + populated", func(t *testing.T) {
		if got, want := stream.Empty[int]().Concat(stream.Of(1, 2, 3)).Collect(),
			[]int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + empty", func(t *testing.T) {
		if got, want := stream.Of(1, 2, 3).Concat(stream.Empty[int]()).Collect(),
			[]int{1, 2, 3}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("populated + populated", func(t *testing.T) {
		if got, want := stream.Of(1, 2).Concat(stream.Of(3, 4, 5)).Collect(),
			[]int{1, 2, 3, 4, 5}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("early termination stops the second source", func(t *testing.T) {
		secondVisited := 0
		second := func(yield func(int) bool) {
			for i := 100; ; i++ {
				secondVisited++
				if !yield(i) {
					return
				}
			}
		}
		got := stream.Of(1, 2, 3).Concat(stream.From(second)).Take(4).Collect()
		if want := []int{1, 2, 3, 100}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		if secondVisited != 1 {
			t.Fatalf("second source visited %d times, want 1", secondVisited)
		}
	})
	t.Run("early termination stops the first source", func(t *testing.T) {
		firstVisited := 0
		first := func(yield func(int) bool) {
			for i := 0; ; i++ {
				firstVisited++
				if !yield(i) {
					return
				}
			}
		}
		got := stream.From(first).Concat(stream.Of(99, 100)).Take(2).Collect()
		if want := []int{0, 1}; !equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		if firstVisited != 2 {
			t.Fatalf("first source visited %d times, want 2", firstVisited)
		}
	})
}

// TestStreamCount covers both empty and non-empty sources.
func TestStreamCount(t *testing.T) {
	if got := stream.Empty[int]().Count(); got != 0 {
		t.Fatalf("Empty.Count: got %d, want 0", got)
	}
	if got := stream.Of(1, 2, 3, 4, 5).Count(); got != 5 {
		t.Fatalf("Of(5).Count: got %d, want 5", got)
	}
	// Count after intermediate operations must reflect the post-op size.
	if got := stream.Of(1, 2, 3, 4, 5).Filter(func(n int) bool { return n > 2 }).Count(); got != 3 {
		t.Fatalf("Filter.Count: got %d, want 3", got)
	}
}

// TestStreamLast covers the empty, single-element, and multi-element
// cases. Last is non-short-circuiting and must walk the full source.
func TestStreamLast(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if stream.Empty[int]().Last().IsPresent() {
			t.Fatal("Empty.Last: expected absent")
		}
	})
	t.Run("single", func(t *testing.T) {
		if got := stream.Of(42).Last().OrElse(0); got != 42 {
			t.Fatalf("Of(42).Last: got %d, want 42", got)
		}
	})
	t.Run("multi", func(t *testing.T) {
		if got := stream.Of(1, 2, 3, 4, 5).Last().OrElse(0); got != 5 {
			t.Fatalf("Of(1..5).Last: got %d, want 5", got)
		}
	})
	t.Run("chained", func(t *testing.T) {
		// Last on a filtered pipeline must reflect the last value the
		// filter lets through, not the last value of the source.
		got := stream.Of(1, 2, 3, 4, 5).
			Filter(func(n int) bool { return n%2 == 0 }).
			Last().OrElse(0)
		if got != 4 {
			t.Fatalf("filtered Last: got %d, want 4", got)
		}
	})
}

// TestStreamAll covers the short-circuiting positive and negative
// paths plus the vacuous-true-on-empty case.
func TestStreamAll(t *testing.T) {
	t.Run("all match", func(t *testing.T) {
		if !stream.Of(2, 4, 6, 8).All(func(n int) bool { return n%2 == 0 }) {
			t.Fatal("All: expected true")
		}
	})
	t.Run("mismatch at start short-circuits", func(t *testing.T) {
		visited := 0
		// All must stop as soon as one mismatch is found.
		result := stream.Of(1, 2, 3, 4, 5).
			Peek(func(int) { visited++ }).
			All(func(n int) bool { return n > 0 })
		// All elements are > 0, so it must walk the full stream and return true.
		if !result {
			t.Fatal("All: expected true for n > 0")
		}
		if visited != 5 {
			t.Fatalf("Peek visited %d values, want 5", visited)
		}

		// Now the short-circuit: a mismatch at the second element must
		// stop consumption after two elements.
		visited = 0
		result = stream.Of(1, 2, 3, 4, 5).
			Peek(func(int) { visited++ }).
			All(func(n int) bool { return n < 3 })
		if result {
			t.Fatal("All: expected false for n < 3")
		}
		if visited != 3 {
			t.Fatalf("All short-circuit visited %d values, want 3", visited)
		}
	})
	t.Run("empty is vacuously true", func(t *testing.T) {
		if !stream.Empty[int]().All(func(int) bool { return false }) {
			t.Fatal("Empty.All: expected vacuously true")
		}
	})
}

// TestStreamNone covers both the truthy (no element matches) and
// falsy (at least one element matches) outcomes. None is the dual
// of Any; both are tested here so the symmetry is documented.
func TestStreamNone(t *testing.T) {
	t.Run("no match", func(t *testing.T) {
		if !stream.Of(1, 2, 3).None(func(n int) bool { return n > 100 }) {
			t.Fatal("None: expected true when nothing matches")
		}
	})
	t.Run("match exists", func(t *testing.T) {
		if stream.Of(1, 2, 3).None(func(n int) bool { return n == 2 }) {
			t.Fatal("None: expected false when something matches")
		}
	})
	t.Run("empty is vacuously true", func(t *testing.T) {
		if !stream.Empty[int]().None(func(int) bool { return true }) {
			t.Fatal("Empty.None: expected true")
		}
	})
	t.Run("none short-circuits via Any", func(t *testing.T) {
		// None uses !Any(p) internally; once Any finds a match it
		// returns true, so None must also stop early. We verify via
		// a Peek counter on an infinite source.
		visited := 0
		infinite := func(yield func(int) bool) {
			for i := 0; ; i++ {
				visited++
				if !yield(i) {
					return
				}
			}
		}
		// None(n == 5) is false because 5 matches; it must stop
		// after i==5 has been produced (six values: 0..5).
		if stream.From(infinite).None(func(n int) bool { return n == 5 }) {
			t.Fatal("None: expected false")
		}
		if visited != 6 {
			t.Fatalf("None visited %d values from infinite source, want 6", visited)
		}
	})
}

// TestStreamFind covers both the found and not-found paths. Find is
// short-circuiting: once a match is found, the underlying source
// must stop.
func TestStreamFind(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		if got := stream.Of(10, 20, 30, 40).Find(func(n int) bool { return n >= 30 }).OrElse(0); got != 30 {
			t.Fatalf("Find: got %d, want 30", got)
		}
	})
	t.Run("first element matches", func(t *testing.T) {
		if got := stream.Of(7, 8, 9).Find(func(int) bool { return true }).OrElse(0); got != 7 {
			t.Fatalf("Find first: got %d, want 7", got)
		}
	})
	t.Run("not found", func(t *testing.T) {
		if stream.Of(1, 2, 3).Find(func(n int) bool { return n > 100 }).IsPresent() {
			t.Fatal("Find: expected absent when no match")
		}
	})
	t.Run("empty source", func(t *testing.T) {
		if stream.Empty[int]().Find(func(int) bool { return true }).IsPresent() {
			t.Fatal("Empty.Find: expected absent")
		}
	})
	t.Run("short-circuits source", func(t *testing.T) {
		visited := 0
		infinite := func(yield func(int) bool) {
			for i := 0; ; i++ {
				visited++
				if !yield(i) {
					return
				}
			}
		}
		if got := stream.From(infinite).Find(func(n int) bool { return n == 3 }).OrElse(-1); got != 3 {
			t.Fatalf("Find on infinite: got %d, want 3", got)
		}
		if visited != 4 { // i = 0, 1, 2, 3
			t.Fatalf("Find visited %d values, want 4", visited)
		}
	})
}

// TestStreamReduce covers same-type accumulation (the only signature
// the Stream method exposes: init T, f func(acc, v T) T), the empty
// source (returns init), and a single-element source.
func TestStreamReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		got := stream.Of(1, 2, 3, 4, 5).Reduce(0, func(a, v int) int { return a + v })
		if got != 15 {
			t.Fatalf("Reduce sum: got %d, want 15", got)
		}
	})
	t.Run("empty returns init", func(t *testing.T) {
		got := stream.Empty[int]().Reduce(99, func(a, v int) int { return a + v })
		if got != 99 {
			t.Fatalf("Reduce on empty: got %d, want 99 (init)", got)
		}
	})
	t.Run("single element", func(t *testing.T) {
		got := stream.Of(42).Reduce(0, func(a, v int) int { return a + v })
		if got != 42 {
			t.Fatalf("Reduce single: got %d, want 42", got)
		}
	})
	t.Run("fold with a string source", func(t *testing.T) {
		// Reduce works on any T, not just numbers. Build a joined
		// string from a stream of strings.
		joined := stream.Of("a", "b", "c").Reduce("", func(a, v string) string {
			if a == "" {
				return v
			}
			return a + "," + v
		})
		if joined != "a,b,c" {
			t.Fatalf("Reduce strings: got %q, want %q", joined, "a,b,c")
		}
	})
}

// TestStreamForEach covers empty (no calls), non-empty (every element
// visited in order), and a downstream chain that observes the same
// values ForEach was supposed to consume.
func TestStreamForEach(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		called := 0
		stream.Empty[int]().ForEach(func(int) { called++ })
		if called != 0 {
			t.Fatalf("ForEach on empty: called %d times, want 0", called)
		}
	})
	t.Run("order", func(t *testing.T) {
		var got []int
		stream.Of(5, 4, 3, 2, 1).ForEach(func(n int) { got = append(got, n) })
		if want := []int{5, 4, 3, 2, 1}; !equal(got, want) {
			t.Fatalf("ForEach order: got %v, want %v", got, want)
		}
	})
}

// TestStreamFirstEmpty covers the empty-source branch of First: the
// underlying for-range body never runs and the function returns
// option.Empty[T]().
func TestStreamFirstEmpty(t *testing.T) {
	if stream.Empty[int]().First().IsPresent() {
		t.Fatal("First on empty: expected absent")
	}
	// Of() with no args has no inferable T, so it cannot be called
	// directly; use an explicit type argument.
	if stream.Of[int]().First().IsPresent() {
		t.Fatal("First on stream.Of[int](): expected absent")
	}
}

// TestStreamAnyFalse covers the no-match branch of Any: the for-range
// body runs to completion without finding a match and the function
// returns false.
func TestStreamAnyFalse(t *testing.T) {
	if stream.Of(1, 2, 3).Any(func(n int) bool { return n > 100 }) {
		t.Fatal("Any: expected false when nothing matches")
	}
	if stream.Empty[int]().Any(func(int) bool { return true }) {
		t.Fatal("Any on empty: expected false")
	}
}

// TestStreamSortByEmpty covers the empty-source branch of SortBy.
// The inner Collect() returns an empty slice, slices.SortFunc is a
// no-op, and the resulting Stream yields nothing.
func TestStreamSortByEmpty(t *testing.T) {
	got := stream.Empty[int]().SortBy(func(a, b int) int { return a - b }).Collect()
	if len(got) != 0 {
		t.Fatalf("SortBy on empty: got %v, want []", got)
	}
}

// TestStreamMinByMaxByEmpty covers the empty-source branches of
// MinBy and MaxBy. The for-range body never runs, found stays false,
// and the function returns option.Empty[T]().
func TestStreamMinByMaxByEmpty(t *testing.T) {
	if stream.Empty[int]().MinBy(func(a, b int) int { return a - b }).IsPresent() {
		t.Fatal("MinBy on empty: expected absent")
	}
	if stream.Empty[int]().MaxBy(func(a, b int) int { return a - b }).IsPresent() {
		t.Fatal("MaxBy on empty: expected absent")
	}

	// After filtering down to an empty result, MinBy / MaxBy must
	// also return absent.
	if stream.Of(1, 2, 3).Filter(func(n int) bool { return n > 100 }).MinBy(func(a, b int) int { return a - b }).IsPresent() {
		t.Fatal("MinBy on filtered-empty: expected absent")
	}
	if stream.Of(1, 2, 3).Filter(func(n int) bool { return n > 100 }).MaxBy(func(a, b int) int { return a - b }).IsPresent() {
		t.Fatal("MaxBy on filtered-empty: expected absent")
	}
}

// TestStreamFilterEarlyTermination covers the downstream-short-circuit
// branch of Filter: when the consumer of the filtered stream signals
// "stop" via yield returning false, the underlying source must also
// stop. This is the 20% gap in the Filter coverage report.
func TestStreamFilterEarlyTermination(t *testing.T) {
	produced := 0
	src := func(yield func(int) bool) {
		for i := 0; ; i++ {
			produced++
			if !yield(i) {
				return
			}
		}
	}

	// Take(2) after Filter must stop the source after two elements
	// have been yielded by Filter — which means at most two elements
	// are produced by the source.
	got := stream.From(src).Filter(func(n int) bool { return true }).Take(2).Collect()
	if want := []int{0, 1}; !equal(got, want) {
		t.Fatalf("Filter+Take: got %v, want %v", got, want)
	}
	if produced != 2 {
		t.Fatalf("Filter+Take produced %d source values, want 2", produced)
	}
}

// TestStreamMapEarlyTermination covers the same downstream stop
// branch for Map. The 25% gap in Map's coverage is exactly this
// branch.
func TestStreamMapEarlyTermination(t *testing.T) {
	produced := 0
	src := func(yield func(int) bool) {
		for i := 0; ; i++ {
			produced++
			if !yield(i) {
				return
			}
		}
	}

	got := stream.From(src).Map(func(n int) int { return n * 2 }).Take(3).Collect()
	if want := []int{0, 2, 4}; !equal(got, want) {
		t.Fatalf("Map+Take: got %v, want %v", got, want)
	}
	if produced != 3 {
		t.Fatalf("Map+Take produced %d source values, want 3", produced)
	}
}

// TestStreamFlatMapEarlyTerminationAndEmptyInner covers both gaps in
// FlatMap: the downstream short-circuit branch and the case where
// the inner stream is empty.
func TestStreamFlatMapEarlyTerminationAndEmptyInner(t *testing.T) {
	t.Run("downstream early termination", func(t *testing.T) {
		produced := 0
		src := func(yield func(int) bool) {
			for i := 0; ; i++ {
				produced++
				if !yield(i) {
					return
				}
			}
		}
		// Each outer element produces one inner element. Take(2)
		// must stop the outer source after two outer elements.
		got := stream.From(src).FlatMap(func(n int) stream.Stream[int] {
			return stream.Of(n)
		}).Take(2).Collect()
		if want := []int{0, 1}; !equal(got, want) {
			t.Fatalf("FlatMap+Take: got %v, want %v", got, want)
		}
		if produced != 2 {
			t.Fatalf("FlatMap+Take produced %d outer values, want 2", produced)
		}
	})

	t.Run("empty inner stream is a no-op", func(t *testing.T) {
		// f returns an empty Stream for every element: the result
		// must be empty regardless of the source.
		got := stream.Of(1, 2, 3).FlatMap(func(int) stream.Stream[int] {
			return stream.Empty[int]()
		}).Collect()
		if len(got) != 0 {
			t.Fatalf("FlatMap with empty inner: got %v, want []", got)
		}
	})

	t.Run("inner early-termination short-circuits the inner", func(t *testing.T) {
		// Outer source yields 3 elements, but the inner stream
		// yields 5 and a downstream Take(7) must stop the inner
		// of the second outer element after one more value.
		innerProduced := 0
		innerFactory := func(int) stream.Stream[int] {
			return stream.From(func(yield func(int) bool) {
				for j := 0; j < 5; j++ {
					innerProduced++
					if !yield(j) {
						return
					}
				}
			})
		}
		got := stream.Of(1, 2, 3).FlatMap(innerFactory).Take(7).Collect()
		if want := []int{0, 1, 2, 3, 4, 0, 1}; !equal(got, want) {
			t.Fatalf("FlatMap inner early-stop: got %v, want %v", got, want)
		}
		// We expected 5 (outer #0) + 2 (outer #1, stopped at index 1) = 7
		// inner yields. 5+2 == 7, and produced records the third
		// element being never reached.
		if innerProduced != 7 {
			t.Fatalf("inner produced %d, want 7", innerProduced)
		}
	})
}

// TestStreamAnyStopsAfterMatch complements the TestStreamEarlyTermination
// coverage in the existing test by also covering the false branch via
// an infinite source — Any must walk the whole source if nothing matches.
func TestStreamAnyWalksAllWhenNoMatch(t *testing.T) {
	visited := 0
	infinite := func(yield func(int) bool) {
		for i := 0; i < 1000; i++ {
			visited++
			if !yield(i) {
				return
			}
		}
	}
	if stream.From(infinite).Any(func(n int) bool { return n == 999999 }) {
		t.Fatal("Any: expected false for a value beyond the source")
	}
	if visited != 1000 {
		t.Fatalf("Any visited %d values, want 1000", visited)
	}
}

// TestStreamEmptyPreservesType verifies the type parameter on Empty[T].
// T must be carried through; downstream Map must be able to call a
// function with that exact T.
func TestStreamEmptyPreservesType(t *testing.T) {
	type point struct{ x, y int }
	got := stream.Empty[point]().Map(func(p point) int { return p.x + p.y }).Collect()
	if len(got) != 0 {
		t.Fatalf("Empty[point].Map.Collect: got %v, want []", got)
	}
}

// TestStreamIntermediateEarlyTermination covers the downstream-stop
// branch inside each intermediate that materializes a downstream
// yield returning false. The pattern is the same in every case:
// chain the intermediate with Take(N) where N is smaller than the
// upstream would produce, so the inner `yield` callback returns
// false and the intermediate's early-exit return fires.
//
// Without these sub-tests, the 6 remaining uncovered statements are
// the `return` lines inside the `if !yield(v) { ... }` blocks of
// Take, Drop, Distinct, and SortBy.
func TestStreamIntermediateEarlyTermination(t *testing.T) {
	t.Run("Take stops when downstream stops", func(t *testing.T) {
		// Take(10) on a 5-element source would walk the whole source
		// without an early-exit. Chaining Take(2) on the outside
		// forces the outer Take's yield to return false on the 3rd
		// element, firing its early-exit branch.
		produced := 0
		src := func(yield func(int) bool) {
			for i := 0; i < 5; i++ {
				produced++
				if !yield(i) {
					return
				}
			}
		}
		got := stream.From(src).Take(10).Take(2).Collect()
		if want := []int{0, 1}; !equal(got, want) {
			t.Fatalf("Take(Take): got %v, want %v", got, want)
		}
		if produced != 2 {
			t.Fatalf("source produced %d, want 2", produced)
		}
	})

	t.Run("Drop stops when downstream stops", func(t *testing.T) {
		// Drop(1) on [0,1,2,3,4] would normally yield [1,2,3,4].
		// Chaining Take(2) forces Drop's yield to return false
		// after [1, 2], firing its early-exit branch.
		produced := 0
		src := func(yield func(int) bool) {
			for i := 0; i < 5; i++ {
				produced++
				if !yield(i) {
					return
				}
			}
		}
		got := stream.From(src).Drop(1).Take(2).Collect()
		if want := []int{1, 2}; !equal(got, want) {
			t.Fatalf("Drop(Take): got %v, want %v", got, want)
		}
		if produced != 3 {
			t.Fatalf("source produced %d, want 3 (0 dropped, 1, 2 yielded)", produced)
		}
	})

	t.Run("Distinct stops when downstream stops", func(t *testing.T) {
		// Distinct on [1,1,2,2,3,3] would normally yield [1,2,3].
		// Chaining Take(2) forces Distinct's yield to return false
		// after the 2nd element, firing its early-exit branch.
		produced := 0
		src := func(yield func(int) bool) {
			for _, v := range []int{1, 1, 2, 2, 3, 3} {
				produced++
				if !yield(v) {
					return
				}
			}
		}
		got := stream.From(src).Distinct(func(a, b int) bool { return a == b }).Take(2).Collect()
		if want := []int{1, 2}; !equal(got, want) {
			t.Fatalf("Distinct(Take): got %v, want %v", got, want)
		}
		// 1, 1 (dup), 2, (stop) — produced 3 source values.
		if produced != 3 {
			t.Fatalf("source produced %d, want 3", produced)
		}
	})

	t.Run("SortBy stops when downstream stops", func(t *testing.T) {
		// SortBy on [3,1,4,1,5,9,2,6] would yield [1,1,2,3,4,5,6,9]
		// after sorting. Chaining Take(3) forces SortBy's yield to
		// return false after the 3rd sorted element, firing its
		// early-exit branch.
		produced := 0
		src := func(yield func(int) bool) {
			for _, v := range []int{3, 1, 4, 1, 5, 9, 2, 6} {
				produced++
				if !yield(v) {
					return
				}
			}
		}
		got := stream.From(src).
			SortBy(func(a, b int) int { return a - b }).
			Take(3).Collect()
		if want := []int{1, 1, 2}; !equal(got, want) {
			t.Fatalf("SortBy(Take): got %v, want %v", got, want)
		}
		// SortBy materializes the entire source first, so the source
		// must be fully consumed even if Take(3) cuts the output.
		if produced != 8 {
			t.Fatalf("source produced %d, want 8 (SortBy is eager on Collect)", produced)
		}
	})
}
