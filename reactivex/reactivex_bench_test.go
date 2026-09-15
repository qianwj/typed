package reactivex_test

import (
	"context"
	"testing"

	"github.com/qianwj/typed/reactivex"
)

// BenchmarkObservable_JustToSlice measures the round-trip throughput
// of a finite cold source: subscribe, request everything, collect.
// This is the "every-value-matters" baseline for the toolkit.
func BenchmarkObservable_JustToSlice(b *testing.B) {
	const n = 1024
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		values, err := reactivex.Just(seed...).ToSlice(ctx)
		if err != nil || len(values) != n {
			b.Fatalf("got %d values, err=%v", len(values), err)
		}
	}
}

// BenchmarkObservable_FilterMapToSlice measures the wrapping-style
// pipeline cost: each operator is a Subscribe adapter that forwards
// notifications without adding queues or goroutines. The total
// per-element cost is the sum of operator dispatch overhead.
func BenchmarkObservable_FilterMapToSlice(b *testing.B) {
	const n = 1024
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := reactivex.Just(seed...).
			Filter(func(v int) bool { return v%2 == 0 }).
			Map(func(_ context.Context, v int) (int, error) { return v * 3, nil }).
			ToSlice(ctx)
		if err != nil {
			b.Fatalf("err=%v", err)
		}
	}
}

// BenchmarkObservable_TakeEarlyTermination exercises the early-termination
// path of Take. Even though the source is large, the collector
// should not pull more than the requested count.
func BenchmarkObservable_TakeEarlyTermination(b *testing.B) {
	const n = 1 << 16 // large source
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		values, err := reactivex.Just(seed...).Take(10).ToSlice(ctx)
		if err != nil || len(values) != 10 {
			b.Fatalf("got %d values, err=%v", len(values), err)
		}
	}
}

// BenchmarkObservable_ForEachNoCallbacks measures the overhead of the
// callback-style subscriber with all callbacks nil. This is the
// lowest-overhead way to drive a cold source to completion.
func BenchmarkObservable_ForEachNoCallbacks(b *testing.B) {
	const n = 256
	seed := make([]int, n)
	for i := range seed {
		seed[i] = i
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reactivex.Just(seed...).ForEach(ctx, nil, nil, nil)
	}
}
