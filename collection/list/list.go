package list

type List[T any] interface {
	Get(index int) T
}
