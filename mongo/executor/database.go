package executor

import (
	"github.com/qianwj/typed/mongo/model/aggregate/pipe"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	primary   *raw.Database
	secondary *raw.Database
}

func NewDatabase(primary, defaultReadpref *raw.Database) *Database {
	return &Database{
		primary:   primary,
		secondary: defaultReadpref,
	}
}

func (d *Database) Aggregate(pipe *pipe.Pipeline) *DatabaseAggregateExecutor {
	return &DatabaseAggregateExecutor{
		db:   d,
		pipe: pipe,
		opts: options.Aggregate(),
	}
}
