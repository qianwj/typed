package filters

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestExists(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.D{{Key: operator.Exists, Value: true}}}})
	actual, _ := bson.Marshal(Exists("a", true))
	assert.Equal(t, expected, actual)
}

func TestType(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.D{{Key: operator.Type, Value: 9}}}})
	actual, _ := bson.Marshal(Type("a", model.DataTypeDate))
	assert.Equal(t, expected, actual)
}
