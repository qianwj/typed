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
	internal *mongo.Client
}

func NewData(opts ...*options.ClientOptions) (tfx.DataAccess, error) {
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &mongoClient{
		internal: client,
	}, nil
}

func Apply(uri string) (tfx.DataAccess, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("uri: %s", uri))
	}
	return &mongoClient{
		internal: client,
	}, nil
}

func (m *mongoClient) Provide(name ...string) fx.Option {
	if len(name) == 0 {
		return fx.Provide(fx.Annotate(func() *mongo.Client {
			return m.internal
		}, fx.ResultTags(`name:"mongo"`)))
	} else {
		return fx.Provide(fx.Annotate(func() *mongo.Client {
			return m.internal
		}, fx.ResultTags(fmt.Sprintf(`name:"%s_mongo"`, name[0]))))
	}
}

func (m *mongoClient) Connect(ctx context.Context) error {
	return m.internal.Connect(ctx)
}

func (m *mongoClient) Close(ctx context.Context) error {
	return m.internal.Disconnect(ctx)
}
