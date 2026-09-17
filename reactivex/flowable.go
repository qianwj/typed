package reactivex

import "context"

// Flowable is a typed stream of zero or more values with explicit demand
// and backpressure. Values arrive through OnNext, followed by OnComplete
// or OnError when the stream terminates. Subscribers request values through
// Subscription.Request. Constructing or copying a Flowable does not start
// consumption; subscribing does.
//
// Flowable is a lightweight, immutable pipeline description, so its methods
// use value receivers. Copying it shares the source function and its captured
// resources; it does not clone those resources or an active subscription.
// Each subscription creates new operator state. Source resources are not
// necessarily duplicated: FromChannel subscriptions share their input, and
// FromSlice subscriptions share the original backing array. The zero value
// has no source and cannot be subscribed to; use a constructor such as Just.
type Flowable[T any] struct {
	// subscribe is private so sources control lifecycle setup.
	subscribe func(context.Context, Subscriber[T]) Subscription
}

var _ Publisher[any] = Flowable[any]{}

// Subscribe starts this source and its operator chain for subscriber. The
// source invokes OnSubscribe during setup so initial demand can be requested
// before values arrive. No automatic demand is added by this method.
//
// A nil ctx is treated as context.Background. subscriber must be non-nil, and
// the Flowable must have been constructed with a source. The returned handle
// controls only this subscription. Callback execution may begin before this
// method returns, so initialize shared callback state before subscribing.
func (o Flowable[T]) Subscribe(ctx context.Context, subscriber Subscriber[T]) Subscription {
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
func (o Flowable[T]) ForEach(ctx context.Context, onNext func(T), onError func(error), onComplete func()) Subscription {
	s := &callbackSubscriber[T]{onNext: onNext, onError: onError, onComplete: onComplete}
	sub := o.Subscribe(ctx, s)
	if sub != nil {
		sub.Request(^uint64(0))
	}
	return sub
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
