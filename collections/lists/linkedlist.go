package lists

import (
	"slices"

	"github.com/qianwj/typed/collections/stream"
	"github.com/qianwj/typed/utils/option"
)

// ---------- LinkedList ----------

// node is one element of a doubly-linked list. The fields are unexported
// because users should not manipulate the linkage directly.
type node[T any] struct {
	value T
	prev  *node[T]
	next  *node[T]
}

// LinkedList is a doubly-linked list of T.
//
// As a concrete generic type, its methods (including type-changing ones
// such as Map[R]) can use Go 1.27's generic methods.
//
// All methods are defined on *LinkedList, not LinkedList, and constructors
// return *LinkedList, so the list is treated as a single reference-shared
// instance. Copying a LinkedList value copies only the head/tail/size fields
// and shares the underlying node chain; use the pointer returned by the
// constructor and do not copy a LinkedList after it has been put into use.
type LinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

// NewLinkedList returns a new empty LinkedList.
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// LinkedListOf returns a LinkedList containing the given values in order.
func LinkedListOf[T any](values ...T) *LinkedList[T] {
	l := NewLinkedList[T]()
	for _, v := range values {
		l.Add(v)
	}
	return l
}

// Add appends v to the end of the list.
func (l *LinkedList[T]) Add(v T) {
	n := &node[T]{value: v}
	if l.head == nil {
		l.head = n
	} else {
		n.prev = l.tail
		l.tail.next = n
	}
	l.tail = n
	l.size++
}

// AddFirst prepends v to the start of the list.
func (l *LinkedList[T]) AddFirst(v T) {
	n := &node[T]{value: v}
	if l.head == nil {
		l.tail = n
	} else {
		n.next = l.head
		l.head.prev = n
	}
	l.head = n
	l.size++
}

// RemoveFirst removes and returns the first element wrapped in a present
// option.Optional[T], or an absent Optional if the list is empty.
//
// RemoveFirst returns option.Optional[T] rather than (T, bool) so the
// "remove and get" path is symmetric with First / Last / Find, and
// callers can chain OrElse / Map / FlatMap on the removed value
// without first unpacking the result.
func (l *LinkedList[T]) RemoveFirst() option.Optional[T] {
	if l.head == nil {
		return option.Empty[T]()
	}
	v := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}
	l.size--
	return option.Of(v)
}

// RemoveLast removes and returns the last element wrapped in a present
// option.Optional[T], or an absent Optional if the list is empty.
//
// See RemoveFirst for the rationale behind returning Optional[T]
// rather than (T, bool).
func (l *LinkedList[T]) RemoveLast() option.Optional[T] {
	if l.tail == nil {
		return option.Empty[T]()
	}
	v := l.tail.value
	l.tail = l.tail.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	l.size--
	return option.Of(v)
}

// Insert inserts value at the given index. Elements at index and after are
// shifted one position to the right. If index == 0, value is prepended; if
// index == l.Size(), value is appended. Panics if index < 0 or
// index > l.Size().
func (l *LinkedList[T]) Insert(index int, value T) {
	if index < 0 || index > l.size {
		panic("LinkedList.Insert: index out of range")
	}
	switch index {
	case 0:
		l.AddFirst(value)
	case l.size:
		l.Add(value)
	default:
		// Walk to the node just before the insertion point.
		prev := l.head
		for range index - 1 {
			prev = prev.next
		}
		n := &node[T]{value: value, prev: prev, next: prev.next}
		prev.next.prev = n
		prev.next = n
		l.size++
	}
}

// RemoveAt removes and returns the element at the given index. Panics if
// index < 0 or index >= l.Size().
func (l *LinkedList[T]) RemoveAt(index int) T {
	if index < 0 || index >= l.size {
		panic("LinkedList.RemoveAt: index out of range")
	}
	switch index {
	case 0:
		// The bounds check at the top of RemoveAt rules out
		// an empty list at this point, so RemoveFirst always
		// returns a present value. Get() panics if that
		// invariant is ever broken, which is the right
		// behaviour for a logic error.
		return l.RemoveFirst().Get()
	case l.size - 1:
		return l.RemoveLast().Get()
	default:
		n := l.head
		for range index {
			n = n.next
		}
		v := n.value
		n.prev.next = n.next
		n.next.prev = n.prev
		l.size--
		return v
	}
}

