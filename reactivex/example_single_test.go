package reactivex_test

import (
	"errors"
	"fmt"

	"github.com/qianwj/typed/reactivex"
)

func ExampleSingle_Await() {
	// Single[T] is a typed Future: fn runs once, every subscriber
	// sees the same cached result.
	v, err := reactivex.NewSingle(func() (int, error) {
		return 42, nil
	}).Await()

	fmt.Println(v, err)
	// Output: 42 <nil>
}

func ExampleSingle_Map() {
	out, err := reactivex.NewSingle(func() (int, error) { return 3, nil }).
		Map(func(n int) string { return fmt.Sprintf("v=%d", n) }).
		Await()

	fmt.Println(out, err)
	// Output: v=3 <nil>
}

func ExampleSingle_Zip() {
	a := reactivex.NewSingle(func() (int, error) { return 3, nil })
	b := reactivex.NewSingle(func() (int, error) { return 4, nil })
	v, err := a.Zip(b, func(x, y int) int { return x*x + y*y }).Await()

	fmt.Println(v, err)
	// Output: 25 <nil>
}

func ExampleSingle_errorPropagation() {
	s := reactivex.NewSingle(func() (int, error) {
		return 0, errors.New("boom")
	})
	mapped := s.Map(func(n int) int { return n * 2 }) // Map's f does not run

	_, err := mapped.Await()
	fmt.Println(err)
	// Output: boom
}