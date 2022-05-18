package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Filter           bson.M
	Update           bson.M
	FindOneAndUpdate struct {
		Filter Filter
		Update Update
	}
	UpdateMany struct {
		Filter Filter
		Update Update
	}
	UpdateById struct {
		Id     primitive.ObjectID
		Update Update
	}
)
