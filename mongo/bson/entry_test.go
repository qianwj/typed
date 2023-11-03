package bson

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestEntry_MarshalBJSON(t *testing.T) {
	expected, _ := bson.Marshal(bson.E{Key: "abc", Value: "123"})
	actual, _ := bson.Marshal(Entry{Key: "abc", Value: "123"})
	assert.Equal(t, expected, actual)
}

func TestEntry_UnmarshalBJSON(t *testing.T) {
	expected := Entry{Key: "abc", Value: "123"}
	bytes, _ := bson.Marshal(bson.E{Key: "abc", Value: "123"})
	var actual Entry
	err := bson.Unmarshal(bytes, &actual)
	if err != nil {
		t.Errorf("unmarshal entry failed. %+v", err)
		t.FailNow()
	}
	assert.Equal(t, expected, actual)
}
