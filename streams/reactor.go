package streams

import "streams/item"

type Mono interface {
	Subscribe(actual Subscriber)
	//Block() *item.Item
	//BlockWith(d time.Duration) *item.Item
}

type Flux interface {
	Subscribe(actual Subscriber)
	Consume(func(any))
}

type QueueSubscription interface {
	Subscription
	Poll() *item.Item
}
