# `typed/collections` — Typed

[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=collections)](https://codecov.io/gh/qianwj/typed)

Generic collections and a synchronous data-flow layer. `collections` ships four families of containers and one lazy stream tool:

- `Stack[T]` / `Queue[T]` / `Deque[T]` — basic linear containers; "take one" returns `option.Option[T]`.
- `queues.PriorityQueue[T]` — min-heap-backed priority queue with `Push` / `PushWithPriority` / `Pop` / `Peek` / `Len`.
- `lists.ArrayList[T]` / `lists.LinkedList[T]` — lists with immutable-style transforms (`Filter` / `Map` / `Take` / `Drop` / `Concat` / `Distinct` / `SortBy`).
- `maps.HashMap[K, V]` — hash table with `Keys` / `Values` / `Entries` / `Filter*` / `MapValues` / `Concat`.
- `sets.HashSet[T]` — hash set with set algebra (`Union` / `Intersect` / `Difference` / `SymmetricDifference`) and same-kind transforms.
- `stream.Stream[T]` — a lazy stream over `iter.Seq[T]`, chainable with `Filter` / `Map` / `FlatMap` / `Take` / `Drop` / `Distinct` / `Concat` / `SortBy` / `Reduce` / `Count` / `Find` and other terminal operations.
- `collections.Range[T]` — integer half-open interval as a `Stream[T]` factory.

Every collection implements `MarshalJSON` / `UnmarshalJSON`. The element type only needs to satisfy `any`; the JSON form is Go-style arrays or objects.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [`iterable.Iterable[T]`](#iterableiterablet)
- [`Stack[T]` / `Queue[T]` / `Deque[T]`](#stackt--queuet--dequet)
- [`queues.PriorityQueue[T]`](#queuespriorityqueuet)
- [`lists.ArrayList[T]`](#listsarraylistt)
- [`lists.LinkedList[T]`](#listslinkedlistt)
- [`maps.HashMap[K, V]`](#mapshashmapk-v)
- [`sets.HashSet[T]`](#setshashsett)
- [`stream.Stream[T]`](#streamstreamt)
- [`collections.Range[T]`](#collectionsranget)
- [See also](#see-also)

## Import

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/iterable"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/collections/maps"
    "github.com/qianwj/typed/collections/queues"
    "github.com/qianwj/typed/collections/sets"
    "github.com/qianwj/typed/collections/stream"
)
```

The four sub-packages are independent `go.mod` modules; import what you use.

---

## `iterable.Iterable[T]`

The `collections/iterable` package defines the shared traversal contract:

```go
type Iterable[T any] interface {
    ForEach(visit func(T))
}
```

`ArrayList`, `LinkedList`, `HashSet`, and `Stream` already satisfy it through their existing `ForEach` methods. The interface lets constructors consume any of these types, or your own implementation, without dependencies between concrete collection types.

| Constructor | Behavior |
|---|---|
| `lists.ArrayListFrom[T](source iterable.Iterable[T])` | Creates a new `*ArrayList[T]`, preserving source order and duplicates. |
| `lists.LinkedListFrom[T](source iterable.Iterable[T])` | Creates a new `*LinkedList[T]`, preserving source order and duplicates. |
| `sets.HashSetFrom[T comparable](source iterable.Iterable[T])` | Creates a new `*HashSet[T]`, eliminating duplicates. |

Each constructor consumes the source once into independent storage; elements are shallow-copied. Converting a set to a list does not establish an order because set iteration is unspecified. Passing a Stream consumes that stream.

```go
list := lists.ArrayListOf(3, 1, 3, 2)
set := sets.HashSetFrom(list)
array := lists.ArrayListFrom(set)
linked := lists.LinkedListFrom(array)
```

---

## `Stack[T]` / `Queue[T]` / `Deque[T]`

All three are constructed as struct pointers, returning `*Stack[T]` and the like. Every "take one" operation returns `option.Option[T]` rather than `(T, bool)`, so it composes naturally with the `Stream` / `Result` chains.

```go
s := collections.NewStack[int]()
s.Push(1); s.Push(2)
v, ok := s.Peek().Get()        // 2, true
top, _ := s.Pop().Get()        // 2
_, present := s.Pop().Get()    // 1, true
s.Pop().IsEmpty()              // false (1 still inside)
s.Pop().IsEmpty()              // true
s.Pop()                        // option.Empty[int](), no panic
```

| Type | Construct | Key methods |
|---|---|---|
| `Stack[T]` | `NewStack[T]()` | `Push(v) / Pop() Option[T] / Peek() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |
| `Queue[T]` | `NewQueue[T]()` | `Push(v) / Pop() Option[T] / Peek() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |
| `Deque[T]` | `NewDeque[T]()` | `PushFront(v) / PushBack(v) / PopFront() Option[T] / PopBack() Option[T] / Front() Option[T] / Back() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |

`Pop*` / `Peek*` on an empty container return `option.Empty[T]()`: they do not panic, and `Peek` / `Front` / `Back` leave the container unchanged, while `Pop*` mutates per the container rules.

---

## `queues.PriorityQueue[T]`

`NewPriorityQueue[T any]() *PriorityQueue[T]` returns an empty priority queue backed by a min-heap (slice-based, in-place sift). Use it to order work items by priority within a single goroutine.

```go
import "github.com/qianwj/typed/collections/queues"

pq := queues.NewPriorityQueue[int]()
pq.PushWithPriority(42, 0)
pq.PushWithPriority(7,  -1) // highest priority
pq.PushWithPriority(99, 5)

for pq.Len() > 0 {
    fmt.Println(pq.Pop().Get()) // 7, 42, 99
}
```

| Method | Behaviour |
| --- | --- |
| `NewPriorityQueue[T]()` | Construct an empty queue. |
| `Push(data T)` | Add with priority 0. Equivalent to `PushWithPriority(data, 0)`. |
| `PushWithPriority(data T, priority int)` | Add with the given priority. Lower priority values are dequeued first. |
| `Pop() adt.Option[T]` | Remove and return the highest-priority element, or `Empty[T]()` if the queue is empty. O(log n). |
| `Peek() adt.Option[T]` | Return the highest-priority element without removing it, or `Empty[T]()` if empty. O(1). |
| `Len() int` | Current size. |

**Tie-breaking.** The heap uses strict less-than comparison; equal-priority elements have no guaranteed relative order. If you need FIFO (or any stable) ordering among ties, encode the tiebreaker in the priority itself — e.g., `priority = scheduledAt.UnixNano()`.

**Concurrency.** `PriorityQueue` is synchronous (no internal locking) and lives in the `collections` package — the same single-goroutine contract as `Queue[T]` and `Stack[T]`. For cross-goroutine use, wrap with a `sync.Mutex` or feed it through a `concurrency.Group`. The roadmap originally sketched `Push` / `Poll` / `TryPoll` (the `BoundedBlockingQueue` verbs), but a synchronous container has no natural blocking `Poll`, so the API mirrors `collections.Queue` instead — `Push` + `Pop` (instead of `TryPoll`).

**Memory model.** Backed by a single `[]priorityItem` that grows via `append` on every push. Freed slots are zeroed in `Pop` so a pointer-typed `T` is not pinned in the backing array after removal, matching `collections.Stack` and `collections.Queue`.

---

## `lists.ArrayList[T]`

`NewArrayList[T any]() *ArrayList[T]` constructs an empty list. `ArrayList` is a single contiguous `[]T` with a `head` offset; `Take` / `Drop` only move `head` and do not copy elements (`O(1)`); only `Collect` actually copies into a new slice.

### Read / write

| Method | Description |
|---|---|
| `Add(v)` / `AddFirst(v)` | Append to the tail / insert at the head. |
| `Insert(i, v)` | Insert at index `i`; `i < 0` or `i > Size()` panics. |
| `Get(i) Option[T]` | Out-of-range returns `Empty`, no panic. |
| `First() Option[T]` / `Last() Option[T]` | Empty list returns `Empty`. |
| `RemoveAt(i) T` | Remove the element at index `i` and return it; out-of-range panics. |
| `RemoveFirst() Option[T]` / `RemoveLast() Option[T]` | Empty list returns `Empty`. |
| `Size() / IsEmpty() / Clear() / Collect() []T` | Cardinality and export. |
| `Stream() stream.Stream[T]` | Convert to a lazy `Stream[T]` (no data copy). |
| `MarshalJSON / UnmarshalJSON` | JSON array. |

### Immutable-style transforms (return a new list; the receiver is not modified)

| Method | Description |
|---|---|
| `Filter(p func(T) bool) *ArrayList[T]` | Keep the elements satisfying `p`. |
| `Map[R](f func(T) R) *ArrayList[R]` | Element type `T → R`. |
| `FlatMap[R](f func(T) *ArrayList[R]) *ArrayList[R]` | Concatenate the sub-lists produced by `f`. |
| `Take(n) / Drop(n)` | Prefix / suffix, `O(1)`. |
| `Distinct(eq func(T, T) bool)` | Dedup by `eq`. |
| `Concat(other)` | Append. |
| `Peek(visit func(T))` | Walk without changing the list; returns the receiver for chaining, useful for inline debugging. |
| `SortBy(less func(x, y T) int) *ArrayList[T]` | Sort by `less`; returns a new list. |
| `MinBy(less func(x, y T) int) option.Option[T]` | Empty list returns `option.Empty[T]()`; otherwise returns the smallest element (present). |
| `MaxBy(less func(x, y T) int) option.Option[T]` | Same as above, returns the largest element. |

### Predicates

`Any(p) / All(p) / None(p) / Find(p) Option[T]` — short-circuiting universal / existential quantifiers; `Find` returns the first matching element.

### Examples

```go
xs := lists.NewArrayList[int]()
xs.Add(3); xs.Add(1); xs.Add(4); xs.Add(1); xs.Add(5)
_ = xs
```

> `Add` is not designed to chain by returning the receiver; if you want chaining, use `Peek` or start from a `Stream`.

```go
even := lists.NewArrayList[int]()
even.Add(2); even.Add(4); even.Add(6)
sorted := even.SortBy(func(a, b int) int { return a - b })
max := sorted.MaxBy(func(a, b int) int { return a - b }).OrElse(0) // 6
```

---

## `lists.LinkedList[T]`

Doubly-linked list with sentinel nodes. `NewLinkedList[T any]() *LinkedList[T]` constructs.

### Read / write

| Method | Description |
|---|---|
| `Add(v)` / `AddFirst(v)` | Append / insert at the head. |
| `Insert(i, v)` | Insert at index `i`; out-of-range panics. |
| `Get(i) Option[T]` | Out-of-range returns `Empty`. |
| `First() / Last() Option[T]` | Empty list returns `Empty`. |
| `RemoveAt(i) T` | Out-of-range panics. |
| `RemoveFirst() / RemoveLast() Option[T]` | Empty list returns `Empty`. |
| `Size() / IsEmpty() / Clear() / Collect() []T` | Cardinality and export. |
| `Stream() stream.Stream[T]` | Lazy `Stream[T]`. |
| `MarshalJSON / UnmarshalJSON` | JSON array. |

### Immutable-style transforms

`Filter(p) *LinkedList[T]` / `Map[R](f) *ArrayList[R]` / `FlatMap[R](f) *ArrayList[R]` / `Take(n) / Drop(n) / Distinct(eq) / Concat(other) / Peek(visit) / SortBy(less) *LinkedList[T]`.

`Map` and `FlatMap` on a linked list return `*ArrayList[R]` because downstream consumers usually want index or slice access; the conversion is a single `O(n)` walk.

`MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` — empty list returns `Empty`, otherwise the extreme.

### Predicates

`Any / All / None / Find(p) Option[T]`, plus `Reduce[U](init U, f func(U, T) U) U`.

---

## `maps.HashMap[K, V]`

`K` must be `comparable`, `V` is any. `NewHashMap[K, V]() *HashMap[K, V]` constructs.

### Basics

| Method | Description |
|---|---|
| `Put(k, v) V` | Set; returns the previous value (zero if none). |
| `PutIfAbsent(k, v)` | Set only when the key is absent. |
| `Get(k) (V, bool)` | Get; returns zero + `false` when absent. |
| `GetOrDefault(k, defaultV) V` | Get or fall back. |
| `Remove(k) (V, bool)` | Remove and return the previous value. |
| `Contains(k) bool` | Membership. |
| `Size() / IsEmpty() / Clear()` | Cardinality. |

### Views

| Method | Description |
|---|---|
| `Keys() *ArrayList[K]` | Keys (unordered). |
| `Values() *ArrayList[V]` | Values (unordered). |
| `Entries() *ArrayList[Entry[K, V]]` | `Entry` is `{K, V}`; `MarshalJSON` produces `{"k": v}`. |
| `ForEach(func(K, V))` | Walk without changing the map. |
| `Stream() stream.Stream[Entry[K, V]]` | Lazy stream. |
| `Collect() map[K]V` | Export as a Go built-in `map`. |
| `MarshalJSON / UnmarshalJSON` | `{"k": v, ...}` form. |

### Immutable-style transforms

`Filter(p) / FilterKeys(p) / FilterValues(p) *HashMap[K, V]` — return a new map.
`MapValues[R](f func(K, V) R) *HashMap[K, R]` — value-type transform.
`Concat(other) *HashMap[K, V]` — same keys are overwritten by `other`.

---

## `sets.HashSet[T]`

`T` must be `comparable`. `NewHashSet[T any]() *HashSet[T]` constructs.

### Basics

`Add(v) / Remove(v) / Contains(v) bool / Size() / IsEmpty() / Clear()` mirror the semantics of Go's `map[T]struct{}`; `Add` on an existing element is a no-op.

### Views and stream

`ForEach(func(T)) / Collect() []T / Stream() stream.Stream[T] / Peek(visit) *HashSet[T] / MarshalJSON / UnmarshalJSON` (JSON array, deduplicated on output).

### Immutable-style transforms

| Method | Description |
|---|---|
| `Filter(p) *HashSet[T]` | Keep elements satisfying `p`. |
| `Map[R](f) *ArrayList[R]` | Element-type transform (no longer deduplicated, typed as list). |
| `FlatMap[R](f) *ArrayList[R]` | Same but sub-results are also lists. |
| `MapSet[R comparable](f) *HashSet[R]` | Element-type transform with dedup. |
| `FlatMapSet[R comparable](f) *HashSet[R]` | Sub-results are sets; flatten and dedup. |
| `Concat(other) *HashSet[T]` | Union (deduplicated); receiver is not modified. |
| `Reduce[R](init R, f func(R, T) R) R` | Fold. |
| `SortBy(less) *ArrayList[T]` | Sort and export as a list. |
| `MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` | Empty set returns `option.Empty[T]()`. |

### Set algebra

`Concat / Union / Intersect / Difference / SymmetricDifference` all return a new `*HashSet[T]`, leaving the receiver unchanged:

| Method | Set notation |
|---|---|
| `Concat(other)` | `s ∪ other` |
| `Union(other)` | `s ∪ other` (same as `Concat`) |
| `Intersect(other)` | `s ∩ other` |
| `Difference(other)` | `s \ other` |
| `SymmetricDifference(other)` | `(s ∪ other) \ (s ∩ other)` |
| `IsSubsetOf(other) / IsSupersetOf(other)` | Containment. |

### Predicates

`Any / All / None / Find(p) option.Option[T]`.

---

## `stream.Stream[T]`

A lazy synchronous stream. `Stream[T]` is a thin wrapper over `iter.Seq[T]`; every transform returns a new `Stream` and does **not** consume the original. `Collect` / `Count` / `First` / `Last` / `Any` / `All` / `None` / `Find` / `Reduce` / `ForEach` / `MinBy` / `MaxBy` / `SortBy` are terminals.

### Construction

| Function | Description |
|---|---|
| `From[T](seq iter.Seq[T]) Stream[T]` | Direct wrap of `iter.Seq[T]`. |
| `FromSlice[T](s []T) Stream[T]` | Lazy iteration over `s` (the slice is held but not copied). |
| `Of[T](values ...T) Stream[T]` | Variadic → `Stream[T]`, sugar for `FromSlice`. |
| `Empty[T]() Stream[T]` | Terminates immediately. |

> Any collection (`ArrayList` / `LinkedList` / `HashSet` / `HashMap.Entries`) can be turned into a `Stream[T]` via its own `Stream()` method — that is the recommended entry point.

### Chained transforms

`Filter(p) / Map[R](f) / FlatMap[R](f) / Peek(visit) / Take(n) / Drop(n) / Distinct(eq) / Concat(other) / SortBy(less)` — all return `Stream[...]` and remain chainable.

`Associate[K, V](f func(T) (K, V)) MapStream[K, V]` lazily produces key-value pairs. `MapStream.Filter(func(K, V) bool)` and `MapValues[R](func(K, V) R)` preserve the map pipeline; `Collect()` returns a plain `map[K]V` (non-nil even for empty input). Duplicate keys remain through transformations; the last value reaching `Collect` wins. `Map[R](func(K, V) R)` explicitly projects to an ordinary `Stream[R]`, whose `Collect()` returns a slice.

```go
var result map[int]string = stream.Of("a", "bb", "ccc", "dd").
    Associate(func(s string) (int, string) { return len(s), s }).
    Filter(func(k int, _ string) bool { return k == 2 }).
    Collect() // map[int]string{2: "dd"}
```

### Terminals

| Method | Returns |
|---|---|
| `Collect() []T` | Materialise into a slice. |
| `Count() int` | Count (counts the post-filter elements). |
| `First() / Last() option.Option[T]` | First / last element. |
| `Any(p) / All(p) / None(p) bool` | Short-circuiting existential / universal. |
| `Find(p) option.Option[T]` | First matching element. |
| `Reduce(init T, f func(acc, v T) T) T` | Fold. |
| `ForEach(visit func(T))` | Pure consumption. |
| `MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` | Extreme; empty stream returns `option.Empty[T]()`. |
| `SortBy(less) Stream[T]` | Terminal: materialise and sort, then return a new stream. |

> `SortBy` on `Stream` is terminal (it walks the stream once) — unlike `ArrayList` / `LinkedList` / `HashSet`, where `SortBy` is a "transform-style" method.

### Example

```go
out := stream.Of(1, 2, 3, 4, 5).
    Filter(func(v int) bool { return v%2 == 1 }).
    Map(func(v int) int { return v * v }).
    Collect() // [1, 9, 25]

first, ok := stream.Of[int]().First().Get() // 0, false
```

---

## `collections.Range[T]`

`Range[T constraints.Integer](start, end T) stream.Stream[T]` — turns an integer half-open interval `[start, end)` into a `Stream[T]`.

- When `start >= end` the stream is empty; it does not panic.
- The stream is lazy: iterating `Range(0, 1_000_000_000)` costs constant memory, not a billion values.
- The type constraint is `constraints.Integer` (`int` / `uint` / `int32` / …); floats and strings are not accepted.

```go
collections.Range[int](1, 4).Collect() // [1, 2, 3]
collections.Range[int](0, 5).
    Map(func(i int) int { return i * i }).
    Take(3).
    Collect() // [0, 1, 4]
```

## See also

- Most return values are `option.Option[T]`; see [`adt`](../adt/README.md).
- The error flow uses `result.Result[T]`; see [`adt`](../adt/README.md).
- For asynchronous / multi-subscriber scenarios, use `reactivex.Observable[T]`; see [`reactivex`](../reactivex/README.md). The `Stream` here is a synchronous, single-consumer model — the two do not interoperate.
