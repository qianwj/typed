package bson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Null = primitive.Null{}
	Nil  = Null
)

type ID interface {
	~string | primitive.ObjectID
}

type Doc[I ID] interface {
	GetID() I
}
