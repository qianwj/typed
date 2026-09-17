package reactivex_test

import (
	"errors"
	"fmt"
	"github.com/qianwj/typed/adt"

	"github.com/qianwj/typed/reactivex"
)

func ExampleSingle_Await() {
	// Single[T] is a typed Future: fn runs once, every subscriber
	// sees the same cached result.
	v, err := reactivex.NewSingle(func() adt.Result[int] {
		return adt.Success[int](42)
	}).Await().Unwrap()

	fmt.Println(v, err)
	// Output: 42 <nil>
}

func ExampleSingle_Map() {
	out, err := reactivex.NewSingle(func() adt.Result[int] { return adt.Success[int](3) }).
		Map(func(n int) string { return fmt.Sprintf("v=%d", n) }).
		Await().Unwrap()

	fmt.Println(out, err)
	// Output: v=3 <nil>
}

func ExampleSingle_Zip() {
	a := reactivex.NewSingle(func() adt.Result[int] { return adt.Success[int](3) })
	b := reactivex.NewSingle(func() adt.Result[int] { return adt.Success[int](4) })
	v, err := a.Zip(b, func(x, y int) int { return x*x + y*y }).Await().Unwrap()

	fmt.Println(v, err)
	// Output: 25 <nil>
}

func ExampleSingle_errorPropagation() {
	s := reactivex.NewSingle(func() adt.Result[int] {
		return adt.Failure[int](errors.New("boom"))
	})
	mapped := s.Map(func(n int) int { return n * 2 }) // Map's f does not run

	_, err := mapped.Await().Unwrap()
	fmt.Println(err)
	// Output: boom
}
