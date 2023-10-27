package bson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Null = primitive.Null{}
	Nil  = Null
)

// ID is a constraint that defines the type of the `_id` field in mongo documents.
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

// Doc is an interface that defines the type of the mongo document.
// If you use `TypedCollection`, your document type must implement this interface.
type Doc[I ID] interface {
	GetID() I
}
