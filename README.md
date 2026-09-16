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

A **fluent, type-safe Go generics toolkit**. Compose concrete generic
types (`ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]`,
`Stream[T]`) through left-to-right chainable transforms (`Filter`,
`Map[R]`, `FlatMap[R]`, `Reduce[R]`, `Take`, `Drop`, `Distinct`,
`SortBy`, `Concat`), backed by the supporting abstractions
(`Option[T]`, `Result[T]`, `Either[L, R]`, `IsNil`, `Equals`) that the
rest of the toolkit is built on. The operator vocabulary follows the
de facto naming used by Java Streams, .NET LINQ, and JavaScript array
pipelines — without inheriting their runtime model. On the async side,
`reactivex.Observable[T]` adds a typed event stream with explicit
demand, configurable backpressure, and a first-class multicast
`Subject[T]`.

## Why Typed

Go's `for` loop is clear and should stay the default for simple logic.
The trouble starts when a collection goes through several transforms:
nested functions, scattered `if err != nil` blocks, `(T, bool)` returns
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

For laziness, the same source becomes a stream explicitly:

```go
profiles := stream.Of(users...).
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## Modules

| Module | Path | What's in it |
|---|---|---|
| `collections` | `github.com/qianwj/typed/collections` | `ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]`, `Stack[T]`, `Stream[T]`, `Range[T]` |
| `collections/lists` | (sub-module) | `ArrayList[T]`, `LinkedList[T]` with fluent transforms |
| `collections/queues` | (sub-module) | `Queue[T]`, `Deque[T]`, `PriorityQueue[T]` |
| `collections/stream` | (sub-module) | `Stream[T]` — the fluent lazy layer over `iter.Seq[T]` |
| `adt` | `github.com/qianwj/typed/adt` | `Option[T]`, `Result[T]`, `Either[L, R]` |
| `reactivex` | `github.com/qianwj/typed/reactivex` | `Observable[T]`, `Subject[T]`, `Single[T]`, `Maybe[T]`, backpressure |
| `concurrency` | `github.com/qianwj/typed/concurrency` | `BoundedBlockingQueue[T]`, `UnboundedBlockingQueue[T]`, `Group`, `Semaphore`, `Pool[T]` |
| `control` | `github.com/qianwj/typed/control` | `If`/`IfGet`, `Repeat`/`RepeatE`, pattern matching |
| `utils/objects` | `github.com/qianwj/typed/utils/objects` | `IsNil[T]`, `Equals[T]` |
| `utils/json` | `github.com/qianwj/typed/utils/json` | `Encode[T]`, `Decode[T]` (Result-style on `encoding/json/v2`) |

Each top-level package and sub-package is its own Go module — pull in
only what you use.

## Highlights

- **Concrete types, not interfaces.** `ArrayList[T]`, `LinkedList[T]`, etc. are concrete generic structs. Go 1.27's generic methods (`Map[R]`, `FlatMap[R]`, `Reduce[R]`) need concrete receivers.
- **Fluent transforms return concrete types.** `Filter`, `Map[R]`, `FlatMap[R]`, `Take`, `Drop`, `Distinct`, `SortBy`, `Concat` chain on the same type — no boxing into `any`.
- **`Option[T]` / `Result[T]` / `Either[L, R]`.** Absence and failure are first-class types. Both `Option` and `Result` bridge to `(T, error)` cleanly.
- **No reflection in the hot path.** Only `utils/objects` uses reflection (for `IsNil`/`Equals`), and only there.
- **Bounded memory.** `ArrayList` uses a head-offset layout with periodic compaction at `head >= 64`. `Stack.Pop` and `Remove*` zero the freed slot so popped references are GC-eligible.
- **Typed async streams.** `reactivex.Observable[T]` with explicit demand (`Request(n)`) and per-subscription backpressure (`OverflowStrategy`).
- **Synchronous concurrency primitives.** `BoundedBlockingQueue[T]`, `UnboundedBlockingQueue[T]`, `Group` (Strict / BestEffort), `Semaphore`, `Pool[T]`. Backed by `chan T` / `sync.WaitGroup` / `sync.Cond` directly — no third-party deps.

## Install

```bash
go get github.com/qianwj/typed/collections
go get github.com/qianwj/typed/adt
go get github.com/qianwj/typed/reactivex
go get github.com/qianwj/typed/control
go get github.com/qianwj/typed/concurrency
go get github.com/qianwj/typed/utils/objects
go get github.com/qianwj/typed/utils/json
```

Requires **Go 1.27+** for generic methods on concrete types.

## Quick start

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/collections/queues"
    "github.com/qianwj/typed/collections/stream"
    "github.com/qianwj/typed/adt"
)

// Eager fluent transform
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

first := adults.First().OrElse("(none)")

// Lazy pipeline over an iter.Seq
profiles := stream.Of(users...).
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()

// Option-based access — no (T, bool) dance
top := collections.NewStack[int]().Push(1).Push(2).Push(3).Pop().OrElse(0)

// Result: thread a (T, error) through combinators
v, err := adt.Wrap(loadConfig(path)).
    MapError(func(err error) error { return fmt.Errorf("config %s: %w", path, err) }).
    Map(func(b []byte) Config { return parseConfig(b) }).
    Unwrap()
if err != nil { return err }
use(v)

// Top-K with a comparator
top3 := queues.NewPriorityQueue(3, cmp.Less[int])
top3.Push(42); top3.Push(7); top3.Push(99); top3.Push(1)
for top3.Size() > 0 {
    fmt.Println(top3.Pop().Get()) // 1, 7, 42
}
```

