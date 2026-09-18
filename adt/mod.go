// Package adt provides typed algebraic data types — value types
// that explicitly model "success vs failure" or "this variant vs
// that variant".
//
// The types in this package are:
//
//   - [Result][T] — success carrying T, or failure carrying a
//     non-nil error. The typed alternative to `(T, error)`.
//   - [Either][L, R] — exactly one of two values. The general
//     tagged-union primitive.
//
// The zero-or-one counterpart lives in the [option] subpackage:
// [option.Option][T]. Option and Result / Either share the same
// vocabulary (present / absent / success / failure / left / right)
// and they cross-reference each other in idiomatic ways —
// [Result.Option] returns an [option.Option], and
// [Either.Right] returns an [option.Option]. Putting Option in
// its own subpackage keeps its bulk idiomatic import path
// (`option.Option`, `option.Empty`, `option.Of`, …) short while
// leaving Result / Either in the `adt` namespace.
//
// Why a separate module? They are value types, not utilities in
// the same sense as a nil-check helper or a JSON codec; the `adt`
// module gives them a distinct import path so consumers can depend
// on the value types without pulling in unrelated `utils/*` code.
// This module depends only on the Go standard library and the
// `option` subpackage in the same module.
package adt

import (
	"fmt"
	"reflect"
)

// Cast asserts v's dynamic type as T and returns the outcome as a Result.
// It performs a type assertion, not a conversion: an int32 cannot be cast
// to int64. For an interface T, v's dynamic type must implement T.
//
// A mismatch or a nil interface returns Failure with the actual and target
// types in the error. A matching typed nil returns Success containing that
// nil value. Cast does not panic on a failed assertion.
func Cast[T any](v any) Result[T] {
	value, ok := v.(T)
	if !ok {
		return Failure[T](fmt.Errorf("adt.Cast: cannot assert %T as %v", v, reflect.TypeFor[T]()))
	}
	return Success(value)
}
