package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestEq(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", "b"}})
	actual, _ := bson.Marshal(Eq("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestNe(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$ne": "b"}}})
	actual, _ := bson.Marshal(Ne("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestGt(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$gt": "b"}}})
	actual, _ := bson.Marshal(Gt("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestGte(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$gte": "b"}}})
	actual, _ := bson.Marshal(Gte("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestLt(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$lt": "b"}}})
	actual, _ := bson.Marshal(Lt("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestLte(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$lte": "b"}}})
	actual, _ := bson.Marshal(Lte("a", "b"))
	assert.Equal(t, expected, actual)
}

func TestIn(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$in": []string{"b", "c"}}}})
	actual, _ := bson.Marshal(In("a", []any{"b", "c"}))
	assert.Equal(t, expected, actual)
}

func TestNin(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{"a", bson.M{"$nin": []string{"b", "c"}}}})
	actual, _ := bson.Marshal(Nin("a", []any{"b", "c"}))
	assert.Equal(t, expected, actual)
}

func TestAll(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{
		{"a", "b"},
		{"p", bson.M{"$ne": "q"}},
		{"r", bson.M{"$gt": 10, "$lt": 20}},
		{"d", bson.M{operator.Gte: 1}},
		{"c", bson.M{operator.In: []string{"666"}}},
	})
	actual, _ := bson.Marshal(
		Eq("a", "b").
			Ne("p", "q").
			Gt("r", 10).Lt("r", 20).
			Gte("d", 1).
			In("c", []any{"666"}),
	)
	assert.Equal(t, expected, actual)
}
