// Package queues provides typed container helpers that complement the
// top-level collections package: heap-backed [PriorityQueue] for now.
//
// All types in this package are synchronous — no internal locking —
// matching the rest of the collections package. For cross-goroutine
// use, layer a channel, a [sync.Mutex], or a concurrency.Group on
// top, or reach for the typed concurrent primitives in the
// concurrency package.
package queues

import "github.com/qianwj/typed/adt"

// PriorityQueue[T] is a binary-heap-backed priority queue. The
// ordering is supplied by the caller as a comparator function, the
// same way [sort.Slice] / [container/heap] do it: Pop / Peek return
// the element for which `less(element, other)` is true for every
// other element.
//
// Why a comparator rather than `PushWithPriority(data, priority)`?
//
// A comparator is strictly more general: integer priority with
// "lower = first" is just one possible `less` (specifically
// `func(a, b T) bool { return a < b }`, or `cmp.Less[T]` for any
// `cmp.Ordered`). A comparator also lets callers encode multi-field
// keys ("earlier deadline wins, then lower id"), domain-specific
// orders ("shortest job first"), and stable FIFO among ties
// (encode a monotonic counter as a tiebreaker) without the type
// having to know any of those rules. The int-priority shortcut
// would have baked in one particular scheme at the API level and
// forced every other scheme through it.
//
// PriorityQueue is synchronous (no internal locking) and lives in
// the collections package rather than concurrency: the roadmap
// use cases — schedulers, leader election, "process the most
// important thing first" — are typically intra-goroutine ordering.
// For cross-goroutine use, wrap with a Mutex or hand the data
// through a channel.
//
// The zero value of PriorityQueue is not usable; construct with
// [NewPriorityQueue].
//
// # API shape
//
// Push / Pop / Peek / Len mirror the collections.Queue surface
// area one-for-one. The roadmap sketched Push / Poll / TryPoll
// (the BoundedBlockingQueue verbs), but in a synchronous container
// Poll has no natural meaning distinct from TryPoll — it cannot
// block. To keep the API consistent with collections.Queue we drop
// Poll and use Pop instead.
type PriorityQueue[T any] struct {
	// heap is a binary heap ordered by `less`. The heap invariant
	// is maintained locally after every Push / Pop — see heapifyUp
	// and heapifyDown. capacity grows on demand via append; there
	// is no separate capacity / size split.
	heap []T

	// less orders the queue: less(a, b) reports that a should be
	// dequeued before b. Stored at construction; never reassigned.
	less func(a, b T) bool
}

// NewPriorityQueue returns an empty PriorityQueue ordered by less:
// Pop / Peek return the element for which less(element, other) is
// true for every other element in the queue.
//
// A classic min-heap for [cmp.Ordered] types is
//
//	NewPriorityQueue(cmp.Less[int])            // ints, smallest first
//	NewPriorityQueue(func(a, b string) bool {  // strings, ...
//	    return len(a) < len(b)                  // shortest first
//	})
//
// To get FIFO among ties (or any other secondary key), encode it
// in a wrapper:
//
//	type byDeadline struct{ deadline time.Time; seq int }
//	NewPriorityQueue(func(a, b byDeadline) bool {
//	    if !a.deadline.Equal(b.deadline) {
//	        return a.deadline.Before(b.deadline)
//	    }
//	    return a.seq < b.seq // monotonic counter, FIFO tiebreak
//	})
//
// less must be a pure function of its arguments (no observable
// side effects, deterministic output for equal inputs). An
// inconsistent comparator will produce an inconsistent heap —
// elements may come out in the wrong order, or never come out at
// all. The standard library's [sort.Slice] / [container/heap]
// docs carry the same warning.
func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{less: less}
}

// Push adds data to the queue.
//
// Push is O(log n), dominated by the sift-up.
func (q *PriorityQueue[T]) Push(data T) {
	q.heap = append(q.heap, data)
	heapifyUp(q.heap, len(q.heap)-1, q.less)
}

// Pop removes and returns the highest-priority element, or an
// absent [Option] if the queue is empty.
//
// Pop is O(log n), dominated by the sift-down.
func (q *PriorityQueue[T]) Pop() adt.Option[T] {
	if len(q.heap) == 0 {
		return adt.Empty[T]()
	}
	top := q.heap[0]
	n := len(q.heap) - 1
	q.heap[0] = q.heap[n]
	// Zero the freed slot so a pointer-typed T is not pinned in
	// the backing array after removal — mirrors collections.Stack
	// / collections.Queue.
	var zero T
	q.heap[n] = zero
	q.heap = q.heap[:n]
	if n > 0 {
		heapifyDown(q.heap, 0, n, q.less)
	}
	return adt.Of(top)
}

// Peek returns the highest-priority element without removing it,
// or an absent [Option] if the queue is empty. Peek is O(1).
func (q *PriorityQueue[T]) Peek() adt.Option[T] {
	if len(q.heap) == 0 {
		return adt.Empty[T]()
	}
	return adt.Of(q.heap[0])
}

// Len returns the current number of elements in the queue.
func (q *PriorityQueue[T]) Len() int {
	return len(q.heap)
}

// heapifyUp restores the heap invariant by sifting the element at
// index i up toward the root. Called after append — the new
// element can only violate the invariant on the path from i to
// the root.
func heapifyUp[T any](heap []T, i int, less func(a, b T) bool) {
	for i > 0 {
		parent := (i - 1) / 2
		if less(heap[i], heap[parent]) {
			heap[i], heap[parent] = heap[parent], heap[i]
			i = parent
			continue
		}
		return
	}
}

// heapifyDown restores the heap invariant by sifting the element
// at index root down toward the leaves. n is the heap size — the
// active range is heap[:n]. Called after the root has been
// replaced by the last element and the last slot zeroed.
func heapifyDown[T any](heap []T, root, n int, less func(a, b T) bool) {
	for {
		smallest := root
		left := 2*root + 1
		right := 2*root + 2
		if left < n && less(heap[left], heap[smallest]) {
			smallest = left
		}
		if right < n && less(heap[right], heap[smallest]) {
			smallest = right
		}
		if smallest == root {
			return
		}
		heap[root], heap[smallest] = heap[smallest], heap[root]
		root = smallest
	}
}
