package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

func NewUpdateMany() *UpdateMany {
	return &UpdateMany{}
}

func (m *UpdateMany) SetFilter(filter Filter) *UpdateMany {
	m.Filter = filter
	return m
}

func (m *UpdateMany) SetUpdate(update Update) *UpdateMany {
	m.Update = update
	return m
}

func NewFindOneAndUpdate() *FindOneAndUpdate {
	return &FindOneAndUpdate{}
}

func (m *FindOneAndUpdate) SetFilter(filter Filter) *FindOneAndUpdate {
	m.Filter = filter
	return m
}

func (m *FindOneAndUpdate) SetUpdate(update Update) *FindOneAndUpdate {
	m.Update = update
	return m
}

func (u Update) Set(update bson.M) Update {
	u["$set"] = update
	return u
}

func (u Update) SetOnInsert(update bson.M) Update {
	u["$setOnInsert"] = update
	return u
}
