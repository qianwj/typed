package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestElemMatch(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{"a", bson.D{
			{operator.ElemMatch, bson.D{
				{"b", 1},
			}},
		}},
	})
	actual, _ := bson.Marshal(
		ElemMatch("a", Eq("b", 1)),
	)
	assert.Equal(t, expected, actual)
	expected, _ = bson.Marshal(bson.D{
		{"a", bson.D{
			{operator.ElemMatch, bson.D{
				{operator.Gt, 1},
				{operator.Lt, 2},
			}},
		}},
	})
	actual, _ = bson.Marshal(
		ElemMatchWithInterval("a", OpenInterval(1, 2)),
	)
	assert.Equal(t, expected, actual)
}
