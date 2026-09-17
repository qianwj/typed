package reactivex_test

import (
	"errors"
	"fmt"
	"github.com/qianwj/typed/adt"

	"github.com/qianwj/typed/reactivex"
)

func ExampleMaybe_Await() {
	// Maybe[T] adds a 'no value' terminal state to Single: the
	// inner Option distinguishes a value from empty completion.

	value, err := reactivex.NewMaybe(func() adt.Result[adt.Option[string]] {
		return adt.Success(adt.Of[string]("value-for-user:42"))
	}).Await().Unwrap()
	fmt.Println(value.OrElse("missing"), value.IsPresent(), err)
	// Output: value-for-user:42 true <nil>
}

func ExampleMaybe_completeWithoutValue() {
	// Empty completion is a successful result containing an empty Option.
	value, err := reactivex.NewMaybe(func() adt.Result[adt.Option[string]] {
		return adt.Success(adt.Empty[string]())
	}).Await().Unwrap()
	fmt.Printf("%q %v %v\n", value.OrElse(""), value.IsPresent(), err)
	// Output: "" false <nil>
}

func ExampleMaybe_errorPropagates() {
	m := reactivex.NewMaybe(func() adt.Result[adt.Option[int]] {
		return adt.Failure[adt.Option[int]](errors.New("missing"))
	})
	mapped := m.Map(func(int) string { return "should not run" })

	result := mapped.Await()
	fmt.Println(result.IsFailure(), result.Error())
	// Output: true missing
}
