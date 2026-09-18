package reactivex_test

import (
	"errors"
	"fmt"
	"github.com/qianwj/typed/adt/option"
	"github.com/qianwj/typed/adt/result"

	"github.com/qianwj/typed/reactivex"
)

func ExampleMaybe_Await() {
	// Maybe[T] adds a 'no value' terminal state to Single: the
	// inner Option distinguishes a value from empty completion.

	value, err := reactivex.NewMaybe(func() result.Result[option.Option[string]] {
		return result.Success(option.Of[string]("value-for-user:42"))
	}).Await().Unwrap()
	fmt.Println(value.OrElse("missing"), value.IsPresent(), err)
	// Output: value-for-user:42 true <nil>
}

func ExampleMaybe_completeWithoutValue() {
	// Empty completion is a successful result containing an empty Option.
	value, err := reactivex.NewMaybe(func() result.Result[option.Option[string]] {
		return result.Success(option.Empty[string]())
	}).Await().Unwrap()
	fmt.Printf("%q %v %v\n", value.OrElse(""), value.IsPresent(), err)
	// Output: "" false <nil>
}

func ExampleMaybe_errorPropagates() {
	m := reactivex.NewMaybe(func() result.Result[option.Option[int]] {
		return result.Failure[option.Option[int]](errors.New("missing"))
	})
	mapped := m.Map(func(int) string { return "should not run" })

	result := mapped.Await()
	fmt.Println(result.IsFailure(), result.Error())
	// Output: true missing
}
