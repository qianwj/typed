package set

import (
	"github.com/qianwj/typed/collection"
	"github.com/qianwj/typed/collection/map"
)

type Set[T comparable] struct {
	collection.Collection[T]
	elements _map.Map[T, struct{}]
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(_map.Map[T, struct{}]),
	}
}

func (s *Set[T]) AddAll(c collection.Collection[T]) {
	c.Foreach(func(e T) {
		s.elements[e] = struct{}{}
	})
}

func (s *Set[T]) Add(e T) {
	s.elements[e] = struct{}{}
}

func (s *Set[T]) Contains(e T) bool {
	_, ok := s.elements[e]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.elements)
}

func (s *Set[T]) Foreach(handle func(e T)) {
	for k := range s.elements {
		handle(k)
	}
}

func (s Set[T]) Slice() []T {
	return s.elements.Keys()
}
