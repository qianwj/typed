package streams

type (
	OnNext[V any]          func(V)
	OnError                func(error)
	OnComplete             func()
	SubscribeOption[V any] func(opts *subscribeOptions[V])
)
