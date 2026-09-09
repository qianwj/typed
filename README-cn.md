# typed

`typed` 是一个基于 Go 泛型的类型安全集合工具集。

项目的第一个目标是提供类似 Java/JavaScript 的具体集合类型，例如 `ArrayList[T]` 和 `HashMap[K, V]`，让集合操作可以自然地进行链式调用。同时提供可选的惰性 `Stream[T]` 层，用于需要延迟执行的数据处理流程。

> 项目目前处于早期设计阶段。`collections/stream` 现在还是 API 占位实现，下面的链式代码代表目标方向，不保证当前版本可以直接编译。

## 目标

Go 的 `for` 循环非常清晰，也应当继续作为简单逻辑的首选。但当一个集合需要经过多个连续变换时，嵌套函数或重复的临时切片会让处理流程不容易阅读：

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

`typed` 希望支持下面这种从左到右的集合处理方式：

```go
result := typed.ArrayList[User](users).
    Filter(func(u User) bool {
        return u.Age >= 18
    }).
    Map(func(u User) string {
        return u.Name
    })
```

如果需要惰性处理，可以显式转换成 Stream：

```go
result := typed.ArrayList[User](users).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## 设计方向

- **类型安全**：依赖 Go 泛型，尽量避免 `any`、反射和运行时类型断言。
- **集合优先**：`ArrayList[T]` 和 `HashMap[K, V]` 是普通集合处理的主要类型。
- **链式调用**：集合和 Stream 操作都返回带类型的值，支持自然的左到右阅读顺序。
- **按需惰性**：集合操作简单直接、默认立即执行；需要惰性执行和提前停止时使用 `Stream[T]`。
- **可组合**：集合、迭代器和其他 Stream 可以组合成新的数据源。
- **可提前停止**：支持 `First`、`Any`、`Take` 等操作，避免不必要地处理剩余元素。
- **Go 风格**：不强行复制 Java Stream 的全部语义；简单逻辑仍然应当可以直接使用 `for range`。

## API 草案

### 集合类型

计划中的集合类型是泛型具体类型，而不是接口：

```go
type ArrayList[T any] []T
type HashMap[K comparable, V any] map[K]V
```

这样可以支持改变元素类型的泛型方法：

```go
names := ArrayList[User](users).
    Filter(isAdult).
    Map(func(u User) string { return u.Name })
```

`ArrayList[T]` 面向有顺序的、类似 slice 的数据；`HashMap[K, V]` 面向 `Filter`、`MapValues`、`Keys`、`Values` 和 `ToSlice` 等键值操作。具体命名以及立即执行和惰性执行的边界仍在设计中。

### 惰性 Stream

`Stream[T]` 是基于 Go 迭代器约定的可选惰性层，适合大型数据、单次数据源、channel 数据源或潜在的无限序列：

```go
profiles := ArrayList[User](users).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Collect()
```

Stream 不会取代集合类型，也不会取代普通的 `for range` 循环。

计划中的操作大致分为三类。

### 中间操作

根据接收者类型的不同，这些操作通常返回另一个集合或 `Stream`：

- `Filter`
- `Map`
- `FlatMap`
- `Distinct`
- `Take` / `Drop`
- `Skip` / `Limit`
- `Concat`
- `Peek`

### 终止操作

这些操作会真正消费 Stream：

- `Collect`
- `Count`
- `First` / `Last`
- `Any` / `All` / `None`
- `Find`
- `Reduce`
- `ToMap`
- `GroupBy`
- `ForEach`

### 排序与聚合

- `Sort` / `SortBy`
- `Min` / `Max`
- `Sum` / `Average`
- `GroupBy`
- `PartitionBy`
- `Join`

具体 API 会根据 Go 的错误处理、泛型方法能力和惰性迭代器语义逐步确定。

## 与 go-linq 的关系

[`go-linq`](https://github.com/ahmetb/go-linq) 是一个重要的参考项目，并且已经为 Go 1.27 提供了完整的类型安全 LINQ 实现。它的核心抽象是惰性的 `Query[T]`，提供 `Where`、`Select`、`GroupBy`、`Join` 和 `Aggregate` 等丰富的 LINQ 风格操作。

`typed` 在顶层设计上有意采用不同方向：

| 对比项 | `go-linq` | `typed` 方向 |
| --- | --- | --- |
| 核心抽象 | 惰性 `Query[T]` | 先提供具体的 `ArrayList[T]` 和 `HashMap[K, V]` |
| 命名风格 | LINQ 风格：`Where`、`Select` | 集合风格：`Filter`、`Map` |
| 执行方式 | 默认惰性 | 集合默认立即执行，`Stream[T]` 显式惰性 |
| JavaScript 数组体验 | 通过 Query 适配 | 将 `ArrayList[T]` 作为一等目标 |
| Map 操作 | 将键值对转换为 Query | 将 `HashMap[K, V]` 作为一等目标 |
| 错误处理 | 不是主要模型 | 计划支持显式错误传播操作 |

目标不是重复实现 `go-linq`，而是探索一种更适合 Go、Java 和 JavaScript 开发者迁移的集合模型，同时兼容 `iter.Seq[T]`，并在真正需要时提供惰性 Stream。

## Java / JavaScript 对照

| Java / JavaScript | `typed` 方向 |
| --- | --- |
| `stream()` | `ArrayList(values).Stream()` |
| `filter` | `Filter` |
| `map` | `Map` |
| `flatMap` | `FlatMap` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` | `Limit` |
| `skip` | `Skip` |
| `findFirst` | `First` |
| `anyMatch` | `Any` |
| `allMatch` | `All` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |

