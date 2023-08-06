package builder

import (
	"github.com/qianwj/typed/mongo/executor"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type CollectionBuilder[D model.Document[I], I model.DocumentId] struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func NewCollection[D model.Document[I], I model.DocumentId](db *mongo.Database, name string) *CollectionBuilder[D, I] {
	return &CollectionBuilder[D, I]{
		db:   db,
		name: name,
		opts: options.Collection(),
	}
}

// ReadConcern sets the value for the ReadConcern field.
func (b *CollectionBuilder[D, I]) ReadConcern(rc *readconcern.ReadConcern) *CollectionBuilder[D, I] {
	b.opts.SetReadConcern(rc)
	return b
}

// WriteConcern sets the value for the WriteConcern field.
func (b *CollectionBuilder[D, I]) WriteConcern(wc *writeconcern.WriteConcern) *CollectionBuilder[D, I] {
	b.opts.SetWriteConcern(wc)
	return b
}

// ReadPreference sets the value for the ReadPreference field.
func (b *CollectionBuilder[D, I]) ReadPreference(rp *readpref.ReadPref) *CollectionBuilder[D, I] {
	b.opts.SetReadPreference(rp)
	return b
}

// Registry sets the value for the Registry field.
func (b *CollectionBuilder[D, I]) Registry(r *bsoncodec.Registry) *CollectionBuilder[D, I] {
	b.opts.SetRegistry(r)
	return b
}

func (b *CollectionBuilder[D, I]) build() *executor.Collection[D, I] {
	primary := b.db.Collection(b.name, b.opts, options.Collection().SetReadPreference(readpref.Primary()))
	defaultNode := b.db.Collection(b.name, b.opts)
	return executor.NewCollection[D, I](primary, defaultNode)
}
