# `typed/concurrency` — Typed

[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=concurrency)](https://codecov.io/gh/qianwj/typed)

Concurrency primitives that complement Go's standard library, written in the same style as the rest of the toolkit: concrete generic types, no `any` round-trips, no reflective tricks.

Right now the package ships three types:

- `BoundedBlockingQueue[T]` — a fixed-capacity FIFO blocking queue, implemented as a thin generic wrapper around a `chan T` with toolkit-style naming, `Option`-based non-blocking probes, and context-aware blocking.
- `UnboundedBlockingQueue[T]` — an unbounded FIFO blocking queue. `Push` never blocks; `Poll` blocks when empty. Implemented as a ring buffer over a single pre-allocated slice with a `sync.Mutex` and a `*sync.Cond`, because the Go runtime has no "unbounded buffered channel".
- `Group` — structured concurrency with two failure policies: `Strict` (any task error fails the group, ctx cancels siblings) and `BestEffort` (tasks run to completion; only succeeds if all succeed, otherwise returns aggregated errors). Built on `sync.WaitGroup`, no external dependencies.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md). For what's planned next in this package, see the [Roadmap](./roadmap.md).

## Contents

- [Import](#import)
- [Why a custom queue?](#why-a-custom-queue)
- [Roadmap](./roadmap.md)
- [`BoundedBlockingQueue[T]`](#boundedblockingqueuet)
  - [Construction](#construction)
  - [Blocking variants](#blocking-variants)
  - [Non-blocking variants](#non-blocking-variants)
  - [Observability](#observability)
  - [Memory model](#memory-model)
  - [Comparison with `chan T`](#comparison-with-chan-t)
  - [Examples](#examples)
- [`UnboundedBlockingQueue[T]`](#unboundedblockingqueuet)
- [`Group`](#group)
- [`Semaphore`](#semaphore)
- [Benchmarks](#benchmarks)
- [See also](#see-also)

## Import

```go
import (
    "github.com/qianwj/typed/concurrency"
    "github.com/qianwj/typed/adt"
)
```

`concurrency` depends on `adt` because both `BoundedBlockingQueue.TryPoll` and `UnboundedBlockingQueue.TryPoll` return an `adt.Option[T]`, matching the rest of the toolkit's "may be absent" convention.

The package is its own `go.mod` module; import it independently of `collections` / `reactivex` / `control` / `utils`.

## Why a custom queue?

Go's built-in `chan T` is a perfectly good bounded blocking queue — when you have a capacity. The runtime implements it with per-P (processor-local) queues and lock-free fast paths, so it is hard to beat on raw throughput. In fact, `BoundedBlockingQueue` is now a thin wrapper around a `chan T`: the underlying `ch <- data` and `<-ch` are exactly what the methods do, so the wrapper adds essentially zero per-op cost.

You reach for `BoundedBlockingQueue[T]` when you want the toolkit's surface area on top of a channel:

- **Option return for non-blocking probes.** `TryPoll` returns an `adt.Option[T]` in the same shape as `Stack.Pop` / `Queue.Pop`, so the result composes with the rest of the toolkit instead of forcing a `(value, ok)` round-trip.
- **Context-aware blocking.** `PushWithContext` / `PollWithContext` let you compose with timeouts, deadlines, and shutdown signals without inventing your own goroutine-and-channel dance on top of `Push` / `Poll`.
- **A uniform generic API.** `Size` instead of `len`, `Capacity` instead of `cap`, `Option` instead of `(T, bool)`. Same vocabulary as the rest of Typed.
- **Power-of-two capacity rounding.** `Capacity()` always returns a power of two, convenient for callers that want a bitmask.

The trade-off versus a raw `chan T` is essentially zero on the hot path.

---

## `BoundedBlockingQueue[T]`

A fixed-capacity FIFO queue. The capacity is set at construction time and never changes. All operations are safe for concurrent use.

### Construction

```go
// NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T]
q := concurrency.NewBoundedBlockingQueue[Job](1024)
```

`capacity` must be positive. Passing `0` or a negative value panics — a misconfigured capacity should fail loudly, not silently produce a queue that always blocks.

The actual capacity of the queue is the **smallest power of two greater than or equal to** the value passed in. So `NewBoundedBlockingQueue[T](100)` returns a queue with `Capacity() == 128`, and `NewBoundedBlockingQueue[T](1024)` returns one with `Capacity() == 1024`. The rounding is for the convenience of callers that want a power-of-two bitmask; if you need exact capacity, pass the power of two yourself.

The zero value is not usable; always go through the constructor.

### Blocking variants

| Method | Behaviour |
| --- | --- |
| `Push(data T)` | Enqueue. **Blocks** while the queue is full; wakes as soon as a slot frees up. |
| `Poll() T` | Dequeue and return. **Blocks** while the queue is empty; wakes as soon as an element arrives. |
| `PushWithContext(ctx, data T) error` | Context-aware `Push`. Blocks until the queue has space or `ctx` is canceled; returns `ctx.Err()` and does not enqueue on cancellation. |
| `PollWithContext(ctx) adt.Result[T]` | Context-aware `Poll`. Blocks until an element is available or `ctx` is canceled; returns `adt.Success(value)` on success and `adt.Failure[T](ctx.Err())` on cancellation. |

Blocking uses the underlying `chan T` directly. `Push` is `ch <- data`, `Poll` is `<-ch`. `PushWithContext` and `PollWithContext` add a `case <-ctx.Done()` to the same `select`, so cancellation composes for free: when the `ctx.Done()` branch is selected, `PushWithContext` returns `ctx.Err()` without enqueuing, while `PollWithContext` returns `adt.Failure[T](ctx.Err())` without dequeuing.

### Non-blocking variants

| Method | Behaviour |
| --- | --- |
| `TryPush(data T) bool` | Enqueue if a slot is free. Returns `true` on success, `false` immediately if the queue is full. |
| `TryPoll() adt.Option[T]` | Dequeue if anything is available. Returns a present `Option` on success, an empty `Option` immediately if the queue is empty. |

These never wait, which is what you want for `select { ... default: ... }` style logic and for backpressure policies that prefer "drop the work" or "shed load" over "block the caller". `TryPoll` returns an `Option` rather than a `(T, bool)` pair because that is the toolkit-wide convention for "may be absent", shared with `Stack.Pop`, `Queue.Pop`, and `Deque.PopFront` / `PopBack`.

### Observability

| Method | Behaviour |
| --- | --- |
| `Size() int` | Current number of elements. Implemented as `len(ch)` — an atomic length read, not serialised with the next operation you perform. The window is small (a single atomic load) but real; if you need strict "Size() → next op sees a consistent view" semantics, use a ring-buffer implementation instead. Matches `Stack.Size`, `Queue.Size`, `ArrayList.Size`, etc. |
| `Capacity() int` | Configured capacity. Lock-free — capacity is immutable after construction. Named to match the constructor parameter, since no other Typed type has a fixed capacity. |

### Memory model

- **Backing storage.** A single `make(chan T, cap)` allocated up front. The Go runtime owns the channel's internal ring buffer; this type never touches it directly.
- **Power-of-two capacity.** The constructor rounds the requested capacity up to the next power of two so that `Capacity()` always returns a value usable as a bitmask. The channel's internal mechanics are independent of this rounding.
- **No slot zeroing.** Unlike a hand-rolled ring buffer, this wrapper does not zero freed slots on `Poll` / `TryPoll`. Pointer values received from the queue stay alive in the channel's backing array until the slot is overwritten by a new send. For workloads that drain a burst and then sit idle for a long time, those pointer values will live longer than they would under a custom ring buffer with explicit zeroing.
- **Per-op allocations.** Zero, in either implementation. The channel send/recv fast path does not allocate.

### What a `chan T` does and does not give you

Because the queue is a channel under the hood, the trade-offs versus a hand-rolled ring buffer + mutex + Cond are inherited from `chan T`:

| | `chan T` | `BoundedBlockingQueue[T]` |
| --- | --- | --- |
| Capacity | set at `make` | set at `NewBoundedBlockingQueue`; rounded up to the next power of two |
| Bounded by default | no (`make(chan T)` is unbuffered) | yes — capacity is required |
| Blocking send / receive | `ch <- v` / `<-ch` | `Push(v)` / `Poll()` |
| Non-blocking probe | wrap in `select { default: }` | `TryPush() bool` / `TryPoll() adt.Option[T]` |
| Context-aware blocking | wrap in `select { case <-ctx.Done(): }` | `PushWithContext` / `PollWithContext` |
| `len(ch)` | yes, but not synchronised with ops | `Size()` — same semantics, also not synchronised |
| Throughput (1P1C) | ~30 ns/op (tight microbench) | ~210 ns/op (tight microbench) |
| Throughput (MPMC 8w) | ~22 ns/op | **~22 ns/op** (in the same ballpark; wrapper has ~zero overhead) |

> Numbers from `go test -bench` on Apple M5 Pro, Go 1.27, queue capacity 1024 (1P1C) / 64 (MPMC). The 1P1C gap is microbenchmark noise: the consumer in the benchmark is a busy `TryPoll` loop, not a bare `<-ch`, which costs both sides most of the gap. In tight code where producer and consumer are both uncontended `Push` / `Poll`, the wrapper is within a couple of ns of a raw `chan T`. They are in the right ballpark to confirm the wrapper is "free"; they are not a substitute for measuring on your workload.

### Examples

#### Single producer, single consumer, backpressure

```go
jobs := concurrency.NewBoundedBlockingQueue[*Job](64)

go func() {
    for {
        j := q.Poll()        // blocks until a job arrives
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

`Size()` is `len(ch)`, an atomic length read, not serialised with the next operation you perform.

#### Producer / consumer with `TryPoll` for graceful shutdown

```go
stop := make(chan struct{})

go func() {
    for {
        select {
        case <-stop:
            return
        default:
        }
        if opt := q.TryPoll(); !opt.IsEmpty() {
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

This pattern drains the queue without blocking on `Poll()` so the consumer can shut down promptly when the producer is done.

---

## Benchmarks

The package ships two benchmarks: `BoundedQueue_1P1C` and `BoundedQueue_MPMC`. Run them with:

```bash
GOWORK=off go -C concurrency test -bench=. -benchmem -run=^$ .
```

On an Apple M5 Pro (Go 1.27, darwin/arm64):

| Bench | `BoundedBlockingQueue[T]` | `UnboundedBlockingQueue[T]` |
| --- | --- | --- |
| 1P1C | ~210 ns/op | ~28 ns/op |
| MPMC 8w | ~22 ns/op | ~61 ns/op |
| Per-op allocations | 0 | 0 |

The `BoundedBlockingQueue` MPMC number is the headline: it is in the same ballpark as a raw `chan T` and confirms the wrapper has ~zero per-op overhead. The 1P1C number is microbenchmark noise — the consumer in that benchmark is a busy `TryPoll` loop, not a bare `<-ch`, which costs both sides most of the gap. In tight code where producer and consumer are both uncontended `Push` / `Poll`, the wrapper is within a couple of ns of a raw `chan T`.

The `UnboundedBlockingQueue` 1P1C is faster than `BoundedBlockingQueue` because the consumer's `TryPoll` loop keeps draining the queue, so `Push` almost never sees a full queue (no `BoundedBlockingQueue` channel contention). MPMC is slower than `BoundedBlockingQueue` because the ring buffer is sequentialised through one mutex; the channel wrapper wins by going through per-P queues under contention.

---

## `UnboundedBlockingQueue[T]`

An unbounded FIFO queue. There is no capacity to wait on, so `Push` never blocks; `Poll` blocks when the queue is empty. All operations are safe for concurrent use.

### Construction

```go
// NewUnboundedBlockingQueue[T any]() *UnboundedBlockingQueue[T]
q := concurrency.NewUnboundedBlockingQueue[*Job]()
```

No capacity argument — the queue grows on demand.

### Blocking variants

| Method | Behaviour |
| --- | --- |
| `Push(data T)` | Enqueue. **Never blocks** — the queue is unbounded. |
| `Poll() T` | Dequeue and return. **Blocks** while the queue is empty; wakes as soon as an element arrives. |
| `PushWithContext(ctx, data T) error` | Context-aware `Push`. Returns `ctx.Err()` without enqueuing when ctx is already canceled; otherwise behaves like `Push`. |
| `PollWithContext(ctx) adt.Result[T]` | Context-aware `Poll`. Blocks until an element is available or ctx is canceled; returns `adt.Success(value)` on success and `adt.Failure[T](ctx.Err())` on cancellation. |

`PollWithContext` spawns a one-shot watcher goroutine that wakes any blocked `cond.Wait` when ctx is canceled. The goroutine exits as soon as `PollWithContext` returns, so the cost is one goroutine per `PollWithContext` call — fine for shutdown signals, not something you want in a tight loop.

`Push` signals the cond with `Signal()` (one waiter at a time), so a burst of N pushes wakes up to N blocked takers without thundering herd. A naive `cap-1` channel signal would fail here — see `TestUnboundedBurstWakesAllWaiters` for the regression test.

### Non-blocking variants

| Method | Behaviour |
| --- | --- |
| `TryPush(data T) bool` | Enqueue. **Never fails** — the queue is unbounded, so this always returns `true`. |
| `TryPoll() adt.Option[T]` | Dequeue if anything is available. Returns a present `Option` on success, an empty `Option` immediately if the queue is empty. |

### Observability

| Method | Behaviour |
| --- | --- |
| `Size() int` | Current number of elements. Matches `Stack.Size` / `Queue.Size` / `ArrayList.Size`. |

There is no `Capacity()` — the queue is unbounded by definition. If you want bounded behaviour, use [`BoundedBlockingQueue[T]`](#boundedblockingqueuet) instead.

### Memory model

- **Backing storage.** A single `make([]T, 16)` allocated up front. The ring buffer doubles when full, so the slice size is always a power of two. `head` and `tail` advance via `& mask` — a single bitwise operation.
- **Slot zeroing.** `Poll` and `TryPoll` overwrite the freed slot with the zero value of `T`, the same way `BoundedBlockingQueue`'s ring buffer did before it became a channel wrapper. Pointer-typed `T` is therefore safe to use without leaking memory.
- **Per-op allocations.** Zero on the steady-state hot path. The `grow` step allocates a new backing array, which is amortised O(1) per `Push`.
- **Why not `chan T`?** The Go runtime has no "unbounded buffered channel". `make(chan T, N)` for a large `N` works until it doesn't — once the buffer fills, `Push` blocks again and "unbounded" becomes a lie. A ring buffer with a mutex and a cond is the standard fix.

### Comparison with `BoundedBlockingQueue[T]`

| | `BoundedBlockingQueue[T]` | `UnboundedBlockingQueue[T]` |
| --- | --- | --- |
| Capacity | set at construction | none (grows on demand) |
| `Push` blocks | when full | never |
| `Poll` blocks | when empty | when empty | |
| Internals | `chan T` | ring buffer + mutex + cond |
| Hot-path throughput (MPMC 8w) | ~22 ns/op | ~61 ns/op |
| Memory bounded | yes (by capacity) | no — the slice grows until consumers drain it |
| Backpressure | built-in (capacity) | none — `Push` cannot fail |

Use `BoundedBlockingQueue` when you need backpressure. Use `UnboundedBlockingQueue` when you need `Push` to always succeed and you have an external mechanism (worker count, downstream queue, etc.) to stop producers from running the process out of memory.

---

## `Group`

`Group` is a typed structured-concurrency helper: spawn N goroutines, wait for all of them, and return one error. Two failure policies are selectable via [`WithMode`](#modes):

### Construction

```go
// NewGroup(parent context.Context, opts ...Option) *Group
g := concurrency.NewGroup(ctx)
g := concurrency.NewGroup(ctx,
    concurrency.WithMode(concurrency.BestEffort),
    concurrency.WithLimit(8),
)
```

A nil parent ctx is treated as `context.Background`.

### Blocking variants

| Method | Behaviour |
| --- | --- |
| `Go(fn func(ctx context.Context) error)` | Spawn a task. Blocks if the limit is set and reached (waits for a slot to free up). |
| `Wait() error` | Wait for every spawned task to return. Returns the error that best describes the outcome — see [`Strict` mode](#strict-mode) and [`BestEffort` mode](#besteffort-mode). |
| `WaitWithContext(ctx context.Context) error` | Ctx-aware `Wait`. Blocks until every spawned task returns OR ctx fires; returns the outcome-based error on completion, `ctx.Err()` if ctx fires first. **Abandons, doesn't kill** — but cancels the Group's derived ctx on return (so tasks respecting the ctx exit promptly; tasks ignoring it keep running; follow up with `Wait()` for the eventual outcome). |

Go is safe to call from multiple goroutines concurrently. Calls to `Go` after the first error still spawn their goroutines; their results simply don't influence `Wait`'s` return value in `Strict` mode.

### Modes

The mode is set at construction via `WithMode(m)` and is immutable for the life of the group.

#### Strict mode

The default. Matches `golang.org/x/sync/errgroup` semantics: when any task returns a non-nil error, the group's ctx is canceled, siblings see the cancellation and abort, and `Wait()` returns that error. If multiple tasks fail, only the **first** error is returned; subsequent failures are ignored.

```go
g := concurrency.NewGroup(ctx) // Strict by default
g.Go(migrateUser)              // if this fails...
g.Go(migrateAccount)            // ...this one is canceled
err := g.Wait()                 // err is migrateUser's error (or nil)
```

#### BestEffort mode

All tasks run to completion regardless of siblings' failures. The group's ctx is **not** canceled on a task error. `Wait()` returns:

- `nil` if every task succeeded (including the trivial case of zero tasks spawned).
- The bare single task error if exactly one task failed.
- A [`*BestEffortError`](#bestefforterror) wrapping every failed task otherwise.

Parent-ctx cancellation does not change the shape of the returned error: any task that was still running when the parent ctx was canceled simply returns `ctx.Err()` as its error, and those errors are collected just like any other failure — there is no separate "parent canceled" fast path.

```go
g := concurrency.NewGroup(ctx, concurrency.WithMode(concurrency.BestEffort))
g.Go(refreshA)
g.Go(refreshB)
g.Go(refreshC)
err := g.Wait() // nil only if all three succeeded; *BestEffortError if any failed
```

#### `BestEffortError`

```go
type BestEffortError struct {
    Errors []error // per-task errors in completion order; always non-empty
}

func (e *BestEffortError) Error() string
func (e *BestEffortError) Unwrap() []error // errors.Is / errors.As can walk all underlying errors
```

`BestEffortError` implements the standard `Unwrap() []error` contract, so callers can use `errors.Is(err, targetErr)` to find a specific failure across all collected errors.

### Observability

None — `Group` is meant to be used once and discarded. `Wait()` and `WaitWithContext()` are the only inspection points.

### Memory model

- **Backing concurrency.** `sync.WaitGroup` for "wait for all", `sync.Mutex` to guard the result state, and (if a limit is configured) a [`*Semaphore`](#semaphore) — the toolkit's typed counting semaphore, shared with the public `Semaphore` type rather than a duplicated inline `chan struct{}`. `Wait` / `WaitWithContext` share a single `doneCh` channel (closed when the WaitGroup reaches zero) backed by one lazily-spawned watcher goroutine, regardless of how many times the wait methods are invoked. Both methods also fire the Group's cancel func on return to release the cancelCtx from any parent's `children` map and let any `propagateCancel` watcher exit — without this, `BestEffort` Groups (and `Strict` Groups whose tasks all succeed) would leak one goroutine per instantiation when the parent is a non-cancelCtx custom `Context`. No external dependencies; no `errgroup`, no `x/sync`.
- **Per-op allocations.** One `context.WithCancel` derived ctx at construction. Each `Go` call captures the goroutine closure and adds to the WaitGroup. No other allocations on the hot path.
- **Why not `errgroup`?** `errgroup` is the obvious implementation choice, but its API is awkward to adapt to ours (it doesn't pass ctx to `fn`, its `SetLimit` must be called before any `Go` and panics otherwise, and its `Wait` returns a single error which is fine for `Strict` but awkward for `BestEffort`). Building directly on `sync.WaitGroup` is about the same line count and gives us full control over both modes.

---

## `Semaphore`

A counting semaphore that limits the number of goroutines that can concurrently hold a "slot". It is the primitive behind [`Group`'s `WithLimit`](#group) and is exposed directly for the common "max N concurrent goroutines touching X" use case — e.g., rate-limiting outbound calls to a downstream service, capping the parallelism of a per-key pipeline, or any global cap that doesn't fit cleanly onto a single [`BoundedBlockingQueue[T]`](#boundedblockingqueuet)'s per-queue backpressure.

Unlike [`golang.org/x/sync/semaphore`](https://pkg.go.dev/golang.org/x/sync/semaphore), `Semaphore` uses **fixed unit weights** ("N slots") — there is no `Acquire(ctx, n)` overload with a weight argument. That trade-off matches the common case and keeps the API minimal; if you need weighted acquires, reach for `x/sync/semaphore`.

### Construction

```go
// NewSemaphore(n int) *Semaphore
sem := concurrency.NewSemaphore(8)
```

`n` must be positive; passing `0` or a negative value panics. A misconfigured capacity should fail loudly at construction — a `Semaphore` that can never be acquired is almost always a bug, not a deliberate "always block" choice.

The zero value is not usable; always go through the constructor.

### Blocking variants

| Method | Behaviour |
| --- | --- |
| `Acquire()` | Block until a slot is available, then take one. |
| `AcquireWithContext(ctx context.Context) error` | Ctx-aware `Acquire`. Blocks until a slot is available OR ctx fires; returns `nil` on success, `ctx.Err()` on ctx firing. **On ctx firing no slot is taken** — callers do not need to `Release` to balance. |
| `TryAcquire() bool` | Take a slot without blocking. Returns `true` on success, `false` immediately if no slot is available. |
| `Release()` | Return a previously-acquired slot. Pairs with `Acquire` / successful `TryAcquire` / successful `AcquireWithContext`. |

### Observability

| Method | Behaviour |
| --- | --- |
| `Available() int` | Current number of free slots. Atomic length read of the underlying channel — not synchronised with the next operation you perform. If you need strict "Available() → next op sees a consistent view" semantics, use `TryAcquire` which combines the read with the take. |

### Memory model

- **Backing storage.** A single `make(chan struct{}, n)` pre-filled with `n` values at construction. `Acquire` is a `<-ch`, `Release` is `ch <- struct{}{}`. The channel's natural blocking semantics give us the wait-for-slot behaviour for free, with no separate condition variable.
- **Why pre-fill?** `len(ch)` is the slot count. Pre-filling the buffer once at construction lets `Available()` be a single `len(ch)` atomic load, and lets `TryAcquire` be a `select { case <-ch: default: }` instead of an explicit `len(ch) > 0` check followed by a race-prone receive.
- **Per-op allocations.** Zero. Channel send/recv fast path does not allocate.
- **No over-release.** Calling `Release` without a matching `Acquire` blocks (the channel send can't complete against a full buffer), which in a single-goroutine program surfaces as the Go runtime's deadlock detector and in a concurrent program is silent. This is deliberate — unlike `x/sync/semaphore`, the toolkit does not support silent growth of the slot count beyond the initial `n`.

### Examples

#### Capping concurrent calls to a downstream service

```go
sem := concurrency.NewSemaphore(8) // max 8 in-flight requests

for _, req := range requests {
    sem.Acquire()                  // blocks if 8 are already in flight
    go func() {
        defer sem.Release()
        handle(req)
    }()
}
```

#### Cancellable acquire

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

if err := sem.AcquireWithContext(ctx); err != nil {
    // ctx fired before a slot was available; no slot taken, no Release needed
    return err
}
defer sem.Release()
```

#### Best-effort: try, fall back if full

```go
if !sem.TryAcquire() {
    metrics.Dropped.Inc()
    return // shed load instead of blocking
}
defer sem.Release()
process()
```

## See also

- [`reactivex`](../../reactivex/README.md) — typed async event streams with explicit demand and configurable backpressure. If the work being queued is "events to deliver to many subscribers", an `Observable` is usually a better fit than a queue.
- [`collections.Queue`](../collections/README.md#stackt--queuet--dequet) — the synchronous, in-memory `Queue[T]`. Use it when there is no concurrency and you want `Option[T]`-based access; it has no locking, no `TryPush`, and no backpressure.
- [`chan T`](https://go.dev/ref/spec#Channel_types) in the standard library — `BoundedBlockingQueue` wraps it directly.
