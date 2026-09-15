package stream_test

import (
	"testing"

	"github.com/qianwj/typed/collections/lists"
)

// BenchmarkStream_FilterMapCollect measures the lazy pipeline
// Filter -> Map -> Collect end-to-end. The lazy layer adds an
// indirection per element compared to the eager ArrayList pipeline;
// the ratio between the two is the cost of laziness on this workload.
func BenchmarkStream_FilterMapCollect(b *testing.B) {
	const n = 1024
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lists.ArrayListOf(seed...).
			Stream().
			Filter(func(v int) bool { return v%2 == 0 }).
			Map(func(v int) int { return v * 3 }).
			Collect()
	}
}

// BenchmarkStream_TakeEarlyTermination exercises the early-termination
// behaviour of Stream. The source exposes a finite sequence; after
// Take(10) is satisfied, the rest of the source must not be pulled.
func BenchmarkStream_TakeEarlyTermination(b *testing.B) {
	const n = 1 << 16 // large source
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lists.ArrayListOf(seed...).
			Stream().
			Take(10).
			Collect()
	}
}
