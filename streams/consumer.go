package streams

import (
	"context"
)

type FunctionalSubscriber struct {
	consumer    func(any)
	errConsumer func(error)
	ctx         context.Context
}

func NewFunctionalSubscriber(consumer func(any), errConsumer func(error)) Subscriber {
	return &FunctionalSubscriber{
		consumer:    consumer,
		errConsumer: errConsumer,
	}
}

func (f *FunctionalSubscriber) OnNext(val any) {
	f.consumer(val)
}

func (f *FunctionalSubscriber) OnError(e error) {
	if f.errConsumer != nil {
		f.errConsumer(e)
	}
}
