package bson

import (
	"encoding/json"
)

type Map interface {
	Bson
	json.Marshaler
	Put(string, Bson) Map
	PutAll(other Map)
	Get(key string) (Bson, bool)
	Entries() []Entry
}
