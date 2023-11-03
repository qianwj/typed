package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnorderedMap bson.M

func NewUnorderdMap() UnorderedMap {
	return UnorderedMap{}
}

func (u UnorderedMap) Put(key string, val any) UnorderedMap {
	u[key] = val
	return u
}

func (u UnorderedMap) Get(key string) (any, bool) {
	val, exist := u[key]
	return val, exist
}

func (u UnorderedMap) Primitive() primitive.M {
	return primitive.M(u)
}
