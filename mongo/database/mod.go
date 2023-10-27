package database

import (
	"context"
	"github.com/qianwj/typed/mongo/collection"
	"github.com/qianwj/typed/mongo/executor"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	primary         *mongo.Database
	defaultReadpref *mongo.Database
}

func New(primary, defaultReadpref *mongo.Database) *Database {
	return &Database{
		primary:         primary,
		defaultReadpref: defaultReadpref,
	}
}

func (d *Database) Raw() *mongo.Database {
	return d.defaultReadpref
}

func (d *Database) Aggregate(pipe *aggregates.Pipeline) *executor.DatabaseAggregateExecutor {
	return executor.NewDatabaseAggregateExecutor(d.primary, d.defaultReadpref, pipe)
}

func (d *Database) Collection(name string) *collection.Builder {
	return collection.NewBuilder(d.defaultReadpref, name)
}

func (d *Database) ListCollections(filter *filters.Filter) *executor.ListCollectionsExecutor {
	return executor.NewListCollectionsExecutor(d.defaultReadpref, filter)
}

func (d *Database) Drop(ctx context.Context) error {
	return d.primary.Drop(ctx)
}
