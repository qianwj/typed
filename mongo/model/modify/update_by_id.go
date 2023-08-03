package modify

import (
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateById[T model.DocumentId] struct {
	id T
	Update
}

func NewUpdateById[T model.DocumentId](id T) *UpdateById[T] {
	return &UpdateById[T]{id: id}
}

func (f *UpdateById[T]) Set(key string, val any) *UpdateById[T] {
	f.set[key] = val
	return f
}

func (f *UpdateById[T]) SetAll(pairs bson.M) *UpdateById[T] {
	for k, v := range pairs {
		f.set[k] = v
	}
	return f
}

func (f *UpdateById[T]) SetOnInsert(key string, val any) *UpdateById[T] {
	f.setOnInsert[key] = val
	return f
}

func (f *UpdateById[T]) SetOnInsertAll(pairs bson.M) *UpdateById[T] {
	for k, v := range pairs {
		f.setOnInsert[k] = v
	}
	return f
}

func (f *UpdateById[T]) Inc(key string, delta int) *UpdateById[T] {
	f.increments[key] = delta
	return f
}

func (f *UpdateById[T]) Dec(key string, delta int) *UpdateById[T] {
	f.decrements[key] = -delta
	return f
}
