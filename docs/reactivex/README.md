# typed/reactivex 设计说明

## 项目定位

`typed/reactivex` 是 `typed` 项目中的独立模块，目标是在 Go 1.27 上提供类型安全的异步事件流处理能力。

它不是旧版 RxGo 的兼容实现，也不是把 `collections/stream` 改造成异步版本。它会参考 ReactiveX 的概念和 RxGo 的实践经验，重新设计自己的类型、生命周期、背压和并发边界。

预计独立发布为：

```text
github.com/qianwj/typed/reactivex
```

它不应依赖 `typed/collections`。两个模块可以通过标准库的 `iter.Seq` 进行可选互操作，但应该保持独立的依赖和演进节奏。

## 为什么需要 ReactiveX

Go 已经有 goroutine、channel 和 `context.Context`，简单的异步处理通常不需要额外抽象：

```go
for value := range input {
    if value.Valid {
        output <- transform(value)
    }
}
```

但当系统包含多个异步数据源和多个处理阶段时，仅使用 channel 会让以下问题分散在业务代码中：

- 多个数据源的合并、拼接和切换
- 异步操作的错误传播和恢复
- 订阅取消后的 goroutine 和资源回收
- 生产速度大于消费速度时的背压策略
- 多个订阅者共享同一个数据源
- 定时器、窗口、采样和防抖
- 并行处理后的顺序恢复
- 一次性数据源和可重复数据源的区别

ReactiveX 的价值不是替代 channel，而是为这些问题提供统一的组合模型。

## 与 collections 的边界

`typed` 中的集合和流分为三个层次：

```text
ArrayList[T] / HashMap[K, V]
    同步、具体集合、通常立即执行

Stream[T]
    同步、惰性、基于 iter.Seq

Observable[T]
    异步、推送、订阅、取消和背压
```

`collections` 适合处理已经存在的数据，`stream` 适合处理同步迭代器，`reactivex` 适合处理持续产生、异步到达或需要订阅生命周期的数据。

ReactiveX 不应被设计成集合操作的另一套别名，它解决的是异步生产和消费之间的协调问题。

## 核心设计原则

### 类型安全

核心数据类型使用泛型具体类型，不使用 `interface{}` 作为数据通道：

```go
type Observable[T any] struct {
    // subscription implementation
}
```

类型转换通过 Go 1.27 的泛型方法表达：

```go
func (o Observable[T]) Map[R any](
    f func(context.Context, T) (R, error),
) Observable[R]
```

这样 `Observable[User]` 可以安全地转换为 `Observable[Profile]`，而不需要类型断言。

### 具体类型负责链式操作

Go 支持泛型接口，但接口方法不能声明自己的类型参数。因此不能把 `Map[R]` 放进 `Observable[T]` 接口：

```go
// 不支持
type Observable[T any] interface {
    Map[R any](func(T) R) Observable[R]
}
```

`Observable[T]` 应该是具体泛型类型。接口只用于稳定的协作边界：

```go
type Subscriber[T any] interface {
    OnSubscribe(Subscription)
    OnNext(T)
    OnError(error)
    OnComplete()
}
```

### 错误不是普通数据

不采用旧 RxGo 中的 `Item{V interface{}, E error}`。错误应该通过独立的终止通知传播：

```text
OnNext(T)       数据
OnError(error)  异常终止
OnComplete()    正常完成
```

转换函数可以显式返回错误：

```go
Map(func(ctx context.Context, value User) (Profile, error) {
    return loadProfile(ctx, value)
})
```

### 订阅必须有生命周期

每一次订阅都应该有明确的取消和完成路径：

```go
type Subscription interface {
    Request(n uint64)
    Cancel()
    Done() <-chan struct{}
}
```

任何 source 或 operator 创建的 goroutine，都必须能够通过取消信号退出。订阅者提前停止消费时，底层数据源也必须收到取消通知。

## 建议的核心 API

```go
type Observable[T any] struct {
    subscribe func(context.Context, Subscriber[T]) Subscription
}

func (o Observable[T]) Subscribe(
    ctx context.Context,
    subscriber Subscriber[T],
) Subscription

func (o Observable[T]) Map[R any](
    f func(context.Context, T) (R, error),
) Observable[R]

func (o Observable[T]) Filter(
    predicate func(T) bool,
) Observable[T]

func (o Observable[T]) FlatMap[R any](
    f func(T) Observable[R],
) Observable[R]
```

