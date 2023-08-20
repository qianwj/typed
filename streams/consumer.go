package streams

import (
	"context"
	"time"
)

type FunctionalSubscriber struct {
	consumer          func(any)
	errConsumer       func(error)
	completionHandler func()
	consumes          int
	complete          bool
	ctx               context.Context
}

func NewFunctionalSubscriber(consumer func(any), errConsumer func(error)) Subscriber {
	return &FunctionalSubscriber{
		consumer:    consumer,
		errConsumer: errConsumer,
	}
}

func (f *FunctionalSubscriber) OnSubScribe(sub Subscription) {
	tick := time.NewTicker(time.Second * 2)
	for !f.complete {
		<-tick.C
		if f.consumes < 1 {
			// slow
			sub.Request(1)
		} else {
			sub.Request(10)
		}
		f.consumes = 0
	}
	tick.Stop()
}

func (f *FunctionalSubscriber) OnNext(val any) {
	f.consumer(val)
	f.consumes++
}

func (f *FunctionalSubscriber) OnError(e error) {
	if f.errConsumer != nil {
		f.errConsumer(e)
	}
}

func (f *FunctionalSubscriber) OnComplete() {
	f.complete = true
	if f.completionHandler != nil {
		f.completionHandler()
	}
}
