package database

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Builder struct {
	cli  *mongo.Client
	name string
	opts *options.DatabaseOptions
}

func NewBuilder(cli *mongo.Client, name string) *Builder {
	return &Builder{
		cli:  cli,
		name: name,
		opts: options.Database(),
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

// BSONOptions configures optional BSON marshaling and unmarshaling behavior.
func (b *Builder) BSONOptions(opts *options.BSONOptions) *Builder {
	b.opts.SetBSONOptions(opts)
	return b
}

func (b *Builder) Build() *Database {
	primary := b.cli.Database(b.name, options.Database().SetReadPreference(readpref.Primary()), b.opts)
	defaultReadpref := b.cli.Database(b.name, b.opts)
	return New(primary, defaultReadpref)
}
