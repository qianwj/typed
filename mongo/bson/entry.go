package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry bson.E

func E(key string, val any) Entry {
	return Entry{Key: key, Value: val}
}

func (e Entry) Primitive() primitive.E {
	return primitive.E(e)
}

func DocID[I ID](id I) Entry {
	return E("_id", id)
}
