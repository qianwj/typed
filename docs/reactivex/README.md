# `typed/reactivex` — Typed

[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=reactivex)](https://codecov.io/gh/qianwj/typed)

Typed, subscription-explicit, demand-driven, configurable-backpressure asynchronous event streams. This package is unrelated to the `Stream` in [`collections`](../collections/README.md): `Stream` is synchronous, single-consumer; `reactivex` handles subscription lifetimes, asynchronous inputs, and multicast.

Go 1.27+ is required because methods like `Flowable.Map[R]` / `Flowable.Scan[R]` declare their own type parameters.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [Core types](#core-types)
  - [`Flowable[T]`](#flowablet)
  - [`Publisher[T]`](#publishert)
  - [`Subscriber[T]`](#subscribert)
  - [`Subscription`](#subscription)
- [Subscribing and consuming](#subscribing-and-consuming)
- [Sources](#sources)
- [Operators](#operators)
- [Backpressure](#backpressure)
  - [`OverflowStrategy`](#overflowstrategy)
  - [`BackpressureOption`](#backpressureoption)
- [`Subject[T]`](#subjectt)
- [`Single[T]`](#singlet)
- [`Maybe[T]`](#maybet)
- [Type conversions](#type-conversions)
- [Examples](#examples)
- [See also](#see-also)

## Import

```go
import "github.com/qianwj/typed/reactivex"
```

Blocking consumption of `Single` and `Maybe` returns types from `github.com/qianwj/typed/adt`; this module depends on `adt`.

## Core types

### `Flowable[T]`

`Flowable[T]` is a concrete type, not an interface — that lets transform operators (`Map[R]`, `Scan[R]`, ...) declare their own result type `R`. Method-level type parameters like these do not fit on the `Publisher` interface under Go 1.27.

`Flowable`, `Single`, and `Maybe` all use value receivers. A Flowable is an immutable pipeline description: copies share its source function and captured resources, while each subscription creates fresh operator state. Single/Maybe copies instead share one execution and its cached terminal result. Copying any of these handles does not start consumption.

The zero value of `Flowable` has no source and cannot be subscribed to. Always construct via `Just` / `FromSlice` / `FromChannel` / `FromSeq` / `Create` / `Interval`.

```go
var _ reactivex.Publisher[int] = reactivex.Flowable[int]{}
```

### `Publisher[T]`

A narrow interface for components that only care about "can I subscribe?". Both `Flowable[T]` and `Subject[T]` implement it. Use `Publisher` when the consumer only wants notifications; use `Flowable` when building a fluent operator chain.

```go
type Publisher[T any] interface {
    Subscribe(ctx context.Context, sub Subscriber[T]) Subscription
}
```

### `Subscriber[T]`

```go
type Subscriber[T any] interface {
    OnSubscribe(Subscription)
    OnNext(T)
    OnError(error)
    OnComplete()
}
```

- `OnSubscribe` is called with the subscription handle before any value is delivered. The subscriber may call `Request` or `Cancel` from inside it. Retain the handle if you will replenish demand as processing finishes.
- `OnError` and `OnComplete` are terminal signals — they do **not** consume demand and must not be reinterpreted as a data item to request again.
- Callbacks may run on a producer goroutine or directly on a caller's goroutine; subscribers must not depend on a particular execution thread. State shared by multiple subscriptions needs its own synchronisation.
- **Callback panics are not converted into `OnError` notifications.**

### `Subscription`

```go
type Subscription interface {
    Request(n uint64)   // accumulate demand; n==0 is a no-op; saturates at uint64 max
    Cancel()            // idempotent; does not wait for an in-flight callback; does not interrupt user code
    Done() <-chan struct{} // closes on cancellation or termination
}
```

`Done` is a **lifecycle signal**, not a join on every piece of work started by a producer or callback.

## Subscribing and consuming

| Method | Use |
|---|---|
| `Subscribe(ctx, sub) Subscription` | Sets up the subscription and supplies its handle through `OnSubscribe`. Callback execution depends on the source; no demand is requested automatically. |
| `ForEach(ctx, onNext, onError, onComplete) Subscription` | Registers callbacks and requests `^uint64(0)` demand. Any callback may be `nil`; **does not wait for termination**. |
| `ToSlice(ctx) ([]T, error)` | Blocks until completion, error or context cancellation. Returns collected values together with the error, including partial values on failure or cancellation. |

`ForEach` / `ToSlice` request maximum demand; use `Subscribe` to control how many values may be delivered. `ToSlice` needs a finite source and returns a non-nil empty slice on normal empty completion.

Package-level `Collect` retains its blocking signature:

```go
func Collect[T, R any](ctx context.Context, source Flowable[T], initial R, f func(R, T) R) (R, error)
```

It first calls `ToSlice`, then folds the collected values into `initial`. On failure it returns `initial` and the error without calling `f`; it does not fold partial values. For an incremental aggregate that can become a Single, use `Reduce(...).FirstOrError(ctx)` (see [Type conversions](#type-conversions)).

## Sources

```go
func Just[T any](values ...T) Flowable[T]
func FromSlice[T any](values []T) Flowable[T]
func FromChannel[T any](ch <-chan T) Flowable[T]
func FromChannelWithOptions[T any](ch <-chan T, opts ...BackpressureOption) Flowable[T]
func FromSeq[T any](seq iter.Seq[T]) Flowable[T]
func Create[T any](run func(ctx context.Context, emit func(T) bool, complete func())) Flowable[T]
func Interval(ctx context.Context, period time.Duration) Flowable[uint64]
```

| Source | Key properties |
|---|---|
| `Just(vs...)` | Delegates to `FromSlice(vs...)`; every subscription starts at the first value. |
| `FromSlice(vs)` | One producer goroutine per subscription. **Shares the backing array of `vs`** (no copy). Mutating the backing array concurrently with an active subscription is the caller's responsibility. A slow callback slows the source. |
| `FromChannel(ch)` | Equivalent to `FromChannelWithOptions(ch)`: default blocking, no queueing. Multiple subscriptions **compete** for values from the same `ch` (use a `Subject` for broadcast). The source does not own or close `ch`. |
| `FromChannelWithOptions(ch, opts...)` | Each subscription gets its own per-subscription queue; subscriptions still compete for the same `ch`. Closing `ch` schedules normal completion after queued values drain; draining requires demand. |
| `FromSeq(seq)` | `seq` is invoked once per subscription. Reentrancy / concurrent validity is the caller's responsibility. Returning `false` from the emit callback stops iteration. |
| `Create(run)` | General-purpose source. `emit` waits for demand; `complete` runs at most once. |
| `Interval(ctx, period)` | One ticker per subscription, emitting counters starting at zero. The first value follows the first tick, not subscription setup. `period <= 0` makes `time.NewTicker` panic. The constructor's `ctx` argument is currently unused; the `ctx` supplied to `Subscribe` / `ForEach` controls the loop. |

> There is no `Empty` / `Error` / `Never` / `Merge` / `Concat` / `Debounce` / `Throttle` / `Sample` source or operator in this package.

## Operators

| Operator | Signature | Semantics |
|---|---|---|
| `Map[R]` | `Map[R any](f func(context.Context, T) (R, error)) Flowable[R]` | Calls `f` on each `OnNext`. If `f` returns an error, the chain terminates with `OnError(err)` and cancels upstream. `f` receives the `context`, which differs from the slicing-style `collections` API. |
| `Filter` | `Filter(predicate func(T) bool) Flowable[T]` | Drops values for which the predicate returns `false`. |
| `Take` | `Take(n uint64) Flowable[T]` | Take the first `n` values. |
| `Skip` | `Skip(n uint64) Flowable[T]` | Skip the first `n` values. |
| `Scan[R]` | `Scan[R any](initial R, f func(R, T) R) Flowable[R]` | Emits the accumulator at each step; the initial value is **not** emitted before the first input. |
| `Reduce` | `Reduce(f func(T, T) T) Flowable[T]` | Folds a finite source to one value. The first positive request requests all upstream inputs. Empty input completes without a value. |

`Map` / `Filter` / `Take` / `Skip` / `Scan` / `Reduce` are all **wrapping** operators: they layer over the downstream `Subscriber` and do **not** start their own goroutine or maintain their own queue.

## Backpressure

Demand and buffering are **two different things**:

- `Subscription.Request(n)` — unlock `n` values, letting the producer continue sending.
- `WithBuffer` / `WithOverflow` on a `Subject` or channel source — cap the number of "arrived but not yet delivered" values pending at the asynchronous boundary.

### `OverflowStrategy`

```go
type OverflowStrategy uint8

const (
    OverflowBlock       OverflowStrategy = iota // wait for demand or buffer room
    OverflowDropLatest                            // drop the incoming value, keep the older ones
    OverflowDropOldest                            // drop the oldest, keep the incoming one; requires WithBuffer > 0
    OverflowKeepLatest                            // keep exactly one pending value, replaced by newer ones; WithBuffer is ignored
    OverflowError                                 // terminate that subscription and emit ErrBackpressureOverflow
)
```

```go
var ErrBackpressureOverflow = errors.New("reactivex: backpressure buffer overflow")
```

`OverflowError` affects **only that one subscription**; the `Subject` and sibling subscriptions continue.

### `BackpressureOption`

```go
type BackpressureOption func(*backpressureConfig)

func WithBuffer(size int) BackpressureOption   // default 0; negative panics
func WithOverflow(strategy OverflowStrategy) BackpressureOption
```

- Options apply in order; later options override earlier ones; `nil` is skipped.
- `OverflowKeepLatest` always uses one pending slot, regardless of `WithBuffer`.
- `OverflowDropOldest` combined with `WithBuffer(0)` panics at construction.

```go
subject := reactivex.NewSubject[int](
    reactivex.WithBuffer(128),
    reactivex.WithOverflow(reactivex.OverflowDropOldest),
)
```

## `Subject[T]`

A hot multicast publisher and subscriber. Each subscription has its own demand and a buffer configured at construction time. **No replay** — `OnNext` with no subscribers discards the value, and a new subscriber only participates in later publications. The terminal state is retained, so late subscribers receive completion or the stored error directly.

The default unbuffered blocking mode invokes callbacks synchronously when demand is available. A buffer or a non-blocking overflow strategy dispatches notifications asynchronously, but callbacks are still serialised per subscription. Across different subscriptions callbacks may run concurrently. Construct a `Subject` with `NewSubject` and do not copy it after use.

```go
func NewSubject[T any](options ...BackpressureOption) *Subject[T]
func (s *Subject[T]) Subscribe(ctx context.Context, out Subscriber[T]) Subscription
func (s *Subject[T]) ForEach(ctx context.Context, onNext func(T), onError func(error), onComplete func()) Subscription
```

`Subject` also implements `Subscriber[T]` and can be used as a bridge from another stream:

```go
func (s *Subject[T]) OnSubscribe(sub Subscription)
func (s *Subject[T]) OnNext(v T)
func (s *Subject[T]) OnError(err error)
func (s *Subject[T]) OnComplete()
```

> `OnSubscribe` currently does **not** automatically request from upstream, and downstream demand is not aggregated into upstream requests. When using `Subject` as a `Subscriber`, call `Request` on the returned upstream subscription explicitly. Only the most recently supplied upstream handle is retained.

## `Single[T]`

`Single` uses value receivers; constructors and composition operators return values. Copies share the same execution and cached result through internal state. Copying does not start the source. The zero value is unusable; construct a handle with `NewSingle` or a Flowable conversion.

`Single[T]` is a reactive container that emits **exactly one** value or **exactly one** error. It is the typed equivalent of a Future combined with a reactive subscription model.

Compared to `Flowable[T]`:

- **Cardinality is fixed at 1.** A `Single` terminates with `OnSuccess(T)` or `OnError(error)`; there is no "no value arrived" state, so callers never have to distinguish "the result is still pending" from "no result will ever arrive".
- **No demand tracking.** At most one value is delivered, so the subscriber does not call `Request`.
- **Shared execution.** `fn` runs once and the result is cached; subscribers arriving after completion also receive that outcome. Use a cold Flowable source such as `Just` or `FromSlice` for a fresh iteration per subscriber; channel sources share input and `ToFlowable` preserves its Single/Maybe cache.

```go
type Single[T any] struct { /* ... */ }

func NewSingle[T any](fn func() adt.Result[T]) Single[T]

// Reactive subscription
func (s Single[T]) Subscribe(onSuccess func(T), onError func(error)) Subscription

// Blocking consumption
func (s Single[T]) Await() adt.Result[T]
func (s Single[T]) AwaitWithContext(ctx context.Context) adt.Result[T]

// Non-blocking check
func (s Single[T]) Done() bool

// Composition
func (s Single[T]) Map[R any](f func(T) R) Single[R]
func (s Single[T]) FlatMap[R any](f func(T) Single[R]) Single[R]
func (s Single[T]) Zip[U, R any](other Single[U], combine func(T, U) R) Single[R]
func (s Single[T]) AndThen[R any](next Single[R]) Single[R]
```

The source and `Await` both return `Result[T]`: `Success(value)` or `Failure[T](err)`. For an existing `(T, error)` function, return `adt.Wrap(loadConfig())` from the source closure. To use Go's `(T, error)` form, call `single.Await().Unwrap()`; otherwise compose directly with Result methods such as `Map` and `OrElse`.

`AwaitWithContext` returns `Failure(ctx.Err())` when the wait is canceled. An already completed result takes precedence over cancellation; otherwise an already canceled context prevents starting the source. Canceling a wait does not cancel a running source or overwrite its eventual result.

`Subscribe` returns without waiting for user code, including when the result is cached. Its `Subscription.Done()` closes after callback delivery or immediately on cancellation. `Cancel` suppresses a callback whose delivery has not started, without stopping the shared source or interrupting an in-flight callback. `Single.Done()` and `Await` observe the terminal result independently of subscriber callbacks.

### Composition

Composition is lazy and connects completion callbacks: operators do not call `Await` internally or park a goroutine for each stage. Starting consumption uses an asynchronous entry point, and each source runs in its own goroutine at most once. Transforms and pending subscription callbacks execute outside internal locks on the completing goroutine; cached subscriptions are delivered asynchronously. Slow transforms or callbacks delay other continuations on the same goroutine, so they should not synchronously wait for a dependent continuation to run.

| Operator | Behaviour |
| --- | --- |
| `Map` | Apply `f(value)` on success; pass errors through unchanged. `f` must be infallible. |
| `FlatMap` | On success, run `f(value)` and return the resulting `Single`. The chain produces a fresh `Single`; the inner `Single` runs only if the outer succeeded. |
| `Zip` | Observe this `Single` first, then `other` on success, and run `combine(left, right)`. A left error short-circuits the right source. |
| `AndThen` | On success, run `next` and return its result. The original value is discarded; use `FlatMap` if `next` depends on it. |

### Comparison with `Flowable`

| | `Flowable[T]` | `Single[T]` |
| --- | --- | --- |
| Cardinality | 0..N | exactly 1 |
| Terminal states | `OnNext*` then `OnComplete` or `OnError` | `OnSuccess(T)` or `OnError(error)` |
| Demand tracking | yes (`Request`) | no |
| Repeated subscription | Source-dependent: fresh slice iteration, competing channel reads, or shared converted result | Same cached result |
| Blocking consumption | `ToSlice`, package-level `Collect` | `Await`, `AwaitWithContext` |

### Use cases

- **Cache warm-up**: `var cfg = reactivex.NewSingle(func() adt.Result[Config] { return adt.Wrap(loadConfig()) })` — multiple callers `Await` the same `Single`, the underlying fn runs once.
- **One-shot async computation** with a typed return (HTTP fetch on demand, expensive calculation).
- **Bridge from a callback-based API** to a typed value.

## `Maybe[T]`

`Maybe` uses value receivers; constructors and composition operators return values. Copies share the same execution and cached result through internal state. Copying does not start the source. The zero value is unusable; construct a handle with `NewMaybe` or a Flowable conversion.

`Maybe[T]` generalises `Single` by adding a third terminal state: the producer may legitimately complete **without a value**. It models "looked up the key, no entry" or "scanned the queue, no message pending" — cases where absence is a normal outcome rather than an error.

Terminal states:

- `OnSuccess(T)` — exactly one value delivered.
- `OnComplete()` — no value delivered; the absence is normal.
- `OnError(err)` — abnormal termination.

```go
type Maybe[T any] struct { /* ... */ }

func NewMaybe[T any](fn func() adt.Result[adt.Option[T]]) Maybe[T]
// Success(Of(value)) → OnSuccess; Success(Empty[T]()) → OnComplete.
// Failure[Option[T]](err) → OnError.

func (m Maybe[T]) Subscribe(
    onSuccess func(T),
    onComplete func(),
    onError func(error),
) Subscription

func (m Maybe[T]) Await() adt.Result[adt.Option[T]]
func (m Maybe[T]) AwaitWithContext(ctx context.Context) adt.Result[adt.Option[T]]

func (m Maybe[T]) Done() bool

func (m Maybe[T]) Map[R any](f func(T) R) Maybe[R]
func (m Maybe[T]) FlatMap[R any](f func(T) Maybe[R]) Maybe[R]
func (m Maybe[T]) Zip[U, R any](other Maybe[U], combine func(T, U) R) Maybe[R]
func (m Maybe[T]) AndThen[R any](next Maybe[R]) Maybe[R]
```

The outer Result describes success or failure; the inner Option describes whether a successful completion emitted a value:

| Outcome | Await result |
| --- | --- |
| `OnSuccess(value)` | `Success(Of(value))` |
| `OnComplete()` | `Success(Empty[T]())` |
| `OnError(err)` | `Failure[Option[T]](err)` |
| Wait canceled or deadline exceeded | `Failure[Option[T]](ctx.Err())` |

Presence follows the inner Option: `Success(Of(value))` preserves `0` and nil as present values. Context handling and cached-result precedence follow Single's rules above.

```go
value, err := reactivex.NewMaybe(func() adt.Result[adt.Option[string]] {
    return adt.Success(adt.Empty[string]())
}).Await().Unwrap()
if err != nil {
    return err
}
name := value.OrElse("anonymous") // empty completion uses the fallback
```

### Composition

All operators propagate the three states: `Map` / `FlatMap` / `Zip` / `AndThen` pass `OnComplete` through unchanged (the user's `f` is not invoked on the complete-without-value branch). Errors always short-circuit the chain.

Execution and subscription cancellation follow Single's rules above. Maybe also composes through completion callbacks, preserving empty completion without waiting goroutines. `Zip` observes the left side first; a left error or empty completion skips the right source.

### Comparison with `Single` and `Flowable`

| | `Flowable[T]` | `Single[T]` | `Maybe[T]` |
| --- | --- | --- | --- |
| Cardinality | 0..N | exactly 1 | 0 or 1 |
| Terminal states | complete / error | success / error | success / complete / error |
| Await return | n/a (use `ToSlice`) | `Result[T]` | `Result[Option[T]]` |
| Best for | streams of events | a guaranteed result | an optional result |

### Use cases

- **Cache lookup** — hit returns `Success(Of(value))`, miss returns `Success(Empty[T]())`, failure returns `Failure[Option[T]](err)`.
- **Pull from a queue** with timeout — got a message or got nothing.
- **Database row fetch** by primary key — row exists or doesn't.

## Type conversions

```go
func (s Single[T]) ToFlowable() Flowable[T]
func (m Maybe[T]) ToFlowable() Flowable[T]
func (o Flowable[T]) FirstElement(ctx context.Context) Maybe[T]
func (o Flowable[T]) FirstOrError(ctx context.Context) Single[T]
```

| Method | Result | Semantics |
| --- | --- | --- |
| `single.ToFlowable()` | `Flowable[T]` | Emit the Single's value, then complete; propagate failure. |
| `maybe.ToFlowable()` | `Flowable[T]` | Emit a present value or complete empty; propagate failure. |
| `flow.FirstElement(ctx)` | `Maybe[T]` | Select the first value; an empty stream completes normally. |
| `flow.FirstOrError(ctx)` | `Single[T]` | Select the first value; an empty stream fails with `ErrNoElements`. |

All conversions are lazy. `ToFlowable` subscriptions share the original cached computation, but have independent demand and cancellation. Each subscription can hold one pending value until `Request` supplies demand. Empty completion and errors need no demand; present zero and nil values remain values. Canceling a converted stream subscription does not stop its Single or Maybe source.

`FirstElement` and `FirstOrError` subscribe once when the returned container is consumed and cache its result. They request one output and cancel upstream on the first value, without waiting for stream completion or checking for further elements. Errors before that value propagate unchanged. The supplied `ctx` controls the shared upstream subscription; cancellation settles the container with `Failure(ctx.Err())`. Canceling an individual `AwaitWithContext` wait only stops that wait. A nil source context means `context.Background()`.

```go
first := reactivex.Just(3, 5).FirstElement(ctx) // Maybe[int]
value, err := first.Await().Unwrap()          // Option[int], error

required := reactivex.Just[int]().FirstOrError(ctx)
missing := errors.Is(required.Await().Error(), reactivex.ErrNoElements)

total := reactivex.Just(1, 2, 3).
    Reduce(func(sum, value int) int { return sum + value }).
    FirstOrError(ctx) // Single[int], successful value 6
```

`Reduce` translates a request for its one result into demand for all upstream inputs, so it can feed these first-element conversions. It still needs a finite source to complete. Existing `ToSlice(ctx)` and package-level `Collect` retain their blocking return types.

## Examples

```go
// A cold, finite stream
values, err := reactivex.Just(1, 2, 3, 4).
    Filter(func(v int) bool { return v%2 == 0 }).
    ToSlice(ctx)
// err == nil ⇒ values == []int{2, 4}
```

```go
// Asynchronous channel source, buffered with drop-oldest
src := reactivex.FromChannelWithOptions(ch,
    reactivex.WithBuffer(64),
    reactivex.WithOverflow(reactivex.OverflowDropOldest),
)
sub := src.ForEach(ctx,
    func(v int) { consume(v) },
    func(err error) { log.Println(err) },
    nil,
)
<-sub.Done()
```

```go
// Hot multicast
subj := reactivex.NewSubject[string]()
sub := subj.ForEach(ctx, onMsg, onErr, onDone)
go func() {
    defer subj.OnComplete()
    for _, v := range source {
        subj.OnNext(v)
    }
}()
<-sub.Done() // The subscriber was attached before publication started.
```

## See also

- For synchronous, single-consumer iteration use the `Stream` in [`collections`](../collections/README.md).
- [`adt`](../adt/README.md) provides `Result[T]` and `Option[T]`. Single producers and Await return `Result[T]`; Maybe uses `Result[Option[T]]` to keep empty completion distinct from failure.
