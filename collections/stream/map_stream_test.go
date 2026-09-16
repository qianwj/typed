package stream_test

import (
	"slices"
	"testing"

	"github.com/qianwj/typed/collections/stream"
)

func TestMapStreamLazyTransforms(t *testing.T) {
	sourceCalls, mapCalls := 0, 0
	s := stream.Of(1, 2, 3).
		Associate(func(n int) (int, int) {
			sourceCalls++
			return n, n * 10
		}).
		Filter(func(k, v int) bool { return k > 1 && v < 30 }).
		MapValues(func(k, v int) []int {
			mapCalls++
			return []int{k, v}
		})
	if sourceCalls != 0 || mapCalls != 0 {
		t.Fatal("pipeline ran before Collect")
	}
	var got map[int][]int = s.Collect()
	if len(got) != 1 || !slices.Equal(got[2], []int{2, 20}) || sourceCalls != 3 || mapCalls != 1 {
		t.Fatalf("got=%v sourceCalls=%d mapCalls=%d", got, sourceCalls, mapCalls)
	}
}

func TestMapStreamResolvesDuplicatesAfterFiltering(t *testing.T) {
	got := stream.Of(1, 2, 3, 4).
		Associate(func(n int) (int, int) { return n % 2, n }).
		Filter(func(_, v int) bool { return v != 3 }).
		Collect()
	if len(got) != 2 || got[1] != 1 || got[0] != 4 {
		t.Fatalf("got %v, want map[0:4 1:1]", got)
	}
}

func TestMapStreamProjectionStopsUpstream(t *testing.T) {
	calls := 0
	got := stream.Of(1, 2, 3, 4).
		Associate(func(n int) (int, int) {
			calls++
			return n, n * 10
		}).
		Filter(func(k, _ int) bool { return k%2 == 0 }).
		MapValues(func(_, v int) int { return v + 1 }).
		Map(func(k, v int) int { return k + v }).
		Take(1).
		Collect()
	if !slices.Equal(got, []int{23}) || calls != 2 {
		t.Fatalf("got=%v calls=%d, want [23] and 2 calls", got, calls)
	}
}

func TestMapStreamEmptyCollect(t *testing.T) {
	got := stream.Empty[int]().
		Associate(func(n int) (int, int) {
			t.Fatal("Associate called on empty stream")
			return n, n
		}).
		Filter(func(int, int) bool {
			t.Fatal("Filter called on empty stream")
			return true
		}).
		MapValues(func(int, int) string {
			t.Fatal("MapValues called on empty stream")
			return ""
		}).Collect()
	if got == nil || len(got) != 0 {
		t.Fatalf("got %v, want non-nil empty map", got)
	}
}
