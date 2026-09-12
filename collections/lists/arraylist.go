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
	"github.com/qianwj/typed/utils/json"
	"github.com/qianwj/typed/utils/option"
)

// ---------- ArrayList ----------

// ArrayList is an ordered, resizable collection of T backed by a private
// []T and a head offset.
//
// As a concrete generic type, ArrayList's methods (including type-changing
// ones such as Map[R]) can take advantage of Go 1.27's generic methods.
//
// # Internal layout
//
// The backing slice items is laid out as [discarded prefix | live range |
// trailing zeroed slots]. head is the index of the first live element; the
// live count is len(items) - head. All public methods address elements by
// their *logical* position, which is translated to a physical position via
// head + i.
//
// # Why a head offset
//
// Before the head-offset refactor, ArrayList.RemoveFirst used
// slices.Delete(items, 0, 1), which shifts every remaining element down by
// one. That meant two things:
//
//  1. RemoveFirst was O(n): each call copied the entire tail. AddFirst was
//     the same shape via slices.Insert.
//  2. The slice's capacity never shrank after a head removal, so an
//     ArrayList that grew to a million elements and was then drained
//     from the front still held a million-element backing array for the
//     last few elements it carried.
//
// The head offset makes both operations O(1): RemoveFirst just bumps head
// past the discarded slot, and AddFirst writes at items[head-1] and
// decrements head. The capacity of the backing array is then bounded by
// the high-water mark of in-flight elements, not by the all-time maximum
// the ArrayList ever saw, because the array is reused and reset rather
// than grown monotonically.
//
// # Periodic compaction
//
// When the discarded prefix grows past a threshold (head >= 64 and head
// is at least half of len(items)), the live range is copied down to the
// start of items, head is reset to zero, and the now-vacated tail is
// zeroed. The threshold of 64 prevents short-lived lists from paying a
// copy on every RemoveFirst. Tests in arraylist_test.go drive a long
// sequence of push / pop-from-front cycles and assert that retained
// capacity drops back to the in-flight size, not the all-time maximum.
type ArrayList[T any] struct {
	items []T
	head  int
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

// size returns the number of live elements, which is the public Size.
// It is kept as an unexported helper to make the head-offset translation
// obvious in every method that needs it.
func (a *ArrayList[T]) size() int {
	return len(a.items) - a.head
}

// compact folds the live range down to the start of items and resets
// head to zero, releasing the discarded prefix back to the runtime.
//
// compact is called automatically when the discarded prefix grows past
// a threshold (see the type doc). It is a no-op when head is already
// zero. The vacating tail is zeroed so that any references the
// displaced slots held are released for GC.
func (a *ArrayList[T]) compact() {
	if a.head == 0 {
		return
	}
	n := copy(a.items, a.items[a.head:])
	var zero T
	for i := n; i < len(a.items); i++ {
		a.items[i] = zero
	}
	a.items = a.items[:n]
	a.head = 0
}

// compactIfNeeded is the threshold-driven entry point. It avoids
// paying the copy cost on lists that have not yet grown a meaningful
// discarded prefix.
//
// The threshold is "head >= 64": once the discarded prefix reaches
// 64 elements, the wasted capacity is large enough that folding
// it back to zero is worth the O(n) copy. With this rule the
// backing array's wasted prefix is bounded by a constant (64),
// so the capacity of a long-running head-drained list is bounded
// by the high-water mark of in-flight elements plus a constant,
// not by the all-time maximum. There is no second "head*2 >= len"
// check: the absolute threshold alone is enough, and a per-list
// ratio would over-compact small lists that have a small
// high-water mark anyway.
func (a *ArrayList[T]) compactIfNeeded() {
	if a.head >= 64 {
		a.compact()
	}
}

// Add appends value to the end of the list. Add is amortised O(1):
// append may grow the backing slice, but the amortised cost is
// constant per call.
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
	out := make([]T, a.size())
	copy(out, a.items[a.head:])
	return stream.FromSlice(out)
}

