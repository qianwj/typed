package update

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestSet(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Set: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(Set("a", 1).Set("b", 2))
	assert.Equal(t, expected, actual)
}

func TestSetAll(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Set: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(SetAll(bson.M{
		"a": 1,
		"b": 2,
	}))
	assert.Equal(t, expected, actual)
}

func TestSetOnInsert(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.SetOnInsert: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(SetOnInsert("a", 1).SetOnInsert("b", 2))
	assert.Equal(t, expected, actual)
}

func TestSetOnInsertAll(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.SetOnInsert: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(SetOnInsertAll(bson.M{
		"a": 1,
		"b": 2,
	}))
	assert.Equal(t, expected, actual)
}

func TestUnset(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Unset: bson.M{
			"b": "",
			"c": "",
		},
	})
	actual, _ := bson.Marshal(Unset("b", "c"))
	assert.Equal(t, expected, actual)
}

func TestInc(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Inc: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(Inc("a", 1).Inc("b", 2))
	assert.Equal(t, expected, actual)
}

func TestAddToSet(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.AddToSet: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(AddToSet("a", 1).AddToSet("b", 2))
	assert.Equal(t, expected, actual)
}

func TestAddEachToSet(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.AddToSet: bson.M{
			"a": bson.M{
				operator.Each: []int{1, 2},
			},
		},
	})
	actual, _ := bson.Marshal(AddEachToSet[int]("a", []int{1, 2}))
	assert.Equal(t, expected, actual)
}

func TestPull(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Pull: bson.M{
			"a": 1,
			"b": 2,
		},
	})
	actual, _ := bson.Marshal(Pull("a", 1).Pull("b", 2))
	assert.Equal(t, expected, actual)
}

func TestPullConditioned(t *testing.T) {
	expected, _ := bson.Marshal(bson.M{
		operator.Pull: bson.M{
			"a": bson.M{
				operator.Gt: 10,
			},
		},
	})
	actual, _ := bson.Marshal(PullConditioned(filters.Gt("a", 10)))
	assert.Equal(t, expected, actual)
}
