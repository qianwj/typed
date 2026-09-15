package reactivex

import (
	"context"
	"sync"
	"sync/atomic"
)

// maybeKind distinguishes the three terminal states a Maybe can
// reach. We store it explicitly rather than overload the zero value
// of T (or the nil-ness of err) so callers don't have to peek at
// runtime state to know which state was reached.
type maybeKind uint8

const (
	maybeSettled maybeKind = iota // zero value: producer hasn't fired yet
	maybeSuccess
	maybeComplete
	maybeError
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
// [Observable] when there can be many values.
//
// Construct with [NewMaybe]; consume via [Subscribe], [Await] /
// [AwaitWithContext], or [Done] for a non-blocking check. Compose
// with [Map] / [FlatMap] / [Zip] / [AndThen].
//
// # Comparison with Single and Observable
//
//   - Single: always emits exactly one terminal event (value or
//     error); never a 'no value' state.
//   - Maybe: emits one of three terminal events (value, no-value,
//     error).
//   - Observable: emits 0..N values, then complete or error.
//
// Like [Single], fn runs at most once and the result is shared
// across every subscriber.
type Maybe[T any] struct {
	fn func() (T, bool, error)

	subsMu      sync.Mutex
	subs        []func(T, bool, error)
	startedOnce sync.Once
	doneCh      chan struct{}

	done     atomic.Bool
	kind     maybeKind
	val      T
	err      error
}

// NewMaybe returns a Maybe whose source function produces a
// (value, present, error) triple. present=true means the value is
// meaningful; present=false means "no value, complete normally" and
// triggers the OnComplete branch.
//
// If err != nil, the present flag is ignored and OnError fires.
// fn is invoked at most once, no matter how many subscribers attach.
func NewMaybe[T any](fn func() (T, bool, error)) *Maybe[T] {
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
	cb := func(v T, present bool, err error) {
		switch {
		case err != nil:
			if onError != nil {
				onError(err)
			}
		case !present:
			if onComplete != nil {
				onComplete()
			}
		default:
			if onSuccess != nil {
				onSuccess(v)
			}
		}
	}
	m.subsMu.Lock()
	m.subs = append(m.subs, cb)
	m.subsMu.Unlock()
	m.startProducer()
	return &maybeSubscription[T]{done: m.doneCh, cb: cb}
}

// Await blocks until the Maybe has terminated and returns
// (value, present, error). present=true means v holds a meaningful
// value; present=false means the producer completed with no value.
// On error, err is non-nil and present has no meaning.
//
// Await cannot be cancelled; use [AwaitWithContext] for cancellation
// or deadlines.
func (m *Maybe[T]) Await() (T, bool, error) {
	return m.await(context.Background())
}

// AwaitWithContext blocks until the Maybe has terminated or ctx is
// canceled. On terminal event it returns the cached result. On ctx
// cancellation it returns (zero, false, ctx.Err()) and does not
// start the producer if ctx is already done.
func (m *Maybe[T]) AwaitWithContext(ctx context.Context) (T, bool, error) {
	var zero T
	if m.Done() {
		return m.val, m.kind == maybeSuccess, m.err
	}
	if err := ctx.Err(); err != nil {
		return zero, false, err
	}
	m.startProducer()
	select {
	case <-m.doneCh:
		return m.val, m.kind == maybeSuccess, m.err
	case <-ctx.Done():
		return zero, false, ctx.Err()
	}
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
	return NewMaybe(func() (R, bool, error) {
		v, present, err := m.Await()
		if err != nil {
			var zero R
			return zero, false, err
		}
		if !present {
			var zero R
			return zero, false, nil
		}
		return f(v), true, nil
	})
}

// FlatMap returns a new Maybe that, on this Maybe's success, runs
// f(value) and returns the resulting Maybe. On this Maybe's
// completion-without-value or error, f is not invoked and that
// outcome is propagated.
func (m *Maybe[T]) FlatMap[R any](f func(T) *Maybe[R]) *Maybe[R] {
	return NewMaybe(func() (R, bool, error) {
		v, present, err := m.Await()
		if err != nil {
			var zero R
			return zero, false, err
		}
		if !present {
			var zero R
			return zero, false, nil
		}
		return f(v).Await()
	})
}

// Zip returns a new Maybe that succeeds with combine(left, right)
// only if both this Maybe and other succeed with values. If either
// completes without a value, the result completes without a value.
// If either errors, the error is propagated (this Maybe's error
// wins when both error).
func (m *Maybe[T]) Zip[U, R any](other *Maybe[U], combine func(T, U) R) *Maybe[R] {
	return NewMaybe(func() (R, bool, error) {
		var zero R
		l, lPresent, lErr := m.Await()
		if lErr != nil {
			return zero, false, lErr
		}
		if !lPresent {
			return zero, false, nil
		}
		r, rPresent, rErr := other.Await()
		if rErr != nil {
			return zero, false, rErr
		}
		if !rPresent {
			return zero, false, nil
		}
		return combine(l, r), true, nil
	})
}

// AndThen returns a new Maybe that, on this Maybe's success, runs
// next and returns its outcome. On this Maybe's completion-without-
// value or error, next is not invoked and that outcome is propagated.
// AndThen ignores the value of this Maybe — use FlatMap if next
// depends on it.
func (m *Maybe[T]) AndThen[R any](next *Maybe[R]) *Maybe[R] {
	return NewMaybe(func() (R, bool, error) {
		_, present, err := m.Await()
		if err != nil {
			var zero R
			return zero, false, err
		}
		if !present {
			var zero R
			return zero, false, nil
		}
		return next.Await()
	})
}

func (m *Maybe[T]) startProducer() {
	m.startedOnce.Do(func() {
		go m.runProducer()
	})
}

func (m *Maybe[T]) runProducer() {
	v, present, err := m.fn()
	m.val = v
	m.err = err
	switch {
	case err != nil:
		m.kind = maybeError
	case !present:
		m.kind = maybeComplete
	default:
		m.kind = maybeSuccess
	}
	m.done.Store(true)

	m.subsMu.Lock()
	subs := m.subs
	m.subs = nil
	m.subsMu.Unlock()

	for _, cb := range subs {
		cb(v, m.kind == maybeSuccess, err)
	}
	close(m.doneCh)
}

func (m *Maybe[T]) await(ctx context.Context) (T, bool, error) {
	var zero T
	if m.Done() {
		return m.val, m.kind == maybeSuccess, m.err
	}
	m.startProducer()
	select {
	case <-m.doneCh:
		return m.val, m.kind == maybeSuccess, m.err
	case <-ctx.Done():
		return zero, false, ctx.Err()
	}
}

// maybeSubscription is the Subscription returned by Maybe.Subscribe.
// As with [singleSubscription], Cancel does not stop the producer;
// it only signals that the caller is no longer interested.
type maybeSubscription[T any] struct {
	done <-chan struct{}
	cb   func(T, bool, error)
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