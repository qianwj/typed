// Package adt provides typed algebraic data types — value types
// that explicitly model "present vs absent", "success vs failure",
// or "this variant vs that variant".
//
// The three ADT types each live in their own subpackage:
//
//   - [option.Option][T] — zero-or-one of a single type. The typed
//     alternative to `(T, bool)`. Lives in the [option] subpackage.
//   - [result.Result][T] — success carrying T, or failure carrying a
//     non-nil error. The typed alternative to `(T, error)`. Lives in
//     the [result] subpackage.
//   - [either.Either][L, R] — exactly one of two values. The general
//     tagged-union primitive. Lives in the [either] subpackage.
//
// The three cross-reference each other in idiomatic ways:
// [result.Result.Option] returns an [option.Option],
// [either.Either.Left] returns an [option.Option], and
// [either.Either.Right] returns an [option.Option]. Splitting them
// into one subpackage each keeps the call sites short
// (`option.Of`, `result.Wrap`, `either.Left[…]`) without an
// import-by-import cross-reference tax.
//
// The [Cast] helper stays in the umbrella `adt` package: it is not
// a value type, just a typed-alternative to Go's type assertion that
// returns a Result on outcome. It is the only thing this top-level
// package exports.
//
// Why a separate module? They are value types, not utilities in
// the same sense as a nil-check helper or a JSON codec; the `adt`
// module gives them a distinct import path so consumers can depend
// on the value types without pulling in unrelated `utils/*` code.
// This module depends only on the Go standard library and the
// `option` / `result` / `either` subpackages in the same module.
package adt

import (
	"fmt"
	"reflect"

	"github.com/qianwj/typed/adt/result"
)

// Cast asserts v's dynamic type as T and returns the outcome as a Result.
// It performs a type assertion, not a conversion: an int32 cannot be cast
// to int64. For an interface T, v's dynamic type must implement T.
//
// A mismatch or a nil interface returns Failure with the actual and target
// types in the error. A matching typed nil returns Success containing that
// nil value. Cast does not panic on a failed assertion.
func Cast[T any](v any) result.Result[T] {
	value, ok := v.(T)
	if !ok {
		return result.Failure[T](fmt.Errorf("adt.Cast: cannot assert %T as %v", v, reflect.TypeFor[T]()))
	}
	return result.Success(value)
}
