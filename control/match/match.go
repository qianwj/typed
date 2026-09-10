// Package match provides a small, type-safe control-flow API for
// first-match-wins pattern matching.
//
// The package covers two matching modes:
//
//   - value matching with Pattern[T], useful for equality, membership,
//     range, and arbitrary predicate matches
//   - type matching over an any value, useful when a value may hold one
//     of several concrete types
//
// The API is deliberately explicit: every case returns the same result
// type R, cases are evaluated in declaration order, and once a case
// matches, later patterns and handlers are not evaluated.
//
// The package requires Go 1.27 or later because the fluent chain methods
// declare their own type parameters.
package match

import (
	"cmp"

	"github.com/qianwj/typed/utils/objects"
)

// Pattern decides whether a value should match a case.
type Pattern[T any] interface {
	Match(T) bool
}

// PatternFunc adapts a function into a Pattern.
type PatternFunc[T any] func(T) bool

// Match reports whether value satisfies the pattern.
func (p PatternFunc[T]) Match(value T) bool {
	if p == nil {
		panic("match.PatternFunc: nil pattern")
	}
	return p(value)
}

// Predicate returns a Pattern backed by predicate.
func Predicate[T any](predicate func(T) bool) Pattern[T] {
	if predicate == nil {
		panic("match.Predicate: nil predicate")
	}
	return PatternFunc[T](predicate)
}

// Any returns a Pattern that matches every value.
func Any[T any]() Pattern[T] {
	return PatternFunc[T](func(T) bool {
		return true
	})
}

// Eq returns a Pattern that matches values equal to want according to
// objects.Equals.
//
// Eq accepts any T, including structs containing slices or maps. The shared
// equality helper performs deep comparison and honors an Equal method when the
// value type provides one; use Predicate when the domain needs a per-pattern
// rule.
func Eq[T any](want T) Pattern[T] {
	return PatternFunc[T](func(value T) bool {
		return objects.Equals(value, want)
	})
}

