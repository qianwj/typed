# typed/collections 使用说明

`typed/collections` 是一个基于 Go 泛型的集合工具模块，提供有序列表、无序键值表、无序集合，以及基于 `iter.Seq` 的惰性数据流。

模块路径：

```text
github.com/qianwj/typed/collections
```

当前模块使用 Go 1.27。集合类型是具体的泛型类型，而不是接口，这样 `Map[R]`、`FlatMap[R]` 等类型变换方法可以保留完整的静态类型。`ArrayList`、`LinkedList` 和 `HashSet` 的构造器都返回指针，便于统一使用可变集合。

## 模块结构

```text
collections/
├── lists/    ArrayList[T]、LinkedList[T]
├── maps/     HashMap[K, V]
├── sets/     HashSet[T]
└── stream/   Stream[T]
```

常用导入：

```go
import (
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/collections/maps"
    "github.com/qianwj/typed/collections/sets"
)
```

## 统一 API

`ArrayList`、`LinkedList` 和 `HashSet` 都提供下面这组通用操作；`HashMap` 也使用相同的 `Size()` 数量命名。

| 能力 | API |
| --- | --- |
| 构造 | `NewX[T]()`、`XOf(values...)`，返回 `*X[T]` |
| 数量 | `Size()` |
| 生命周期 | `IsEmpty()`、`Clear()` |
| 遍历和导出 | `ForEach`、`Collect`、`Stream` |
| 条件查询 | `Any`、`All`、`None`、`Find` |
| 筛选 | `Filter`，返回同一种集合 |
| 类型变换 | `Map[R](func(T) R)`，返回 `*ArrayList[R]` |
| 扁平化 | `FlatMap[R](func(T) *ArrayList[R])`，返回 `*ArrayList[R]` |
| 折叠 | `Reduce[R](init, func(R, T) R)` |
| 观察和拼接 | `Peek`、`Concat`，返回同一种集合 |

所有容器统一使用 `Size()` 查询元素或条目数量，容器类型不再提供 `Len()` 或 `Count()`。`Stream` 的 `Count()` 是消费流的终止操作，语义独立。

`Map` 和 `FlatMap` 统一返回 `ArrayList`，因为结果类型 `R` 可以是 slice、map、函数等不可比较类型。HashSet 另提供 `MapSet` 和 `FlatMapSet`，在结果满足 `comparable` 时保留去重语义。

列表拥有顺序和索引，因此额外提供 `Get`、`Insert`、`RemoveAt`、`First`、`Last`、`Take`、`Drop`、`Distinct` 和 `SortBy`。HashSet 没有这些依赖稳定顺序的操作，额外提供 `Union`、`Intersect`、`Difference`、`SymmetricDifference`、`IsSubsetOf` 和 `IsSupersetOf`。

## 设计边界

`ArrayList`、`LinkedList`、`HashMap` 和 `HashSet` 是立即执行的具体集合。`Filter`、`Map`、`FlatMap` 等操作调用后就会生成结果。

`Stream` 是显式的惰性层，底层使用 Go 的 `iter.Seq[T]`。调用 `Collect`、`Count`、`First`、`Any` 等终止操作时，数据才会真正被消费。

集合的 `Stream()` 会在调用时创建快照，因此之后对集合的修改不会影响已经创建的 Stream。Stream 本身是单次消费的，终止操作执行后不应再次使用同一个 Stream。

## ArrayList[T]

`ArrayList[T]` 是有序、可变长度的集合，内部使用私有 slice 保存数据。构造器和返回集合的变换方法都返回 `*ArrayList[T]`。

```go
users := lists.ArrayListOf(
    User{Name: "Alice", Age: 17},
    User{Name: "Bob", Age: 21},
    User{Name: "Carol", Age: 30},
)

names := users.
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()

// names == []string{"Bob", "Carol"}
```

主要方法：

