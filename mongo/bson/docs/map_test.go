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
		},
	})
	bytes, err := m.MarshalJSON()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("marshal: %s", string(bytes))
}
