# `github.com/qianwj/typed/control`

控制流辅助。两个子包：

- [`control`](#control) — `Repeat` / `RepeatE`，简洁的"做 N 次"循环原语。
- [`control/match`](#controlmatch) — 类型安全、首个匹配获胜的模式匹配（`Pattern[T]` + 链式 `Case` / `Type`）。

需要 Go 1.27+（链式方法带自己的类型参数）。

## 包导入

```go
import (
    "github.com/qianwj/typed/control"
    "github.com/qianwj/typed/control/match"
)
```

---

## `control`

```go
func Repeat(times int, f func())
func RepeatE(times int, f func() error) (int, error)
```

### `Repeat(times, f)`

调用 `f` 恰好 `times` 次，不返回值。`times <= 0` 是空操作（body 不会执行，`f` 也不会被调用）。等价于 `for i := 0; i < times; i++ { f() }`，但更显式。

```go
control.Repeat(3, func() { fmt.Println("tick") })
```

### `RepeatE(times, f) (int, error)`

错误感知的变体。`f` 返回非 nil 错误时立即停掉循环，返回 `(i, err)` —— 其中 `i` 是**已经成功完成**的迭代次数（错误是第 `i+1` 次调用的结果）。

- 全部成功 → `(times, nil)`。
- 第 `i+1` 次出错（0 索引 `i`）→ `(i, err)`，循环停在该次。
- `times <= 0` → `(0, nil)`（body 不进）；`times < 0` 时返回 `(times, nil)`，这是为了和入参对称而非有效计数。

适用场景：重试循环、批处理"遇到第一个失败就停"。

```go
n, err := control.RepeatE(5, func() error {
    return doOne()
})
if err != nil {
    log.Printf("after %d successful attempts: %v", n, err)
}
```

### 与 `collections.Range` 的关系

`collections.Range(0, n).ForEach(f)` 会先构造一个 `Stream`，再迭代它。`control.Repeat(n, f)` 是省掉 `Stream` 分配的等价命令式写法。选 `Repeat` 用于一次性副作用；选 `Stream` 用于"链式 pipeline 的起点"。

---

## `control/match`

首个匹配获胜的链式模式匹配。两种模式：

1. **值匹配** `Match(v)` / `Value(v)` —— 按 `Pattern[T]` 测试输入值。
2. **类型匹配** `Type(v)` —— 测试一个 `any` 的动态类型 / 类型守卫。

API 是显式的：所有 `Case` 必须返回相同的 `R`，按声明顺序求值，第一个匹配的 `Case` 之后的所有 `Pattern` / handler 都不会被执行。

### `Pattern[T]`

```go
type Pattern[T any] interface {
    Match(value T) bool
}
type PatternFunc[T any] func(T) bool
func (p PatternFunc[T]) Match(value T) bool
```

### 工厂函数

| 工厂 | 语义 |
|---|---|
| `Predicate[T](p func(T) bool) Pattern[T]` | 任意谓词。 |
| `Any[T]() Pattern[T]` | 匹配任何值。 |
| `Eq[T](want T) Pattern[T]` | 按 `objects.Equals` 深度比较；带 `Equal(T) bool` 的类型会被特化调用。 |
| `In[T comparable](values ...T) Pattern[T]` | 命中给定集合之一。 |
| `Between[T cmp.Ordered](min, max T) Pattern[T]` | 闭区间 `[min, max]`。 |
| `Not[T](pattern Pattern[T]) Pattern[T]` | 否定。 |
| `And[T](patterns ...Pattern[T]) Pattern[T]` | 全配；无参时匹配任何值。 |
| `Or[T](patterns ...Pattern[T]) Pattern[T]` | 至少一配；无参时永假。 |

### 值匹配入口

```go
func Value[T any](value T) Matcher[T]   // 入口
func Match[T any](value T) Matcher[T]   // Value 的别名
```

`Matcher[T]` / `Chain[T, R]` 上的方法：

| 方法 | 说明 |
|---|---|
| `Case[R](pattern Pattern[T], then func(T) R) Chain[T, R]` | 模式匹配时执行 `then` 并进入下一链节。 |
| `When[R](predicate func(T) bool, then func(T) R) Chain[T, R]` | `Case(Predicate(predicate), then)` 的简写。 |
| `Default[R](then func(T) R) R` | 没有 case 匹配时执行 `then`，终结。 |
| `OrElse[R](fallback R) R` | 没有 case 匹配时返回字面量。 |
| `OrElseGet[R](fallback func(T) R) R` | 没有 case 匹配时调用 `fallback`。 |

`Chain[T, R]` 重复以上 API（除 `Case`/`When` 之外，因为已经链式起来了），外加：

- `Matched() bool` —— 是否已有 case 匹配。
- `Unwrap() (R, bool)` —— 已匹配返回 `(R, true)`，未匹配返回 `(R 零值, false)`。

### 类型匹配入口

```go
func Type(value any) TypeMatcher
```

`TypeMatcher` / `TypeChain[R]` 上的方法：

| 方法 | 说明 |
|---|---|
| `Case[T, R](then func(T) R) TypeChain[R]` | 动态类型为 `T` 时执行 `then`（`T` 为接口时表示"实现 `T`"）。 |
| `CaseWhen[T, R](guard func(T) bool, then func(T) R) TypeChain[R]` | 类型 + 守卫。 |
| `Nil[R](then func() R) TypeChain[R]` | 匹配 untyped nil 或 typed nil 指针 / map / slice / chan / func / interface。 |
| `When[R](predicate func(any) bool, then func(any) R) TypeChain[R]` | 对原始 `any` 自定义谓词。 |
| `Default[R](then func(any) R) R` / `OrElse[R](fallback R) R` / `OrElseGet[R](fallback func(any) R) R` | 终结。 |

`TypeChain` 也提供 `Matched() bool` / `Unwrap() (R, bool)`。

### 例子

```go
import "github.com/qianwj/typed/control/match"
import "github.com/qianwj/typed/utils/option"

desc := match.Value(httpStatus).
    Case(match.In(200, 201, 204), func(code int) string { return "ok" }).
    Case(match.Between(400, 499), func(code int) string { return "client" }).
    Case(match.Between(500, 599), func(code int) string { return "server" }).
    Default(func(code int) string { return "other" })
```

```go
kind := match.Type(payload).
    Case(func(s string) string { return "string:" + s }).
    Case(func(n int) string { return fmt.Sprintf("int:%d", n) }).
    Case(func(b []byte) string { return "bytes" }).
    Nil(func() string { return "<nil>" }).
    Default(func(v any) string { return fmt.Sprintf("%T", v) })
```

> 类型匹配的 `Case[T]` 用 `func(T) R` 而不是直接传 `T`，是因为 `T` 在编译期是泛型参数、运行时是反射的动态类型 —— 用 handler 函数让 Go 编译器在每个分支都内联一次断言。

## 与其他包的关系

- `control.Repeat(n, f)` 是 `collections.Range(0, n).ForEach(f)` 的命令式等价（少一次 `Stream` 分配），见 [`collections`](../collections/README-cn.md)。
- `match.Case` 与 `Option` / `Result` 互补：拿不准"到底命中什么"用模式匹配，单纯"值可能是 X"用 `Optional`，"操作可能失败"用 `Result`，见 [`option`](../option/README-cn.md) / [`result`](../result/README-cn.md)。
