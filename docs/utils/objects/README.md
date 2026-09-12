# `typed/utils/objects` — Typed

Helpers for working with dynamically typed Go values.

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [`IsNil[T]`](#isnilt-value-t-bool)
- [`Equals[T]`](#equalst-a-b-t-bool)
- [See also](#see-also)

## Import

```go
import "github.com/qianwj/typed/utils/objects"
```

## `IsNil[T any](value T) bool`

Reports whether `value` is nil in the Go sense:

- An untyped nil interface.
- A typed nil: pointer / map / slice / channel / function / interface.
- `int` / `string` / struct / array — never nil; returns `false`.

```go
var p *Foo
objects.IsNil(p)             // true
objects.IsNil(nil)           // true
objects.IsNil((*Foo)(nil))   // true
objects.IsNil([]int(nil))    // true
objects.IsNil(0)             // false
```

## `Equals[T any](a, b T) bool`

A more humane "deep equality" than `reflect.DeepEqual`. Three intentional differences:

1. **Custom `Equal` method dispatch.** If either operand's dynamic type defines a method with signature `func (T) Equal(T) bool`, that method is called **instead of** a recursive walk. The operand whose dynamic type has the `Equal` method is treated as the authoritative source. Both sides are checked, so `Equals(a, b) == Equals(b, a)` when only one side has a custom `Equal` method: that side's method is consulted in both orderings. Types like `time.Time` that do not declare any `Equaler` interface still get picked up — the lookup is by signature shape, not by interface name.

2. **Nil-vs-empty collection normalisation, recursively.** A nil slice and an empty non-nil slice are considered equal at every level (not just the top). The same rule applies to maps. `reflect.DeepEqual` returns `false` for these pairs, which is rarely what callers want.

3. **Nil-safe custom `Equal` methods.** If `Equal` is defined on a pointer receiver and either operand is a typed nil pointer, both-nil returns `true` and exactly-one-nil returns `false` — **without ever calling the nil-receiver method**. This avoids the panic a naive "call `Equal` first, then check nil" implementation would produce.

All other aspects mirror `reflect.DeepEqual`: pointers are followed, structs are compared field by field (including unexported fields), arrays element by element, cyclic structures do not blow the stack.

### Opting in to a custom `Equal`

Add a method with the exact signature `func (T) Equal(T) bool`. You do **not** need to implement the `objects.Equaler` interface (it is kept for documentation and as a type hint; `Equals` does not require it). Mismatched signatures such as `Equal(Interface) bool` are ignored — `Equals` falls back to the structural walk.

```go
type Vec2 struct{ X, Y int }
func (a Vec2) Equal(b Vec2) bool { return a.X == b.X && a.Y == b.Y }

objects.Equals(Vec2{1, 2}, Vec2{1, 2}) // true (via Equal)
objects.Equals([]int(nil), []int{})    // true (nil-vs-empty normalisation)
```

## See also

- `control/match.Eq` uses `objects.Equals` as its "equal to" semantics, so `match.Eq(v).Match(x)` enjoys the same `Equal(T) bool` dispatch, nil-vs-empty normalisation, and typed-nil safety. See [`control/match`](../../control/README.md).
