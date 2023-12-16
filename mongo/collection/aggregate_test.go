package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/model/aggregates/operators"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type fruit struct {
	model.Doc[primitive.ObjectID]
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Type string             `bson:"type"`
	Name string             `bson:"name"`
	Rest int                `bson:"rest"`
}

func (f *fruit) GetID() primitive.ObjectID {
	return f.Id
}

type order struct {
	FruitId string `bson:"fruitId"`
	User    string `bson:"user"`
	Count   int    `bson:"count"`
}

func TestAggregateExecutor_Collect(t *testing.T) {
	ctx := context.TODO()
	fruits := []*fruit{
		{Id: primitive.NewObjectID(), Type: "apple", Name: "red_fuji", Rest: 100},
		{Id: primitive.NewObjectID(), Type: "apple", Name: "gala", Rest: 150},
		{Id: primitive.NewObjectID(), Type: "apple", Name: "jazz", Rest: 0},
	}
	fruitColl := testDB.Collection("fruits")
	_, _ = fruitColl.InsertMany(ctx, util.ToAny(fruits))

	t.Run("testCountDocuments", func(t *testing.T) {
		pipe := aggregates.Group(1, operators.Sum("n", 1))
		t.Logf("pipe: %+v", pipe.Stages())
		var result []primitive.M
		err := newAggregateExecutor[*fruit, primitive.ObjectID](fruitColl, fruitColl, pipe).Collect(ctx, &result)
		if err != nil {
			t.Errorf("agg error: %+v", err)
			t.FailNow()
		}
		t.Logf("result: %+v", result)
		assert.Equal(t, int32(3), result[0]["n"])
	})

	t.Run("testCountWithMatch", func(t *testing.T) {
		pipe := aggregates.Match(filters.Eq("name", "gala")).Group(1, operators.Sum("n", 1))
		t.Logf("pipe: %+v", pipe.Stages())
		var result []primitive.M
		err := newAggregateExecutor[*fruit, primitive.ObjectID](fruitColl, fruitColl, pipe).Collect(ctx, &result)
		if err != nil {
			t.Errorf("agg error: %+v", err)
			t.FailNow()
		}
		t.Logf("result: %+v", result)
		assert.Equal(t, int32(1), result[0]["n"])
	})
}
