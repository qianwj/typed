package database

import (
	"context"
	"github.com/qianwj/typed/mongo/collection"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	db      *mongo.Database
	primary *mongo.Database
}

func (d *Database) Raw() *mongo.Database {
	return d.db
}

func (d *Database) Aggregate(pipe *aggregates.Pipeline) *AggregateExecutor {
	return newAggregateExecutor(d.primary, d.db, pipe)
}

func (d *Database) Collection(name string) *collection.Builder {
	return collection.NewBuilder(d.db, name)
}

func (d *Database) ListCollections(filter *filters.Filter) *ListCollectionsExecutor {
	return newListCollectionsExecutor(d.db, filter)
}

func (d *Database) Drop(ctx context.Context) error {
	return d.primary.Drop(ctx)
}
