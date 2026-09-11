package collections

import (
	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/utils/option"
)

// Deque[T] is a double-ended queue of T, implemented as a thin
// wrapper over LinkedList[T].
//
// A doubly-linked list gives O(1) head and tail operations
// without any auxiliary bookkeeping: AddFirst and RemoveFirst
// are constant-time pointer updates around the head pointer;
// Add and RemoveLast are constant-time updates around the tail
// pointer; First and Last are direct reads of the head and
// tail pointers. There is no periodic compaction, no amortised
// cost, and no worst-case O(n) operation: every public method
// on Deque is exactly O(1).
//
// # Why LinkedList and not ArrayList
//
// ArrayList's head-offset layout also gives O(1) on both ends,
// but only amortised: RemoveFirst is constant, but the
// associated periodic compaction is O(n) for that one call. A
// Deque whose API documents "O(1) on both ends" should not
// hide an occasional O(n) call behind the word "amortised".
// LinkedList avoids that concern by being O(1) worst-case at
// both ends, which is what callers expect from a deque.
//
// The trade-off is per-node memory: each LinkedList node
// carries two *node[T] pointers and the value, which is
// 16–24 bytes of overhead on 64-bit platforms in addition to
// the value's size. For Deque's typical workloads (BFS,
// monotonic queues, sliding windows, work-stealing buffers)
// the per-element cost is usually negligible against the
// per-call simplicity, and iteration is rare compared to
// push/pop. Queue is still ArrayList-backed because Queue only
// touches one end, so ArrayList's per-element memory wins are
// kept; Deque touches both, so LinkedList's worst-case O(1)
// wins are kept.
//
// # Naming
//
// Deque uses the standard "Front" / "Back" qualifiers rather
// than the bare "Push" / "Pop" / "Peek" used by Stack and
// Queue: those single-ended collections do not need to
// disambiguate which end, and a Deque would otherwise have no
// way to ask "push to the front" or "pop from the back".
//
// All accessors (Front, Back, PopFront, PopBack) return
// option.Optional[T] rather than (T, bool), matching the
// convention established by Stack, Queue, and the lists
// package: a present Optional on a hit, an absent one on a
// miss. Callers can chain OrElse / OrElseGet / Map on the
// return value without first unpacking.
//
// # Receivers and the zero value
//
// All methods use pointer receivers so the embedded LinkedList
// is shared between the caller's handle and the receiver;
// copying a Deque value by value would produce a second,
// independent deque. The "do not copy a Deque after first use"
// rule mirrors the one in Queue and lists.
//
// The zero value of Deque[T] is also usable: a freshly
// declared 'var d Deque[int]' accepts PushFront / PushBack /
// PopFront / PopBack without going through NewDeque. The
// embedded LinkedList is itself usable as a zero value (its
// head / tail pointers are nil, but every method handles that
// case).
type Deque[T any] struct {
	items lists.LinkedList[T]
}

// NewDeque returns a new, empty Deque[T].
//
// The zero value of Deque[T] is also usable: a freshly
// declared 'var d Deque[int]' accepts PushFront / PushBack /
// PopFront / PopBack without going through NewDeque. Internally,
// the first call lazily writes to d.items, so the zero value
// is treated as 'a Deque whose LinkedList has not been touched
// yet'.
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{}
}

// PushFront prepends value to the front of the Deque. PushFront
// is O(1): it allocates a new node and links it in front of
// the current head.
func (d *Deque[T]) PushFront(value T) {
	d.items.AddFirst(value)
}

// PushBack appends value to the back of the Deque. PushBack is
// O(1): it allocates a new node and links it behind the
// current tail.
func (d *Deque[T]) PushBack(value T) {
	d.items.Add(value)
}

// PopFront removes and returns the front element wrapped in a
// present option.Optional[T], or an absent Optional if the
// Deque is empty.
//
// PopFront is O(1): it unlinks the head node, releases the
// reference the node held, and returns the value. The node is
// then unreachable and eligible for GC.
func (d *Deque[T]) PopFront() option.Optional[T] {
	return d.items.RemoveFirst()
}

// PopBack removes and returns the back element wrapped in a
// present option.Optional[T], or an absent Optional if the
// Deque is empty.
//
// PopBack is O(1): it unlinks the tail node, releases the
// reference the node held, and returns the value.
func (d *Deque[T]) PopBack() option.Optional[T] {
	return d.items.RemoveLast()
}

// Front returns the front element wrapped in a present
// option.Optional[T], or an absent Optional if the Deque is
// empty. Front is O(1) and does not modify the Deque: the same
// call repeated yields the same value, and Size is unchanged.
func (d *Deque[T]) Front() option.Optional[T] {
	return d.items.First()
}

// Back returns the back element wrapped in a present
// option.Optional[T], or an absent Optional if the Deque is
// empty. Back is O(1) and does not modify the Deque.
func (d *Deque[T]) Back() option.Optional[T] {
	return d.items.Last()
}

// Size returns the number of elements currently in the Deque.
// Size is O(1) and delegates to LinkedList.Size.
func (d *Deque[T]) Size() int {
	return d.items.Size()
}

// IsEmpty reports whether the Deque has no elements. IsEmpty is
// the O(1) companion to Size and the standard guard before
// PopFront, PopBack, Front, or Back for callers that prefer it
// to the absent-Optional return.
func (d *Deque[T]) IsEmpty() bool {
	return d.items.IsEmpty()
}

// Clear removes all elements and resets the Deque to the empty
// state. Clear delegates to LinkedList.Clear, which sets both
// head and tail to nil and resets the size counter. After
// Clear, the Deque behaves exactly as one returned by NewDeque.
func (d *Deque[T]) Clear() {
	d.items.Clear()
}