// Filter returns a new ArrayList containing only the elements for which p
// returns true.
func (a *ArrayList[T]) Filter(p func(T) bool) *ArrayList[T] {
	out := make([]T, 0, a.size())
	for _, v := range a.items[a.head:] {
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
	out := make([]R, a.size())
	for i, v := range a.items[a.head:] {
		out[i] = f(v)
	}
	return &ArrayList[R]{items: out}
}

// FlatMap applies f to every element and concatenates the resulting
// ArrayLists.
func (a *ArrayList[T]) FlatMap[R any](f func(T) *ArrayList[R]) *ArrayList[R] {
	var out []R
	for _, v := range a.items[a.head:] {
		mapped := f(v)
		if mapped == nil {
			continue
		}
		out = append(out, mapped.items[mapped.head:]...)
	}
	return &ArrayList[R]{items: out}
}

// Take returns a new ArrayList with at most the first n elements.
func (a *ArrayList[T]) Take(n int) *ArrayList[T] {
	size := a.size()
	if n <= 0 {
		return &ArrayList[T]{}
	}
	if n >= size {
		out := make([]T, size)
		copy(out, a.items[a.head:])
		return &ArrayList[T]{items: out}
	}
	out := make([]T, n)
	copy(out, a.items[a.head:a.head+n])
	return &ArrayList[T]{items: out}
}

// Drop returns a new ArrayList with the first n elements removed.
func (a *ArrayList[T]) Drop(n int) *ArrayList[T] {
	size := a.size()
	if n <= 0 {
		out := make([]T, size)
		copy(out, a.items[a.head:])
		return &ArrayList[T]{items: out}
	}
	if n >= size {
		return &ArrayList[T]{}
	}
	out := make([]T, size-n)
	copy(out, a.items[a.head+n:])
	return &ArrayList[T]{items: out}
}

// Distinct returns a new ArrayList keeping only the first occurrence of each
// element under eq.
func (a *ArrayList[T]) Distinct(eq func(T, T) bool) *ArrayList[T] {
	out := make([]T, 0, a.size())
outer:
	for _, v := range a.items[a.head:] {
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
	out := make([]T, 0, a.size()+other.size())
	out = append(out, a.items[a.head:]...)
	out = append(out, other.items[other.head:]...)
	return &ArrayList[T]{items: out}
}

// Peek calls visit on each element and returns a unchanged. Useful for
// debugging or observing a pipeline without modifying it.
func (a *ArrayList[T]) Peek(visit func(T)) *ArrayList[T] {
	for _, v := range a.items[a.head:] {
		visit(v)
	}
	out := make([]T, a.size())
	copy(out, a.items[a.head:])
	return &ArrayList[T]{items: out}
}

// AddFirst prepends value to the list. AddFirst is O(1) under the
// head-offset layout: the value is written at items[head-1] and head
// is decremented. (If head is already 0 the slice is grown to make
// room, which is amortised O(1).)
func (a *ArrayList[T]) AddFirst(value T) {
	if a.head > 0 {
		a.head--
		a.items[a.head] = value
		return
	}
	// head is 0: prepend by growing the slice and shifting the
	// live range right by one. The wasted slot at items[0] is
	// reclaimed on the next compaction.
	a.items = slices.Insert(a.items, 0, value)
	// After slices.Insert, head is still 0; the new element is
	// at items[0] and the live range is items[0:size+1]. No
	// translation needed.
}

// Insert inserts value at the given index. Elements at index and after
// are shifted one position to the right. If index == size, value is
// appended. If index == 0, value becomes the new first element.
// Panics if index < 0 or index > size.
func (a *ArrayList[T]) Insert(index int, value T) {
	size := a.size()
	if index < 0 || index > size {
		panic("ArrayList.Insert: index out of range")
	}
	if index == 0 {
		a.AddFirst(value)
		return
	}
	if index == size {
		a.Add(value)
		return
	}
	// index is strictly inside the live range. Grow by one
	// (with a zero fill — Go's copy uses memmove, so shifting
	// the live tail right by one in place is safe even when
	// the source and destination ranges overlap). The new
	// logical slot at items[head+index] is then filled with
	// value.
	var zero T
	a.items = append(a.items, zero)
	copy(a.items[a.head+index+1:], a.items[a.head+index:a.head+size])
	a.items[a.head+index] = value
}

// RemoveAt removes and returns the element at the given index. Subsequent
// elements are shifted one position to the left. Panics if index < 0 or
// index >= size.
func (a *ArrayList[T]) RemoveAt(index int) T {
	size := a.size()
	if index < 0 || index >= size {
		panic("ArrayList.RemoveAt: index out of range")
	}
	phys := a.head + index
	v := a.items[phys]
	if index == 0 {
		// Common case: head-side removal. Bump head and zero
		// the freed slot to release any reference it held.
		var zero T
		a.items[phys] = zero
		a.head++
		a.compactIfNeeded()
		return v
	}
	if index == size-1 {
		// Tail-side removal: just shrink the slice.
		var zero T
		a.items[phys] = zero
		a.items = a.items[:phys]
		return v
	}
	// Middle removal: shift the live tail left by one and zero
	// the vacated tail slot.
	var zero T
	copy(a.items[phys:phys+size-index-1], a.items[phys+1:phys+size-index])
	a.items[phys+size-index-1] = zero
	a.items = a.items[:size-1+a.head]
	return v
}

// RemoveFirst removes and returns the first element wrapped in a present
// option.Optional[T], or an absent Optional if the list is empty.
//
// RemoveFirst is O(1) under the head-offset layout: it advances
// head, zeroes the freed slot to release any reference it held,
// and triggers a periodic compaction when the discarded prefix
// grows past a threshold.
//
// RemoveFirst returns option.Optional[T] rather than (T, bool) so
// the result is symmetric with First / Last / Find, with
// LinkedList.RemoveFirst, and with Stack.Pop and Queue.Pop.
func (a *ArrayList[T]) RemoveFirst() option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	v := a.items[a.head]
	var zero T
	a.items[a.head] = zero
	a.head++
	a.compactIfNeeded()
	return option.Of(v)
}

// RemoveLast removes and returns the last element wrapped in a
// present option.Optional[T], or an absent Optional if the list
// is empty.
//
// RemoveLast is O(1): it shrinks the slice by one and zeroes
// the vacated slot. See RemoveFirst for the Optional rationale.
func (a *ArrayList[T]) RemoveLast() option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	last := a.head + a.size() - 1
	v := a.items[last]
	var zero T
	a.items[last] = zero
	a.items = a.items[:last]
	return option.Of(v)
}

