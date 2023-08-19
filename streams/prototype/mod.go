package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Flux interface {
	Map(mapper func(any) any) Flux
	OnError(errConsumer func(error)) Flux
	Subscribe(consumer func(any))
}

type Publisher interface {
	Subscribe(sub Subscriber)
}

type Subscriber interface {
	OnNext(e any)
	OnError(e error)
}

type flux struct {
	source      Publisher
	errConsumer func(error)
}

type fluxSubscriber struct {
	consumer    func(any)
	errConsumer func(error)
}

func (f *fluxSubscriber) OnNext(e any) {
	f.consumer(e)
}

func (f *fluxSubscriber) OnError(e error) {
	if f.errConsumer != nil {
		f.errConsumer(e)
	}
}

func (f *flux) Subscribe(consumer func(any)) {
	f.source.Subscribe(&fluxSubscriber{
		consumer:    consumer,
		errConsumer: f.errConsumer,
	})
}

func (f *flux) OnError(errConsumer func(error)) Flux {
	f.errConsumer = errConsumer
	return f
}

func (f *flux) Map(mapper func(any) any) Flux {
	return &flux{
		source: &fluxMap{f.source, mapper},
	}
}

type fluxMap struct {
	upstream Publisher
	mapper   func(any) any
}

func (f *fluxMap) Subscribe(sub Subscriber) {
	f.upstream.Subscribe(&mapSubscriber{actual: sub, mapper: f.mapper})
}

type mapSubscriber struct {
	actual Subscriber
	mapper func(any) any
}

func (m *mapSubscriber) OnNext(it any) {
	m.actual.OnNext(m.mapper(it))
}

func (m *mapSubscriber) OnError(err error) {
	m.actual.OnError(err)
}

type fluxJust struct {
	data chan any
}

func just(data ...any) Flux {
	ch := make(chan any)
	go func() {
		for _, it := range data {
			ch <- it
			time.Sleep(time.Second * 1)
		}
		close(ch)
	}()
	return &flux{
		source: &fluxJust{data: ch},
	}
}

func (f *fluxJust) Subscribe(sub Subscriber) {
	for it := range f.data {
		switch it.(type) {
		case error:
			sub.OnError(it.(error))
		default:
			sub.OnNext(it)
		}
	}
}

func main() {
	f := just(1, 2, errors.New("an error"), 3, 4, 5).Map(func(a any) any {
		return strconv.Itoa(a.(int))
	}).Map(func(a any) any {
		return a.(string) + " mapped"
	}).OnError(func(e error) {
		fmt.Printf("error: %+v\r\n", e)
	})
	f.Subscribe(func(d any) {
		fmt.Printf("data: %v\r\n", d)
	})
}
