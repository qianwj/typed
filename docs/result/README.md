# `github.com/qianwj/typed/utils/result`

`Result[T]`, an explicit container for the outcome of an operation that may fail: a success carrying a value of type `T`, or a failure carrying a non-nil `error`. It is a struct, not an interface — that lets `Map[R]` / `FlatMap[R]` declare their own type parameters (Go 1.27+ generic methods).

> Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Import

```go
import "github.com/qianwj/typed/utils/result"
```

## Why not `(T, error)`?

- The call site is self-documenting; combinators chain without scattering `if err != nil` blocks at every step.
- The error channel is fixed to `error` (no custom `E`), so interop with ordinary Go code is zero-cost.
- A whole pipeline can chain combinators all the way to the end before deciding whether to fall back to `(T, error)`.

## Three bridges

| Bridge | Direction | Function |
|---|---|---|
| `result.Wrap` | `(T, error) → Result[T]` | `Wrap[T](value T, err error) Result[T]` |
| `result.Unwrap` | `Result[T] → (T, error)` | `(Result[T]) Unwrap() (T, error)` |
| `Result.Optional` | `Result[T] → option.Optional[T]` | `(Result[T]) Optional() option.Optional[T]` |

- `Wrap` and `Unwrap` are inverses, so `Result` and `(T, error)` chains can switch back and forth.
- `Optional` is a **one-way** bridge: success becomes a present `Optional`; failure becomes an absent one (and the error is dropped). Use it only when the caller has already decided the error channel can be discarded.

## Construction

| Factory | Semantics |
|---|---|
| `Success[T](value T) Result[T]` | Success. |
| `Failure[T](err error) Result[T]` | Failure. **Panics if `err == nil`.** |
| `Wrap[T](value T, err error) Result[T]` | Standard `(T, error)` adapter: `err == nil` becomes success; otherwise failure (`value` is still stored on the struct but observationally equivalent to `Wrap(zero, err)`). `Wrap` does **not** panic when `err == nil`. |

The zero value of `Result` is "a failure with a nil error" — do not construct it directly. **Always** use `Success` / `Failure` / `Wrap`.

## State queries

| Method | Description |
|---|---|
| `IsSuccess() bool` | Success. |
| `IsFailure() bool` | Failure. |
| `Error() error` | The failure error when failed; `nil` when successful. |

## Value access

| Method | Behavior on failure |
|---|---|
| `Value() T` | **Panics**, panic value is the underlying `error`. |
| `Unwrap() (T, error)` | Returns `(zero, error)`. Does **not** panic. The exit at the end of a chain when you want to hand control back to ordinary Go error handling. |
| `OrElse(default T) T` | Returns `default`. `default` is always evaluated. |
| `OrElseGet(f func() T) T` | Returns `f()`. Runs only on failure. |
| `Recover(f func(error) T) T` | Returns `f(err)`. Runs only on failure, **cannot re-fail** — to re-fail, use `Unwrap`. |
| `Optional() option.Optional[T]` | Convert to `Optional[T]` (see above). |

## Chained transforms

| Method | Description |
|---|---|
| `Map[R](f func(T) R) Result[R]` | `f(value)` on success; failure is propagated unchanged. `f` is **not called** on failure. |
| `FlatMap[R](f func(T) Result[R]) Result[R]` | `f(value)` on success; failure is propagated unchanged. `f` is **not called** on failure. |
| `MapError(f func(error) error) Result[T]` | `f(err)` on failure only; success is returned unchanged. Useful for attaching context to an error before handing the result up the call stack. |

## Examples

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
// Recover: turn a recoverable error into a default value
port, err := result.Wrap(lookupPort()).
    Recover(func(err error) int { return 8080 }).
    Unwrap() // err is always nil
_ = port
```

## See also

- `Result[T]` bridges to `option.Optional[T]` through `Optional()`. See [`option`](../option/README.md).
- `utils/json`'s `Encode` / `Decode` both return `Result[...]`, adapting standard-library `(T, error)` calls into the `Result` chain. See [`utils/json`](../utils/json/README.md).
- `reactivex` reports subscription errors through `OnError`, a different error channel from `Result`. See [`reactivex`](../reactivex/README.md).
