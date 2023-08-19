package streams

type Flux interface {
	Map(mapper func(any) any) Flux
	MapT(mapFn any) Flux
	OnError(errConsumer func(error)) Flux
	Subscribe(consumer func(any))
}
