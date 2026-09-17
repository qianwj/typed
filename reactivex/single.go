package reactivex

import (
	"context"

	"github.com/qianwj/typed/adt"
)

// Single[T] is a lazy asynchronous computation that produces one value or
// one error. Its terminal result is an adt.Result[T]: Success(value) or
// Failure[T](err). A successful zero or nil value is still a value.
// Use [Maybe] for an optional result and [Flowable] for multiple values.
//
// # Lifecycle
//
// Create a Single with [NewSingle]; the zero value is not usable. The first
// [Single.Subscribe], [Single.Await], or [Single.AwaitWithContext] starts
// consumption. Each source runs in its own goroutine at most once, and the
// terminal result is cached for all consumers, including late subscribers.
//
// Single is a lightweight handle with value receivers. Copies share the
// same completion state, including execution and the cached result. Copying
// a handle does not start the source or create an independent computation.
//
// Subscribe delivers the result through callbacks without blocking its caller.
// Await returns the result synchronously; AwaitWithContext allows the wait to
// be canceled. [Single.Done] reports whether the result is available. Neither
// Done nor Await waits for subscription callbacks to finish. Demand tracking
// is unnecessary because a Single produces at most one value.
//
// Canceling a wait or subscription does not cancel the shared computation.
// Subscription cancellation suppresses delivery that has not started; it
// cannot interrupt an in-flight callback. See [Single.Subscribe] and
// [Single.AwaitWithContext] for their cancellation contracts.
//
// # Composition
//
// [Single.Map], [Single.FlatMap], [Single.Zip], and [Single.AndThen] return
// lazy Singles that propagate errors and transform successful values.
// They connect completion callbacks instead of blocking on upstream Await
// calls, so a chain does not need a waiting goroutine for each operator.
//
// Transforms and pending subscription callbacks run outside internal locks
// on the completing goroutine. Slow callbacks delay other continuations on
// that goroutine. Cached subscriptions are delivered asynchronously.
// Concurrent consumers share the same result; callers must synchronize any
// mutation of data referenced by the result or shared between callbacks.
type Single[T any] struct {
	state *completion[T]
}

// NewSingle returns a Single that will execute fn when its first
// Await / Subscribe call starts the source. fn returns Success(value)
// or Failure[T](err); the same Result is cached and returned by Await.
//
// fn is invoked at most once, no matter how many subscribers attach.
// fn is not retried on error; callers must compose with another
// source if retry is desired.
func NewSingle[T any](fn func() adt.Result[T]) Single[T] {
	return newSingle(func(complete func(adt.Result[T])) {
		go func() { complete(fn()) }()
	})
}

// newSingle builds a lazy node whose starter connects continuations or
// executes a source. It does not add a goroutine to an internal chain.
func newSingle[T any](start func(func(adt.Result[T]))) Single[T] {
	return Single[T]{state: newCompletion(start)}
}

// Subscribe registers onSuccess and onError as the terminal callbacks
// for this Single. Subscribe starts consumption asynchronously and returns
// a Subscription whose Done channel closes after its callback returns or
// when canceled. Late subscribers receive the cached outcome asynchronously.
// Cancel suppresses a callback that has not started; it does not stop the
// shared source or a callback already in progress.
//
// Unless canceled before delivery, exactly one terminal callback is invoked.
// Nil callbacks are dropped silently; callers who only care about one outcome
// can pass nil for the other.
// Callbacks run outside internal locks, on the completing goroutine or an
// asynchronous cached-result delivery; do not depend on a specific goroutine.
func (s Single[T]) Subscribe(onSuccess func(T), onError func(error)) Subscription {
	cb := func(result adt.Result[T]) {
		if result.IsFailure() {
			if onError != nil {
				onError(result.Error())
			}
			return
		}
		if onSuccess != nil {
			onSuccess(result.Value())
		}
	}
	return s.state.subscribe(cb)
}

// Await blocks until the Single has terminated and returns its result
// as adt.Result[T]. Source errors produce Failure; successful values,
// including nil, produce Success. Repeated calls return the same cached result
// without re-running fn. Await does not wait for subscription callbacks.
//
// Await cannot be cancelled. Use [AwaitWithContext] when the caller
// needs cancellation or a deadline.
func (s Single[T]) Await() adt.Result[T] {
	return s.state.await(context.Background())
}

// AwaitWithContext blocks until the Single has terminated or ctx is
// canceled. It returns the cached Result on termination or Failure(ctx.Err())
// on cancellation. An already completed Single takes precedence over ctx.
// Otherwise, an already canceled ctx prevents the source from starting.
// Canceling a wait does not cancel an already running source.
func (s Single[T]) AwaitWithContext(ctx context.Context) adt.Result[T] {
	return s.state.await(ctx)
}

// Done reports whether the Single has terminated (either with a value
// or an error). Done is non-blocking and safe to call concurrently
// with Await / Subscribe.
func (s Single[T]) Done() bool {
	return s.state.isDone()
}

// Map returns a new Single that applies f to this Single's value on
// successful termination. If this Single errors, f is not invoked and
// the error is propagated unchanged.
//
// f must not return an error — Map is for value-only transforms. Use
// FlatMap when the transformation can itself fail.
func (s Single[T]) Map[R any](f func(T) R) Single[R] {
	return newSingle(func(complete func(adt.Result[R])) {
		s.state.onComplete(func(result adt.Result[T]) {
			complete(result.Map(f))
		})
	})
}

// FlatMap returns a new Single that, on this Single's success, runs
// f(value) and returns the resulting Single. If this Single errors,
// f is not invoked and the error is propagated unchanged.
//
// This is the canonical "chain another Single" operator — it lets a
// transformation itself produce an async value (e.g. another Single
// from a cache lookup).
func (s Single[T]) FlatMap[R any](f func(T) Single[R]) Single[R] {
	return newSingle(func(complete func(adt.Result[R])) {
		s.state.onComplete(func(result adt.Result[T]) {
			if result.IsFailure() {
				complete(adt.Failure[R](result.Error()))
				return
			}
			f(result.Value()).state.onComplete(complete)
		})
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
func (s Single[T]) Zip[U, R any](other Single[U], combine func(T, U) R) Single[R] {
	return s.FlatMap(func(left T) Single[R] {
		return other.Map(func(right U) R {
			return combine(left, right)
		})
	})
}

// AndThen returns a new Single that, on this Single's success, runs
// next and returns its result. On this Single's failure, next is not
// invoked and the error is propagated unchanged. AndThen ignores the
// value of this Single — use FlatMap if the next Single depends on it.
//
// R can differ from T; the return type is the next Single's element
// type.
func (s Single[T]) AndThen[R any](next Single[R]) Single[R] {
	return s.FlatMap(func(T) Single[R] {
		return next
	})
}
