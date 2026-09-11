# typed

`typed` 是一个基于 Go 泛型的类型安全集合工具集。它提供具体的、立即执行的集合类型(`ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]`),一个轻量的惰性 `Stream[T]` 层,以及专门的线性数据结构(`Stack[T]`、`Queue[T]`、`Deque[T]`),和支撑这些集合的抽象(`Option[T]`、`Result[T]`、`Equaler[T]`)。

集合的设计参考 Java / JavaScript:具体泛型类型,带链式方法如 `Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`。所有可能"缺失"的访问器(`Get`、`First`、`Last`、`Find`、`MinBy`、`MaxBy`、`Pop`、`Peek`、`Front`、`Back`)返回 `Option[T]` 而不是 `(T, bool)`,调用方可以在结果上直接链式 `OrElse` / `OrElseGet` / `Map`。

## 当前状态

核心类型和 `Option[T]` / `Result[T]` 抽象已经落地并经过测试。下文「开发路线」部分列出了尚未启动的工作。

## 目标

Go 的 `for` 循环非常清晰,应当继续作为简单逻辑的首选。但当一个集合需要经过多个连续变换时,嵌套函数或重复的临时切片会让处理流程变得不易阅读:

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

`typed` 支持下面这种从左到右的集合处理方式:

```go
result := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()
```

如果需要惰性处理,可以显式转换成 Stream:

