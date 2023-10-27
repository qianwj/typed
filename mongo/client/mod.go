package client

import (
	"context"
	"github.com/qianwj/typed/mongo/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	internal            *mongo.Client
	defaultDatabaseName string
}

func newClient(ctx context.Context, pingReadpref *readpref.ReadPref, defaultDB string, opts *options.ClientOptions) (*Client, error) {
	internal, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if pingReadpref != nil {
		err := internal.Ping(ctx, pingReadpref)
		if err != nil {
			return nil, err
		}
	}
	return &Client{
		internal:            internal,
		defaultDatabaseName: defaultDB,
	}, nil
}

// DefaultDatabase the database which specified in uri, if not present, the default database is admin.
func (c *Client) DefaultDatabase() *database.Builder {
	return database.NewBuilder(c.internal, c.defaultDatabaseName)
}

func (c *Client) Database(name string) *database.Builder {
	return database.NewBuilder(c.internal, name)
}

func (c *Client) Transaction() *TxSessionBuilder {
	return NewTxSessionBuilder(c.internal)
}

func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.internal.Ping(ctx, rp)
}

func (c *Client) Disconnect(ctx context.Context) error {
	if c != nil && c.internal != nil {
		return c.internal.Disconnect(ctx)
	}
	return nil
}
