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

type Builder[D model.Document[I], I model.DocumentId] struct {
	db   *mongo.Database
	name string
	opts *options.CollectionOptions
}

func FromNamespace[D model.Document[I], I model.DocumentId](cli *client.Client, ns string) (*Builder[D, I], error) {
	pair := strings.Split(ns, ".")
	if len(pair) != 2 {
		return nil, fmt.Errorf("invalid ns: %s", ns)
	}
	db := cli.Database(pair[0]).Build().Raw()
	return NewCollection[D, I](db, pair[1]), nil
}

func FromDatabase[D model.Document[I], I model.DocumentId](db *database.Database, name string) *Builder[D, I] {
	return NewCollection[D, I](db.Raw(), name)
}

func NewCollection[D model.Document[I], I model.DocumentId](db *mongo.Database, name string) *Builder[D, I] {
	return &Builder[D, I]{
		db:   db,
		name: name,
		opts: options.Collection(),
	}
}

// ReadConcern sets the value for the ReadConcern field.
func (b *Builder[D, I]) ReadConcern(rc *readconcern.ReadConcern) *Builder[D, I] {
	b.opts.SetReadConcern(rc)
	return b
}

// WriteConcern sets the value for the WriteConcern field.
func (b *Builder[D, I]) WriteConcern(wc *writeconcern.WriteConcern) *Builder[D, I] {
	b.opts.SetWriteConcern(wc)
	return b
}

// ReadPreference sets the value for the ReadPreference field.
func (b *Builder[D, I]) ReadPreference(rp *readpref.ReadPref) *Builder[D, I] {
	b.opts.SetReadPreference(rp)
	return b
}

// Registry sets the value for the Registry field.
func (b *Builder[D, I]) Registry(r *bsoncodec.Registry) *Builder[D, I] {
	b.opts.SetRegistry(r)
	return b
}

func (b *Builder[D, I]) Build() *collection.Collection[D, I] {
	primary := b.db.Collection(b.name, options.Collection().SetReadPreference(readpref.Primary()), b.opts)
	defaultNode := b.db.Collection(b.name, b.opts)
	return collection.NewCollection[D, I](primary, defaultNode)
}
