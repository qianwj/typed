package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAnd(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{Key: operator.And, Value: bson.A{
			bson.D{{Key: "a", Value: "b"}},
			bson.D{{Key: "p", Value: "q"}},
		}},
	})
	actual, _ := bson.Marshal(And(Eq("a", "b"), Eq("p", "q")))
	assert.Equal(t, expected, actual)
}

func TestOr(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{Key: operator.Or, Value: bson.A{
			bson.D{{Key: "a", Value: "b"}},
			bson.D{{Key: "p", Value: "q"}},
		}},
	})
	actual, _ := bson.Marshal(Or(Eq("a", "b"), Eq("p", "q")))
	assert.Equal(t, expected, actual)
}
