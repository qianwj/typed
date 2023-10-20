package client

import (
	"context"
	"github.com/qianwj/typed/mongo/database"
	"github.com/qianwj/typed/mongo/transaction"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	internal            *mongo.Client
	defaultDatabaseName string
}

func New(ctx context.Context, defaultDB string, opts *options.ClientOptions) (*Client, error) {
	internal, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal:            internal,
		defaultDatabaseName: defaultDB,
	}, nil
}

func (c *Client) DefaultDatabase() *database.Builder {
	return database.NewBuilder(c.internal, c.defaultDatabaseName)
}

func (c *Client) Database(name string) *database.Builder {
	return database.NewBuilder(c.internal, name)
}

func (c *Client) Transaction() *transaction.TxSessionBuilder {
	return transaction.NewTxSessionBuilder(c.internal)
}

func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.internal.Ping(ctx, rp)
}
