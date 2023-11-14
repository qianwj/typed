package bson

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestArray_Append(t *testing.T) {
	expected := A(1, "2")
	expected = expected.Append(E("3", 4))
	actual := bson.A{
		1, "2", Entry{Key: "3", Value: 4},
	}
	assert.Equal(t, expected.Primitive(), actual)
}
