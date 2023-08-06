package executor

import (
	"github.com/qianwj/typed/mongo/model/aggregate/pipe"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (d *Database) Aggregate(pipe *pipe.Pipeline) *DatabaseAggregateExecutor {
	return &DatabaseAggregateExecutor{
		db:   d,
		pipe: pipe,
		opts: options.Aggregate(),
	}
}
