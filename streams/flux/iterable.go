package flux

import (
	"context"
	"streams"
)

type iterable[T any] struct {
	iter streams.Iterable[T]
	ctx  context.Context
}

type iterSubscription[T any] struct {
	iter   streams.Iterable[T]
	actual streams.Subscriber
	ctx    context.Context
	cancel func()
}

func FromIterable[T any](iter streams.Iterable[T]) streams.Flux {
	return &flux{
		source: &iterable[T]{iter: iter, ctx: context.Background()},
	}
}

func (i *iterable[T]) Subscribe(actual streams.Subscriber) {
	actual.OnSubScribe(newIterSubscription(i.iter, actual, i.ctx))
}

func newIterSubscription[T any](iter streams.Iterable[T], actual streams.Subscriber, ctx context.Context) streams.Subscription {
	sub, cancel := context.WithCancel(ctx)
	return &iterSubscription[T]{
		iter:   iter,
		actual: actual,
		ctx:    sub,
		cancel: cancel,
	}
}

func (f *iterSubscription[T]) Request(n int64) {
	for f.iter.HasNext(f.ctx) {
		val, err := f.iter.Next(f.ctx)
		if err != nil {
			f.actual.OnError(err)
		} else {
			f.actual.OnNext(val)
		}
	}
}

func (f *iterSubscription[T]) Cancel() {
	f.cancel()
}
