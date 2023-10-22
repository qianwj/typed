package sorts

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Order int

const (
	Desc Order = -1
	Asc  Order = 1
)

type Options struct {
	fields bson.D
}

func New() *Options {
	return &Options{fields: bson.D{}}
}

func Ascending(field string) *Options {
	return New().Ascending(field)
}

func Descending(field string) *Options {
	return New().Descending(field)
}

func Meta(field string) *Options {
	return New().Meta(field)
}

func (s *Options) Ascending(field string) *Options {
	s.fields = append(s.fields, bson.E{Key: field, Value: Asc})
	return s
}

func (s *Options) Descending(field string) *Options {
	s.fields = append(s.fields, bson.E{Key: field, Value: Desc})
	return s
}

func (s *Options) Meta(field string) *Options {
	s.fields = append(s.fields, bson.E{Key: field, Value: bson.M{operator.Meta: "textScore"}})
	return s
}

func (s *Options) MarshalBSON() ([]byte, error) {
	return bson.Marshal(s.fields)
}

func (s *Options) UnmarshalBSON(bytes []byte) error {
	var fields bson.D
	if err := bson.Unmarshal(bytes, fields); err != nil {
		return err
	}
	s.fields = fields
	return nil
}

func (s *Options) ToMap() bson.M {
	res := bson.M{}
	for _, field := range s.fields {
		res[field.Key] = field.Value
	}
	return res
}
