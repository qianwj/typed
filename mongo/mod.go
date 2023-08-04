package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FromUri(ctx context.Context, uri string, opts ...*options.ClientOptions) (*Client, error) {
	return newClient(ctx, uri, opts...)
}
