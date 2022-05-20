package collection

import "github.com/qianwj/typed/collection/list"

type Collection[T any] interface {
	Add(e T)
	AddAll(c Collection[T])
	Foreach(handle func(e T))
	Size() int
}

func ArrayListOf[T any](data []T) *list.ArrayList[T] {
	return list.NewArrayList[T](data...)
}
