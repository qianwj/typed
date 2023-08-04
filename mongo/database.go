package mongo

import (
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	primary   *raw.Database
	secondary *raw.Database
}

func (d *Database) Primary(coll string, opts ...*options.CollectionOptions) *raw.Collection {
	return d.primary.Collection(coll, opts...)
}

func (d *Database) Secondary(coll string, opts ...*options.CollectionOptions) *raw.Collection {
	return d.secondary.Collection(coll, opts...)
}
