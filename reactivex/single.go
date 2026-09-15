package reactivex

import (
	"context"
	"sync"
	"sync/atomic"
)

// Single[T] is a reactive container that emits exactly one value or
// one error. It is the typed equivalent of Go's `func() (T, error)`
// combined with a Future-like subscription model.
//
// Compared to Observable[T]:
//
//   - Cardinality is fixed at 1: a Single terminates with either
//     OnSuccess(T) or OnError(error). There is no OnComplete-without-
//     value state, so callers never have to distinguish "no value
//     arrived yet" from "no value will ever arrive".
//   - There is no demand tracking because at most one value will be
//     delivered; the Subscriber does not call Request.
//   - There is no per-subscription replay: a Single runs its source
//     function once and caches the result, so every subscriber sees
//     the same outcome. Use Observable if you want a fresh execution
//     per subscriber.
//
// Construct with [NewSingle]; consume either reactively via [Subscribe]
// or synchronously via [Await] / [AwaitWithContext].
//
// # Naming
//
//   - [Subscribe] — register callbacks for the terminal event.
//   - [Await] / [AwaitWithContext] — block until the terminal event
//     and return the result.
//   - [Done] — non-blocking check that the terminal event has fired.
//   - [Map] / [FlatMap] / [Zip] / [AndThen] — compose multiple
//     Singles into a new Single.
//
// Single values must not be copied after construction.
type Single[T any] struct {
	fn func() (T, error)

	// subs holds callbacks for every active subscriber. Protected
	// by subsMu. On terminal event each callback is invoked once
	// and the list is dropped.
	subsMu sync.Mutex
	subs   []func(T, error)

	// startedOnce ensures fn is invoked at most once. The producer
	// goroutine is started lazily on the first Subscribe / Await.
	startedOnce sync.Once

	// doneCh closes once on terminal event; Await / AwaitWithContext
	// select on it.
	doneCh chan struct{}

	// Cached terminal result.
	done atomic.Bool
	val  T
	err  error
}

// NewSingle returns a Single that will execute fn when its first
// Await / Subscribe call observes an un-completed terminal event.
//
// fn is invoked at most once, no matter how many subscribers attach.
// fn is not retried on error; callers must compose with another
// source if retry is desired.
func NewSingle[T any](fn func() (T, error)) *Single[T] {
	return &Single[T]{fn: fn, doneCh: make(chan struct{})}
}

// Subscribe registers onSuccess and onError as the terminal callbacks
// for this Single. Subscribe is non-blocking: it starts the source
// goroutine (if not already started) and returns a Subscription whose
// Done channel closes on terminal event.
//
// Exactly one of onSuccess / onError is invoked. If either is nil,
// the corresponding branch is dropped silently — callers who only
// care about one outcome can pass nil for the other.
// Callbacks run on the goroutine that produced the result; do not
// rely on them being on a specific thread.
func (s *Single[T]) Subscribe(onSuccess func(T), onError func(error)) Subscription {
	cb := func(v T, err error) {
		if err != nil {
			if onError != nil {
				onError(err)
			}
			return
		}
		if onSuccess != nil {
			onSuccess(v)
		}
	}
	s.subsMu.Lock()
	s.subs = append(s.subs, cb)
	s.subsMu.Unlock()
	s.startProducer()
	return &singleSubscription[T]{done: s.doneCh, cb: cb}
}

// Await blocks until the Single has terminated and returns its result
// as (T, error). It is safe to call Await multiple times; subsequent
// calls return the same cached result without re-running fn.
//
// Await cannot be cancelled. Use [AwaitWithContext] when the caller
// needs cancellation or a deadline.
func (s *Single[T]) Await() (T, error) {
	return s.await(context.Background())
}

