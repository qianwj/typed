package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestBsonMap_MarshalJSON(t *testing.T) {
	m := NewMap()
	m.Put("aa", "bb")
	m.Put("bb", 1)
	m.Put("_id", primitive.NewObjectID())
	m.Put("time", time.Now())
	m.Put("cc", Null)
	m.Put("dd", bson.A{
		bson.D{
			{Key: "sdd", Value: "ddd"},
			{Key: "sdd2", Value: "ddd1"},
		},
		bson.M{
			"aaa": nil,
		},
		NewMap().Put("sub1", 222).Put("sub2", NewMap().Put("ddv", false)),
	})
	bytes, err := m.MarshalJSON()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("map: %+v", string(bytes))
}
