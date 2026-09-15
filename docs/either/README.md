# `typed/utils/either` — Typed

`Either[L, R]` is a tagged-union value type that holds exactly one of two values: a `Left` of type `L` or a `Right` of type `R`. The conventional reading is **"Left = failure, Right = success"**, so `Right[error, T]` is the typed equivalent of Go's idiomatic `(T, error)` pair, but with explicit handling instead of a hidden zero-value convention.

```go
result := either.Right[error, int](42)
val, ok := result.Right().Get() // val == 42, ok == true
```

## Contents

- [Why a tagged union?](#why-a-tagged-union)
- [API](#api)
- [Comparison with `Optional[T]`](#comparison-with-optionalt)
- [When to use Either](#when-to-use-either)
- [See also](#see-also)

## Why a tagged union?

Go's `(T, error)` shape is convenient for one-shot call sites but awkward once a value has to flow through several intermediate functions: every step has to thread an `error` variable and repeat the nil check. `Either` keeps both sides first-class so a function that may fail in two distinct ways can return:

```go
either.Either[error, T]   // generic "value or error"
either.Either[errA, errB] // two error categories
either.Either[T, U]       // "either a T or a U"
```

…and the call site picks the right branch with `IsLeft` / `IsRight` or the `Fold` combinator.

## API

```go
type Either[L, R any] struct { /* ... */ }

// Constructors
func Left[L, R any](l L) Either[L, R]
func Right[L, R any](r R) Either[L, R]

// Predicates
func (e Either[L, R]) IsLeft() bool
func (e Either[L, R]) IsRight() bool

// Safe accessors returning Optional
func (e Either[L, R]) Left() option.Optional[L]
func (e Either[L, R]) Right() option.Optional[R]

// Zero-value fallbacks (no allocation)
func (e Either[L, R]) LeftOrZero() L
func (e Either[L, R]) RightOrZero() R

// Combinators
func (e Either[L, R]) MapLeft[L2 any](f func(L) L2) Either[L2, R]
func (e Either[L, R]) MapRight[R2 any](f func(R) R2) Either[L, R2]
func (e Either[L, R]) Fold[T any](onLeft func(L) T, onRight func(R) T) T
```

### Construction

`Left` and `Right` are the only constructors. The zero value of `Either` is treated as the Left with both fields zero — prefer the explicit constructors.

### Safe accessors vs `OrZero`

`Left()` / `Right()` return `Optional[L]` / `Optional[R]`, so call sites can chain with the same combinators they already use for `Optional`:

```go
result := either.Right[error, int](42)
opt := result.Right()              // option.Optional[int], present
val := opt.Map(func(n int) int { return n * 2 }).OrElse(0) // 84
```

`LeftOrZero` / `RightOrZero` skip the `Optional` allocation when reading the wrong side is harmless (e.g. logging).

### `Fold`

`Fold` is the canonical "pick the right branch" combinator. Exactly one of `onLeft` / `onRight` is invoked:

```go
result.Fold(
    func(e error) string { return "err: " + e.Error() },
    func(n int) string { return fmt.Sprintf("val: %d", n) },
)
```

## Comparison with `Optional[T]`

`Optional[T]` holds zero or one value of a single type. `Either[L, R]` holds exactly one value, but the value's type depends on which side was taken. Optional answers "is this here?"; Either answers "which of these two is it?".

The safe accessors (Left / Right) return `Optional[L]` / `Optional[R]` so call sites can chain with the same combinators they already use for Optional.

## When to use Either

Reach for `Either` when:

- A function returns one of two structurally different value types (e.g. parse success vs. parse error details).
- You want explicit error categories rather than a single `error` interface value.
- You are composing functions that each can fail and want to thread the failure type through without scattering `nil` checks.

Reach for `(T, error)` when:

- The call site is a single one-shot `result, err := ...` and the error handling is a simple early return.

Reach for [`Optional[T]`](../option/README.md) when:

- The "missing" case has no associated failure detail.

## See also

- [`Optional[T]`](../option/README.md) — zero-or-one of a single type.
- [`Result`](../result/README.md) — success / failure container with a different surface.
- [`reactivex.Single[T]`](../reactivex/README.md#singlet) / [`reactivex.Maybe[T]`](../reactivex/README.md#maybet) — reactive containers with `Either`-shaped terminal events.