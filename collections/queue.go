package collections

import (
	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/utils/option"
)

// Queue[T] is a First-In-First-Out container of T, implemented as
// a thin wrapper over ArrayList[T]. The ArrayList already has a
// head-offset layout with periodic compaction, so Queue inherits
// O(1) Push, O(1) Pop, and bounded memory without duplicating
// the storage machinery.
//
// All methods use pointer receivers so the Queue's embedded
// ArrayList is shared between the caller's handle and the
// receiver; copying a Queue value by value would produce a
// second, independent queue. The "do not copy a Queue after
// first use" rule mirrors the one in lists.
//
// # Memory model
//
// Queue delegates to ArrayList, which keeps the live range in
// items[head:head+size] and folds the discarded prefix back to
// zero when head grows past a threshold of 64. The retained
// capacity of a long-running Queue is therefore bounded by the
// high-water mark of in-flight elements plus a constant, not by
// the all-time maximum the Queue ever saw. See ArrayList's type
// doc for the full rationale.
type Queue[T any] struct {
	items lists.ArrayList[T]
}

// NewQueue returns a new, empty Queue[T].
//
// The zero value of Queue[T] is also usable: a freshly declared
// 'var q Queue[int]' accepts Push / Pop / Peek without going
// through NewQueue. Internally, the first call lazily writes to
// q.items, so the zero value is treated as 'a Queue whose
// ArrayList has not been touched yet'.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Push appends value to the tail of the Queue. Push is amortised
// O(1) and delegates to ArrayList.Add.
func (q *Queue[T]) Push(value T) {
	q.items.Add(value)
}

// Pop removes and returns the front element wrapped in a present
// option.Optional[T], or an absent Optional if the Queue is
// empty.
//
// Pop delegates to ArrayList.RemoveFirst, which under the
// head-offset layout is O(1): the front slot is zeroed, head is
// bumped, and the discarded prefix is folded back to zero once
// it grows past 64 elements.
func (q *Queue[T]) Pop() option.Optional[T] {
	v, ok := q.items.RemoveFirst()
	if !ok {
		return option.Empty[T]()
	}
	return option.Of(v)
}

// Peek returns the front element wrapped in a present
// option.Optional[T], or an absent Optional if the Queue is
// empty. Peek delegates to ArrayList.Get(0) and is O(1).
//
// Peek does not modify the Queue: the same call repeated yields
// the same value, and Size is unchanged.
func (q *Queue[T]) Peek() option.Optional[T] {
	if q.items.IsEmpty() {
		return option.Empty[T]()
	}
	v, _ := q.items.Get(0)
	return option.Of(v)
}

// Size returns the number of elements currently in the Queue.
// Size is O(1) and delegates to ArrayList.Size.
func (q *Queue[T]) Size() int {
	return q.items.Size()
}

// IsEmpty reports whether the Queue has no elements. IsEmpty is
// the O(1) companion to Size and the standard guard before
// Pop or Peek for callers that prefer it to the absent-Optional
// return.
func (q *Queue[T]) IsEmpty() bool {
	return q.items.IsEmpty()
}

// Clear removes all elements and resets the Queue to the empty
// state. Clear delegates to ArrayList.Clear, which resets the
// head offset to zero and re-allocates a new (empty) backing
// slice. After Clear, the Queue behaves exactly as one
// returned by NewQueue.
func (q *Queue[T]) Clear() {
	q.items.Clear()
}
