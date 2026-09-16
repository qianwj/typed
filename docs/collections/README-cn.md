# `typed/collections` — Typed

泛型集合与同步数据流。`collections` 分为四类容器和一个惰性流工具：

- `Stack[T]` / `Queue[T]` / `Deque[T]` — 基本线性容器，返回 `option.Option[T]` 表示"无值"。
- `lists.ArrayList[T]` / `lists.LinkedList[T]` — 列表，提供不可变式 transform（`Filter`/`Map`/`Take`/`Drop`/`Concat`/`Distinct`/`SortBy`）。
- `maps.HashMap[K, V]` — 哈希表，提供 `Keys`/`Values`/`Entries`/`Filter*`/`MapValues`/`Concat`。
- `maps.TreeMap[K, V]` — AVL 平衡树映射，支持比较器排序、邻近键查询和区间查询。
- `sets.HashSet[T]` — 哈希集合，提供集合代数（`Union`/`Intersect`/`Difference`/`SymmetricDifference`）以及同型 transform。
- `stream.Stream[T]` — `iter.Seq[T]` 之上的惰性流，可链接 `Filter`/`Map`/`FlatMap`/`Take`/`Drop`/`Distinct`/`Concat`/`SortBy`/`Reduce`/`Count`/`Find` 等终结操作。
- `collections.Range[T]` — 整数半开区间到 `Stream[T]` 的工厂。

所有集合都实现了 `MarshalJSON` / `UnmarshalJSON`，元素类型必须满足 `any` 约束即可，JSON 形式是 Go 风格的数组或对象。

> **Typed** 工具集的一部分。Looking for the English version? See [README.md](./README.md)。

## 目录

