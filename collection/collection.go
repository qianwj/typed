package collection

import (
	_map "github.com/qianwj/typed/collection/map"
)

type Collection[T any] interface {
	Add(e T)
	AddAll(c Collection[T])
	Foreach(handle func(e T))
	ForeachIndexed(handle func(index int, e T))
	Slice() []T
	Size() int
	IsEmpty() bool
	Clear()
}

func NewMap[K comparable, V any]() _map.Map[K, V] {
	return make(_map.Map[K, V])
}
