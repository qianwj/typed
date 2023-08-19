package flux

import "streams"

type fluxMap struct {
	upstream streams.Publisher
	mapper   func(any) any
}

func (f *fluxMap) Subscribe(sub streams.Subscriber) {
	f.upstream.Subscribe(&mapSubscriber{actual: sub, mapper: f.mapper})
}

type mapSubscriber struct {
	actual streams.Subscriber
	mapper func(any) any
}

func (m *mapSubscriber) OnNext(it any) {
	m.actual.OnNext(m.mapper(it))
}

func (m *mapSubscriber) OnError(err error) {
	m.actual.OnError(err)
}