也可以提供函数式订阅的便利 API：

```go
func (o Observable[T]) ForEach(
    ctx context.Context,
    onNext func(T),
    onError func(error),
    onComplete func(),
) Subscription
```

当前模块已经提供上述核心类型和 API。`Observable[T]` 是具体的 cold publisher，`Publisher[T]` 是稳定的订阅边界：

```go
type Publisher[T any] interface {
    Subscribe(context.Context, Subscriber[T]) Subscription
}
```

终止收集可以使用：

```go
values, err := observable.ToSlice(ctx)
total, err := reactivex.Collect(ctx, observable, 0,
    func(sum, value int) int { return sum + value })
```

首批数据源包括 `Just`、`FromSlice`、`FromSeq`、`FromChannel`、`FromChannelWithOptions`、`Create` 和 `Interval`。

## 背压模型

channel 的发送阻塞可以提供一种背压，但 channel 容量本身并不等于消费者 demand。设计需要明确区分：阻塞、缓冲、丢弃最新值、丢弃旧值、只保留最新值以及溢出报错。

如果目标是接近 Reactive Streams 规范，应该使用基于需求量的订阅：

```go
type Subscription interface {
    Request(n uint64)
    Cancel()
    Done() <-chan struct{}
}
```

背压配置使用函数式选项，不把配置塞进一个大型参数列表：

```go
subject := reactivex.NewSubject[int](
    reactivex.WithBuffer(128),
    reactivex.WithOverflow(reactivex.OverflowDropOldest),
)

source := reactivex.FromChannelWithOptions(
    input,
    reactivex.WithBuffer(64),
    reactivex.WithOverflow(reactivex.OverflowError),
)
```

支持的溢出策略有：

| 策略 | 行为 |
| --- | --- |
| `OverflowBlock` | 等待 demand 或缓冲空间，保留所有值 |
| `OverflowDropLatest` | 缓冲区满时丢弃刚到达的值 |
| `OverflowDropOldest` | 缓冲区满时丢弃最早的待处理值 |
| `OverflowKeepLatest` | 只保留一个最新值 |
| `OverflowError` | 终止订阅并发送 `ErrBackpressureOverflow` |

`Request(n)` 表示订阅者可以接收的数量，`WithBuffer` 和 `WithOverflow` 只控制异步边界溢出，两者不是同一个概念。默认策略是 `OverflowBlock`，默认缓冲区大小为零。

如果目标是更贴近 RxGo，可以提供阻塞、缓冲和丢弃策略作为高层便利模式，但必须在 API 文档中说明每种策略的语义。

## Hot 与 Cold Observable

Cold Observable 为每个订阅者创建独立的数据生产过程，适合请求、文件读取和数据库查询。Hot Observable 独立于订阅者持续产生事件，适合行情、日志、设备事件和 WebSocket 消息。

建议将多播能力集中在 `Subject`、`Publish`、`Replay`、`Share` 等明确 API 中，而不是让普通 `Observable` 在不同场景下改变行为。

当前的 `Subject[T]` 同时实现 `Publisher[T]` 和 `Subscriber[T]`：

```go
subject := reactivex.NewSubject[int]()
subject.Subscribe(ctx, subscriber)
subject.OnNext(1)
subject.OnComplete()
```

每个订阅者拥有独立的 demand 和缓冲区。来自同一个 channel 的多个订阅会竞争消费；`Subject` 才是面向多个订阅者的广播入口。

## Pull 与 Push 的桥接

`iter.Seq[T]` 是同步 pull 风格的迭代器，`Observable[T]` 是异步 push 风格的数据源：

```text
iter.Seq[T]       调用者请求下一个值
Observable[T]     数据源主动推送值
```

建议提供：

```go
func FromSeq[T any](seq iter.Seq[T]) Observable[T]
func (o Observable[T]) ToSeq(ctx context.Context) iter.Seq[T]
```

`ToSeq` 必须保证停止 `range` 时取消底层订阅，错误不会被静默丢弃，channel 数据源不会因为没有消费者而永久阻塞，并在文档中明确说明它通常是单次消费。

## 操作符实现思路

操作符不应该默认创建一个新的公开 channel 和永久 goroutine。更合理的实现方向是：

