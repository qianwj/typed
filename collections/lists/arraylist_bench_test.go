package lists_test

import (
	"testing"

	"github.com/qianwj/typed/collections/lists"
)

// BenchmarkArrayList_Add measures tail-side amortised growth.
func BenchmarkArrayList_Add(b *testing.B) {
	var s lists.ArrayList[int]
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

// BenchmarkArrayList_AddFirst measures the head-side AddFirst path
// without RemoveFirst; the head-offset layout keeps the cost stable
// even as the high-water mark grows.
func BenchmarkArrayList_AddFirst(b *testing.B) {
	const cap = 1 << 16
	s := lists.ArrayListOf(make([]int, cap)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.AddFirst(i)
	}
}

// BenchmarkArrayList_HeadDrain measures the head-offset layout
// under steady-state drain: RemoveFirst on a long-running head-drained
// list must stay bounded by a constant rather than the high-water mark.
func BenchmarkArrayList_HeadDrain(b *testing.B) {
	const cap = 1 << 16
	s := lists.ArrayListOf(make([]int, cap)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.AddFirst(i)
		_ = s.RemoveFirst()
	}
}

// BenchmarkArrayList_Peek measures the callback-style iteration path
// that Peek exposes (the simplest full-collection traversal).
func BenchmarkArrayList_Peek(b *testing.B) {
	const n = 1024
	s := lists.ArrayListOf(make([]int, n)...)
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		s.Peek(func(v int) { sink += v })
	}
	_ = sink
}

// BenchmarkArrayList_FluentFilterMap measures the eager fluent
// pipeline Filter -> Map -> Collect on a non-trivial input.
func BenchmarkArrayList_FluentFilterMap(b *testing.B) {
	const n = 1024
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lists.ArrayListOf(seed...).
			Filter(func(v int) bool { return v%2 == 0 }).
			Map(func(v int) int { return v * 3 }).
			Collect()
	}
}