## Design principles

- **Type safety first.** Generics, no `any`, no reflection in the hot path. `Option[T]` carries an explicit `present` flag rather than a nil check, so it works for value types like `int`, `string`, and `struct{}`.
- **Concrete types, not interfaces.** Generic methods (`Map[R]`, `FlatMap[R]`, `Reduce[R]`) require concrete receivers. `ArrayList[T]` is a struct, not an `interface{}`.
- **`Option`-based access.** Every accessor that can fail by absence returns `Option[T]`, not `(T, bool)`. Same convention across `ArrayList`, `LinkedList`, `Stack`, `Queue`, `Deque`, `PriorityQueue`.
- **Bounded memory.** `ArrayList` head-offset layout with compaction at `head >= 64` bounds retained capacity by the high-water mark of in-flight elements plus a constant, not by all-time maximum. `Pop` / `Remove*` zero freed slots so popped references are GC-eligible.
- **Eager collections, lazy streams.** `ArrayList.Filter` runs immediately. `ArrayList.Stream()` (and `stream.Of(iter.Seq)`) gives a `Stream[T]`-backed lazy layer.
- **Composable but opt-in.** Simple logic stays as `for range`. Typed is a fluent layer, not a replacement for idiomatic Go.

## What ships

### Collection types

| Type | Source | Notes |
|---|---|---|
| `ArrayList[T]` | `collections/lists` | Head-offset `[]T`; O(1) `Add` / `AddFirst` / `RemoveFirst` / `RemoveLast`; periodic compaction at `head >= 64`. |
| `LinkedList[T]` | `collections/lists` | Doubly-linked; O(1) head/tail. |
| `HashMap[K, V]` | `collections/maps` | Open-addressed hash table; `K comparable`. |
| `HashSet[T]` | `collections/sets` | `T comparable`; same hash-table machinery. |
| `Stack[T]` | `collections` | LIFO; `[]T` with explicit slot zeroing on `Pop`. |
| `Queue[T]` | `collections/queues` | FIFO; thin wrapper over `ArrayList[T]`. |
| `Deque[T]` | `collections/queues` | Double-ended; thin wrapper over `LinkedList[T]`; strict O(1) on both ends. |
| `PriorityQueue[T]` | `collections/queues` | Heap-backed; comparator-ordered; bounded or unbounded top-K. |
| `Stream[T]` | `collections/stream` | Lazy single-use pipeline over `iter.Seq[T]`. |
| `Range[T]` | `collections` | Integer half-open interval as a `Stream[T]` factory. |

### Abstraction types

| Type | Source | Notes |
|---|---|---|
| `adt.Option[T]` | `adt` | Present / absent value; explicit `present` flag, no nil check on `T`. |
| `adt.Result[T]` | `adt` | Success / failure. `Wrap` / `Unwrap` bridge to `(T, error)`. `Recover` for error-aware fallback. |
| `adt.Either[L, R]` | `adt` | Tagged-union; safe accessors return `Option`. `Fold` for branch-pick. |
| `json.Encode[T]` / `Decode[T]` | `utils/json` | `encoding/json/v2`-backed `Result`-style codec. |
| `objects.IsNil[T]`, `objects.Equals[T]` | `utils/objects` | Reflection-based nil check and equality dispatch. |

### Reactive streams

