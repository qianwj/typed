package streams

import "sync/atomic"

const (
	defaultRequestNum = 1
)

type FixedSubscriber[V any] struct {
	complete      atomic.Bool
	subscription  Subscription[V]
	requestNum    int
	cancelOnError bool
	onNext        OnNext[V]
	onErr         OnError
	onComplete    OnComplete
}

func NewFixedSubscriber[V any](opts ...SubscribeOption[V]) Subscriber[V] {
	merged := merge(opts...)
	return &FixedSubscriber[V]{
		onNext:        merged.onNext,
		onErr:         merged.onErr,
		onComplete:    merged.onComplete,
		requestNum:    merged.requestNum,
		cancelOnError: merged.cancelOnError,
	}
}

func (s *FixedSubscriber[V]) OnSubscribe(subscription Subscription[V]) {
	s.subscription = subscription
	for s.complete.Load() {
		subscription.Request(s.requestNum)
	}
}

func (s *FixedSubscriber[V]) OnNext(val V) {
	s.onNext(val)
}

func (s *FixedSubscriber[V]) OnError(err error) {
	if s.cancelOnError && s.subscription != nil {
		s.subscription.Cancel()
	}
	s.onErr(err)
}

func (s *FixedSubscriber[V]) OnComplete() {
	s.complete.Store(true)
	s.onComplete()
}
