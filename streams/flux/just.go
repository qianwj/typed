package flux

import (
	"context"
	"streams"
)

func Just(data ...any) streams.Flux {
	ch := make(chan any)
	go func() {
		for _, it := range data {
			ch <- it
		}
		close(ch)
	}()
	return &flux{
		source: &just{data: ch, ctx: context.Background()},
	}
}

type just struct {
	data chan any
	ctx  context.Context
}

func (f *just) Subscribe(actual streams.Subscriber) {
	actual.OnSubScribe(newJustSubscription(f.data, actual, f.ctx))
}

type justSubscription struct {
	data   chan any
	index  int
	actual streams.Subscriber
	ctx    context.Context
	cancel func()
}

func newJustSubscription(data chan any, actual streams.Subscriber, ctx context.Context) streams.Subscription {
	sub, cancel := context.WithCancel(ctx)
	return &justSubscription{
		data:   data,
		actual: actual,
		ctx:    sub,
		cancel: cancel,
	}
}

func (f *justSubscription) Request(n int64) {
	for {
		select {
		case it, ok := <-f.data:
			if !ok {
				f.actual.OnComplete()
				return
			}
			switch it.(type) {
			case error:
				f.actual.OnError(it.(error))
			default:
				f.actual.OnNext(it)
			}
		case <-f.ctx.Done():
			f.actual.OnComplete()
			return
		}
	}
}

func (f *justSubscription) Cancel() {
	f.cancel()
}
