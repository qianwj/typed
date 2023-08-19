package flux

import (
	"streams"
	"streams/item"
	"streams/publisher"
)

type fluxJust struct {
	val item.Item
}

func (f *fluxJust) Subscribe(actual streams.Subscriber) {
	actual.OnSubscribe(publisher.NewScalarSubscription(actual, f.val))
}

func (f *fluxJust) Consume(consumer func(it any)) {
	consumer(f.val)
}
