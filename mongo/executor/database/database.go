package database

import (
	"github.com/qianwj/typed/mongo/builder/collection"
	"github.com/qianwj/typed/mongo/executor"
	"github.com/qianwj/typed/mongo/model/aggregate"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	primary         *mongo.Database
	defaultReadpref *mongo.Database
}

func NewDatabase(primary, defaultReadpref *mongo.Database) *Database {
	return &Database{
		primary:         primary,
		defaultReadpref: defaultReadpref,
	}
}

func (d *Database) Raw() *mongo.Database {
	return d.defaultReadpref
}

func (d *Database) Aggregate(pipe aggregate.Pipeline) *executor.DatabaseAggregateExecutor {
	return executor.NewDatabaseAggregateExecutor(d.primary, d.defaultReadpref, pipe)
}

func (d *Database) Collection(name string) *collection.Builder {
	return collection.NewCollection(d.defaultReadpref, name)
}
