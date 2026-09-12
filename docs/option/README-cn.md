# `github.com/qianwj/typed/utils/option`

`Optional[T]`，显式的"可能为 `T`"容器。结构体而非接口，因此方法可以声明自己的类型参数（`Map[R]`、`FlatMap[R]`），是 Go 1.27+ 的泛型方法。

## 包导入

```go
import "github.com/qianwj/typed/utils/option"
```

## 设计动机

为什么不直接用 `(T, bool)`？

- 调用点自我解释；组合子可以链式拼装，不用每一步散落 `if ok { ... }`。
- `Optional` 对值类型（`int` / `string` / `struct{}`）与指针类型一视同仁 —— 它携带独立的 `present` 标志位，而不是依赖"对值做 nil 检查"。`int` 不存在"零值还是缺省"的歧义。
- 不会让 `result.Result[T]` 之类需要 `(T, error)` 形态的组合子和"是否缺席"耦合到一起。

## 构造

| 工厂 | 语义 |
|---|---|
| `Empty[T any]() Optional[T]` | 不携带值；零值 `Optional[T]{}` 等价。 |
| `Of[T any](value T) Optional[T]` | 携带 `value`；**对 typed nil（如 `nil *Foo` 作为指针型 `T`）会 panic**。 |
| `OfNullable[T any](value T) Optional[T]` | 用 Go 的 nil 语义判断：对指针 / map / slice / chan / func / interface 的 nil 返回 `Empty`，否则 `Of(value)`。 |

> 选 `Of` 还是 `OfNullable` 取决于 `T` 的语义：值类型用 `Of`；指针型 / 接口型用 `OfNullable`。

## 状态查询

| 方法 | 说明 |
|---|---|
| `IsPresent() bool` | 携带值。 |
| `IsEmpty() bool` | `!IsPresent()`。 |

## 取值

| 方法 | 缺席时行为 |
|---|---|
| `Get() T` | **panic**（无消息）。 |
| `OrElse(default T) T` | 返回 `default`；`default` 总是被求值。 |
| `OrElseGet(f func() T) T` | 返回 `f()`；仅缺席时调用。 |
| `OrElseThrow(errMsg string) (T, error)` | 缺席时返回 `(零值, errors.New(errMsg))`；出席时返回 `(value, nil)`。**不 panic**，名字沿用 Java `Optional.orElseThrow` 的约定，但遵循 Go 显式 `(T, error)`。 |

## 副作用

| 方法 | 说明 |
|---|---|
| `IfPresent(f func(T))` | 出席时调 `f`；缺席时不动。 |
| `IfPresentOrElse(present func(T), absent func())` | 二选一，恰好一个会被调。 |

## 链式 transform

| 方法 | 说明 |
|---|---|
| `Filter(predicate func(T) bool) Optional[T]` | 出席且 `predicate(value) == true` 时返回自身；否则返回 `Empty`。缺席时不会调用 `predicate`。 |
| `Map[R](f func(T) R) Optional[R]` | 缺席返回 `Empty[R]`，**不会调用** `f`。 |
| `FlatMap[R](f func(T) Optional[R]) Optional[R]` | 缺席时**不会调用** `f`；`f` 自身返回 `Optional[R]` 时直接采用。 |

## 例子

```go
import "github.com/qianwj/typed/utils/option"

o := option.OfNullable(findUser(id))
name, err := o.Map(func(u User) string { return u.Name }).
    OrElseThrow("user not found")
if err != nil {
    return err
}
useName(name)
```

```go
opt := option.Of(42)
opt.IfPresent(func(v int) { fmt.Println(v) }) // 42

empty := option.Empty[int]()
v := empty.OrElseGet(func() int { return compute() }) // 仅此处调用 compute
```

## 与其他包的关系

- `collections` 里所有"可能缺席"的访问器（`ArrayList.Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy`，`Stack` / `Queue` / `Deque` 的 `Pop` / `Peek` / `Front` / `Back`，`Stream.First` / `Last` / `Find`）都返回 `option.Optional[T]`，详见 [`collections`](../collections/README-cn.md)。
- `Result[T]` 通过 `Optional()` 桥到 `Optional[T]`，见 [`result`](../result/README-cn.md)。
