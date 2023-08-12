package docs

import (
	"github.com/qianwj/typed/mongo/bson"
	"testing"
)

func TestBsonMap_MarshalJSON(t *testing.T) {
	m := New()
	m.Put("aa", bson.String("bb"))
	m.Put("bb", bson.Int(1))
	m.Put("_id", bson.NewObjectId())
	m.Put("time", bson.Now())
	m.Put("cc", bson.Null{})
	m.Put("dd", bson.A{
		D{
			{Key: "sdd", Value: bson.String("ddd")},
			{Key: "sdd2", Value: "ddd1"},
		},
		New().Put("sub1", bson.Long(222)).Put("sub2", New().Put("ddv", bson.Bool(false))),
	})
	bytes, err := m.MarshalJSON()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("marshal: %s", string(bytes))
}
