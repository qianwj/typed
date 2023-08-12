package collection

import (
	"fmt"
	"github.com/qianwj/typed/mongo/executor/client"
	"github.com/qianwj/typed/mongo/executor/collection"
	"github.com/qianwj/typed/mongo/executor/database"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strings"
)

type TypedBuilder[D model.Document[I], I model.DocumentId] struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func FromNamespace[D model.Document[I], I model.DocumentId](cli *client.Client, ns string) (*TypedBuilder[D, I], error) {
	pair := strings.Split(ns, ".")
	if len(pair) != 2 {
		return nil, fmt.Errorf("invalid ns: %s", ns)
	}
	db := cli.Database(pair[0]).Build().Raw()
	return NewTypedCollection[D, I](db, pair[1]), nil
}

func FromDatabase[D model.Document[I], I model.DocumentId](db *database.Database, name string) *TypedBuilder[D, I] {
	return NewTypedCollection[D, I](db.Raw(), name)
}

func NewTypedCollection[D model.Document[I], I model.DocumentId](db *mongo.Database, name string) *TypedBuilder[D, I] {
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

func (b *TypedBuilder[D, I]) Build() *collection.TypedCollection[D, I] {
	primary := b.db.Collection(b.name, options.Collection().SetReadPreference(readpref.Primary()), b.opts)
	defaultNode := b.db.Collection(b.name, b.opts)
	return collection.NewTypedCollection[D, I](primary, defaultNode)
}
