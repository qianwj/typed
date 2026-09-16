package adt_test

import (
	"fmt"

	"github.com/qianwj/typed/adt"
)

func ExampleCast() {
	var value any = 42
	fmt.Println(adt.Cast[int](value).Map(func(n int) int { return n * 2 }).OrElse(0))
	fmt.Println(adt.Cast[int64](value).IsFailure())
	fmt.Println(adt.Cast[int](nil).IsFailure())
	fmt.Println(adt.Cast[*int]((*int)(nil)).IsSuccess())
	// Output:
	// 84
	// true
	// true
	// true
}
