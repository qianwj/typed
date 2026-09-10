// Package result provides Result[T], a small value type that
// represents the outcome of an operation that may fail: a success
// carrying a value of type T, or a failure carrying a non-nil error.
//
// Result is not an interface, which lets its methods declare their
// own type parameters (Map[R]) under Go 1.27's generic methods.
//
// Why not the standard `(T, error)` shape?
//
//   - An explicit Result makes the call site self-documenting and
//     lets combinators chain without scattering `if err != nil`
//     blocks at every step.
//   - The shape is intentionally simpler than a generic Result[T, E]:
//     the error channel is always a standard error, which keeps the
//     interop with regular Go code free of custom-error boxing.
//
// # Bridge to Optional
//
// Result.Optional returns an option.Optional[T]: a success becomes a
// present Optional, a failure becomes an absent one (and the error is
// dropped). This is the one-way bridge from this package to
// github.com/qianwj/typed/utils/option; Result depends on option,
// not the other way around.
package result

import (
	"github.com/qianwj/typed/utils/option"
)

// Result[T] is the outcome of an operation that may fail.
//
// A Result is either a success carrying a value of type T, or a
// failure carrying a non-nil error. The zero value of Result is a
// failure with a nil error; callers should always build results
// through Success / Failure and inspect them through the methods.
//
// Result is meant to be used where the caller wants to thread both
// the value and the error through a chain of combinators (Map,
// FlatMap, OrElse) without writing repeated `if err != nil` blocks.
type Result[T any] struct {
	value T
	err   error
}

// Success returns a successful Result carrying value.
func Success[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Failure returns a failed Result carrying err. Failure panics if
// err is nil; a Result with a nil error must be built via Success.
func Failure[T any](err error) Result[T] {
	if err == nil {
		panic("result.Failure: nil error")
	}
	return Result[T]{err: err}
}

// IsSuccess reports whether the Result is a success.
func (r Result[T]) IsSuccess() bool {
	return r.err == nil
}

// IsFailure reports whether the Result is a failure.
func (r Result[T]) IsFailure() bool {
	return r.err != nil
}

// Value returns the success value.
//
// Value panics if the Result is a failure; the panic value is the
// underlying error. Use Value together with IsSuccess, or prefer
// the OrElse family when the failure branch is expected.
func (r Result[T]) Value() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// Error returns the failure error, or nil if the Result is a success.
func (r Result[T]) Error() error {
	return r.err
}

// Unwrap returns the success value and a nil error if r is a
// success, otherwise the zero value of T and the failure error.
//
// Unwrap is the bridge from Result back to Go's standard (T, error)
// return shape. It is the right method to call when a chain of
// combinators ends and the caller wants to hand control back to
// ordinary Go error handling:
//
//	v, err := loadProfile(id).FlatMap(validate).Unwrap()
//	if err != nil {
//	    return fallbackProfile(id)
//	}
//	return v
//
// Unlike Value, Unwrap does not panic on failure: it surfaces the
// error to the caller, who can then decide what to do (log, fall
// back, wrap, return, …) using the full (T, error) vocabulary
// instead of being limited to a single fallback combinator.
func (r Result[T]) Unwrap() (T, error) {
	if r.err != nil {
		var zero T
		return zero, r.err
	}
	return r.value, nil
}

// Optional converts a successful Result into a present Optional and
// a failed Result into an absent Optional. The error is dropped, so
// Optional is only appropriate when the caller has already decided
// that the error channel can be discarded.
func (r Result[T]) Optional() option.Optional[T] {
	if r.err != nil {
		return option.Empty[T]()
	}
	return option.Of(r.value)
}

// OrElse returns the success value if the Result is a success,
// otherwise defaultValue. defaultValue is evaluated eagerly; use
// OrElseGet when it is expensive to compute.
func (r Result[T]) OrElse(defaultValue T) T {
	if r.err == nil {
		return r.value
	}
	return defaultValue
}

// OrElseGet returns the success value if the Result is a success,
// otherwise the result of calling f. f is only invoked on failure.
func (r Result[T]) OrElseGet(f func() T) T {
	if r.err == nil {
		return r.value
	}
	return f()
}

// Recover returns the success value if the Result is a success,
// otherwise the result of calling f with the failure error.
//
// Recover is the error-aware counterpart of OrElse and OrElseGet:
// while OrElse / OrElseGet hand the fallback no information, Recover
// gives the fallback the underlying error so it can decide what
// value to produce:
//
//	v := loadProfile(id).Recover(func(err error) Profile {
//	    if errors.Is(err, ErrNotFound) {
//	        return Profile{}               // missing is a zero profile
//	    }
//	    return Profile{Source: "cache"}   // anything else falls back to cache
//	})
//
// f is only invoked on failure. The success value is returned
// untouched. Unlike Unwrap, Recover always returns a T — it cannot
// re-fail, because that is what Unwrap is for.
func (r Result[T]) Recover(f func(error) T) T {
	if r.err == nil {
		return r.value
	}
	return f(r.err)
}

// Map applies f to the success value if the Result is a success and
// returns a new Result[R] with the result. A failure is propagated
// unchanged and f is not invoked.
//
// Map changes the success type via its own R parameter, which is only
// possible because Result is a concrete generic type.
func (r Result[T]) Map[R any](f func(T) R) Result[R] {
	if r.err != nil {
		return Result[R]{err: r.err}
	}
	return Result[R]{value: f(r.value)}
}

// FlatMap applies f to the success value if the Result is a success
// and returns the Result[R] produced by f. A failure is propagated
// unchanged and f is not invoked.
//
// FlatMap is the natural combinator for "if successful, run a
// function that itself may fail".
func (r Result[T]) FlatMap[R any](f func(T) Result[R]) Result[R] {
	if r.err != nil {
		return Result[R]{err: r.err}
	}
	return f(r.value)
}

// MapError applies f to the error if the Result is a failure and
// returns a new Result with the transformed error. A success is
// returned unchanged and f is not invoked.
//
// MapError is useful for wrapping low-level errors with a higher-level
// description before handing the Result up the call stack:
//
//	Failure[int](errors.New("disk full")).
//	    MapError(func(e error) error { return fmt.Errorf("save profile: %w", e) })
func (r Result[T]) MapError(f func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return Result[T]{err: f(r.err)}
}