// In returns a Pattern that matches any value contained in values.
func In[T comparable](values ...T) Pattern[T] {
	set := make(map[T]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return PatternFunc[T](func(value T) bool {
		_, ok := set[value]
		return ok
	})
}

// Between returns a Pattern that matches values in the inclusive range
// [min, max].
func Between[T cmp.Ordered](min, max T) Pattern[T] {
	return PatternFunc[T](func(value T) bool {
		return value >= min && value <= max
	})
}

// Not returns a Pattern that negates pattern.
func Not[T any](pattern Pattern[T]) Pattern[T] {
	requirePattern(pattern, "Not")
	return PatternFunc[T](func(value T) bool {
		return !pattern.Match(value)
	})
}

// And returns a Pattern that matches only when all patterns match. With no
// patterns it matches every value.
func And[T any](patterns ...Pattern[T]) Pattern[T] {
	return PatternFunc[T](func(value T) bool {
		for _, pattern := range patterns {
			requirePattern(pattern, "And")
			if !pattern.Match(value) {
				return false
			}
		}
		return true
	})
}

// Or returns a Pattern that matches when at least one pattern matches. With
// no patterns it matches no value.
func Or[T any](patterns ...Pattern[T]) Pattern[T] {
	return PatternFunc[T](func(value T) bool {
		for _, pattern := range patterns {
			requirePattern(pattern, "Or")
			if pattern.Match(value) {
				return true
			}
		}
		return false
	})
}

// Matcher is a value matcher for one input value of type T.
type Matcher[T any] struct {
	value T
}

// Value starts value-based pattern matching.
func Value[T any](value T) Matcher[T] {
	return Matcher[T]{value: value}
}

// Match is an alias for Value and starts value-based pattern matching.
func Match[T any](value T) Matcher[T] {
	return Value(value)
}

// Case evaluates pattern against the input value. If it matches, then is
// called and the returned Chain is already matched.
func (m Matcher[T]) Case[R any](pattern Pattern[T], then func(T) R) Chain[T, R] {
	return applyValueCase[T, R](Chain[T, R]{value: m.value}, pattern, then)
}

// When is shorthand for Case(Predicate(predicate), then).
func (m Matcher[T]) When[R any](predicate func(T) bool, then func(T) R) Chain[T, R] {
	return applyValueWhen[T, R](Chain[T, R]{value: m.value}, predicate, then)
}

// Default returns the result of then for a matcher with no prior cases.
func (m Matcher[T]) Default[R any](then func(T) R) R {
	if then == nil {
		panic("match.Default: nil handler")
	}
	return then(m.value)
}

// OrElse returns fallback for a matcher with no prior cases.
func (m Matcher[T]) OrElse[R any](fallback R) R {
	return fallback
}

// OrElseGet returns the result of fallback for a matcher with no prior cases.
func (m Matcher[T]) OrElseGet[R any](fallback func(T) R) R {
	if fallback == nil {
		panic("match.OrElseGet: nil handler")
	}
	return fallback(m.value)
}

// Chain is the in-progress result of value-based pattern matching.
type Chain[T any, R any] struct {
	value   T
	matched bool
	result  R
}

// Case evaluates pattern only if no previous case matched.
func (c Chain[T, R]) Case(pattern Pattern[T], then func(T) R) Chain[T, R] {
	return applyValueCase[T, R](c, pattern, then)
}

func applyValueCase[T any, R any](c Chain[T, R], pattern Pattern[T], then func(T) R) Chain[T, R] {
	if c.matched {
		return c
	}
	requirePattern(pattern, "Case")
	if pattern.Match(c.value) {
		if then == nil {
			panic("match.Case: nil handler")
		}
		c.result = then(c.value)
		c.matched = true
	}
	return c
}

// When is shorthand for Case(Predicate(predicate), then).
func (c Chain[T, R]) When(predicate func(T) bool, then func(T) R) Chain[T, R] {
	return applyValueWhen[T, R](c, predicate, then)
}

func applyValueWhen[T any, R any](c Chain[T, R], predicate func(T) bool, then func(T) R) Chain[T, R] {
	if c.matched {
		return c
	}
	return applyValueCase[T, R](c, Predicate(predicate), then)
}

// Default returns the matched result, or evaluates then when no case matched.
func (c Chain[T, R]) Default(then func(T) R) R {
	if c.matched {
		return c.result
	}
	if then == nil {
		panic("match.Default: nil handler")
	}
	return then(c.value)
}

// OrElse returns the matched result, or fallback when no case matched.
func (c Chain[T, R]) OrElse(fallback R) R {
	if c.matched {
		return c.result
	}
	return fallback
}

// OrElseGet returns the matched result, or evaluates fallback when no case
// matched.
func (c Chain[T, R]) OrElseGet(fallback func(T) R) R {
	if c.matched {
		return c.result
	}
	if fallback == nil {
		panic("match.OrElseGet: nil handler")
	}
	return fallback(c.value)
}

// Matched reports whether any case has matched.
func (c Chain[T, R]) Matched() bool {
	return c.matched
}

// Unwrap returns the matched result and true, or the zero value of R and
// false when no case matched.
func (c Chain[T, R]) Unwrap() (R, bool) {
	if c.matched {
		return c.result, true
	}
	var zero R
	return zero, false
}

// TypeMatcher is a matcher for dynamic type checks over an any value.
type TypeMatcher struct {
	value any
}

// Type starts type-based pattern matching.
func Type(value any) TypeMatcher {
	return TypeMatcher{value: value}
}

// Case matches when the input value has dynamic type T, or when the dynamic
// type implements T if T is an interface type.
func (m TypeMatcher) Case[T any, R any](then func(T) R) TypeChain[R] {
	return applyTypeCase[T, R](TypeChain[R]{value: m.value}, then)
}

// CaseWhen is a type case with an additional guard. then runs only when the
// dynamic type matches T and guard returns true.
func (m TypeMatcher) CaseWhen[T any, R any](guard func(T) bool, then func(T) R) TypeChain[R] {
	return applyTypeCaseWhen[T, R](TypeChain[R]{value: m.value}, guard, then)
}

// Nil matches an untyped nil interface or a typed nil pointer, map, slice,
// channel, function, or interface.
func (m TypeMatcher) Nil[R any](then func() R) TypeChain[R] {
	return applyTypeNil(TypeChain[R]{value: m.value}, then)
}

// When matches on an arbitrary predicate over the raw input value.
func (m TypeMatcher) When[R any](predicate func(any) bool, then func(any) R) TypeChain[R] {
	return applyTypeWhen(TypeChain[R]{value: m.value}, predicate, then)
}

// Default returns the result of then for a matcher with no prior cases.
func (m TypeMatcher) Default[R any](then func(any) R) R {
	if then == nil {
		panic("match.Type.Default: nil handler")
	}
	return then(m.value)
}

// OrElse returns fallback for a matcher with no prior cases.
func (m TypeMatcher) OrElse[R any](fallback R) R {
	return fallback
}

// OrElseGet returns the result of fallback for a matcher with no prior cases.
func (m TypeMatcher) OrElseGet[R any](fallback func(any) R) R {
	if fallback == nil {
		panic("match.Type.OrElseGet: nil handler")
	}
	return fallback(m.value)
}

// TypeChain is the in-progress result of type-based pattern matching.
type TypeChain[R any] struct {
	value   any
	matched bool
	result  R
}

// Case evaluates the type case only if no previous case matched.
func (c TypeChain[R]) Case[T any](then func(T) R) TypeChain[R] {
	return applyTypeCase[T, R](c, then)
}

func applyTypeCase[T any, R any](c TypeChain[R], then func(T) R) TypeChain[R] {
	if c.matched {
		return c
	}
	typed, ok := c.value.(T)
	if ok {
		if then == nil {
			panic("match.Type.Case: nil handler")
		}
		c.result = then(typed)
		c.matched = true
	}
	return c
}

// CaseWhen evaluates the type case and guard only if no previous case
// matched.
func (c TypeChain[R]) CaseWhen[T any](guard func(T) bool, then func(T) R) TypeChain[R] {
	return applyTypeCaseWhen[T, R](c, guard, then)
}

func applyTypeCaseWhen[T any, R any](c TypeChain[R], guard func(T) bool, then func(T) R) TypeChain[R] {
	if c.matched {
		return c
	}
	typed, ok := c.value.(T)
	if !ok {
		return c
	}
	if guard == nil {
		panic("match.Type.CaseWhen: nil guard")
	}
	if guard(typed) {
		if then == nil {
			panic("match.Type.CaseWhen: nil handler")
		}
		c.result = then(typed)
		c.matched = true
	}
	return c
}

// Nil evaluates then only when the input is nil.
func (c TypeChain[R]) Nil(then func() R) TypeChain[R] {
	return applyTypeNil(c, then)
}

func applyTypeNil[R any](c TypeChain[R], then func() R) TypeChain[R] {
	if c.matched {
		return c
	}
	if objects.IsNil(c.value) {
		if then == nil {
			panic("match.Type.Nil: nil handler")
		}
		c.result = then()
		c.matched = true
	}
	return c
}

// When evaluates predicate against the raw input value only if no previous
// case matched.
func (c TypeChain[R]) When(predicate func(any) bool, then func(any) R) TypeChain[R] {
	return applyTypeWhen(c, predicate, then)
}

func applyTypeWhen[R any](c TypeChain[R], predicate func(any) bool, then func(any) R) TypeChain[R] {
	if c.matched {
		return c
	}
	if predicate == nil {
		panic("match.Type.When: nil predicate")
	}
	if predicate(c.value) {
		if then == nil {
			panic("match.Type.When: nil handler")
		}
		c.result = then(c.value)
		c.matched = true
	}
	return c
}

// Default returns the matched result, or evaluates then when no case matched.
func (c TypeChain[R]) Default(then func(any) R) R {
	if c.matched {
		return c.result
	}
	if then == nil {
		panic("match.Type.Default: nil handler")
	}
	return then(c.value)
}

// OrElse returns the matched result, or fallback when no case matched.
func (c TypeChain[R]) OrElse(fallback R) R {
	if c.matched {
		return c.result
	}
	return fallback
}

// OrElseGet returns the matched result, or evaluates fallback when no case
// matched.
func (c TypeChain[R]) OrElseGet(fallback func(any) R) R {
	if c.matched {
		return c.result
	}
	if fallback == nil {
		panic("match.Type.OrElseGet: nil handler")
	}
	return fallback(c.value)
}

// Matched reports whether any case has matched.
func (c TypeChain[R]) Matched() bool {
	return c.matched
}

// Unwrap returns the matched result and true, or the zero value of R and
// false when no case matched.
func (c TypeChain[R]) Unwrap() (R, bool) {
	if c.matched {
		return c.result, true
	}
	var zero R
	return zero, false
}

func requirePattern[T any](pattern Pattern[T], op string) {
	if pattern == nil {
		panic("match." + op + ": nil pattern")
	}
}
