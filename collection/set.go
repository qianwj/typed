package collection

type Set[T comparable] struct {
	Collection[T]
	elements Map[T, struct{}]
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(Map[T, struct{}]),
	}
}

func (s *Set[T]) AddAll(c Collection[T]) *Set[T] {
	c.Foreach(func(e T) {
		s.elements[e] = struct{}{}
	})
	return s
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
