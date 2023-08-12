package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//var (
//	Null = primitive.Null{}
//	Nil  = Null
//)

type Bson interface {
	Cast() any
}

type Entry interface {
	Bson
	Marshal() primitive.E
	GetKey() string
	GetVal() Bson
}

type Document interface {
	Bson
	Marshal() primitive.D
	ToMap() Map
}

type Array interface {
	Bson
	Marshal() primitive.A
	Size() int
	Iter() ArrayIterator
}

type ArrayIterator interface {
	Next() (Bson, bool)
	HasNext() bool
}

type E bson.E

func (e E) Cast() any {
	return e
}

func (e E) Marshal() primitive.E {
	return (bson.E)(e)
}

func (e E) GetKey() string {
	return e.Key
}

func (e E) GetVal() Bson {
	return e.Value.(Bson)
}

type A []Bson

func (a A) Cast() any {
	return a
}

func (a A) Marshal() bson.A {
	res := make(bson.A, len(a))
	for i, e := range a {
		res[i] = e.Cast()
	}
	return res
}

func (a A) Size() int {
	return len(a)
}

func (a A) Iter() ArrayIterator {
	return &bsonAIter{
		data:  a,
		cur:   0,
		total: len(a),
	}
}

type bsonAIter struct {
	data  A
	cur   int
	total int
}

func (b *bsonAIter) HasNext() bool {
	return b.cur < b.total
}

func (b *bsonAIter) Next() (Bson, bool) {
	ele, ok := b.data[b.cur].(Bson)
	b.cur++
	return ele, ok
}
