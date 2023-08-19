package flux

import (
	"streams"
	"streams/item"
)

func Just[T any](value ...T) streams.Flux {
	if len(value) == 0 {
		return &fluxEmpty{}
	}
	if len(value) == 1 {
		return &fluxJust{val: item.Ok(value[0])}
	}
	items := make([]item.Item, len(value))
	for i, val := range value {
		items[i] = item.Ok(val)
	}
	return &fluxArray{val: items}
}

func FromChannel(ch chan<- any) streams.Flux {
	return nil
}

type fluxEmpty struct{}

func (e *fluxEmpty) Subscribe(streams.Subscriber) {}

func (e *fluxEmpty) Consume(func(any)) {

}
