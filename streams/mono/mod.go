package mono

import (
	"streams"
	"streams/item"
)

func Just(val any) streams.Mono {
	return &monoJust{
		val: item.Ok(val),
	}
}