```go
result := lists.ArrayListOf(users...).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## 设计原则

- **类型安全优先。** 优先使用 Go 泛型,尽量避免 `any`、反射和运行时类型断言。`Option[T]` 使用显式的 `present` 标记,而不是依赖 nil 检查,所以它对值类型(`int`、`string`、`struct{}` 等)也有效。
- **具体类型优先于接口。** `ArrayList[T]` 等都是具体泛型类型,不是接口。Go 1.27 的泛型方法(`Map[R]`、`FlatMap[R]`、`Reduce[R]`)只能在具体接收者上工作;接口无法让这些方法返回接口类型,会破坏链式调用。
- **基于 Optional 的访问。** 任何可能"缺失"的访问器都返回 `option.Optional[T]` 而不是 `(T, bool)`。这一约定在 `ArrayList.Get / First / Last / Find / MinBy / MaxBy / RemoveFirst / RemoveLast`、`LinkedList`、`Stack`、`Queue`、`Deque` 中保持一致。
- **有界内存。** `ArrayList` 使用 head offset 布局并周期性压缩,所以长期从头部排空的列表的保留容量,被"在飞元素的高水位 + 64"这个常数所界定,而不是被该列表曾经达到过的历史最大值所界定。`Stack.Pop`、`ArrayList.RemoveFirst / RemoveLast`、`Deque.PopFront / PopBack` 会把释放的槽位清零,以便运行时的可达性扫描不会让被弹出元素的引用继续存活。
- **可选惰性。** 集合操作是立即执行的;`Stream[T]` 是显式的惰性层,基于 Go 的 `iter.Seq[T]` 实现。
- **可组合。** 集合、迭代器、`Stream` 可以组合成新的数据源;`Stream()` 的快照语义保证后续对源集合的修改不会泄漏到进行中的管道中。
- **提前停止。** `First`、`Any`、`All`、`Find`、`Take`,以及返回 `Optional[T]` 的访问器,在答案确定时会立刻停止。
- **Go 风格。** 不盲目复制 Java Stream 的每一条语义;简单逻辑仍然应该用 `for range` 写。工具集是 opt-in 的,偏好普通 slice 和 map 的现有代码不受影响。

## 已交付的内容

### 集合类型

| 类型 | 类别 | 源码 | 说明 |
| --- | --- | --- | --- |
| `ArrayList[T]` | 具体泛型 struct | `collections/lists` | head-offset `[]T`;`Add` / `AddFirst` / `RemoveFirst` / `RemoveLast` 都是 O(1);head ≥ 64 时周期性压缩。 |
| `LinkedList[T]` | 具体泛型 struct | `collections/lists` | 双向链表;头尾 O(1),随机访问 O(i)。 |
| `HashMap[K, V]` | 具体泛型 struct | `collections/maps` | 开放寻址哈希表;`K comparable`。 |
| `HashSet[T]` | 具体泛型 struct | `collections/sets` | `T comparable`;复用同一套哈希表实现。 |
| `Stack[T]` | 具体泛型 struct | `collections` | 单端 LIFO;`[]T` 加 `Pop` 时显式清零槽位。 |
| `Queue[T]` | 具体泛型 struct | `collections` | 单端 FIFO;`ArrayList[T]` 的薄包装。 |
| `Deque[T]` | 具体泛型 struct | `collections` | 双端;`LinkedList[T]` 的薄包装。所有操作都是严格 O(1)。 |
| `Stream[T]` | 具体泛型 struct | `collections/stream` | 惰性、单次消费的管道,基于 `iter.Seq[T]`。 |

### 父包构造函数

| 函数 | 源码 | 说明 |
| --- | --- | --- |
| `Range[T constraints.Integer](start, end T) Stream[T]` | `collections` | iota 风格的 `[start, end)` 惰性整数流。 |

### 抽象类型

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `option.Optional[T]` | `utils/option` | 存在 / 缺失的值,不依赖对 `T` 的 nil 检查。 |
| `result.Result[T]` | `utils/result` | 成功 / 失败;`Unwrap` 返回 `(T, error)`,`Recover` 实现错误感知的 fallback。 |
| `objects.Equaler[T]` | `utils/objects` | 静态接口 `Equal(T) bool`。 |
| `objects.IsNil[T]`、`objects.Equals[T]` | `utils/objects` | 基于反射的 nil 检查和相等分派(支持 typed nil、`Equaler[T]`、`time.Time.Equal`)。 |

### 控制流

| 类型 | 源码 | 说明 |
| --- | --- | --- |
| `match.Case[T, R]` | `control/match` | Go 1.27 泛型方法下的模式匹配原语。 |

## 速览

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/utils/option"
)

// 立即执行的转换。
adults := lists.ArrayListOf(users...).
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// 基于 Optional 的访问:不用写 (T, bool) 那一套。
firstAdult := adults.First().OrElse("(none)")

// 线性结构。
s := collections.NewStack[int]()
s.Push(1); s.Push(2); s.Push(3)
top := s.Pop().OrElse(0)  // 3

q := collections.NewQueue[int]()
q.Push(1); q.Push(2)
front := q.Pop().OrElse(0)  // 1

d := collections.NewDeque[int]()
d.PushBack(1); d.PushFront(0); d.PushBack(2)
// [0, 1, 2]
left  := d.PopFront().OrElse(-1)  // 0
right := d.PopBack().OrElse(-1)   // 2

// Optional 链式。
v := option.Of(7).Map(func(x int) int { return x * 2 }).OrElse(0)  // 14
```

## 集合的链式 API

`ArrayList[T]`、`LinkedList[T]`、`HashSet[T]` 都提供同一组核心操作:

| 能力 | API |
| --- | --- |
| 构造 | `NewX[T]()`、`XOf(values...)`,都返回 `*X[T]` |
| 数量 | `Size()` |
| 生命周期 | `IsEmpty()`、`Clear()` |
| 遍历 | `ForEach`、`Collect`、`Stream` |
| 条件查询 | `Any`、`All`、`None`、`Find`(返回 `Optional[T]`) |
| 筛选 | `Filter`,返回同种类型 |
| 类型变换 | `Map[R](func(T) R)`,返回 `*ArrayList[R]` |
| 扁平化 | `FlatMap[R](func(T) *ArrayList[R])`,返回 `*ArrayList[R]` |
| 折叠 | `Reduce[R](init, func(R, T) R)` |
| 观察和拼接 | `Peek`、`Concat`,返回同种类型 |

