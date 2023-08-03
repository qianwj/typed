package modify

import (
	"errors"
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Update struct {
	set         bson.M
	setOnInsert bson.M
	unset       bson.M
	addToSets   []bson.M
	increments  bson.M
	decrements  bson.M
}

func NewUpdate() Update {
	return Update{
		set:         bson.M{},
		setOnInsert: bson.M{},
		unset:       bson.M{},
		increments:  bson.M{},
		decrements:  bson.M{},
	}
}

func (f Update) Marshal() (bson.D, error) {
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