这不是把 Java 或 JavaScript 的运行时模型原样搬到 Go，而是借鉴它们在集合处理上的表达方式，同时保留 Go 的静态类型、显式错误和简单控制流。

## 惰性与执行边界

一个典型的 Stream 处理流程会分为：

```text
数据源 -> 中间操作 -> 中间操作 -> 终止操作
```

例如：

```go
adults := stream.From(users).
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

在 `Collect` 之前，`Filter`、`Map` 和 `Take` 只描述处理流程。终止操作开始消费数据，并且 `Take(100)` 可以让底层数据源在满足数量后提前停止。

## 错误处理方向

Go 的函数经常返回 `(value, error)`，因此 Stream API 不会简单地隐藏错误。对于可能失败的转换，计划支持显式的错误传播，例如：

```go
profiles, err := stream.From(users).
    MapE(loadProfile).
    Collect()
```

错误处理的具体形式仍在设计中，优先保证错误不会被静默丢弃，并且能够在链路中尽早停止。

## 并发边界

Stream 默认不自动并行执行。Go 中的 goroutine、channel、锁和取消信号都有明确的并发语义，自动并行化可能引入不可预测的开销和生命周期问题。

未来可以考虑显式的并发操作，但不会把普通的 `Map` 默认变成并发版本。

## Go 版本

当前模块使用 Go 1.27：

```text
go 1.27.1
```

目标 API 会使用 Go 泛型和标准迭代器能力。Go 1.23 引入了 `iter.Seq`、`iter.Seq2` 以及对函数迭代器的 `for range` 支持；Go 1.27 的泛型方法进一步使类似 `Stream[T].Map[R]` 的链式 API 成为可能。

## 项目结构

```text
typed/
└── collections/
    ├── go.mod
    ├── arraylist/
    ├── hashmap/
    └── stream/
        └── mod.go
```

## 开发路线

1. 确定 `ArrayList[T]` 和 `HashMap[K, V]` 的语义与命名。
2. 实现 `Filter`、`Map`、`FlatMap` 和 `Collect` 等立即执行的集合操作。
3. 增加基于 `iter.Seq[T]` 的 `Stream()` 适配器。
4. 增加 `First`、`Any`、`All`、`Take` 等可提前停止的 Stream 操作。
5. 设计错误传播 API，包括 `MapE` 和 `FilterE` 的取舍。
6. 增加排序、分组、聚合和 Map 处理能力。
7. 为立即执行集合、惰性 Stream、空集合和提前停止补充测试与基准测试。

## 许可证

本项目使用 [MIT License](LICENSE)。
