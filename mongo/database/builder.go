// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	db := b.cli.Database(b.name, b.opts)
	return &Database{db: db, primary: primary}
}
