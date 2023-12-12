package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestDistinctExecutor_Execute(t *testing.T) {
	ctx := context.TODO()
	fruits := []*fruit{
		{Id: primitive.NewObjectID(), Type: "apple", Name: "red_fuji", Rest: 100},
		{Id: primitive.NewObjectID(), Type: "apple", Name: "gala", Rest: 150},
		{Id: primitive.NewObjectID(), Type: "apple", Name: "jazz", Rest: 0},
	}
	fruitColl := testDB.Collection("fruits")
	_, _ = fruitColl.InsertMany(ctx, util.ToAny(fruits))
	exec := newDistinctExecutor(fruitColl, fruitColl, "type")
	res, err := exec.Execute(ctx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v", res)
}
