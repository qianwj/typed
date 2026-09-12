# `typed/utils/option` — Typed

`Optional[T]`, an explicit container for a value of `T` that may be absent. It is a struct, not an interface, so its methods can declare their own type parameters (`Map[R]`, `FlatMap[R]`) — a Go 1.27+ generic-method feature.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [Why not `(T, bool)`?](#why-not-t-bool)
- [Construction](#construction)
- [State queries](#state-queries)
- [Value access](#value-access)
- [Side effects](#side-effects)
- [Chained transforms](#chained-transforms)
- [Examples](#examples)
- [See also](#see-also)

## Import

```go
import "github.com/qianwj/typed/utils/option"
```

## Why not `(T, bool)`?

- The call site is self-documenting; combinators chain without scattering `if ok { ... }` blocks at every step.
- `Optional` treats value types (`int` / `string` / `struct{}`) and pointer-like types uniformly — it carries an explicit `present` flag instead of relying on a nil check on the value. There is no "is it zero or is it absent" ambiguity for `int`.
- It keeps `result.Result[T]` and other `(T, error)`-shaped combinators decoupled from "is it absent?".

## Construction

| Factory | Semantics |
|---|---|
| `Empty[T any]() Optional[T]` | Carries no value. The zero value `Optional[T]{}` is equivalent. |
| `Of[T any](value T) Optional[T]` | Carries `value`. **Panics on a typed nil** (e.g. a `nil *Foo` passed as a pointer-typed `T`). |
| `OfNullable[T any](value T) Optional[T]` | Applies Go's nil semantics: a nil pointer / map / slice / chan / func / interface becomes `Empty`; anything else is treated as `Of(value)`. |

> Pick `Of` vs `OfNullable` based on `T`: value types use `Of`; pointer-like and interface types use `OfNullable`.

## State queries

| Method | Description |
|---|---|
| `IsPresent() bool` | Carries a value. |
| `IsEmpty() bool` | `!IsPresent()`. |

## Value access

| Method | Behavior when absent |
|---|---|
| `Get() T` | **Panics** (no message). |
| `OrElse(default T) T` | Returns `default`. `default` is always evaluated. |
| `OrElseGet(f func() T) T` | Returns `f()`. `f` runs only when absent. |
| `OrElseThrow(errMsg string) (T, error)` | Returns `(zero, errors.New(errMsg))` when absent; `(value, nil)` when present. **Does not panic** — the name follows Java's `Optional.orElseThrow` convention, adapted to Go's explicit `(T, error)`. |

## Side effects

| Method | Description |
|---|---|
| `IfPresent(f func(T))` | Calls `f` when present; does nothing when absent. |
| `IfPresentOrElse(present func(T), absent func())` | Exactly one of the two callbacks runs. |

## Chained transforms

| Method | Description |
|---|---|
| `Filter(predicate func(T) bool) Optional[T]` | Returns this `Optional` when present and `predicate(value) == true`; otherwise `Empty`. `predicate` is not called when the receiver is absent. |
| `Map[R](f func(T) R) Optional[R]` | Returns `Empty[R]` when absent; **does not call** `f`. |
| `FlatMap[R](f func(T) Optional[R]) Optional[R]` | **Does not call** `f` when absent. When called, the `Optional[R]` returned by `f` is taken as the result. |

## Examples

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
v := empty.OrElseGet(func() int { return compute() }) // compute runs only here
```

## See also

- Every "may be absent" accessor in `collections` returns `option.Optional[T]`: `ArrayList.Get` / `First` / `Last` / `Find` / `MinBy` / `MaxBy`, the `Stack` / `Queue` / `Deque` `Pop` / `Peek` / `Front` / `Back`, and `Stream.First` / `Last` / `Find`. See [`collections`](../collections/README.md).
- `Result[T]` bridges to `Optional[T]` through `Optional()`. See [`result`](../result/README.md).
