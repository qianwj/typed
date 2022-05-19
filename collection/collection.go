package collection

type Collection[T any] interface {
	Contains(e T) bool
	Foreach(handle func(e T))
}
