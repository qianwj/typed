package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnorderedMap bson.M

func M(e ...Entry) UnorderedMap {
	m := UnorderedMap{}
	for _, entry := range e {
		m[entry.Key] = entry.Value
	}
	return m
}

func (u UnorderedMap) Primitive() primitive.M {
	return primitive.M(u)
}

func (u UnorderedMap) Marshal() ([]byte, error) {
	return bson.Marshal(u)
}
