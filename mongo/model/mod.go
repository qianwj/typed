package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Filter     bson.M
	Update     bson.M
	UpdateSome struct {
		Filter Filter
		Update Update
	}
	FindOneAndUpdate UpdateSome
	UpdateOne        UpdateSome
	UpdateMany       UpdateSome
	UpdateById       struct {
		Id     primitive.ObjectID
		Update Update
	}
	FieldType int32
)
