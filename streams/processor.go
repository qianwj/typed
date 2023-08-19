package streams

type Publisher interface {
	Subscribe(sub Subscriber)
}

type Subscriber interface {
	OnNext(item any)
	OnError(e error)
}

type Subscription interface {
	Cancel()
	Request(n int64)
}
