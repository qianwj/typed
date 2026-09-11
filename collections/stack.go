package collections

import "github.com/qianwj/typed/utils/option"

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

// Pop removes and returns the top element wrapped in a present
// Optional, or an absent Optional if the Stack is empty.
//
// Pop explicitly zeroes the popped slot before shortening the
// underlying slice. Without this, a Stack[T] where T holds
// references (pointers, slices, maps, channels, strings, funcs,
// or types with finalizers) would leak the popped element's
// references until the backing array is reallocated or the
// Stack itself becomes unreachable: the slot is still part of
// the backing array even after s.items shrinks past it, so the
// runtime's reachability walk keeps the reference alive.
func (s *Stack[T]) Pop() option.Optional[T] {
	n := len(s.items) - 1
	if n < 0 {
		return option.Empty[T]()
	}
	item := s.items[n]
	var zero T
	s.items[n] = zero
	s.items = s.items[:n]
	return option.Of(item)
}

// Peek returns the top element wrapped in a present Optional,
// or an absent Optional if the Stack is empty. Peek does not
// modify the Stack.
func (s *Stack[T]) Peek() option.Optional[T] {
	if len(s.items) == 0 {
		return option.Empty[T]()
	}
	return option.Of(s.items[len(s.items)-1])
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Clear() {
	s.items = make([]T, 0)
}
