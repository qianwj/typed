package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFrom(t *testing.T) {
	m := bson.M{
		"a": "a",
		"b": 1,
	}
	f := From(m)
	val, ok := f.get("a")
	if !ok || !assert.Equal(t, "a", val) {
		t.FailNow()
	}
}

func TestCast(t *testing.T) {
	m := bson.D{
		{"a", "a"},
		{"a", "b"},
		{"b", 1},
	}
	f := Cast(m)
	val, ok := f.get("a")
	if !ok || !assert.Equal(t, "b", val) {
		t.FailNow()
	}
}

func TestFilter_And(t *testing.T) {
	oid := primitive.NewObjectID()
	expected := bson.M{
		operator.And: []bson.M{
			{
				"a": "b",
				"b": 1,
				"c": "null",
				"d": "null",
				"e": oid.String(),
			},
		},
	}
	f := New().And(Cast(bson.D{
		{"a", "a"},
		{"a", "b"},
		{"b", 1},
		{"c", Null},
		{"d", Nil},
		{"e", oid},
	}))
	if !assert.Equal(t, expected, f.ToMap()) {
		t.FailNow()
	}
}
