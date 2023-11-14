package aggregates

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestGroup(t *testing.T) {
	eDoc := bson.D{{Key: "$group", Value: bson.D{{"_id", "$state"}}}}
	t.Logf("e stages: %+v\r\n", eDoc)
	expected, e := bson.Marshal(eDoc)
	if e != nil {
		t.Error(e)
	}
	aDoc := Group("$state").Stages()[0]
	t.Logf("a stages: %+v\r\n", aDoc)
	actual, e := bson.Marshal(aDoc)
	if e != nil {
		t.Error(e)
	}
	assert.Equal(t, expected, actual)
}
