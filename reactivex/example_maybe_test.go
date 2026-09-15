package reactivex_test

import (
	"errors"
	"fmt"

	"github.com/qianwj/typed/reactivex"
)

func ExampleMaybe_Await() {
	// Maybe[T] adds a 'no value' terminal state to Single: the
	// presence flag distinguishes success-with-value from
	// complete-without-value.

	v, present, err := reactivex.NewMaybe(func() (string, bool, error) {
		return "value-for-user:42", true, nil
	}).Await()
	fmt.Println(v, present, err)
	// Output: value-for-user:42 true <nil>
}

func ExampleMaybe_completeWithoutValue() {
	// Terminal 'no value' state: present=false, err=nil.
	v, present, err := reactivex.NewMaybe(func() (string, bool, error) {
		return "", false, nil
	}).Await()
	fmt.Printf("%q %v %v\n", v, present, err)
	// Output: "" false <nil>
}

func ExampleMaybe_errorPropagates() {
	m := reactivex.NewMaybe(func() (int, bool, error) {
		return 0, false, errors.New("missing")
	})
	mapped := m.Map(func(int) string { return "should not run" })

	_, present, err := mapped.Await()
	fmt.Println(present, err)
	// Output: false missing
}