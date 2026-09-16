# `typed/adt` — Typed

[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=adt)](https://codecov.io/gh/qianwj/typed)

强类型的代数数据类型 —— 显式建模「存在 / 缺失」「成功 / 失败」「这一变体 / 那一变体」的值类型。本包提供三种类型：

- [`Option[T]`](#optionalt) —— 单一类型的零或一个。`(T, bool)` 的类型化替代。
- [`Result[T]`](#resultt) —— 携带 `T` 的成功，或携带非 nil `error` 的失败。`(T, error)` 的类型化替代。
- [`Either[L, R]`](#eitherl-r) —— 两种值中恰好取一。前两种类型的通用 tagged-union 原语。

> **Typed** 工具集的一部分。Looking for the English version? See [README.md](./README.md)。

## 目录

- [包导入](#包导入)
- [为什么放在同一个包](#为什么放在同一个包)
- [`Option[T]`](#optionalt)
- [`Result[T]`](#resultt)
- [`Either[L, R]`](#eitherl-r)
- [`Cast[T]`](#castt)
- [如何选择](#如何选择)
- [相关阅读](#相关阅读)

## 包导入

```go
import "github.com/qianwj/typed/adt"
```

三种类型共享同一个 import —— `adtOption`、`adtResult`、`adtEither`。

## 为什么放在同一个包

三种类型共享同一套词汇（present / absent / success / failure / left / right），并以惯用方式相互引用：

- `Result.Option` 返回 `Option[T]` —— 成功转 present，失败转 absent。
- `Either.Left` / `Either.Right` 返回 `Option[L]` / `Option[R]` —— tagged union 的安全访问器。

放进同一个包让这些关系在调用点显而易见，省去多 import 的摩擦。独立的 module（`typed/adt` 而非 `typed/utils/adt`）表明这些是值类型，不是 `objects.IsNil` 或 JSON codec 那种意义上的工具——消费者可以只依赖值类型，不用连带拉入 `utils/*` 其他代码。

## `Option[T]`

`Option[T]` 是一个容器，可能持有也可能不持有类型 `T` 的值。`Option` 要么 present（携带 `T`），要么 absent（无值）。`Option` 的零值就是 absent，等价于 `Empty[T]()`。

`Option` 不是接口，所以方法可以声明自己的类型参数（`Map[R]`、`FlatMap[R]`）—— Go 1.27+ 的泛型方法特性。

### 为什么不直接用 `(T, bool)`？

- 调用点自我解释；组合子可以链式拼装，不用每一步散落 `if ok { ... }`。
- `Option` 对值类型（`int` / `string` / `struct{}`）与指针类型一视同仁 —— 它携带独立的 `present` 标志位，而不是依赖"对值做 nil 检查"。`int` 不存在"零值还是缺省"的歧义。
- 不会让 `Result[T]` 之类需要 `(T, error)` 形态的组合子和"是否缺席"耦合到一起。

### 构造

| 工厂 | 语义 |
|---|---|
| `Empty[T any]() Option[T]` | 不携带值；零值 `Option[T]{}` 等价。 |
| `Of[T any](value T) Option[T]` | 携带 `value`；**对 typed nil（如 `nil *Foo` 作为指针型 `T`）会 panic**。 |
| `OfNullable[T any](value T) Option[T]` | 用 Go 的 nil 语义判断：对指针 / map / slice / chan / func / interface 的 nil 返回 `Empty`，否则 `Of(value)`。 |

> 选 `Of` 还是 `OfNullable` 取决于 `T` 的语义：值类型用 `Of`；指针型 / 接口型用 `OfNullable`。

### 状态查询

| 方法 | 说明 |
|---|---|
| `IsPresent() bool` | 携带值。 |
| `IsEmpty() bool` | `!IsPresent()`。 |

### 取值

| 方法 | 缺席时行为 |
|---|---|
| `Get() T` | **panic**（「Get on empty Option」）。 |
| `OrElse(default T) T` | 返回 `default`；`default` 总是被求值。 |
| `OrElseGet(f func() T) T` | 返回 `f()`；仅缺席时调用。 |
| `OrElseThrow(errMsg string) (T, error)` | 缺席时返回 `(零值, errors.New(errMsg))`；出席时返回 `(value, nil)`。**不 panic**，名字沿用 Java `Option.orElseThrow` 的约定，但遵循 Go 显式 `(T, error)`。 |

### 副作用

| 方法 | 说明 |
|---|---|
| `IfPresent(f func(T))` | 出席时调 `f`；缺席时不动。 |
| `IfPresentOrElse(present func(T), absent func())` | 二选一，恰好一个会被调。 |

### 链式 transform

| 方法 | 说明 |
|---|---|
| `Filter(predicate func(T) bool) Option[T]` | 出席且 `predicate(value) == true` 时返回自身；否则返回 `Empty`。缺席时不会调用 `predicate`。 |
| `Map[R](f func(T) R) Option[R]` | 缺席返回 `Empty[R]`，**不会调用** `f`。 |
| `FlatMap[R](f func(T) Option[R]) Option[R]` | 缺席时**不会调用** `f`；`f` 自身返回 `Option[R]` 时直接采用。 |

### 例子

```go
import "github.com/qianwj/typed/adt"

o := adt.OfNullable(findUser(id))
name, err := o.Map(func(u User) string { return u.Name }).
    OrElseThrow("user not found")
if err != nil {
    return err
}
useName(name)
```

```go
opt := adt.Of(42)
opt.IfPresent(func(v int) { fmt.Println(v) }) // 42

empty := adt.Empty[int]()
v := empty.OrElseGet(func() int { return compute() }) // 仅此处调用 compute
```

## `Result[T]`

`Result[T]` 是可能失败的操作结果。`Result` 要么是携带 `T` 的成功，要么是携带非 nil error 的失败。

`Result` 的零值是「携带 nil error 的失败」；调用方应该始终通过 `Success` / `Failure` 构造并通过方法检查。

### 为什么不直接用 `(T, error)`？

- 显式 `Result` 让调用点自我解释，组合子可以链式拼接，不用每一步散落 `if err != nil`。
- 形态刻意比泛型 `Result[T, E]` 简单：错误通道始终是标准 `error`，保持和普通 Go 代码的互操作不需要自定义 error 装箱。

### 构造

| 工厂 | 语义 |
|---|---|
| `Success[T any](value T) Result[T]` | 携带 `value`。 |
| `Failure[T any](err error) Result[T]` | 携带 `err`。**对 nil error 会 panic**——nil error 必须走 `Success`。 |
| `Wrap[T any](value T, err error) Result[T]` | Go 惯用 `(T, error)` 的桥。`err == nil` 时成功，否则失败。`value` 两种情况都会存，但只有成功时能通过方法读到。 |

`Wrap` 的两个适用场景：`(T, error)` API 边界的末端，调用方想继续在 `Result` 流水线（Map / FlatMap / OrElse / ...）里跑；以及从 `Result` 链回到普通 Go 错误处理的桥。

### 状态查询

| 方法 | 说明 |
|---|---|
| `IsSuccess() bool` | `err == nil`。 |
| `IsFailure() bool` | `err != nil`。 |

### 取值

| 方法 | 失败时行为 |
|---|---|
| `Value() T` | **panic**，panic 值就是底层 error。 |
| `Unwrap() (T, error)` | 返回 `(零值, err)`，不 panic，把 error 交给调用方处理。是 `Wrap` 的反方向桥。 |
| `Error() error` | 返回 error，成功时为 nil。 |
| `OrElse(default T) T` | 返回 `default`；`default` 总是被求值。 |
| `OrElseGet(f func() T) T` | 返回 `f()`，仅失败时调用。 |
| `Recover(f func(error) T) T` | 失败时调 `f(err)` 并返回其结果；成功时原值不动。错误感知的 fallback。 |

### 桥接

| 方法 | 说明 |
|---|---|
| `Option() Option[T]` | 成功 → present `Option`；失败 → absent `Option`，error 被丢掉。仅在调用方已经决定可以丢弃 error 时使用。 |

### 链式 transform

| 方法 | 说明 |
|---|---|
| `Map[R](f func(T) R) Result[R]` | 失败时返回带相同 error 的 `Result[R]`，**不会调用** `f`。 |
| `FlatMap[R](f func(T) Result[R]) Result[R]` | 失败时不会调用 `f`。调用时直接采用 `f` 返回的 `Result[R]`。 |
| `MapError(f func(error) error) Result[T]` | 失败时调 `f(err)`；成功时原值不动。把底层 error 包成高层描述再上抛很有用。 |

### 例子

```go
import "github.com/qianwj/typed/adt"

r := adt.Wrap(loadProfile(id))
// r 是 Result[Profile]，err 为 nil 即成功，否则失败。

name := r.Map(func(p Profile) string { return p.Name }).
    Recover(func(err error) string {
        if errors.Is(err, ErrNotFound) {
            return "anonymous"
        }
        return "" // 任何其他 error → 空名
    })
```

## `Either[L, R]`

`Either[L, R]` 是一个 tagged-union 值类型，恰好持有两种值中的一种：类型 `L` 的 `Left` 或类型 `R` 的 `Right`。约定俗成的读法是**「Left = 失败，Right = 成功」**——比如 `Right[error, T]` 是 Go 惯用 `(T, error)` 元组的类型化等价物。

### 为什么用 tagged union？

Go 的 `(T, error)` 在单点调用处很顺手，但一旦值要在多个中间函数之间流转就显得笨拙：每一步都得传递 `error` 变量并重复 nil 检查。`Either` 把两侧都变成一等公民，因此一个可能以两种方式失败的函数可以返回 `Either[errA, errB]`，一个可能产生两种值类型的函数可以返回 `Either[T, U]`。调用方用 `IsLeft` / `IsRight` 或 `Fold` 组合子选对应分支。

### 构造

| 工厂 | 语义 |
|---|---|
| `Left[L, R any](l L) Either[L, R]` | 携带 `l`。`Either` 的零值是持有 `L` 零值的 `Left`。 |
| `Right[L, R any](r R) Either[L, R]` | 携带 `r`。 |

### 谓词

| 方法 | 说明 |
|---|---|
| `IsLeft() bool` | 携带 `Left`。 |
| `IsRight() bool` | 携带 `Right`。 |

### 安全访问器

| 方法 | 说明 |
|---|---|
| `Left() Option[L]` | `IsLeft()` 时 present，`IsRight()` 时 absent。 |
| `Right() Option[R]` | `IsRight()` 时 present，`IsLeft()` 时 absent。 |
| `LeftOrZero() L` | 返回 `Left` 值，`Right` 时返回 `L` 零值。读错一侧无所谓时用这个。 |
| `RightOrZero() R` | `LeftOrZero` 的镜像。 |

安全访问器返回 `Option`，这样调用方可以和别处已有的同套组合子链式拼装。

### 组合子

`MapLeft` / `MapRight` 和 `FlatMapLeft` / `FlatMapRight` 显式选择要操作的分支；遇到另一侧时均原样透传，不调用回调。`FlatMapLeft` 保持 `R` 不变，允许改变 `L`；`FlatMapRight` 保持 `L` 不变，允许改变 `R`。两者的回调均可返回 `Left` 或 `Right`，因此链式调用可以切换分支。用 `Fold` 汇合两个分支。

| 方法 | 说明 |
|---|---|
| `MapLeft[L2](f func(L) L2) Either[L2, R]` | 只在 `Left` 上调 `f`；`Right` 分支原样透传。 |
| `MapRight[R2](f func(R) R2) Either[L, R2]` | 只在 `Right` 上调 `f`；`Left` 原样透传。 |
| `FlatMapLeft[L2](f func(L) Either[L2, R]) Either[L2, R]` | `Left` 时直接返回 `f` 产生的 `Either`；`Right` 原样透传，不调用 `f`。 |
| `FlatMapRight[R2](f func(R) Either[L, R2]) Either[L, R2]` | `Right` 时直接返回 `f` 产生的 `Either`；`Left` 原样透传，不调用 `f`。 |
| `Fold[T](onLeft func(L) T, onRight func(R) T) T` | 二选一，恰好一个回调会被调用并返回其结果。「挑对分支」的规范模式。 |

### 例子

```go
import "github.com/qianwj/typed/adt"

result := adt.Right[error, int](42)
val := result.Right().Get() // val == 42
```

```go
result := divide(10, 2)
msg := result.Fold(
    func(e error) string { return "err: " + e.Error() },
    func(n int) string { return fmt.Sprintf("val: %d", n) },
)
```

```go
result := adt.Left[string, int]("missing").
    FlatMapLeft(func(problem string) adt.Either[error, int] {
        if problem == "missing" {
            return adt.Right[error, int](21)
        }
        return adt.Left[error, int](errors.New(problem))
    }).
    FlatMapRight(func(n int) adt.Either[error, string] {
        return adt.Right[error, string](fmt.Sprintf("value: %d", n*2))
    })
fmt.Println(result.Right().Get()) // value: 42
```

## `Cast[T]`

`Cast[T any](v any) Result[T]` 执行类型断言，类型不匹配时不会 panic。目标为具体类型时，动态类型必须与其相同；目标为接口时，动态类型必须实现该接口。它不会执行 `int32` 到 `int64` 这样的类型转换。

```go
adt.Cast[int](42).OrElse(0)           // 42
adt.Cast[int64](int32(42)).IsFailure() // true
adt.Cast[int](nil).IsFailure()        // true
adt.Cast[*int]((*int)(nil)).IsSuccess() // true：类型匹配的 typed nil
```

断言失败时，错误包含实际类型与目标类型。nil 接口断言失败；类型匹配的 typed nil 断言成功并保留原值，这与将 typed nil 视为缺失的 `OfNullable` 不同。返回结果可继续使用 `Map`、`FlatMap`、`OrElse` 或 `Unwrap`。

## 如何选择

| 场景 | 用 |
| --- | --- |
| 可能缺值，没有附带失败细节 | `Option[T]` |
| 成功返回值，或失败返回 error | `Result[T]` |
| 两种不同的成功 / 失败类别，或两种值类型 | `Either[L, R]` |

`Option` 是零或一。`Result` 是成功或失败（失败总带标准 `error`）。`Either` 是通用的——两个分支可以是任何类型，包括两种不同的 error 类型，如果想在标准 `error` 之上做错误分类。

实际使用：

- `Map[K]V.Get(key)` 可能 miss → `Option[V]`。
- 数据库查询返回一行或错误 → `Result[Row]`。
- 可能返回 `T` 或配置错误的函数 → `Either[ConfigError, T]`。

## 相关阅读

- [`collections`](../collections/README-cn.md) 里所有「可能缺席」的访问器都返回 `adtOption[T]`：`ArrayList.Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy`，`Stack` / `Queue` / `Deque` 的 `Pop` / `Peek` / `Front` / `Back`，`Stream.First` / `Last` / `Find`。
- [`utils/json`](../utils/json/README-cn.md) 用 `adtResult[T]` 作为编解码返回值。
- [`reactivex.Single[T]`](../reactivex/README-cn.md#singlet) / [`reactivex.Maybe[T]`](../reactivex/README-cn.md#maybet) —— 终止事件是 `Either` 形态的 reactive 容器。