`ArrayList[T]` 和 `LinkedList[T]` 额外提供依赖顺序的操作(`Get`、`Insert`、`RemoveAt`、`First`、`Last`、`Take`、`Drop`、`Distinct`、`SortBy`、`MinBy`、`MaxBy`)。`HashMap[K, V]` 提供 `Get`、`GetOrDefault`、`Put`、`Remove`、`Keys`、`Values`、`Entries`、`Filter`、`MapValues[R]`、`Concat`。`HashSet[T]` 提供 `Union`、`Intersect`、`Difference`、`SymmetricDifference`、`IsSubsetOf`、`IsSupersetOf`,以及当结果满足 `comparable` 时保留去重语义的 `MapSet` / `FlatMapSet`。

`Map` 和 `FlatMap` 始终返回 `*ArrayList[R]`,因为结果类型 `R` 可能是 slice、map、函数或任何不可比较的类型。`HashSet.MapSet` 和 `HashSet.FlatMapSet` 是保留去重语义的版本。

## Optional 访问

`Optional[T]` 是任何"可能缺失"的访问器的标准返回类型:

| 操作 | 缺失时的结果 |
| --- | --- |
| `Optional.Get()`(缺失时) | panic |
| `Optional.OrElse(default)` | `default` |
| `Optional.OrElseGet(f)` | `f()` |
| `Optional.OrElseThrow(msg)` | `(zero, errors.New(msg))` |
| `Optional.IsPresent` / `IsEmpty` | bool |
| `Optional.Map[R](f)` | 缺失的 `Optional[R]`,不会调用 `f` |
| `Optional.FlatMap[R](f)` | 缺失的 `Optional[R]`,不会调用 `f` |
| `Optional.Filter(predicate)` | 谓词为 false 时缺失 |
| `Optional.IfPresent(f)` / `IfPresentOrElse(p, a)` | no-op 或 `a()` |

存在 / 缺失标志是显式存储的,而不是从 `T` 上推断 nil,所以 `Optional[int]` 也能用于值类型,避免了只能用 `(int, bool)` 的处境。

`Result[T]` 是与之配套的 `(T, error)` 形态。`Unwrap` 在成功路径上返回值和 `nil` 错误,在失败路径上返回零值和捕获的错误。`Recover(f func(error) Result[T])` 允许针对特定错误做 fallback,不必先 unwrap:

```go
v, err := loadProfile(id).Unwrap()
if err != nil { return err }
```

`Result.MapError` 就地变换捕获的错误;`Map[R]` / `FlatMap[R]` 在成功路径上传递值,失败路径上的错误保持不变。

## 内存模型

`ArrayList` 把活跃区间放在 `items[head:head+size]`,当 `head >= 64` 时把丢弃的前缀向前折叠归零。长期从头部排空的列表,保留容量被"在飞元素的高水位 + 64"这个常数所界定,而不是被该列表曾经达到过的历史最大值所界定。从头部弹出的引用类型元素,只要 `RemoveFirst` 把释放的槽位清零,就立即可以被 GC 回收。`collections/lists/arraylist_test.go` 和 `collections/{stack,queue,deque}_test.go` 中的引用处理测试用 `runtime.SetFinalizer` 断言所有弹出 box 的 finalizer 都会运行。

`Stack.Pop` 是 slice-backed 的等价物:它截断底层数组并显式清零弹出的槽位,所以 `Stack[T]`(`T` 是指针)不会因为越界槽位而泄漏引用。

## 并发

所有集合类型都不支持并发修改。Go 的标准模式(单 goroutine 拥有集合,通信通过 channel 发生)继续适用。`Stream` 默认不并行;普通的 `Map` 不会静默变成并发版本。工具集不引入新的并发模型;它遵循 Go 显式并发的约定。

## 项目结构

```text
typed/
├── go.work
├── collections/
│   ├── go.mod
│   ├── mod.go
│   ├── stack.go, stack_test.go
│   ├── queue.go, queue_test.go
│   ├── deque.go, deque_test.go
│   ├── lists/
│   │   ├── arraylist.go, arraylist_test.go
│   │   ├── linkedlist.go
│   │   └── lists_test.go
│   ├── maps/
│   │   └── hashmap.go, hashmap_test.go
│   ├── sets/
│   │   └── hashset.go, hashset_test.go
│   └── stream/
│       └── mod.go, stream_test.go
├── utils/
│   ├── go.mod
│   ├── option/
│   ├── result/
│   └── objects/
├── control/
│   ├── go.mod
│   └── match/
└── docs/
    ├── collections/README.md
    ├── option/README.md
    ├── result/README.md
    ├── control/README.md
    └── reactivex/README.md
```

