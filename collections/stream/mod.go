// Package stream provides Stream[T], a type-safe, lazy, fluent
// data-processing pipeline backed by Go's iter.Seq[T].
//
// # Stream is a concrete generic type, not an interface
//
// Go 1.27 added type parameters to methods of concrete types. Interface
// methods still cannot declare their own type parameters, because doing so
// would force the compiler to instantiate a concrete type with every
// possible type argument reachable through the interface value, which is
// impractical with Go's per-instantiation code-generation approach.
//
// Because fluent operations such as Map change the element type and need a
// method-level type parameter (e.g. Map[R any](f func(T) R) Stream[R]),
// Stream must be a concrete generic type, not an interface, to support the
// intended API.
//
// Internally, a Stream wraps a single-use iter.Seq[T] (a Go 1.23 function
// iterator). Intermediate operations return a new Stream that wraps the
// previous sequence with a transformation; terminal operations consume the
// sequence and produce a concrete result.
//
// Typical usage:
//
//	profiles := arraylist.Of(users).
//	    Stream().
//	    Filter(isAdult).
//	    Map(toProfile).
//	    Take(100).
//	    Collect()
//
// A Stream is single-use: each terminal operation consumes the underlying
// iterator. Save the result of Collect / Count / Reduce rather than reusing
// the Stream.
package stream

import (
	"iter"
	"slices"
)

// Stream is a lazy, fluent pipeline of values of type T.
//
// A Stream is built from a single iter.Seq[T] and is single-use. Each
// intermediate operation returns a new Stream; each terminal operation
// consumes the underlying sequence.
type Stream[T any] struct {
	seq iter.Seq[T]
}

// ---------- Constructors ----------

// From builds a Stream from a function iterator.
func From[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}

// FromSlice builds a Stream from a slice.
func FromSlice[T any](s []T) Stream[T] {
	return From(slices.Values(s))
}

// Of builds a Stream from individual values.
func Of[T any](values ...T) Stream[T] {
	return FromSlice(values)
}

// Empty returns a Stream that yields no values.
func Empty[T any]() Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {}}
}

// ---------- Intermediate operations (return Stream) ----------

// Filter keeps only the elements for which p returns true.
func (s Stream[T]) Filter(p func(T) bool) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range s.seq {
			if p(v) {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// Map applies f to every element and yields the results.
//
// This is the operation that motivates Stream being a concrete generic type
// in Go 1.27: the method declares its own type parameter R, which an
// interface method cannot do.
func (s Stream[T]) Map[R any](f func(T) R) Stream[R] {
	return Stream[R]{seq: func(yield func(R) bool) {
		for v := range s.seq {
			if !yield(f(v)) {
				return
			}
		}
	}}
}

// FlatMap applies f to every element and concatenates the resulting Streams.
//
// f returns a Stream[R]; to FlatMap with a raw iter.Seq[R], wrap it with
// stream.From first.
func (s Stream[T]) FlatMap[R any](f func(T) Stream[R]) Stream[R] {
	return Stream[R]{seq: func(yield func(R) bool) {
		for v := range s.seq {
			for r := range f(v).seq {
				if !yield(r) {
					return
				}
			}
		}
	}}
}

// Peek yields every element unchanged but also calls visit on each one.
// Useful for debugging or observing a pipeline without modifying it.
func (s Stream[T]) Peek(visit func(T)) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range s.seq {
			visit(v)
			if !yield(v) {
				return
			}
		}
	}}
}

// Take keeps at most the first n elements.
func (s Stream[T]) Take(n int) Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}
	return Stream[T]{seq: func(yield func(T) bool) {
		taken := 0
		for v := range s.seq {
			if !yield(v) {
				return
			}
			taken++
			if taken >= n {
				return
			}
		}
	}}
}

