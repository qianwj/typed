# `github.com/qianwj/typed/control`

Control-flow helpers. Two sub-packages:

- [`control`](#control) — `Repeat` / `RepeatE`, concise "do this N times" loop primitives.
- [`control/match`](#controlmatch) — type-safe, first-match-wins pattern matching (`Pattern[T]` plus chained `Case` / `Type`).

Go 1.27+ is required because the fluent chain methods declare their own type parameters.

> Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Import

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

Calls `f` exactly `times` times and discards the result. `times <= 0` is a no-op (the body is never entered and `f` is not called). It is equivalent to `for i := 0; i < times; i++ { f() }`, just more explicit.

```go
control.Repeat(3, func() { fmt.Println("tick") })
```

### `RepeatE(times, f) (int, error)`

The error-aware variant. When `f` returns a non-nil error the loop stops immediately and returns `(i, err)` — where `i` is the count of **successfully completed** iterations (the error came from the `i+1`-th call).

- All success → `(times, nil)`.
- Failure on the `i+1`-th call (zero-indexed `i`) → `(i, err)`; the loop stops at that call.
- `times <= 0` → `(0, nil)` (the body is not entered). For `times < 0` the return is `(times, nil)` — that is symmetric with the input rather than a meaningful count. Negative `times` is a caller error; the function does not panic so it stays safe to call with computed counts, but treat a negative return as a signal to fix the input, not as a real count.

Use it for retry loops, batch operations that should stop on the first failure, and any "do this N times unless something goes wrong" workflow. A nil error from `f` is treated as success; only a non-nil error short-circuits the loop.

```go
n, err := control.RepeatE(5, func() error {
    return doOne()
})
if err != nil {
    log.Printf("after %d successful attempts: %v", n, err)
}
```

### Relationship to `collections.Range`

`collections.Range(0, n).ForEach(f)` constructs a `Stream` and then iterates it. `control.Repeat(n, f)` is the imperative equivalent that skips the `Stream` allocation. Pick `Repeat` for one-off side effects; pick `Stream` when the iteration is the start of a larger fluent pipeline (`Map`, `Filter`, `Take`, ...).

---

## `control/match`

First-match-wins chained pattern matching. Two modes:

1. **Value matching** — `Match(v)` / `Value(v)`, testing the input value against `Pattern[T]`.
2. **Type matching** — `Type(v)`, testing the dynamic type of an `any`.

The API is deliberately explicit: every `Case` must return the same `R`, cases are evaluated in declaration order, and once a case matches, later patterns and handlers are **not** evaluated.

### `Pattern[T]`

```go
type Pattern[T any] interface {
    Match(value T) bool
}
type PatternFunc[T any] func(T) bool
func (p PatternFunc[T]) Match(value T) bool
```

### Factories

| Factory | Semantics |
|---|---|
| `Predicate[T](p func(T) bool) Pattern[T]` | Arbitrary predicate. |
| `Any[T]() Pattern[T]` | Matches every value. |
| `Eq[T](want T) Pattern[T]` | Deep compare with `want` per `objects.Equals`; types with an `Equal(T) bool` method are dispatched specially. |
| `In[T comparable](values ...T) Pattern[T]` | Matches any value in the given set. |
| `Between[T cmp.Ordered](min, max T) Pattern[T]` | Closed interval `[min, max]`. |
| `Not[T](pattern Pattern[T]) Pattern[T]` | Negation. |
| `And[T](patterns ...Pattern[T]) Pattern[T]` | All match. With no arguments, matches every value. |
| `Or[T](patterns ...Pattern[T]) Pattern[T]` | At least one matches. With no arguments, matches nothing. |

### Value matching entry

```go
func Value[T any](value T) Matcher[T]   // entry
func Match[T any](value T) Matcher[T]   // alias for Value
```

Methods on `Matcher[T]` / `Chain[T, R]`:

| Method | Description |
|---|---|
| `Case[R](pattern Pattern[T], then func(T) R) Chain[T, R]` | Run `then` when the pattern matches and advance the chain. |
| `When[R](predicate func(T) bool, then func(T) R) Chain[T, R]` | Shorthand for `Case(Predicate(predicate), then)`. |
| `Default[R](then func(T) R) R` | Terminal: run `then` when no case matched. |
| `OrElse[R](fallback R) R` | Terminal: return a literal when no case matched. |
| `OrElseGet[R](fallback func(T) R) R` | Terminal: call `fallback` when no case matched. |

`Chain[T, R]` repeats the same API (except the first `Case` / `When`, because the chain is already started), and additionally:

- `Matched() bool` — whether any case has matched.
- `Unwrap() (R, bool)` — `(R, true)` on match, `(zero R, false)` on miss.

### Type matching entry

```go
func Type(value any) TypeMatcher
```

Methods on `TypeMatcher` / `TypeChain[R]`:

| Method | Description |
|---|---|
| `Case[T, R](then func(T) R) TypeChain[R]` | Run `then` when the dynamic type is `T` (or, if `T` is an interface, when the dynamic type implements `T`). |
| `CaseWhen[T, R](guard func(T) bool, then func(T) R) TypeChain[R]` | Type plus a guard. |
| `Nil[R](then func() R) TypeChain[R]` | Matches untyped nil or typed nil pointer / map / slice / chan / func / interface. |
| `When[R](predicate func(any) bool, then func(any) R) TypeChain[R]` | Arbitrary predicate on the raw `any`. |
| `Default[R](then func(any) R) R` / `OrElse[R](fallback R) R` / `OrElseGet[R](fallback func(any) R) R` | Terminals. |

`TypeChain` also provides `Matched() bool` / `Unwrap() (R, bool)`.

### Examples

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

> The `Case[T]` of type matching takes `func(T) R` rather than a bare `T` because `T` is a generic parameter at compile time but a reflective dynamic type at run time — having each branch receive a typed handler function lets the Go compiler inline a type assertion per branch.

## See also

- `control.Repeat(n, f)` is the imperative equivalent of `collections.Range(0, n).ForEach(f)` (one fewer `Stream` allocation). See [`collections`](../collections/README.md).
- `match.Case` complements `Option` and `Result`: use pattern matching when "what exactly did I match?" is the question, `Optional` when "the value may be X", and `Result` when "the operation may fail". See [`option`](../option/README.md) and [`result`](../result/README.md).
