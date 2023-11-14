package mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/client"
	"github.com/qianwj/typed/mongo/collection"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

type testDoc struct {
	model.Doc[primitive.ObjectID]
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (t *testDoc) GetId() primitive.ObjectID {
	return t.Id
}

func TestClientBuilder_ApplyUri(t *testing.T) {
	cli, err := client.NewBuilder().ApplyUri("mongodb://localhost:27017").Build(context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Error(err)
		t.FailNow()
	}
	db := cli.DefaultDatabase().Build().Raw()
	coll := collection.NewTypedBuilder[testDoc, primitive.ObjectID](db, "test").Build()
	iter, err := coll.Find(filters.Eq("name", "aaa")).Cursor(context.TODO())
	if err != nil {
		t.Errorf("curosr error: %+v", err)
		t.FailNow()
	}
	for iter.HasNext(context.TODO()) {
		doc, err := iter.Next()
		if err != nil {
			t.Errorf("doc error: %+v", err)
			t.FailNow()
		}
		t.Logf("current: %+v", doc)
	}
}
