package filters

import (
	"encoding/json"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFrom(t *testing.T) {
	m := rawbson.M{
		"a": "a",
		"b": 1,
	}
	f := From(m)
	val, ok := f.Get("a")
	if !ok || !assert.Equal(t, "a", val) {
		t.FailNow()
	}
}

func TestCast(t *testing.T) {
	m := rawbson.D{
		{"a", "a"},
		{"a", "b"},
		{"b", 1},
	}
	f := Cast(m)
	val, ok := f.Get("a")
	if !ok || !assert.Equal(t, "b", val) {
		t.FailNow()
	}
}

func TestFilter_And(t *testing.T) {
	oid := primitive.NewObjectID()
	expected := map[string]any{
		operator.And: []*bson.Map{
			bson.FromM(rawbson.M{
				"a": "b",
				"b": 1,
				"c": nil,
				"d": nil,
				"e": oid,
			}),
		},
	}
	eBytes, _ := json.Marshal(expected)
	t.Logf("expect: %s", string(eBytes))
	actual := New().And(Cast(rawbson.D{
		{"a", "a"},
		{"a", "b"},
		{"b", 1},
		{"c", bson.Null},
		{"d", bson.Nil},
		{"e", oid},
	}))
	aBytes, _ := json.Marshal(actual)
	t.Logf("actual: %s", string(aBytes))
	if !assert.Equal(t, string(eBytes), string(aBytes)) {
		t.FailNow()
	}
}
