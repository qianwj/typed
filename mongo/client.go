package mongo

import (
	"context"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Client struct {
	internal            *raw.Client
	clusterMode         bool
	defaultDatabaseName string
}

func newClient(ctx context.Context, uri string, opts ...*options.ClientOptions) (*Client, error) {
	connStr, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}
	opts = append(opts, options.Client().ApplyURI(uri))
	internal, err := raw.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal:            internal,
		clusterMode:         connStr.ReplicaSet != "",
		defaultDatabaseName: connStr.Database,
	}, nil
}

func (c *Client) DefaultDatabase(opts ...*options.DatabaseOptions) *Database {
	db := Database{
		primary: c.internal.Database(c.defaultDatabaseName),
	}
	if c.clusterMode {
		db.secondary = c.internal.Database(
			c.defaultDatabaseName,
			append(opts, options.Database().SetReadPreference(readpref.Secondary()))...,
		)
	} else {
		db.secondary = db.primary
	}
	return &db
}
