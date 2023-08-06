package builder

import (
	"github.com/qianwj/typed/mongo/executor"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type DatabaseBuilder struct {
	cli  *mongo.Client
	name string
	opts *options.DatabaseOptions
}

func NewDatabase(cli *mongo.Client, name string) *DatabaseBuilder {
	return &DatabaseBuilder{
		cli:  cli,
		name: name,
		opts: options.Database(),
	}
}

// ReadConcern sets the value for the ReadConcern field.
func (b *DatabaseBuilder) ReadConcern(rc *readconcern.ReadConcern) *DatabaseBuilder {
	b.opts.SetReadConcern(rc)
	return b
}

// WriteConcern sets the value for the WriteConcern field.
func (b *DatabaseBuilder) WriteConcern(wc *writeconcern.WriteConcern) *DatabaseBuilder {
	b.opts.SetWriteConcern(wc)
	return b
}

// ReadPreference sets the value for the ReadPreference field.
func (b *DatabaseBuilder) ReadPreference(rp *readpref.ReadPref) *DatabaseBuilder {
	b.opts.SetReadPreference(rp)
	return b
}

// Registry sets the value for the Registry field.
func (b *DatabaseBuilder) Registry(r *bsoncodec.Registry) *DatabaseBuilder {
	b.opts.SetRegistry(r)
	return b
}

// BSONOptions configures optional BSON marshaling and unmarshaling behavior.
func (b *DatabaseBuilder) BSONOptions(opts *options.BSONOptions) *DatabaseBuilder {
	b.opts.SetBSONOptions(opts)
	return b
}

func (b *DatabaseBuilder) build() *executor.Database {
	primary := b.cli.Database(b.name, options.Database().SetReadPreference(readpref.Primary()), b.opts)
	defaultReadpref := b.cli.Database(b.name, b.opts)
	return executor.NewDatabase(primary, defaultReadpref)
}