// Drop discards the first n elements.
func (s Stream[T]) Drop(n int) Stream[T] {
	if n <= 0 {
		return s
	}
	return Stream[T]{seq: func(yield func(T) bool) {
		skipped := 0
		for v := range s.seq {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}}
}

// Distinct keeps only the first occurrence of each element under eq.
// eq reports whether two values are considered equal.
func (s Stream[T]) Distinct(eq func(T, T) bool) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		var seen []T
		for v := range s.seq {
			duplicate := false
			for _, x := range seen {
				if eq(x, v) {
					duplicate = true
					break
				}
			}
			if duplicate {
				continue
			}
			seen = append(seen, v)
			if !yield(v) {
				return
			}
		}
	}}
}

// Concat appends the elements of other to this Stream.
func (s Stream[T]) Concat(other Stream[T]) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range s.seq {
			if !yield(v) {
				return
			}
		}
		for v := range other.seq {
			if !yield(v) {
				return
			}
		}
	}}
}

// ---------- Terminal operations ----------

// Collect materializes the Stream into a slice.
func (s Stream[T]) Collect() []T {
	out := make([]T, 0)
	for v := range s.seq {
		out = append(out, v)
	}
	return out
}

// Count returns the number of elements.
func (s Stream[T]) Count() int {
	n := 0
	for range s.seq {
		n++
	}
	return n
}

// First returns the first element, or the zero value and false if the Stream
// is empty.
func (s Stream[T]) First() (T, bool) {
	for v := range s.seq {
		return v, true
	}
	var zero T
	return zero, false
}

// Last returns the last element, or the zero value and false if the Stream
// is empty.
func (s Stream[T]) Last() (T, bool) {
	var (
		zero  T
		last  T
		found bool
	)
	for v := range s.seq {
		last = v
		found = true
	}
	if !found {
		return zero, false
	}
	return last, true
}

// Any reports whether at least one element satisfies p.
// Stops as soon as one match is found.
func (s Stream[T]) Any(p func(T) bool) bool {
	for v := range s.seq {
		if p(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies p.
// Stops as soon as one mismatch is found.
func (s Stream[T]) All(p func(T) bool) bool {
	for v := range s.seq {
		if !p(v) {
			return false
		}
	}
	return true
}

// None reports whether no element satisfies p. Equivalent to !Any(p).
func (s Stream[T]) None(p func(T) bool) bool {
	return !s.Any(p)
}

// Find returns the first element for which p returns true, or the zero value
// and false if none match.
func (s Stream[T]) Find(p func(T) bool) (T, bool) {
	for v := range s.seq {
		if p(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Reduce folds the elements left-to-right using f, starting from init.
func (s Stream[T]) Reduce(init T, f func(acc, v T) T) T {
	acc := init
	for v := range s.seq {
		acc = f(acc, v)
	}
	return acc
}

// ForEach invokes visit on every element.
func (s Stream[T]) ForEach(visit func(T)) {
	for v := range s.seq {
		visit(v)
	}
}

// ---------- Sort and aggregation (comparator-provided variants) ----------

// SortBy returns a Stream that yields the elements ordered by less.
//
// The comparator is provided by the caller, so this method works for any T
// (not just T cmp.Ordered). A no-arg Sort() that uses cmp.Compare would
// require a tighter constraint on T, which a single Stream[T any] type
// cannot express on its method set.
func (s Stream[T]) SortBy(less func(x, y T) int) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		buf := s.Collect()
		slices.SortFunc(buf, less)
		for _, v := range buf {
			if !yield(v) {
				return
			}
		}
	}}
}

// MinBy returns the smallest element under less, or the zero value and false
// if the Stream is empty.
func (s Stream[T]) MinBy(less func(x, y T) int) (T, bool) {
	var (
		zero  T
		best  T
		found bool
	)
	for v := range s.seq {
		if !found {
			best = v
			found = true
			continue
		}
		if less(v, best) < 0 {
			best = v
		}
	}
	if !found {
		return zero, false
	}
	return best, true
}

// MaxBy returns the largest element under less, or the zero value and false
// if the Stream is empty.
func (s Stream[T]) MaxBy(less func(x, y T) int) (T, bool) {
	var (
		zero  T
		best  T
		found bool
	)
	for v := range s.seq {
		if !found {
			best = v
			found = true
			continue
		}
		if less(v, best) > 0 {
			best = v
		}
	}
	if !found {
		return zero, false
	}
	return best, true
}
