package flux

import (
	"streams"
	"streams/item"
	"streams/publisher"
)

type fluxArray struct {
	val []item.Item
}

func (f *fluxArray) Subscribe(actual streams.Subscriber) {
	actual.OnSubscribe(publisher.NewArraySubscription(actual, f.val))
}

func (f *fluxArray) Consume(consumer func(it any)) {
	f.Subscribe(streams.NewFunctionalSubscriber(consumer, nil))
}
