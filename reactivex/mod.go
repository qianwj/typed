// Package reactivex provides typed asynchronous event streams with explicit
// subscriptions, consumer demand and configurable overflow behavior.
//
// # Observable and Publisher
//
// Observable[T] is the concrete pipeline type. Its operators construct new
// observables without starting a producer; Subscribe, ForEach and ToSlice
// start consumption. Publisher[T] is the smaller interface for components
// that only need to subscribe. Subject[T] implements both Publisher[T] and
// Subscriber[T], providing a hot multicast entry point.
//
// Observable is concrete so transformations such as Map[R] can declare their
// own result type under Go 1.27's generic methods. These method-level type
// parameters do not belong on the Publisher interface.
//
// A finite pipeline can be collected directly:
//
//	values, err := reactivex.Just(1, 2, 3, 4).
//	    Filter(func(v int) bool { return v%2 == 0 }).
//	    ToSlice(ctx)
//	// On successful completion, values is []int{2, 4}.
//
// # Demand and buffering
//
// Subscribe gives the Subscriber a Subscription through OnSubscribe. The
// subscriber must call Request(n) to allow values to flow. ForEach and ToSlice
// request the maximum uint64 demand for callers that want all available values.
//
// Demand and buffer capacity serve different purposes. Request(n) permits
// delivery; WithBuffer limits pending values at a Subject or channel boundary.
// WithOverflow chooses whether a full boundary blocks, discards values or
// terminates the affected subscription. Ordinary Map and Filter operators
// wrap subscribers without adding their own queue or goroutine.
//
//	subject := reactivex.NewSubject[int](
//	    reactivex.WithBuffer(128),
//	    reactivex.WithOverflow(reactivex.OverflowDropOldest),
//	)
//
// # Ownership and lifecycle
//
// FromSlice iterates the supplied slice separately for each subscription but
// does not copy its backing array. FromSeq invokes the supplied sequence again;
// whether it can be reused depends on that sequence. FromChannel reads from
// the original channel, so multiple subscriptions compete for values rather
// than receiving copies. Use Subject when every active subscriber should have
// its own opportunity to receive each value.
//
// Cancel signals that consumption should stop. Done reports cancellation or
// termination, not that every user callback or producer goroutine has exited.
// Blocking callbacks and custom producers must cooperate with cancellation;
// the package cannot interrupt arbitrary user code or close caller-owned input
// channels. Keep producer resources and callback state under explicit ownership.
//
// This package is independent of typed/collections. Use its synchronous stream
// package for in-memory iteration, and reactivex when subscription lifetimes,
// asynchronous inputs or multicast behavior are part of the problem.
package reactivex
