package model

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/sorts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fill struct {
	bson.Bson
	partitionBy       string   //: <expression>,
	partitionByFields []string //: [ <field 1>, <field 2>, ... , <field n> ],
	sortBy            *sorts.SortOptions
	output            bson.Document
}

func New(output bson.Document) *Fill {
	return &Fill{
		output: output,
	}
}

func (f *Fill) PartitionBy(partitionBy string) *Fill {
	f.partitionBy = partitionBy
	return f
}

func (f *Fill) Tag() {}

func (f *Fill) Marshal() primitive.D {
	return primitive.D{
		{Key: "partitionBy", Value: f.partitionBy},
		{Key: "partitionByFields", Value: f.partitionByFields},
		{Key: "sortBy", Value: f.sortBy.Marshal()},
		{Key: "output", Value: f.output.Marshal()},
	}
}

func (f *Fill) ToMap() primitive.M {
	return primitive.M{
		"partitionBy":       f.partitionBy,
		"partitionByFields": f.partitionByFields,
		"sortBy":            f.sortBy.Marshal(),
		"output":            f.output.Marshal(),
	}
}
