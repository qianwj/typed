# `github.com/qianwj/typed/utils/result`

`Result[T]`，显式的"成功携带 `T` 或失败携带 `error`"容器。结构体而非接口 —— 这让 `Map[R]` / `FlatMap[R]` 等方法可以声明自己的类型参数（Go 1.27+ 泛型方法）。

## 包导入

```go
import "github.com/qianwj/typed/utils/result"
```

## 设计动机

为什么不直接用 `(T, error)`？

- 调用点自我解释；组合子可以链式拼装，不用每一步散落 `if err != nil`。
- 错误通道固定为 `error`（不是自定义 `E`），与标准 Go 互操作零成本。
- 整条链可以一路组合到结尾，再决定要不要回到 `(T, error)` 形态。

## 三个桥

| 桥 | 方向 | 函数 |
|---|---|---|
| `result.Wrap` | `(T, error) → Result[T]` | `Wrap[T](value T, err error) Result[T]` |
| `result.Unwrap` | `Result[T] → (T, error)` | `(Result[T]) Unwrap() (T, error)` |
| `Result.Optional` | `Result[T] → option.Optional[T]` | `(Result[T]) Optional() option.Optional[T]` |

- `Wrap` 与 `Unwrap` 互为逆操作，让 `Result` 链与 `(T, error)` 链可以无缝切换。
- `Optional` 是**单向**桥：成功 → 出席 Optional；失败 → 缺席 Optional。错误被丢弃，仅在"调用点已经决定不再关心错误"时使用。

## 构造

| 工厂 | 语义 |
|---|---|
| `Success[T](value T) Result[T]` | 成功。 |
| `Failure[T](err error) Result[T]` | 失败；**`err == nil` 会 panic**。 |
| `Wrap[T](value T, err error) Result[T]` | 标准 `(T, error)` 适配器：`err == nil` → 成功；否则失败（`value` 在结构体里仍被存储但观测上等价于 `Wrap(zero, err)`）。`Wrap` 在 `err == nil` 时**不** panic。 |

`Result` 的零值是"带 nil error 的失败" —— 不要直接构造零值，**总是**用 `Success` / `Failure` / `Wrap`。

## 状态查询

| 方法 | 说明 |
|---|---|
| `IsSuccess() bool` | 成功。 |
| `IsFailure() bool` | 失败。 |
| `Error() error` | 失败时返回错误；成功时返回 nil。 |

## 取值

| 方法 | 失败时行为 |
|---|---|
| `Value() T` | **panic**，panic 值是底层 `error`。 |
| `Unwrap() (T, error)` | 返回 `(零值, error)`，**不 panic**。链末"回到普通 Go 错误处理"的出口。 |
| `OrElse(default T) T` | 返回 `default`；`default` 总是被求值。 |
| `OrElseGet(f func() T) T` | 返回 `f()`；仅失败时调用。 |
| `Recover(f func(error) T) T` | 返回 `f(err)`；仅失败时调用，**不**再失败 —— 想再失败就用 `Unwrap`。 |
| `Optional() option.Optional[T]` | 转 `Optional[T]`（见上）。 |

## 链式 transform

| 方法 | 说明 |
|---|---|
| `Map[R](f func(T) R) Result[R]` | 成功时 `f(value)`；失败时原样传播，**不调用** `f`。 |
| `FlatMap[R](f func(T) Result[R]) Result[R]` | 成功时 `f(value)`；失败时原样传播，**不调用** `f`。 |
| `MapError(f func(error) error) Result[T]` | 仅失败时 `f(err)`；成功时不变。常用于在错误上附加上下文。 |

## 例子

```go
import "github.com/qianwj/typed/utils/result"

val, err := result.Wrap(loadConfig(path)).
    MapError(func(err error) error { return fmt.Errorf("config %s: %w", path, err) }).
    Map(func(b []byte) Config { return parseConfig(b) }).
    FlatMap(func(c Config) result.Result[Config] {
        if c.Version == 0 {
            return result.Failure[Config](errors.New("version required"))
        }
        return result.Success(c)
    }).
    Unwrap()
if err != nil {
    return err
}
use(val)
```

```go
// Recover: 把可恢复的错误映射回一个默认值
port, err := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err 始终是 nil
_ = port
```

## 与其他包的关系

- `Result[T]` 通过 `Optional()` 桥到 `option.Optional[T]`，见 [`option`](../option/README.md)。
- `utils/json` 的 `Encode` / `Decode` 都返回 `Result[...]`，把 `(T, error)` 风格的标准库调用适配进 `Result` 链，见 [`utils/json`](../utils/json/README.md)。
- `reactivex` 的订阅级错误通过 `OnError` 通知，与 `Result` 是不同的错误通道，见 [`reactivex`](../reactivex/README.md)。
