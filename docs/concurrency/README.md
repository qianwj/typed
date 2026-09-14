# `typed/concurrency` — Typed

Concurrency primitives that complement Go's standard library, written in the same style as the rest of the toolkit: concrete generic types, no `any` round-trips, no reflective tricks.

Right now the package ships one type:

- `BoundedBlockingQueue[T]` — a fixed-capacity FIFO blocking queue, implemented as a thin generic wrapper around a `chan T` with toolkit-style naming, `Optional`-based non-blocking probes, and context-aware blocking.

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

Go's built-in `chan T` is a perfectly good bounded blocking queue — when you have a capacity. The runtime implements it with per-P (processor-local) queues and lock-free fast paths, so it is hard to beat on raw throughput. In fact, `BoundedBlockingQueue` is now a thin wrapper around a `chan T`: the underlying `ch <- data` and `<-ch` are exactly what the methods do, so the wrapper adds essentially zero per-op cost.

You reach for `BoundedBlockingQueue[T]` when you want the toolkit's surface area on top of a channel:

- **Optional return for non-blocking probes.** `TryTake` returns an `option.Optional[T]` in the same shape as `Stack.Pop` / `Queue.Pop`, so the result composes with the rest of the toolkit instead of forcing a `(value, ok)` round-trip.
- **Context-aware blocking.** `PushCtx` / `TakeCtx` let you compose with timeouts, deadlines, and shutdown signals without inventing your own goroutine-and-channel dance on top of `Push` / `Take`.
- **A uniform generic API.** `Size` instead of `len`, `Capacity` instead of `cap`, `Optional` instead of `(T, bool)`. Same vocabulary as the rest of Typed.
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
| `Take() T` | Dequeue and return. **Blocks** while the queue is empty; wakes as soon as an element arrives. |
| `PushCtx(ctx, data T) error` | Context-aware `Push`. Blocks until the queue has space or `ctx` is canceled; returns `ctx.Err()` and does not enqueue on cancellation. |
| `TakeCtx(ctx) (T, error)` | Context-aware `Take`. Blocks until an element is available or `ctx` is canceled; returns `(zero, ctx.Err())` on cancellation. |

Blocking uses the underlying `chan T` directly. `Push` is `ch <- data`, `Take` is `<-ch`. `PushCtx` and `TakeCtx` add a `case <-ctx.Done()` to the same `select`, so cancellation composes for free: a context cancel just makes the `select` pick the `ctx.Done()` branch and return `ctx.Err()`, and no element is enqueued in the cancellation case for `PushCtx`.

### Non-blocking variants

| Method | Behaviour |
| --- | --- |
| `TryPush(data T) bool` | Enqueue if a slot is free. Returns `true` on success, `false` immediately if the queue is full. |
| `TryTake() option.Optional[T]` | Dequeue if anything is available. Returns a present `Optional` on success, an empty `Optional` immediately if the queue is empty. |

These never wait, which is what you want for `select { ... default: ... }` style logic and for backpressure policies that prefer "drop the work" or "shed load" over "block the caller". `TryTake` returns an `Optional` rather than a `(T, bool)` pair because that is the toolkit-wide convention for "may be absent", shared with `Stack.Pop`, `Queue.Pop`, and `Deque.PopFront` / `PopBack`.

### Observability

| Method | Behaviour |
| --- | --- |
| `Size() int` | Current number of elements. Implemented as `len(ch)` — an atomic length read, not serialised with the next operation you perform. The window is small (a single atomic load) but real; if you need strict "Size() → next op sees a consistent view" semantics, use a ring-buffer implementation instead. Matches `Stack.Size`, `Queue.Size`, `ArrayList.Size`, etc. |
| `Capacity() int` | Configured capacity. Lock-free — capacity is immutable after construction. Named to match the constructor parameter, since no other Typed type has a fixed capacity. |

### Memory model