| 方法 | 行为 |
| --- | --- |
| `NewArrayList[T]()` | 创建空列表 |
| `ArrayListOf(values...)` | 按传入顺序创建列表 |
| `Add(value)` | 追加到列表尾部 |
| `AddFirst(value)` | 插入到列表头部 |
| `Insert(index, value)` | 在指定位置插入；使用 `Insert(list.Size(), value)` 追加，越界时 panic |
| `RemoveFirst` / `RemoveLast` | 删除并返回首/尾元素，没有元素时返回 `(zero, false)` |
| `RemoveAt(index)` | 删除并返回指定位置的元素，越界时 panic |
| `Get(index)` | 返回 `(value, ok)` |
| `Size()` | 返回元素数量 |
| `IsEmpty` / `Clear` | 查询或清空列表 |
| `Filter(p)` | 返回满足条件的新列表 |
| `Map[R](f)` | 映射为任意类型的新列表 |
| `FlatMap[R](f)` | 拼接每个元素返回的 `ArrayList` |
| `Take(n)` / `Drop(n)` | 截取或跳过元素 |
| `Distinct(eq)` | 根据相等函数保留第一次出现的元素 |
| `SortBy(less)` | 根据比较函数返回排序后的新列表 |
| `First` / `Last` / `Find` | 查询元素 |
| `Any` / `All` / `None` | 条件判断 |
| `ForEach` / `Peek` / `Concat` | 遍历、观察或拼接 |
| `Reduce[R](init, f)` | 按顺序折叠，可改变结果类型 |
| `Collect()` | 返回与列表内部存储解耦的 `[]T` |
| `Stream()` | 创建快照 Stream |

`ArrayListOf` 会直接保存收到的 variadic slice。需要把数据与外部 slice 完全隔离时，可以通过 `Collect` 取出副本后再创建列表。

## LinkedList[T]

`LinkedList[T]` 是双向链表，适合需要频繁在头尾或指定位置插入、删除的场景。它的构造器和返回集合的变换方法都返回指针，方法使用指针接收者：

```go
numbers := lists.LinkedListOf(1, 2, 3)
numbers.AddFirst(0)
numbers.Add(4)
numbers.RemoveAt(2)

got := numbers.Collect()
// got == []int{0, 1, 3, 4}
```

主要方法：

```text
NewLinkedList[T]、LinkedListOf
Add、AddFirst、RemoveFirst、RemoveLast
Insert、RemoveAt、Get
First、Last、Find
Filter、Map、FlatMap、Reduce
Take、Drop、Distinct、SortBy
Peek、Concat、ForEach、Any、All、None
Size()、IsEmpty、Clear
Collect、Stream
```

`LinkedList` 的 `Collect` 按头到尾顺序返回副本。复制一个已经使用中的 `LinkedList` 值会共享底层节点链，使用时应保留原实例，不要复制结构体值。

`LinkedList.Map` 和 `LinkedList.FlatMap` 返回 `*ArrayList[R]`，这样它们和 `ArrayList`、`HashSet` 的类型变换 API 一致，同时允许 `R` 为任意类型。

## HashMap[K, V]

`HashMap[K, V]` 是无序的键值集合，键必须满足 `comparable`，值可以是任意类型。

```go
scores := maps.HashMapOf(
    maps.Entry[string, int]{Key: "alice", Value: 88},
    maps.Entry[string, int]{Key: "bob", Value: 95},
)

scores.Put("carol", 91)
if score, ok := scores.Get("bob"); ok {
    _ = score
}

top := scores.
    FilterValues(func(score int) bool { return score >= 90 }).
    Collect()
```

主要方法：

| 方法 | 行为 |
| --- | --- |
| `NewHashMap[K, V]()` | 创建空 HashMap |
| `HashMapOf(entries...)` | 创建 HashMap，重复键以后面的值为准 |
| `HashMapFromMap(src)` | 从普通 map 拷贝 |
| `Put` / `PutIfAbsent` | 写入或条件写入 |
| `Get` / `GetOrDefault` | 读取 |
| `Remove` | 删除并返回 `(value, ok)` |
| `Contains` / `Size()` / `IsEmpty` / `Clear` | 基础查询和管理 |
| `ForEach` | 遍历键值对 |
| `Keys` / `Values` / `Entries` | 投影为 `ArrayList` |
| `Filter` / `FilterKeys` / `FilterValues` | 筛选并返回新 HashMap |
| `MapValues[R](f)` | 变换 value 类型，key 保持不变 |
| `Concat(other)` | 合并两个 HashMap，后者覆盖同键值 |
| `Collect()` | 返回与内部存储解耦的普通 map |
| `Stream()` | 返回 `Stream[Entry[K, V]]` 快照 |

HashMap 的遍历顺序没有保证，不要依赖 `Keys`、`Values`、`Entries` 或 `Stream` 的顺序。