// AwaitWithContext blocks until the Single has terminated or ctx is
// canceled. On normal termination it returns (value, nil). On ctx
// cancellation it returns (zero, ctx.Err()); the underlying source
// is not started if ctx is already done.
func (s *Single[T]) AwaitWithContext(ctx context.Context) (T, error) {
	if s.Done() {
		return s.val, s.err
	}
	if err := ctx.Err(); err != nil {
		var zero T
		return zero, err
	}
	s.startProducer()
	select {
	case <-s.doneCh:
		return s.val, s.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Done reports whether the Single has terminated (either with a value
// or an error). Done is non-blocking and safe to call concurrently
// with Await / Subscribe.
func (s *Single[T]) Done() bool {
	return s.done.Load()
}

// Map returns a new Single that applies f to this Single's value on
// successful termination. If this Single errors, f is not invoked and
// the error is propagated unchanged.
//
// f must not return an error — Map is for value-only transforms. Use
// FlatMap when the transformation can itself fail.
func (s *Single[T]) Map[R any](f func(T) R) *Single[R] {
	return NewSingle(func() (R, error) {
		v, err := s.Await()
		if err != nil {
			var zero R
			return zero, err
		}
		return f(v), nil
	})
}

// FlatMap returns a new Single that, on this Single's success, runs
// f(value) and returns the resulting Single. If this Single errors,
// f is not invoked and the error is propagated unchanged.
//
// This is the canonical "chain another Single" operator — it lets a
// transformation itself produce an async value (e.g. another Single
// from a cache lookup).
func (s *Single[T]) FlatMap[R any](f func(T) *Single[R]) *Single[R] {
	return NewSingle(func() (R, error) {
		v, err := s.Await()
		if err != nil {
			var zero R
			return zero, err
		}
		return f(v).Await()
	})
}

// Zip returns a new Single that succeeds when both this Single and
// other succeed, with the result produced by combine(left, right).
// If either source errors, the new Single errors with that error
// (this Single's error takes precedence if both error).
//
// combine is invoked exactly once, when both values are available.
// It must not return an error — wrap the result in FlatMap if the
// combination can fail.
func (s *Single[T]) Zip[U, R any](other *Single[U], combine func(T, U) R) *Single[R] {
	return NewSingle(func() (R, error) {
		l, lErr := s.Await()
		if lErr != nil {
			var zero R
			return zero, lErr
		}
		r, rErr := other.Await()
		if rErr != nil {
			var zero R
			return zero, rErr
		}
		return combine(l, r), nil
	})
}

// AndThen returns a new Single that, on this Single's success, runs
// next and returns its result. On this Single's failure, next is not
// invoked and the error is propagated unchanged. AndThen ignores the
// value of this Single — use FlatMap if the next Single depends on it.
//
// R can differ from T; the return type is the next Single's element
// type.
func (s *Single[T]) AndThen[R any](next *Single[R]) *Single[R] {
	return NewSingle(func() (R, error) {
		if _, err := s.Await(); err != nil {
			var zero R
			return zero, err
		}
		return next.Await()
	})
}

// startProducer launches the producer goroutine that runs fn. Idempotent.
func (s *Single[T]) startProducer() {
	s.startedOnce.Do(func() {
		go s.runProducer()
	})
}

// runProducer is the single producer goroutine. It runs fn exactly
// once, caches the result, then invokes every subscriber callback
// and closes doneCh.
func (s *Single[T]) runProducer() {
	v, err := s.fn()
	s.val = v
	s.err = err
	s.done.Store(true)

	s.subsMu.Lock()
	subs := s.subs
	s.subs = nil
	s.subsMu.Unlock()

	for _, cb := range subs {
		cb(v, err)
	}
	close(s.doneCh)
}

func (s *Single[T]) await(ctx context.Context) (T, error) {
	if s.Done() {
		return s.val, s.err
	}
	s.startProducer()
	select {
	case <-s.doneCh:
		return s.val, s.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// singleSubscription is the Subscription returned by Single.Subscribe.
// It wraps the Single's doneCh with Cancel semantics: Cancel does not
// stop the producer (a Single always completes its fn); it merely
// signals the caller that it is no longer interested.
type singleSubscription[T any] struct {
	done <-chan struct{}
	cb   func(T, error)
	once sync.Once
}

func (sub *singleSubscription[T]) Request(uint64) {
	// Single emits at most one value; demand is implicitly satisfied.
}

func (sub *singleSubscription[T]) Cancel() {
	// Removal is best-effort: if the callback has already fired, the
	// cb reference is stale and Cancel is a no-op for the caller. We
	// only nill out our local reference so the closure can be GC'd.
	sub.once.Do(func() {
		sub.cb = nil
	})
}

func (sub *singleSubscription[T]) Done() <-chan struct{} {
	return sub.done
}