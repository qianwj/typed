package assert

import "testing"

func TestAssert(t *testing.T) {
	i1, i2 := 1, 1
	if err := Assert[int](i1, i2); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
