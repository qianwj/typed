package reactivex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFromSliceMapFilterTake(t *testing.T) {
	var got []int
	done := make(chan struct{})
	FromSlice([]int{1, 2, 3, 4}).Map(func(_ context.Context, v int) (int, error) { return v * 2, nil }).Filter(func(v int) bool { return v > 2 }).Take(2).ForEach(context.Background(), func(v int) { got = append(got, v) }, func(error) { t.Fail() }, func() { close(done) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
	if len(got) != 2 || got[0] != 4 || got[1] != 6 {
		t.Fatalf("got %v", got)
	}
}

func TestMapError(t *testing.T) {
	want := errors.New("bad")
	got := make(chan error, 1)
	FromSlice([]int{1}).Map(func(context.Context, int) (int, error) { return 0, want }).ForEach(context.Background(), nil, func(err error) { got <- err }, nil)
	if err := <-got; !errors.Is(err, want) {
		t.Fatalf("got %v", err)
	}
}
