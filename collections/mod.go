// Package collections holds the parent-package helpers that the
// sibling subpackages (lists, maps, sets, stream) build on, plus
// the stand-alone linear data structures that do not warrant
// their own subpackage yet.
//
// # Linear data structures
//
// Stack[T], Queue[T], and Deque[T] are the linear structures
// defined directly in this package. All three are concrete
// generic types (not interfaces) for the same reason the lists
// package uses concrete types: Go 1.27's generic methods work
// on concrete receivers, and a useful Stack / Queue / Deque
// interface that carried every combinator would force
// intermediate results to escape to heap.
//
// # Optional-based access
//
// Stack.Pop, Stack.Peek, Queue.Pop, Queue.Peek, Deque.Front,
// Deque.Back, Deque.PopFront, and Deque.PopBack return
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
// RemoveFirst / Get(0) and now also return option.Optional[T];
// Queue's are named Push / Pop / Peek and return the same
// option.Optional[T]. After the Get → Optional change, the
// distinction is mostly naming, but Queue still pays its way:
// it is the natural focal point for FIFO-only consumers and
// the place to add FIFO-specific helpers in the future.
//
// # When to use a Deque
//
// Deque[T] is a thin wrapper over lists.LinkedList[T] that
// exposes both ends: PushFront / PopFront / Front and PushBack
// / PopBack / Back. The doubly-linked backing list gives every
// operation an exact O(1) cost — including PushFront and
// RemoveFirst, which are amortised O(1) in the head-offset
// ArrayList that backs Queue. The "Front" / "Back" qualifiers
// are needed because Deque is two-ended: the bare Push / Pop /
// Peek used by Stack and Queue would not disambiguate which
// end. The trade-off is per-node memory overhead: each
// LinkedList node carries two *node[T] pointers in addition
// to the value, which is acceptable for Deque's typical
// workloads (BFS, monotonic queues, sliding windows).
//
// Stack is still a standalone type with its own []T field, not
// a wrapper. Stack is the only collection in the package that
// operates exclusively at the tail, so it does not need the
// head-offset pattern and is cheaper than a wrapper would be.
package collections

import (
	"golang.org/x/exp/constraints"

	"github.com/qianwj/typed/collections/stream"
)

// Range returns a lazy stream.Stream[T] of consecutive integer
// values starting at start (inclusive) and ending at end
// (exclusive). The values are produced one at a time on
// iteration, in order, with step 1. The Stream is the standard
// building block for iota-style "do something n times" or
// "iterate 0..n-1" pipelines:
//
//	indices := collections.Range(0, 10).Collect()  // [0, 1, ..., 9]
//	squares := collections.Range(1, 6).
//	    Map(func(x int) int { return x * x }).
//	    Collect()                                    // [1, 4, 9, 16, 25]
//
// # Half-open interval
//
// The interval is half-open: [start, end). This matches Rust's
// 0..n, Python's range(0, n), and Java's IntStream.range. The
// result has end - start elements when start < end, and zero
// elements when start >= end. start == end produces an empty
// Stream and start > end also produces an empty Stream; neither
// is an error.
//
// # Why half-open and not closed
//
// A half-open interval makes the result length equal to
// end - start, which is the most useful length for downstream
// combinators (Take, ForEach, Count). With a closed [start, end]
// interval the length would be end - start + 1, and callers would
// have to remember the off-by-one. Half-open is the convention
// used by every standard library "range" in modern languages.
//
// # Memory
//
// Range wraps a counter in the underlying iter.Seq[T]. The
// Stream's memory cost is independent of end - start: iterating
// a Range(0, 1_000_000_000) costs a few bytes, not a billion
// values.
//
// # Type constraint
//
// Range is generic over any constraints.Integer (int, int8/16/32/64,
// uint, uint8/16/32/64, uintptr). The "Integer" constraint is
// chosen deliberately over cmp.Ordered: a Range of floats would
// be ambiguous without a step, and a Range of strings is just
// a slice, so neither fits the iota-style use case the function
// is named for.
func Range[T constraints.Integer](start, end T) stream.Stream[T] {
	return stream.From(func(yield func(T) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	})
}
