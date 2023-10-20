package sorts

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SortOrder int

const (
	Desc SortOrder = -1
	Asc  SortOrder = 1
)

type SortOptions struct {
	fields bson.D
}

func New() *SortOptions {
	return &SortOptions{fields: bson.D{}}
}

func Ascending(field string) *SortOptions {
	return New().Ascending(field)
}

func Descending(field string) *SortOptions {
	return New().Descending(field)
}

func Meta(field string) *SortOptions {
	return New().Meta(field)
}

func (s *SortOptions) Ascending(field string) *SortOptions {
	s.fields = append(s.fields, bson.E{Key: field, Value: Asc})
	return s
}

func (s *SortOptions) Descending(field string) *SortOptions {
	s.fields = append(s.fields, bson.E{Key: field, Value: Desc})
	return s
}

func (s *SortOptions) Meta(field string) *SortOptions {
	s.fields = append(s.fields, bson.E{Key: field, Value: bson.M{operator.Meta: "textScore"}})
	return s
}

func (s *SortOptions) Marshal() primitive.D {
	if s == nil {
		return bson.D{}
	}
	return s.fields
}

func (s *SortOptions) ToMap() primitive.M {
	res := bson.M{}
	for _, field := range s.fields {
		res[field.Key] = field.Value
	}
	return res
}
