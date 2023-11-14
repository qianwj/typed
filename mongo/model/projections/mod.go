package projections

import (
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Include = true
	Exclude = false
)

type Options struct {
	fields bson.D
}

func New() *Options {
	return &Options{
		fields: primitive.D{},
	}
}

func Includes(fields ...string) *Options {
	return New().Includes(fields...)
}

func ExcludeId() *Options {
	return New().ExcludeId()
}

func (p *Options) Includes(fields ...string) *Options {
	p.fields = append(p.fields, util.OrderedMap(fields, func(f string) bson.E {
		return bson.E{Key: f, Value: Include}
	})...)
	return p
}

func (p *Options) AddField(field string, expression any) *Options {
	p.fields = append(p.fields, primitive.E{Key: field, Value: expression})
	return p
}

func (p *Options) ExcludeId() *Options {
	p.fields = append(p.fields, bson.E{Key: "_id", Value: Exclude})
	return p
}

func (p *Options) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.fields)
}

func (p *Options) UnmarshalBSON(bytes []byte) error {
	var fields bson.D
	if err := bson.Unmarshal(bytes, fields); err != nil {
		return err
	}
	p.fields = fields
	return nil
}
