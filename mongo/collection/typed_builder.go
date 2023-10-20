package collection

import (
	"fmt"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/client"
	"github.com/qianwj/typed/mongo/database"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strings"
)

type TypedBuilder[D bson.Doc[I], I bson.ID] struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func FromNamespace[D bson.Doc[I], I bson.ID](cli *client.Client, ns string) (*TypedBuilder[D, I], error) {
	pair := strings.Split(ns, ".")
	if len(pair) != 2 {
		return nil, fmt.Errorf("invalid ns: %s", ns)
	}
	db := cli.Database(pair[0]).Build().Raw()
	return NewTypedBuilder[D, I](db, pair[1]), nil
}

func FromDatabase[D bson.Doc[I], I bson.ID](db *database.Database, name string) *TypedBuilder[D, I] {
	return NewTypedBuilder[D, I](db.Raw(), name)
}

func NewTypedBuilder[D bson.Doc[I], I bson.ID](db *mongo.Database, name string) *TypedBuilder[D, I] {
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
