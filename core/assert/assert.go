package assert

import (
	"errors"
	"fmt"
	"github.com/qianwj/typed/core/object"
	"testing"
)

func Assert[T any](expect, actual T) error {
	if !object.Equals[T](expect, actual) {
		return errors.New(fmt.Sprintf("expect: %+v, but actual: %+v", expect, actual))
	}
	return nil
}

func Equal[T any](t *testing.T, expect, actual T) {
	if !object.Equals[T](expect, actual) {
		t.Error(errors.New(fmt.Sprintf("expect: %+v, but actual: %+v", expect, actual)))
		t.FailNow()
	}
}
