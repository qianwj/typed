package updates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestSet(t *testing.T) {
	expected := bson.M{
		operator.Set: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := Set("a", 1).Set("b", 2).Document()
	assert.Equal(t, expected, actual)
}

func TestSetAll(t *testing.T) {
	expected := bson.M{
		operator.Set: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := SetAll(bson.M{
		"a": 1,
		"b": 2,
	}).Document()
	assert.Equal(t, expected, actual)
}

func TestSetOnInsert(t *testing.T) {
	expected := bson.M{
		operator.SetOnInsert: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := SetOnInsert("a", 1).SetOnInsert("b", 2).Document()
	assert.Equal(t, expected, actual)
}

func TestSetOnInsertAll(t *testing.T) {
	expected := bson.M{
		operator.SetOnInsert: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := SetOnInsertAll(bson.M{
		"a": 1,
		"b": 2,
	}).Document()
	assert.Equal(t, expected, actual)
}

func TestUnset(t *testing.T) {
	expected := bson.M{
		operator.Unset: bson.M{
			"b": "",
			"c": "",
		},
	}
	actual := Unset("b", "c").Document()
	assert.Equal(t, expected, actual)
}

func TestInc(t *testing.T) {
	expected := bson.M{
		operator.Inc: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := Inc("a", 1).Inc("b", 2).Document()
	assert.Equal(t, expected, actual)
}

func TestAddToSet(t *testing.T) {
	expected := bson.M{
		operator.AddToSet: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := AddToSet("a", 1).AddToSet("b", 2).Document()
	assert.Equal(t, expected, actual)
}

func TestAddEachToSet(t *testing.T) {
	expected := bson.M{
		operator.AddToSet: bson.M{
			"a": bson.M{
				operator.Each: []any{1, 2},
			},
		},
	}
	actual := AddEachToSet[int]("a", []int{1, 2}).Document()
	assert.Equal(t, expected, actual)
}

func TestPull(t *testing.T) {
	expected := bson.M{
		operator.Pull: bson.M{
			"a": 1,
			"b": 2,
		},
	}
	actual := Pull("a", 1).Pull("b", 2).Document()
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
