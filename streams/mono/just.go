package mono

import (
	"streams"
	"streams/item"
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

func (m *monoJust) Subscribe(sub streams.Subscriber) {
	//ScalarSubscription
	sub.OnSubscribe(nil)
}
