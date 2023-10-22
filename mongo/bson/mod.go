package bson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Null = primitive.Null{}
	Nil  = Null
)

type ID interface {
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

type Doc[I ID] interface {
	GetID() I
}
