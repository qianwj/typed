package database

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListCollectionsExecutor struct {
	db     *mongo.Database
	filter *filters.Filter
	opts   *options.ListCollectionsOptions
}

func newListCollectionsExecutor(db *mongo.Database, filter *filters.Filter) *ListCollectionsExecutor {
	return &ListCollectionsExecutor{
		db:     db,
		filter: filter,
		opts:   options.ListCollections(),
	}
}

// NameOnly sets the value for the NameOnly field.
func (l *ListCollectionsExecutor) NameOnly() *ListCollectionsExecutor {
	l.opts.SetNameOnly(true)
	return l
}

// BatchSize sets the value for the BatchSize field.
func (l *ListCollectionsExecutor) BatchSize(size int32) *ListCollectionsExecutor {
	l.opts.SetBatchSize(size)
	return l
}

// AuthorizedCollections sets the value for the AuthorizedCollections field. This option is only valid for MongoDB server versions >= 4.0. Server
// versions < 4.0 ignore this option.
func (l *ListCollectionsExecutor) AuthorizedCollections() *ListCollectionsExecutor {
	l.opts.SetAuthorizedCollections(true)
	return l
}

func (l *ListCollectionsExecutor) Name(ctx context.Context) ([]string, error) {
	return l.db.ListCollectionNames(ctx, l.filter, l.opts)
}

func (l *ListCollectionsExecutor) Specification(ctx context.Context) ([]*mongo.CollectionSpecification, error) {
	return l.db.ListCollectionSpecifications(ctx, l.filter, l.opts)
}

func (l *ListCollectionsExecutor) Collect(ctx context.Context, result any) error {
	cursor, err := l.db.ListCollections(ctx, l.filter, l.opts)
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}
