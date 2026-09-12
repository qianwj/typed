# `typed/concurrency` — Typed

Concurrency primitives that complement Go's standard library, written in the same style as the rest of the toolkit: concrete generic types, no `any` round-trips, no reflective tricks.

Right now the package ships one type:

- `BoundedBlockingQueue[T]` — a fixed-capacity FIFO queue with blocking and non-blocking variants, backed by a single ring buffer, a single `sync.Mutex`, and a single `*sync.Cond`.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [Why a custom queue?](#why-a-custom-queue)
- [`BoundedBlockingQueue[T]`](#boundedblockingqueuet)
  - [Construction](#construction)
  - [Blocking variants](#blocking-variants)
  - [Non-blocking variants](#non-blocking-variants)
  - [Observability](#observability)
  - [Memory model](#memory-model)
  - [Comparison with `chan T`](#comparison-with-chan-t)
  - [Examples](#examples)
- [Benchmarks](#benchmarks)
- [See also](#see-also)

## Import

```go
import (
    "github.com/qianwj/typed/concurrency"
    "github.com/qianwj/typed/utils/option"
)
```

`concurrency` depends on `utils/option` because [`BoundedBlockingQueue.TryTake`](#non-blocking-variants) returns an `option.Optional[T]`, matching the rest of the toolkit's "may be absent" convention.

The package is its own `go.mod` module; import it independently of `collections` / `reactivex` / `control` / `utils`.

## Why a custom queue?

Go's built-in `chan T` is a perfectly good bounded blocking queue — when you have a capacity. The runtime implements it with per-P (processor-local) queues and lock-free fast paths, so it is hard to beat on raw throughput.

You reach for `BoundedBlockingQueue[T]` when a channel's surface area is not enough:

- **Non-blocking probes.** `TryPush` / `TryTake` let you peek-and-decide without spawning a goroutine to time out a `select`. `TryTake` returns an `option.Optional[T]` in the same shape as `Stack.Pop` / `Queue.Pop`, so the result composes with the rest of the toolkit.
- **Synchronous observability.** `Size()` and `Capacity()` give you a consistent snapshot under the same lock that serialises the queue, so you can branch on backpressure without a separate counter.
- **A uniform generic API.** When the rest of the codebase already imports the Typed toolkit, the queue's signature looks and feels like the rest of it: `Size` instead of `len`, `Optional` instead of `(T, bool)`.

The trade-off is throughput. On a single-producer / single-consumer microbenchmark the native channel runs roughly **5–7× faster**. The array-backed ring buffer keeps the gap from being worse: it has **zero per-operation allocations** in either implementation, so GC pressure is not the differentiator — sync overhead is.

`BoundedBlockingQueue` is therefore positioned as: **"when `chan T` is not expressive enough, and the throughput tax is acceptable"**.

---

## `BoundedBlockingQueue[T]`

A fixed-capacity FIFO queue. The capacity is set at construction time and never changes. All operations are safe for concurrent use.

### Construction

```go
// NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T]
q := concurrency.NewBoundedBlockingQueue[Job](1024)
```

`capacity` must be positive. Passing `0` or a negative value panics — a misconfigured capacity should fail loudly, not silently produce a queue that always blocks.

The actual capacity of the queue is the **smallest power of two greater than or equal to** the value passed in. So `NewBoundedBlockingQueue[T](100)` returns a queue with `Capacity() == 128`, and `NewBoundedBlockingQueue[T](1024)` returns one with `Capacity() == 1024`. The rounding lets `Push` / `Take` advance head and tail with a single bitwise `&` against a mask instead of a modulo, which is the reason for the policy. If you need an exact capacity, pass the power of two yourself.

The zero value is not usable; always go through the constructor.

### Blocking variants

| Method | Behaviour |
| --- | --- |
| `Push(data T)` | Enqueue. **Blocks** while the queue is full; wakes as soon as a slot frees up. |
| `Take() T` | Dequeue and return. **Blocks** while the queue is empty; wakes as soon as an element arrives. |

Blocking uses a single `*sync.Cond` shared with the mutex. Every state change (`Push`, `Take`, `TryPush`, `TryTake`, `DrainTo`) calls `cond.Broadcast()`, which is the conservative choice: it always wakes the right kind of waiter (pushers after a `Take` or `DrainTo`, takers after a `Push`) without needing two separate conditions or a thundering-herd-friendly `Signal()`. For a queue that takes a single mutex anyway, this is the right amount of cleverness.

### Non-blocking variants

| Method | Behaviour |
| --- | --- |
| `TryPush(data T) bool` | Enqueue if a slot is free. Returns `true` on success, `false` immediately if the queue is full. |
| `TryTake() option.Optional[T]` | Dequeue if anything is available. Returns a present `Optional` on success, an empty `Optional` immediately if the queue is empty. |
| `DrainTo(dst []T) int` | Dequeue up to `len(dst)` elements and write them into `dst` in FIFO order. Returns the number of elements actually written. `dst` is not grown; if the queue has fewer elements than `len(dst)`, the remaining slots in `dst` are left untouched. |

These never wait, which is what you want for `select { ... default: ... }` style logic and for backpressure policies that prefer "drop the work" or "shed load" over "block the caller". `TryTake` returns an `Optional` rather than a `(T, bool)` pair because that is the toolkit-wide convention for "may be absent", shared with `Stack.Pop`, `Queue.Pop`, and `Deque.PopFront` / `PopBack`.

`DrainTo` exists for the case where you want to flush a batch in a single critical section — taking N elements one at a time would take the mutex N times and would not give you a consistent snapshot across the batch. `DrainTo` does both in one pass.

### Observability

| Method | Behaviour |
| --- | --- |
| `Size() int` | Current number of elements. Takes the same mutex as the queue, so the value is consistent with the next operation you perform after `Size()` returns. Matches `Stack.Size`, `Queue.Size`, `ArrayList.Size`, etc. |
| `Capacity() int` | Configured capacity. Lock-free — capacity is immutable after construction. Named to match the constructor parameter, since no other Typed type has a fixed capacity. |

### Memory model

- **Backing storage.** A single `make([]T, capacity)` allocated up front. The head and tail indices wrap around it. There are no per-element allocations on the hot path.
- **Power-of-two capacity.** The constructor rounds the requested capacity up to the next power of two and stores `cap - 1` as a precomputed mask. `head` and `tail` advance via `& mask` — a single bitwise operation — rather than `% cap`. The `count` field disambiguates "empty" from "full", so the head-equals-tail case is never ambiguous.
- **Slot zeroing.** `Take`, `TryTake`, and `DrainTo` overwrite the freed slots with the zero value of `T` before advancing `head`. This is purely a GC hint: if `T` contains pointers, the popped values can be collected even while the queue still owns its backing array. The slot is never read again until a future `Push` overwrites it.
- **No pointers to internal state escape.** The slice is held by value inside the struct; concurrent callers go through the mutex.

### Comparison with `chan T`

| | `chan T` | `BoundedBlockingQueue[T]` |
| --- | --- | --- |
| Capacity | set at `make` | set at `NewBoundedBlockingQueue`; rounded up to the next power of two |
| Bounded by default | no (`make(chan T)` is unbuffered) | yes — capacity is required |
| Blocking send / receive | `ch <- v` / `<-ch` | `Push(v)` / `Take()` |
| Non-blocking probe | none — wrap in `select { default: }` | `TryPush() bool` / `TryTake() option.Optional[T]` |
| Batch dequeue | none — drain in a loop | `DrainTo(dst []T) int` |
| `len(ch)` | yes, but not synchronised with ops | `Size()` under the same lock |
| Per-op allocations | 0 | 0 |
| Throughput (1P1C) | ~30 ns/op | ~200 ns/op |
| Throughput (MPMC 8w) | ~22 ns/op | ~118 ns/op |
| Batch dequeue (64) | n/a | ~310 ns/op (zero allocs) |

> Numbers from `go test -bench` on Apple M4 Pro, Go 1.27.1, queue capacity 1024 (1P1C) / 64 (MPMC). They are in the right ballpark to confirm the queue is competitive with hand-rolled `sync.Cond` implementations; they are not a substitute for measuring on your workload.

### Examples

#### Single producer, single consumer, backpressure

```go
jobs := concurrency.NewBoundedBlockingQueue[*Job](64)

go func() {
    for {
        j := q.Take()        // blocks until a job arrives
        handle(j)
    }
}()

for _, raw := range sources {
    jobs.Push(raw)            // blocks if the consumer falls behind
}
```

The blocking `Push` provides natural backpressure: a slow consumer throttles the producer without any extra protocol.

#### Non-blocking offer

```go
q := concurrency.NewBoundedBlockingQueue[Event](1024)

func submit(e Event) bool {
    if !q.TryPush(e) {
        metrics.DroppedCounter.Inc()
        return false
    }
    return true
}
```

`TryPush` lets you implement "drop on overflow" or "shed load" policies without spawning a watcher goroutine.

#### Inspecting the queue

```go
if q.Size() > q.Capacity() * 9 / 10 {
    log.Printf("queue is %.0f%% full", float64(q.Size()) / float64(q.Capacity()) * 100)
}
```

`Size()` is taken under the same mutex as `Push` / `Take`, so the value is consistent with the operation that follows.

#### Producer / consumer with `TryTake` for graceful shutdown

```go
stop := make(chan struct{})

go func() {
    for {
        select {
        case <-stop:
            return
        default:
        }
        if opt := q.TryTake(); !opt.IsEmpty() {
            process(opt.Get())
            continue
        }
        // queue empty: yield and let producers make progress
        runtime.Gosched()
    }
}()

// later, to stop the producer:
close(stop)
```

This pattern drains the queue without blocking on `Take()` so the consumer can shut down promptly when the producer is done.

#### Batched drain

When you want to process a whole batch under a single critical section — for example, flushing a worker before a flush-to-disk tick — `DrainTo` is cheaper and more consistent than looping `TryTake`:

```go
buf := make([]Event, 64) // size the buffer for your average batch
for {
    n := q.DrainTo(buf)
    if n == 0 {
        time.Sleep(idleInterval) // or block on Take if shutdown isn't a concern
        continue
    }
    flushBatch(buf[:n])
}
```

Two things to know:

- `DrainTo` does not grow `dst` — if the queue has more elements than `len(dst)`, the surplus stays in the queue. Pick `len(dst)` to match the batch size you actually want to process.
- `DrainTo` calls `cond.Broadcast()`, so any `Push` blocked on a full queue wakes up as soon as the drain finishes.

---

## Benchmarks

The package ships two pairs of benchmarks: `BoundedBlockingQueue_*` versus the equivalent `chan int` baseline. Run them with:

```bash
GOWORK=off go -C concurrency test -bench=. -benchmem -run=^$ .
```

On an Apple M4 Pro (Go 1.27.1, darwin/arm64) the 1P1C case lands around **195 ns/op** for `BoundedBlockingQueue` versus **28 ns/op** for `chan int`; the 8-worker MPMC case lands around **118 ns/op** versus **22 ns/op**. Both implementations report `0 B/op` and `0 allocs/op`.

The point of these benchmarks is to confirm the array-backed design does not regress on GC pressure. They are not a substitute for measuring on your own workload.

## Future work cancellation support

`BoundedBlockingQueue` does not currently expose context-aware variants of `Push` and `Take`. This section captures the design so the decision is not relitigated the next time someone needs cancellation.

### Proposed signatures

```go
// PushCtx blocks until there is space or ctx is canceled.
// Returns nil on success, ctx.Err() on cancellation. The element is
// not enqueued in the cancellation case.
func (q *BoundedBlockingQueue[T]) PushCtx(ctx context.Context, data T) error

// TakeCtx blocks until an element is available or ctx is canceled.
// Returns (zero, ctx.Err()) on cancellation.
func (q *BoundedBlockingQueue[T]) TakeCtx(ctx context.Context) (T, error)
```

These mirror the established Go pattern (`http.Request.WithContext`, `sql.QueryContext`, `(*sync.WaitGroup).Wait` analogously via `chan`) and let callers compose with timeouts, deadlines, and shutdown signals without inventing their own goroutine-and-channel dance on top of `Push` / `Take`.

### The blocker: `sync.Cond` does not speak `context.Context`

`Push` and `Take` use `sync.Cond.Wait()` to suspend a goroutine until a `Broadcast()` arrives. `Cond.Wait` is the cleanest possible primitive for "go to sleep on a mutex-protected condition variable" but it has no cancellation hook: the only way out of `cond.Wait()` is a `Broadcast()` / `Signal()`. You cannot `select` on it. So adding `PushCtx` / `TakeCtx` requires either of:

**Option A — Bridge with a per-call watcher goroutine.**

Each `PushCtx` call spawns a goroutine that does `<-ctx.Done(); q.cond.Broadcast()` and exits. A `Broadcast` is cheap, but every blocking call now costs one extra goroutine and one extra channel — both of which show up in `pprof` and as allocator pressure under load. Acceptable for occasional cancellable calls, wasteful as the default path.

**Option B — Replace `sync.Cond` with channel-based signalling.**

The queue grows two buffered (cap 1) `chan struct{}` signals: `notEmpty` and `notFull`. `Push` / `PushCtx` send on `notEmpty`; `Take` / `TakeCtx` send on `notFull`. Blocking then becomes a `select` on the signal channel plus `ctx.Done()`. This is what Go's own `sync/semaphore` does and is the right shape for a context-aware API.

The catch: the basic `Push` / `Take` paths pay a small per-op cost for going through a channel send instead of a `Cond.Broadcast` (a buffered cap-1 channel send is essentially free when the slot is empty, so the steady-state cost is one cache-line write per op). The bigger cost is the rewrite — `Push` / `Take` / `TryPush` / `TryTake` / `DrainTo` all interact with the wait queue, and the lock-then-signal-then-unlock ordering has to be reviewed carefully to avoid missed wakeups and the "thundering herd on Broadcast" pattern.

### Why it is deferred

The package currently has zero callers of the ctx-aware variants. Adding them under Option A would put every future reader of this code through the "why is there an extra goroutine per call?" detour. Option B is the right answer but is large enough that doing it without a real use case risks getting the broadcast ordering subtly wrong (which then fails only on the contended test that nobody runs in CI).

The plan is to land Option B as a single rewrite of the synchronisation core, with a test that exercises the "many pushers, cancel the consumer's context mid-flight, all pushers must observe cancellation without leaking the broadcast" case, before exposing `PushCtx` / `TakeCtx` publicly. Until then, callers that need cancellation can wrap a regular `Push` / `Take` in their own goroutine and select on `ctx.Done()`.

## See also

- [`reactivex`](../../reactivex/README.md) — typed async event streams with explicit demand and configurable backpressure. If the work being queued is "events to deliver to many subscribers", an `Observable` is usually a better fit than a queue.
- [`collections.Queue`](../collections/README.md#stackt--queuet--dequet) — the synchronous, in-memory `Queue[T]`. Use it when there is no concurrency and you want `Optional[T]`-based access; it has no locking, no `TryPush`, and no backpressure.
- [`sync.Cond`](https://pkg.go.dev/sync#Cond) and [`chan T`](https://go.dev/ref/spec#Channel_types) in the standard library — the building blocks this type composes.
