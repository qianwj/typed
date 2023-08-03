package modify

import (
	"errors"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type findAndModify struct {
	filter      *filter.Filter
	set         bson.M
	setOnInsert bson.M
	unset       bson.M
	addToSets   []bson.M
	increments  bson.M
	decrements  bson.M
}

type (
	UpdateOne        findAndModify
	UpdateMany       findAndModify
	FindOneAndUpdate findAndModify
)

func New(f *filter.Filter) *findAndModify {
	return &findAndModify{
		filter:      f,
		set:         bson.M{},
		setOnInsert: bson.M{},
		unset:       bson.M{},
		increments:  bson.M{},
		decrements:  bson.M{},
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

func (f *findAndModify) Marshal() (bson.D, error) {
	update, holder := bson.D{}, false
	if len(f.set) > 0 {
		holder = true
		update = append(update, bson.E{
			Key:   operator.Set,
			Value: f.set,
		})
	}
	if len(f.setOnInsert) > 0 {
		holder = true
		update = append(update, bson.E{
			Key:   operator.SetOnInsert,
			Value: f.setOnInsert,
		})
	}
	if len(f.increments) > 0 {
		holder = true
		update = append(update, bson.E{
			Key:   operator.Inc,
			Value: f.increments,
		})
	}
	if len(f.decrements) > 0 {
		holder = true
		update = append(update, bson.E{
			Key:   operator.Inc,
			Value: f.decrements,
		})
	}
	// complete this function
	if !holder {
		return nil, errors.New("at least have one update operation")
	}
	return update, nil
}