// Clear removes all elements from the list. Clear resets the head
// offset to 0 and re-allocates a new (empty) backing slice, releasing
// every reference the live range held.
func (a *ArrayList[T]) Clear() {
	a.items = nil
	a.head = 0
}

// Collect returns the ArrayList's values as a freshly allocated []T.
//
// The returned slice is a copy, decoupled from the ArrayList's internal
// storage; mutating it does not affect the ArrayList. Collect does
// not include the discarded prefix; only the live range is copied.
func (a *ArrayList[T]) Collect() []T {
	out := make([]T, a.size())
	copy(out, a.items[a.head:])
	return out
}

// MarshalJSON encodes the ArrayList's live range as a JSON
// array in logical order. The discarded prefix is not part of
// the output: an ArrayList that has been drained from the
// front 100 times and then has 3 elements marshals to a
// 3-element array, not a 103-element array with 100 stale
// slots. This matches what Collect() returns and is the
// user-visible content of the list.
//
// An empty ArrayList marshals to "[]" rather than "null": a
// nil live range is normalised to the literal "[]" before
// being passed to utils/json.Encode. This matches the
// project's convention that "absent" and "empty" collections
// are indistinguishable in the JSON output.
//
// MarshalJSON delegates to utils/json.Encode, the project-
// wide wrapper around encoding/json/v2 that returns a
// result.Result[[]byte]. The wrapper is invoked via Unwrap to
// recover the ([]byte, error) shape that the standard
// json.Marshaler interface requires. An ArrayList whose
// element type T satisfies json.Marshaler (or v2's marshaler
// variant) is marshaled element-by-element; for T that does
// not implement MarshalJSON the default encoding for T
// applies.
func (a *ArrayList[T]) MarshalJSON() ([]byte, error) {
	live := a.items[a.head:]
	if live == nil {
		// A nil slice would marshal to "null"; normalise it
		// to the empty-array literal so the output is "[]".
		// This keeps the "absent" and "empty" cases
		// indistinguishable in the JSON.
		return []byte("[]"), nil
	}
	return json.Encode(live).Unwrap()
}