// Get returns the value at index i wrapped in a present
// option.Optional[T], or an absent Optional if i is out of
// range. The walk is O(i), so random access on a linked list
// is slower than on an ArrayList.
//
// Get returns option.Optional[T] rather than (T, bool) so the
// result is symmetric with ArrayList.Get and with First / Last
// / Find on the same list, and so callers can chain OrElse /
// Map / FlatMap on the result.
func (l *LinkedList[T]) Get(i int) option.Optional[T] {
	if i < 0 || i >= l.size {
		return option.Empty[T]()
	}
	n := l.head
	for range i {
		n = n.next
	}
	return option.Of(n.value)
}

// First returns the first element wrapped in a present Optional,
// or an absent Optional if the LinkedList is empty.
//
// First returns option.Optional[T] rather than the (T, bool) shape
// so callers can chain the standard optional combinators
// (OrElse, Map, FlatMap, …) without first unpacking the result.
func (l *LinkedList[T]) First() option.Optional[T] {
	if l.head == nil {
		return option.Empty[T]()
	}
	return option.Of(l.head.value)
}

// Last returns the last element wrapped in a present Optional,
// or an absent Optional if the LinkedList is empty.
//
// See First for the rationale behind returning Optional[T] rather
// than (T, bool).
func (l *LinkedList[T]) Last() option.Optional[T] {
	if l.tail == nil {
		return option.Empty[T]()
	}
	return option.Of(l.tail.value)
}

// Size returns the number of elements.
func (l *LinkedList[T]) Size() int {
	return l.size
}

// IsEmpty reports whether the list contains no elements.
func (l *LinkedList[T]) IsEmpty() bool {
	return l.Size() == 0
}

