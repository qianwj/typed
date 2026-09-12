# typed

> [中文版 README](./README-cn.md)

`typed` is a type-safe collection toolkit built with Go generics. It provides concrete, eager collection types (`ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]`), a small lazy `Stream[T]` layer, dedicated linear data structures (`Stack[T]`, `Queue[T]`, `Deque[T]`), and the supporting abstractions (`Option[T]`, `Result[T]`, `Equaler[T]`) that the collections are built on.

The collections are Java / JavaScript-style: concrete generic types with fluent methods such as `Filter`, `Map[R]`, `FlatMap[R]`, and `Reduce[R]`. Optional accessors (`Get`, `First`, `Last`, `Find`, `MinBy`, `MaxBy`, `Pop`, `Peek`, `Front`, `Back`) return `Option[T]` rather than `(T, bool)` so callers can chain `OrElse` / `OrElseGet` / `Map` on the result.

## Status

The core types and the `Option[T]` / `Result[T]` abstractions are in place and tested. Roadmap items below describe work that is not yet started.

## Why

Go's `for` loop is clear and should remain the preferred choice for simple logic. But when a collection goes through several transformations, nested functions or repeated temporary slices can make the processing flow harder to read:

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

`typed` supports a left-to-right collection style:

```go
result := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()
```

For lazy processing, the same collection can become a stream explicitly:

```go
result := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## Design Principles

- **Type safety first.** Use Go generics and avoid `any`, reflection, and runtime type assertions wherever possible. `Option[T]` carries an explicit `present` flag rather than relying on a nil check, so it works for any `T` including value types such as `int`, `string`, and `struct{}`.
- **Concrete types over interfaces.** `ArrayList[T]` and the rest are concrete generic types, not interfaces. Go 1.27's generic methods (`Map[R]`, `FlatMap[R]`, `Reduce[R]`) only work on concrete receivers; an interface would have to break fluent chaining by returning the interface type from methods that today return the concrete type.
- **Optional-based access.** Every accessor that can fail by absence returns `option.Optional[T]` rather than `(T, bool)`. The convention is uniform across `ArrayList.Get / First / Last / Find / MinBy / MaxBy / RemoveFirst / RemoveLast`, `LinkedList`, `Stack`, `Queue`, and `Deque`.
- **Bounded memory.** `ArrayList` uses a head-offset layout with periodic compaction so the retained capacity of a long-running head-drained list is bounded by the high-water mark of in-flight elements plus a constant, not by the all-time maximum the list ever saw. `Stack.Pop`, `ArrayList.RemoveFirst / RemoveLast`, and `Deque.PopFront / PopBack` zero the freed slot so the runtime's reachability walk does not keep popped references alive.
- **Optional laziness.** Collection operations are eager; `Stream[T]` is the explicit lazy layer, backed by Go's `iter.Seq[T]`.
- **Composability.** Collections, iterators, and `Stream` can be combined into new data sources; the snapshot contract on `Stream()` keeps mutations from leaking into in-flight pipelines.
- **Early termination.** `First`, `Any`, `All`, `Find`, `Take`, and the `Optional`-returning accessors stop as soon as the answer is known.
- **Go style.** Do not copy every Java Stream semantic blindly; simple logic should remain easy to write with `for range`. The toolkit is opt-in: existing code that prefers slices and maps is unaffected.

## What ships today

### Collection types

| Type | Kind | Source | Notes |
| --- | --- | --- | --- |
| `ArrayList[T]` | Concrete generic struct | `collections/lists` | Head-offset `[]T`; O(1) `Add` / `AddFirst` / `RemoveFirst` / `RemoveLast`; periodic compaction at head ≥ 64. |
| `LinkedList[T]` | Concrete generic struct | `collections/lists` | Doubly-linked; O(1) head / tail, O(i) random access. |
| `HashMap[K, V]` | Concrete generic struct | `collections/maps` | Open-addressed hash table; `K comparable`. |
| `HashSet[T]` | Concrete generic struct | `collections/sets` | `T comparable`; built on top of the same hash-table machinery. |
| `Stack[T]` | Concrete generic struct | `collections` | Single-ended LIFO; `[]T` with explicit slot zeroing on `Pop`. |
| `Queue[T]` | Concrete generic struct | `collections` | Single-ended FIFO; thin wrapper over `ArrayList[T]`. |
| `Deque[T]` | Concrete generic struct | `collections` | Double-ended; thin wrapper over `LinkedList[T]`. All operations are strict O(1). |
| `Stream[T]` | Concrete generic struct | `collections/stream` | Lazy, single-use pipeline over `iter.Seq[T]`. |

### Parent-package constructors

| Function | Source | Notes |
| --- | --- | --- |
| `Range[T constraints.Integer](start, end T) Stream[T]` | `collections` | Lazy iota-style `[start, end)` stream. |

### Abstraction types

| Type | Source | Notes |
| --- | --- | --- |
| `option.Optional[T]` | `utils/option` | Present / absent value, no nil check on `T`. |
| `result.Result[T]` | `utils/result` | Success / failure; `Unwrap` returns `(T, error)`, `Wrap` is the forward bridge from `(T, error)`, `Recover` does error-aware fallback that always returns a `T`. |
| `json.Encode[T] / Decode[T]` | `utils/json` | `encoding/json/v2`-backed `Result`-style codec. |
| `objects.Equaler` (interface, optional) | `utils/objects` | Hint interface `Equal(any) bool`; not required for `Equals` to dispatch. |
| `objects.IsNil[T]`, `objects.Equals[T]` | `utils/objects` | Reflection-based nil check and equality dispatch (handles typed nil, `func (T) Equal(T) bool`, `time.Time.Equal`). |

### Control

| Type / function | Source | Notes |
| --- | --- | --- |
| `control.Repeat(times, f)` / `RepeatE(times, f) (int, error)` | `control` | "Do this N times" loops; `RepeatE` stops at the first non-nil error and returns the number of successful iterations. |
| `match.Pattern[T]`, `match.Value[T]`, `match.Type(any)` | `control/match` | First-match-wins pattern matching: `Pattern[T]` for value tests, `Type` for dynamic-type dispatch. |

## Quick tour

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/utils/option"
)

// Eager transformation.
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// Optional-based access: no (T, bool) dance.
firstAdult := adults.First().OrElse("(none)")

// Linear structures.
s := collections.NewStack[int]()
s.Push(1); s.Push(2); s.Push(3)
top := s.Pop().OrElse(0)  // 3

q := collections.NewQueue[int]()
q.Push(1); q.Push(2)
front := q.Pop().OrElse(0)  // 1

d := collections.NewDeque[int]()
d.PushBack(1); d.PushFront(0); d.PushBack(2)
// [0, 1, 2]
left  := d.PopFront().OrElse(-1)  // 0
right := d.PopBack().OrElse(-1)   // 2

// Optional chaining.
v := option.Of(7).Map(func(x int) int { return x * 2 }).OrElse(0)  // 14
```

