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
// Stack.Pop, Stack.Peek, Queue.Pop, and Queue.Peek return
// option.Optional[T] rather than the (T, bool) shape, matching
// the convention established by ArrayList.First / Last / Find
// and the rest of the project: a present Optional on a hit,
// an absent one on a miss. This lets callers chain OrElse /
// OrElseGet / Map directly on the return value without
// unpacking.
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
// Queue[T] is a thin wrapper over lists.ArrayList[T]. The
// ArrayList's head-offset layout makes Push, Pop, and Peek all
// O(1) and bounds the retained memory to the high-water mark
// of in-flight elements; Queue inherits all of that without
// duplicating the storage machinery.
//
// Why have a separate Queue at all if ArrayList does the same
// thing? Two reasons: naming, and return-type uniformity.
// ArrayList's queue-style operations are named Add /
// RemoveFirst / Get(0) and return (T, bool); Queue's are named
// Push / Pop / Peek and return option.Optional[T], so callers
// can chain OrElse or Map directly. Queue is the "I want FIFO
// and nothing else" entry point.
//
// Stack is still a standalone type with its own []T field, not
// a wrapper. Stack is the only collection in the package that
// operates exclusively at the tail, so it does not need the
// head-offset pattern and is cheaper than a wrapper would be.
package collections
