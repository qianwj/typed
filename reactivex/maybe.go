package reactivex

import (
	"context"

	"github.com/qianwj/typed/adt"
)

// Maybe[T] is a lazy asynchronous computation that produces an optional
// value or one error. Its terminal result is an adt.Result[adt.Option[T]]:
//
//   - Success(Of(value)): a value, delivered to onSuccess.
//   - Success(Empty[T]()): normal completion without a value, delivered to onComplete.
//   - Failure[Option[T]](err): an error, delivered to onError.
//
// Presence is determined by the Option, so Of preserves zero and nil values
// as present. Use Maybe when absence is meaningful, such as a cache miss;
// use [Single] for a required result and [Flowable] for multiple values.
//
// # Lifecycle
//
// Create a Maybe with [NewMaybe]; the zero value is not usable. The first
// [Maybe.Subscribe], [Maybe.Await], or [Maybe.AwaitWithContext] starts
// consumption. Each source runs in its own goroutine at most once, and the
// terminal result is cached for all consumers, including late subscribers.
//
// Maybe is a lightweight handle with value receivers. Copies share the
// same completion state, including execution and the cached result. Copying
// a handle does not start the source or create an independent computation.
//
// Subscribe delivers the matching terminal callback without blocking its
// caller. Await returns the result synchronously; AwaitWithContext allows
// the wait to be canceled. [Maybe.Done] reports whether the result is
// available. Neither Done nor Await waits for subscription callbacks to
// finish. Demand tracking is unnecessary because a Maybe emits at most one value.
//
// Canceling a wait or subscription does not cancel the shared computation.
// Subscription cancellation suppresses delivery that has not started; it
// cannot interrupt an in-flight callback. See [Maybe.Subscribe] and
// [Maybe.AwaitWithContext] for their cancellation contracts.
//
// # Composition
//
// [Maybe.Map], [Maybe.FlatMap], [Maybe.Zip], and [Maybe.AndThen] return lazy
// Maybes that propagate errors and empty completion, applying transforms
// only when a value is present. They connect completion callbacks instead
// of blocking on upstream Await calls, so a chain does not need a waiting
// goroutine for each operator.
//
// Transforms and pending subscription callbacks run outside internal locks
// on the completing goroutine. Slow callbacks delay other continuations on
// that goroutine. Cached subscriptions are delivered asynchronously.
// Concurrent consumers share the same result; callers must synchronize any
// mutation of data referenced by the result or shared between callbacks.
type Maybe[T any] struct {
	state *completion[adt.Option[T]]
}

// NewMaybe returns a lazy Maybe whose source produces a Result[Option[T]].
// Success(Of(value)) emits a value, Success(Empty[T]()) completes without
// a value, and Failure[Option[T]](err) emits an error. Of preserves zero
// and nil values as present.
//
// fn is invoked at most once, on the first Await or Subscribe, and its
// result is cached for subsequent waits.
func NewMaybe[T any](fn func() adt.Result[adt.Option[T]]) Maybe[T] {
	return newMaybe(func(complete func(adt.Result[adt.Option[T]])) {
		go func() { complete(fn()) }()
	})
}

// newMaybe builds a lazy node whose starter connects continuations or
// executes a source. It does not add a goroutine to an internal chain.
func newMaybe[T any](start func(func(adt.Result[adt.Option[T]]))) Maybe[T] {
	return Maybe[T]{state: newCompletion(start)}
}

// Subscribe registers the terminal callbacks for this Maybe.
// Subscribe starts consumption asynchronously and returns a Subscription whose
// Done channel closes after its callback returns, or when canceled.
// Late subscribers receive the cached outcome asynchronously. Cancel suppresses
// a callback that has not started; it does not stop the shared source.
//
// Unless canceled before delivery, exactly one terminal callback is invoked. Nil
// callbacks are dropped silently — callers who only care about one
// outcome can pass nil for the other two.
func (m Maybe[T]) Subscribe(
	onSuccess func(T),
	onComplete func(),
	onError func(error),
) Subscription {
	cb := func(result adt.Result[adt.Option[T]]) {
		switch {
		case result.IsFailure():
			if onError != nil {
				onError(result.Error())
			}
		case result.Value().IsEmpty():
			if onComplete != nil {
				onComplete()
			}
		default:
			if onSuccess != nil {
				onSuccess(result.Value().Get())
			}
		}
	}
	return m.state.subscribe(cb)
}

