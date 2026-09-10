package sets

import (
	"slices"

	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/maps"
	"github.com/qianwj/typed/collections/stream"
)

// HashSet is an unordered collection of unique comparable values.
//
// HashSet is backed by a HashMap[T, struct{}]. Its iteration order is
// unspecified, matching Go's map iteration semantics. Methods that return a
// new HashSet leave the receiver unchanged.
type HashSet[T comparable] struct {
	container *maps.HashMap[T, struct{}]
}

// NewHashSet returns a new empty HashSet.
func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{container: maps.NewHashMap[T, struct{}]()}
}

// HashSetOf returns a HashSet containing values. Duplicate values are
// retained only once.
func HashSetOf[T comparable](values ...T) *HashSet[T] {
	set := NewHashSet[T]()
	for _, value := range values {
		set.Add(value)
	}
	return set
}

func (s *HashSet[T]) ensureContainer() {
	if s.container == nil {
		s.container = maps.NewHashMap[T, struct{}]()
	}
}

// Add inserts value. Adding a value already in the set has no effect.
func (s *HashSet[T]) Add(value T) {
	if s == nil {
		return
	}
	s.ensureContainer()
	s.container.Put(value, struct{}{})
}

// Remove deletes value. Removing a value that is not present has no effect.
func (s *HashSet[T]) Remove(value T) {
	if s == nil || s.container == nil {
		return
	}
	s.container.Remove(value)
}

// Contains reports whether value is in the set.
func (s *HashSet[T]) Contains(value T) bool {
	if s == nil || s.container == nil {
		return false
	}
	return s.container.Contains(value)
}

// Size returns the number of unique values.
func (s *HashSet[T]) Size() int {
	if s == nil || s.container == nil {
		return 0
	}
	return s.container.Size()
}