## HashSet[T]

`HashSet[T]` 是无序且不重复的集合，因此 `T` 必须满足 `comparable`。

```go
left := sets.HashSetOf(1, 2, 3)
right := sets.HashSetOf(3, 4)

union := left.Union(right).Collect()
intersection := left.Intersect(right).Collect()
difference := left.Difference(right).Collect()
```

基础操作包括：

```text
NewHashSet、HashSetOf
Add、Remove、Contains
Size()、IsEmpty、Clear
ForEach、Collect、Stream、Peek
Filter
Reduce、SortBy、MinBy、MaxBy
Union、Concat、Intersect、Difference、SymmetricDifference
IsSubsetOf、IsSupersetOf
Any、All、None、Find
```

HashSet 的映射操作需要区分结果是否仍然要保持集合语义：

```go
set := sets.HashSetOf(1, 2, 3)

// R 可以是任意类型，包括 []int、map 或 func。
listsResult := set.Map(func(v int) []int {
    return []int{v, v * 10}
}).Collect()

// FlatMap 的 mapper 返回 *ArrayList[R]，结果保留所有展开后的元素。
flatResult := set.FlatMap(func(v int) *lists.ArrayList[int] {
    return lists.ArrayListOf(v, v*10)
}).Collect()

// 需要映射后继续去重时，使用 MapSet。
mappedSet := set.MapSet(func(v int) int {
    return v % 2
})

// 需要展开后继续去重时，使用 FlatMapSet。
flatSet := set.FlatMapSet(func(v int) *sets.HashSet[int] {
    return sets.HashSetOf(v, v%2)
})
```

`Map` 和 `FlatMap` 返回 `ArrayList`，所以结果类型不需要可比较；`MapSet` 和 `FlatMapSet` 返回 `HashSet`，要求结果类型满足 `comparable`，并会自动去重。

HashSet 的遍历顺序没有保证。`Find` 在多个值满足条件时返回哪个值是不确定的。

## Stream[T]

`Stream[T]` 是基于 `iter.Seq[T]` 的单次消费、惰性处理流水线。

直接创建 Stream：

```go
stream.From(seq)             // 从 iter.Seq[T] 创建
stream.FromSlice(values)     // 从 []T 创建
stream.Of(value1, value2)    // 从若干值创建
stream.Empty[T]()            // 创建空 Stream
```

```go
import (
    "fmt"

    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/collections/stream"
)

numbers := lists.ArrayListOf(1, 2, 3, 4, 5)

result := numbers.
    Stream().
    Filter(func(n int) bool { return n%2 == 1 }).
    Map(func(n int) string { return fmt.Sprint(n * 10) }).
    Take(2).
    Collect()

// result == []string{"10", "30"}
```

可用的中间操作：

```text
Filter、Map、FlatMap、Peek
Take、Drop、Distinct、Concat
SortBy
```

可用的终止操作：

```text
Collect、Count、First、Last
Any、All、None、Find
Reduce、ForEach
MinBy、MaxBy
```

`Any`、`All`、`First` 和 `Take` 会尽早停止消费上游数据，适合处理大型或潜在无限的数据源。

## 组合示例

集合操作默认立即执行，需要提前停止或延迟处理时再切换到 Stream：

```go
type User struct {
    Name string
    Age  int
}

users := lists.ArrayListOf(
    User{Name: "Alice", Age: 17},
    User{Name: "Bob", Age: 21},
    User{Name: "Carol", Age: 30},
)

adultNames := users.
    Stream().
    Filter(func(u User) bool { return u.Age >= 18 }).
    Map(func(u User) string { return u.Name }).
    Collect()
```

简单的一次性逻辑仍然适合使用普通 `for range`。集合和 Stream 主要用于需要连续变换、筛选、扁平化或提前终止的流程。

## 与其他模块的关系

```text
ArrayList[T] / LinkedList[T] / HashMap[K, V] / HashSet[T]
    同步、具体集合、立即执行

Stream[T]
    同步、惰性、基于 iter.Seq

reactivex
    异步、推送、订阅和取消
```

`collections` 处理已经存在的数据和同步迭代；异步事件、订阅生命周期、背压和多播属于 `reactivex` 的职责。

## 许可

本模块使用 [MIT License](../../LICENSE)。
