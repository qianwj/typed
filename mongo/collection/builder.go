package collection

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Builder struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func NewBuilder(db *mongo.Database, name string) *Builder {
	return &Builder{
		db:   db,
		name: name,
		opts: options.Collection(),
	}
}

// ReadConcern sets the value for the ReadConcern field.
func (b *Builder) ReadConcern(rc *readconcern.ReadConcern) *Builder {
	b.opts.SetReadConcern(rc)
	return b
}

// WriteConcern sets the value for the WriteConcern field.
func (b *Builder) WriteConcern(wc *writeconcern.WriteConcern) *Builder {
	b.opts.SetWriteConcern(wc)
	return b
}

// ReadPreference sets the value for the ReadPreference field.
func (b *Builder) ReadPreference(rp *readpref.ReadPref) *Builder {
	b.opts.SetReadPreference(rp)
	return b
}

// Registry sets the value for the Registry field.
func (b *Builder) Registry(r *bsoncodec.Registry) *Builder {
	b.opts.SetRegistry(r)
	return b
}

//
//func (b *Builder) Build() *collection.Collection {
//	primary := b.db.Collection(b.name, options.Collection().SetReadPreference(readpref.Primary()), b.opts)
//	defaultNode := b.db.Collection(b.name, b.opts)
//	return collection.NewCollection[D, I](primary, defaultNode)
//}
