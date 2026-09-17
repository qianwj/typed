package reactivex_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/qianwj/typed/adt"
	"github.com/qianwj/typed/reactivex"
)

func ExampleSingle_ToFlowable() {
	source := reactivex.NewSingle(func() adt.Result[int] { return adt.Success(42) })
	values, err := source.ToFlowable().ToSlice(context.Background())
	fmt.Println(values, err)
	// Output: [42] <nil>
}

func ExampleMaybe_ToFlowable() {
	source := reactivex.NewMaybe(func() adt.Result[adt.Option[int]] {
		return adt.Success(adt.Empty[int]())
	})
	values, err := source.ToFlowable().ToSlice(context.Background())
	fmt.Println(values, err)
	// Output: [] <nil>
}

func ExampleFlowable_FirstElement() {
	value := reactivex.Just(3, 5).FirstElement(context.Background()).Await().Value()
	fmt.Println(value.Get())
	empty := reactivex.Just[int]().FirstElement(context.Background()).Await().Value()
	fmt.Println(empty.IsEmpty())
	// Output:
	// 3
	// true
}

func ExampleFlowable_FirstOrError() {
	first := reactivex.Just(3, 5).FirstOrError(context.Background()).Await()
	fmt.Println(first.Value())
	empty := reactivex.Just[int]().FirstOrError(context.Background()).Await()
	fmt.Println(errors.Is(empty.Error(), reactivex.ErrNoElements))
	// Output:
	// 3
	// true
}

func ExampleFlowable_FirstOrError_aggregate() {
	total := reactivex.Just(1, 2, 3).
		Reduce(func(sum, value int) int { return sum + value }).
		FirstOrError(context.Background())
	fmt.Println(total.Await().Value())
	// Output: 6
}