1. `Observable` 保存一个订阅函数。
2. `Map`、`Filter` 等操作返回新的 `Observable`。
3. 新 Observable 在订阅时包装下游 Subscriber。
4. 上游事件通过包装的 Subscriber 流向下游。
5. 取消和错误沿订阅链反向传播。
6. 只有真正需要异步边界的 source 或 operator 才创建 goroutine。

首批操作符建议包括 `Map`、`Filter`、`Take`、`Skip`、`Scan`、`Reduce`、`Concat`、`Merge`、`Zip`、`FlatMap`、`Retry`、`Catch` 和 `Timeout`。时间窗口、调度器、多播和复杂并发操作放在后续阶段。

## 与旧 RxGo 的关系

旧 RxGo 是重要的参考实现，但 `typed/reactivex` 不以兼容旧 API 为目标。

值得保留的经验包括 context 取消、hot/cold Observable、connectable 和 multicast、pool、serialize、retry，以及 goleak 形式的 goroutine 泄漏测试。

需要重新设计的部分包括：

- 用 `Observable[T]` 替代非泛型 `Observable` 接口
- 用 `T` 替代 `interface{}` 和类型断言
- 用独立的 `OnError` / `OnComplete` 替代 `Item` 错误包装
- 用明确的 `Subscription` 替代分散的 `Option` 生命周期控制
- 将 source、operator、subject 和 scheduler 分开
- 将背压策略从普通操作符配置中独立出来

因此，该模块应当是一个新的设计，而不是旧 RxGo 的机械泛型化。

## 与 collections/stream 的关系

| 特性 | `collections/stream` | `reactivex` |
| --- | --- | --- |
| 数据模型 | 同步迭代 | 异步事件推送 |
| 消费方式 | `for range` / collect | subscribe |
| 错误 | 普通返回值或显式错误 | `OnError` 终止通知 |
| 生命周期 | 迭代结束 | cancel / error / complete |
| 背压 | 提前停止迭代 | demand、buffer、drop 等策略 |
| 多播 | 通常不支持 | Subject、Share、Replay |
| 典型数据源 | slice、map、iter.Seq | channel、timer、网络和事件源 |

不要为了统一命名而强行统一实现。两个模块应该在边界处互操作，在内部保持各自的模型清晰。

## 非目标

第一阶段不考虑：

- 兼容旧 RxGo 的全部方法和签名
- 用反射支持任意运行时类型
- 默认并行执行所有操作符
- 把 channel 的所有用法都包装成 Observable
- 同时实现完整的 Java Reactor 调度器体系
- 在没有明确语义的情况下自动判断 hot/cold
- 用一个巨大接口列出所有操作符

## 实施路线

### 阶段一：核心协议

- `Observable[T]`
- `Subscriber[T]`
- `Subscription`
- `Just`
- `FromSlice`
- `FromSeq`
- `Map`
- `Filter`
- `Take`
- `ForEach`
- context 取消

### 阶段二：错误和组合

- `OnError`
- `OnComplete`
- `FlatMap`
- `Concat`
- `Merge`
- `Zip`
- `Retry`
- `Timeout`

### 阶段三：背压和并发

- `Request(n)`
- bounded buffer
- drop/latest 策略
- 并发 Map
- 顺序恢复
- goroutine 生命周期测试

### 阶段四：多播和时间

- `Subject[T]`
- `Publish`
- `Replay`
- `Share`
- `Debounce`
- `Throttle`
- `Buffer`
- `Window`
- `Interval`

每个阶段都需要配套测试：正常完成、错误终止、取消、提前停止、慢消费者、多个订阅者和 goroutine 泄漏。

## 待决定问题

以下问题在实现前需要形成明确决策：

1. 是否严格实现 Reactive Streams 的 `Request(n)`，还是先提供阻塞式背压？
2. `OnNext` 是否允许并发调用，还是默认保证串行通知？
3. `Subscribe` 是否立即启动 source，还是由 `Request` 启动？
4. `FromChannel` 的取消是否负责关闭输入 channel？
5. 错误发生后是否允许恢复为新的 Observable？
6. `Subject` 是否保证订阅者之间的顺序一致？
7. 是否提供独立 scheduler，还是优先使用显式 goroutine 和 context？

这些问题会直接影响 API 兼容性和资源安全，因此应当在第一版实现前确定，而不是隐藏在 `Option` 中。
