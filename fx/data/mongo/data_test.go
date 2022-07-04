package mongo

import (
	"context"
	"fmt"
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
	app := fx.New(
		data.Provide(),
		data2.Provide(),
		fx.Provide(fx.Annotate(newRepo, fx.ParamTags(`name:"mongo"`, `name:"test_mongo"`))),
		fx.Provide(newMQ),
		fx.Provide(newServer),
		fx.Invoke(func(r *testServer, m *testMQ) {
			t.Logf("datasource is same?%v", r.dataSources == m.dataSources)
			t.Logf("server inject %d mongo client\r\n", len(r.dataSources.DataSources))
			t.Logf("mq inject %d mongo client\r\n", len(m.dataSources.DataSources))
			t.Log("invoke function\r\n")
		}))
	app.Run()
}

type testRepo struct {
	data *mongo.Client `name:"mongo"`
	test *mongo.Client `name:"test_mongo"`
}

type testServer struct {
	dataSources *Controller
	repo        *testRepo
}

type testMQ struct {
	dataSources *Controller
	repo        *testRepo
}

type Controller struct {
	fx.In
	DataSources []tfx.DataSource `group:"data_sources"`
}

func newRepo(c *mongo.Client, test *mongo.Client) *testRepo {
	return &testRepo{data: c, test: test}
}

func newMQ(r *testRepo, c Controller) *testMQ {
	fmt.Println("mq")
	fmt.Println(c)
	return &testMQ{
		dataSources: &c,
		repo:        r,
	}
}

func newServer(r *testRepo, c Controller) *testServer {
	fmt.Println("server")
	fmt.Println(c)
	return &testServer{
		dataSources: &c,
		repo:        r,
	}
}

func (r *testRepo) Get(ctx context.Context) ([]string, error) {
	return r.data.ListDatabaseNames(ctx, bson.M{})
}
