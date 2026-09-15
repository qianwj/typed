package adt_test

import (
	"errors"
	"fmt"

	"github.com/qianwj/typed/adt"
)

func ExampleEither_typedErrorResult() {
	// Either[error, T] is the typed version of Go's idiomatic
	// (T, error) pair. Build one explicitly and consume with Fold.
	result := divide(10, 2)
	fmt.Println(result.Fold(
		func(e error) string { return "err: " + e.Error() },
		func(n int) string { return fmt.Sprintf("val: %d", n) },
	))

	bad := divide(10, 0)
	fmt.Println(bad.Fold(
		func(e error) string { return "err: " + e.Error() },
		func(n int) string { return fmt.Sprintf("val: %d", n) },
	))
	// Output:
	// val: 5
	// err: division by zero
}

func divide(a, b int) adt.Either[error, int] {
	if b == 0 {
		return adt.Left[error, int](errors.New("division by zero"))
	}
	return adt.Right[error, int](a / b)
}