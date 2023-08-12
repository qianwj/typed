package streams

type Publisher interface {
	Subscribe(sub Subscriber)
}

type Subscriber interface {
	OnNext(item any)
	OnError(e error)
	OnComplete()
	OnSubscribe(s Subscription)
}

type Subscription interface {
	Cancel()
	Request(n int64)
}

type Processor interface {
	Publisher
	Subscriber
}
