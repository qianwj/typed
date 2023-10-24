package group

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID interface {
	isGroupId()
}

type StringGroupId string

func (s StringGroupId) isGroupId() {}

func (s StringGroupId) MarshalBSON() ([]byte, error) {
	return bson.Marshal(string(s))
}

type ComplexGroupId primitive.M

func (c ComplexGroupId) isGroupId() {}

func (c ComplexGroupId) MarshalBSON() ([]byte, error) {
	return bson.Marshal(primitive.M(c))
}
