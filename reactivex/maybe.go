package reactivex

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/qianwj/typed/adt"
)

// Maybe[T] is a reactive container that emits zero or one value, or
// one error. It generalises [Single] by adding a 'no value' terminal
// state — the producer can decide that the absence of a value is a
// legitimate outcome rather than an error.
//
// Terminal states:
//
//   - OnSuccess(T)  — exactly one value delivered.
//   - OnComplete()  — no value delivered; the absence is normal.
//   - OnError(err)  — abnormal termination.
//
// Use Maybe when the absence of a value is meaningful, e.g. "looked
// up the cache key, no entry" or "scanned the topic, no message
// pending". Use [Single] when missing data is always an error, and
// [Flowable] when there can be many values.
//
// Construct with [NewMaybe]; consume via [Subscribe], [Await] /
// [AwaitWithContext], or [Done] for a non-blocking check. Compose
// with [Map] / [FlatMap] / [Zip] / [AndThen].
//
// # Comparison with Single and Flowable
//
//   - Single: always emits exactly one terminal event (value or
//     error); never a 'no value' state.
//   - Maybe: emits one of three terminal events (value, no-value,
//     error).
//   - Flowable: emits 0..N values, then complete or error.
//
// Like [Single], fn runs at most once and the result is shared
// across every subscriber.
type Maybe[T any] struct {
	fn func() adt.Result[adt.Option[T]]

	subsMu      sync.Mutex
	subs        []func(adt.Result[adt.Option[T]])
	startedOnce sync.Once
	doneCh      chan struct{}

	done   atomic.Bool
	result adt.Result[adt.Option[T]]
}

// NewMaybe returns a lazy Maybe whose source produces a Result[Option[T]].
// Success(Of(value)) emits a value, Success(Empty[T]()) completes without
// a value, and Failure[Option[T]](err) emits an error. Of preserves zero
// and nil values as present.
//
// fn is invoked at most once, on the first Await or Subscribe, and its
// result is cached for subsequent waits.
func NewMaybe[T any](fn func() adt.Result[adt.Option[T]]) *Maybe[T] {
	return &Maybe[T]{fn: fn, doneCh: make(chan struct{})}
}

// Subscribe registers the terminal callbacks for this Maybe.
// Subscribe is non-blocking: it starts the source goroutine (if not
// already started) and returns a Subscription whose Done channel
// closes on terminal event.
//
// Exactly one of onSuccess / onComplete / onError is invoked. Nil
// callbacks are dropped silently — callers who only care about one
// outcome can pass nil for the other two.
func (m *Maybe[T]) Subscribe(
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
	m.subsMu.Lock()
	m.subs = append(m.subs, cb)
	m.subsMu.Unlock()
	m.startProducer()
	return &maybeSubscription[T]{done: m.doneCh, cb: cb}
}

// Await blocks until the Maybe terminates. A value produces Success(Of(value)),
// empty completion produces Success(Empty[T]()), and a source error produces
// Failure[Option[T]](err). Presence is determined by the inner Option,
// so an emitted zero or nil value can remain present.
// Repeated calls return the cached outcome without re-running the source.
//
// Await cannot be cancelled; use [AwaitWithContext] for cancellation
// or deadlines.
func (m *Maybe[T]) Await() adt.Result[adt.Option[T]] {
	return m.AwaitWithContext(context.Background())
}

// AwaitWithContext blocks until the Maybe has terminated or ctx is
// canceled. Termination uses the same three outcomes as Await; cancellation
// returns Failure[Option[T]](ctx.Err()), not successful empty completion.
// An already completed Maybe takes precedence over ctx. Otherwise, an already
// canceled ctx prevents the source from starting. Canceling a wait does not
// cancel an already running source.
func (m *Maybe[T]) AwaitWithContext(ctx context.Context) adt.Result[adt.Option[T]] {
	return m.await(ctx)
}

// Done reports whether the Maybe has terminated. Done is non-blocking
// and safe to call concurrently with Await / Subscribe.
func (m *Maybe[T]) Done() bool {
	return m.done.Load()
}

// Map returns a new Maybe that applies f to this Maybe's value on
// the success branch. If this Maybe completes without a value or
// errors, f is not invoked and that outcome is propagated
// unchanged.
func (m *Maybe[T]) Map[R any](f func(T) R) *Maybe[R] {
	return NewMaybe(func() adt.Result[adt.Option[R]] {
		return m.Await().Map(func(value adt.Option[T]) adt.Option[R] {
			return value.Map(f)
		})
	})
}

// FlatMap returns a new Maybe that, on this Maybe's success, runs
// f(value) and returns the resulting Maybe. On this Maybe's
// completion-without-value or error, f is not invoked and that
// outcome is propagated.
func (m *Maybe[T]) FlatMap[R any](f func(T) *Maybe[R]) *Maybe[R] {
	return NewMaybe(func() adt.Result[adt.Option[R]] {
		return m.Await().FlatMap(func(value adt.Option[T]) adt.Result[adt.Option[R]] {
			if value.IsEmpty() {
				return adt.Success(adt.Empty[R]())
			}
			return f(value.Get()).Await()
		})
	})
}

// Zip returns a new Maybe that succeeds with combine(left, right)
// only if both this Maybe and other succeed with values. If either
// completes without a value, the result completes without a value.
// If either errors, the error is propagated (this Maybe's error
// wins when both error).
func (m *Maybe[T]) Zip[U, R any](other *Maybe[U], combine func(T, U) R) *Maybe[R] {
	return m.FlatMap(func(left T) *Maybe[R] {
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
func (m *Maybe[T]) AndThen[R any](next *Maybe[R]) *Maybe[R] {
	return m.FlatMap(func(T) *Maybe[R] {
		return next
	})
}

func (m *Maybe[T]) startProducer() {
	m.startedOnce.Do(func() {
		go m.runProducer()
	})
}

func (m *Maybe[T]) runProducer() {
	m.result = m.fn()
	m.done.Store(true)

	m.subsMu.Lock()
	subs := m.subs
	m.subs = nil
	m.subsMu.Unlock()

	for _, cb := range subs {
		cb(m.result)
	}
	close(m.doneCh)
}

func (m *Maybe[T]) await(ctx context.Context) adt.Result[adt.Option[T]] {
	if m.Done() {
		return m.result
	}
	if err := ctx.Err(); err != nil {
		return adt.Failure[adt.Option[T]](err)
	}
	m.startProducer()
	select {
	case <-m.doneCh:
		return m.result
	case <-ctx.Done():
		return adt.Failure[adt.Option[T]](ctx.Err())
	}
}

// maybeSubscription is the Subscription returned by Maybe.Subscribe.
// As with [singleSubscription], Cancel does not stop the producer;
// it only signals that the caller is no longer interested.
type maybeSubscription[T any] struct {
	done <-chan struct{}
	cb   func(adt.Result[adt.Option[T]])
	once sync.Once
}

func (sub *maybeSubscription[T]) Request(uint64) {
	// Maybe emits at most one value; demand is implicitly satisfied.
}

func (sub *maybeSubscription[T]) Cancel() {
	sub.once.Do(func() {
		sub.cb = nil
	})
}

func (sub *maybeSubscription[T]) Done() <-chan struct{} {
	return sub.done
}
