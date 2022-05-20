package collection

type Collection[T any] interface {
	Add(e T)
	AddAll(c Collection[T])
	Foreach(handle func(e T))
	ForeachIndexed(handle func(index int, e T))
	Filter(filter func(e T) bool) Collection[T]
	Slice() []T
	Size() int
	IsEmpty() bool
	Clear()
}
