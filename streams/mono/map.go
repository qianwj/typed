package mono

import (
	"streams"
	"streams/item"
)

type monoMap struct {
	source streams.Mono
	mapper func(item.Item) item.Item
}

func (m *monoMap) Subscribe(sub streams.Subscriber) {

}
