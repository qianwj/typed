package streams

type Publisher interface {
	Subscribe(sub Subscriber)
}

type Subscriber interface {
	OnSubScribe(sub Subscription)
	OnNext(item any)
	OnError(e error)
	OnComplete()
}

type Subscription interface {
	Cancel()
	Request(n int64)
}
