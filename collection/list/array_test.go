package list

import (
	"github.com/qianwj/typed/core/assert"
	"testing"
)

func TestArrayList_Add(t *testing.T) {
	l := NewArrayList[int]()
	l.Add(1)
	if err := assert.Assert[int](1, l.Get(0)); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
