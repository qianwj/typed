# `typed/adt` — Typed

[![codecov](https://codecov.io/gh/qianwj/typed/graph/badge.svg?flag=adt)](https://codecov.io/gh/qianwj/typed)

Typed algebraic data types — value types that explicitly model "present vs absent", "success vs failure", or "this variant vs that variant". The package ships three types:

- [`Option[T]`](#optionalt) — zero-or-one of a single type. The typed alternative to `(T, bool)`.
- [`Result[T]`](#resultt) — success carrying `T`, or failure carrying a non-nil `error`. The typed alternative to `(T, error)`.
- [`Either[L, R]`](#eitherl-r) — exactly one of two values. The general tagged-union primitive the first two specialise.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [Why one package?](#why-one-package)
- [`Option[T]`](#optionalt)
- [`Result[T]`](#resultt)
- [`Either[L, R]`](#eitherl-r)
- [`Cast[T]`](#castt)
- [Choosing between them](#choosing-between-them)
- [See also](#see-also)

## Import

```go
import "github.com/qianwj/typed/adt"
```

All three types share the same import — `adt.Option`, `adt.Result`, `adt.Either`.

## Why one package?

The three types share the same vocabulary (present / absent / success / failure / left / right) and cross-reference each other in idiomatic ways:

- `Result.Option` returns an `Option[T]` — a success becomes present, a failure becomes absent.
- `Either.Left` / `Either.Right` return `Option[L]` / `Option[R]` — safe accessors over the tagged union.

Co-locating them in one package makes those relationships visible at the call site without the import-by-import friction of separate sub-packages. A separate module (`typed/adt`, not `typed/utils/adt`) signals that they are value types, not utilities in the same sense as `objects.IsNil` or a JSON codec — consumers can depend on the value types without pulling in unrelated `utils/*` code.

## `Option[T]`

`Option[T]` is a container that may or may not hold a value of type `T`. An `Option` is either present (carrying a `T`) or absent (no value). The zero value of `Option` is absent and is equivalent to `Empty[T]()`.

`Option` is not an interface, which lets its methods declare their own type parameters (`Map[R]`, `FlatMap[R]`) — a Go 1.27+ generic-method feature.

### Why not `(T, bool)`?

- The call site is self-documenting; combinators chain without scattering `if ok { ... }` blocks at every step.
- `Option` treats value types (`int` / `string` / `struct{}`) and pointer-like types uniformly — it carries an explicit `present` flag instead of relying on a nil check on the value. There is no "is it zero or is it absent" ambiguity for `int`.
- It keeps `Result[T]` and other `(T, error)`-shaped combinators decoupled from "is it absent?".

### Construction

| Factory | Semantics |
|---|---|
| `Empty[T any]() Option[T]` | Carries no value. The zero value `Option[T]{}` is equivalent. |
| `Of[T any](value T) Option[T]` | Carries `value`. **Panics on a typed nil** (e.g. a `nil *Foo` passed as a pointer-typed `T`). |
| `OfNullable[T any](value T) Option[T]` | Applies Go's nil semantics: a nil pointer / map / slice / chan / func / interface becomes `Empty`; anything else is treated as `Of(value)`. |

> Pick `Of` vs `OfNullable` based on `T`: value types use `Of`; pointer-like and interface types use `OfNullable`.

### State queries

| Method | Description |
|---|---|
| `IsPresent() bool` | Carries a value. |
| `IsEmpty() bool` | `!IsPresent()`. |

### Value access

| Method | Behaviour when absent |
|---|---|
| `Get() T` | **Panics** ("Get on empty Option"). |
| `OrElse(default T) T` | Returns `default`. `default` is always evaluated. |
| `OrElseGet(f func() T) T` | Returns `f()`. `f` runs only when absent. |
| `OrElseThrow(errMsg string) (T, error)` | Returns `(zero, errors.New(errMsg))` when absent; `(value, nil)` when present. **Does not panic** — the name follows Java's `Option.orElseThrow` convention, adapted to Go's explicit `(T, error)`. |

### Side effects

| Method | Description |
|---|---|
| `IfPresent(f func(T))` | Calls `f` when present; does nothing when absent. |
| `IfPresentOrElse(present func(T), absent func())` | Exactly one of the two callbacks runs. |

### Chained transforms

| Method | Description |
|---|---|
| `Filter(predicate func(T) bool) Option[T]` | Returns this `Option` when present and `predicate(value) == true`; otherwise `Empty`. `predicate` is not called when the receiver is absent. |
| `Map[R](f func(T) R) Option[R]` | Returns `Empty[R]` when absent; **does not call** `f`. |
| `FlatMap[R](f func(T) Option[R]) Option[R]` | **Does not call** `f` when absent. When called, the `Option[R]` returned by `f` is taken as the result. |

### Examples

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
v := empty.OrElseGet(func() int { return compute() }) // compute runs only here
```

## `Result[T]`

`Result[T]` is the outcome of an operation that may fail. A `Result` is either a success carrying a value of type `T`, or a failure carrying a non-nil error.

The zero value of `Result` is a failure with a nil error; callers should always build results through `Success` / `Failure` and inspect them through the methods.

### Why not `(T, error)`?

- An explicit `Result` makes the call site self-documenting and lets combinators chain without scattering `if err != nil` blocks at every step.
- The shape is intentionally simpler than a generic `Result[T, E]` — the error channel is always a standard `error`, which keeps the interop with regular Go code free of custom-error boxing.

### Construction

| Factory | Semantics |
|---|---|
| `Success[T any](value T) Result[T]` | Carries `value`. |
| `Failure[T any](err error) Result[T]` | Carries `err`. **Panics on a nil error** — a nil error must go through `Success`. |
| `Wrap[T any](value T, err error) Result[T]` | The bridge from Go's idiomatic `(T, error)`. Success when `err == nil`, failure otherwise. `value` is stored verbatim in either case but only reachable through methods on success. |

`Wrap` is the right shape for two specific call sites: the end of a `(T, error)` API boundary where the caller wants to keep going inside a `Result`-based pipeline (Map, FlatMap, OrElse, ...), and the bridge back from a `Result` chain to plain Go error handling.

### State queries

| Method | Description |
|---|---|
| `IsSuccess() bool` | `err == nil`. |
| `IsFailure() bool` | `err != nil`. |

### Value access

| Method | Behaviour on failure |
|---|---|
| `Value() T` | **Panics** with the underlying error. |
| `Unwrap() (T, error)` | Returns `(zero, err)`. Does not panic — surfaces the error to the caller. The reverse bridge of `Wrap`. |
| `Error() error` | Returns the error, or nil on success. |
| `OrElse(default T) T` | Returns `default`. `default` is always evaluated. |
| `OrElseGet(f func() T) T` | Returns `f()`. `f` runs only on failure. |
| `Recover(f func(error) T) T` | Calls `f(err)` on failure, returning its result; the success value is returned untouched. Error-aware fallback. |

### Bridges

| Method | Description |
|---|---|
| `Option() Option[T]` | Success → present `Option`; failure → absent `Option`. The error is dropped — use only when the caller has decided the error channel can be discarded. |

### Chained transforms

| Method | Description |
|---|---|
| `Map[R](f func(T) R) Result[R]` | Returns a failure with the same error when failed; **does not call** `f`. |
| `FlatMap[R](f func(T) Result[R]) Result[R]` | Does not call `f` on failure. When called, the `Result[R]` returned by `f` is taken as the result. |
| `MapError(f func(error) error) Result[T]` | Calls `f(err)` on failure; the success is returned unchanged. Useful for wrapping low-level errors with a higher-level description before handing the `Result` up the call stack. |

### Examples

```go
import "github.com/qianwj/typed/adt"

r := adt.Wrap(loadProfile(id))
// r is a Result[Profile] that is success if err was nil, failure otherwise.

name := r.Map(func(p Profile) string { return p.Name }).
    Recover(func(err error) string {
        if errors.Is(err, ErrNotFound) {
            return "anonymous"
        }
        return "" // any other error → empty name
    })
```

## `Either[L, R]`

`Either[L, R]` is a tagged-union value type that holds exactly one of two values: a `Left` of type `L` or a `Right` of type `R`. The conventional reading is **"Left = failure, Right = success"** — for example `Right[error, T]` is the typed equivalent of Go's idiomatic `(T, error)` pair.

### Why a tagged union?

Go's `(T, error)` shape is convenient for one-shot call sites but awkward once a value has to flow through several intermediate functions: every step has to thread an `error` variable and repeat the nil check. `Either` keeps both sides first-class so a function that may fail in two distinct ways can return `Either[errA, errB]`, or a function that may produce one of two value types can return `Either[T, U]`. The call site picks the right branch with `IsLeft` / `IsRight` or the `Fold` combinator.

### Construction

| Factory | Semantics |
|---|---|
| `Left[L, R any](l L) Either[L, R]` | Carries `l`. The zero value of `Either` holds the zero value of `L` as a `Left`. |
| `Right[L, R any](r R) Either[L, R]` | Carries `r`. |

### Predicates

| Method | Description |
|---|---|
| `IsLeft() bool` | Carries a `Left`. |
| `IsRight() bool` | Carries a `Right`. |

### Safe accessors

| Method | Description |
|---|---|
| `Left() Option[L]` | Present when `IsLeft()`; absent when `IsRight()`. |
| `Right() Option[R]` | Present when `IsRight()`; absent when `IsLeft()`. |
| `LeftOrZero() L` | Returns the `Left` value, or the zero value of `L` when on `Right`. Use when reading the wrong side is harmless. |
| `RightOrZero() R` | Mirror of `LeftOrZero`. |

The safe accessors return `Option` so call sites can chain with the same combinators they already use elsewhere.

### Combinators

`MapLeft` / `MapRight` and `FlatMapLeft` / `FlatMapRight` explicitly select a branch. Each method passes the opposite branch through without invoking its callback. `FlatMapLeft` keeps `R` fixed and allows `L` to change; `FlatMapRight` keeps `L` fixed and allows `R` to change. Either callback may return a `Left` or a `Right`, so chains can switch branches. Use `Fold` to consume both alternatives.

| Method | Description |
|---|---|
| `MapLeft[L2](f func(L) L2) Either[L2, R]` | Calls `f` only on `Left`; the `Right` branch passes through unchanged. |
| `MapRight[R2](f func(R) R2) Either[L, R2]` | Calls `f` only on `Right`; passes `Left` through unchanged. |
| `FlatMapLeft[L2](f func(L) Either[L2, R]) Either[L2, R]` | On `Left`, returns the `Either` produced by `f`; passes `Right` through without calling `f`. |
| `FlatMapRight[R2](f func(R) Either[L, R2]) Either[L, R2]` | On `Right`, returns the `Either` produced by `f`; passes `Left` through without calling `f`. |
| `Fold[T](onLeft func(L) T, onRight func(R) T) T` | Calls exactly one of the two callbacks and returns its result. The canonical "pick the right branch" pattern. |

### Examples

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

`Cast[T any](v any) Result[T]` performs a type assertion without panicking on a mismatch. A concrete target requires an identical dynamic type; an interface target requires that the dynamic type implement it. This does not perform conversions such as `int32` to `int64`.

```go
adt.Cast[int](42).OrElse(0)           // 42
adt.Cast[int64](int32(42)).IsFailure() // true
adt.Cast[int](nil).IsFailure()        // true
adt.Cast[*int]((*int)(nil)).IsSuccess() // true: matching typed nil
```

A failed assertion returns an error describing the actual and target types. A nil interface fails; a matching typed nil succeeds and is preserved, unlike `OfNullable`, which treats typed nil as absent. Results compose with `Map`, `FlatMap`, `OrElse`, and `Unwrap`.

## Choosing between them

| Question | Reach for |
| --- | --- |
| May or may not have a value, with no associated failure detail | `Option[T]` |
| Always returns a value on success or an error on failure | `Result[T]` |
| Two distinct success / failure categories, or two value types | `Either[L, R]` |

`Option` is zero-or-one. `Result` is success-or-failure (the failure always carries a standard `error`). `Either` is general — both branches can be any type, including two distinct error types if you want error categorisation beyond a single `error` interface.

In practice:

- A `Map[K]V.Get(key)` that may miss → `Option[V]`.
- A database query that returns a row or errors → `Result[Row]`.
- A function that may produce `T` or a configuration error → `Either[ConfigError, T]`.

## See also

- Every "may be absent" accessor in [`collections`](../collections/README.md) returns `adt.Option[T]`: `ArrayList.Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy`, the `Stack` / `Queue` / `Deque` `Pop` / `Peek` / `Front` / `Back`, and `Stream.First` / `Last` / `Find`.
- [`utils/json`](../utils/json/README.md) uses `adt.Result[T]` for codec return values.
- [`reactivex.Single[T]`](../reactivex/README.md#singlet) / [`reactivex.Maybe[T]`](../reactivex/README.md#maybet) — reactive containers whose terminal events are `Either`-like.
