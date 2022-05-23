package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

func NewUpdate() Update {
	return Update{}
}

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

func (u Update) AddToSet(field string, value interface{}) Update {
	u["$addToSet"] = bson.M{
		field: value,
	}
	return u
}

func (u Update) AddEachToSet(field string, elements []interface{}) Update {
	u["$addToSet"] = bson.M{
		field: bson.M{
			"$each": elements,
		},
	}
	return u
}

func (u Update) Pull(filter Filter) Update {
	u["$pull"] = filter
	return u
}

func (u Update) PopFirst(field string) Update {
	return u.pop(field, -1)
}

func (u Update) PopLast(field string) Update {
	return u.pop(field, 1)
}

func (u Update) pop(field string, pos int) Update {
	u["$pop"] = bson.M{
		field: pos,
	}
	return u
}
