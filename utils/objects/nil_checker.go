// Package objects provides helpers for working with dynamically typed Go
// values.
package objects

import "reflect"

// IsNil reports whether value is nil in the Go sense.
//
// In addition to an untyped nil interface, IsNil recognizes typed nil
// pointers, maps, slices, channels, functions, and interfaces. Values such
// as integers, strings, structs, and arrays are never nil.
func IsNil[T any](value T) bool {
	if any(value) == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}
