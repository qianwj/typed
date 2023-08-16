package mono

import (
	"streams"
	"streams/item"
	"streams/publisher"
	"time"
)

type monoJust struct {
	streams.Publisher
	val item.Item
}

func (m *monoJust) Block() *item.Item {
	return &m.val
}

func (m *monoJust) BlockWith(time.Duration) *item.Item {
	return &m.val
}

func (m *monoJust) Subscribe(actual streams.Subscriber) {
	actual.OnSubscribe(publisher.NewScalarSubscription(actual, m.val))
}
