package stream

import "iter"

// MapStream is a lazy, single-use pipeline of key-value pairs, constructed
// with FromPairs or Stream.Associate. Filter and MapValues preserve its map shape;
// Collect materializes it into a plain map[K]V.
//
// Pairs follow source order and may contain duplicate keys. Transforms
// process each pair; duplicate keys are resolved only by Collect.
// The zero value is not usable.
type MapStream[K comparable, V any] struct {
	seq iter.Seq2[K, V]
}

// FromPairs builds a MapStream from a key-value iterator without consuming it
// or copying its source. The iterator controls pair order and must stop when
// yield returns false.
func FromPairs[K comparable, V any](seq iter.Seq2[K, V]) MapStream[K, V] {
	return MapStream[K, V]{seq: seq}
}

// Filter keeps the key-value pairs for which p returns true.
func (s MapStream[K, V]) Filter(p func(K, V) bool) MapStream[K, V] {
	return MapStream[K, V]{seq: func(yield func(K, V) bool) {
		for k, v := range s.seq {
			if p(k, v) && !yield(k, v) {
				return
			}
		}
	}}
}

// MapValues transforms each value while preserving its key.
func (s MapStream[K, V]) MapValues[R any](f func(K, V) R) MapStream[K, R] {
	return MapStream[K, R]{seq: func(yield func(K, R) bool) {
		for k, v := range s.seq {
			if !yield(k, f(k, v)) {
				return
			}
		}
	}}
}

// Map projects each pair into an ordinary Stream whose Collect returns a slice.
// Use MapValues to preserve keys and collect into a map.
func (s MapStream[K, V]) Map[R any](f func(K, V) R) Stream[R] {
	return From(func(yield func(R) bool) {
		for k, v := range s.seq {
			if !yield(f(k, v)) {
				return
			}
		}
	})
}

// Collect consumes the stream into a new map. For duplicate keys, the last
// value reaching Collect wins. An empty stream returns a non-nil empty map.
func (s MapStream[K, V]) Collect() map[K]V {
	out := make(map[K]V)
	for k, v := range s.seq {
		out[k] = v
	}
	return out
}