## Java / JavaScript 映射

| Java / JavaScript | `typed` |
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
| 惰性 Stream | `Stream[T]` |
| `IntStream.range` | `Range(start, end) Stream[T]` |
| `Deque`(Java) | `Deque[T]`(LinkedList-backed) |

项目不试图复制 Java 或 JavaScript 的运行时模型。它借鉴了它们的集合处理风格,同时保留 Go 的静态类型、显式错误和直接的控制流。

## 惰性与执行边界

典型的 `Stream` 管道如下:

```text
source → 中间操作 → 中间操作 → 终止操作
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

## 错误处理

`Result[T]` 是类型化的 `(T, error)` 载体:

```go
v, err := result.Success(42).Unwrap()
if err != nil { return err }

fallback := result.Failure[int](errors.New("missing")).
    Recover(func(err error) result.Result[int] {
        if errors.Is(err, ErrMissing) { return result.Success(0) }
        return result.Failure[int](err)
    })
```

`Result.MapError` 就地变换捕获的错误;`Map[R]` / `FlatMap[R]` 在成功路径上传递值,失败路径上的错误保持不变。

对于需要暴露错误的集合管道,推荐的做法是把错误路径转成缺失的 `Optional[T]`(例如对返回 `(T, error)` 的查询用 `OfNullable`),并把成功路径留在普通链式 API 上。

## Go 版本

当前模块使用 Go 1.27:

```text
go 1.27.1
```

Go 1.23 引入了 `iter.Seq`、`iter.Seq2` 以及对函数迭代器的 `for range` 支持;Go 1.27 的泛型方法使 `Stream[T].Map[R]`、`Optional[T].Map[R]`、`Result[T].Map[R]` 等链式 API 成为可能。

## 开发路线

已完成:

- [x] `ArrayList[T]`、`LinkedList[T]`、`HashMap[K, V]`、`HashSet[T]` 及它们的链式方法。
- [x] 立即执行的集合操作:`Filter`、`Map[R]`、`FlatMap[R]`、`Reduce[R]`、`Collect`、`ForEach`、`Peek`、`Concat`。
- [x] 列表的顺序相关操作:`Get`、`Insert`、`RemoveAt`、`First`、`Last`、`Find`、`Take`、`Drop`、`Distinct`、`SortBy`、`MinBy`、`MaxBy`。
- [x] 基于 Optional 的访问:`Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy` / `RemoveFirst` / `RemoveLast` 返回 `Optional[T]`。
- [x] `Stream[T]` 适配器,基于 `iter.Seq[T]`,带可提前停止的终止操作。
- [x] 线性结构:`Stack[T]`、`Queue[T]`、`Deque[T]`。
- [x] 父包辅助函数:`Range[T constraints.Integer](start, end T) Stream[T]`,iota 风格的整数序列。
- [x] `Option[T]` / `Result[T]` / `Equaler[T]` / `IsNil[T]` / `Equals[T]` 工具集。
- [x] 有界内存:head-offset `ArrayList` 加周期性压缩,`Stack.Pop` 与列表的 `Remove*` 路径清零释放的槽位。
- [x] 启用 race detector 的测试,活跃开发的文件 100% 语句覆盖。

待完成:

- [ ] 错误感知的集合操作(`MapE`、`FilterE`、`CollectE`),作为 `Result[T]`-per-element 的一等管道替代品。
- [ ] 基准测试:`ArrayList` 头部操作、`LinkedList` 迭代、`HashMap` rehash 行为、`Stream` 管道开销。
- [ ] 迭代器(`iter.Seq[T]`)作为集合类型的一等输出,与 `Stream()` 并列。

## 许可证

本项目使用 [MIT License](LICENSE)。