| Type / function | Source | Notes |
|---|---|---|
| `Observable[T]`, `Publisher[T]`, `Subscriber[T]` | `reactivex` | Typed async stream with explicit demand (`Request(n)`) and `OnError` / `OnComplete` terminals. |
| `Subject[T]` | `reactivex` | Hot multicast publisher; `WithBuffer` / `WithOverflow` for backpressure. |
| `Single[T]`, `Maybe[T]` | `reactivex` | Reactive containers for "exactly one" and "zero or one" emissions. Await returns `Result[T]` / `Result[Option[T]]`, respectively. |
| `OverflowStrategy` | `reactivex` | `OverflowBlock` / `OverflowDropLatest` / `OverflowDropOldest` / `OverflowKeepLatest` / `OverflowError`. |

### Control flow

| Type / function | Source | Notes |
|---|---|---|
| `control.If[T]`, `IfGet[T]` | `control` | `If` selects eager values; `IfGet` evaluates only the chosen callback. |
| `control.Repeat`, `RepeatE` | `control` | "Do N times"; `RepeatE` stops at the first non-nil error. |
| `match.Pattern[T]`, `match.Value[T]`, `match.Type` | `control/match` | First-match-wins pattern matching. |

### Concurrency primitives

| Type | Source | Notes |
|---|---|---|
| `BoundedBlockingQueue[T]` | `concurrency` | Fixed-capacity FIFO (`chan T` wrapper); `Push` / `Poll` block, `TryPush` / `TryPoll` non-blocking. |
| `UnboundedBlockingQueue[T]` | `concurrency` | Ring buffer + `sync.Mutex` + `*sync.Cond`; `Push` never blocks. |
| `Group` | `concurrency` | `errgroup`-style structured concurrency; `Strict` / `BestEffort` modes. |
| `Semaphore` | `concurrency` | Counting semaphore (fixed unit weights); `Acquire` / `Release` / `TryAcquire` / `AcquireWithContext`. |
| `Pool[T]` | `concurrency` | Typed `sync.Pool` with `NewPool(creator)`, `Get() Option[T]`, `Put(T)`. |

## Operator reference

Operator names follow the de facto vocabulary shared by Java Streams,
.NET LINQ, and JavaScript array methods, so the API reads naturally if
you've used any of them.

| Operation | Typed |
|---|---|
| `stream()` | `ArrayListOf(...).Stream()` / `stream.Of(iter.Seq[T])` |
| `filter` / `where` | `Filter` |
| `map` / `select` | `Map[R]` |
| `flatMap` / `selectMany` | `FlatMap[R]` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` / `take` | `Take` |
| `skip` / `drop` | `Drop` |
| `findFirst` / `find` | `First` / `Find` (returns `Option[T]`) |
| `anyMatch` / `some` | `Any` |
| `allMatch` / `every` | `All` |
| `noneMatch` | `None` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |
| `IntStream.range` | `Range(start, end)` |

Typed is not a port of any one library. The fluent API is built on
top of Go's own generics, `iter.Seq[T]`, and explicit `error` returns;
the operator names are simply the de facto vocabulary engineers
already know.

## Memory model

`ArrayList` stores the live range in `items[head:head+size]` and folds
the discarded prefix back to zero once `head >= 64`. The retained
capacity of a long-running head-drained list is therefore bounded by
the high-water mark of in-flight elements plus 64, not by the all-time
maximum the list ever saw. Reference elements popped from the front are
eligible for GC as soon as `RemoveFirst` zeros the freed slot. The
reference-handling tests in `collections/lists/arraylist_test.go`,
`collections/stack_test.go`, and `collections/queues/queue_test.go` use
`runtime.SetFinalizer` to assert that all popped boxes' finalizers run.

`Stack.Pop`, `Queue.Pop`, `Deque.PopFront` / `PopBack`, and the
`ArrayList.Remove*` paths all zero the freed slot — a
`*collections.Stack[int]` of pointers does not leak references
through out-of-range slots.

## Concurrency

None of the collection types are safe for concurrent mutation. The
standard Go pattern — a single goroutine owns the collection,
communication happens over channels — applies unchanged. `Stream` is
not parallel by default; ordinary `Map` will not silently become
concurrent. The toolkit does not introduce a new concurrency model;
it follows Go's explicit one.

For the cases where the toolkit's own surface area is a better fit than
`chan T` — non-blocking probes (`TryPush` / `TryPoll`), context-aware
blocking (`PushWithContext` / `PollWithContext`), an `Option`-based
return for non-blocking reads, a single generic type for "a
fixed-capacity blocking queue", or `errgroup`-style structured
concurrency with a typed `fn(ctx)` signature — the
[`concurrency` package](./docs/concurrency/README.md) ships
[`BoundedBlockingQueue[T]`](./docs/concurrency/README.md#boundedblockingqueuet)
(a thin wrapper around `chan T`, per-op overhead ≈ 0),
[`UnboundedBlockingQueue[T]`](./docs/concurrency/README.md#unboundedblockingqueuet)
(ring buffer + mutex + cond, since the Go runtime has no "unbounded
buffered channel"),
[`Group`](./docs/concurrency/README.md#group) (strict or best-effort
structured concurrency on top of `sync.WaitGroup`),
[`Semaphore`](./docs/concurrency/README.md#semaphore) (counting
semaphore with fixed unit weights), and
[`Pool[T]`](./docs/concurrency/README.md#poolt) (typed `sync.Pool`).

Reach for `BoundedBlockingQueue` when you want backpressure via
capacity; `UnboundedBlockingQueue` when `Push` must never block;
`Group` when you have a fan-out of goroutines and want the first
error or aggregated failures; `Semaphore` to cap concurrent access
to a resource; `Pool[T]` for temporary object reuse.

## Error handling

`Result[T]` is the typed `(T, error)` carrier:

```go
v, err := adt.Success(42).Unwrap()
if err != nil { return err }

