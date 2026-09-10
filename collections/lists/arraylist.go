// Package lists provides ArrayList[T] and LinkedList[T], two type-safe
// concrete generic collections with fluent operations.
//
// Both types are concrete (not interfaces) so that their methods, including
// type-changing ones such as Map[R], can take advantage of Go 1.27's
// generic methods. Interface methods still cannot declare type parameters,
// and a useful List interface (one carrying Map, Filter, Reduce, …) would
// have to break fluent chaining by returning the interface type from
// methods that today return the concrete type. The package therefore
// exposes the concrete types only; callers pick *ArrayList[T] or
// *LinkedList[T] directly.
//
// All methods on both types use pointer receivers (*ArrayList[T] and
// *LinkedList[T]). Receivers are never value-copied, so the rule "do not
// copy a lists value after putting it into use" is uniform across the
// package. ArrayList is a struct that wraps a private []T; LinkedList is a
// struct with private node pointers. Neither exposes its internal slice
// or node chain, so callers cannot accidentally bypass the API to mutate
// state.
//
// Both collections are eager: their intermediate operations return new
// collections immediately. For lazy pipelines, call Stream() to obtain a
// stream.Stream[T].
package lists

import (
	"slices"

	"github.com/qianwj/typed/collections/stream"
)

// ---------- ArrayList ----------

// ArrayList is an ordered, resizable collection of T backed by a private
// []T.
//
// The underlying slice is not exported; callers cannot index into it or
// otherwise bypass the API. All reads and writes go through the methods
// defined on *ArrayList. To obtain the values as a plain []T, call Collect
// (which returns a copy).
//
// As a concrete generic type, ArrayList's methods (including type-changing
// ones such as Map[R]) can take advantage of Go 1.27's generic methods.
type ArrayList[T any] struct {
	items []T
}

// NewArrayList returns a new empty ArrayList.
func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{}
}

// ArrayListOf returns an ArrayList containing the given values in order.
//
// A plain slice can be passed via the spread operator: ArrayListOf(s...).
// The values are stored without copying; the backing array is the variadic
// slice that ArrayListOf received.
func ArrayListOf[T any](values ...T) *ArrayList[T] {
	return &ArrayList[T]{items: values}
}

// Add appends value to the end of the list.
func (a *ArrayList[T]) Add(value T) {
	a.items = append(a.items, value)
}

// Stream returns a lazy stream.Stream[T] that is a snapshot of the
// ArrayList at the time Stream() is called. Mutations to the ArrayList
// after Stream() is called do not affect the Stream.
//
// To make the snapshot contract hold, Stream copies the values into a
// fresh backing slice. The copy is O(n); iteration of the Stream is then
// the same cost as iterating a plain slice. This cost is paid once at the
// ArrayList → Stream boundary, after which the Stream itself remains lazy.
//
// The Stream is single-use.
func (a *ArrayList[T]) Stream() stream.Stream[T] {
	out := make([]T, len(a.items))
	copy(out, a.items)
	return stream.FromSlice(out)
}

// Filter returns a new ArrayList containing only the elements for which p
// returns true.
func (a *ArrayList[T]) Filter(p func(T) bool) *ArrayList[T] {
	out := make([]T, 0, len(a.items))
	for _, v := range a.items {
		if p(v) {
			out = append(out, v)
		}
	}
	return &ArrayList[T]{items: out}
}

// Map applies f to every element and returns a new ArrayList of the results.
//
// This is the operation that motivates ArrayList being a concrete generic
// type in Go 1.27: the method declares its own type parameter R, which an
// interface method cannot do.
func (a *ArrayList[T]) Map[R any](f func(T) R) *ArrayList[R] {
	out := make([]R, len(a.items))
	for i, v := range a.items {
		out[i] = f(v)
	}
	return &ArrayList[R]{items: out}
}

// FlatMap applies f to every element and concatenates the resulting
// ArrayLists.
func (a *ArrayList[T]) FlatMap[R any](f func(T) *ArrayList[R]) *ArrayList[R] {
	var out []R
	for _, v := range a.items {
		mapped := f(v)
		if mapped == nil {
			continue
		}
		out = append(out, mapped.items...)
	}
	return &ArrayList[R]{items: out}
}

// Take returns a new ArrayList with at most the first n elements.
func (a *ArrayList[T]) Take(n int) *ArrayList[T] {
	if n <= 0 {
		return &ArrayList[T]{}
	}
	if n >= len(a.items) {
		out := make([]T, len(a.items))
		copy(out, a.items)
		return &ArrayList[T]{items: out}
	}
	out := make([]T, n)
	copy(out, a.items[:n])
	return &ArrayList[T]{items: out}
}

// Drop returns a new ArrayList with the first n elements removed.
func (a *ArrayList[T]) Drop(n int) *ArrayList[T] {
	if n <= 0 {
		out := make([]T, len(a.items))
		copy(out, a.items)
		return &ArrayList[T]{items: out}
	}
	if n >= len(a.items) {
		return &ArrayList[T]{}
	}
	out := make([]T, len(a.items)-n)
	copy(out, a.items[n:])
	return &ArrayList[T]{items: out}
}

// Distinct returns a new ArrayList keeping only the first occurrence of each
// element under eq.
func (a *ArrayList[T]) Distinct(eq func(T, T) bool) *ArrayList[T] {
	out := make([]T, 0, len(a.items))
outer:
	for _, v := range a.items {
		for _, x := range out {
			if eq(x, v) {
				continue outer
			}
		}
		out = append(out, v)
	}
	return &ArrayList[T]{items: out}
}

