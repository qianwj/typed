package streams

import (
	"context"
)

type FunctionalSubscriber struct {
	consumer          func(any)
	errConsumer       func(error)
	completionHandler func()
	ctx               context.Context
}

func NewFunctionalSubscriber(consumer func(any), errConsumer func(error)) Subscriber {
	return &FunctionalSubscriber{
		consumer:    consumer,
		errConsumer: errConsumer,
	}
}

func (f *FunctionalSubscriber) OnSubScribe(sub Subscription) {
	sub.Request(1)
}

func (f *FunctionalSubscriber) OnNext(val any) {
	f.consumer(val)
}

func (f *FunctionalSubscriber) OnError(e error) {
	if f.errConsumer != nil {
		f.errConsumer(e)
	}
}

func (f *FunctionalSubscriber) OnComplete() {
	if f.completionHandler != nil {
		f.completionHandler()
	}
}
