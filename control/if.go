package control

// If returns onTrue when condition is true, otherwise onFalse.
// Both arguments are evaluated before If is called, unlike Java's ternary
// operator. This is Go's language-defined argument evaluation rule, enforced
// by the compiler, not a compiler bug or an implementation choice in If.
// An ordinary library function cannot defer argument evaluation; a direct
// ternary expression with short-circuit semantics requires language/compiler
// support. Use IfGet with callbacks to evaluate only the selected branch.
//
// In particular, If(user != nil, user.Name, "anonymous") still panics when
// user is nil. For an optional value, use adt.OfNullable from
// github.com/qianwj/typed/adt instead:
//
//	name := adt.OfNullable(user).
//		Map(func(u *User) string { return u.Name }).
//		OrElse("anonymous")
//
// Map skips its callback when the Option is absent. Use OrElseGet instead
// of OrElse when the fallback should also be evaluated only when needed.
func If[T any](condition bool, onTrue, onFalse T) T {
	if condition {
		return onTrue
	}
	return onFalse
}

// IfGet calls onTrue when condition is true, otherwise onFalse, and returns
// its result. The selected callback runs exactly once; the other is not called.
// A nil selected callback panics; an unselected callback may be nil.
//
// Go evaluates arguments before entering a function. Passing callbacks defers
// their bodies, allowing IfGet to select a branch before computing its value:
//
//	name := IfGet(user != nil,
//		func() string { return user.Name },
//		func() string { return "anonymous" },
//	)
//
// This provides the selected-branch evaluation of Java's ternary operator
// without special compiler support. Expressions that create the callbacks
// are still evaluated eagerly; put deferred work inside the callback bodies.
// For optional values, adt.OfNullable(...).Map(...).OrElse(...) is another
// option; see If for an example.
func IfGet[T any](condition bool, onTrue, onFalse func() T) T {
	if condition {
		return onTrue()
	}
	return onFalse()
}
