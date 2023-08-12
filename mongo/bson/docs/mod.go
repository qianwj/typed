package docs

import (
	"github.com/qianwj/typed/mongo/bson"
	raw "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type D raw.D

func (d D) Cast() any {
	return d
}

func (d D) Marshal() primitive.D {
	return (raw.D)(d)
}

func (d D) ToMap() bson.Map {
	m := New()
	for _, e := range d {
		m.Put(e.Key, e.Value.(bson.Bson))
	}
	return m
}

func (d D) IsExpression() {}

type M raw.M

func (m M) Cast() any {
	return m
}

func (m M) Marshal() raw.D {
	d := raw.D{}
	for k, v := range m {
		d = append(d, raw.E{Key: k, Value: v})
	}
	return d
}

func (m M) ToMap() bson.Map {
	res := New()
	for k, v := range m {
		res.Put(k, v.(bson.Bson))
	}
	return res
}
