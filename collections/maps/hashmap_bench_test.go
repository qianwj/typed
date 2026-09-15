package maps_test

import (
	"testing"

	"github.com/qianwj/typed/collections/maps"
)

// BenchmarkHashMap_Put measures the put path on a fresh HashMap.
// The per-op cost includes the hash computation and the open-addressed
// probe; with b.N small (no resize triggered repeatedly) it should
// be roughly constant.
func BenchmarkHashMap_Put(b *testing.B) {
	m := maps.NewHashMap[int, int]()
	for i := 0; i < b.N; i++ {
		m.Put(i, i*10)
	}
}

// BenchmarkHashMap_Get measures the get path. The probe length
// stays bounded by the load factor, so this should be O(1).
func BenchmarkHashMap_Get(b *testing.B) {
	const n = 4096
	m := maps.NewHashMap[int, int]()
	for i := 0; i < n; i++ {
		m.Put(i, i*10)
	}
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		if v, ok := m.Get(i % n); ok {
			sink += v
		}
	}
	_ = sink
}

// BenchmarkHashMap_HitMiss measures both Get outcomes: a present
// key (longest probe) and a missing key (full probe to a tombstone
// or empty slot). The two are reported separately because the
// miss path is the upper bound on cost.
func BenchmarkHashMap_HitMiss(b *testing.B) {
	const n = 4096
	m := maps.NewHashMap[int, int]()
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Hit: keys 0..n/2 are guaranteed to exist.
		_, _ = m.Get(i % (n / 2))
		// Miss: keys far above the populated range.
		_, _ = m.Get(n + i)
	}
}

// BenchmarkHashMap_Remove measures the remove path, including the
// tombstone write that keeps the open-addressed probe sequence
// well-defined for subsequent inserts.
func BenchmarkHashMap_Remove(b *testing.B) {
	const n = 4096
	m := maps.NewHashMap[int, int]()
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(i % n)
		m.Put(i%n, i) // reinsert so the next iteration has something to remove
	}
}
