package mongo

import (
	"context"
	tfx "github.com/qianwj/typed/fx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"testing"
)

func TestApply(t *testing.T) {
	data, err := Apply("mongodb://localhost:27017")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := data.Connect(context.TODO()); err != nil {
		t.Error(err)
		t.FailNow()
	}
	var data2 tfx.DataAccess
	data2, err = Apply("mongodb://localhost:27017")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := data2.Connect(context.TODO()); err != nil {
		t.Error(err)
		t.FailNow()
	}
	fx.New(data.Provide(), data2.Provide("test"), fx.Provide(fx.Annotate(newRepo, fx.ParamTags(`name:"mongo"`, `name:"test_mongo"`))), fx.Invoke(func(r *testRepo) {
		data, err := r.Get(context.TODO())
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log(data)
	}))
}

type testRepo struct {
	data *mongo.Client `name:"mongo"`
	test *mongo.Client `name:"test_mongo"`
}

func newRepo(c *mongo.Client, test *mongo.Client) *testRepo {
	return &testRepo{data: c, test: test}
}

func (r *testRepo) Get(ctx context.Context) ([]string, error) {
	return r.data.ListDatabaseNames(ctx, bson.M{})
}
