package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAll(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{Key: "a", Value: bson.D{{Key: operator.All, Value: []int{12, 14, 13}}}},
	})
	actual, _ := bson.Marshal(All[int]("a", []int{12, 14, 13}))
	assert.Equal(t, expected, actual)
}

func TestElemMatch(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{Key: "a", Value: bson.D{
			{Key: operator.ElemMatch, Value: bson.D{
				{Key: "b", Value: 1},
			}},
		}},
	})
	actual, _ := bson.Marshal(
		ElemMatch("a", Eq("b", 1)),
	)
	assert.Equal(t, expected, actual)
	expected, _ = bson.Marshal(bson.D{
		{Key: "a", Value: bson.D{
			{Key: operator.ElemMatch, Value: bson.D{
				{Key: operator.Gt, Value: 1},
				{Key: operator.Lt, Value: 2},
			}},
		}},
	})
	actual, _ = bson.Marshal(
		ElemMatchWithInterval("a", OpenInterval(1, 2)),
	)
	assert.Equal(t, expected, actual)
}