## Fluent collection API

`ArrayList[T]`, `LinkedList[T]`, and `HashSet[T]` all provide a common core of operations:

| Capability | API |
| --- | --- |
| Construction | `NewX[T]()`, `XOf(values...)` — both return `*X[T]` |
| Cardinality | `Size()` |
| Lifecycle | `IsEmpty()`, `Clear()` |
| Traversal | `ForEach`, `Collect`, `Stream` |
| Conditional | `Any`, `All`, `None`, `Find` (returns `Optional[T]`) |
| Selection | `Filter` — returns the same kind |
| Type-changing | `Map[R](func(T) R)` — returns `*ArrayList[R]` |
| Flattening | `FlatMap[R](func(T) *ArrayList[R])` — returns `*ArrayList[R]` |
| Folding | `Reduce[R](init, func(R, T) R)` |
| Observation | `Peek`, `Concat` — returns the same kind |

`ArrayList[T]` and `LinkedList[T]` additionally provide order-dependent operations (`Get`, `Insert`, `RemoveAt`, `First`, `Last`, `Take`, `Drop`, `Distinct`, `SortBy`, `MinBy`, `MaxBy`). `HashMap[K, V]` provides `Get`, `GetOrDefault`, `Put`, `Remove`, `Keys`, `Values`, `Entries`, `Filter`, `MapValues[R]`, and `Concat`. `HashSet[T]` provides `Union`, `Intersect`, `Difference`, `SymmetricDifference`, `IsSubsetOf`, `IsSupersetOf`, plus `MapSet` and `FlatMapSet` that preserve deduplication when the result is `comparable`.

`Map` and `FlatMap` always return `*ArrayList[R]` because the result type `R` can be a slice, map, function, or any other non-comparable type. `HashSet.MapSet` and `HashSet.FlatMapSet` are the deduplicating variants.

## Optional access

`Optional[T]` is the standard return shape for any accessor that can fail by absence:

