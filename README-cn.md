<div align="center">

<img src="./docs/assets/typed-logo-gopher-official-light.png" alt="Typed — 基于 Go 泛型的类型安全集合工具集" width="180" />

# Typed

**基于 Go 泛型的类型安全集合工具集。**

[English](./README.md) · [文档索引](./docs/README-cn.md)

[![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=root)](https://codecov.io/gh/qianwj/typed)
[![Module](https://img.shields.io/badge/module-github.com%2Fqianwj%2Ftyped-6f42c1)](#安装)
[![Made with Go](https://img.shields.io/badge/made%20with-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)

</div>

Typed 把 Java / JavaScript 风格的集合操作体验带到 Go，但不放弃静态类型。具体的泛型类型（`ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]`、`Stream[T]`、`Subject[T]`），链式 transform（`Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Take`、`Drop`、`SortBy`、`Distinct`、`Concat`），以及支撑一切的抽象（`Optional[T]`、`Result[T]`、`IsNil`、`Equals`）。再加上强类型的异步事件流（`reactivex.Observable[T]`），显式需求、可配置背压，以及一等公民的多播 `Subject[T]`。

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
- 🟢 **Optional / Result** —— `Optional[T]` 表示"可能缺席",`Result[T]` 表示"可能失败"。两者都能干净地与 `(T, error)` 互转,也能和集合 API 拼装。
- 📡 **强类型异步流** —— `reactivex.Observable[T]`,显式需求(`Request(n)`)、按订阅配置背压(`WithBuffer`、`OverflowStrategy`),以及热多播的 `Subject[T]`。
- 🧠 **更智能的相等** —— `objects.Equals[T]` 理解 `func (T) Equal(T) bool`、对 nil/空集合做归一化,对 typed nil 指针 nil 安全。
- 🪶 **有界内存** —— `ArrayList` 用 head-offset 布局并周期性压缩,`Stack.Pop` 和 `Remove*` 路径把释放的槽位清零,被弹出的引用不会因底层数组残留。
- 🧵 **有界阻塞队列** —— `concurrency.BoundedBlockingQueue[T]` 是固定容量的 FIFO,`Push` / `Take` 阻塞,`TryPush` / `TryTake` 不阻塞,底层是单个环形缓冲区。

## 安装

```bash
go get github.com/qianwj/typed/collections
go get github.com/qianwj/typed/utils/option
go get github.com/qianwj/typed/utils/result
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
    "github.com/qianwj/typed/utils/option"
    "github.com/qianwj/typed/utils/result"
)

// 立即执行的 transform
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// 基于 Optional 的访问 —— 不用写 (T, bool) 那一套
first := adults.First().OrElse("(none)")

// 线性结构
s := collections.NewStack[int]()
s.Push(1); s.Push(2); s.Push(3)
top := s.Pop().OrElse(0) // 3

// Result: 把 (T, error) 串成组合子链
v, err := result.Wrap(loadConfig(path)).
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
- [Java / JavaScript 对照表](#java--javascript-对照表)
- [惰性与执行边界](#惰性与执行边界)
- [内存模型](#内存模型)
- [并发模型](#并发模型)
- [错误处理](#错误处理)
- [Go 版本要求](#go-版本要求)
- [开发路线](#开发路线)
- [许可证](#许可证)

## 设计原则

- **类型安全优先。** 优先使用 Go 泛型,尽量避免 `any`、反射和运行时类型断言。`Optional[T]` 用显式的 `present` 标志位,而不是依赖 nil 检查,所以它对值类型(`int`、`string`、`struct{}` 等)也有效。
- **具体类型优先于接口。** `ArrayList[T]` 等都是具体泛型类型,不是接口。Go 1.27 的泛型方法(`Map[R]`、`FlatMap[R]`、`Reduce[R]`)只能在具体接收者上工作;接口无法让这些方法返回接口类型,会破坏链式调用。
- **基于 Optional 的访问。** 任何可能"缺失"的访问器都返回 `option.Optional[T]` 而不是 `(T, bool)`。这一约定在 `ArrayList.Get / First / Last / Find / MinBy / MaxBy / RemoveFirst / RemoveLast`、`LinkedList`、`Stack`、`Queue`、`Deque` 中保持一致。
- **有界内存。** `ArrayList` 使用 head offset 布局并周期性压缩,所以长期从头部排空的列表的保留容量,被"在飞元素的高水位 + 64"这个常数所界定,而不是被该列表曾经达到过的历史最大值所界定。`Stack.Pop`、`ArrayList.RemoveFirst / RemoveLast`、`Deque.PopFront / PopBack` 会把释放的槽位清零,以便运行时的可达性扫描不会让被弹出元素的引用继续存活。
- **可选惰性。** 集合操作是立即执行的;`Stream[T]` 是显式的惰性层,基于 Go 的 `iter.Seq[T]` 实现。
- **可组合性。** 集合、迭代器和 `Stream` 可以组合成新的数据源;`Stream()` 的快照契约保证源集合的后续变更不会泄漏到进行中的管道。
- **提前结束。** `First`、`Any`、`All`、`Find`、`Take` 以及返回 `Optional` 的访问器一旦得到答案就立刻停下。
- **保持 Go 风格。** 不照搬 Java Stream 的全部语义;简单逻辑仍应能用 `for range` 轻松写出来。Typed 是可选的:仍然偏好内置 slice 和 map 的代码不受影响。

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
| `option.Optional[T]` | `utils/option` | 存在 / 缺失的值,不依赖对 `T` 的 nil 检查。 |
| `result.Result[T]` | `utils/result` | 成功 / 失败;`Unwrap` 返回 `(T, error)`,`Wrap` 是从 `(T, error)` 到 `Result[T]` 的正向桥,`Recover` 是错误感知的 fallback,总返回 `T`。 |
| `json.Encode[T] / Decode[T]` | `utils/json` | 基于 `encoding/json/v2` 的 `Result` 风格编解码。 |
| `objects.Equaler`(接口,可选用) | `utils/objects` | 提示接口 `Equal(any) bool`;`Equals` 不要求类型实现它。 |
| `objects.IsNil[T]`、`objects.Equals[T]` | `utils/objects` | 基于反射的 nil 检查和相等分派(支持 typed nil、`func (T) Equal(T) bool`、`time.Time.Equal`)。 |

### 控制流

| 类型 / 函数 | 源码 | 说明 |
| --- | --- | --- |
| `control.Repeat(times, f)` / `RepeatE(times, f) (int, error)` | `control` | "做 N 次"循环;`RepeatE` 在首次非 nil 错误处停,返回已成功次数。 |
| `match.Pattern[T]`、`match.Value[T]`、`match.Type(any)` | `control/match` | 首个匹配获胜的模式匹配:`Pattern[T]` 测试值,`Type` 按动态类型分派。 |

### 响应式流

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `reactivex.Observable[T]`、`Publisher[T]`、`Subscriber[T]` | `reactivex` | 强类型异步流,显式需求(`Subscription.Request(n)`)与 `OnError` / `OnComplete` 终止信号。 |
| `reactivex.Subject[T]` | `reactivex` | 热多播发布者 + 订阅者;通过 `WithBuffer` / `WithOverflow` 配置。 |
| `reactivex.OverflowStrategy` | `reactivex` | `OverflowBlock` / `OverflowDropLatest` / `OverflowDropOldest` / `OverflowKeepLatest` / `OverflowError`。 |
| 源:`Just` / `FromSlice` / `FromChannel` / `FromChannelWithOptions` / `FromSeq` / `Create` / `Interval` | `reactivex` | 冷、热源构造器。 |
| 算子:`Map[R]`、`Filter`、`Take`、`Skip`、`Scan[R]`、`Reduce` | `reactivex` | 全部为包装型,自身不启 goroutine、不带队列。 |

### 并发原语

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `BoundedBlockingQueue[T]` | `concurrency` | 固定容量 FIFO(容量向上取整到 2 的幂);`Push` / `Take` 阻塞,`TryPush` / `TryTake` / `DrainTo` 不阻塞;`TryTake` 返回 `option.Optional[T]`;基于 `[]T` 的环形缓冲区,单 `sync.Mutex` + `*sync.Cond`。热路径上 `0 B/op`、`0 allocs/op`。 |

## 文档导航

按包组织的 API 参考与示例,中英双语:

- [docs/README-cn.md](./docs/README-cn.md) —— 索引
- [collections](./docs/collections/README-cn.md) —— `ArrayList` / `LinkedList` / `HashMap` / `HashSet` / `Stack` / `Queue` / `Deque` / `Stream` / `Range`
- [control](./docs/control/README-cn.md) —— `Repeat` / `RepeatE` 与 `control/match`
- [reactivex](./docs/reactivex/README-cn.md) —— `Observable` / `Subject` / 背压 / 算子
- [concurrency](./docs/concurrency/README-cn.md) —— `BoundedBlockingQueue[T]`(阻塞 + 非阻塞,固定容量)
- [utils/option](./docs/option/README-cn.md) —— `Optional[T]`
- [utils/result](./docs/result/README-cn.md) —— `Result[T]`
- [utils/objects](./docs/utils/objects/README-cn.md) —— `IsNil` / `Equals`
- [utils/json](./docs/utils/json/README-cn.md) —— `Encode` / `Decode`,基于 `encoding/json/v2`

## Java / JavaScript 对照表

| Java / JavaScript | Typed |
| --- | --- |
| `stream()` | `ArrayListOf(...).Stream()` |
| `filter` / `where` | `Filter` |
| `map` / `select` | `Map` |
| `flatMap` / `selectMany` | `FlatMap` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` / `take` | `Take` |
| `skip` / `drop` | `Drop` |
| `findFirst` / `find` | `First` / `Find`(返回 `Optional[T]`) |
| `anyMatch` / `some` | `Any` |
| `allMatch` / `every` | `All` |
| `noneMatch` | `None` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |
| `Optional.of` / `ofNullable` | `option.Of` / `option.OfNullable` |
| `Optional.orElse` | `OrElse` |
| `Stream` lazy | `Stream[T]` |
| `IntStream.range` | `Range(start, end) Stream[T]` |
| `Deque` (Java) | `Deque[T]`(LinkedList-backed) |

Typed 不复制 Java 或 JavaScript 的运行时模型,只借鉴它们的集合处理风格,保留 Go 的静态类型、显式错误和直观的控制流。

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

少数场景下,工具集自己的 API 比 `chan T` 更顺手:非阻塞探测(`TryPush` / `TryTake`)、与下一次操作一致的 `Size()`、或者一个泛型签名就能讲清“固定容量的阻塞队列”。这时 [`concurrency` 包](./docs/concurrency/README-cn.md) 提供了 [`BoundedBlockingQueue[T]`](./docs/concurrency/README-cn.md#boundedblockingqueuet)。原生的 `chan` 在裸吞吐上仍然快约 5–7×,只有额外的 API 表达力能换回这部分成本时,再选 `BoundedBlockingQueue`。

## 错误处理

`Result[T]` 是类型化的 `(T, error)` 载体:

```go
v, err := result.Success(42).Unwrap()
if err != nil { return err }

port, _ := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err 始终为 nil
```

`Recover` 总会返回 `T`;需要把新错误往上传,请用 `Unwrap` 收尾而不是 `Recover`。

对需要暴露错误的集合管道,推荐把错误路径转成缺失的 `Optional[T]`(例如对返回 `(T, error)` 的查询用 `OfNullable`),并把成功路径留在普通链式 API 上。

## Go 版本要求

当前 module 使用 **Go 1.27**:

```text
go 1.27.1
```

Go 1.23 引入了 `iter.Seq`、`iter.Seq2` 以及对函数迭代器的 `for range` 支持;Go 1.27 的泛型方法使 `Stream[T].Map[R]`、`Optional[T].Map[R]`、`Result[T].Map[R]` 等链式 API 成为可能。

## 开发路线

已完成:

- [x] `ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]` 及它们的链式方法。
- [x] 立即执行的集合操作:`Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Collect`、`ForEach`、`Peek`、`Concat`。
- [x] 列表的顺序相关操作:`Get`、`Insert`、`RemoveAt`、`First`、`Last`、`Find`、`Take`、`Drop`、`Distinct`、`SortBy`、`MinBy`、`MaxBy`。
- [x] 基于 Optional 的访问:`Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy` / `RemoveFirst` / `RemoveLast` 以及 `Stack` / `Queue` / `Deque` 的 `Pop` / `Peek` / `Front` / `Back` 都返回 `Optional[T]`。
- [x] `Stream[T]` 适配器,基于 `iter.Seq[T]`,带可提前停止的终止操作。
- [x] 线性结构:`Stack[T]`、`Queue[T]`、`Deque[T]`。
- [x] 父包辅助函数:`Range[T constraints.Integer](start, end T) Stream[T]`,iota 风格的整数序列。
- [x] `Option[T]` / `Result[T]` / `Equaler` / `IsNil[T]` / `Equals[T]` 工具集。
- [x] `control.Repeat` / `RepeatE` 与 `control/match` 模式匹配。
- [x] `reactivex` 包:`Observable[T]` / `Publisher[T]` / `Subscriber[T]`、`Subject[T]`、`WithBuffer` / `WithOverflow` 背压、`Map` / `Filter` / `Take` / `Skip` / `Scan` / `Reduce` 算子。
- [x] `concurrency` 包:`BoundedBlockingQueue[T]`(数组环形缓冲,阻塞 + 非阻塞双 API,通过 race 测试)。
- [x] 基于 `encoding/json/v2` 的 `utils/json` `Result` 风格编解码。
- [x] 有界内存:head-offset `ArrayList` 加周期性压缩,`Stack.Pop` 与列表的 `Remove*` 路径清零释放的槽位。
- [x] 启用 race detector 的测试,覆盖每个子包。各模块的实时覆盖率通过 Codecov 报告 —— 见上方 badge。

待完成:

- [ ] 错误感知的集合操作(`MapE`、`FilterE`、`CollectE`),作为 `Result[T]`-per-element 的一等管道替代品。
- [ ] 基准测试:`ArrayList` 头部操作、`LinkedList` 迭代、`HashMap` rehash 行为、`Stream` 管道开销。
- [ ] 迭代器(`iter.Seq[T]`)作为集合类型的一等输出,与 `Stream()` 并列。
- [ ] `concurrency` 包:`PushCtx` / `TakeCtx` 支持取消。

## 许可证

本项目使用 [MIT License](./LICENSE)。
