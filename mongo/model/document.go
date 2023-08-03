package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentId interface {
	~string | Number | primitive.ObjectID
}

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

type Document[T DocumentId] interface {
	GetId() T
}
