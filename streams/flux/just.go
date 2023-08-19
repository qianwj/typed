package flux

import (
	"streams"
	"time"
)

type just struct {
	data chan any
}

func (f *just) Subscribe(sub streams.Subscriber) {
	for it := range f.data {
		switch it.(type) {
		case error:
			sub.OnError(it.(error))
		default:
			sub.OnNext(it)
		}
	}
}

func Just(data ...any) streams.Flux {
	ch := make(chan any)
	go func() {
		for _, it := range data {
			ch <- it
			time.Sleep(time.Second * 1)
		}
		close(ch)
	}()
	return &flux{
		source: &just{data: ch},
	}
}