| Operation | Empty result |
| --- | --- |
| `Optional.Get()` (when absent) | panic |
| `Optional.OrElse(default)` | `default` |
| `Optional.OrElseGet(f)` | `f()` |
| `Optional.OrElseThrow(msg)` | `(zero, errors.New(msg))` |
| `Optional.IsPresent` / `IsEmpty` | bool |
| `Optional.Map[R](f)` | absent `Optional[R]`; `f` is not called |
| `Optional.FlatMap[R](f)` | absent `Optional[R]`; `f` is not called |
| `Optional.Filter(predicate)` | absent if predicate is false |
| `Optional.IfPresent(f)` / `IfPresentOrElse(p, a)` | no-op or `a()` |

The present / absent flag is stored explicitly, not inferred from a nil check on `T`, so `Optional[int]` works for value types where `(int, bool)` would be the only alternative.

`Result[T]` is the `(T, error)`-shaped companion. `Unwrap` returns the success value and a `nil` error on the happy path, or the zero value and the captured error on failure. `Recover(f func(error) T) T` takes a fallback function that receives the captured error and produces a replacement value; the chain always returns a `T`, so to surface a new error use `Unwrap`. `Wrap` is the forward bridge from a `(T, error)` return into `Result[T]`.

```go
v, err := loadProfile(id).Unwrap()
if err != nil { return err }

port, _ := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err is always nil
```

`Result.MapError` transforms the captured error in place; `Map[R]` / `FlatMap[R]` thread the success value through a transformation while preserving the failure on the error path.

## Memory model

`ArrayList` stores the live range in `items[head:head+size]` and folds the discarded prefix back to zero once `head >= 64`. The retained capacity of a long-running head-drained list is therefore bounded by the high-water mark of in-flight elements plus 64, not by the all-time maximum the list ever saw. Reference elements popped from the front are eligible for GC as soon as `RemoveFirst` zeros the freed slot. The reference-handling tests in `collections/lists/arraylist_test.go` and `collections/{stack,queue,deque}_test.go` use `runtime.SetFinalizer` to assert that all popped boxes' finalizers run.

`Stack.Pop` is the slice-backed equivalent: it shrinks the backing array and explicitly zeros the popped slot, so a `Stack[T]` of pointers does not leak references through out-of-range slots.

## Concurrency

None of the collection types are safe for concurrent mutation. The standard Go pattern — a single goroutine owns the collection, communication happens over channels — applies unchanged. `Stream` is not parallel by default; ordinary `Map` will not silently become concurrent. The toolkit does not introduce a new concurrency model; it follows Go's explicit one.

## Project structure

```text
typed/
├── collections/
│   ├── go.mod
│   ├── mod.go
│   ├── stack.go, stack_test.go
│   ├── queue.go, queue_test.go
│   ├── deque.go, deque_test.go
│   ├── lists/
│   │   ├── arraylist.go, arraylist_test.go
│   │   ├── linkedlist.go
│   │   └── lists_test.go
│   ├── maps/
│   │   └── hashmap.go, hashmap_test.go
│   ├── sets/
│   │   └── hashset.go, hashset_test.go
│   └── stream/
│       └── mod.go, stream_test.go
├── utils/
│   ├── go.mod
│   ├── option/
│   ├── result/
│   ├── objects/
│   └── json/
├── control/
│   ├── go.mod
│   ├── mod.go             # Repeat / RepeatE
│   └── match/             # Pattern matching
├── reactivex/
│   ├── go.mod
│   ├── mod.go
│   ├── interfaces.go      # Publisher / Subscriber / Subscription / Observable
│   ├── sources.go         # Just / FromSlice / FromChannel / FromSeq / Create / Interval
│   ├── operators.go       # Map / Filter / Take / Skip / Scan / Reduce
│   ├── subject.go         # NewSubject + Subscriber-side bridge
│   ├── options.go         # WithBuffer / WithOverflow / OverflowStrategy
│   ├── backpressure.go
│   └── collect.go         # Subscribe / ForEach / ToSlice
└── docs/                     # 英文 README.md + 中文 README-cn.md 每篇都存在
    ├── README.md            # 英文索引
    ├── README-cn.md         # 中文索引
    ├── collections/{README.md, README-cn.md}
    ├── option/{README.md, README-cn.md}
    ├── result/{README.md, README-cn.md}
    ├── control/{README.md, README-cn.md}
    ├── reactivex/{README.md, README-cn.md}
    └── utils/
        ├── objects/{README.md, README-cn.md}
        └── json/{README.md, README-cn.md}
```

## Java / JavaScript mapping

