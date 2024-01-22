package streams

type Item[V any] struct {
	val V
	err error
}

func NewItem[V any](val V, err error) *Item[V] {
	return &Item[V]{val: val, err: err}
}

func (i *Item[V]) IsError() bool {
	return i.err != nil
}

func (i *Item[V]) Err() error {
	return i.err
}

func (i *Item[V]) IsOk() bool {
	return i.val != nil
}

func (i *Item[V]) Ok() V {
	return i.val
}
