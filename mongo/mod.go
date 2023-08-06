package mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/executor"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FromUri(ctx context.Context, uri string, opts ...*options.ClientOptions) (*executor.Client, error) {
	return newClient(ctx, uri, opts...)
}
