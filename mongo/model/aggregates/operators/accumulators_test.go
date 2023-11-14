package operators

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestSum(t *testing.T) {
	expected, _ := bson.Marshal(bson.E{Key: "sum", Value: bson.M{"$sum": 1}})
	actual, err := bson.Marshal(Sum("sum", 1))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.Equal(t, expected, actual)
}
