package flux

import (
	"go.uber.org/atomic"
	"streams"
)

type fluxTakeSubscriber struct {
	actual       streams.Subscriber
	n            int64
	remaining    int64
	subscription streams.Subscription
	done         bool
	wip          atomic.Int32
}

func NewFluxTakeSubscriber(actual streams.Subscriber, n int64) streams.Subscriber {
	return &fluxTakeSubscriber{
		actual: actual,
		n:      n,
	}
}

func (f *fluxTakeSubscriber) OnSubscribe(subscription streams.Subscription) {

}

func (f *fluxTakeSubscriber) OnNext(it any) {

}

func (f *fluxTakeSubscriber) OnError(err error) {
	if f.done {
		// on error dropped
		return
	}
	f.done = true
	f.actual.OnError(err)
}

func (f *fluxTakeSubscriber) OnComplete() {
	if f.done {
		return
	}
	f.done = true
	f.actual.OnComplete()
}

func (f *fluxTakeSubscriber) Request(n int64) {

}

func (f *fluxTakeSubscriber) Cancel() {
	f.subscription.Cancel()
}
