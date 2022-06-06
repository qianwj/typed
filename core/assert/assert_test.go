package assert

import (
	"errors"
	"testing"
)

func TestAssert(t *testing.T) {
	i1, i2 := 1, 1
	if err := Assert[int](i1, i2); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestEqual(t *testing.T) {
	Equal[error](t, errors.New("test"), errors.New("test"))
	Equal[error](t, errors.New("test"), nil)
}
