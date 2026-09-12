// Package concurrency provides synchronization primitives that complement
// the standard library.
package concurrency

import (
	"sync"

	"github.com/qianwj/typed/utils/option"
)

// BoundedBlockingQueue is a fixed-capacity FIFO blocking queue.
//
// Internally it is backed by a ring buffer over a single pre-allocated slice
// and a single mutex plus a single sync.Cond. The choice of array over a
// linked list trades a small amount of index bookkeeping for significantly
// better cache locality, lower GC pressure, and a smaller per-element memory
// footprint — which is exactly what you want from a queue whose capacity is
// already fixed.
//
// Capacity is rounded up to the next power of two (so that head / tail
// advancement can use a single AND with a mask instead of a modulo). The
// actual capacity of the queue — and therefore the value returned by
// [BoundedBlockingQueue.Capacity] — may be larger than the value passed to
// [NewBoundedBlockingQueue].
//
// The zero value is not usable; construct one with [NewBoundedBlockingQueue].
//
// Blocking semantics:
//
//   - [BoundedBlockingQueue.Push] blocks while the queue is full.
//   - [BoundedBlockingQueue.Take] blocks while the queue is empty.
//   - [BoundedBlockingQueue.TryPush] / [BoundedBlockingQueue.TryTake] / [BoundedBlockingQueue.DrainTo] never block.
//
// All operations are safe for concurrent use.
//
// # Cancellation
//
// There is no PushCtx / TakeCtx yet. Adding them requires replacing
// sync.Cond with channel-based signalling so the wait can participate in
// a `select` with ctx.Done(). The design is captured in
// docs/concurrency/README.md (see "Future work: PushCtx / TakeCtx").
//
// # Naming
//
// The methods follow the conventions used elsewhere in the Typed toolkit:
// [BoundedBlockingQueue.Size] matches `Stack.Size` / `Queue.Size` / `ArrayList.Size`;
// [BoundedBlockingQueue.TryTake] returns an [option.Optional] in the same style
// as `Stack.Pop` / `Queue.Pop` / `Deque.PopFront`. Blocking reads are named
// [BoundedBlockingQueue.Take] (matching `BlockingQueue.take` in Java and the
// Go channel idiom) rather than `Pop`, because `Pop` in this toolkit means
// "non-blocking take that returns an Optional".
type BoundedBlockingQueue[T any] struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []T

	// head is the index of the next element to be returned by Take.
	// tail is the index of the next slot Push will write to.
	// count is the number of elements currently in the queue.
	//
	// head and tail advance via (i + 1) & mask rather than (i + 1) % cap.
	// That requires cap to be a power of two, which is why the constructor
	// rounds the requested capacity up.
	head  int
	tail  int
	count int
	cap   int
	mask  int
}

// NewBoundedBlockingQueue returns a queue whose actual capacity is the
// smallest power of two greater than or equal to the given capacity.
// The requested capacity must be positive; passing 0 or a negative value
// panics.
//
// Because the capacity is rounded up, a queue constructed with
// `NewBoundedBlockingQueue[T](100)` will actually hold up to 128 elements
// and report `Capacity() == 128`. If you need exact capacity, pick the
// power of two yourself.
func NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T] {
	if capacity <= 0 {
		panic("concurrency: BoundedBlockingQueue capacity must be positive")
	}
	cap := nextPowerOfTwo(capacity)
	q := &BoundedBlockingQueue[T]{
		items: make([]T, cap),
		cap:   cap,
		mask:  cap - 1,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push enqueues data, blocking until the queue has space.
func (q *BoundedBlockingQueue[T]) Push(data T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.count == q.cap {
		q.cond.Wait()
	}
	q.items[q.tail] = data
	q.tail = q.next(q.tail)
	q.count++
	q.cond.Broadcast()
}

// TryPush attempts to enqueue data without blocking.
// It returns true on success, or false immediately if the queue is full.
func (q *BoundedBlockingQueue[T]) TryPush(data T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == q.cap {
		return false
	}
	q.items[q.tail] = data
	q.tail = q.next(q.tail)
	q.count++
	q.cond.Broadcast()
	return true
}

// Take dequeues and returns the next element, blocking until one is available.
func (q *BoundedBlockingQueue[T]) Take() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.count == 0 {
		q.cond.Wait()
	}
	v := q.items[q.head]
	// Clear the slot so the previous value can be collected by the GC if T
	// contains pointers. This is purely an optimisation; it does not affect
	// correctness because we never re-read this slot until it is overwritten
	// by a future Push.
	var zero T
	q.items[q.head] = zero
	q.head = q.next(q.head)
	q.count--
	q.cond.Broadcast()
	return v
}

// TryTake attempts to dequeue without blocking.
// It returns the dequeued value as an [option.Optional]; the result is empty
// when the queue is empty.
func (q *BoundedBlockingQueue[T]) TryTake() option.Optional[T] {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		return option.Empty[T]()
	}
	v := q.items[q.head]
	var zero T
	q.items[q.head] = zero
	q.head = q.next(q.head)
	q.count--
	q.cond.Broadcast()
	return option.Of(v)
}

// DrainTo removes up to len(dst) elements from the queue and writes them into
// dst in FIFO order, then returns the number of elements actually written.
// The slice dst is not grown; if it is shorter than the queue, the surplus
// stays in the queue. If the queue is empty, DrainTo returns 0 and dst is
// left untouched.
//
// DrainTo takes a single snapshot under the queue's mutex, so the dequeued
// elements are written as a contiguous, atomically-observed batch — which
// is the main reason to prefer it over a loop of TryTake calls.
func (q *BoundedBlockingQueue[T]) DrainTo(dst []T) int {
	if len(dst) == 0 {
		return 0
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	n := q.count
	if n > len(dst) {
		n = len(dst)
	}
	if n == 0 {
		return 0
	}
	// Two contiguous ranges when the queue wraps, one otherwise.
	first := q.cap - q.head
	if first > n {
		first = n
	}
	copy(dst, q.items[q.head:q.head+first])
	if n > first {
		copy(dst[first:], q.items[:n-first])
	}
	// Clear the slots we just consumed to release references for GC.
	for i := range n {
		var zero T
		q.items[(q.head+i)&q.mask] = zero
	}
	q.head = (q.head + n) & q.mask
	q.count -= n
	q.cond.Broadcast()
	return n
}

// Size returns the current number of elements in the queue.
func (q *BoundedBlockingQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count
}

// Capacity returns the actual fixed capacity of the queue, which is the
// smallest power of two greater than or equal to the capacity passed to
// [NewBoundedBlockingQueue].
func (q *BoundedBlockingQueue[T]) Capacity() int {
	return q.cap
}

// next wraps an index around the ring buffer. The caller must hold q.mu.
func (q *BoundedBlockingQueue[T]) next(i int) int {
	return (i + 1) & q.mask
}

// nextPowerOfTwo returns the smallest power of two greater than or equal to n.
// It is only defined for n >= 1; callers must validate.
func nextPowerOfTwo(n int) int {
	// n is at least 1 here (constructor rejects ≤ 0).
	// For n == 1 we want to return 1; the bit-twiddling below would return
	// (0 | 0) + 1 = 1 too, so it happens to work without a special case, but
	// we handle it explicitly for clarity.
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}
