package database

import (
	"github.com/qianwj/typed/mongo/builder/database"
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

func (d *Database) Aggregate(pipe aggregate.Pipeline) *database.AggregateExecutor {
	return database.NewAggregateExecutor(d.primary, d.defaultReadpref, pipe)
}
