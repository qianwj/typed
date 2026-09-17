package adt

import (
	"errors"
	"reflect"
)

// Option[T] is a container that may or may not hold a value of type T.
//
// An Option is either present (carrying a T) or absent (no value).
// The zero value of Option is absent and is equivalent to Empty[T]().
// Option values must not be copied after creation; the safe usage is
// to obtain them through one of the constructors (Empty, Of, OfNullable)
// and then consume them through the methods.
//
// Option is not an interface, which lets its methods declare
// their own type parameters (Map[R]) under Go 1.27's generic
// methods.
//
// Why not the standard `(T, bool)` shape?
//
//   - An explicit Option makes the call site self-documenting
//     and lets combinators chain without scattering boolean
//     checks at every step.
//   - Option is designed to work for any T, including value
//     types such as int, string and struct{}, not just
//     pointer-like ones. This is why Option carries a separate
//     present flag rather than relying on a nil check on the
//     value.
type Option[T any] struct {
	value   T
	present bool
}

// Empty returns an absent Option[T].
//
// Calling Empty on a value type such as int is the idiomatic way to say
// "no result" and is preferred over using a sentinel like 0 or "".
func Empty[T any]() Option[T] {
	return Option[T]{}
}

// Of returns a present Option[T] holding value.
//
// Of panics if value is a typed nil (e.g. a nil *Foo passed as a
// pointer-typed T). Use OfNullable when the underlying T is a
// reference type and nil is a meaningful value.
func Of[T any](value T) Option[T] {
	return Option[T]{value: value, present: true}
}

// OfNullable returns an Option[T] that is absent when value is nil
// (in the sense of Go's nil: nil pointer, nil slice, nil map, nil
// channel, nil function, nil interface) and present otherwise.
//
// OfNullable is the right choice for pointer-like Ts. For value types
// (int, string, struct, …) use Of directly: there is no nil to test.
func OfNullable[T any](value T) Option[T] {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return Option[T]{}
	}
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		if rv.IsNil() {
			return Option[T]{}
		}
	}
	return Option[T]{value: value, present: true}
}

// IsPresent reports whether the Option holds a value.
func (o Option[T]) IsPresent() bool {
	return o.present
}

// IsEmpty reports whether the Option is absent.
// IsEmpty is the logical negation of IsPresent.
func (o Option[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the underlying value.
//
// Get panics if the Option is absent. Use Get together with
// IsPresent, or prefer OrElse / OrElseGet / OrElseThrow when the
// absent branch is expected.
func (o Option[T]) Get() T {
	if !o.present {
		panic("adt.Get on empty Option")
	}
	return o.value
}

// OrElse returns the underlying value if present, otherwise defaultValue.
//
// OrElse evaluates defaultValue eagerly. Use OrElseGet if computing the
// fallback is expensive.
func (o Option[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet returns the underlying value if present, otherwise the
// result of calling f. f is only invoked when the Option is absent.
func (o Option[T]) OrElseGet(f func() T) T {
	if o.present {
		return o.value
	}
	return f()
}

// OrElseThrow returns the underlying value if present, otherwise an
// error built from errMsg via errors.New. The second return value is
// nil when the Option is present.
//
// OrElseThrow does not panic; the name follows the Java Option
// convention of "throw if empty" but is adapted to Go's explicit
// (T, error) return shape.
func (o Option[T]) OrElseThrow(errMsg string) (T, error) {
	if o.present {
		return o.value, nil
	}
	var zero T
	return zero, errors.New(errMsg)
}

// IfPresent invokes f with the underlying value if the Option is
// present. If the Option is absent f is not called.
func (o Option[T]) IfPresent(f func(T)) {
	if o.present {
		f(o.value)
	}
}

// IfPresentOrElse invokes present with the value if the Option is
// present, otherwise invokes absent. Exactly one of the two callbacks
// runs.
func (o Option[T]) IfPresentOrElse(present func(T), absent func()) {
	if o.present {
		present(o.value)
		return
	}
	absent()
}

// Filter returns this Option if it is present and the value matches
// predicate, otherwise an absent Option. Filter does not call
// predicate when the receiver is absent.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if !o.present {
		return Option[T]{}
	}
	if predicate(o.value) {
		return o
	}
	return Option[T]{}
}

// Map applies f to the underlying value if present and returns a
// present Option[R] with the result. If the receiver is absent the
// returned Option[R] is also absent and f is not invoked.
//
// Map changes the element type via its own R parameter, which is only
// possible because Option is a concrete generic type (Go 1.27
// generic methods).
func (o Option[T]) Map[R any](f func(T) R) Option[R] {
	if !o.present {
		return Option[R]{}
	}
	return Option[R]{value: f(o.value), present: true}
}

// FlatMap applies f to the underlying value if present and returns
// the Option[R] produced by f. If the receiver is absent the
// returned Option[R] is absent and f is not invoked.
//
// FlatMap is the natural combinator for "if present, run a function
// that itself returns an Option".
func (o Option[T]) FlatMap[R any](f func(T) Option[R]) Option[R] {
	if !o.present {
		return Option[R]{}
	}
	return f(o.value)
}
