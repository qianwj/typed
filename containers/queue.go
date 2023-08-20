package containers

type Queue[T any] interface {
	Enqueue(e T) bool
	Dequeue() (T, bool)
}
