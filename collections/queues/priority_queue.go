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
	"fmt"

	"github.com/qianwj/typed/adt/option"
)

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
// # Capacity
//
// The queue is bounded at construction time by the capacity
// argument to [NewPriorityQueue]:
//
//   - capacity > 0: the queue holds at most that many elements.
//     Push follows the standard top-K rule — see Push. This is
//     the common case (bounded caches, top-N tracking, etc.).
//   - capacity == 0: the queue is unbounded. Push always accepts.
//     Use this when you want a heap that grows on demand.
//   - capacity < 0: NewPriorityQueue panics — a misconfigured
//     capacity should fail loudly at construction.
//
// # API shape
//
// Push / Pop / Peek / Size mirror the [Queue] surface area
// one-for-one. The roadmap sketched Push / Poll / TryPoll (the
// concurrency.BoundedBlockingQueue verbs), but in a synchronous
// container Poll has no natural meaning distinct from TryPoll
// — it cannot block. To keep the API consistent with [Queue] we
// drop Poll and use Pop instead.
type PriorityQueue[T any] struct {
	// heap is a binary heap ordered by `less`. The heap invariant
	// is maintained locally after every Push / Pop — see heapifyUp
	// and heapifyDown. capacity grows on demand via append; there
	// is no separate capacity / size split.
	heap []T

	// less orders the queue: less(a, b) reports that a should be
	// dequeued before b. Stored at construction; never reassigned.
	less func(a, b T) bool

	// capacity is the upper bound on len(heap). 0 means unbounded;
	// >0 means bounded. Stored at construction; never reassigned.
	capacity int
}

// NewPriorityQueue returns an empty PriorityQueue ordered by less:
// Pop / Peek return the element for which less(element, other) is
// true for every other element in the queue.
//
// capacity sets the upper bound on the queue's size:
//
//   - capacity > 0 — the queue holds at most capacity elements.
//     Push follows the top-K rule: the new element is accepted
//     when the queue is below capacity, or when the queue is at
//     capacity and `less(new, root)` is true (in which case the
//     root is replaced by the new element). Push returns false
//     when the new element is dropped. See [PriorityQueue.Push]
//     for the full contract.
//   - capacity == 0 — the queue is unbounded. Push always
//     accepts. Use this for the "grow on demand" case.
//   - capacity < 0 — panics. A misconfigured capacity should
//     fail loudly at construction, not silently produce a queue
//     that drops everything or accepts everything.
//
// A classic min-heap of [cmp.Ordered] types is
//
//	NewPriorityQueue(0, cmp.Less[int])            // ints, unbounded
//	NewPriorityQueue(100, cmp.Less[int])          // top-100 ints
//	NewPriorityQueue(0, func(a, b string) bool {  // strings, unbounded
//	    return len(a) < len(b)                    // shortest first
//	})
//
// To get FIFO (or any stable order) among ties, encode the
// tiebreaker in a wrapper and pass a composite-key comparator —
// see the PriorityQueue type doc for an example.
//
// less must be a pure function of its arguments (no observable
// side effects, deterministic output for equal inputs). An
// inconsistent comparator will produce an inconsistent heap —
// elements may come out in the wrong order, or never come out at
// all. The standard library's [sort.Slice] / [container/heap]
// docs carry the same warning.
func NewPriorityQueue[T any](capacity int, less func(a, b T) bool) *PriorityQueue[T] {
	if capacity < 0 {
		panic(fmt.Sprintf("queues: NewPriorityQueue(capacity=%d, ...); capacity must be non-negative", capacity))
	}
	return &PriorityQueue[T]{
		less:     less,
		capacity: capacity,
	}
}

