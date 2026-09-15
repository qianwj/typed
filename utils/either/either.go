// Package either provides Either[L, R], a tagged-union value type that
// holds exactly one of two values: a Left of type L or a Right of type R.
//
// The conventional reading is "Left = failure, Right = success" —
// for example Right[error, T] is the typed equivalent of Go's
// idiomatic `(T, error)` pair, but with explicit handling instead of
// a hidden zero-value convention. A typical use:
//
//	val, err := compute()
//	if err != nil {
//	    return either.Left[error, int](err)
//	}
//	return either.Right[error, int](val)
//
// Or, building it directly:
//
//	result := either.Right[error, int](42)
//	val, ok := result.Right().Get() // val == 42, ok == true
//
// # Why a tagged union instead of (T, error)?
//
// Go's `(T, error)` shape is convenient for one-shot call sites but
// awkward once a value has to flow through several intermediate
// functions, because every step has to thread an `error` variable and
// the "is this nil?" check has to be repeated. Either keeps both
// sides first-class so a function that may fail in two distinct ways
// can return:
//
//	either.Either[error, T]   // generic "value or error"
//	either.Either[errA, errB] // two error categories
//	either.Either[T, U]       // "either a T or a U"
//
// and the call site picks the right branch with IsLeft / IsRight or
// the Fold combinator.
//
// # Comparison with Optional[T]
//
// Optional[T] holds zero or one value of a single type. Either[L, R]
// holds exactly one value, but the value's type depends on which
// side was taken. Optional answers "is this here?"; Either answers
// "which of these two is it?".
//
// The safe accessors (Left / Right) return Optional[L] / Optional[R]
// so call sites can chain with the same combinators they already use
// for Optional.
package either

import "github.com/qianwj/typed/utils/option"

// Either[L, R] holds either a Left value of type L or a Right value of
// type R, exactly one of which is present.
//
// An Either is constructed through [Left] or [Right]; the zero value
// of Either is treated as the Left with both fields zero, which is
// rarely what callers want — prefer the explicit constructors.
// Either values must not be copied after creation; the safe usage is
// to obtain them through [Left] / [Right] and then consume them
// through the methods.
type Either[L, R any] struct {
	left    L
	right   R
	isRight bool
}

// Left returns an Either holding the given Left value. L is the
// "failure" or "first alternative" type by convention; nothing in the
// implementation enforces this, callers are free to flip the
// convention if it suits their domain.
func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{left: l}
}

// Right returns an Either holding the given Right value. R is the
// "success" or "second alternative" type by convention.
func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{right: r, isRight: true}
}

// IsLeft reports whether this Either holds a Left value.
func (e Either[L, R]) IsLeft() bool {
	return !e.isRight
}

// IsRight reports whether this Either holds a Right value.
func (e Either[L, R]) IsRight() bool {
	return e.isRight
}

// Left returns the Left value as an [option.Optional]. The result is
// present when IsLeft() is true and absent when IsRight() is true.
// Use Fold or the OrZero variant when a zero value is acceptable.
func (e Either[L, R]) Left() option.Optional[L] {
	if e.isRight {
		return option.Empty[L]()
	}
	return option.Of(e.left)
}

// Right returns the Right value as an [option.Optional]. The result is
// present when IsRight() is true and absent when IsLeft() is true.
func (e Either[L, R]) Right() option.Optional[R] {
	if e.isRight {
		return option.Of(e.right)
	}
	return option.Empty[R]()
}

// LeftOrZero returns the Left value, or the zero value of L when this
// Either is Right. Use this when the Left branch is informational
// only and the wrong-side fallback is acceptable.
//
// Prefer Left().Get() / Right().Get() (with an IsLeft / IsRight
// guard) when reading the wrong side would be a programming error.
func (e Either[L, R]) LeftOrZero() L {
	if e.isRight {
		var zero L
		return zero
	}
	return e.left
}

// RightOrZero returns the Right value, or the zero value of R when
// this Either is Left.
func (e Either[L, R]) RightOrZero() R {
	if e.isRight {
		return e.right
	}
	var zero R
	return zero
}

// MapLeft applies f to the Left value when this Either is Left and
// returns a new Either with the mapped Left. The Right branch is
// passed through unchanged.
func (e Either[L, R]) MapLeft[L2 any](f func(L) L2) Either[L2, R] {
	if e.isRight {
		return Right[L2, R](e.right)
	}
	return Left[L2, R](f(e.left))
}

// MapRight applies f to the Right value when this Either is Right
// and returns a new Either with the mapped Right. The Left branch is
// passed through unchanged.
func (e Either[L, R]) MapRight[R2 any](f func(R) R2) Either[L, R2] {
	if e.isRight {
		return Right[L, R2](f(e.right))
	}
	return Left[L, R2](e.left)
}

// Fold applies onLeft when this Either is Left and onRight when it is
// Right, returning the result. Exactly one of the two callbacks is
// invoked; this is the canonical way to consume an Either when both
// branches produce the same result type.
func (e Either[L, R]) Fold[T any](onLeft func(L) T, onRight func(R) T) T {
	if e.isRight {
		return onRight(e.right)
	}
	return onLeft(e.left)
}