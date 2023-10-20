package bson

import "go.mongodb.org/mongo-driver/bson"

type Entry bson.E

func (e Entry) MarshalBJSON() ([]byte, error) {
	return bson.Marshal(e)
}