// Push inserts data into the queue and reports whether the queue
// accepted it.
//
//   - In an unbounded queue (capacity == 0), Push always appends
//     the new element and returns true.
//   - In a bounded queue (capacity > 0), Push returns true when
//     the queue is below capacity (the new element is appended)
//     or when the queue is at capacity and `less(data, boundary)`
//     is true, where `boundary` is the current lowest-priority
//     element in the heap (the K-th highest, the eviction
//     candidate). In the replacement case the boundary slot is
//     overwritten with the new element and the heap is sifted
//     up to maintain the min-heap invariant. Push returns false
//     when the queue is at capacity and the new element is not
//     higher priority than the boundary — the new element is
//     dropped in that case.
//
// The bounded case is the standard top-K rule: at any point the
// queue holds the capacity-many highest-priority elements that
// have been Pushed so far, regardless of insertion order. See
// the type doc for the rationale and typical use cases.
//
// Why scan for the boundary instead of comparing against the
// root? The heap is a min-heap by `less`, so the root is the
// HIGHEST-priority element — replacing the root with a new
// higher-priority one would discard the current best, not the
// current worst. The K-th highest (the eviction candidate) sits
// somewhere in the leaves; we scan O(K) to find it. The total
// Push cost in bounded mode is O(K) — acceptable for the
// typical small K (top-N queries, bounded caches).
//
// Push is O(log n) in the unbounded case. In the bounded case
// Push is O(K + log K) (linear scan for the boundary plus sift
// up), and O(log K) when below capacity (append + sift up).
func (q *PriorityQueue[T]) Push(data T) bool {
	if q.capacity > 0 && len(q.heap) >= q.capacity {
		// At capacity. Find the boundary — the element with the
		// LOWEST priority, i.e. the MAX under `less` — by linear
		// scan. This is O(K); see the doc above.
		boundary := 0
		for i := 1; i < len(q.heap); i++ {
			// heap[i] has lower priority than heap[boundary] iff
			// !less(heap[i], heap[boundary]) (i.e., heap[i] is not
			// higher priority than heap[boundary]). The boundary
			// tracks the worst-in-heap (the max by less), so we
			// update whenever heap[i] is worse than (or tied with)
			// the current boundary.
			if !q.less(q.heap[i], q.heap[boundary]) {
				boundary = i
			}
		}

		// If the new element is higher priority than the boundary,
		// replace and sift up to maintain the min-heap invariant.
		// Sift up only — the boundary sits at (or near) a leaf, so
		// no sift-down is needed for the replacement itself, only
		// propagation up if the new value out-prioritises its parent.
		if q.less(data, q.heap[boundary]) {
			q.heap[boundary] = data
			heapifyUp(q.heap, boundary, q.less)
			return true
		}
		return false
	}
	q.heap = append(q.heap, data)
	heapifyUp(q.heap, len(q.heap)-1, q.less)
	return true
}

// Pop removes and returns the highest-priority element, or an
// absent [Option] if the queue is empty.
//
// Pop is O(log n), dominated by the sift-down.
func (q *PriorityQueue[T]) Pop() option.Option[T] {
	if len(q.heap) == 0 {
		return option.Empty[T]()
	}
	top := q.heap[0]
	n := len(q.heap) - 1
	q.heap[0] = q.heap[n]
	// Zero the freed slot so a pointer-typed T is not pinned in
	// the backing array after removal — mirrors [collections.Stack]
	// / [Queue].
	var zero T
	q.heap[n] = zero
	q.heap = q.heap[:n]
	if n > 0 {
		heapifyDown(q.heap, 0, n, q.less)
	}
	return option.Of(top)
}

// Peek returns the highest-priority element without removing it,
// or an absent [Option] if the queue is empty. Peek is O(1).
func (q *PriorityQueue[T]) Peek() option.Option[T] {
	if len(q.heap) == 0 {
		return option.Empty[T]()
	}
	return option.Of(q.heap[0])
}

// Size returns the current number of elements in the queue.
// Named to match the [collections.Stack] / [Queue] / [Deque] /
// concurrency.BoundedBlockingQueue convention (all the queue
// types in the toolkit, including the ones in the concurrency
// package); the underlying stdlib [container/heap] uses Len,
// which the toolkit deliberately deviates from for naming
// consistency across all queue-like types.
func (q *PriorityQueue[T]) Size() int {
	return len(q.heap)
}

// Capacity returns the upper bound configured at [NewPriorityQueue]
// time. Returns 0 if the queue is unbounded — the "0 = unbounded"
// reading follows [container/list]'s idiom (an "n element" of 0
// means unlimited) and is documented in the constructor.
func (q *PriorityQueue[T]) Capacity() int {
	return q.capacity
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