// UnmarshalJSON decodes a JSON array into the ArrayList,
// replacing any existing contents. The new elements become
// the live range starting at index 0; the head offset is
// reset to zero so the next Add / AddFirst / Get operates on
// the freshly decoded data without any leftover prefix.
//
// UnmarshalJSON delegates to utils/json.Decode and then
// plumbs the resulting []T directly into the ArrayList's
// backing slice. JSON null is accepted and treated as an
// empty array, matching the v1 json package's behaviour for
// slices: a nil-decoded ArrayList is empty (head == 0,
// items == nil).
//
// The return value is the underlying v2 error if data is not
// a JSON array or if any element fails to decode into T.
func (a *ArrayList[T]) UnmarshalJSON(data []byte) error {
	r := json.Decode[[]T](data)
	if err := r.Error(); err != nil {
		return err
	}
	a.items = r.Value()
	a.head = 0
	return nil
}

// Size returns the number of live elements in the list, not the
// capacity of the backing slice. Size is O(1).
func (a *ArrayList[T]) Size() int {
	return a.size()
}

// IsEmpty reports whether the list contains no elements.
func (a *ArrayList[T]) IsEmpty() bool {
	return a.size() == 0
}

// Get returns the value at index i, or the zero value and false if i is
// out of range. Get is O(1) on ArrayList since the data is contiguous.
func (a *ArrayList[T]) Get(i int) option.Optional[T] {
	if i < 0 || i >= a.size() {
		return option.Empty[T]()
	}
	return option.Of(a.items[a.head+i])
}

// ForEach invokes visit on every element in the live range.
func (a *ArrayList[T]) ForEach(visit func(T)) {
	for _, v := range a.items[a.head:] {
		visit(v)
	}
}

// First returns the first element wrapped in a present Optional,
// or an absent Optional if the ArrayList is empty.
//
// First returns option.Optional[T] rather than the (T, bool) shape
// so callers can chain the standard optional combinators
// (OrElse, Map, FlatMap, …) without first unpacking the result.
func (a *ArrayList[T]) First() option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	return option.Of(a.items[a.head])
}

// Last returns the last element wrapped in a present Optional,
// or an absent Optional if the ArrayList is empty.
//
// See First for the rationale behind returning Optional[T] rather
// than (T, bool).
func (a *ArrayList[T]) Last() option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	return option.Of(a.items[a.head+a.size()-1])
}

// Any reports whether at least one element satisfies p.
func (a *ArrayList[T]) Any(p func(T) bool) bool {
	for _, v := range a.items[a.head:] {
		if p(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies p.
func (a *ArrayList[T]) All(p func(T) bool) bool {
	for _, v := range a.items[a.head:] {
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

// Find returns the first element for which p returns true, wrapped
// in a present Optional, or an absent Optional if no element matches.
//
// See First for the rationale behind returning Optional[T] rather
// than (T, bool). Find is a short-circuiting terminal-style
// operation: it stops at the first match.
func (a *ArrayList[T]) Find(p func(T) bool) option.Optional[T] {
	for _, v := range a.items[a.head:] {
		if p(v) {
			return option.Of(v)
		}
	}
	return option.Empty[T]()
}

// Reduce folds the elements left-to-right using f, starting from init.
func (a *ArrayList[T]) Reduce[R any](init R, f func(acc R, v T) R) R {
	acc := init
	for _, v := range a.items[a.head:] {
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
	out := make([]T, a.size())
	copy(out, a.items[a.head:])
	slices.SortFunc(out, less)
	return &ArrayList[T]{items: out}
}

// MinBy returns the smallest element under less wrapped in a present
// option.Optional[T], or an absent Optional when the ArrayList is empty.
//
// MinBy uses option.Optional[T] rather than (T, bool) so the
// "find and get" path is symmetric with First / Last / Find and
// chains naturally with option.Map / option.FlatMap.
func (a *ArrayList[T]) MinBy(less func(x, y T) int) option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	best := a.items[a.head]
	for _, v := range a.items[a.head+1 : a.head+a.size()] {
		if less(v, best) < 0 {
			best = v
		}
	}
	return option.Of(best)
}

// MaxBy returns the largest element under less wrapped in a present
// option.Optional[T], or an absent Optional when the ArrayList is empty.
//
// See MinBy for the rationale.
func (a *ArrayList[T]) MaxBy(less func(x, y T) int) option.Optional[T] {
	if a.size() == 0 {
		return option.Empty[T]()
	}
	best := a.items[a.head]
	for _, v := range a.items[a.head+1 : a.head+a.size()] {
		if less(v, best) > 0 {
			best = v
		}
	}
	return option.Of(best)
}
