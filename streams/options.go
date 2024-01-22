package streams

type subscribeOptions[V any] struct {
	cancelOnError bool
	requestNum    int
	onNext        OnNext[V]
	onErr         OnError
	onComplete    OnComplete
}

func CancelOnError[V any]() SubscribeOption[V] {
	return func(opts *subscribeOptions[V]) {
		if !opts.cancelOnError {
			opts.cancelOnError = true
		}
	}
}

func RequestNum[V any](n int) SubscribeOption[V] {
	return func(opts *subscribeOptions[V]) {
		if opts.requestNum == 0 && n > 0 {
			opts.requestNum = n
		}
	}
}

func DoOnError[V any](onErr OnError) SubscribeOption[V] {
	return func(opts *subscribeOptions[V]) {
		if opts.onErr == nil {
			opts.onErr = onErr
		}
	}
}

func DoOnComplete[V any](onComplete OnComplete) SubscribeOption[V] {
	return func(opts *subscribeOptions[V]) {
		if opts.onComplete == nil {
			opts.onComplete = onComplete
		}
	}
}

func merge[V any](opts ...SubscribeOption[V]) *subscribeOptions[V] {
	merged := &subscribeOptions[V]{}
	for _, opt := range opts {
		opt(merged)
	}
	if merged.onNext == nil {
		merged.onNext = func(V) {}
	}
	if merged.onErr == nil {
		merged.onErr = func(error) {}
	}
	if merged.onComplete == nil {
		merged.onComplete = func() {}
	}
	if merged.requestNum == 0 {
		merged.requestNum = defaultRequestNum
	}
	return merged
}
