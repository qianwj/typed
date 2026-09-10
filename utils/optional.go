package utils

// Optional[T] is a container that may or may not hold a value of type T.
//
// An Optional is either present (carrying a T) or absent (no value).
// The zero value of Optional is absent and is equivalent to Empty[T]().
// Optional values must not be copied after creation; the safe usage is
// to obtain them through one of the constructors (Empty, Of, OfNullable)
// and then consume them through the methods.
type Optional[T any] struct {
	value   T
	present bool
}

// Empty returns an absent Optional[T].
//
// Calling Empty on a value type such as int is the idiomatic way to say
// "no result" and is preferred over using a sentinel like 0 or "".
func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Of returns a present Optional[T] holding value.
//
// Of panics if value is a typed nil (e.g. a nil *Foo passed as a
// pointer-typed T). Use OfNullable when the underlying T is a
// reference type and nil is a meaningful value.
func Of[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// OfNullable returns an Optional[T] that is absent when value is nil
// (in the sense of Go's nil: nil pointer, nil slice, nil map, nil
// channel, nil function, nil interface) and present otherwise.
//
// OfNullable is the right choice for pointer-like Ts. For value types
// (int, string, struct, …) use Of directly: there is no nil to test.
func OfNullable[T any](value T) Optional[T] {
	if isNil(value) {
		return Optional[T]{}
	}
	return Optional[T]{value: value, present: true}
}

// IsPresent reports whether the Optional holds a value.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty reports whether the Optional is absent.
// IsEmpty is the logical negation of IsPresent.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the underlying value.
//
// Get panics if the Optional is absent. Use Get together with
// IsPresent, or prefer OrElse / OrElseGet / OrElseThrow when the
// absent branch is expected.
func (o Optional[T]) Get() T {
	if !o.present {
		panic("Optional.Get on empty Optional")
	}
	return o.value
}

// OrElse returns the underlying value if present, otherwise defaultValue.
//
// OrElse evaluates defaultValue eagerly. Use OrElseGet if computing the
// fallback is expensive.
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet returns the underlying value if present, otherwise the
// result of calling f. f is only invoked when the Optional is absent.
func (o Optional[T]) OrElseGet(f func() T) T {
	if o.present {
		return o.value
	}
	return f()
}

// OrElseThrow returns the underlying value if present, otherwise an
// error built from errMsg via errors.New. The second return value is
// nil when the Optional is present.
//
// OrElseThrow does not panic; the name follows the Java Optional
// convention of "throw if empty" but is adapted to Go's explicit
// (T, error) return shape.
func (o Optional[T]) OrElseThrow(errMsg string) (T, error) {
	if o.present {
		return o.value, nil
	}
	var zero T
	return zero, errFromMsg(errMsg)
}

// IfPresent invokes f with the underlying value if the Optional is
// present. If the Optional is absent f is not called.
func (o Optional[T]) IfPresent(f func(T)) {
	if o.present {
		f(o.value)
	}
}

// IfPresentOrElse invokes present with the value if the Optional is
// present, otherwise invokes absent. Exactly one of the two callbacks
// runs.
func (o Optional[T]) IfPresentOrElse(present func(T), absent func()) {
	if o.present {
		present(o.value)
		return
	}
	absent()
}

// Filter returns this Optional if it is present and the value matches
// predicate, otherwise an absent Optional. Filter does not call
// predicate when the receiver is absent.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if !o.present {
		return Optional[T]{}
	}
	if predicate(o.value) {
		return o
	}
	return Optional[T]{}
}

// Map applies f to the underlying value if present and returns a
// present Optional[R] with the result. If the receiver is absent the
// returned Optional[R] is also absent and f is not invoked.
//
// Map changes the element type via its own R parameter, which is only
// possible because Optional is a concrete generic type (Go 1.27
// generic methods).
func (o Optional[T]) Map[R any](f func(T) R) Optional[R] {
	if !o.present {
		return Optional[R]{}
	}
	return Optional[R]{value: f(o.value), present: true}
}

// FlatMap applies f to the underlying value if present and returns
// the Optional[R] produced by f. If the receiver is absent the
// returned Optional[R] is absent and f is not invoked.
//
// FlatMap is the natural combinator for "if present, run a function
// that itself returns an Optional".
func (o Optional[T]) FlatMap[R any](f func(T) Optional[R]) Optional[R] {
	if !o.present {
		return Optional[R]{}
	}
	return f(o.value)
}
