package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

func (u Update) Set(update bson.M) Update {
	u["$set"] = update
	return u
}

func (u Update) SetOnInsert(update bson.M) Update {
	u["$setOnInsert"] = update
	return u
}
