package collections

import (
	"encoding/json"

	"github.com/qianwj/typed/utils/option"
)

// Stack[T] is a Last-In-First-Out container of T, backed by a
// single private []T. The zero value is ready to use: NewStack
// is the canonical constructor but is not required, because the
// methods handle a nil backing slice the same way they handle an
// empty one.
//
// All methods use pointer receivers so the backing slice is
// shared between the caller's handle and the receiver; copying
// a Stack value by value would produce a second, independent
// stack. The "do not copy a Stack after first use" rule mirrors
// the one in lists.
type Stack[T any] struct {
	items []T
}

// NewStack returns a new, empty Stack[T] with a non-nil backing
// slice. Callers that want to avoid the make call can use the
// zero value directly:
//
//	var s Stack[int] // also valid; ready for Push / Pop / Peek
//	s.Push(1)
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// Push appends value to the top of the Stack. The backing
// slice grows as needed; Push is amortised O(1).
func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

// Pop removes and returns the top element wrapped in a present
// option.Optional[T], or an absent Optional if the Stack is
// empty.
//
// Pop explicitly zeroes the popped slot before shortening the
// underlying slice. Without this, a Stack[T] where T holds
// references (pointers, slices, maps, channels, strings, funcs,
// or types with finalizers) would leak the popped element's
// references until the backing array is reallocated or the
// Stack itself became unreachable: the slot is still part of
// the backing array even after s.items shrinks past it, so the
// runtime's reachability walk would keep the reference alive.
//
// The reference-handling test in stack_test.go drives a sequence
// of *finalizableBox elements through Push / Pop and asserts
// that all finalizers run. With the previous (non-zeroing) Pop
// the finalizers would never fire, because the backing array
// kept the popped pointers alive through out-of-range slots.
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

// Peek returns the top element wrapped in a present
// option.Optional[T], or an absent Optional if the Stack is
// empty. Peek does not modify the Stack: the same call
// repeated yields the same value, and Size is unchanged.
func (s *Stack[T]) Peek() option.Optional[T] {
	if len(s.items) == 0 {
		return option.Empty[T]()
	}
	return option.Of(s.items[len(s.items)-1])
}

// Size returns the number of elements currently in the Stack.
// Size is O(1).
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty reports whether the Stack has no elements. IsEmpty
// is the O(1) companion to Size and the standard guard before
// Pop or Peek for callers that prefer it to the absent-Optional
// return.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements and resets the Stack to the empty
// state. Clear re-allocates a new backing slice rather than
// reusing the existing one, so any references held by elements
// in the old slice are released. After Clear, the Stack
// behaves exactly as one returned by NewStack.
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0)
}

// MarshalJSON encodes the Stack's elements as a JSON array. The
// top of the Stack is the last element of the array, matching
// the natural "stack-as-list" view: a Stack built by Push(1),
// Push(2), Push(3) marshals to [1, 2, 3] and unmarshals back to
// a Stack where Pop() yields 3, then 2, then 1.
//
// MarshalJSON uses encoding/json under the hood. A Stack whose
// element type T does not satisfy json.Marshaler falls back to
// the default encoding for T, exactly as json.Marshal on a
// []T would.
func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.items)
}

// UnmarshalJSON decodes a JSON array into the Stack, replacing
// any existing contents. The first element of the array becomes
// the bottom of the Stack; the last element becomes the top.
// This is the inverse of MarshalJSON: a round-trip of a Stack
// preserves its push order.
//
// UnmarshalJSON returns the underlying json error if data is
// not a JSON array or if any element fails to decode into T.
// On error the Stack is left in an unspecified state; callers
// that need to recover should treat the Stack as invalid and
// replace it.
func (s *Stack[T]) UnmarshalJSON(data []byte) error {
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	s.items = items
	return nil
}
