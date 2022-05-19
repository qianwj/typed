package collection

type Collection[T any] interface {
	Add(e T)
	AddAll(c Collection[T])
	Foreach(handle func(e T))
	Size() int
}
