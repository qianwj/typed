package typed_mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	typedColl := NewTypedCollection[testDoc](db, "test_doc")
	res, err := typedColl.FindOne(context.TODO(), model.NewFilter().Eq("name", "abc"))
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
	id := primitive.NewObjectID()
	coll.InsertOne(context.TODO(), bson.M{
		"_id":   id,
		"name":  "abc",
		"value": 1,
	})
	typedColl := NewTypedCollection[testDoc](db, "test_doc")
	res, err := typedColl.FindByDocIds(context.TODO(), []primitive.ObjectID{id})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, item := range res {
		t.Logf("%+v", item.Name)
	}
}
