# `typed/concurrency` — Roadmap

This is the plan for the [`concurrency`](./README.md) package. The
package is intentionally small — Go's [`sync`](https://pkg.go.dev/sync)
already covers the basics — so each entry below has to earn its keep.

The conventions used by every new type here:

- Generic over `T any` (no `any`-round-trips at the API surface).
- `Push` / `Poll`-style verbs for queue-style ops, `Acquire` /
  `Release` for lock-style ops, `Await` for one-shot waits.
- `WithContext` suffix for ctx-aware blocking variants.
- `Try` prefix for the non-blocking variant.
- Optional returns use `option.Optional[T]` (matching
  [`collections.Queue`](../collections/README.md#stackt--queuet--dequet)
  and the rest of the toolkit).
- `*Or[T]` helpers for tests where the convention is "(value, ok)".
- Zero value not usable; construct with `New*`.

## Done

The package ships these types today:

- **`BoundedBlockingQueue[T]`** — fixed-capacity FIFO blocking queue,
  implemented as a thin generic wrapper around `chan T`. `Push` /
  `PushWithContext` / `TryPush` to enqueue; `Poll` /
  `PollWithContext` / `TryPoll` to dequeue. `Capacity()` returns the
  power-of-two rounded upper bound.
- **`UnboundedBlockingQueue[T]`** — unbounded FIFO blocking queue.
  `Push` never blocks; `Poll` blocks when empty. Ring buffer +
  `sync.Mutex` + `*sync.Cond` under the hood, because the Go runtime
  has no "unbounded buffered channel". `PushWithContext` is a thin
  ctx-check over `Push`; `PollWithContext` uses a one-shot watcher
  goroutine that broadcasts the cond on `ctx.Done()`.
- **`Group`** — structured concurrency with two failure policies.
  Built directly on `sync.WaitGroup` (no `errgroup` dependency). In
  `Strict` mode (default) any task error cancels the group's ctx and
  `Wait()` returns that error; in `BestEffort` mode tasks run to
  completion and `Wait()` returns nil if every task succeeded, the
  bare single error if exactly one task failed, or a
  `*BestEffortError` aggregating all failures otherwise.
  `WithLimit(n)` (set at construction) caps concurrent tasks via a
  semaphore channel.

All three share the same API conventions and pass `-race` clean.

## Tier 1 — synchronization primitives

These two are small and address gaps in the standard library's
"goroutine coordination" surface.

### `Semaphore` — counting semaphore

Limit concurrent access to a resource (a downstream service, a
DB, a remote API). `golang.org/x/sync/semaphore` exists but uses
`int64` weights and is awkward for the simple "N slots" case. A
typed version with `Acquire` / `Release` is more idiomatic.

```go
type Semaphore struct { /* ... */ }

func NewSemaphore(n int) *Semaphore
func (s *Semaphore) Acquire()
func (s *Semaphore) AcquireWithContext(ctx context.Context) error
func (s *Semaphore) TryAcquire() bool
func (s *Semaphore) Release()
func (s *Semaphore) Available() int
```

Typical use: rate limit at the goroutine level instead of the
request level. The toolkit's existing `BoundedBlockingQueue` already
gives per-queue backpressure; `Semaphore` is for the broader
"max N concurrent goroutines touching X" case.

### `CountDownLatch` — one-shot N-party coordination

"Wait for N events to occur, then proceed." Java has this; Go's
`sync` does not. The primitive is small (one atomic counter + one
cond + a `Done` flag) and slots in alongside `Semaphore` in the
package.

```go
type CountDownLatch struct { /* ... */ }

func NewCountDownLatch(count int) *CountDownLatch
func (l *CountDownLatch) Await()
func (l *CountDownLatch) AwaitWithContext(ctx context.Context) error
func (l *CountDownLatch) CountDown()
func (l *CountDownLatch) Done() bool
```

Typical use: "wait for N subsystems to be ready before opening
traffic", "wait for the last writer to finish before a barrier
flush".

## Tier 2 — worker pattern

### `WorkerPool[T]` — bounded N workers consuming from a channel

A "batteries included" wrapper around `BoundedBlockingQueue[T]`
plus N goroutines that consume from it. Users can build this from
`BoundedBlockingQueue` + a `Group` themselves, but the dedicated
type saves boilerplate for the common "N workers + 1 (or many)
dispatchers" pattern.

```go
type WorkerPool[T any] struct { /* ... */ }

func NewWorkerPool[T any](workers int, handler func(T)) *WorkerPool[T]
func (p *WorkerPool[T]) Submit(data T) bool
func (p *WorkerPool[T]) SubmitWithContext(ctx context.Context, data T) error
func (p *WorkerPool[T]) Close() error // drains in-flight, then stops
```

`WorkerPool` will be implemented on top of `BoundedBlockingQueue` +
`Group` once both are in; the implementation is short and the test
surface is mostly inherited.

## Tier 3 — specialised queues

Both of these are useful but narrow. They will be added when a real
use case shows up, not preemptively — preemptive additions tend to
over-design the API and under-exercise the tests.

### `PriorityQueue[T]` — heap-backed priority queue

Same `Push` / `Poll` / `TryPoll` verbs as `BoundedBlockingQueue`,
plus `PushWithPriority(data T, priority int)`. Heap implementation;
the interface stays small so the package doesn't grow a
container-style toolkit by accident.

Use case: schedulers, leader election, anything with "process the
most important thing first".

### `DelayQueue[T]` — items become available after a delay

Elements carry an `availableAt time.Time`; `Poll` blocks until the
head element is ready.

Use case: scheduled tasks, retry with backoff.

## Out of scope

Considered and explicitly rejected:

- **`Future[T]` — typed async result.** Go's `go` keyword + a
  buffered channel already gives you one-shot async value delivery
  in two lines; `Group` with one task adds cancellation and error
  aggregation on top. A dedicated `Future` type would just be a
  re-skinned wrapper around the same machinery, with no new
  capability. Skip unless a real ergonomic gap appears.
- **Rate limiter / token bucket.** Different domain (time math,
  refill semantics), different testability needs (clock injection).
  Belongs in its own package — a future `typed/ratelimit` is the
  natural home.
- **Pipeline composition.** Overlaps with
  [`reactivex`](../../reactivex/README.md). The
  [`Stream`](../collections/README.md#streamt--pipelines) type
  already covers the common "compose stages of work" use case.
- **Mutex with timeout.** Complex, error-prone, rarely needed.
  `Semaphore` covers the "I want bounded contention" use case; if you
  need a timed lock, that's a different primitive and probably not
  what you actually want.
- **`Pool[T]` object pool.** Mostly subsumed by
  [`sync.Pool`](https://pkg.go.dev/sync#Pool) with a tiny generic
  helper, and the toolkit's existing `Pool` in
  [`collections`](../collections/README.md). Adding a third one
  here would just split the same idea across packages.

## Update policy

This roadmap is reviewed when:

- A tier is completed (move it from "Next" to "Done").
- A real use case appears for a "Tier 3+" item (promote it).
- A previously rejected item turns out to be needed (re-open it).

The expectation is that each tier adds ~100–200 lines of code +
tests + docs, and the whole package stays under ~1k lines including
tests. Beyond that, it's a sign that something belongs in its own
package, not in `concurrency`.