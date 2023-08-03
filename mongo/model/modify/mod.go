package modify

import (
	"github.com/qianwj/typed/mongo/model/filter"
	"go.mongodb.org/mongo-driver/bson"
)

type findAndModify struct {
	filter *filter.Filter
	Update
}

type (
	UpdateOne        findAndModify
	UpdateMany       findAndModify
	FindOneAndUpdate findAndModify
)

func New(f *filter.Filter) *findAndModify {
	return &findAndModify{
		filter: f,
		Update: NewUpdate(),
	}
}

func (f *findAndModify) Set(key string, val any) *findAndModify {
	f.set[key] = val
	return f
}

func (f *findAndModify) SetAll(pairs bson.M) *findAndModify {
	for k, v := range pairs {
		f.set[k] = v
	}
	return f
}

func (f *findAndModify) SetOnInsert(key string, val any) *findAndModify {
	f.setOnInsert[key] = val
	return f
}

func (f *findAndModify) SetOnInsertAll(pairs bson.M) *findAndModify {
	for k, v := range pairs {
		f.setOnInsert[k] = v
	}
	return f
}

func (f *findAndModify) Inc(key string, delta int) *findAndModify {
	f.increments[key] = delta
	return f
}

func (f *findAndModify) Dec(key string, delta int) *findAndModify {
	f.decrements[key] = -delta
	return f
}

func (f *findAndModify) UpdateOne() *UpdateOne {
	return (*UpdateOne)(f)
}

func (f *findAndModify) UpdateMany() *UpdateMany {
	return (*UpdateMany)(f)
}

func (f *findAndModify) FindOneAndUpdate() *FindOneAndUpdate {
	return (*FindOneAndUpdate)(f)
}

func (f *findAndModify) Filter() *filter.Filter {
	return f.filter
}
