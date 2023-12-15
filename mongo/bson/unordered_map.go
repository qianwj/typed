package bson

import (
	"go.mongodb.org/mongo-driver/bson"
)

type UnorderedMap bson.M

func M(e ...Entry) UnorderedMap {
	m := UnorderedMap{}
	for _, entry := range e {
		m[entry.Key] = entry.Value
	}
	return m
}

func (u UnorderedMap) Marshal() ([]byte, error) {
	return bson.Marshal(u)
}
