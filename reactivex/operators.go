package reactivex

import (
	"context"
	"sync"
)

// Map applies f to each upstream value and forwards its result as an R.
// It is lazy: f is not invoked until the returned Observable is subscribed to.
// On a non-nil error, Map forwards that error and cancels the upstream handle.
// Panics from f are not recovered or converted into errors.
//
// f receives the subscription context and executes inline in the upstream
// notification path. Map adds neither a queue nor a worker goroutine, and
// forwards demand unchanged. The result type R may differ from T because
// Observable is a concrete type with a Go 1.27 generic method.
func (o Observable[T]) Map[R any](f func(context.Context, T) (R, error)) Observable[R] {
	return Observable[R]{subscribe: func(ctx context.Context, out Subscriber[R]) Subscription {
		return o.Subscribe(ctx, &mapSubscriber[T, R]{out: out, ctx: ctx, f: f})
	}}
}

// mapSubscriber carries the context and cancellation handle for one Map
// subscription; once prevents repeated transformation-error notifications.
type mapSubscriber[T, R any] struct {
	out  Subscriber[R]
	ctx  context.Context
	f    func(context.Context, T) (R, error)
	sub  Subscription
	once sync.Once
}

func (s *mapSubscriber[T, R]) OnSubscribe(sub Subscription) { s.sub = sub; s.out.OnSubscribe(sub) }

// Map's adapter forwards lifecycle events and keeps transformation errors
// terminal for this subscription.
func (s *mapSubscriber[T, R]) OnNext(v T) {
	r, err := s.f(s.ctx, v)
	if err != nil {
		s.once.Do(func() { s.out.OnError(err); s.sub.Cancel() })
		return
	}
	s.out.OnNext(r)
}
func (s *mapSubscriber[T, R]) OnError(err error) { s.out.OnError(err) }
func (s *mapSubscriber[T, R]) OnComplete()       { s.out.OnComplete() }

// Filter forwards values for which predicate returns true, preserving their
// upstream order. The predicate runs during delivery, not construction, and
// this operator creates no queue or goroutine.
//
// For every rejected value, Filter requests one replacement from upstream so
// rejected values do not exhaust downstream demand. A request for one match
// can therefore inspect many upstream values, or wait indefinitely on an
// infinite source with no matches. Predicate panics are not recovered.
func (o Observable[T]) Filter(predicate func(T) bool) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		return o.Subscribe(ctx, &filterSubscriber[T]{out: out, predicate: predicate})
	}}
}

// filterSubscriber compensates upstream demand for values it does not forward.
type filterSubscriber[T any] struct {
	out       Subscriber[T]
	sub       Subscription
	predicate func(T) bool
}

func (s *filterSubscriber[T]) OnSubscribe(sub Subscription) { s.sub = sub; s.out.OnSubscribe(sub) }

// A rejected value is invisible downstream but still consumes one upstream
// demand unit; OnNext replenishes it before waiting for the next match.
func (s *filterSubscriber[T]) OnNext(v T) {
	if s.predicate(v) {
		s.out.OnNext(v)
	} else {
		s.sub.Request(1)
	}
}
func (s *filterSubscriber[T]) OnError(err error) { s.out.OnError(err) }
func (s *filterSubscriber[T]) OnComplete()       { s.out.OnComplete() }

// Take forwards at most n values before completing and cancelling upstream.
// If upstream completes sooner, that completion is forwarded. Take(0) still
// subscribes upstream but cancels during setup and emits no values.
//
// The count is independent for each subscription. Take does not add demand or
// clamp forwarded requests to n; the limit is enforced as values arrive.
// Cancelling upstream does not forcibly interrupt work already executing there.
func (o Observable[T]) Take(n uint64) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		return o.Subscribe(ctx, &takeSubscriber[T]{out: out, left: n})
	}}
}

// takeSubscriber stores the remaining delivery limit for a single subscription.
type takeSubscriber[T any] struct {
	out  Subscriber[T]
	sub  Subscription
	left uint64
	once sync.Once
}

func (s *takeSubscriber[T]) OnSubscribe(sub Subscription) {
	s.sub = sub
	s.out.OnSubscribe(sub)
	if s.left == 0 {
		s.sub.Cancel()
		s.out.OnComplete()
	}
}

// OnNext counts delivered values, so values after the limit are ignored even
// if an upstream source races with cancellation.
func (s *takeSubscriber[T]) OnNext(v T) {
	if s.left == 0 {
		return
	}
	s.out.OnNext(v)
	s.left--
	if s.left == 0 {
		s.once.Do(func() { s.out.OnComplete(); s.sub.Cancel() })
	}
}
func (s *takeSubscriber[T]) OnError(err error) {
	if s.left > 0 {
		s.out.OnError(err)
	}
}
func (s *takeSubscriber[T]) OnComplete() {
	if s.left > 0 {
		s.out.OnComplete()
	}
}