// Concat returns a new ArrayList that appends other to a.
func (a *ArrayList[T]) Concat(other *ArrayList[T]) *ArrayList[T] {
	out := make([]T, 0, len(a.items)+len(other.items))
	out = append(out, a.items...)
	out = append(out, other.items...)
	return &ArrayList[T]{items: out}
}

// Peek calls visit on each element and returns a unchanged. Useful for
// debugging or observing a pipeline without modifying it.
func (a *ArrayList[T]) Peek(visit func(T)) *ArrayList[T] {
	for _, v := range a.items {
		visit(v)
	}
	out := make([]T, len(a.items))
	copy(out, a.items)
	return &ArrayList[T]{items: out}
}

// AddFirst prepends value to the list.
func (a *ArrayList[T]) AddFirst(value T) {
	a.items = slices.Insert(a.items, 0, value)
}

// Insert inserts value at the given index. Elements at index and after are
// shifted one position to the right. If index == len(a), value is appended.
// Panics if index < 0 or index > len(a).
func (a *ArrayList[T]) Insert(index int, value T) {
	if index < 0 || index > len(a.items) {
		panic("ArrayList.Insert: index out of range")
	}
	a.items = slices.Insert(a.items, index, value)
}

// RemoveAt removes and returns the element at the given index. Subsequent
// elements are shifted one position to the left. Panics if index < 0 or
// index >= len(a).
func (a *ArrayList[T]) RemoveAt(index int) T {
	if index < 0 || index >= len(a.items) {
		panic("ArrayList.RemoveAt: index out of range")
	}
	v := a.items[index]
	a.items = slices.Delete(a.items, index, index+1)
	return v
}

// RemoveFirst removes and returns the first element, or the zero value and
// false if the list is empty.
func (a *ArrayList[T]) RemoveFirst() (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	v := a.items[0]
	a.items = slices.Delete(a.items, 0, 1)
	return v, true
}

// RemoveLast removes and returns the last element, or the zero value and
// false if the list is empty.
func (a *ArrayList[T]) RemoveLast() (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	last := len(a.items) - 1
	v := a.items[last]
	a.items = slices.Delete(a.items, last, len(a.items))
	return v, true
}

// Clear removes all elements from the list.
func (a *ArrayList[T]) Clear() {
	a.items = nil
}

// Collect returns the ArrayList's values as a freshly allocated []T.
//
// The returned slice is a copy, decoupled from the ArrayList's internal
// storage; mutating it does not affect the ArrayList.
func (a *ArrayList[T]) Collect() []T {
	out := make([]T, len(a.items))
	copy(out, a.items)
	return out
}

// Size returns the number of elements.
func (a *ArrayList[T]) Size() int {
	return len(a.items)
}

// IsEmpty reports whether the list contains no elements.
func (a *ArrayList[T]) IsEmpty() bool {
	return a.Size() == 0
}

// Get returns the value at index i, or the zero value and false if i is
// out of range. O(1) on ArrayList since the data is contiguous.
func (a *ArrayList[T]) Get(i int) (T, bool) {
	if i < 0 || i >= len(a.items) {
		var zero T
		return zero, false
	}
	return a.items[i], true
}

// ForEach invokes visit on every element.
func (a *ArrayList[T]) ForEach(visit func(T)) {
	for _, v := range a.items {
		visit(v)
	}
}

// First returns the first element, or the zero value and false if empty.
func (a *ArrayList[T]) First() (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	return a.items[0], true
}

// Last returns the last element, or the zero value and false if empty.
func (a *ArrayList[T]) Last() (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	return a.items[len(a.items)-1], true
}

// Any reports whether at least one element satisfies p.
func (a *ArrayList[T]) Any(p func(T) bool) bool {
	return slices.ContainsFunc(a.items, p)
}

// All reports whether every element satisfies p.
func (a *ArrayList[T]) All(p func(T) bool) bool {
	for _, v := range a.items {
		if !p(v) {
			return false
		}
	}
	return true
}

// None reports whether no element satisfies p. Equivalent to !Any(p).
func (a *ArrayList[T]) None(p func(T) bool) bool {
	return !a.Any(p)
}

// Find returns the first element for which p returns true, or the zero value
// and false if none match.
func (a *ArrayList[T]) Find(p func(T) bool) (T, bool) {
	for _, v := range a.items {
		if p(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Reduce folds the elements left-to-right using f, starting from init.
func (a *ArrayList[T]) Reduce[R any](init R, f func(acc R, v T) R) R {
	acc := init
	for _, v := range a.items {
		acc = f(acc, v)
	}
	return acc
}

// SortBy returns a new ArrayList with the elements ordered by less.
//
// A no-arg Sort() that uses cmp.Compare would require T cmp.Ordered, which
// a single ArrayList[T any] type cannot express on its method set. Provide
// the comparator explicitly to keep the fluent API type-parameter-free.
func (a *ArrayList[T]) SortBy(less func(x, y T) int) *ArrayList[T] {
	out := make([]T, len(a.items))
	copy(out, a.items)
	slices.SortFunc(out, less)
	return &ArrayList[T]{items: out}
}

// MinBy returns the smallest element under less, or the zero value and false
// if the ArrayList is empty.
func (a *ArrayList[T]) MinBy(less func(x, y T) int) (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	best := a.items[0]
	for _, v := range a.items[1:] {
		if less(v, best) < 0 {
			best = v
		}
	}
	return best, true
}

// MaxBy returns the largest element under less, or the zero value and false
// if the ArrayList is empty.
func (a *ArrayList[T]) MaxBy(less func(x, y T) int) (T, bool) {
	if len(a.items) == 0 {
		var zero T
		return zero, false
	}
	best := a.items[0]
	for _, v := range a.items[1:] {
		if less(v, best) > 0 {
			best = v
		}
	}
	return best, true
}
