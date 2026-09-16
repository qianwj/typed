// Package iterable defines the shared traversal contract for collections.
// It has no collection dependencies, so lists and sets can consume the same
// interface without depending on each other's concrete types.
package iterable

// Iterable visits its elements synchronously, calling visit once per element.
// Traversal order and whether traversal can be repeated are defined by the
// implementation. The callback must not structurally mutate the source.
//
// ArrayList, LinkedList, HashSet, and Stream satisfy this interface through
// their existing ForEach methods. It deliberately requires no size, mutation,
// or type-changing operations.
type Iterable[T any] interface {
	ForEach(visit func(T))
}
