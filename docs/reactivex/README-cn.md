# `typed/reactivex` — Typed

强类型、显式订阅、按需推送、可配置背压的异步事件流。本包与 [`collections`](../collections/README-cn.md) 的 `Stream` 无关：`Stream` 是同步单消费，`reactivex` 处理订阅生命周期、异步输入与多播。

需要 Go 1.27+（`Observable.Map[R]` / `Observable.Scan[R]` 等方法自带类型参数）。

> **Typed** 工具集的一部分。Looking for the English version? See [README.md](./README.md)。

## 目录

- [包导入](#包导入)
- [核心类型](#核心类型)
  - [`Observable[T]`](#observablet)
  - [`Publisher[T]`](#publishert)
  - [`Subscriber[T]`](#subscribert)
  - [`Subscription`](#subscription)
- [订阅与消费](#订阅与消费)
- [源](#源)
- [算子](#算子)
- [背压](#背压)
  - [`OverflowStrategy`](#overflowstrategy)
  - [`BackpressureOption`](#backpressureoption)
- [`Subject[T]`](#subjectt)
- [例子](#例子)
- [与其他包的关系](#与其他包的关系)

## 包导入

```go
import "github.com/qianwj/typed/reactivex"
```

## 核心类型

### `Observable[T]`

`Observable[T]` 是具体类型而不是接口 —— 这样转换算子（`Map[R]`、`Scan[R]` 等）可以在自己的方法上声明结果类型 `R`，这些方法级类型参数在 Go 1.27 之前无法放在 `Publisher` 接口上。

`Observable` 的零值没有源、不能订阅。构造必须用 `Just` / `FromSlice` / `FromChannel` / `FromSeq` / `Create` / `Interval` 等。

```go
var _ reactivex.Publisher[int] = reactivex.Observable[int]{}
```

### `Publisher[T]`

只关心"能否订阅"的窄接口。`Observable[T]` 和 `Subject[T]` 都实现了它。组件只消费通知时用 `Publisher`；构建流式算子链用 `Observable`。

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

- `OnSubscribe` 必须在值传递之前拿到 handle；可以同步调 `Request` / `Cancel`。
- `OnError` / `OnComplete` 不消耗 demand —— 它们是终止信号，不是数据项。
- 回调可能在生产者 goroutine 上被调用，**不要**依赖具体线程；多个订阅共享状态时自己加锁。
- **回调 panic 不会被转成 `OnError`**。

### `Subscription`

```go
type Subscription interface {
    Request(n uint64)   // 累加需求；n==0 不动；上界 uint64 最大
    Cancel()            // 幂等；不等回调返回，也不强制打断用户代码
    Done() <-chan struct{} // 取消或终止时关闭
}
```

`Done` 是**生命周期信号**，不是"所有派生工作已完成"的 join。

## 订阅与消费

| 方法 | 用途 |
|---|---|
| `Subscribe(ctx, sub) Subscription` | 同步设置；回调在生产者 goroutine 上发生。返回该订阅的 `Subscription`。 |
| `ForEach(ctx, onNext, onError, onComplete) Subscription` | 一次性回调；任意回调可为 `nil`；`ForEach` 内部申请 `^uint64` 满需求，所以会拿到所有可用值；**不等终止**。 |
| `ToSlice(ctx) ([]T, error)` | 收集到 `[]T`；`OnError` 立即返回 `([]T(nil), err)`，否则在 `OnComplete` 后返回整片。 |

`ForEach` / `ToSlice` / `ToSlice` 会申请"无限"需求；想限流请用 `Subscribe` 自己控制 `Request`。

## 源

```go
func Just[T any](values ...T) Observable[T]
func FromSlice[T any](values []T) Observable[T]
func FromChannel[T any](ch <-chan T) Observable[T]
func FromChannelWithOptions[T any](ch <-chan T, opts ...BackpressureOption) Observable[T]
func FromSeq[T any](seq iter.Seq[T]) Observable[T]
func Create[T any](run func(ctx context.Context, emit func(T) bool, complete func())) Observable[T]
func Interval(ctx context.Context, period time.Duration) Observable[uint64]
```

| 源 | 关键性质 |
|---|---|
| `Just(vs...)` | 委托给 `FromSlice(vs...)`；每次订阅从第一个值开始。 |
| `FromSlice(vs)` | 每次订阅一个生产者 goroutine；**共享 `vs` 的底层数组**（不复制），并发修改需调用方保证。慢回调会拖慢源。 |
| `FromChannel(ch)` | 等价 `FromChannelWithOptions(ch)`：默认阻塞、无限缓冲 = 1。多个订阅**竞争**消费同一条 `ch`（要广播请用 `Subject`）。不负责关闭 `ch`。 |
| `FromChannelWithOptions(ch, opts...)` | 每个订阅独立队列；同一 `ch` 仍然被多个订阅竞争。关闭 `ch` 会让该订阅在排空后正常完成。 |
| `FromSeq(seq)` | 每次订阅调一次 `seq`；`seq` 的可重入性由调用方负责。`emit` 返回 `false` 停止迭代。 |
| `Create(run)` | 通用源；`emit` 等待 demand，`complete` 至多一次。 |
| `Interval(ctx, period)` | 每订阅一个 ticker，从 0 开始发计数器；**第一个值在第一个 tick 之后**；`period <= 0` 时 `time.NewTicker` panic。`ctx` 参数当前未使用，订阅时的 `ctx` 才控制循环。 |

> 不存在 `Empty` / `Error` / `Never` / `Merge` / `Concat` / `Debounce` / `Throttle` / `Sample` 源或算子。

## 算子

| 算子 | 签名 | 语义 |
|---|---|---|
| `Map[R]` | `Map[R any](f func(context.Context, T) (R, error)) Observable[R]` | 每次 `OnNext` 调 `f`；`f` 返回 error → 整条链发 `OnError(err)` 并完成。`f` 必须能感知 `context`，与切片风格的 `collections` 不同。 |
| `Filter` | `Filter(predicate func(T) bool) Observable[T]` | 谓词为 `false` 时跳过。 |
| `Take` | `Take(n uint64) Observable[T]` | 取前 `n` 个。 |
| `Skip` | `Skip(n uint64) Observable[T]` | 跳过头 `n` 个。 |
| `Scan[R]` | `Scan[R any](initial R, f func(R, T) R) Observable[R]` | 每一步发累加器；初始值在收到第一个值之前**不**发。 |
| `Reduce` | `Reduce(f func(T, T) T) Observable[T]` | 折成单个值；空流 `OnComplete` 而不补发。 |

`Map` / `Filter` / `Take` / `Skip` / `Scan` / `Reduce` 都是**包装型**算子 —— 它们用下游 `Subscriber` 包一层，**不**自己开 goroutine、**不**自带队列。

## 背压

需求和缓冲区是**两件不同的事**：

- `Subscription.Request(n)` —— 解锁 `n` 个值，让生产者可以继续送。
- `Subject` / 通道源上的 `WithBuffer` / `WithOverflow` —— 限制在异步边界上"已到但未投递"的最大挂起值数。

### `OverflowStrategy`

```go
type OverflowStrategy uint8

const (
    OverflowBlock       OverflowStrategy = iota // 等需求或缓冲位
    OverflowDropLatest                            // 丢新值，保留旧值
    OverflowDropOldest                            // 丢最旧，留新值；需要 WithBuffer > 0
    OverflowKeepLatest                            // 永远只留 1 个最新值，WithBuffer 被忽略
    OverflowError                                 // 终止该订阅并发 ErrBackpressureOverflow
)
```

```go
var ErrBackpressureOverflow = errors.New("reactivex: backpressure buffer overflow")
```

`OverflowError` 只影响**这一个订阅**，`Subject` 与同辈订阅继续。

### `BackpressureOption`

```go
type BackpressureOption func(*backpressureConfig)

func WithBuffer(size int) BackpressureOption   // 默认 0；负数 panic
func WithOverflow(strategy OverflowStrategy) BackpressureOption
```

- 选项按顺序生效，后写的覆盖先写的；`nil` 跳过。
- `OverflowKeepLatest` 永远使用 1 个挂起位，与 `WithBuffer` 无关。
- `OverflowDropOldest` 与 `WithBuffer(0)` 组合会在构造时 panic。

```go
subject := reactivex.NewSubject[int](
    reactivex.WithBuffer(128),
    reactivex.WithOverflow(reactivex.OverflowDropOldest),
)
```

## `Subject[T]`

热多播发布者 + 订阅者。每个订阅独立需求、独立缓冲（构造时配置），**不**回放 —— 没有订阅者时 `OnNext` 丢弃，新订阅只参与之后的发布；终止状态被保留，迟到订阅直接收到终止通知。

默认无缓冲阻塞模式：有需求时同步调回调；缓冲或非阻塞策略异步派发，但单订阅内仍串行；订阅之间可能并发。

```go
func NewSubject[T any](options ...BackpressureOption) *Subject[T]
func (s *Subject[T]) Subscribe(ctx context.Context, out Subscriber[T]) Subscription
func (s *Subject[T]) ForEach(ctx context.Context, onNext func(T), onError func(error), onComplete func()) Subscription
```

`Subject` 也实现 `Subscriber[T]`，可作为另一个流的桥：

```go
func (s *Subject[T]) OnSubscribe(sub Subscription)
func (s *Subject[T]) OnNext(v T)
func (s *Subject[T]) OnError(err error)
func (s *Subject[T]) OnComplete()
```

> `OnSubscribe` 当前**不**自动向上游 `Request`；下游需求也不会被自动汇总到上游。用 `Subject` 当 `Subscriber` 时，记得自己 `Request`。

## 例子

```go
// 冷的有限流
values, err := reactivex.Just(1, 2, 3, 4).
    Filter(func(v int) bool { return v%2 == 0 }).
    ToSlice(ctx)
// err == nil 时 values == []int{2, 4}
```

```go
// 异步 channel 源，缓冲 + 丢最旧
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
// 热多播
subj := reactivex.NewSubject[string]()
go func() {
    defer subj.OnComplete()
    for _, v := range source {
        subj.OnNext(v)
    }
}()
subj.ForEach(ctx, onMsg, onErr, onDone) // 启动一个订阅
```

## 与其他包的关系

- 同步、单次消费请用 [`collections`](../collections/README-cn.md) 的 `Stream`。
- 一次性的成功 / 失败用 [`result`](../result/README-cn.md)；订阅级错误通过 `OnError` 报出。
- 单值"可能缺席"用 [`option`](../option/README-cn.md)。
