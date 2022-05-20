package list

import "github.com/qianwj/typed/collection"

type ArrayList[T any] struct {
	elements []T
}

const defaultInitialCapacity = 10

func NewArrayList[T any](e ...T) *ArrayList[T] {
	var elements []T
	if len(e) == 0 {
		elements = make([]T, 0, defaultInitialCapacity)
	} else {
		elements = e
	}
	return &ArrayList[T]{elements: elements}
}

func (a ArrayList[T]) Add(e T) {
	a.elements = append(a.elements, e)
}

func (a ArrayList[T]) AddAll(c collection.Collection[T]) {
	c.Foreach(func(e T) {
		a.Add(e)
	})
}

func (a ArrayList[T]) Foreach(handle func(e T)) {
	for _, e := range a.elements {
		handle(e)
	}
}

func (a ArrayList[T]) Size() int {
	return len(a.elements)
}
