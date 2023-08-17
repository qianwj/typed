package sink

import (
	"streams/item"
)

func Empty() func(item.Item) item.Item {
	return func(it item.Item) item.Item {
		return it
	}
}

func Ignore() func(item.Item) {
	return func(it item.Item) {}
}
