package sets

import "github.com/qianwj/typed/collections/iterable"

// HashSetFrom consumes source once into a new HashSet. Duplicate values are
// retained only once, and iteration order is unspecified. The new set has
// independent storage; elements themselves are not deep-copied.
func HashSetFrom[T comparable](source iterable.Iterable[T]) *HashSet[T] {
	out := NewHashSet[T]()
	source.ForEach(out.Add)
	return out
}
