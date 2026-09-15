<div align="center">

<img src="./docs/assets/typed-logo-gopher-official-light.png" alt="Typed — a fluent, type-safe Go generics toolkit" width="180" />

# Typed

**A fluent, type-safe Go generics toolkit.**

[中文版](./README-cn.md) · [Documentation](./docs/README.md)

[![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=root)](https://codecov.io/gh/qianwj/typed)
[![Module](https://img.shields.io/badge/module-github.com%2Fqianwj%2Ftyped-6f42c1)](#install)
[![Made with Go](https://img.shields.io/badge/made%20with-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)

</div>

Typed is a **fluent, type-safe Go generics toolkit**.
Concrete generic types (`ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`,
`HashSet[T]`, `Stream[T]`, `Subject[T]`) compose through left-to-right
chainable transforms (`Filter`, `Map[R]`, `FlatMap[R]`, `Reduce[R]`,
`Take`, `Drop`, `SortBy`, `Distinct`, `Concat`), backed by the supporting
abstractions (`Optional[T]`, `Result[T]`, `IsNil`, `Equals`) that the rest
of the toolkit is built on. The operator vocabulary is intentionally
familiar to anyone who has used Java Streams, .NET LINQ, or JavaScript
array pipelines — without inheriting their runtime model. On the async
side, `reactivex.Observable[T]` adds a typed event stream with explicit
demand, configurable backpressure, and a first-class multicast
`Subject[T]`.

## Why Typed

Go's `for` loop is clear and should remain the default for simple logic.
The trouble starts when a collection goes through several transforms:
nested functions, scattered `if err != nil` blocks, and `(T, bool)` returns
that do not chain.

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

Typed supports a left-to-right, chainable style:

```go
names := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()
```

For laziness, the same collection becomes a stream explicitly:

```go
profiles := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## Highlights

- 🧱 **Concrete generic types** — `ArrayList[T]`, `LinkedList[T]`,
  `HashMap[K, V]`, `HashSet[T]`, `Stack[T]`, `Queue[T]`, `Deque[T]`,
  `Stream[T]`, `Subject[T]`. No runtime type assertions, no `any`-shaped
  surprises.
- 🪄 **Fluent transforms** — `Filter`, `Map[R]`, `FlatMap[R]`, `Reduce[R]`,
  `Take`, `Drop`, `Distinct`, `SortBy`, `Concat`. All return concrete
  types under Go 1.27's generic methods.
- 🟢 **Optional / Result** — `Optional[T]` for "may be absent",
  `Result[T]` for "may fail". Both bridge cleanly to `(T, error)` and
  compose with the collection API.
- 📡 **Typed async streams** — `reactivex.Observable[T]` with explicit
  demand (`Request(n)`), per-subscription backpressure (`WithBuffer`,
  `OverflowStrategy`), and a hot multicast `Subject[T]`.
- 🧠 **Smart equality** — `objects.Equals[T]` understands
  `func (T) Equal(T) bool`, normalises nil-vs-empty slices / maps, and
  is nil-safe for typed-nil pointers.
- 🪶 **Bounded memory** — head-offset `ArrayList` with periodic
  compaction, slot zeroing on `Stack.Pop` and `Remove*`, so popped
  references do not leak through the backing array.
- 🧵 **Bounded & unbounded blocking queues** —
  `concurrency.BoundedBlockingQueue[T]` is a fixed-capacity FIFO with
  `Push` / `Poll` (blocking) and `TryPush` / `TryPoll` (non-blocking),
  backed by a single `chan T`. `concurrency.UnboundedBlockingQueue[T]`
  is the sibling that never blocks `Push`; ring buffer + mutex + cond
  under the hood. Both expose context-aware variants
  (`PushWithContext` / `PollWithContext`).
- 🧵 **Structured concurrency** — `concurrency.Group` is an
  `errgroup`-style helper with two failure policies: `Strict` (any
  task error fails the group, ctx cancels siblings) and `BestEffort`
  (tasks run to completion; only succeeds if all succeed, otherwise
  returns the aggregated errors). Built on `sync.WaitGroup` directly,
  no external dependencies.

## Install

```bash
go get github.com/qianwj/typed/collections
go get github.com/qianwj/typed/utils/option
go get github.com/qianwj/typed/utils/either
go get github.com/qianwj/typed/utils/result
go get github.com/qianwj/typed/reactivex
go get github.com/qianwj/typed/control
go get github.com/qianwj/typed/concurrency
go get github.com/qianwj/typed/utils/objects
go get github.com/qianwj/typed/utils/json
```

Each sub-package is its own Go module, so you can pull in only what you
use. Requires **Go 1.27+** for generic methods on concrete types.

## Quick start

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/utils/option"
    "github.com/qianwj/typed/utils/result"
)

// Eager transform
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// Optional-based access — no (T, bool) dance
first := adults.First().OrElse("(none)")

// Linear structures
s := collections.NewStack[int]()
s.Push(1); s.Push(2); s.Push(3)
top := s.Pop().OrElse(0) // 3

// Result: thread a (T, error) through combinators
v, err := result.Wrap(loadConfig(path)).
    MapError(func(err error) error { return fmt.Errorf("config %s: %w", path, err) }).
    Map(func(b []byte) Config { return parseConfig(b) }).
    Unwrap()
if err != nil { return err }
use(v)
```

## Table of contents

- [Design principles](#design-principles)
- [What ships today](#what-ships-today)
- [Documentation](#documentation)
- [Operator reference](#operator-reference)
- [Laziness and execution boundaries](#laziness-and-execution-boundaries)
- [Memory model](#memory-model)
- [Concurrency](#concurrency)
- [Error handling](#error-handling)
- [Go version](#go-version)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Bug reports](#bug-reports)
- [License](#license)

## Design principles

- **Type safety first.** Use Go generics and avoid `any`, reflection,
  and runtime type assertions wherever possible. `Optional[T]` carries
  an explicit `present` flag rather than relying on a nil check, so it
  works for any `T` including value types such as `int`, `string`, and
  `struct{}`.
- **Concrete types over interfaces.** `ArrayList[T]` and the rest are
  concrete generic types, not interfaces. Go 1.27's generic methods
  (`Map[R]`, `FlatMap[R]`, `Reduce[R]`) only work on concrete
  receivers; an interface would break fluent chaining.
- **Optional-based access.** Every accessor that can fail by absence
  returns `option.Optional[T]` rather than `(T, bool)`. The convention
  is uniform across `ArrayList.Get / First / Last / Find / MinBy / MaxBy / RemoveFirst / RemoveLast`, `LinkedList`, `Stack`, `Queue`, and `Deque`.
- **Bounded memory.** `ArrayList` uses a head-offset layout with
  periodic compaction so the retained capacity of a long-running
  head-drained list is bounded by the high-water mark of in-flight
  elements plus a constant, not by the all-time maximum the list ever
  saw. `Stack.Pop`, `ArrayList.RemoveFirst / RemoveLast`, and
  `Deque.PopFront / PopBack` zero the freed slot so the runtime's
  reachability walk does not keep popped references alive.
- **Optional laziness.** Collection operations are eager; `Stream[T]`
  is the explicit lazy layer, backed by Go's `iter.Seq[T]`.
- **Composability.** Collections, iterators, and `Stream` can be
  combined into new data sources; the snapshot contract on `Stream()`
  keeps mutations from leaking into in-flight pipelines.
- **Early termination.** `First`, `Any`, `All`, `Find`, `Take`, and
  the `Optional`-returning accessors stop as soon as the answer is
  known.
- **Opt-in, not a replacement.** Typed is a fluent layer on top of
  idiomatic Go; simple logic should remain easy to write with
  `for range` over a slice or map. Existing code that prefers plain
  slices, maps, and `chan T` is unaffected.

## What ships today

### Collection types

| Type | Kind | Source | Notes |
| --- | --- | --- | --- |
| `ArrayList[T]` | Concrete generic struct | `collections/lists` | Head-offset `[]T`; O(1) `Add` / `AddFirst` / `RemoveFirst` / `RemoveLast`; periodic compaction at head ≥ 64. |
| `LinkedList[T]` | Concrete generic struct | `collections/lists` | Doubly-linked; O(1) head / tail, O(i) random access. |
| `HashMap[K, V]` | Concrete generic struct | `collections/maps` | Open-addressed hash table; `K comparable`. |
| `HashSet[T]` | Concrete generic struct | `collections/sets` | `T comparable`; built on the same hash-table machinery. |
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
| `either.Either[L, R]` | `utils/either` | Tagged-union value type; `Left` = failure, `Right` = success by convention. Safe accessors return `Optional`. `Fold` for the canonical branch-pick pattern. |
| `result.Result[T]` | `utils/result` | Success / failure; `Unwrap` returns `(T, error)`, `Wrap` is the forward bridge from `(T, error)`, `Recover` does error-aware fallback that always returns a `T`. |
| `json.Encode[T] / Decode[T]` | `utils/json` | `encoding/json/v2`-backed `Result`-style codec. |
| `objects.Equaler` (interface, optional) | `utils/objects` | Hint interface `Equal(any) bool`; not required for `Equals` to dispatch. |
| `objects.IsNil[T]`, `objects.Equals[T]` | `utils/objects` | Reflection-based nil check and equality dispatch (handles typed nil, `func (T) Equal(T) bool`, `time.Time.Equal`). |

### Control flow

| Type / function | Source | Notes |
| --- | --- | --- |
| `control.Repeat(times, f)` / `RepeatE(times, f) (int, error)` | `control` | "Do this N times" loops; `RepeatE` stops at the first non-nil error and returns the number of successful iterations. |
| `match.Pattern[T]`, `match.Value[T]`, `match.Type(any)` | `control/match` | First-match-wins pattern matching: `Pattern[T]` for value tests, `Type` for dynamic-type dispatch. |

### Reactive streams

| Type | Source | Notes |
| --- | --- | --- |
| `reactivex.Observable[T]`, `Publisher[T]`, `Subscriber[T]` | `reactivex` | Typed async stream with explicit demand (`Subscription.Request(n)`) and `OnError` / `OnComplete` terminals. |
| `reactivex.Subject[T]` | `reactivex` | Hot multicast publisher and subscriber; configured with `WithBuffer` / `WithOverflow`. |
| `reactivex.Single[T]` | `reactivex` | Reactive container that emits exactly one value or one error. `Await` / `Subscribe` for blocking / callback consumption; `Map` / `FlatMap` / `Zip` / `AndThen` for composition. |
| `reactivex.Maybe[T]` | `reactivex` | Reactive container that emits zero-or-one value or one error — three terminal states (success / complete / error). |
| `reactivex.OverflowStrategy` | `reactivex` | `OverflowBlock` / `OverflowDropLatest` / `OverflowDropOldest` / `OverflowKeepLatest` / `OverflowError`. |
| Sources: `Just`, `FromSlice`, `FromChannel`, `FromChannelWithOptions`, `FromSeq`, `Create`, `Interval` | `reactivex` | Cold and hot source constructors. |
| Operators: `Map[R]`, `Filter`, `Take`, `Skip`, `Scan[R]`, `Reduce` | `reactivex` | All wrapping-style; no goroutines or queues of their own. |

### Concurrency primitives

| Type | Source | Notes |
| --- | --- | --- |
| `BoundedBlockingQueue[T]` | `concurrency` | Fixed-capacity FIFO (capacity rounded up to a power of two); `Push` / `Poll` block, `TryPush` / `TryPoll` do not; `TryPoll` returns `option.Optional[T]`. A thin generic wrapper around `chan T`. `0 B/op`, `0 allocs/op` on the hot path. |
| `UnboundedBlockingQueue[T]` | `concurrency` | Unbounded FIFO; `Push` never blocks, `Poll` blocks when empty. Ring buffer + `sync.Mutex` + `*sync.Cond` under the hood. Same `Push` / `Poll` / `WithContext` / `Try*` surface as `BoundedBlockingQueue`. |
| `Group` | `concurrency` | `errgroup`-style structured concurrency with two failure policies — `Strict` (any task error fails the group, ctx cancels siblings) and `BestEffort` (tasks run to completion; only succeeds if all succeed, otherwise returns `*BestEffortError`). Built on `sync.WaitGroup` directly, no external dependencies. |

## Documentation

Per-package API reference and examples, in English and Chinese:

- [docs/README.md](./docs/README.md) — index
- [collections](./docs/collections/README.md) — `ArrayList`, `LinkedList`, `HashMap`, `HashSet`, `Stack`, `Queue`, `Deque`, `Stream`, `Range`
- [control](./docs/control/README.md) — `Repeat` / `RepeatE` and `control/match`
- [reactivex](./docs/reactivex/README.md) — `Observable`, `Subject`, backpressure, operators
- [concurrency](./docs/concurrency/README.md) — `BoundedBlockingQueue[T]` (blocking + non-blocking, fixed capacity)
- [utils/option](./docs/option/README.md) — `Optional[T]`
- [utils/result](./docs/result/README.md) — `Result[T]`
- [utils/objects](./docs/utils/objects/README.md) — `IsNil`, `Equals`
- [utils/json](./docs/utils/json/README.md) — `Encode` / `Decode` on `encoding/json/v2`

## Operator reference

Operator names follow the de facto vocabulary shared by Java Streams,
.NET LINQ, and JavaScript array methods, so the API reads naturally
if you've used any of them. The table below maps the common aliases
to their Typed counterpart.

| Operation | Typed |
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

Typed is not a port of any one library. The fluent API is built on
top of Go's own generics, `iter.Seq[T]`, and explicit `error` returns;
the operator names are simply the de facto vocabulary engineers
already know.

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

Before `Collect`, `Filter`, `Map`, and `Take` only describe the
pipeline. The terminal operation starts consumption, and `Take(100)`
can stop the underlying source as soon as enough values have been
produced. `Stream()` creates a snapshot of the source at call time, so
later mutations to the source collection do not affect the in-flight
pipeline.

## Memory model

`ArrayList` stores the live range in `items[head:head+size]` and folds
the discarded prefix back to zero once `head >= 64`. The retained
capacity of a long-running head-drained list is therefore bounded by
the high-water mark of in-flight elements plus 64, not by the all-time
maximum the list ever saw. Reference elements popped from the front
are eligible for GC as soon as `RemoveFirst` zeros the freed slot. The
reference-handling tests in `collections/lists/arraylist_test.go` and
`collections/{stack,queue,deque}_test.go` use `runtime.SetFinalizer`
to assert that all popped boxes' finalizers run.

`Stack.Pop` is the slice-backed equivalent: it shrinks the backing
array and explicitly zeros the popped slot, so a `Stack[T]` of
pointers does not leak references through out-of-range slots.

## Concurrency

None of the collection types are safe for concurrent mutation. The
standard Go pattern — a single goroutine owns the collection,
communication happens over channels — applies unchanged. `Stream` is
not parallel by default; ordinary `Map` will not silently become
concurrent. The toolkit does not introduce a new concurrency model; it
follows Go's explicit one.

For the rare case where the toolkit's own surface area is a better fit
than `chan T` — non-blocking probes (`TryPush` / `TryPoll`), context-aware
blocking (`PushWithContext` / `PollWithContext`), an `Optional`-based
return for non-blocking reads, a single generic type to express "a
fixed-capacity blocking queue", or `errgroup`-style structured
concurrency with a typed `fn(ctx)` signature — the [`concurrency`
package](./docs/concurrency/README.md) ships
[`BoundedBlockingQueue[T]`](./docs/concurrency/README.md#boundedblockingqueuet)
(a thin wrapper around `chan T`, so per-op overhead is essentially
zero), [`UnboundedBlockingQueue[T]`](./docs/concurrency/README.md#unboundedblockingqueuet)
(ring buffer + mutex + cond, since the Go runtime has no "unbounded
buffered channel"), and [`Group`](./docs/concurrency/README.md#group)
(strict or best-effort structured concurrency on top of
`sync.WaitGroup`). Reach for `BoundedBlockingQueue` when you want
backpressure via capacity; reach for `UnboundedBlockingQueue` when
`Push` must never block; reach for `Group` when you have a fan-out
of goroutines and want the first error or the aggregated failures
without writing the boilerplate yourself.

## Error handling

`Result[T]` is the typed `(T, error)` carrier:

```go
v, err := result.Success(42).Unwrap()
if err != nil { return err }

port, _ := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err is always nil
```

`Recover` always returns a `T`; surface a fresh error by ending the
chain with `Unwrap` instead.

For collection pipelines that need to surface errors, the recommended
pattern is to convert the error path into an absent `Optional[T]`
(e.g. `OfNullable` on a lookup that returns `(T, error)`) and keep the
success path on the regular fluent API.

## Stability and SemVer

Every module is published under [Semantic Versioning 2.0.0](https://semver.org/).
The pre-`v1.0.0` numbering is chosen on purpose:

- **`0.0.x` (current).** No API stability guarantee. The toolkit is
  intentionally small enough that the cost of editing a method
  signature is low, and the audience is small enough that the cost of
  breaking an early adopter is also low. Patch releases (`0.0.x`)
  contain only bug fixes; anything that touches a public signature,
  exported type, or the documented behaviour of an existing
  operator is a **minor** bump (`0.0.x` → `0.(x+1).0`).
- **`0.y.0` (planned once a module reaches feature freeze).** Public
  types, exported function signatures, and the documented
  behaviour of an operator are frozen. Patch releases contain only
  bug fixes and documentation; no new API surface. A module is
  promoted to this tier when its API has been exercised for at
  least one minor cycle without changes.
- **`v1.0.0`.** Reserved for the modules the maintainer is willing
  to backport bug fixes to. Until then, "use at HEAD" is the
  recommended installation mode.

### What changes between minor releases will look like

- Renaming an exported type or method.
- Changing a generic receiver type (e.g. swapping `Map(func(T) R)`
  for `Map(func(context.Context, T) (R, error))`).
- Adding new required fields to a public struct.
- Changing a documented invariant (e.g. "`Map` returns an empty
  observable on a nil slice", "`MinBy` panics on an empty list").

### What does not count as a breaking change

- Adding a new method to an existing type.
- Adding a new top-level function or sub-package.
- Adding a new option to a variadic options function.
- Performance improvements that change big-O only on inputs the
  previous contract already declared out of scope (for example, the
  "no concurrency safety" note on every collection type).
- Bug fixes that change observable behaviour in ways that the
  documentation already said could happen.

### Per-module release tag format

Each module publishes its own tag prefix:

```text
utils/v0.0.1
collections/v0.0.1
control/v0.0.1
reactivex/v0.0.1
concurrency/v0.0.1
```

A change to one module's tag never implies a change to another.
`utils/v0.0.2` may ship the same week as `reactivex/v0.0.1`, or
years later; the prefixes are independent.

### How to find a breaking change

Every release notes a `BREAKING:` line for any signature or
behaviour change. The diff between the previous and current
`vX.Y.Z` tag is the authoritative changelog; release notes
summarise it but do not replace it.

## Go version

The current modules use **Go 1.27**:

```text
go 1.27.1
```

Go 1.23 introduced `iter.Seq`, `iter.Seq2`, and `for range` support
for function iterators. Go 1.27's generic methods make fluent APIs
such as `Stream[T].Map[R]`, `Optional[T].Map[R]`, and
`Result[T].Map[R]` possible.

## Roadmap

Done:

- [x] `ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]` with
      fluent methods.
- [x] Eager collection operations: `Filter`, `Map[R]`, `FlatMap[R]`,
      `Reduce[R]`, `Collect`, `ForEach`, `Peek`, `Concat`.
- [x] Order-dependent operations on lists: `Get`, `Insert`, `RemoveAt`,
      `First`, `Last`, `Find`, `Take`, `Drop`, `Distinct`, `SortBy`,
      `MinBy`, `MaxBy`.
- [x] Optional-based access: `Get` / `First` / `Last` / `Find` /
      `MinBy` / `MaxBy` / `RemoveFirst` / `RemoveLast` and `Stack` /
      `Queue` / `Deque` `Pop` / `Peek` / `Front` / `Back` return
      `Optional[T]`.
- [x] `Stream[T]` adapter backed by `iter.Seq[T]`, with
      early-terminating terminals.
- [x] Linear structures: `Stack[T]`, `Queue[T]`, `Deque[T]`.
- [x] Parent-package helpers: `Range[T constraints.Integer](start, end T)
      Stream[T]` for iota-style integer sequences.
- [x] `Option[T]` / `Result[T]` / `Equaler` / `IsNil[T]` / `Equals[T]`
      utilities.
- [x] `control.Repeat` / `RepeatE` and `control/match` pattern matching.
- [x] `reactivex` package: `Observable[T]` / `Publisher[T]` /
      `Subscriber[T]`, `Subject[T]`, `WithBuffer` / `WithOverflow`
      backpressure, `Map` / `Filter` / `Take` / `Skip` / `Scan` /
      `Reduce` operators.
- [x] `concurrency` package: `BoundedBlockingQueue[T]` (array-backed
      ring buffer, blocking + non-blocking variants, race-tested).
- [x] `utils/json` Result-style codec on top of `encoding/json/v2`.
- [x] Bounded memory: head-offset `ArrayList` with periodic compaction,
      slot-zeroing on `Stack.Pop` and the list `Remove*` paths.
- [x] Tests with the race detector on every module. Live coverage is
      reported per module via Codecov — see the badge above.

Open:

- [ ] Error-aware collection operations (`MapE`, `FilterE`, `CollectE`)
      as a first-class pipeline alternative to `Result[T]`-per-element.
- [ ] Benchmarks for `ArrayList` head-side ops, `LinkedList`
      iteration, `HashMap` resize behaviour, and `Stream` pipeline
      overhead.
- [ ] Iterators (`iter.Seq[T]`) as a first-class output of the
      collection types, parallel to `Stream()`.

## Contributing

The toolkit is intentionally small, but contributions are welcome. Before
opening a pull request, please skim the principles below so the review goes
faster for everyone.

- **Scope per module.** Each top-level package (`collections`, `control`,
  `reactivex`, `concurrency`, `utils/option`, `utils/result`,
  `utils/objects`, `utils/json`) ships as an **independent Go module** and
  has its own version tag prefix (see [Stability and SemVer](#stability-and-semver)).
  A PR that touches more than one module should call that out explicitly
  in the description and explain why cross-module coordination is needed;
  otherwise the maintainer will ask you to split it.
- **Discuss before you build.** For anything beyond a small bug fix or
  documentation typo, file an issue first so we can agree on the shape of
  the change. Things that look like "just adding a method" often turn into
  SemVer / API-surface decisions (see [the Stability chapter](#stability-and-semver)
  for what counts as breaking).
- **Tests are required.** A PR that changes behaviour must come with
  `go test ./...` passing **with `-race`** on the affected module. The
  toolkit's correctness story depends on it; PRs without a regression test
  will be asked to add one.
- **Linting must stay clean.** Run `golangci-lint run ./...` against the
  module you changed before pushing. The CI workflow runs the same
  configuration (`./.golangci.yml`) on every push and PR.
- **Public API is sacred.** If your change touches an exported type,
  method, function, or constant, the PR description must call it out
  explicitly and tag it as `BREAKING:`, `feat:`, or `fix:` so the release
  notes can be written correctly. See the [Per-module release tag
  format](#per-module-release-tag-format) section for the tagging
  conventions.
- **Commit hygiene.** One logical change per commit. Bug fixes and
  refactors should not be mixed with feature work in the same commit, and
  each commit message should make sense on its own.
- **Coding style.** Match the file you are editing. The codebase does not
  introduce new dependencies without discussion; do the same. When in
  doubt, read two neighbouring files first.

The standard flow:

1. Fork the repository and create a topic branch off `main`.
2. Make your change, run `go test -race ./...` and `golangci-lint run`
   on the affected module, push the branch.
3. Open a pull request against `qianwj/typed:main` with a clear
   description, a link to the issue it closes (if any), and the
   `BREAKING:` / `feat:` / `fix:` tag if applicable.

## Bug reports

Please use the [issue tracker](https://github.com/qianwj/typed/issues) on
GitHub. Before opening a new issue, search to make sure it is not already
filed (including closed ones — the fix may have landed on a different
module's branch and not yet been tagged).

A good bug report includes:

- **Goal.** One sentence: what were you trying to do with Typed?
- **Module and version.** Which module is affected (`collections`,
  `reactivex`, `concurrency`, ...), and at which commit / tag? `git rev-parse HEAD`
  inside the module directory, or the module's `vX.Y.Z` tag, is enough.
- **Environment.** `go version`, `go env GOOS GOARCH`, and (if relevant)
  the OS / kernel version.
- **Reproduction.** The smallest self-contained snippet that triggers the
  bug. The preferred form is a Go test or a `go run`-able file; the
  `gist`-style "here are ten lines of code that compile" report is
  usually not enough to reproduce.
- **Expected vs actual.** What you expected to happen, what actually
  happened, and the exact error / panic / log output. For panics,
  include the full stack trace.
- **Race / data-race reports.** If you suspect a data race, say so
  explicitly and include the `-race` output verbatim. Most concurrency
  bugs in this toolkit only surface with `go test -race` or the runtime
  race detector attached to a long-running process.
- **Workarounds.** Anything you tried that got you past the bug, even if
  it is ugly — it often reveals the underlying assumption.

Security issues should **not** go through the public tracker. Use GitHub's
[private security reporting](https://github.com/qianwj/typed/security/advisories/new)
flow so the maintainer can coordinate a fix before disclosure.

## License

This project is licensed under the [MIT License](./LICENSE).
