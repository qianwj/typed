package collection

import (
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type TypedBuilder[D model.Doc[I], I model.ID] struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func NewTypedBuilder[D model.Doc[I], I model.ID](db *mongo.Database, name string) *TypedBuilder[D, I] {
	return &TypedBuilder[D, I]{
		db:   db,
		name: name,
		opts: options.Collection(),
	}
}

// ReadConcern sets the value for the ReadConcern field.
func (b *TypedBuilder[D, I]) ReadConcern(rc *readconcern.ReadConcern) *TypedBuilder[D, I] {
	b.opts.SetReadConcern(rc)
	return b
}

// WriteConcern sets the value for the WriteConcern field.
func (b *TypedBuilder[D, I]) WriteConcern(wc *writeconcern.WriteConcern) *TypedBuilder[D, I] {
	b.opts.SetWriteConcern(wc)
	return b
}

// ReadPreference sets the value for the ReadPreference field.
func (b *TypedBuilder[D, I]) ReadPreference(rp *readpref.ReadPref) *TypedBuilder[D, I] {
	b.opts.SetReadPreference(rp)
	return b
}

// Registry sets the value for the Registry field.
func (b *TypedBuilder[D, I]) Registry(r *bsoncodec.Registry) *TypedBuilder[D, I] {
	b.opts.SetRegistry(r)
	return b
}

func (b *TypedBuilder[D, I]) Build() *TypedCollection[D, I] {
	primary := b.db.Collection(b.name, options.Collection().SetReadPreference(readpref.Primary()), b.opts)
	defaultNode := b.db.Collection(b.name, b.opts)
	return NewTypedCollection[D, I](primary, defaultNode)
}
