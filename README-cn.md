<div align="center">

<img src="./docs/assets/typed-logo-gopher-official-light.png" alt="Typed — 面向 Go 泛型的链式类型安全工具集" width="180" />

# Typed

**面向 Go 泛型的链式类型安全工具集。**

[English](./README.md) · [文档索引](./docs/README-cn.md)

[![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=root)](https://codecov.io/gh/qianwj/typed)
[![Module](https://img.shields.io/badge/module-github.com%2Fqianwj%2Ftyped-6f42c1)](#安装)
[![Made with Go](https://img.shields.io/badge/made%20with-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)

</div>

Typed 是一个**面向 Go 泛型的链式类型安全工具集**。具体的泛型类型（`ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]`、`Stream[T]`、`Subject[T]`）通过从左到右、可一路链下去的 transform（`Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Take`、`Drop`、`SortBy`、`Distinct`、`Concat`）组合而成，背后是支撑整个工具集的抽象（`Option[T]`、`Result[T]`、`IsNil`、`Equals`）。操作符的命名是有意沿用业界通用词汇，对熟悉 Java Streams、.NET LINQ 或 JavaScript 数组管道的工程师会很顺手，但**不沿用它们各自的运行时模型**。异步一侧，`reactivex.Flowable[T]` 提供了强类型的事件流：显式需求、可配置背压，以及一等公民的多播 `Subject[T]`。

## 为什么选 Typed

Go 的 `for` 循环本身已经足够清晰,简单逻辑应当继续用它。问题是当一个集合需要经过多个连续变换时:嵌套函数、散落的 `if err != nil`、无法链式拼装的 `(T, bool)` 返回值,会让处理流程变得难以阅读。

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

Typed 支持从左到右、可以一路链下去的处理方式:

```go
names := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()
```

需要惰性求值时,显式转成 `Stream` 即可:

```go
profiles := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## 一眼概览

- 🧱 **具体泛型类型** —— `ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]`、`Stack[T]`、`Queue[T]`、`Deque[T]`、`Stream[T]`、`Subject[T]`。无运行时类型断言,没有 `any` 带来的意外。
- 🪄 **链式 transform** —— `Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Take`、`Drop`、`Distinct`、`SortBy`、`Concat`。借 Go 1.27 泛型方法,全部返回具体类型。
- 🟢 **Option / Result** —— `Option[T]` 表示"可能缺席",`Result[T]` 表示"可能失败"。两者都能干净地与 `(T, error)` 互转,也能和集合 API 拼装。
- 📡 **强类型异步流** —— `reactivex.Flowable[T]`,显式需求(`Request(n)`)、按订阅配置背压(`WithBuffer`、`OverflowStrategy`),以及热多播的 `Subject[T]`。
- 🧠 **更智能的相等** —— `objects.Equals[T]` 理解 `func (T) Equal(T) bool`、对 nil/空集合做归一化,对 typed nil 指针 nil 安全。
- 🪶 **有界内存** —— `ArrayList` 用 head-offset 布局并周期性压缩,`Stack.Pop` 和 `Remove*` 路径把释放的槽位清零,被弹出的引用不会因底层数组残留。
- 🧵 **有界与无界阻塞队列** —— `concurrency.BoundedBlockingQueue[T]` 是固定容量的 FIFO,`Push` / `Poll` 阻塞,`TryPush` / `TryPoll` 不阻塞,底层是单个 `chan T`。`concurrency.UnboundedBlockingQueue[T]` 是它的兄弟类型,`Push` 永不阻塞,底层是环形缓冲区 + mutex + cond。两者都暴露了带 context 的变体(`PushWithContext` / `PollWithContext`)。
- 🧵 **结构化并发** —— `concurrency.Group` 是 `errgroup` 风格的 helper,带两种失败策略:`Strict`(任何任务失败即失败,ctx 取消兄弟) 和 `BestEffort`(任务全部跑完,只有全部成功才成功,否则返回聚合后的 error)。直接构建在 `sync.WaitGroup` 上,零外部依赖。

## 安装

```bash
go get github.com/qianwj/typed/collections
go get github.com/qianwj/typed/adt
go get github.com/qianwj/typed/adt
go get github.com/qianwj/typed/adt
go get github.com/qianwj/typed/reactivex
go get github.com/qianwj/typed/control
go get github.com/qianwj/typed/concurrency
go get github.com/qianwj/typed/utils/objects
go get github.com/qianwj/typed/utils/json
```

每个子包是独立的 Go module,按需引入即可。需要 **Go 1.27+** 才能使用具体类型上的泛型方法。

## 速览

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/adt"
    "github.com/qianwj/typed/adt"
)

// 立即执行的 transform
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// 基于 Option 的访问 —— 不用写 (T, bool) 那一套
first := adults.First().OrElse("(none)")

// 线性结构
s := collections.NewStack[int]()
s.Push(1); s.Push(2); s.Push(3)
top := s.Pop().OrElse(0) // 3

// Result: 把 (T, error) 串成组合子链
v, err := adt.Wrap(loadConfig(path)).
    MapError(func(err error) error { return fmt.Errorf("config %s: %w", path, err) }).
    Map(func(b []byte) Config { return parseConfig(b) }).
    Unwrap()
if err != nil { return err }
use(v)
```

## 目录

- [设计原则](#设计原则)
- [当前已经发布的能力](#当前已经发布的能力)
- [文档导航](#文档导航)
- [操作符参考](#操作符参考)
- [惰性与执行边界](#惰性与执行边界)
- [内存模型](#内存模型)
- [并发模型](#并发模型)
- [错误处理](#错误处理)
- [稳定性与 SemVer](#稳定性与-semver)
- [Go 版本要求](#go-版本要求)
- [开发路线](#开发路线)
- [贡献指南](#贡献指南)
- [Bug 反馈](#bug-反馈)
- [许可证](#许可证)

## 设计原则

- **类型安全优先。** 优先使用 Go 泛型,尽量避免 `any`、反射和运行时类型断言。`Option[T]` 用显式的 `present` 标志位,而不是依赖 nil 检查,所以它对值类型(`int`、`string`、`struct{}` 等)也有效。
- **具体类型优先于接口。** `ArrayList[T]` 等都是具体泛型类型,不是接口。Go 1.27 的泛型方法(`Map[R]`、`FlatMap[R]`、`Reduce[R]`)只能在具体接收者上工作;接口无法让这些方法返回接口类型,会破坏链式调用。
- **基于 Option 的访问。** 任何可能"缺失"的访问器都返回 `adt.Option[T]` 而不是 `(T, bool)`。这一约定在 `ArrayList.Get / First / Last / Find / MinBy / MaxBy / RemoveFirst / RemoveLast`、`LinkedList`、`Stack`、`Queue`、`Deque` 中保持一致。
- **有界内存。** `ArrayList` 使用 head offset 布局并周期性压缩,所以长期从头部排空的列表的保留容量,被"在飞元素的高水位 + 64"这个常数所界定,而不是被该列表曾经达到过的历史最大值所界定。`Stack.Pop`、`ArrayList.RemoveFirst / RemoveLast`、`Deque.PopFront / PopBack` 会把释放的槽位清零,以便运行时的可达性扫描不会让被弹出元素的引用继续存活。
- **可选惰性。** 集合操作是立即执行的;`Stream[T]` 是显式的惰性层,基于 Go 的 `iter.Seq[T]` 实现。
- **可组合性。** 集合、迭代器和 `Stream` 可以组合成新的数据源;`Stream()` 的快照契约保证源集合的后续变更不会泄漏到进行中的管道。
- **提前结束。** `First`、`Any`、`All`、`Find`、`Take` 以及返回 `Option` 的访问器一旦得到答案就立刻停下。
- **可选,不替代。** Typed 是建立在 Go 惯用法之上的链式层;简单逻辑仍应能用 `for range` 配合 slice、map 轻松写出来。仍然偏好内置 slice、map 和 `chan T` 的代码完全不受影响。

## 当前已经发布的能力

### 集合类型

| 类型 | 形态 | 源码 | 说明 |
| --- | --- | --- | --- |
| `ArrayList[T]` | 具体泛型结构体 | `collections/lists` | Head-offset `[]T`;`Add` / `AddFirst` / `RemoveFirst` / `RemoveLast` 都是 O(1);`head >= 64` 时周期性压缩。 |
| `LinkedList[T]` | 具体泛型结构体 | `collections/lists` | 双向链表;头尾 O(1),按下标 O(i)。 |
| `HashMap[K, V]` | 具体泛型结构体 | `collections/maps` | 开地址哈希表;`K comparable`。 |
| `HashSet[T]` | 具体泛型结构体 | `collections/sets` | `T comparable`;构建在同款哈希表之上。 |
| `Stack[T]` | 具体泛型结构体 | `collections` | 单端 LIFO;`[]T` 存储,`Pop` 时显式清零。 |
| `Queue[T]` | 具体泛型结构体 | `collections` | 单端 FIFO;`ArrayList[T]` 的薄包装。 |
| `Deque[T]` | 具体泛型结构体 | `collections` | 双端;`LinkedList[T]` 的薄包装。所有操作严格 O(1)。 |
| `Stream[T]` | 具体泛型结构体 | `collections/stream` | 基于 `iter.Seq[T]` 的惰性、单次消费管道。 |

### 父包构造器

| 函数 | 源码 | 说明 |
| --- | --- | --- |
| `Range[T constraints.Integer](start, end T) Stream[T]` | `collections` | 整数半开区间 `[start, end)` 的惰性 iota。 |

### 抽象类型

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `adt.Option[T]` | `adt` | 存在 / 缺失的值,不依赖对 `T` 的 nil 检查。 |
| `adt.Either[L, R]` | `adt` | tagged-union 值类型;约定 `Left` = 失败、`Right` = 成功。安全访问器返回 `Option`;`Fold` 是选分支的规范组合子。 |
| `adt.Result[T]` | `adt` | 成功 / 失败;`Unwrap` 返回 `(T, error)`,`Wrap` 是从 `(T, error)` 到 `Result[T]` 的正向桥,`Recover` 是错误感知的 fallback,总返回 `T`。 |
| `json.Encode[T] / Decode[T]` | `utils/json` | 基于 `encoding/json/v2` 的 `Result` 风格编解码。 |
| `objects.Equaler`(接口,可选用) | `utils/objects` | 提示接口 `Equal(any) bool`;`Equals` 不要求类型实现它。 |
| `objects.IsNil[T]`、`objects.Equals[T]` | `utils/objects` | 基于反射的 nil 检查和相等分派(支持 typed nil、`func (T) Equal(T) bool`、`time.Time.Equal`)。 |

### 控制流

| 类型 / 函数 | 源码 | 说明 |
| --- | --- | --- |
| `control.Repeat(times, f)` / `RepeatE(times, f) (int, error)` | `control` | "做 N 次"循环;`RepeatE` 在首次非 nil 错误处停,返回已成功次数。 |
| `control.If[T](condition, onTrue, onFalse)` / `IfGet[T](condition, onTrue, onFalse)` | `control` | `If` 从提前求值的参数中取值；`IfGet` 只执行选中的回调，恰好一次。 |
| `match.Pattern[T]`、`match.Value[T]`、`match.Type(any)` | `control/match` | 首个匹配获胜的模式匹配:`Pattern[T]` 测试值,`Type` 按动态类型分派。 |

受 Go 实参求值规则限制，`If(user != nil, user.Name, "匿名")` 在 user 为 nil 时仍会 panic。可用 `IfGet` 回调，或通过 `adt.OfNullable(user).Map(...).OrElse(...)` 表达可选值。示例与编译器语义说明见[条件取值文档](./docs/control/README-cn.md#if-与-ifget)。

### 响应式流

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `reactivex.Flowable[T]`、`Publisher[T]`、`Subscriber[T]` | `reactivex` | 强类型异步流,显式需求(`Subscription.Request(n)`)与 `OnError` / `OnComplete` 终止信号。 |
| `reactivex.Subject[T]` | `reactivex` | 热多播发布者 + 订阅者;通过 `WithBuffer` / `WithOverflow` 配置。 |
| `reactivex.Single[T]` | `reactivex` | reactive 容器，恰好发一个值或一个 error。`Await` / `AwaitWithContext` 返回 `adt.Result[T]`，`Subscribe` 用于回调消费。 |
| `reactivex.Maybe[T]` | `reactivex` | reactive 容器，包含有值、空完成、错误三种终态。`Await` / `AwaitWithContext` 返回 `adt.Result[adt.Option[T]]`，空完成是成功结果。 |
| `reactivex.OverflowStrategy` | `reactivex` | `OverflowBlock` / `OverflowDropLatest` / `OverflowDropOldest` / `OverflowKeepLatest` / `OverflowError`。 |
| 源:`Just` / `FromSlice` / `FromChannel` / `FromChannelWithOptions` / `FromSeq` / `Create` / `Interval` | `reactivex` | 冷、热源构造器。 |
| 算子:`Map[R]`、`Filter`、`Take`、`Skip`、`Scan[R]`、`Reduce` | `reactivex` | 全部为包装型,自身不启 goroutine、不带队列。 |

### 并发原语

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `BoundedBlockingQueue[T]` | `concurrency` | 固定容量 FIFO(容量向上取整到 2 的幂);`Push` / `Poll` 阻塞,`TryPush` / `TryPoll` 不阻塞;`TryPoll` 返回 `adt.Option[T]`。`chan T` 的薄泛型包装。热路径上 `0 B/op`、`0 allocs/op`。 |
| `UnboundedBlockingQueue[T]` | `concurrency` | 无界 FIFO;`Push` 永不阻塞,`Poll` 在空队列时阻塞。底层是环形缓冲区 + `sync.Mutex` + `*sync.Cond`。API 表面和 `BoundedBlockingQueue` 一致。 |
| `Group` | `concurrency` | `errgroup` 风格的结构化并发,两种失败策略 —— `Strict`(任何任务失败即失败,ctx 取消兄弟)与 `BestEffort`(任务跑完,只有全部成功才成功,否则返回 `*BestEffortError`)。直接构建在 `sync.WaitGroup` 上,零外部依赖。 |
| `Pool[T]` | `concurrency` | `sync.Pool` 的泛型封装，通过 `NewPool(creator)`、`Get() adt.Option[T]`、`Put(T)` 复用临时对象。无值或 nil 结果返回空 Option；调用方负责重置，并在归还后停止访问。 |

## 文档导航

按包组织的 API 参考与示例,中英双语:

- [docs/README-cn.md](./docs/README-cn.md) —— 索引
- [collections](./docs/collections/README-cn.md) —— `ArrayList` / `LinkedList` / `HashMap` / `HashSet` / `Stack` / `Queue` / `Deque` / `Stream` / `Range`
- [control](./docs/control/README-cn.md) —— `If` / `IfGet`、`Repeat` / `RepeatE` 与 `control/match`
- [reactivex](./docs/reactivex/README-cn.md) —— `Flowable` / `Subject` / 背压 / 算子
- [concurrency](./docs/concurrency/README-cn.md) —— 阻塞队列、`Group`、`Semaphore` 与 `Pool[T]`
- [adt](./docs/adt/README-cn.md) —— `Option[T]`
- [adt](./docs/adt/README-cn.md) —— `Result[T]`
- [utils/objects](./docs/utils/objects/README-cn.md) —— `IsNil` / `Equals`
- [utils/json](./docs/utils/json/README-cn.md) —— `Encode` / `Decode`,基于 `encoding/json/v2`

## 操作符参考

操作符的命名沿用 Java Streams、.NET LINQ 与 JavaScript 数组方法之间共通的"业界通用词汇",所以只要用过其中任意一套,这套 API 读起来就很自然。下表把常见的别称映射到 Typed 对应的接口。

| 操作 | Typed |
| --- | --- |
| `stream()` | `ArrayListOf(...).Stream()` |
| `filter` / `where` | `Filter` |
| `map` / `select` | `Map` |
| `flatMap` / `selectMany` | `FlatMap` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` / `take` | `Take` |
| `skip` / `drop` | `Drop` |
| `findFirst` / `find` | `First` / `Find`(返回 `Option[T]`) |
| `anyMatch` / `some` | `Any` |
| `allMatch` / `every` | `All` |
| `noneMatch` | `None` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |
| `Option.of` / `ofNullable` | `adt.Of` / `adt.OfNullable` |
| `Option.orElse` | `OrElse` |
| `Stream` lazy | `Stream[T]` |
| `IntStream.range` | `Range(start, end) Stream[T]` |
| `Deque` (Java) | `Deque[T]`(LinkedList-backed) |

Typed 不是任何一个库的移植品。链式 API 直接构建在 Go 自带的泛型、`iter.Seq[T]` 与显式 `error` 返回之上;操作符的名字之所以眼熟,只是因为沿用了工程师们本来就会的那套"业界通用词汇"。

## 惰性与执行边界

一个典型的 `Stream` 管道如下:

```text
source → intermediate operation → intermediate operation → terminal operation
```

例如:

```go
adults := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

`Collect` 之前,`Filter`、`Map`、`Take` 只描述管道。终止操作开始消费数据,`Take(100)` 可以在产出足够数据时立即停止底层 source。`Stream()` 在调用时创建源集合的快照,所以之后对源集合的修改不会影响进行中的管道。

## 内存模型

`ArrayList` 把活区放在 `items[head:head+size]`,`head >= 64` 时把已丢弃的前缀归零并折叠。长期从头部排空的列表的保留容量,因此被"在飞元素的高水位 + 64"这个常数所界定,而不是曾经达到过的历史最大值。引用元素从头部弹出时,只要 `RemoveFirst` 把对应槽位清零,就立刻可被 GC。`collections/lists/arraylist_test.go` 和 `collections/{stack,queue,deque}_test.go` 里的引用回收测试用 `runtime.SetFinalizer` 验证被弹出的 box 都会触发 finalizer。

`Stack.Pop` 是 slice 版的对等实现:它收缩底层数组并把弹出的槽位显式清零,所以 `Stack[T]` 持有指针时不会因为越界槽位而泄漏引用。

## 并发模型

任何集合类型都不支持并发修改。Go 的标准做法 —— 单 goroutine 持有集合,通信走 channel —— 继续适用。`Stream` 默认不并行;普通的 `Map` 不会悄悄变成并发的。Typed 不引入新的并发模型,沿用 Go 显式并发的路子。

少数场景下,工具集自己的 API 比 `chan T` 更顺手:非阻塞探测(`TryPush` / `TryPoll`)、带 context 的阻塞(`PushWithContext` / `PollWithContext`)、非阻塞读返回 `Option`、一个泛型签名就能讲清“固定容量的阻塞队列”、或者 `errgroup` 风格的 typed 结构化并发。这时 [`concurrency` 包](./docs/concurrency/README-cn.md) 提供了 [`BoundedBlockingQueue[T]`](./docs/concurrency/README-cn.md#boundedblockingqueuet)(`chan T` 的薄包装,每操作开销基本为零)、[`UnboundedBlockingQueue[T]`](./docs/concurrency/README-cn.md#unboundedblockingqueuet)(环形缓冲区 + mutex + cond,因为 Go runtime 没有“无界 buffered channel”)和 [`Group`](./docs/concurrency/README-cn.md#group)(直接构建在 `sync.WaitGroup` 上的 strict 或 best-effort 结构化并发)。需要容量做背压用 `BoundedBlockingQueue`;需要 `Push` 永不阻塞用 `UnboundedBlockingQueue`;需要并发编排且不想自己写 `sync.WaitGroup + context + recover` 模板时用 `Group`。

## 错误处理

`Result[T]` 是类型化的 `(T, error)` 载体:

```go
v, err := adt.Success(42).Unwrap()
if err != nil { return err }

port, _ := adt.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err 始终为 nil
```

`Recover` 总会返回 `T`;需要把新错误往上传,请用 `Unwrap` 收尾而不是 `Recover`。

对需要暴露错误的集合管道,推荐把错误路径转成缺失的 `Option[T]`(例如对返回 `(T, error)` 的查询用 `OfNullable`),并把成功路径留在普通链式 API 上。

## 稳定性与 SemVer

每个 module 都遵循 [Semantic Versioning 2.0.0](https://semver.org/)。当前刻意停在 `v1.0.0` 之前:

- **`0.0.x`(当前)。** 没有 API 稳定性承诺。工具集刻意保持小巧,改一个方法签名成本低、早期受众也小,所以破坏性变更的代价同样低。Patch 发布(`0.0.x`)只包含 bug 修复;只要涉及公开签名、导出类型、或者某个算子已有文档行为,一律算作 **minor** 升级(`0.0.x` → `0.(x+1).0`)。
- **`0.y.0`(模块进入 feature freeze 后规划)。** 公开类型、导出函数签名、算子的文档行为全部冻结。Patch 发布只修 bug 和文档,不再加新的 API 面。模块在 API 经历过至少一个 minor 周期且没有再改动时升到这一档。
- **`v1.0.0`。** 只给"维护者愿意 backport bug fix"的模块预留。当前推荐用 `use at HEAD` 安装。

### 哪些改动算 minor 之间的破坏性变更

- 重命名导出类型或方法。
- 修改泛型接收者(比如 `Map(func(T) R)` 改成 `Map(func(context.Context, T) (R, error))`)。
- 给公开 struct 加新的必填字段。
- 改写文档里已有的不变量(例如 "`Map` 对 nil slice 返回空 flowable"、"`MinBy` 在空列表上 panic")。

### 不算破坏性变更

- 给已有类型加新方法。
- 加新的顶层函数或子包。
- 给 variadic options 函数加新选项。
- 在"文档已声明不在范围内"的输入上改大 O(例如每个集合类型头上都标了"非并发安全")。
- bug fix 改写了文档里已经说"可能发生"的可观察行为。

### 每个 module 独立的发布 tag

每个 module 用自己的 tag 前缀:

```text
utils/v0.0.1
collections/v0.0.1
control/v0.0.1
reactivex/v0.0.1
concurrency/v0.0.1
```

某个 module tag 的变动不会带动另一个;`utils/v0.0.2` 和 `reactivex/v0.0.1` 可能同年发布,也可能隔几年。前缀彼此独立。

### 怎么查破坏性变更

每次发布的 release notes 里都会给一行 `BREAKING:` 摘要签名或行为变动。上一版和当前 `vX.Y.Z` tag 之间的 diff 才是权威 changelog,release notes 只是摘要。

## Go 版本要求

当前 module 使用 **Go 1.27**:

```text
go 1.27.1
```

Go 1.23 引入了 `iter.Seq`、`iter.Seq2` 以及对函数迭代器的 `for range` 支持;Go 1.27 的泛型方法使 `Stream[T].Map[R]`、`Option[T].Map[R]`、`Result[T].Map[R]` 等链式 API 成为可能。

## 开发路线

已完成:

- [x] `ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]` 及它们的链式方法。
- [x] 立即执行的集合操作:`Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Collect`、`ForEach`、`Peek`、`Concat`。
- [x] 列表的顺序相关操作:`Get`、`Insert`、`RemoveAt`、`First`、`Last`、`Find`、`Take`、`Drop`、`Distinct`、`SortBy`、`MinBy`、`MaxBy`。
- [x] 基于 Option 的访问:`Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy` / `RemoveFirst` / `RemoveLast` 以及 `Stack` / `Queue` / `Deque` 的 `Pop` / `Peek` / `Front` / `Back` 都返回 `Option[T]`。
- [x] `Stream[T]` 适配器,基于 `iter.Seq[T]`,带可提前停止的终止操作。
- [x] 线性结构:`Stack[T]`、`Queue[T]`、`Deque[T]`。
- [x] 父包辅助函数:`Range[T constraints.Integer](start, end T) Stream[T]`,iota 风格的整数序列。
- [x] `Option[T]` / `Result[T]` / `Equaler` / `IsNil[T]` / `Equals[T]` 工具集。
- [x] `control.If` / `IfGet`、`Repeat` / `RepeatE` 与 `control/match` 模式匹配。
- [x] `reactivex` 包:`Flowable[T]` / `Publisher[T]` / `Subscriber[T]`、`Subject[T]`、`WithBuffer` / `WithOverflow` 背压、`Map` / `Filter` / `Take` / `Skip` / `Scan` / `Reduce` 算子。
- [x] `concurrency` 包:`BoundedBlockingQueue[T]`(数组环形缓冲,阻塞 + 非阻塞双 API,通过 race 测试)。
- [x] 基于 `encoding/json/v2` 的 `utils/json` `Result` 风格编解码。
- [x] `concurrency.Pool[T]`：基于 `sync.Pool` 的泛型临时对象复用。
- [x] 有界内存:head-offset `ArrayList` 加周期性压缩,`Stack.Pop` 与列表的 `Remove*` 路径清零释放的槽位。
- [x] 启用 race detector 的测试,覆盖每个子包。各模块的实时覆盖率通过 Codecov 报告 —— 见上方 badge。

待完成:

- [ ] 错误感知的集合操作(`MapE`、`FilterE`、`CollectE`),作为 `Result[T]`-per-element 的一等管道替代品。
- [ ] 基准测试:`ArrayList` 头部操作、`LinkedList` 迭代、`HashMap` rehash 行为、`Stream` 管道开销。
- [ ] 迭代器(`iter.Seq[T]`)作为集合类型的一等输出,与 `Stream()` 并列。

## 贡献指南

工具集刻意保持小巧,但欢迎贡献。在开 PR 之前,请先扫一眼下面的原则,
review 流程会顺畅得多。

- **按 module 划定范围。** 每个顶层包(`collections`、`control`、
  `reactivex`、`concurrency`、`adt`、`adt`、
  `utils/objects`、`utils/json`)都是 **独立的 Go module**,拥有各自的
  tag 前缀(见 [稳定性与 SemVer](#稳定性与-semver))。跨多个 module 的
  修改,必须在 PR 描述里明确点出,并说明为什么需要跨 module 联动;否则
  维护者会让你拆 PR。
- **先讨论再动手。** bug fix、文案小修之外的工作,请先开 issue 沟通,
  把改动形态对齐。表面看像"加一个方法"的改动,经常会演变成 SemVer /
  API 面决策(见 [稳定性与 SemVer](#稳定性与-semver) 里关于"什么算破坏性
  变更"的说明)。
- **测试是硬性要求。** 改变行为的 PR 必须附带在对应 module 上能通过
  `go test -race ./...` 的测试。工具集的正确性故事依赖这条;
  没有回归测试的 PR 会被要求补上。
- **Lint 必须保持干净。** 推送前在你改动的 module 下跑一次
  `golangci-lint run ./...`。CI 在每次 push 和 PR 上跑的是同一份配置
  (`./.golangci.yml`)。
- **公开 API 是禁区。** 只要动了导出类型、方法、函数或常量,PR 描述
  里必须显式说明,并打上 `BREAKING:` / `feat:` / `fix:` 标签,这样 release
  notes 才能写对。具体的 tag 约定见 [每个 module 独立的发布 tag](#每个-module-独立的发布-tag)。
- **提交粒度。** 一次提交只做一件事。bug fix 与重构不要和功能改动
  混在同一个 commit 里;每条 commit message 单独看也要讲得通。
- **代码风格。** 跟随你正在编辑的文件。仓库不引入未经讨论的新依赖,
  请你也按这条走;拿不准时,先读两个相邻的文件再下笔。

标准流程:

1. 在 GitHub 上 fork 仓库,从 `main` 切出一个 topic 分支。
2. 写代码,在对应 module 下跑 `go test -race ./...` 与
   `golangci-lint run`,再 push 分支。
3. 在 `qianwj/typed:main` 上开 PR,描述里写清改了什么、关联了哪个
   issue(若有),以及是否带 `BREAKING:` / `feat:` / `fix:` 标签。

## Bug 反馈

请用 GitHub 上的 [issue 列表](https://github.com/qianwj/typed/issues)。
开新 issue 之前先搜一下,确认没被人提过(包括已关闭的——修复可能
已经落在别的 module 分支上,只是还没打 tag)。

一份合格的 bug 报告应当包含:

- **目标。** 一句话讲清楚:你在用 Typed 做什么。
- **Module 与版本。** 受影响的 module(`collections` / `reactivex` /
  `concurrency` / ...),以及对应的 commit / tag。在 module 目录下
  `git rev-parse HEAD`,或直接给出该 module 的 `vX.Y.Z` tag。
- **环境。** `go version`、`go env GOOS GOARCH`,必要时附上操作系统 /
  内核版本。
- **复现步骤。** 能触发 bug 的最小可独立运行片段。优先用 Go test,
  或者 `go run` 能直接跑的 file;只有"几行代码拼一起"的描述通常
  不够复现。
- **期望 vs 实际。** 期望发生什么、实际发生什么,以及完整报错 / panic /
  日志输出。panic 一定要附上完整堆栈。
- **数据竞争相关。** 如果怀疑是数据竞争,请显式说出来,并把
  `-race` 的输出原样贴上来。仓库里大部分并发 bug 只能在
  `go test -race` 或对长跑进程挂了 runtime race detector 后才能复现。
- **绕开方法。** 你试过的任何能让你继续往下走的方法,哪怕很丑——
  它通常能反推出底层的假设。

安全问题 **不要** 走公开 issue 通道。请用 GitHub 的
[私密安全上报](https://github.com/qianwj/typed/security/advisories/new)
流程,以便维护者在公开披露之前协调修复。

## 许可证

本项目使用 [MIT License](./LICENSE)。
