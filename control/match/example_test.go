package match_test

import (
	"fmt"

	"github.com/qianwj/typed/control/match"
)

// ExampleMatch shows first-match-wins value matching: every When
// is (predicate, then) and the result is whatever the matching then
// returns. Otherwise supplies the fallback when no When fires.
func ExampleMatch() {
	classify := func(n int) string {
		return match.Match(n).
			When(func(v int) bool { return v < 0 }, func(int) string { return "negative" }).
			When(func(v int) bool { return v == 0 }, func(int) string { return "zero" }).
			OrElse("positive")
	}

	fmt.Println(classify(-5))
	fmt.Println(classify(0))
	fmt.Println(classify(7))
	// Output:
	// negative
	// zero
	// positive
}