// Await blocks until the Maybe terminates. A value produces Success(Of(value)),
// empty completion produces Success(Empty[T]()), and a source error produces
// Failure[Option[T]](err). Presence is determined by the inner Option,
// so an emitted zero or nil value can remain present.
// Repeated calls return the cached outcome without re-running the source.
// Await does not wait for subscription callbacks.
//
// Await cannot be cancelled; use [AwaitWithContext] for cancellation
// or deadlines.
func (m Maybe[T]) Await() adt.Result[adt.Option[T]] {
	return m.AwaitWithContext(context.Background())
}

// AwaitWithContext blocks until the Maybe has terminated or ctx is
// canceled. Termination uses the same three outcomes as Await; cancellation
// returns Failure[Option[T]](ctx.Err()), not successful empty completion.
// An already completed Maybe takes precedence over ctx. Otherwise, an already
// canceled ctx prevents the source from starting. Canceling a wait does not
// cancel an already running source.
func (m Maybe[T]) AwaitWithContext(ctx context.Context) adt.Result[adt.Option[T]] {
	return m.state.await(ctx)
}

// Done reports whether the Maybe has terminated. Done is non-blocking
// and safe to call concurrently with Await / Subscribe.
func (m Maybe[T]) Done() bool {
	return m.state.isDone()
}

// Map returns a new Maybe that applies f to this Maybe's value on
// the success branch. If this Maybe completes without a value or
// errors, f is not invoked and that outcome is propagated
// unchanged.
func (m Maybe[T]) Map[R any](f func(T) R) Maybe[R] {
	return newMaybe(func(complete func(adt.Result[adt.Option[R]])) {
		m.state.onComplete(func(result adt.Result[adt.Option[T]]) {
			complete(result.Map(func(value adt.Option[T]) adt.Option[R] {
				return value.Map(f)
			}))
		})
	})
}

// FlatMap returns a new Maybe that, on this Maybe's success, runs
// f(value) and returns the resulting Maybe. On this Maybe's
// completion-without-value or error, f is not invoked and that
// outcome is propagated.
func (m Maybe[T]) FlatMap[R any](f func(T) Maybe[R]) Maybe[R] {
	return newMaybe(func(complete func(adt.Result[adt.Option[R]])) {
		m.state.onComplete(func(result adt.Result[adt.Option[T]]) {
			if result.IsFailure() {
				complete(adt.Failure[adt.Option[R]](result.Error()))
				return
			}
			if result.Value().IsEmpty() {
				complete(adt.Success(adt.Empty[R]()))
				return
			}
			f(result.Value().Get()).state.onComplete(complete)
		})
	})
}

// Zip returns a new Maybe that succeeds with combine(left, right)
// only if both this Maybe and other succeed with values. If either
// completes without a value, the result completes without a value.
// If either errors, the error is propagated (this Maybe's error
// wins when both error).
func (m Maybe[T]) Zip[U, R any](other Maybe[U], combine func(T, U) R) Maybe[R] {
	return m.FlatMap(func(left T) Maybe[R] {
		return other.Map(func(right U) R {
			return combine(left, right)
		})
	})
}

// AndThen returns a new Maybe that, on this Maybe's success, runs
// next and returns its outcome. On this Maybe's completion-without-
// value or error, next is not invoked and that outcome is propagated.
// AndThen ignores the value of this Maybe — use FlatMap if next
// depends on it.
func (m Maybe[T]) AndThen[R any](next Maybe[R]) Maybe[R] {
	return m.FlatMap(func(T) Maybe[R] {
		return next
	})
}
