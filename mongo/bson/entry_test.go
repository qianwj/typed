package bson

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_E(t *testing.T) {
	expected, _ := bson.Marshal(bson.E{Key: "abc", Value: "123"})
	actual, _ := bson.Marshal(E("abc", "123"))
	assert.Equal(t, expected, actual)
}

func TestEntry_Primitive(t *testing.T) {
	expected := bson.E{Key: "abc", Value: "123"}
	actual := E("abc", "123").Primitive()
	assert.Equal(t, expected, actual)
}