- [包导入](#包导入)
- [`iterable.Iterable[T]`](#iterableiterablet)
- [`Stack[T]` / `Queue[T]` / `Deque[T]`](#stackt--queuet--dequet)
- [`lists.ArrayList[T]`](#listsarraylistt)
- [`lists.LinkedList[T]`](#listslinkedlistt)
- [`maps.HashMap[K, V]`](#mapshashmapk-v)
- [`maps.TreeMap[K, V]`](#mapstreemapk-v)
- [`sets.HashSet[T]`](#setshashsett)
- [`stream.Stream[T]`](#streamstreamt)
- [`collections.Range[T]`](#collectionsranget)
- [与其他包的关系](#与其他包的关系)

## 包导入

```go
import (
    "github.com/qianwj/typed/collections"
    "github.com/qianwj/typed/collections/iterable"
    "github.com/qianwj/typed/collections/lists"
    "github.com/qianwj/typed/collections/maps"
    "github.com/qianwj/typed/collections/sets"
    "github.com/qianwj/typed/collections/stream"
)
```

四个子包独立成 `go.mod`，按需引入。

---

## `iterable.Iterable[T]`

`collections/iterable` 包定义共享的遍历接口：

```go
type Iterable[T any] interface {
    ForEach(visit func(T))
}
```

`ArrayList`、`LinkedList`、`HashSet` 和 `Stream` 通过已有的 `ForEach` 方法自动满足接口。转换构造函数接受这些类型，也接受自定义实现，具体集合之间无需互相依赖。

| 构造函数 | 行为 |
|---|---|
| `lists.ArrayListFrom[T](source iterable.Iterable[T])` | 创建新的 `*ArrayList[T]`，保留源遍历顺序和重复值。 |
| `lists.LinkedListFrom[T](source iterable.Iterable[T])` | 创建新的 `*LinkedList[T]`，保留源遍历顺序和重复值。 |
| `sets.HashSetFrom[T comparable](source iterable.Iterable[T])` | 创建新的 `*HashSet[T]`，自动去重。 |

每个构造函数只遍历源一次，结果具有独立存储，元素本身为浅拷贝。set 的遍历顺序未指定，因此转为 list 后也不保证顺序。传入 Stream 会消费该流。

```go
list := lists.ArrayListOf(3, 1, 3, 2)
set := sets.HashSetFrom(list)
array := lists.ArrayListFrom(set)
linked := lists.LinkedListFrom(array)
```

---

## `Stack[T]` / `Queue[T]` / `Deque[T]`

三者都是结构体指针构造，返回 `*Stack[T]` 等。所有"取一个元素"的操作都返回 `option.Option[T]`，而不是 `(T, bool)`，以便和 `Stream` / `Result` 链路自然拼装。

```go
s := collections.NewStack[int]()
s.Push(1); s.Push(2)
v, ok := s.Peek().Get()        // 2, true
top, _ := s.Pop().Get()        // 2
_, present := s.Pop().Get()    // 1, true
s.Pop().IsEmpty()              // false（还有 1）
s.Pop().IsEmpty()              // true
s.Pop()                        // option.Empty[int]()，不 panic
```

| 类型 | 构造 | 关键方法 |
|---|---|---|
| `Stack[T]` | `NewStack[T]()` | `Push(v) / Pop() Option[T] / Peek() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |
| `Queue[T]` | `NewQueue[T]()` | `Push(v) / Pop() Option[T] / Peek() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |
| `Deque[T]` | `NewDeque[T]()` | `PushFront(v) / PushBack(v) / PopFront() Option[T] / PopBack() Option[T] / Front() Option[T] / Back() Option[T] / Size() / IsEmpty() / Clear() / MarshalJSON / UnmarshalJSON` |

`Pop*` / `Peek*` 在容器为空时返回 `option.Empty[T]()`，不会 panic，也不会修改容器（`Peek`/`Front`/`Back`），或者按对应规则修改（`Pop*`）。

---

## `lists.ArrayList[T]`

`NewArrayList[T any]() *ArrayList[T]` 构造一个空表。`ArrayList` 是单段连续存储的 `[]T`，头部有一个 `head` 偏移；`Take`/`Drop` 只移动 `head` 而不复制元素（`O(1)`），`Collect` 之后才真正复制到新切片。

### 读写

| 方法 | 说明 |
|---|---|
| `Add(v)` / `AddFirst(v)` | 追加到末尾 / 插入到头部。 |
| `Insert(i, v)` | 在下标 `i` 插入 `v`；`i < 0` 或 `i > Size()` 会 panic。 |
| `Get(i) Option[T]` | 下标越界返回 `Empty`，不 panic。 |
| `First() Option[T]` / `Last() Option[T]` | 空表返回 `Empty`。 |
| `RemoveAt(i) T` | 删除下标 `i` 的元素并返回它；越界 panic。 |
| `RemoveFirst() Option[T]` / `RemoveLast() Option[T]` | 空表返回 `Empty`。 |
| `Size() / IsEmpty() / Clear() / Collect() []T` | 基础量与导出。 |
| `Stream() stream.Stream[T]` | 转成惰性 `Stream[T]`（不复制数据）。 |
| `MarshalJSON / UnmarshalJSON` | JSON 数组。 |

### 不可变 transform（返回新表，原表不变）

| 方法 | 说明 |
|---|---|
| `Filter(p func(T) bool) *ArrayList[T]` | 保留满足 `p` 的元素。 |
| `Map[R](f func(T) R) *ArrayList[R]` | 元素类型 `T → R`。 |
| `FlatMap[R](f func(T) *ArrayList[R]) *ArrayList[R]` | `f` 返回的子表按顺序拼接。 |
| `Take(n) / Drop(n)` | 前缀 / 后缀，复杂度 `O(1)`。 |
| `Distinct(eq func(T, T) bool)` | 按 `eq` 去重。 |
| `Concat(other)` | 末尾拼接。 |
| `Peek(visit func(T))` | 不改表地遍历；返回值仍是原表，用于链式调试。 |
| `SortBy(less func(x, y T) int) *ArrayList[T]` | 按 `less` 排序，返回新表。 |
| `MinBy(less func(x, y T) int) option.Option[T]` | 空表返回 `option.Empty[T]()`；否则返回最小元素（出席）。 |
| `MaxBy(less func(x, y T) int) option.Option[T]` | 同上，返回最大元素。 |

### 谓词

`Any(p) / All(p) / None(p) / Find(p) Option[T]` — 短路的全称 / 存在量词；`Find` 返回首个匹配元素。

### 例子

```go
xs := lists.NewArrayList[int]()
xs.Add(3).Add(1).Add(4).Add(1).Add(5) // 注：Add 返回 *ArrayList 是设计待定；以源码为准
_ = xs
```

> 上面这行 `Add` 链式调用仅为示意；当前实现里 `Add` 不返回 `*ArrayList[T]`，如需链式请改用 `Peek` 或 `stream.Stream`。

```go
even := lists.NewArrayList[int]()
even.Add(2); even.Add(4); even.Add(6)
sorted := even.SortBy(func(a, b int) int { return a - b })
max := sorted.MaxBy(func(a, b int) int { return a - b }).OrElse(0) // 6
```

---

## `lists.LinkedList[T]`

双向链表 + 哨兵节点。`NewLinkedList[T any]() *LinkedList[T]` 构造。

### 读写

| 方法 | 说明 |
|---|---|
| `Add(v)` / `AddFirst(v)` | 追加到末尾 / 插入到头部。 |
| `Insert(i, v)` | 在下标 `i` 插入；越界 panic。 |
| `Get(i) Option[T]` | 越界返回 `Empty`。 |
| `First() / Last() Option[T]` | 空表返回 `Empty`。 |
| `RemoveAt(i) T` | 越界 panic。 |
| `RemoveFirst() / RemoveLast() Option[T]` | 空表返回 `Empty`。 |
| `Size() / IsEmpty() / Clear() / Collect() []T` | 基础量与导出。 |
| `Stream() stream.Stream[T]` | 惰性 `Stream[T]`。 |
| `MarshalJSON / UnmarshalJSON` | JSON 数组。 |

### 不可变 transform

`Filter(p) *LinkedList[T]` / `Map[R](f) *ArrayList[R]` / `FlatMap[R](f) *ArrayList[R]` / `Take(n) / Drop(n) / Distinct(eq) / Concat(other) / Peek(visit) / SortBy(less) *LinkedList[T]`。

注意 `Map` / `FlatMap` 从链表到数组是顺序遍历，复杂度 `O(n)`，但结果类型是 `*ArrayList[R]`，因为下游多半要按下标或切片处理。

`MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` — 空表返回 `Empty`，否则返回极值。

### 谓词

`Any / All / None / Find(p) Option[T]`，以及 `Reduce[U](init U, f func(U, T) U) U`。

---

## `maps.HashMap[K, V]`

`K` 必须 `comparable`，`V` 任意。`NewHashMap[K, V]() *HashMap[K, V]` 构造。

### 基础

| 方法 | 说明 |
|---|---|
| `Put(k, v) V` | 设值并返回旧值（不存在时返回零值）。 |
| `PutIfAbsent(k, v)` | 仅在不存在时设值。 |
| `Get(k) (V, bool)` | 取值；不存在返回零值 + `false`。 |
| `GetOrDefault(k, defaultV) V` | 取值或回退。 |
| `Remove(k) (V, bool)` | 删除并返回旧值。 |
| `Contains(k) bool` | 键存在性。 |
| `Size() / IsEmpty() / Clear()` | 基础量。 |

### 视图

| 方法 | 说明 |
|---|---|
| `Keys() *ArrayList[K]` | 键集合（无序）。 |
| `Values() *ArrayList[V]` | 值集合（无序）。 |
| `Entries() *ArrayList[Entry[K, V]]` | `Entry` 是 `{K, V}` 结构体，`MarshalJSON` 产出 `{"k": v}`。 |
| `ForEach(func(K, V))` | 不改表遍历。 |
| `Stream() stream.MapStream[K, V]` | 基于键值快照的惰性流；`Collect()` 返回 `map[K]V`，遍历顺序不确定。 |
| `Collect() map[K]V` | 导出为 Go 内置 `map`。 |
| `MarshalJSON / UnmarshalJSON` | `{"k": v, ...}` 形式。 |

### 不可变 transform

`Filter(p) / FilterKeys(p) / FilterValues(p) *HashMap[K, V]` — 返回新表。
`MapValues[R](f func(K, V) R) *HashMap[K, R]` — 值类型变换。
`Concat(other) *HashMap[K, V]` — 同键时 `other` 覆盖当前表。

---

## `maps.TreeMap[K, V]`

基于 AVL 平衡树的有序映射。用 `NewTreeMap[K comparable, V any](compare func(K, K) int)` 或 `TreeMapOf(compare, entries...)` 构造。自然顺序可传标准库 `cmp` 包的 `cmp.Compare[K]`，交换比较参数即可降序。

比较器必须定义一致的全序。比较结果为零就视为同一个键，即使 Go 的 `==` 判断不同；更新时保留最初存入的键。键存入后，其比较结果必须保持稳定。nil 比较器会 panic，TreeMap 零值不可直接使用。

```go
m := maps.NewTreeMap[int, string](cmp.Compare[int])
m.Put(30, "thirty")
m.Put(10, "ten")
m.Put(20, "twenty")

m.Keys().Collect()      // [10, 20, 30]
m.Floor(25).Get()       // Entry{Key: 20, Value: "twenty"}
m.Higher(20).Get()      // Entry{Key: 30, Value: "thirty"}
m.Range(10, 30).Collect() // entries for 10 and 20
```

| 方法 | 行为 |
|---|---|
| `Put(k, v) V / PutIfAbsent(k, v)` | 插入或替换，返回约定与 HashMap 一致；O(log n)。 |
| `Get(k) (V, bool) / GetOrDefault(k, fallback) / Contains(k)` | 按比较器查询；O(log n)。 |
| `Remove(k) (V, bool)` | 删除并返回旧值及是否存在；O(log n)。 |
| `Size() / IsEmpty() / Clear()` | 大小、是否为空、清空；清空后保留比较器。 |
| `First() / Last()` | 最小/最大键对应的 `adt.Option[Entry[K, V]]`；O(log n)。 |
| `Floor(k) / Ceiling(k)` | 不大于/不小于 k 的最近条目；返回 `adt.Option[Entry[K, V]]`，O(log n)。 |
| `Lower(k) / Higher(k)` | 严格前驱/后继；返回 `adt.Option[Entry[K, V]]`，O(log n)。 |
| `Range(from, to)` | 比较器顺序下 `[from, to)` 的有序快照流；k 项结果耗时 O(log n + k)，相等或反向边界返回空流。 |
| `ForEach(func(K, V))` | 按比较器顺序遍历；回调不得修改树结构。 |
| `Keys() / Values() / Entries()` | 按键顺序返回独立 ArrayList。 |
| `Stream() stream.MapStream[K, V]` | 按比较器顺序提供键值快照；`Collect()` 返回不保留顺序的 `map[K]V`。 |
| `Collect() map[K]V` | 独立的原生 map，不保留顺序。 |
| `Filter(p) / MapValues[R](f)` | 保留比较器，返回新 TreeMap；O(n log n)。 |
| `MarshalJSON / UnmarshalJSON` | JSON 对象，不保证成员按比较器排序。 |

没有匹配项时，邻近键查询返回空 Option。快照及转换结果具有独立容器存储，元素本身为浅拷贝。TreeMap 不支持无同步的并发修改。

JSON 解码前必须先用比较器构造 TreeMap。解码保留比较器，成功时替换条目，`null` 清空，输入非法则保留原数据。若对象中存在 Go 键不同但比较器判断相等的键，会合并且不保证哪一个值胜出。

---

## `sets.HashSet[T]`

`T` 必须 `comparable`。`NewHashSet[T any]() *HashSet[T]` 构造。

### 基础

`Add(v) / Remove(v) / Contains(v) bool / Size() / IsEmpty() / Clear()` 与 Go `map[T]struct{}` 语义一致；`Add` 对已存在元素无副作用。

### 视图与流

`ForEach(func(T)) / Collect() []T / Stream() stream.Stream[T] / Peek(visit) *HashSet[T] / MarshalJSON / UnmarshalJSON`（JSON 数组，去重后输出）。

### 不可变 transform

| 方法 | 说明 |
|---|---|
| `Filter(p) *HashSet[T]` | 保谓词元素。 |
| `Map[R](f) *ArrayList[R]` | 元素类型变换（结果不再去重，类型为列表）。 |
| `FlatMap[R](f) *ArrayList[R]` | 同上但子结果也是列表。 |
| `MapSet[R comparable](f) *HashSet[R]` | 元素类型变换并去重。 |
| `FlatMapSet[R comparable](f) *HashSet[R]` | 子结果为集合，平铺去重。 |
| `Concat(other) *HashSet[T]` | 并集（去重），原表不变。 |
| `Reduce[R](init R, f func(R, T) R) R` | 折叠。 |
| `SortBy(less) *ArrayList[T]` | 排序后导出为列表。 |
| `MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` | 空集返回 `option.Empty[T]()`。 |

### 集合代数

`Concat / Union / Intersect / Difference / SymmetricDifference` 都返回新的 `*HashSet[T]`，原表不变：

| 方法 | 等价集合记号 |
|---|---|
| `Concat(other)` | `s ∪ other` |
| `Union(other)` | `s ∪ other`（与 `Concat` 行为一致） |
| `Intersect(other)` | `s ∩ other` |
| `Difference(other)` | `s \ other` |
| `SymmetricDifference(other)` | `(s ∪ other) \ (s ∩ other)` |
| `IsSubsetOf(other) / IsSupersetOf(other)` | 集合包含关系。 |

### 谓词

`Any / All / None / Find(p) option.Option[T]`。

---

## `stream.Stream[T]`

惰性同步流。`Stream[T]` 是 `iter.Seq[T]` 之上的薄包装，所有 transform 都返回新 `Stream`，**不**在 transform 中消费原流。`Collect` / `Count` / `First` / `Last` / `Any` / `All` / `None` / `Find` / `Reduce` / `ForEach` / `MinBy` / `MaxBy` / `SortBy` 是终结操作。

### 构造

| 函数 | 说明 |
|---|---|
| `From[T](seq iter.Seq[T]) Stream[T]` | 直接包装 `iter.Seq[T]`。 |
| `FromSlice[T](s []T) Stream[T]` | 在 `s` 上做惰性遍历（`s` 会被持有但不会复制）。 |
| `Of[T](values ...T) Stream[T]` | 变参 → `Stream[T]`，`FromSlice` 的语法糖。 |
| `Empty[T]() Stream[T]` | 立刻结束的流。 |

> 任何集合（`ArrayList` / `LinkedList` / `HashSet` / `HashMap.Entries`）都可以通过自身的 `Stream()` 方法直接转 `Stream[T]`，这是推荐的入口。

### 链式 transform

`Filter(p) / Map[R](f) / FlatMap[R](f) / Peek(visit) / Take(n) / Drop(n) / Distinct(eq) / Concat(other) / SortBy(less)` — 全部返回 `Stream[...]`，可继续链式。

`Associate[K, V](f func(T) (K, V)) MapStream[K, V]` 惰性地产生键值对。`MapStream.Filter(func(K, V) bool)` 和 `MapValues[R](func(K, V) R)` 保持键值流类型；`Collect()` 返回原生 `map[K]V`，空流也返回非 nil 的空 map。重复键在中间转换中保留，最终到达 `Collect` 的后值覆盖前值。`Map[R](func(K, V) R)` 显式投影为普通 `Stream[R]`，其 `Collect()` 返回切片。

```go
var result map[int]string = stream.Of("a", "bb", "ccc", "dd").
    Associate(func(s string) (int, string) { return len(s), s }).
    Filter(func(k int, _ string) bool { return k == 2 }).
    Collect() // map[int]string{2: "dd"}
```

### 终结

| 方法 | 返回 |
|---|---|
| `Collect() []T` | 物化为切片。 |
| `Count() int` | 计数（短路不过滤后元素）。 |
| `First() / Last() option.Option[T]` | 第一个 / 最后一个元素。 |
| `Any(p) / All(p) / None(p) bool` | 短路存在 / 全称量词。 |
| `Find(p) option.Option[T]` | 首个匹配元素。 |
| `Reduce(init T, f func(acc, v T) T) T` | 折叠。 |
| `ForEach(visit func(T))` | 纯消费。 |
| `MinBy(less) option.Option[T]` / `MaxBy(less) option.Option[T]` | 极值；空流返回 `option.Empty[T]()`。 |
| `SortBy(less) Stream[T]` | 终结式：先物化排序再返回新流。 |

> `SortBy` 在 `Stream` 上是终结操作（会一次性遍历），与 `ArrayList` / `LinkedList` / `HashSet` 上"transform 风格"的 `SortBy` 不一样。

### 例子

```go
out := stream.Of(1, 2, 3, 4, 5).
    Filter(func(v int) bool { return v%2 == 1 }).
    Map(func(v int) int { return v * v }).
    Collect() // [1, 9, 25]

first, ok := stream.Of[int]().First().Get() // 0, false
```

---

## `collections.Range[T]`

`Range[T constraints.Integer](start, end T) stream.Stream[T]` —— 整数半开区间 `[start, end)` 转 `Stream[T]`。

- 当 `start >= end` 时是空流；不会 panic。
- 流是惰性的：迭代 `Range(0, 1_000_000_000)` 只占常数级内存。
- 类型约束是 `constraints.Integer`（`int` / `uint` / `int32` / …），不接受浮点或字符串。

```go
collections.Range[int](1, 4).Collect() // [1, 2, 3]
collections.Range[int](0, 5).
    Map(func(i int) int { return i * i }).
    Take(3).
    Collect() // [0, 1, 4]
```

## 与其他包的关系

- 返回值大量用 `option.Option[T]`，见 [`adt`](../adt/README-cn.md)。
- 错误流请用 `result.Result[T]`，见 [`adt`](../adt/README-cn.md)。
- 异步 / 多订阅请用 `reactivex.Observable[T]`，见 [`reactivex`](../reactivex/README-cn.md)。本包的 `Stream` 是同步单次消费模型，两者不互通。
