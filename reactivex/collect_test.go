package reactivex

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestToSliceAndCollect(t *testing.T) {
	values, err := FromSlice([]int{1, 2, 3}).ToSlice(context.Background())
	if err != nil || !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("values=%v err=%v", values, err)
	}
	total, err := Collect(context.Background(), FromSlice([]int{1, 2, 3}), 0, func(sum, value int) int { return sum + value })
	if err != nil || total != 6 {
		t.Fatalf("total=%d err=%v", total, err)
	}
}

// TestToSliceCancelRace is a regression test for the cancel-path data race
// in ToSlice (and therefore in Collect, which uses ToSlice internally).
//
// Before the fix, the producer's OnNext append raced with the cancel
// branch's return of the live slice, so `go test -race` would flag the
// library code rather than the test. After the fix, ToSlice returns a
// snapshot on the cancel path and the test passes under -race.
func TestToSliceCancelRace(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	src := Create(func(ctx context.Context, emit func(int) bool, _ func()) {
		for i := 0; ; i++ {
			if !emit(i) {
				return
			}
		}
	})
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	values, err := src.ToSlice(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ToSlice err = %v, want context.Canceled", err)
	}
	// The returned slice is a snapshot; ranging it must not trigger the
	// race detector. We don't assert a specific length (cancellation
	// timing is non-deterministic) but each element should equal its
	// index because the source emits 0, 1, 2, ... in order.
	for i, v := range values {
		if v != i {
			t.Fatalf("values[%d] = %d, want %d (torn snapshot)", i, v, i)
		}
	}
}

// TestCollectCancelRace exercises the same race through Collect (which
// builds on ToSlice) to make sure the fix covers both entry points.
func TestCollectCancelRace(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	src := Create(func(ctx context.Context, emit func(int) bool, _ func()) {
		for i := 0; ; i++ {
			if !emit(i) {
				return
			}
		}
	})
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	// We don't assert anything about total — the point is that the
	// underlying ToSlice must not race, and Collect reads its returned
	// slice on the success / cancel path.
	total, err := Collect(ctx, src, 0, func(acc, v int) int { return acc + v })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Collect err = %v, want context.Canceled", err)
	}
	_ = total
}
