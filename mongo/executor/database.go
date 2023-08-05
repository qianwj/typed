package executor

import (
	"github.com/qianwj/typed/mongo/model/aggregate"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	primary   *raw.Database
	secondary *raw.Database
}

func NewDatabase(cli *raw.Client, dbname string, cluster bool, opts ...*options.DatabaseOptions) *Database {
	db := Database{
		primary: cli.Database(dbname),
	}
	if cluster {
		db.secondary = cli.Database(
			dbname,
			append(opts, options.Database().SetReadPreference(readpref.Secondary()))...,
		)
	} else {
		db.secondary = db.primary
	}
	return &db
}

func (d *Database) Primary(coll string, opts ...*options.CollectionOptions) *raw.Collection {
	return d.primary.Collection(coll, opts...)
}

func (d *Database) Secondary(coll string, opts ...*options.CollectionOptions) *raw.Collection {
	return d.secondary.Collection(coll, opts...)
}

func (d *Database) Aggregate(pipe *aggregate.Pipeline) *AggregateExecutor {
	return &AggregateExecutor{
		db:   d,
		pipe: pipe,
		opts: options.Aggregate(),
	}
}
