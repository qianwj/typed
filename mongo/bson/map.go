package bson

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

type Map interface {
	Bson
	json.Marshaler
	bson.Marshaler
	Put(string, Bson) Map
	PutAll(other Map)
	Get(key string) (Bson, bool)
	Entries() []Entry
}
