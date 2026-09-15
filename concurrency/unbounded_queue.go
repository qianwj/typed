// Package concurrency provides synchronization primitives that complement
// the standard library.
package concurrency

import (
	"context"
	"sync"

	"github.com/qianwj/typed/utils/option"
)

// UnboundedBlockingQueue is an unbounded FIFO blocking queue.
//
// Push never blocks — there is no capacity to wait on. Take blocks when
// the queue is empty.
//
// Internally the queue is backed by a ring buffer over a single
// pre-allocated slice plus a single mutex plus a `*sync.Cond`. The ring
// buffer grows in powers of two when full, so amortised per-op cost is
// O(1). Unlike [BoundedBlockingQueue], which is a thin wrapper around
// `chan T`, this type cannot use a channel because the Go runtime has
// no "unbounded buffered channel". A ring buffer with a mutex and a
// condition variable is the standard fix.
//
// # Signalling
//
// Push signals the cond with `Signal()` (one waiter at a time), so a
// burst of N pushes wakes up to N blocked takers — no thundering herd.
// `Broadcast()` would also be correct but would wake every blocked
// taker on every push, which is wasteful when the worker count is
// large.
//
// `TakeCtx` additionally spawns a one-shot watcher goroutine that calls
// `Broadcast()` on `ctx.Done()` so a canceled context can wake blocked
// takers. The goroutine exits when TakeCtx returns.
//
// # Naming
//
// The methods follow the conventions used elsewhere in the Typed toolkit:
// [UnboundedBlockingQueue.Size] matches `Stack.Size` / `Queue.Size`;
// [UnboundedBlockingQueue.TryTake] returns an [option.Optional] in the
// same style as `Stack.Pop` / `Queue.Pop`. Blocking reads are named
// [UnboundedBlockingQueue.Take] to match the [BoundedBlockingQueue]
// naming. Unlike the bounded sibling, this type has no `Capacity` — the
// queue is unbounded by definition.
//
// The zero value is not usable; construct one with
// [NewUnboundedBlockingQueue].
type UnboundedBlockingQueue[T any] struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []T

	// head is the index of the next element to be returned by Take.
	// tail is the index of the next slot Push will write to.
	// count is the number of elements currently in the queue.
	// mask is items.length - 1, kept as a precomputed bitmask so that
	// head and tail advance via `(i + 1) & mask` instead of `% len`.
	//
	// items grows in powers of two (see grow) so mask is always valid
	// after the first grow.
	head  int
	tail  int
	count int
	mask  int
}

// NewUnboundedBlockingQueue returns a fresh, empty queue. There is no
// capacity argument — the queue grows on demand.
func NewUnboundedBlockingQueue[T any]() *UnboundedBlockingQueue[T] {
	q := &UnboundedBlockingQueue[T]{
		items: make([]T, 16),
		mask:  15,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// grow doubles the backing array and copies the items in order. Caller
// must hold q.mu.
func (q *UnboundedBlockingQueue[T]) grow() {
	oldCap := q.mask + 1
	newCap := max(oldCap*2, 16)
	newItems := make([]T, newCap)
	// Copy items[head:oldCap], then wrap around items[0:tail] if needed.
	n := copy(newItems, q.items[q.head:])
	if n < q.count {
		copy(newItems[n:], q.items[:q.count-n])
	}
	q.items = newItems
	q.head = 0
	q.tail = q.count // items[0:count] are valid (copy preserved order)
	q.mask = newCap - 1
}

// Push enqueues data. It never blocks because the queue is unbounded.
func (q *UnboundedBlockingQueue[T]) Push(data T) {
	q.mu.Lock()
	if q.count == q.mask+1 {
		q.grow()
	}
	q.items[q.tail] = data
	q.tail = (q.tail + 1) & q.mask
	q.count++
	q.cond.Signal() // wake exactly one waiting taker
	q.mu.Unlock()
}

// TryPush is the non-blocking variant of [UnboundedBlockingQueue.Push].
// For an unbounded queue it always succeeds and returns true.
func (q *UnboundedBlockingQueue[T]) TryPush(data T) bool {
	q.Push(data)
	return true
}

// Take dequeues and returns the next element, blocking until one is available.
func (q *UnboundedBlockingQueue[T]) Take() T {
	q.mu.Lock()
	for q.count == 0 {
		q.cond.Wait()
	}
	v := q.items[q.head]
	// Clear the slot so the previous value can be collected by the GC if
	// T contains pointers. This is purely an optimisation; it does not
	// affect correctness because we never re-read this slot until it is
	// overwritten by a future Push.
	var zero T
	q.items[q.head] = zero
	q.head = (q.head + 1) & q.mask
	q.count--
	q.mu.Unlock()
	return v
}

// TryTake is the non-blocking variant of [UnboundedBlockingQueue.Take].
// It returns the dequeued value as an [option.Optional]; the result is
// empty when the queue is empty.
func (q *UnboundedBlockingQueue[T]) TryTake() option.Optional[T] {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		return option.Empty[T]()
	}
	v := q.items[q.head]
	var zero T
	q.items[q.head] = zero
	q.head = (q.head + 1) & q.mask
	q.count--
	return option.Of(v)
}

// PushCtx is the context-aware variant of [UnboundedBlockingQueue.Push].
// Because Push never blocks, PushCtx returns ctx.Err() without enqueuing
// when ctx is already canceled, and otherwise behaves like Push.
func (q *UnboundedBlockingQueue[T]) PushCtx(ctx context.Context, data T) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	q.Push(data)
	return nil
}

// TakeCtx is the context-aware variant of [UnboundedBlockingQueue.Take].
// It blocks until an element is available, or until ctx is canceled —
// whichever happens first. Returns (zero, ctx.Err()) on cancellation.
func (q *UnboundedBlockingQueue[T]) TakeCtx(ctx context.Context) (T, error) {
	var zero T

	// Watch ctx in a one-shot goroutine. When ctx is canceled, broadcast
	// the cond so any blocked taker wakes up and checks ctx.Err(). The
	// `done` channel makes sure the goroutine exits when TakeCtx
	// returns, even if ctx is never canceled.
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			q.mu.Lock()
			q.cond.Broadcast()
			q.mu.Unlock()
		case <-done:
		}
	}()

	q.mu.Lock()
	defer q.mu.Unlock()
	for q.count == 0 {
		if err := ctx.Err(); err != nil {
			return zero, err
		}
		q.cond.Wait()
	}
	v := q.items[q.head]
	var zeroInner T
	q.items[q.head] = zeroInner
	q.head = (q.head + 1) & q.mask
	q.count--
	return v, nil
}

// Size returns the current number of elements in the queue.
func (q *UnboundedBlockingQueue[T]) Size() int {
	q.mu.Lock()
	n := q.count
	q.mu.Unlock()
	return n
}
