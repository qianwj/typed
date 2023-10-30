package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestExists(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.D{{operator.Exists, true}}}})
	actual, _ := bson.Marshal(Exists("a", true))
	assert.Equal(t, expected, actual)
}

func TestType(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.D{{operator.Type, 9}}}})
	actual, _ := bson.Marshal(Type("a", DataTypeDate))
	assert.Equal(t, expected, actual)
}
