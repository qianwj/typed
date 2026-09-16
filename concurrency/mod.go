// Package concurrency provides synchronization primitives that complement
// the standard library: typed generic wrappers around `chan T`,
// `sync.WaitGroup`, and the runtime signals they are built on.
//
// # Queues
//
// BoundedBlockingQueue[T] and UnboundedBlockingQueue[T] are FIFO blocking
// queues written as typed generic wrappers around a single `chan T`. The
// bounded form rounds the requested capacity up to a power of two so
// Capacity() always returns a value usable as a bitmask; it has
// context-aware blocking for cancellation, deadlines, and shutdown. The
// unbounded form uses a ring buffer + mutex + cond because the runtime
// has no "unbounded buffered channel". Both return option.Optional[T]
// from their non-blocking probes to match the project-wide Optional
// convention — see Stack.Pop / Queue.Front / Deque.PopFront for the same
// pattern in collections.
//
// # Structured concurrency
//
// Group is a typed wrapper around "spawn N goroutines, wait for them all"
// with parent-context propagation, per-task panic recovery (best-effort
// by default; strict mode fails the whole group on the first panic), and
// a Go / TryGo / SetLimit pair for admission control. Use it instead of
// raw goroutines when the lifetime of the spawned work should be tied to
// a single enclosing scope, the way errgroup-style structured concurrency
// is meant to work.
//
// # Counting semaphore
//
// Semaphore is a typed counting semaphore: Acquire blocks the caller
// until a slot is free, or the context cancels. Each unit weight is one
// slot; the semaphore is fair under FIFO admission for the blocked
// callers waiting on the same capacity.
//
// # Conventions
//
// None of the types in this package are safe for concurrent mutation of
// the *same* value from multiple goroutines; each primitive owns its own
// internal synchronisation and exposes it through its API surface only.
// Concrete generic types (not interfaces) for the same reason as the
// rest of the toolkit: Go 1.27's generic methods require a concrete
// receiver.
package concurrency
