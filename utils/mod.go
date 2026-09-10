// Package utils provides small, dependency-free generic helpers that are
// used across the typed modules: Optional[T] for "value that may be
// absent" and Result[T] for "value that may have failed".
//
// Both types follow the same idea: a small value type that wraps a
// present/absent (or success/failure) state, plus a few combinators
// (Map, FlatMap, OrElse, …) that keep call sites readable. Neither
// type is an interface, which lets their methods declare their own
// type parameters (Map[R]) under Go 1.27's generic methods.
//
// Why not the standard `(T, bool)` / `(T, error)` shapes?
//
//   - An explicit Optional/Result makes the call site self-documenting
//     and lets combinators chain without scattering boolean or error
//     checks at every step.
//   - Both types are designed to work for any T, including value types
//     such as int, string and struct{}, not just pointer-like ones.
//     This is why Optional carries a separate present flag rather than
//     relying on a nil check on the value.
//
// # File layout
//
// The package is split by responsibility rather than by type:
//
//   - mod.go        — package doc and small internal helpers shared by
//     both Optional and Result (nil check, error wrap)
//   - optional.go   — Optional[T] and its combinators
//   - result.go     — Result[T] and its combinators
//
// The internal helpers in this file are deliberately unexported: they
// are implementation details, not part of the package's public API.
package utils

import (
	"errors"
	"reflect"
)

// isNil reports whether v is a typed nil in the Go sense: a nil
// pointer, nil slice, nil map, nil channel, nil function value, or
// a nil interface. Value types (int, string, struct, …) are never
// nil and always return false.
//
// isNil is the building block for Optional.OfNullable. Reflection is
// used because Go's generic type system cannot express "is this value
// a nil of its type" without a runtime check; the alternative — a
// type switch — is not reachable from a generic function.
//
// The untyped-nil-interface case (`var i any; OfNullable(i)`) is
// caught up front: reflect.ValueOf on an untyped nil returns a
// zero Value with Kind == Invalid, which would otherwise slip
// through the switch below.
func isNil[T any](v T) bool {
	if any(v) == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

// errFromMsg wraps a string into an error. It is a tiny indirection
// over errors.New so the rest of the package does not have to import
// "errors" repeatedly.
func errFromMsg(msg string) error {
	return errors.New(msg)
}
