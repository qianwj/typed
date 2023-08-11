package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bson interface {
	Tag()
}

type Entry interface {
	Bson
	Marshal() primitive.E
}

type Document interface {
	Bson
	Marshal() primitive.D
	ToMap() primitive.M
}

type Array interface {
	Bson
	Marshal() primitive.A
	Size() int
}

type E bson.E

func (e E) Tag() {}

func (e E) Marshal() primitive.E {
	return (bson.E)(e)
}

type D bson.D

func (d D) Tag() {}

func (d D) Marshal() primitive.D {
	return (bson.D)(d)
}

func (d D) ToMap() primitive.M {
	m := bson.M{}
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}

func (d D) IsExpression() {}

type M bson.M

func (m M) Tag() {}

func (m M) Marshal() bson.D {
	d := bson.D{}
	for k, v := range m {
		d = append(d, bson.E{Key: k, Value: v})
	}
	return d
}

func (m M) ToMap() bson.M {
	return (bson.M)(m)
}

type A bson.A

func (a A) Tag() {}

func (a A) Marshal() bson.A {
	return (bson.A)(a)
}

func (a A) Size() int {
	return len(a)
}
