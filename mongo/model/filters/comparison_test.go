package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestEq(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: "b"}})
	actual, _ := bson.Marshal(Eq("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestNe(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$ne": "b"}}})
	actual, _ := bson.Marshal(Ne("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestGt(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$gt": "b"}}})
	actual, _ := bson.Marshal(Gt("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestGte(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$gte": "b"}}})
	actual, _ := bson.Marshal(Gte("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestLt(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$lt": "b"}}})
	actual, _ := bson.Marshal(Lt("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestLte(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$lte": "b"}}})
	actual, _ := bson.Marshal(Lte("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestWithInterval(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{operator.Gte: "b", operator.Lte: "c"}}})
	actual, _ := bson.Marshal(WithInterval("a", ClosedInterval("b", "c")))
	assert.Equal(t, expected, actual)
}

func TestIn(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$in": []string{"b", "c"}}}})
	actual, _ := bson.Marshal(In[string]("a", []string{"b", "c"}))
	assert.Equal(t, expected, actual)
}

func TestNin(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.M{"$nin": []string{"b", "c"}}}})
	actual, _ := bson.Marshal(Nin[string]("a", []string{"b", "c"}))
	assert.Equal(t, expected, actual)
}

func TestMixed(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{Key: "a", Value: "b"},
		{Key: "p", Value: bson.D{{operator.Ne, "q"}}},
		{Key: "r", Value: bson.D{{operator.Gt, 10}, {operator.Lt, 20}}},
		{Key: "d", Value: bson.D{{operator.Gte, 1}}},
		{Key: "c", Value: bson.M{operator.In: []string{"666"}}},
	})
	actual, _ := bson.Marshal(
		Eq("a", "b").
			Ne("p", "q").
			Gt("r", 10).
			WithInterval("r", OpenInterval(10, 20)).
			Gte("d", 1).
			In("c", []any{"666"}),
	)
	assert.Equal(t, expected, actual)
}
