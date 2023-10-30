package projections

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestExcludeId(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "_id", Value: -1}})
	actual, _ := bson.Marshal(ExcludeId())
	assert.Equal(t, expected, actual)
}
