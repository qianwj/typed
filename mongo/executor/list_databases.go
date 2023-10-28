package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListDatabasesExecutor struct {
	cli    *mongo.Client
	filter *filters.Filter
	opts   *options.ListDatabasesOptions
}

func NewListDatabasesExecutor(cli *mongo.Client, filter *filters.Filter) *ListDatabasesExecutor {
	return &ListDatabasesExecutor{
		cli:    cli,
		filter: filter,
		opts:   options.ListDatabases(),
	}
}

// NameOnly sets the value for the NameOnly field.
func (l *ListDatabasesExecutor) NameOnly() *ListDatabasesExecutor {
	l.opts.SetNameOnly(true)
	return l
}

// AuthorizedDatabases sets the value for the AuthorizedDatabases field.
func (l *ListDatabasesExecutor) AuthorizedDatabases() *ListDatabasesExecutor {
	l.opts.SetAuthorizedDatabases(true)
	return l
}

func (l *ListDatabasesExecutor) Name(ctx context.Context) ([]string, error) {
	return l.cli.ListDatabaseNames(ctx, l.filter, l.opts)
}

func (l *ListDatabasesExecutor) Collect(ctx context.Context) (*mongo.ListDatabasesResult, error) {
	result, err := l.cli.ListDatabases(ctx, l.filter, l.opts)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
