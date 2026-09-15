# `typed/utils/either` — Typed

`Either[L, R]` 是一个 tagged-union 值类型，恰好持有两种值中的一种：类型 `L` 的 `Left` 或类型 `R` 的 `Right`。约定俗成的读法是 **"Left = 失败, Right = 成功"**，因此 `Right[error, T]` 是 Go 惯用的 `(T, error)` 元组的类型化等价物，但用显式分支处理代替隐式的零值约定。

```go
result := either.Right[error, int](42)
val, ok := result.Right().Get() // val == 42, ok == true
```

## 目录

- [为什么用 tagged union？](#为什么用-tagged-union)
- [API](#api)
- [和 `Optional[T]` 的对比](#和-optionalt-的对比)
- [何时用 Either](#何时用-either)
- [相关阅读](#相关阅读)

## 为什么用 tagged union？

Go 的 `(T, error)` 在单点调用处很顺手，但一旦值要在多个中间函数之间流转就显得笨拙：每一步都得传递 `error` 变量并重复 nil 检查。`Either` 把两侧都变成一等公民，因此一个可能以两种方式失败的函数可以返回：

```go
either.Either[error, T]   // 泛型的「值或 error」
either.Either[errA, errB] // 两种 error 类别
either.Either[T, U]       // 「T 或 U」
```

调用方用 `IsLeft` / `IsRight` 或者 `Fold` 组合子选对应分支。

## API

```go
type Either[L, R any] struct { /* ... */ }

// 构造
func Left[L, R any](l L) Either[L, R]
func Right[L, R any](r R) Either[L, R]

// 谓词
func (e Either[L, R]) IsLeft() bool
func (e Either[L, R]) IsRight() bool

// 安全访问器，返回 Optional
func (e Either[L, R]) Left() option.Optional[L]
func (e Either[L, R]) Right() option.Optional[R]

// 零值兜底（无分配）
func (e Either[L, R]) LeftOrZero() L
func (e Either[L, R]) RightOrZero() R

// 组合子
func (e Either[L, R]) MapLeft[L2 any](f func(L) L2) Either[L2, R]
func (e Either[L, R]) MapRight[R2 any](f func(R) R2) Either[L, R2]
func (e Either[L, R]) Fold[T any](onLeft func(L) T, onRight func(R) T) T
```

### 构造

`Left` 和 `Right` 是仅有的构造器。`Either` 的零值会被当作「Left 且两个字段都是零值」——优先用显式构造器。

### 安全访问器 vs `OrZero`

`Left()` / `Right()` 返回 `Optional[L]` / `Optional[R]`，因此调用方可以链上和 `Optional` 一样的组合子：

```go
result := either.Right[error, int](42)
opt := result.Right()              // option.Optional[int]，present
val := opt.Map(func(n int) int { return n * 2 }).OrElse(0) // 84
```

`LeftOrZero` / `RightOrZero` 在读错一侧无所谓时（比如日志）跳过 `Optional` 分配。

### `Fold`

`Fold` 是「挑对分支」的规范组合子。`onLeft` / `onRight` 中**恰好一个**会被调用：

```go
result.Fold(
    func(e error) string { return "err: " + e.Error() },
    func(n int) string { return fmt.Sprintf("val: %d", n) },
)
```

## 和 `Optional[T]` 的对比

`Optional[T]` 持有零个或一个单一类型的值。`Either[L, R]` 恰好持有一个值，但值的类型取决于取了哪一侧。Optional 回答「这个在不在？」；Either 回答「是这两个中的哪一个？」。

安全访问器（Left / Right）返回 `Optional[L]` / `Optional[R]`，因此调用方可以链上和 Optional 一样的组合子。

## 何时用 Either

- 函数返回两种结构不同的值（比如解析成功 vs 解析错误详情）
- 想要显式的错误类别，而不是单一的 `error` 接口值
- 串联多个可能失败的函数，想把失败类型一路传下去又不想散布 `nil` 检查

以下场景用 `(T, error)`：

- 单次 `result, err := ...`，错误处理就是简单 early return

以下场景用 [`Optional[T]`](../option/README.md)：

- 「缺失」这种情况没有附带失败细节

## 相关阅读

- [`Optional[T]`](../option/README.md) —— 零或一个单一类型的值
- [`Result`](../result/README.md) —— 另一种成功/失败容器
- [`reactivex.Single[T]`](../reactivex/README.md#singlet) / [`reactivex.Maybe[T]`](../reactivex/README.md#maybet) —— 终止事件是 `Either` 形状的 reactive 容器