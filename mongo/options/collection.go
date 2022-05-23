package options

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type CollectionOptions struct {
	internal *options.CollectionOptions
}

// Collection creates a new CollectionOptions instance.
func Collection() *CollectionOptions {
	return &CollectionOptions{
		internal: options.Collection(),
	}
}

// SetReadConcern sets the value for the ReadConcern field.
func (c *CollectionOptions) SetReadConcern(rc *readconcern.ReadConcern) *CollectionOptions {
	c.internal.ReadConcern = rc
	return c
}

// SetWriteConcern sets the value for the WriteConcern field.
func (c *CollectionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *CollectionOptions {
	c.internal.WriteConcern = wc
	return c
}

// SetReadPreference sets the value for the ReadPreference field.
func (c *CollectionOptions) SetReadPreference(rp *readpref.ReadPref) *CollectionOptions {
	c.internal.ReadPreference = rp
	return c
}

// SetRegistry sets the value for the Registry field.
func (c *CollectionOptions) SetRegistry(r *bsoncodec.Registry) *CollectionOptions {
	c.internal.Registry = r
	return c
}

// MergeCollectionOptions combines the given CollectionOptions instances into a single *CollectionOptions in a
// last-one-wins fashion.
func MergeCollectionOptions(opts ...*CollectionOptions) *options.CollectionOptions {
	c := options.Collection()

	for _, o := range opts {
		opt := o.internal
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			c.ReadConcern = opt.ReadConcern
		}
		if opt.WriteConcern != nil {
			c.WriteConcern = opt.WriteConcern
		}
		if opt.ReadPreference != nil {
			c.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			c.Registry = opt.Registry
		}
	}

	return c
}
