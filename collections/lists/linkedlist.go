package lists

import (
	"github.com/qianwj/typed/collections/stream"
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
// All methods are defined on *LinkedList, not LinkedList, so the list is
// treated as a single reference-shared instance. Copying a LinkedList
// value copies only the head/tail/size fields and shares the underlying
// node chain — do not copy a LinkedList after it has been put into use.
type LinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

// NewLinkedList returns a new empty LinkedList.
func NewLinkedList[T any]() LinkedList[T] {
	return LinkedList[T]{}
}

// LinkedListOf returns a LinkedList containing the given values in order.
func LinkedListOf[T any](values ...T) LinkedList[T] {
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

// RemoveFirst removes and returns the first element, or the zero value and
// false if the list is empty.
func (l *LinkedList[T]) RemoveFirst() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	v := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}
	l.size--
	return v, true
}

// RemoveLast removes and returns the last element, or the zero value and
// false if the list is empty.
func (l *LinkedList[T]) RemoveLast() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	v := l.tail.value
	l.tail = l.tail.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	l.size--
	return v, true
}

// Insert inserts value at the given index. Elements at index and after are
// shifted one position to the right. If index == 0, value is prepended; if
// index == l.Len(), value is appended. Panics if index < 0 or
// index > l.Len().
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
// index < 0 or index >= l.Len().
func (l *LinkedList[T]) RemoveAt(index int) T {
	if index < 0 || index >= l.size {
		panic("LinkedList.RemoveAt: index out of range")
	}
	switch index {
	case 0:
		v, _ := l.RemoveFirst()
		return v
	case l.size - 1:
		v, _ := l.RemoveLast()
		return v
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

// Get returns the value at index i, or the zero value and false if i is
// out of range. The walk is O(i), so random access on a linked list is
// slower than on an ArrayList.
func (l *LinkedList[T]) Get(i int) (T, bool) {
	if i < 0 || i >= l.size {
		var zero T
		return zero, false
	}
	n := l.head
	for range i {
		n = n.next
	}
	return n.value, true
}

// First returns the first element, or the zero value and false if empty.
func (l *LinkedList[T]) First() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// Last returns the last element, or the zero value and false if empty.
func (l *LinkedList[T]) Last() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// Len returns the number of elements.
func (l *LinkedList[T]) Len() int {
	return l.size
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

// Find returns the first element for which p returns true, or the zero value
// and false if none match.
func (l *LinkedList[T]) Find(p func(T) bool) (T, bool) {
	for n := l.head; n != nil; n = n.next {
		if p(n.value) {
			return n.value, true
		}
	}
	var zero T
	return zero, false
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
func (l *LinkedList[T]) Filter(p func(T) bool) LinkedList[T] {
	out := NewLinkedList[T]()
	for n := l.head; n != nil; n = n.next {
		if p(n.value) {
			out.Add(n.value)
		}
	}
	return out
}

// Map applies f to every element and returns a new LinkedList of the
// results.
func (l *LinkedList[T]) Map[R any](f func(T) R) LinkedList[R] {
	out := NewLinkedList[R]()
	for n := l.head; n != nil; n = n.next {
		out.Add(f(n.value))
	}
	return out
}
