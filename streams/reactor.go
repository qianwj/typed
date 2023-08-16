package streams

import (
	"streams/item"
	"time"
)

type Mono interface {
	Block() *item.Item
	BlockWith(d time.Duration) *item.Item
}
