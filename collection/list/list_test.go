package list

import "testing"

func TestArrayList_Contains(t *testing.T) {
	l := NewArrayList[int]()
	l.Add(1)
	l.Add(2)
	if !l.Contains(1) {
		t.Errorf("expact contains 1, but not contains this element")
		t.FailNow()
	}
}
