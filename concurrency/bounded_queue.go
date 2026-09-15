// Package concurrency provides synchronization primitives that complement
// the standard library.
package concurrency

import (
	"context"

	"github.com/qianwj/typed/adt"
)

// BoundedBlockingQueue is a fixed-capacity FIFO blocking queue.
//
// It is a thin generic wrapper around a single bounded `chan T`, adding
// three things on top of the channel primitives:
//
//   - A [BoundedBlockingQueue.TryPush] / [BoundedBlockingQueue.TryPoll] pair
//     that returns an [adt.Option] (matching the toolkit convention
//     from `Stack.Pop` / `Queue.Pop` / `Deque.PopFront`).
//   - Context-aware [BoundedBlockingQueue.PushWithContext] /
//     [BoundedBlockingQueue.PollWithContext] for cancellation, deadlines,
//     and shutdown signalling.
//   - A [BoundedBlockingQueue.Capacity] accessor and a power-of-two rounding
//     policy on the requested capacity, so `Capacity()` always returns a
//     value usable as a bitmask.
//
// Internally the queue is just `make(chan T, cap)`. There is no mutex, no
// ring buffer, no custom condition variable — the Go runtime's per-P
// scheduler-aware channel implementation handles all of the hard parts,
// which is why this wrapper is in the same throughput ballpark as a raw
// `chan T` instead of the 5–10× slower a hand-rolled ring buffer would be.
//
// Capacity is rounded up to the next power of two. The actual capacity of
// the queue — and therefore the value returned by
// [BoundedBlockingQueue.Capacity] — may be larger than the value passed to
// [NewBoundedBlockingQueue].
//
// The zero value is not usable; construct one with [NewBoundedBlockingQueue].
//
// # Naming
//
// The four blocking operations follow a consistent pattern:
//
//   - [BoundedBlockingQueue.Push] / [BoundedBlockingQueue.Poll] — the
//     blocking verbs. `Push` blocks while the queue is full; `Poll` blocks
//     while the queue is empty.
//   - [BoundedBlockingQueue.PushWithContext] / [BoundedBlockingQueue.PollWithContext]
//     — the context-aware variants. They block the same way as Push / Poll
//     and additionally participate in `select { <-ctx.Done(): }` so
//     cancellation, deadlines, and shutdown signals propagate cleanly.
//   - [BoundedBlockingQueue.TryPush] / [BoundedBlockingQueue.TryPoll] —
//     the non-blocking variants. They never wait and return a `bool`
//     (for Push) or an [adt.Option] (for Poll) so the caller can
//     branch on backpressure or absence.
//
// Blocking semantics:
//
//   - [BoundedBlockingQueue.Push] blocks while the queue is full.
//   - [BoundedBlockingQueue.Poll] blocks while the queue is empty.
//   - [BoundedBlockingQueue.PushWithContext] /
//     [BoundedBlockingQueue.PollWithContext] block until the queue is in the
//     corresponding state, or until ctx is canceled — whichever happens first.
//   - [BoundedBlockingQueue.TryPush] / [BoundedBlockingQueue.TryPoll] never block.
//
// All operations are safe for concurrent use.
//
// # What a `chan T` does and does not give you
//
// Because the queue is a channel under the hood, the trade-offs versus a
// hand-rolled ring buffer + mutex + Cond are inherited from `chan T`:
//
//   - The channel's internal ring buffer does not zero freed slots, so
//     pointer values received from the queue stay alive in the channel's
//     backing array until the slot is overwritten by a new send. A custom
//     ring buffer can zero slots on `Poll` to help the GC; this wrapper
//     cannot.
//   - [BoundedBlockingQueue.Size] is `len(ch)`, an atomic length read, but
//     it is not serialised with the next operation you perform. A
//     ring-buffer `Size()` under a mutex gives the stronger "Size() → next
//     op sees a consistent view" property; `len(ch)` does not. In practice
//     the window is a single atomic read, so this rarely matters, but it
//     is a real semantic difference.
type BoundedBlockingQueue[T any] struct {
	ch  chan T
	cap int
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
	n := nextPowerOfTwo(capacity)
	return &BoundedBlockingQueue[T]{
		ch:  make(chan T, n),
		cap: n,
	}
}

// Push enqueues data, blocking until the queue has space.
func (q *BoundedBlockingQueue[T]) Push(data T) {
	q.ch <- data
}

// PushWithContext is the context-aware variant of [BoundedBlockingQueue.Push].
// It blocks until the queue has space, or until ctx is canceled — whichever
// happens first. It returns nil on success and ctx.Err() on cancellation;
// the element is not enqueued in the cancellation case.
func (q *BoundedBlockingQueue[T]) PushWithContext(ctx context.Context, data T) error {
	select {
	case q.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryPush attempts to enqueue data without blocking.
// It returns true on success, or false immediately if the queue is full.
func (q *BoundedBlockingQueue[T]) TryPush(data T) bool {
	select {
	case q.ch <- data:
		return true
	default:
		return false
	}
}

// Poll dequeues and returns the next element, blocking until one is available.
func (q *BoundedBlockingQueue[T]) Poll() T {
	return <-q.ch
}

// PollWithContext is the context-aware variant of [BoundedBlockingQueue.Poll].
// It blocks until an element is available, or until ctx is canceled —
// whichever happens first. It returns (value, nil) on success and
// (zero, ctx.Err()) on cancellation.
func (q *BoundedBlockingQueue[T]) PollWithContext(ctx context.Context) (T, error) {
	select {
	case v := <-q.ch:
		return v, nil
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// TryPoll attempts to dequeue without blocking.
// It returns the dequeued value as an [adt.Option]; the result is empty
// when the queue is empty.
func (q *BoundedBlockingQueue[T]) TryPoll() adt.Option[T] {
	select {
	case v := <-q.ch:
		return adt.Of(v)
	default:
		return adt.Empty[T]()
	}
}

// Size returns the current number of elements in the queue.
//
// This is `len(ch)`, an atomic length read; it is not serialised with the
// next operation you perform. The window is small (a single atomic load)
// but real — if you need strict "Size() → next op sees a consistent view"
// semantics, use a ring-buffer implementation instead.
func (q *BoundedBlockingQueue[T]) Size() int {
	return len(q.ch)
}

// Capacity returns the actual fixed capacity of the queue, which is the
// smallest power of two greater than or equal to the capacity passed to
// [NewBoundedBlockingQueue].
func (q *BoundedBlockingQueue[T]) Capacity() int {
	return q.cap
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