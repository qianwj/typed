package adt_test

import (
	"errors"
	"fmt"
	"strconv"

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

func ExampleEither_FlatMapRight() {
	parse := func(s string) adt.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return adt.Left[string, int]("invalid number")
		}
		return adt.Right[string, int](n)
	}
	for _, input := range []string{"21", "oops"} {
		message := adt.Right[string, string](input).
			FlatMapRight(parse).
			MapRight(func(n int) int { return n * 2 }).
			Fold(
				func(problem string) string { return "error: " + problem },
				func(n int) string { return fmt.Sprintf("value: %d", n) },
			)
		fmt.Println(message)
	}
	// Output:
	// value: 42
	// error: invalid number
}

func ExampleEither_FlatMapLeft() {
	result := adt.Left[string, int]("missing").
		FlatMapLeft(func(problem string) adt.Either[error, int] {
			if problem == "missing" {
				return adt.Right[error, int](21)
			}
			return adt.Left[error, int](errors.New(problem))
		}).
		FlatMapRight(func(n int) adt.Either[error, string] {
			return adt.Right[error, string](fmt.Sprintf("value: %d", n*2))
		})
	fmt.Println(result.Right().Get())
	// Output: value: 42
}
