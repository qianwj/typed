package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/builder"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	internal            *mongo.Client
	defaultDatabaseName string
}

func NewClient(ctx context.Context, defaultDB string, opts *options.ClientOptions) (*Client, error) {
	internal, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal:            internal,
		defaultDatabaseName: defaultDB,
	}, nil
}

func (c *Client) DefaultDatabase() *builder.DatabaseBuilder {
	return builder.NewDatabase(c.internal, c.defaultDatabaseName)
}

func (c *Client) Database(name string) *builder.DatabaseBuilder {
	return builder.NewDatabase(c.internal, name)
}
