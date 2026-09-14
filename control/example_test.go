package control_test

import (
	"errors"
	"fmt"

	"github.com/qianwj/typed/control"
)

// ExampleRepeat shows the imperative "do this N times" loop helper.
// For times <= 0 the body never runs.
func ExampleRepeat() {
	calls := 0
	control.Repeat(3, func() { calls++ })
	fmt.Println("calls:", calls)
	// Output: calls: 3
}

// ExampleRepeatE shows the error-aware variant that stops at the first
// non-nil error and reports how many iterations succeeded.
func ExampleRepeatE() {
	count, err := control.RepeatE(5, func() error {
		// Imagine an action that fails on its third attempt.
		return errors.New("transient")
	})
	fmt.Println("succeeded:", count, "err:", err)
	// Output: succeeded: 0 err: transient
}
