package lists

import "github.com/qianwj/typed/collections/iterable"

// ArrayListFrom consumes source once into a new ArrayList, preserving its
// traversal order and duplicates. The new list has independent storage;
// elements themselves are not deep-copied.
func ArrayListFrom[T any](source iterable.Iterable[T]) *ArrayList[T] {
	out := NewArrayList[T]()
	source.ForEach(out.Add)
	return out
}

// LinkedListFrom consumes source once into a new LinkedList, preserving its
// traversal order and duplicates. The new list has independent storage;
// elements themselves are not deep-copied.
func LinkedListFrom[T any](source iterable.Iterable[T]) *LinkedList[T] {
	out := NewLinkedList[T]()
	source.ForEach(out.Add)
	return out
}
