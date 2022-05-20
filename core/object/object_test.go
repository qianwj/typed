package object

import "testing"

type testObject struct {
	Name  string
	Value int
}

func (t *testObject) Equals(other Object) bool {
	return t.Name == other.(*testObject).Name
}

type testObjectMember struct {
	id     int
	member *testObject
}

func (t *testObjectMember) Equals(other Object) bool {
	o := other.(*testObjectMember)
	return t.id == o.id && Equals[*testObjectMember](t.member, o.member)
}

func TestEquals(t *testing.T) {
	t1 := &testObject{
		Name:  "aa",
		Value: 1,
	}
	t2 := &testObject{
		Name:  "aa",
		Value: 0,
	}
	eq := Equals[testObject](t1, t2)
	if !eq {
		t.Error("want equals, but not equal")
		t.FailNow()
	}
	tm1 := &testObjectMember{id: 1, member: t1}
	tm2 := &testObjectMember{id: 1, member: t2}
	eq = Equals[testObjectMember](tm1, tm2)
	if !eq {
		t.Error("want equals, but not equal")
		t.FailNow()
	}
}
