package reactivex

type subscriberFuncs[T any] struct {
	onSubscribe func(Subscription)
	onNext      func(T)
	onError     func(error)
	onComplete  func()
}

func (s subscriberFuncs[T]) OnSubscribe(sub Subscription) {
	if s.onSubscribe != nil {
		s.onSubscribe(sub)
	}
}
func (s subscriberFuncs[T]) OnNext(v T) {
	if s.onNext != nil {
		s.onNext(v)
	}
}
func (s subscriberFuncs[T]) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}
func (s subscriberFuncs[T]) OnComplete() {
	if s.onComplete != nil {
		s.onComplete()
	}
}
