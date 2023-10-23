package database

import (
	"github.com/qianwj/typed/mongo/collection"
	"github.com/qianwj/typed/mongo/executor"
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

func (d *Database) Aggregate(pipe mongo.Pipeline) *executor.DatabaseAggregateExecutor {
	return executor.NewDatabaseAggregateExecutor(d.primary, d.defaultReadpref, pipe)
}

func (d *Database) Collection(name string) *collection.Builder {
	return collection.NewBuilder(d.defaultReadpref, name)
}
