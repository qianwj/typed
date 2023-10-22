package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAnd(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{operator.And, bson.A{
			bson.D{{"a", "b"}},
			bson.D{{"p", "q"}},
		}},
	})
	actual, _ := bson.Marshal(And(Eq("a", "b"), Eq("p", "q")))
	assert.Equal(t, expected, actual)
}
