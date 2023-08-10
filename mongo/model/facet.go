package model

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Facet struct {
	bson.Entry
	name     string
	pipeline []primitive.D
}

func NewFacet(name string, pipeline ...bson.Document) *Facet {
	return &Facet{
		name: name,
		pipeline: util.Map(pipeline, func(doc bson.Document) primitive.D {
			return doc.Marshal()
		}),
	}
}

func (f *Facet) Tag() {}

func (f *Facet) Marshal() primitive.E {
	return primitive.E{
		Key: f.name, Value: f.pipeline,
	}
}