// Clear removes all elements from the list.
func (l *LinkedList[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}

// Collect returns the list's values as a plain []T in head-to-tail order.
//
// This is the explicit way to take data out of a LinkedList now that all
// methods use *LinkedList receivers: copy the value into a slice the user
// owns, decoupled from the linked nodes.
func (l *LinkedList[T]) Collect() []T {
	out := make([]T, 0, l.size)
	for n := l.head; n != nil; n = n.next {
		out = append(out, n.value)
	}
	return out
}

// Stream returns a lazy stream.Stream[T] that is a snapshot of the list at
// the time Stream() is called. Mutations to the list after Stream() is
// called do not affect the Stream.
//
// The snapshot is materialized eagerly into a backing slice, so creating
// the Stream is O(n) in the list's size. Iteration of the Stream is then
// the same cost as iterating a plain slice. This cost is paid once at the
// LinkedList → Stream boundary, after which the Stream itself remains lazy.
func (l *LinkedList[T]) Stream() stream.Stream[T] {
	buf := make([]T, 0, l.size)
	for n := l.head; n != nil; n = n.next {
		buf = append(buf, n.value)
	}
	return stream.FromSlice(buf)
}

// ForEach invokes visit on every element in head-to-tail order.
func (l *LinkedList[T]) ForEach(visit func(T)) {
	for n := l.head; n != nil; n = n.next {
		visit(n.value)
	}
}

// Peek calls visit on each element and returns an independent copy of the
// list.
func (l *LinkedList[T]) Peek(visit func(T)) *LinkedList[T] {
	out := NewLinkedList[T]()
	for n := l.head; n != nil; n = n.next {
		visit(n.value)
		out.Add(n.value)
	}
	return out
}

// Any reports whether at least one element satisfies p. Stops as soon as a
// match is found.
func (l *LinkedList[T]) Any(p func(T) bool) bool {
	for n := l.head; n != nil; n = n.next {
		if p(n.value) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies p. Stops as soon as a
// mismatch is found.
func (l *LinkedList[T]) All(p func(T) bool) bool {
	for n := l.head; n != nil; n = n.next {
		if !p(n.value) {
			return false
		}
	}
	return true
}

// None reports whether no element satisfies p. Equivalent to !Any(p).
func (l *LinkedList[T]) None(p func(T) bool) bool {
	return !l.Any(p)
}

// Find returns the first element for which p returns true, wrapped
// in a present Optional, or an absent Optional if no element matches.
//
// See First for the rationale behind returning Optional[T] rather
// than (T, bool). Find is a short-circuiting terminal-style
// operation: it stops at the first match.
func (l *LinkedList[T]) Find(p func(T) bool) option.Optional[T] {
	for n := l.head; n != nil; n = n.next {
		if p(n.value) {
			return option.Of(n.value)
		}
	}
	return option.Empty[T]()
}

// Reduce folds the elements left-to-right using f, starting from init.
// The accumulator type U is independent of the element type T, so a
// Reduce can change the result type (e.g. building a string from a list
// of ints). Use U = T for the simple same-type fold.
func (l *LinkedList[T]) Reduce[U any](init U, f func(acc U, v T) U) U {
	acc := init
	for n := l.head; n != nil; n = n.next {
		acc = f(acc, n.value)
	}
	return acc
}

// Filter returns a new LinkedList keeping only the elements for which p
// returns true.
func (l *LinkedList[T]) Filter(p func(T) bool) *LinkedList[T] {
	out := NewLinkedList[T]()
	for n := l.head; n != nil; n = n.next {
		if p(n.value) {
			out.Add(n.value)
		}
	}
	return out
}

// Map applies f to every element and returns an ArrayList of the results.
//
// The mapped type is unconstrained and may be non-comparable. Results keep
// the linked list's head-to-tail order.
func (l *LinkedList[T]) Map[R any](f func(T) R) *ArrayList[R] {
	values := make([]R, 0, l.size)
	for n := l.head; n != nil; n = n.next {
		values = append(values, f(n.value))
	}
	return ArrayListOf(values...)
}

// FlatMap applies f to every value and concatenates the returned ArrayLists.
// A nil mapped list is treated as empty.
func (l *LinkedList[T]) FlatMap[R any](f func(T) *ArrayList[R]) *ArrayList[R] {
	var values []R
	for n := l.head; n != nil; n = n.next {
		mapped := f(n.value)
		if mapped == nil {
			continue
		}
		values = append(values, mapped.items...)
	}
	return ArrayListOf(values...)
}

// Take returns a new LinkedList with at most the first n elements.
func (l *LinkedList[T]) Take(n int) *LinkedList[T] {
	out := NewLinkedList[T]()
	if n <= 0 {
		return out
	}
	for current, taken := l.head, 0; current != nil && taken < n; current, taken = current.next, taken+1 {
		out.Add(current.value)
	}
	return out
}

// Drop returns a new LinkedList with the first n elements removed.
func (l *LinkedList[T]) Drop(n int) *LinkedList[T] {
	out := NewLinkedList[T]()
	if n < 0 {
		n = 0
	}
	for current, index := l.head, 0; current != nil; current, index = current.next, index+1 {
		if index >= n {
			out.Add(current.value)
		}
	}
	return out
}

// Distinct returns a new LinkedList keeping only the first occurrence of
// each element under eq.
func (l *LinkedList[T]) Distinct(eq func(T, T) bool) *LinkedList[T] {
	out := NewLinkedList[T]()
	l.ForEach(func(value T) {
		duplicate := false
		out.ForEach(func(existing T) {
			if eq(existing, value) {
				duplicate = true
			}
		})
		if !duplicate {
			out.Add(value)
		}
	})
	return out
}

// Concat returns a new LinkedList that appends other to l.
func (l *LinkedList[T]) Concat(other *LinkedList[T]) *LinkedList[T] {
	out := NewLinkedList[T]()
	l.ForEach(out.Add)
	if other != nil {
		other.ForEach(out.Add)
	}
	return out
}

// SortBy returns a new LinkedList with the elements ordered by less.
func (l *LinkedList[T]) SortBy(less func(x, y T) int) *LinkedList[T] {
	values := l.Collect()
	slices.SortFunc(values, less)
	return LinkedListOf(values...)
}

// MinBy returns the smallest element under less wrapped in a present
// option.Optional[T], or an absent Optional if the list is empty.
//
// MinBy returns option.Optional[T] rather than (T, bool) so the
// "find and get" path is symmetric with Find and ArrayList.MinBy.
func (l *LinkedList[T]) MinBy(less func(x, y T) int) option.Optional[T] {
	if l.head == nil {
		return option.Empty[T]()
	}
	best := l.head.value
	for n := l.head.next; n != nil; n = n.next {
		if less(n.value, best) < 0 {
			best = n.value
		}
	}
	return option.Of(best)
}

// MaxBy returns the largest element under less wrapped in a present
// option.Optional[T], or an absent Optional if the list is empty.
//
// See MinBy for the rationale.
func (l *LinkedList[T]) MaxBy(less func(x, y T) int) option.Optional[T] {
	if l.head == nil {
		return option.Empty[T]()
	}
	best := l.head.value
	for n := l.head.next; n != nil; n = n.next {
		if less(n.value, best) > 0 {
			best = n.value
		}
	}
	return option.Of(best)
}
