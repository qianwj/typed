package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry bson.E

func NewEntry(key string, val any) Entry {
	return Entry{Key: key, Value: val}
}

func (e Entry) Primitive() primitive.E {
	return primitive.E(e)
}
