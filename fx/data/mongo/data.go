package mongo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	tfx "github.com/qianwj/typed/fx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type mongoClient struct {
	name     string
	internal *mongo.Client
	uri      string
}

func NewData(opts ...*options.ClientOptions) (tfx.DataSource, error) {
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &mongoClient{
		internal: client,
	}, nil
}

func Apply(uri string) (tfx.DataSource, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("uri: %s", uri))
	}
	return &mongoClient{
		internal: client,
		uri:      uri,
	}, nil
}

func (m *mongoClient) Name(name string) tfx.DataSource {
	m.name = name
	return m
}

func (m *mongoClient) Provide() fx.Option {
	data := fx.Provide(fx.Annotate(func() tfx.DataSource {
		return m
	}, fx.ResultTags(`group:"data_sources"`)))
	tag := `name:"mongo"`
	if len(m.name) > 0 {
		tag = fmt.Sprintf(`name:"%s_mongo"`, m.name)
	}
	return fx.Options(data, fx.Provide(fx.Annotate(m.client, fx.ResultTags(tag))))
}

func (m *mongoClient) client() *mongo.Client {
	return m.internal
}

func (m *mongoClient) Connect(ctx context.Context) error {
	err := m.internal.Connect(ctx)
	if err != nil {
		return errors.Wrapf(err, "mongo: connecting [%s] error", m.uri)
	}
	return nil
}

func (m *mongoClient) Close(ctx context.Context) error {
	return m.internal.Disconnect(ctx)
}