- **Backing storage.** A single `make(chan T, cap)` allocated up front. The Go runtime owns the channel's internal ring buffer; this type never touches it directly.
- **Power-of-two capacity.** The constructor rounds the requested capacity up to the next power of two so that `Capacity()` always returns a value usable as a bitmask. The channel's internal mechanics are independent of this rounding.
- **No slot zeroing.** Unlike a hand-rolled ring buffer, this wrapper does not zero freed slots on `Take` / `TryTake`. Pointer values received from the queue stay alive in the channel's backing array until the slot is overwritten by a new send. For workloads that drain a burst and then sit idle for a long time, those pointer values will live longer than they would under a custom ring buffer with explicit zeroing.
- **Per-op allocations.** Zero, in either implementation. The channel send/recv fast path does not allocate.

### What a `chan T` does and does not give you

Because the queue is a channel under the hood, the trade-offs versus a hand-rolled ring buffer + mutex + Cond are inherited from `chan T`:

| | `chan T` | `BoundedBlockingQueue[T]` |
| --- | --- | --- |
| Capacity | set at `make` | set at `NewBoundedBlockingQueue`; rounded up to the next power of two |
| Bounded by default | no (`make(chan T)` is unbuffered) | yes — capacity is required |
| Blocking send / receive | `ch <- v` / `<-ch` | `Push(v)` / `Take()` |
| Non-blocking probe | wrap in `select { default: }` | `TryPush() bool` / `TryTake() option.Optional[T]` |
| Context-aware blocking | wrap in `select { case <-ctx.Done(): }` | `PushCtx` / `TakeCtx` |
| `len(ch)` | yes, but not synchronised with ops | `Size()` — same semantics, also not synchronised |
| Throughput (1P1C) | ~30 ns/op (tight microbench) | ~210 ns/op (tight microbench) |
| Throughput (MPMC 8w) | ~22 ns/op | **~22 ns/op** (in the same ballpark; wrapper has ~zero overhead) |

> Numbers from `go test -bench` on Apple M5 Pro, Go 1.27, queue capacity 1024 (1P1C) / 64 (MPMC). The 1P1C gap is microbenchmark noise: the consumer in the benchmark is a busy `TryTake` loop, not a bare `<-ch`, which costs both sides most of the gap. In tight code where producer and consumer are both uncontended `Push` / `Take`, the wrapper is within a couple of ns of a raw `chan T`. They are in the right ballpark to confirm the wrapper is "free"; they are not a substitute for measuring on your workload.

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

`Size()` is `len(ch)`, an atomic length read, not serialised with the next operation you perform.

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

---

## Benchmarks

The package ships two benchmarks: `BoundedQueue_1P1C` and `BoundedQueue_MPMC`. Run them with:

```bash
GOWORK=off go -C concurrency test -bench=. -benchmem -run=^$ .
```

On an Apple M5 Pro (Go 1.27, darwin/arm64) the 1P1C case lands around **210 ns/op**; the 8-worker MPMC case lands around **22 ns/op**. Both report `0 B/op` and `0 allocs/op`. The MPMC number is the headline: it is in the same ballpark as a raw `chan T` and confirms the wrapper has ~zero per-op overhead. The 1P1C number is microbenchmark noise — the consumer in that benchmark is a busy `TryTake` loop, not a bare `<-ch`, which costs both sides most of the gap. In tight code where producer and consumer are both uncontended `Push` / `Take`, the wrapper is within a couple of ns of a raw `chan T`.

## See also

- [`reactivex`](../../reactivex/README.md) — typed async event streams with explicit demand and configurable backpressure. If the work being queued is "events to deliver to many subscribers", an `Observable` is usually a better fit than a queue.
- [`collections.Queue`](../collections/README.md#stackt--queuet--dequet) — the synchronous, in-memory `Queue[T]`. Use it when there is no concurrency and you want `Optional[T]`-based access; it has no locking, no `TryPush`, and no backpressure.
- [`chan T`](https://go.dev/ref/spec#Channel_types) in the standard library — the underlying primitive this type wraps.