// IsEmpty reports whether the set contains no values.
func (s *HashSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear removes all values from the set.
func (s *HashSet[T]) Clear() {
	if s == nil || s.container == nil {
		return
	}
	s.container.Clear()
}

// ForEach invokes visit once for every value. Iteration order is unspecified.
func (s *HashSet[T]) ForEach(visit func(T)) {
	if s == nil || s.container == nil {
		return
	}
	s.container.ForEach(func(value T, _ struct{}) {
		visit(value)
	})
}

// Collect returns the values as a freshly allocated slice. The order is
// unspecified and the returned slice is independent of the set.
func (s *HashSet[T]) Collect() []T {
	values := make([]T, 0, s.Size())
	s.ForEach(func(value T) {
		values = append(values, value)
	})
	return values
}

// Stream returns a single-use stream over a snapshot of the set at the time
// Stream is called. Mutations to the set after this call do not affect it.
func (s *HashSet[T]) Stream() stream.Stream[T] {
	return stream.FromSlice(s.Collect())
}

// Peek calls visit on every value and returns an independent copy of the set.
func (s *HashSet[T]) Peek(visit func(T)) *HashSet[T] {
	out := NewHashSet[T]()
	s.ForEach(func(value T) {
		visit(value)
		out.Add(value)
	})
	return out
}

// Filter returns a new HashSet containing values for which p returns true.
func (s *HashSet[T]) Filter(p func(T) bool) *HashSet[T] {
	out := NewHashSet[T]()
	s.ForEach(func(value T) {
		if p(value) {
			out.Add(value)
		}
	})
	return out
}

// Map applies f to every value and returns an ArrayList of the results.
//
// The mapped type is unconstrained: it may be a slice, map, function, or any
// other non-comparable type. Results are kept one-for-one, so equal mapped
// values are not implicitly deduplicated.
func (s *HashSet[T]) Map[R any](f func(T) R) *lists.ArrayList[R] {
	values := make([]R, 0, s.Size())
	s.ForEach(func(value T) {
		values = append(values, f(value))
	})
	return lists.ArrayListOf(values...)
}

// FlatMap applies f to every value and concatenates the returned ArrayLists.
//
// The result type is unconstrained for the same reason as Map: a flattened
// value may be non-comparable, so FlatMap returns an ArrayList rather than a
// HashSet.
func (s *HashSet[T]) FlatMap[R any](f func(T) *lists.ArrayList[R]) *lists.ArrayList[R] {
	var values []R
	s.ForEach(func(value T) {
		mapped := f(value)
		if mapped == nil {
			return
		}
		values = append(values, mapped.Collect()...)
	})
	return lists.ArrayListOf(values...)
}

// MapSet applies f to every value and returns a new HashSet. The mapped type
// must be comparable because the result is stored as set members; duplicate
// mapped values are deduplicated.
func (s *HashSet[T]) MapSet[R comparable](f func(T) R) *HashSet[R] {
	out := NewHashSet[R]()
	s.ForEach(func(value T) {
		out.Add(f(value))
	})
	return out
}

// FlatMapSet applies f to every value and returns the union of the returned
// HashSets. The mapped type must be comparable because the result is stored
// as set members.
func (s *HashSet[T]) FlatMapSet[R comparable](f func(T) *HashSet[R]) *HashSet[R] {
	out := NewHashSet[R]()
	s.ForEach(func(value T) {
		mapped := f(value)
		if mapped == nil {
			return
		}
		mapped.ForEach(out.Add)
	})
	return out
}

// Reduce folds the values in unspecified iteration order using f, starting
// from init. The accumulator type may differ from the set's element type.
func (s *HashSet[T]) Reduce[R any](init R, f func(acc R, value T) R) R {
	acc := init
	s.ForEach(func(value T) {
		acc = f(acc, value)
	})
	return acc
}

// Concat returns the union of s and other. Neither input is modified.
func (s *HashSet[T]) Concat(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	s.ForEach(out.Add)
	if other != nil {
		other.ForEach(out.Add)
	}
	return out
}

// Union is an alias for Concat.
func (s *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
	return s.Concat(other)
}

// Intersect returns the values present in both sets. Neither input is
// modified.
func (s *HashSet[T]) Intersect(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	if other == nil {
		return out
	}
	s.ForEach(func(value T) {
		if other.Contains(value) {
			out.Add(value)
		}
	})
	return out
}

// Difference returns the values present in s but absent from other. Neither
// input is modified.
func (s *HashSet[T]) Difference(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	s.ForEach(func(value T) {
		if other == nil || !other.Contains(value) {
			out.Add(value)
		}
	})
	return out
}

// SymmetricDifference returns the values present in exactly one of the two
// sets. Neither input is modified.
func (s *HashSet[T]) SymmetricDifference(other *HashSet[T]) *HashSet[T] {
	out := s.Difference(other)
	if other == nil {
		return out
	}
	other.ForEach(func(value T) {
		if !s.Contains(value) {
			out.Add(value)
		}
	})
	return out
}

// IsSubsetOf reports whether every value in s is also in other.
func (s *HashSet[T]) IsSubsetOf(other *HashSet[T]) bool {
	if other == nil {
		return s.Size() == 0
	}
	if s.Size() > other.Size() {
		return false
	}
	result := true
	s.ForEach(func(value T) {
		if result && !other.Contains(value) {
			result = false
		}
	})
	return result
}

// IsSupersetOf reports whether s contains every value in other.
func (s *HashSet[T]) IsSupersetOf(other *HashSet[T]) bool {
	if other == nil {
		return true
	}
	return other.IsSubsetOf(s)
}

// SortBy returns an ArrayList containing all values ordered by less.
// Sorting an unordered set necessarily produces an ordered list.
func (s *HashSet[T]) SortBy(less func(x, y T) int) *lists.ArrayList[T] {
	values := s.Collect()
	slices.SortFunc(values, less)
	return lists.ArrayListOf(values...)
}

// MinBy returns the smallest value under less, or the zero value and false
// when the set is empty.
func (s *HashSet[T]) MinBy(less func(x, y T) int) (T, bool) {
	var (
		best  T
		found bool
	)
	s.ForEach(func(value T) {
		if !found || less(value, best) < 0 {
			best = value
			found = true
		}
	})
	return best, found
}

// MaxBy returns the largest value under less, or the zero value and false
// when the set is empty.
func (s *HashSet[T]) MaxBy(less func(x, y T) int) (T, bool) {
	var (
		best  T
		found bool
	)
	s.ForEach(func(value T) {
		if !found || less(value, best) > 0 {
			best = value
			found = true
		}
	})
	return best, found
}

// Any reports whether at least one value satisfies p.
func (s *HashSet[T]) Any(p func(T) bool) bool {
	return slices.ContainsFunc(s.Collect(), p)
}

// All reports whether every value satisfies p. It returns true for an empty
// set.
func (s *HashSet[T]) All(p func(T) bool) bool {
	for _, value := range s.Collect() {
		if !p(value) {
			return false
		}
	}
	return true
}

// None reports whether no value satisfies p.
func (s *HashSet[T]) None(p func(T) bool) bool {
	return !s.Any(p)
}

// Find returns an arbitrary matching value and true, or the zero value and
// false when no value satisfies p. Because HashSet is unordered, the
// matching value is not deterministic when multiple values satisfy p.
func (s *HashSet[T]) Find(p func(T) bool) (T, bool) {
	for _, value := range s.Collect() {
		if p(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}
