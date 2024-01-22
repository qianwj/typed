package streams

type Publisher[V any] interface {
	Subscribe(sub Subscriber[V])
}

type Subscriber[V any] interface {
	OnNext(item V)
	OnError(e error)
	OnComplete()
	OnSubscribe(s Subscription[V])
}

type Subscription[V any] interface {
	Cancel()
	Request(n int)
}

//
//type Processor interface {
//	Publisher
//	Subscriber
//}
