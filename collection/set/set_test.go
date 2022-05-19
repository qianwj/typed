package set

import (
	"github.com/qianwj/typed/collection/list"
	"testing"
)

func TestSet_NewSet(t *testing.T) {
	set := NewSet[string]()
	set.Add("a")
	set.Add("b")
	set.Add("b")
	if set.Size() != 2 {
		t.Errorf("want set size is 2, but found: %d", set.Size())
	}
}

func TestSet_AddAll(t *testing.T) {
	set := NewSet[int]()
	coll := list.NewArrayList[int]()
	coll.Add(1)
	coll.Add(1)
	set.AddAll(coll)
	if set.Size() != 1 {
		t.Errorf("want set size is 1, but found: %d", set.Size())
	}
}
