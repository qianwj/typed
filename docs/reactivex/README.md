# `github.com/qianwj/typed/reactivex`

Typed, subscription-explicit, demand-driven, configurable-backpressure asynchronous event streams. This package is unrelated to the `Stream` in [`collections`](../collections/README.md): `Stream` is synchronous, single-consumer; `reactivex` handles subscription lifetimes, asynchronous inputs, and multicast.

Go 1.27+ is required because methods like `Observable.Map[R]` / `Observable.Scan[R]` declare their own type parameters.

> Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Import

```go
import "github.com/qianwj/typed/reactivex"
```

## Core types

### `Observable[T]`

`Observable[T]` is a concrete type, not an interface — that lets transform operators (`Map[R]`, `Scan[R]`, ...) declare their own result type `R`. Method-level type parameters like these do not fit on the `Publisher` interface under Go 1.27.

The zero value of `Observable` has no source and cannot be subscribed to. Always construct via `Just` / `FromSlice` / `FromChannel` / `FromSeq` / `Create` / `Interval`.

```go
var _ reactivex.Publisher[int] = reactivex.Observable[int]{}
```

### `Publisher[T]`

A narrow interface for components that only care about "can I subscribe?". Both `Observable[T]` and `Subject[T]` implement it. Use `Publisher` when the consumer only wants notifications; use `Observable` when building a fluent operator chain.

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
| `Subscribe(ctx, sub) Subscription` | Synchronous setup; callbacks run on the producer goroutine. Returns that subscription's `Subscription`. |
| `ForEach(ctx, onNext, onError, onComplete) Subscription` | One-shot callback; any callback may be `nil`; `ForEach` requests `^uint64` (effectively unbounded) demand, so it receives all available values; **does not wait for termination**. |
| `ToSlice(ctx) ([]T, error)` | Collects into a `[]T`; `OnError` returns `([]T(nil), err)` immediately, otherwise returns the slice after `OnComplete`. |

`ForEach` / `ToSlice` request "unbounded" demand; if you want backpressure, use `Subscribe` and control `Request` yourself.

## Sources

```go
func Just[T any](values ...T) Observable[T]
func FromSlice[T any](values []T) Observable[T]
func FromChannel[T any](ch <-chan T) Observable[T]
func FromChannelWithOptions[T any](ch <-chan T, opts ...BackpressureOption) Observable[T]
func FromSeq[T any](seq iter.Seq[T]) Observable[T]
func Create[T any](run func(ctx context.Context, emit func(T) bool, complete func())) Observable[T]
func Interval(ctx context.Context, period time.Duration) Observable[uint64]
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
| `Map[R]` | `Map[R any](f func(context.Context, T) (R, error)) Observable[R]` | Calls `f` on each `OnNext`. If `f` returns an error, the chain emits `OnError(err)` and completes. `f` receives the `context`, which differs from the slicing-style `collections` API. |
| `Filter` | `Filter(predicate func(T) bool) Observable[T]` | Drops values for which the predicate returns `false`. |
| `Take` | `Take(n uint64) Observable[T]` | Take the first `n` values. |
| `Skip` | `Skip(n uint64) Observable[T]` | Skip the first `n` values. |
| `Scan[R]` | `Scan[R any](initial R, f func(R, T) R) Observable[R]` | Emits the accumulator at each step; the initial value is **not** emitted before the first input. |
| `Reduce` | `Reduce(f func(T, T) T) Observable[T]` | Folds to a single value; an empty stream finishes with `OnComplete` and does not emit a substitute. |

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
src.Subscribe(ctx, reactivex.Subscriber[int]{
    OnSubscribe: func(s reactivex.Subscription) { s.Request(^uint64(0)) },
    OnNext:      func(v int) { consume(v) },
    OnError:     func(err error) { log.Println(err) },
    OnComplete:  func() {},
})
```

```go
// Hot multicast
subj := reactivex.NewSubject[string]()
go func() {
    defer subj.OnComplete()
    for _, v := range source {
        subj.OnNext(v)
    }
}()
subj.ForEach(ctx, onMsg, onErr, onDone) // start one subscription
```

## See also

- For synchronous, single-consumer iteration use the `Stream` in [`collections`](../collections/README.md).
- For one-shot success / failure use [`utils/result`](../result/README.md); subscription-level errors are reported through `OnError`.
- For "single value that may be absent" use [`utils/option`](../option/README.md).