port, _ := adt.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err is always nil
```

`Recover` always returns a `T`; surface a fresh error by ending the
chain with `Unwrap` instead.

For collection pipelines that need to surface errors, the recommended
pattern is to convert the error path into an absent `Option[T]` (e.g.
`OfNullable` on a lookup that returns `(T, error)`) and keep the
success path on the regular fluent API.

## Documentation

Per-package API reference and examples, in English and Chinese:

- [docs/README.md](./docs/README.md) — index, with per-module roadmaps
- [collections](./docs/collections/README.md) — `ArrayList`, `LinkedList`, `HashMap`, `HashSet`, `Stack`, `Queue`, `Deque`, `PriorityQueue`, `Stream`, `Range`
- [adt](./docs/adt/README.md) — `Option[T]`, `Result[T]`, `Either[L, R]`
- [reactivex](./docs/reactivex/README.md) — `Observable`, `Subject`, backpressure, operators
- [concurrency](./docs/concurrency/README.md) — blocking queues, `Group`, `Semaphore`, `Pool[T]`
- [control](./docs/control/README.md) — `If`/`IfGet`, `Repeat`/`RepeatE`, pattern matching
- [utils/objects](./docs/utils/objects/README.md) — `IsNil`, `Equals`
- [utils/json](./docs/utils/json/README.md) — `Encode`/`Decode` on `encoding/json/v2`

## Stability and SemVer

Every module is published under [Semantic Versioning 2.0.0](https://semver.org/).
Pre-`v1.0.0` numbering is chosen on purpose:

- **`0.0.x` (current).** No API stability guarantee. Anything that
  touches a public signature, exported type, or the documented
  behaviour of an operator is a **minor** bump (`0.0.x` →
  `0.(x+1).0`). Patches are bug fixes only.
- **`0.y.0`** (planned once a module reaches feature freeze). Public
  types, function signatures, and operator behaviour are frozen.
  Promoted when the API has been exercised for at least one minor
  cycle without changes.
- **`v1.0.0`** is reserved for modules the maintainer is willing to
  backport bug fixes to.

### What counts as a breaking change

- Renaming an exported type or method.
- Changing a generic receiver type (e.g. swapping `Map(func(T) R)` for
  `Map(func(context.Context, T) (R, error))`).
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
`utils/v0.0.2` may ship the same week as `reactivex/v0.0.1`, or years
later; the prefixes are independent.

Every release notes a `BREAKING:` line for any signature or behaviour
change. The diff between the previous and current `vX.Y.Z` tag is the
authoritative changelog.

## Go version

The current modules use **Go 1.27**:

```text
go 1.27.1
```

Go 1.23 introduced `iter.Seq`, `iter.Seq2`, and `for range` support
for function iterators. Go 1.27's generic methods make fluent APIs
such as `Stream[T].Map[R]`, `Option[T].Map[R]`, and
`Result[T].Map[R]` possible.

## Roadmap

**Done** across modules:

- `ArrayList[T]`, `LinkedList[T]`, `HashMap[K, V]`, `HashSet[T]` with
  fluent methods.
- Eager collection operations: `Filter`, `Map[R]`, `FlatMap[R]`,
  `Reduce[R]`, `Collect`, `ForEach`, `Peek`, `Concat`.
- Order-dependent operations on lists: `Get`, `Insert`, `RemoveAt`,
  `First`, `Last`, `Find`, `Take`, `Drop`, `Distinct`, `SortBy`,
  `MinBy`, `MaxBy`.
- Option-based access on `ArrayList` / `LinkedList` / `Stack` /
  `Queue` / `Deque` / `PriorityQueue` (every "take one" returns
  `Option[T]`).
- `Stream[T]` adapter backed by `iter.Seq[T]`, with early-terminating
  terminals.
- Synchronous queues under `collections/queues`: `Queue[T]`,
  `Deque[T]`, `PriorityQueue[T]` (the latter with bounded top-K
  semantics).
- `Option[T]` / `Result[T]` / `Either[L, R]` / `IsNil[T]` / `Equals[T]`.
- `control.If` / `IfGet`, `Repeat` / `RepeatE`, and `control/match`
  pattern matching.
- `reactivex.Observable[T]` / `Publisher[T]` / `Subscriber[T]`,
  `Subject[T]`, `Single[T]`, `Maybe[T]`, backpressure, operators.
- `concurrency.BoundedBlockingQueue[T]`,
  `UnboundedBlockingQueue[T]`, `Group`, `Semaphore`, `Pool[T]`.
- `utils/json` Result-style codec on top of `encoding/json/v2`.
- Bounded memory: head-offset `ArrayList` with periodic compaction;
  slot-zeroing on `Stack.Pop`, `Queue.Pop`, `Deque.PopFront`,
  `Deque.PopBack`, and the list `Remove*` paths.
- Tests with the race detector on every module. Live coverage is
  reported per module via Codecov — see the badge above.

**Open** (next-tier candidates; promoted on demand, not preemptively):

- Error-aware collection operations (`MapE`, `FilterE`, `CollectE`)
  as a first-class pipeline alternative to `Result[T]`-per-element.
- Iterators (`iter.Seq[T]`) as a first-class output of the collection
  types, parallel to `Stream()`.
- Per-module benchmarks (see each module's `roadmap.md` for the
  specifics).

Per-module roadmaps live in each module's `docs/<module>/roadmap.md`.

## Contributing

The toolkit is intentionally small, but contributions are welcome.

- **Scope per module.** Each top-level package and sub-package ships
  as an independent Go module with its own version tag prefix
  (see [Stability and SemVer](#stability-and-semver)). A PR that
  touches more than one module should call that out explicitly and
  explain why cross-module coordination is needed; otherwise the
  maintainer will ask you to split it.
- **Discuss before you build.** For anything beyond a small bug fix
  or documentation typo, file an issue first so the shape of the
  change can be agreed on. Things that look like "just adding a
  method" often turn into SemVer / API-surface decisions.
- **Tests are required.** A PR that changes behaviour must come with
  `go test ./...` passing **with `-race`** on the affected module.
- **Linting must stay clean.** Run `golangci-lint run ./...` on the
  module you changed before pushing.
- **Public API is sacred.** If your change touches an exported type,
  method, function, or constant, the PR description must call it out
  explicitly and tag it `BREAKING:`, `feat:`, or `fix:` so the release
  notes can be written correctly.
- **Commit hygiene.** One logical change per commit. Bug fixes and
  refactors should not be mixed with feature work in the same commit.

The standard flow: fork the repository and create a topic branch off
`main`; make your change; run `go test -race ./...` and
`golangci-lint run`; push the branch; open a pull request against
`qianwj/typed:main`.

## Bug reports

Use the [issue tracker](https://github.com/qianwj/typed/issues).
Search to make sure the bug is not already filed (including closed
ones — the fix may have landed on a different module's branch and not
yet been tagged).

A good bug report includes:

- **Goal.** One sentence: what were you trying to do with Typed?
- **Module and version.** `git rev-parse HEAD` inside the module
  directory, or the module's `vX.Y.Z` tag.
- **Environment.** `go version`, `go env GOOS GOARCH`.
- **Reproduction.** The smallest self-contained snippet that
  triggers the bug. The preferred form is a Go test or a
  `go run`-able file.
- **Expected vs actual.** What you expected, what actually happened,
  and the exact error / panic / log output. For panics, include the
  full stack trace.
- **Race / data-race reports.** If you suspect a data race, say so
  explicitly and include the `-race` output verbatim.

Security issues should **not** go through the public tracker. Use
GitHub's [private security reporting](https://github.com/qianwj/typed/security/advisories/new)
flow.

## License

[MIT](./LICENSE).