// Skip discards the first n upstream values, then forwards the rest in order.
// Skip(0) forwards all values, and a source that completes within the skipped
// prefix produces no output. Each subscription has its own skip counter.
//
// Skipped values are still consumed from the source. Each one requests a
// replacement so the prefix does not use up downstream's outstanding demand.
// Like Filter, Skip creates no asynchronous boundary of its own.
func (o Observable[T]) Skip(n uint64) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		return o.Subscribe(ctx, &skipSubscriber[T]{out: out, left: n})
	}}
}

// skipSubscriber tracks how much of the prefix remains and replaces demand
// consumed by that prefix.
type skipSubscriber[T any] struct {
	out  Subscriber[T]
	sub  Subscription
	left uint64
}

func (s *skipSubscriber[T]) OnSubscribe(sub Subscription) { s.sub = sub; s.out.OnSubscribe(sub) }

// OnNext replenishes demand while consuming the skipped prefix.
func (s *skipSubscriber[T]) OnNext(v T) {
	if s.left > 0 {
		s.left--
		s.sub.Request(1)
		return
	}
	s.out.OnNext(v)
}
func (s *skipSubscriber[T]) OnError(err error) { s.out.OnError(err) }
func (s *skipSubscriber[T]) OnComplete()       { s.out.OnComplete() }

// Scan emits a running accumulation. For the first input v it emits
// f(initial, v); each later value uses the previous result as the accumulator.
// initial itself is not emitted, and empty input produces no values.
//
// Each subscription initializes an accumulator from initial, but reference
// values such as maps and slices are not cloned. Use immutable accumulators or
// separate observable instances when subscriptions must not share mutable data.
// f runs inline, once per input, and demand is forwarded unchanged.
func (o Observable[T]) Scan[R any](initial R, f func(R, T) R) Observable[R] {
	return Observable[R]{subscribe: func(ctx context.Context, out Subscriber[R]) Subscription {
		return o.Subscribe(ctx, &scanSubscriber[T, R]{out: out, value: initial, f: f})
	}}
}

// scanSubscriber holds one subscription's current accumulator.
type scanSubscriber[T, R any] struct {
	out   Subscriber[R]
	value R
	f     func(R, T) R
}

func (s *scanSubscriber[T, R]) OnSubscribe(sub Subscription) { s.out.OnSubscribe(sub) }
func (s *scanSubscriber[T, R]) OnNext(v T)                   { s.value = s.f(s.value, v); s.out.OnNext(s.value) }
func (s *scanSubscriber[T, R]) OnError(err error)            { s.out.OnError(err) }
func (s *scanSubscriber[T, R]) OnComplete()                  { s.out.OnComplete() }

// Reduce uses the first input as an accumulator, combines subsequent inputs
// with f, and emits the final accumulator on normal completion. Empty input
// completes without a value; a single input is emitted without calling f.
// On upstream error, the partial accumulator is not emitted.
//
// Reduce stores one accumulator per subscription and adds no queue or goroutine.
// Unlike Collect it does not retain every input. Unlike Scan it emits only the
// final result, so a non-terminating input never yields a reduced value.
//
// Currently Reduce forwards downstream requests directly to upstream; it does
// not translate a request for one result into demand for all inputs. Use
// ForEach or ToSlice to request all inputs, or explicitly supply enough demand
// when using Subscribe. A single request can otherwise leave reduction waiting.
func (o Observable[T]) Reduce(f func(T, T) T) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		return o.Subscribe(ctx, &reduceSubscriber[T]{out: out, f: f})
	}}
}

// reduceSubscriber distinguishes an empty input from a zero-valued accumulator
// with hasValue, rather than treating the zero value of T as a sentinel.
type reduceSubscriber[T any] struct {
	out      Subscriber[T]
	f        func(T, T) T
	value    T
	hasValue bool
}

func (s *reduceSubscriber[T]) OnSubscribe(sub Subscription) { s.out.OnSubscribe(sub) }

// OnNext keeps only the accumulator; it does not emit until completion.
func (s *reduceSubscriber[T]) OnNext(v T) {
	if !s.hasValue {
		s.value, s.hasValue = v, true
		return
	}
	s.value = s.f(s.value, v)
}
func (s *reduceSubscriber[T]) OnError(err error) { s.out.OnError(err) }
func (s *reduceSubscriber[T]) OnComplete() {
	if s.hasValue {
		s.out.OnNext(s.value)
	}
	s.out.OnComplete()
}
