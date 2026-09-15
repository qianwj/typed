package sets_test

import (
	"testing"

	"github.com/qianwj/typed/collections/sets"
)

// BenchmarkHashSet_Add measures the add path, including the
// duplicate-check that distinguishes a set from a bag.
func BenchmarkHashSet_Add(b *testing.B) {
	s := sets.NewHashSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

// BenchmarkHashSet_Contains measures membership testing, which is
// the operation HashSet optimises for.
func BenchmarkHashSet_Contains(b *testing.B) {
	const n = 4096
	s := sets.NewHashSet[int]()
	for i := 0; i < n; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Contains(i % (n * 2))
	}
}