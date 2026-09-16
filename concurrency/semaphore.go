package concurrency

import (
	"context"
	"fmt"
)

// Semaphore is a counting semaphore that limits the number of goroutines
// that can concurrently hold a "slot". It is the primitive behind
// [Group]'s WithLimit option and is exposed directly for the common
// "max N concurrent goroutines touching X" use case where the toolkit's
// existing [BoundedBlockingQueue] doesn't apply (queue backpressure is
// per-queue, not global).
//
// Semaphore uses fixed unit weights ("N slots") — there is no
// Acquire(ctx, n) overload with a weight argument. This matches the
// common case and keeps the API minimal; if you need weighted acquires
// (e.g., "this request costs 3 credits"), reach for
// golang.org/x/sync/semaphore.
//
// Use [NewSemaphore] to construct one. The zero value is not usable.
type Semaphore struct {
	// ch is the underlying buffered channel. The buffer capacity is
	// the slot count; the number of values currently in the channel
	// is the number of available slots. Acquire receives one value;
	// Release sends one. The buffer is pre-filled at construction so
	// that len(ch) starts at capacity and Available() == capacity on
	// a freshly-constructed Semaphore.
	ch chan struct{}
}

// NewSemaphore returns a [Semaphore] with n slots. n must be positive;
// passing 0 or a negative value panics. A misconfigured capacity
// should fail loudly at construction — a Semaphore that can never be
// acquired is almost always a bug, not a deliberate "always block"
// choice.
func NewSemaphore(n int) *Semaphore {
	if n <= 0 {
		panic(fmt.Sprintf("concurrency: NewSemaphore(%d); n must be positive", n))
	}
	ch := make(chan struct{}, n)
	// Pre-fill the buffer so the slot count equals n. After this
	// loop, every Receive decrements the available count and every
	// Send increments it back.
	for range n {
		ch <- struct{}{}
	}
	return &Semaphore{ch: ch}
}

// Acquire blocks until a slot is available, then takes one. Pairs with
// [Semaphore.Release]; the available-slot count is decremented by one
// on return.
func (s *Semaphore) Acquire() {
	<-s.ch
}

// AcquireWithContext is the ctx-aware variant of [Semaphore.Acquire].
// It blocks until a slot is available or ctx is canceled or its
// deadline expires; returns nil on success or ctx.Err() on ctx firing.
// On ctx firing, no slot is taken, so callers do not need to
// [Semaphore.Release] to balance.
//
// A nil ctx is treated as [context.Background].
func (s *Semaphore) AcquireWithContext(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-s.ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire attempts to take a slot without blocking. Returns true on
// success, false immediately if no slot is available.
func (s *Semaphore) TryAcquire() bool {
	select {
	case <-s.ch:
		return true
	default:
		return false
	}
}

// Release returns a previously-acquired slot to the semaphore. Each
// Release must be paired with exactly one [Semaphore.Acquire] (or
// successful [Semaphore.TryAcquire] / [Semaphore.AcquireWithContext]).
//
// If Release is called when the semaphore is already at full
// capacity — i.e., more times than Acquire — the underlying channel
// send blocks until a slot becomes available. In a single-goroutine
// program this manifests as the Go runtime's deadlock detector
// ("fatal error: all goroutines are asleep - deadlock!"); in a
// concurrent program it can be silent if other goroutines make
// progress, which is why misuse is hard to diagnose. Always pair
// Release with Acquire, and prefer [Semaphore.TryAcquire] when the
// pairing is conditional.
//
// Unlike golang.org/x/sync/semaphore, over-release is NOT supported:
// the slot count is bounded by the initial [NewSemaphore] capacity
// and cannot grow beyond it. This is intentional — the toolkit's
// policy is "fail at the misuse site", not "silently grow".
func (s *Semaphore) Release() {
	s.ch <- struct{}{}
}

// Available returns the current number of free slots. It is an atomic
// length read on the underlying channel, not serialised with the next
// operation you perform. The window is small (a single atomic load)
// but real; if you need strict "Available() → next op sees a
// consistent view" semantics, use [Semaphore.TryAcquire], which
// combines the read with the take.
func (s *Semaphore) Available() int {
	return len(s.ch)
}
