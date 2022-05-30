package core

import "errors"

type Optional[T any] interface {
	EmptyFunc(emptyFunc func(T) bool) Optional[T]
	IsPresent() bool
	IfPresentThen(process func(T) T) Optional[T]
}

type Option[T any] struct {
	data      T
	emptyFunc func(T) bool
}

func NewOption[T any](data T) *Option[T] {
	return &Option[T]{data: data}
}

func (o *Option[T]) EmptyFunc(emptyFunc func(T) bool) *Option[T] {
	if o.emptyFunc != nil {
		o.emptyFunc = emptyFunc
	}
	return o
}

func (o *Option[T]) IsPresent() bool {
	if o.emptyFunc == nil {
		panic(errors.New("Option.emptyFunc is nil"))
	}
	return !o.emptyFunc(o.data)
}

func (o *Option[T]) IfPresentThen(process func(T) T) *Option[T] {
	if !o.IsPresent() {
		return o
	}
	o.data = process(o.data)
	return o
}

func String(data string) *Option[string] {
	return &Option[string]{data: data, emptyFunc: func(s string) bool {
		return s == ""
	}}
}
