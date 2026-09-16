// Package queues provides typed container helpers that complement the
// top-level collections package: heap-backed [PriorityQueue] for now.
//
// All types in this package are synchronous — no internal locking —
// matching the rest of the collections package. For cross-goroutine
// use, layer a channel, a [sync.Mutex], or a concurrency.Group on
// top, or reach for the typed concurrent primitives in the
// concurrency package.
package queues

import (
	"github.com/qianwj/typed/adt"
)

// PriorityQueue[T] is a min-heap-backed priority queue. Elements
// with the smallest priority value are dequeued first; the
// relative order of equal-priority elements is unspecified —
// a strict-less binary heap has no natural tiebreaker. Callers
// that need a stable order among ties should encode it in the
// priority itself (e.g., use a monotonic counter as the low bits).
//
// PriorityQueue is synchronous (no internal locking) and lives in
// the collections package rather than concurrency: the use cases
// the roadmap lists — "schedulers, leader election, anything with
// 'process the most important thing first'" — are typically
// intra-goroutine ordering operations. For cross-goroutine use,
// wrap with a Mutex or hand the data through a channel.
//
// The zero value of PriorityQueue is not usable; construct with
// [NewPriorityQueue].
//
// # API shape
//
// Push / PushWithPriority / Pop / Peek / Len mirror the
// collections.Queue surface area one-for-one — the only addition
// is PushWithPriority. The roadmap sketched Push / Poll /
// TryPoll (the BoundedBlockingQueue verbs), but in a synchronous
// container Poll has no natural meaning distinct from TryPoll /
// Pop: it cannot block. To keep the API consistent with
// collections.Queue we drop Poll and use Pop, which is the
// toolkit's established "remove-and-return" verb for sync
// collections.
type PriorityQueue[T any] struct {
	// heap is a standard min-heap ordered by .priority. The heap
	// invariant is maintained locally after every Push / Pop — see
	// heapifyUp and heapifyDown. capacity grows on demand via
	// append; there is no separate capacity / size split.
	heap []priorityItem[T]
}

// priorityItem pairs a stored value with its priority. Stored
// contiguously in the heap slice so a single append-and-sift
// keeps both fields coherent.
type priorityItem[T any] struct {
	priority int
	data     T
}

// NewPriorityQueue returns an empty PriorityQueue. The zero value
// is also usable via the heap-package's own lazy-init pattern, but
// going through the constructor is the documented path.
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{}
}

// Push adds data to the queue with priority 0 (mid-range). Use
// [PriorityQueue.PushWithPriority] when the priority matters.
func (q *PriorityQueue[T]) Push(data T) {
	q.PushWithPriority(data, 0)
}

// PushWithPriority adds data with the given priority. Lower
// priority values are dequeued first. Equal-priority elements
// have no guaranteed relative order — see the type doc.
//
// PushWithPriority is O(log n), dominated by the sift-up.
func (q *PriorityQueue[T]) PushWithPriority(data T, priority int) {
	q.heap = append(q.heap, priorityItem[T]{priority: priority, data: data})
	heapifyUp(q.heap, len(q.heap)-1)
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
	var zero priorityItem[T]
	q.heap[n] = zero
	q.heap = q.heap[:n]
	if n > 0 {
		heapifyDown(q.heap, 0, n)
	}
	return adt.Of(top.data)
}

// Peek returns the highest-priority element without removing it,
// or an absent [Option] if the queue is empty. Peek is O(1).
func (q *PriorityQueue[T]) Peek() adt.Option[T] {
	if len(q.heap) == 0 {
		return adt.Empty[T]()
	}
	return adt.Of(q.heap[0].data)
}

// Len returns the current number of elements in the queue.
func (q *PriorityQueue[T]) Len() int {
	return len(q.heap)
}

// heapifyUp restores the heap invariant by sifting the element at
// index i up toward the root. Called after append — the new
// element can only violate the invariant on the path from i to
// the root.
func heapifyUp[T any](heap []priorityItem[T], i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if heap[i].priority < heap[parent].priority {
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
func heapifyDown[T any](heap []priorityItem[T], root, n int) {
	for {
		smallest := root
		left := 2*root + 1
		right := 2*root + 2
		if left < n && heap[left].priority < heap[smallest].priority {
			smallest = left
		}
		if right < n && heap[right].priority < heap[smallest].priority {
			smallest = right
		}
		if smallest == root {
			return
		}
		heap[root], heap[smallest] = heap[smallest], heap[root]
		root = smallest
	}
}
