package list

import "github.com/qianwj/typed/collection"

type ArrayList[T any] struct {
	collection.Collection[T]
	elements []T
}

const defaultInitialCapacity = 10

func ArrayListOf[T any](data []T) *ArrayList[T] {
	return NewArrayList[T](data...)
}

func NewArrayList[T any](e ...T) *ArrayList[T] {
	var elements []T
	if len(e) == 0 {
		elements = make([]T, 0, defaultInitialCapacity)
	} else {
		elements = e
	}
	return &ArrayList[T]{elements: elements}
}

func (a *ArrayList[T]) Add(e T) {
	a.elements = append(a.elements, e)
}

func (a *ArrayList[T]) AddAll(c collection.Collection[T]) {
	c.Foreach(func(e T) {
		a.Add(e)
	})
}

func (a *ArrayList[T]) Foreach(handle func(e T)) {
	for _, e := range a.elements {
		handle(e)
	}
}

func (a *ArrayList[T]) ForeachIndexed(handle func(index int, e T)) {
	for i, e := range a.elements {
		handle(i, e)
	}
}

func (a *ArrayList[T]) Filter(filter func(e T) bool) *ArrayList[T] {
	res := NewArrayList[T]()
	for _, e := range a.elements {
		if filter(e) {
			res.Add(e)
		}
	}
	return res
}

func (a *ArrayList[T]) Get(index int) T {
	return a.elements[index]
}

func (a *ArrayList[T]) Slice() []T {
	copied := a.elements
	return copied
}

func (a *ArrayList[T]) Clear() {
	a.elements = make([]T, 0, defaultInitialCapacity)
}

func (a *ArrayList[T]) Size() int {
	return len(a.elements)
}

func (a *ArrayList[T]) IsEmpty() bool {
	return a.Size() == 0
}
