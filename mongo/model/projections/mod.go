package projections

import (
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Include = 1
	Exclude = -1
)

type ProjectionOptions struct {
	fields bson.D
}

func New() *ProjectionOptions {
	return &ProjectionOptions{
		fields: bson.D{},
	}
}

func Includes(fields ...string) *ProjectionOptions {
	return New().Includes(fields...)
}

func ExcludeId() *ProjectionOptions {
	return New().ExcludeId()
}

func (p *ProjectionOptions) Includes(fields ...string) *ProjectionOptions {
	p.fields = append(p.fields, util.Map(fields, func(f string) bson.E {
		return bson.E{Key: f, Value: Include}
	})...)
	return p
}

func (p *ProjectionOptions) ExcludeId() *ProjectionOptions {
	p.fields = append(p.fields, bson.E{Key: "_id", Value: Exclude})
	return p
}

func (p *ProjectionOptions) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.fields)
}

func (p *ProjectionOptions) UnmarshalBSON(bytes []byte) error {
	var fields bson.D
	if err := bson.Unmarshal(bytes, fields); err != nil {
		return err
	}
	p.fields = fields
	return nil
}
