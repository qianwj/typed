package reactivex

import (
	"context"
	"sync"
)

// Subscribe starts this source and its operator chain for subscriber. The
// source invokes OnSubscribe during setup so initial demand can be requested
// before values arrive. No automatic demand is added by this method.
//
// A nil ctx is treated as context.Background. subscriber must be non-nil, and
// the Observable must have been constructed with a source. The returned handle
// controls only this subscription. Callback execution may begin before this
// method returns, so initialize shared callback state before subscribing.
func (o Observable[T]) Subscribe(ctx context.Context, subscriber Subscriber[T]) Subscription {
	if ctx == nil {
		ctx = context.Background()
	}
	return o.subscribe(ctx, subscriber)
}

// ForEach subscribes with callbacks and requests the maximum uint64 demand.
// It is the convenience form for callers that want every available value
// without manually replenishing demand. Any callback may be nil, in which
// case the corresponding notification is ignored.
//
// The call returns a Subscription; it does not wait for normal completion.
// Call Cancel to stop consumption, or use ToSlice for finite result collection.
// Callbacks execute in the source's delivery context, not necessarily the
// calling goroutine. With nil onError, upstream errors are ignored.
func (o Observable[T]) ForEach(ctx context.Context, onNext func(T), onError func(error), onComplete func()) Subscription {
	s := &callbackSubscriber[T]{onNext: onNext, onError: onError, onComplete: onComplete}
	sub := o.Subscribe(ctx, s)
	if sub != nil {
		sub.Request(^uint64(0))
	}
	return sub
}

// ToSlice starts a subscription and blocks until completion, an upstream error
// or cancellation of ctx. It requests maximum demand and appends received
// values in delivery order. Empty input returns a non-nil empty slice.
//
// On normal completion the returned error is nil. An upstream error returns
// the values collected before that error together with the error. Context
// cancellation calls Cancel and returns the partial result with ctx.Err().
// If termination and context cancellation race, either outcome may be selected.
//
// Collection retains every value in memory and cannot complete normally for
// an infinite source. Bound such sources before collecting. A nil ctx means
// context.Background. Cancellation is not a wait for an already running callback.
func (o Observable[T]) ToSlice(ctx context.Context) ([]T, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	values := make([]T, 0)
	result := make(chan error, 1)
	var once sync.Once
	sub := o.ForEach(ctx, func(value T) { values = append(values, value) }, func(err error) {
		once.Do(func() { result <- err })
	}, func() {
		once.Do(func() { result <- nil })
	})
	select {
	case err := <-result:
		return values, err
	case <-ctx.Done():
		sub.Cancel()
		return values, ctx.Err()
	}
}

// Collect collects o with ToSlice, then folds the resulting values into initial
// by calling f in delivery order. R can differ from T, allowing an aggregate,
// map or caller-owned collection to be built without depending on collections.
//
// On empty input, Collect returns initial with a nil error and does not call f.
// If collection fails, it returns initial and the error without calling f for
// any partial values. It does not accumulate incrementally while data arrives:
// the current implementation retains the entire input slice first.
//
// The accumulator is passed by value; maps, slices and pointers inside it are
// not deep-copied. f runs on the calling goroutine after collection completes.
//
//	total, err := reactivex.Collect(ctx, reactivex.Just(1, 2, 3), 0,
//	    func(sum, value int) int { return sum + value })
//	// On success, total is 6.
func Collect[T, R any](ctx context.Context, o Observable[T], initial R, f func(R, T) R) (R, error) {
	result := initial
	values, err := o.ToSlice(ctx)
	if err != nil {
		return result, err
	}
	for _, value := range values {
		result = f(result, value)
	}
	return result, nil
}

// callbackSubscriber adapts optional callbacks without adding scheduling or
// synchronization. Demand is supplied by the ForEach method that creates it.
type callbackSubscriber[T any] struct {
	sub        Subscription
	onNext     func(T)
	onError    func(error)
	onComplete func()
}

func (s *callbackSubscriber[T]) OnSubscribe(sub Subscription) { s.sub = sub }

// The callback adapter deliberately does not synchronize callbacks. Sources
// define whether callbacks are serialized; callers needing shared state should
// provide their own mutex or channel.
func (s *callbackSubscriber[T]) OnNext(v T) {
	if s.onNext != nil {
		s.onNext(v)
	}
}
func (s *callbackSubscriber[T]) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}
func (s *callbackSubscriber[T]) OnComplete() {
	if s.onComplete != nil {
		s.onComplete()
	}
}
