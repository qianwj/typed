package mono

import (
	"streams"
)

func Create() {

}

func Just[T any](item T) func() streams.Publisher {
	return func() streams.Publisher {
		return &mono{data: item}
	}
}

type mono struct {
	data any
	err  error
}

func (m *mono) Subscribe(subscriber streams.Subscriber) {
	if m.err != nil {
		subscriber.OnError(m.err)
	} else if m.data != nil {
		subscriber.OnNext(m.data)
	}
	subscriber.OnComplete()
}
