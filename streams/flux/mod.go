package flux

import "streams"

type flux struct {
	source      streams.Publisher
	errConsumer func(error)
}

func (f *flux) Subscribe(consumer func(any)) {
	f.source.Subscribe(streams.NewFunctionalSubscriber(consumer, f.errConsumer))
}

func (f *flux) OnError(errConsumer func(error)) streams.Flux {
	f.errConsumer = errConsumer
	return f
}

func (f *flux) Map(mapper func(any) any) streams.Flux {
	return &flux{
		source: &fluxMap{f.source, mapper},
	}
}

func (f *flux) MapT(mapFn any) streams.Flux {
	selectGenericFunc, err := newGenericFunc(
		"MapT", "mapFn", mapFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectGenericFunc.Call(item)
	}
	return f.Map(selectorFunc)
}