| Java / JavaScript | `typed` |
| --- | --- |
| `stream()` | `ArrayListOf(...).Stream()` |
| `filter` / `where` | `Filter` |
| `map` / `select` | `Map` |
| `flatMap` / `selectMany` | `FlatMap` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` / `take` | `Take` |
| `skip` / `drop` | `Drop` |
| `findFirst` / `find` | `First` / `Find` (returns `Optional[T]`) |
| `anyMatch` / `some` | `Any` |
| `allMatch` / `every` | `All` |
| `noneMatch` | `None` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |
| `Optional.of` / `ofNullable` | `option.Of` / `option.OfNullable` |
| `Optional.orElse` | `OrElse` |
| `Stream` lazy | `Stream[T]` |
| `IntStream.range` | `Range(start, end) Stream[T]` |
| `Deque` (Java) | `Deque[T]` (LinkedList-backed) |

The project does not attempt to copy the Java or JavaScript runtime model. It borrows their collection-processing style while preserving Go's static typing, explicit errors, and straightforward control flow.

## Laziness and execution boundaries

A typical `Stream` pipeline looks like:

```text
source → intermediate operation → intermediate operation → terminal operation
```

For example:

```go
adults := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

Before `Collect`, `Filter`, `Map`, and `Take` only describe the pipeline. The terminal operation starts consumption, and `Take(100)` can stop the underlying source as soon as enough values have been produced. `Stream()` creates a snapshot of the source at call time, so later mutations to the source collection do not affect the in-flight pipeline.

## Error handling

`Result[T]` is the typed `(T, error)` carrier:

```go
v, err := result.Success(42).Unwrap()
if err != nil { return err }

port, _ := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err is always nil
```

`Recover` always returns a `T`; surface a fresh error by ending the chain with `Unwrap` instead.

`Result.MapError` transforms the captured error in place; `Map[R]` / `FlatMap[R]` thread the success value through a transformation while preserving the failure on the error path.

For collection pipelines that need to surface errors, the recommended pattern is to convert the error path into an absent `Optional[T]` (e.g. `OfNullable` on a lookup that returns `(T, error)`) and keep the success path on the regular fluent API.

## Go version

The current modules use Go 1.27:

```text
go 1.27.1
```

Go 1.23 introduced `iter.Seq`, `iter.Seq2`, and `for range` support for function iterators. Go 1.27's generic methods make fluent APIs such as `Stream[T].Map[R]`, `Optional[T].Map[R]`, and `Result[T].Map[R]` possible.

## Roadmap

Done:

- [x] `ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]` with fluent methods.
- [x] Eager collection operations: `Filter`, `Map[R]`, `FlatMap[R]`, `Reduce[R]`, `Collect`, `ForEach`, `Peek`, `Concat`.
- [x] Order-dependent operations on lists: `Get`, `Insert`, `RemoveAt`, `First`, `Last`, `Find`, `Take`, `Drop`, `Distinct`, `SortBy`, `MinBy`, `MaxBy`.
- [x] Optional-based access: `Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy` / `RemoveFirst` / `RemoveLast` and `Stack` / `Queue` / `Deque` `Pop` / `Peek` / `Front` / `Back` return `Optional[T]`.
- [x] `Stream[T]` adapter backed by `iter.Seq[T]`, with early-terminating terminals.
- [x] Linear structures: `Stack[T]`, `Queue[T]`, `Deque[T]`.
- [x] Parent-package helpers: `Range[T constraints.Integer](start, end T) Stream[T]` for iota-style integer sequences.
- [x] `Option[T]` / `Result[T]` / `Equaler` / `IsNil[T]` / `Equals[T]` utilities.
- [x] `control.Repeat` / `RepeatE` and `control/match` pattern matching.
- [x] `reactivex` package: `Observable[T]` / `Publisher[T]` / `Subscriber[T]`, `Subject[T]`, `WithBuffer` / `WithOverflow` backpressure, `Map` / `Filter` / `Take` / `Skip` / `Scan` / `Reduce` operators.
- [x] `utils/json` Result-style codec on top of `encoding/json/v2`.
- [x] Bounded memory: head-offset `ArrayList` with periodic compaction, slot-zeroing on `Stack.Pop` and the list `Remove*` paths.
- [x] Tests with race detector, 100% statement coverage on the actively-developed files.

Open:

- [ ] Error-aware collection operations (`MapE`, `FilterE`, `CollectE`) as a first-class pipeline alternative to `Result[T]`-per-element.
- [ ] Benchmarks for `ArrayList` head-side ops, `LinkedList` iteration, `HashMap` resize behaviour, and `Stream` pipeline overhead.
- [ ] Iterators (`iter.Seq[T]`) as a first-class output of the collection types, parallel to `Stream()`.

## License

This project is licensed under the [MIT License](LICENSE).
