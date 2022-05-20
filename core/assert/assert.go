package assert

import (
	"errors"
	"fmt"
	"github.com/qianwj/typed/core/object"
)

func Assert[T any](expect, actual T) error {
	if !object.Equals(expect, actual) {
		return errors.New(fmt.Sprintf("expect: %+v, but actual: %+v", expect, actual))
	}
	return nil
}
