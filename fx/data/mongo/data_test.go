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
	var data2 tfx.DataSource
	data2, err = Apply("mongodb://localhost:27017")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data2.Name("test")
	app := fx.New(data.Provide(), data2.Provide(), fx.Provide(fx.Annotate(newRepo, fx.ParamTags(`group:"data_sources"`, `name:"mongo"`, `name:"test_mongo"`))), fx.Invoke(func(r *testRepo) {
		//data, err := r.Get(context.TODO())
		//if err != nil {
		//	t.Error(err)
		//	t.FailNow()
		//}
		//t.Log(data)
		t.Logf("inject %d mongo client\r\n", len(r.dataSources))
		t.Log("invoke function\r\n")
	}))
	app.Run()
}

type testRepo struct {
	dataSources []tfx.DataSource `group:"data_sources"`
	data        *mongo.Client    `name:"mongo"`
	test        *mongo.Client    `name:"test_mongo"`
}

func newRepo(dataSources []tfx.DataSource, c *mongo.Client, test *mongo.Client) *testRepo {
	return &testRepo{dataSources: dataSources, data: c, test: test}
}

func (r *testRepo) Get(ctx context.Context) ([]string, error) {
	return r.data.ListDatabaseNames(ctx, bson.M{})
}
