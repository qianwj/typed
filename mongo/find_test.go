package typed_mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type testDoc struct {
	Name  string
	Value int
}

func TestFindOne(t *testing.T) {
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}
	cli.Connect(context.TODO())
	db := cli.Database("test_typed")
	coll := db.Collection("test_doc")
	coll.InsertOne(context.TODO(), testDoc{
		Name:  "abc",
		Value: 1,
	})
	res, err := FindOne[testDoc](context.TODO(), coll, model.NewFilter().Eq("name", "abc"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v", res)
}

func TestFind(t *testing.T) {
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}
	cli.Connect(context.TODO())
	db := cli.Database("test_typed")
	coll := db.Collection("test_doc")
	coll.InsertOne(context.TODO(), testDoc{
		Name:  "abc",
		Value: 1,
	})
	res, err := Find[testDoc](context.TODO(), coll, model.NewFilter().Eq("name", "abc"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, item := range res {
		t.Logf("%+v", item.Name)
	}
}
