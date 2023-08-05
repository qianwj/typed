package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentId interface {
	~string | Number | primitive.ObjectID
}

type Number interface {
	Int | Float
}

type Int interface {
	~int | ~int32 | ~int64
}

type Float interface {
	~float32 | ~float64
}

type Document[T DocumentId] interface {
	GetId() T
}

type BsonMap[I DocumentId] bson.M

func (m BsonMap[I]) GetId() I {
	var id I
	_id := bson.M(m)["_id"]
	if id, ok := _id.(I); ok {
		return id
	}
	return id
}

type Pair[V any] struct {
	Key   string
	Value V
}
