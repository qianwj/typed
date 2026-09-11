// Package collections holds the parent-package helpers that the
// sibling subpackages (lists, maps, sets, stream) build on, plus
// the stand-alone linear data structures that do not warrant
// their own subpackage yet.
//
// # Linear data structures
//
// Stack[T] and Queue[T] are the two linear structures defined
// directly in this package. Both are concrete generic types
// (not interfaces) for the same reason the lists package uses
// concrete types: Go 1.27's generic methods work on concrete
// receivers, and a useful Stack / Queue interface that carried
// every combinator would force intermediate results to escape
// to heap.
//
// # Optional-based access
//
// Stack.Pop and Stack.Peek return option.Optional[T] rather
// than the (T, bool) shape, matching the convention established
// by ArrayList.First / Last / Find and the rest of the project:
// a present Optional on a hit, an absent one on a miss. This
// lets callers chain OrElse / OrElseGet / Map directly on the
// return value without unpacking.
//
// # Why slice-backed Stack
//
// A Stack is conceptually just a slice with two operations. The
// implementations below take that literally: items is a []T, and
// Push is a one-liner append. The interesting bit is Pop, which
// has to zero the slot it removes — see the Pop doc comment for
// why.
//
// # When to use a Queue
//
// Queue[T] is a placeholder; its concrete operations are not yet
// implemented in this commit. Until they are, callers needing
// FIFO should use a stream.Stream[T] or a slice with two index
// pointers. The type is exposed so that other code in the
// project can already refer to *Queue[T] without churn.
package collections
