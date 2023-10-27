package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/model/aggregates/lookup"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type fruit struct {
	bson.Doc[primitive.ObjectID]
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

	orders := []*order{
		{FruitId: fruits[0].Id.Hex(), User: "elvis", Count: 10},
		{FruitId: fruits[1].Id.Hex(), User: "qianwj", Count: 12},
	}
	orderColl := testDB.Collection("orders")
	_, _ = orderColl.InsertMany(ctx, util.ToAny(orders))

	pipe := aggregates.Lookup(
		lookup.New("fruits", "fruit_orders").
			Join("fruitId", "fruitId").
			SetPipeline(aggregates.Set(bson.NewMap().Put("fruitId", primitive.M{
				operator.ToObjectID: "$_id",
			})).Stages()),
	)
	//pipe := aggregates.Match(filters.Eq("name", "gala"))
	var result []primitive.M
	err := NewAggregateExecutor[*fruit, primitive.ObjectID](fruitColl, fruitColl, pipe).Collect(ctx, &result)
	if err != nil {
		t.Errorf("agg error: %+v", err)
		t.FailNow()
	}
	t.Logf("result: %+v", result)
}
