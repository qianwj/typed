package reactivex_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/adt/option"
	"github.com/qianwj/typed/adt/result"
	"github.com/qianwj/typed/reactivex"
)

func TestSingleCopiesShareComputation(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("single source failed")
	for _, want := range []result.Result[int]{result.Success(42), result.Failure[int](wantErr)} {
		wantValue, wantError := want.Unwrap()
		var calls atomic.Int32
		source := reactivex.NewSingle(func() result.Result[int] {
			calls.Add(1)
			return want
		})
		copied := source
		handles := map[int]reactivex.Single[int]{0: source, 1: copied}
		if source.Done() || copied.Done() || calls.Load() != 0 {
			t.Fatal("copying started the source")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var workers sync.WaitGroup
		for i := range 32 {
			workers.Go(func() {
				// A map entry is not addressable: this also exercises the value method set.
				got := handles[i%2].AwaitWithContext(ctx)
				if value, err := got.Unwrap(); value != wantValue || !errors.Is(err, wantError) {
					t.Errorf("copied handle result = %v, want %v", got, want)
				}
			})
		}
		workers.Wait()
		cancel()
		if !source.Done() || !copied.Done() || calls.Load() != 1 {
			t.Fatalf("copies did not share completion: source done=%v, copy done=%v, calls=%d",
				source.Done(), copied.Done(), calls.Load())
		}
		if value, err := copied.Await().Unwrap(); value != wantValue || !errors.Is(err, wantError) {
			t.Errorf("cached result = (%v, %v), want %v", value, err, want)
		}
		if calls.Load() != 1 {
			t.Fatal("reading the cached result restarted the source")
		}
	}
}

func TestMaybeCopiesShareComputation(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("maybe source failed")
	for _, want := range []result.Result[option.Option[int]]{
		result.Success(option.Of(42)),
		result.Success(option.Empty[int]()),
		result.Failure[option.Option[int]](wantErr),
	} {
		wantValue, wantError := want.Unwrap()
		var calls atomic.Int32
		source := reactivex.NewMaybe(func() result.Result[option.Option[int]] {
			calls.Add(1)
			return want
		})
		copied := source
		handles := map[int]reactivex.Maybe[int]{0: source, 1: copied}
		if source.Done() || copied.Done() || calls.Load() != 0 {
			t.Fatal("copying started the source")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var workers sync.WaitGroup
		for i := range 32 {
			workers.Go(func() {
				got := handles[i%2].AwaitWithContext(ctx)
				if value, err := got.Unwrap(); value != wantValue || !errors.Is(err, wantError) {
					t.Errorf("copied handle result = %v, want %v", got, want)
				}
			})
		}
		workers.Wait()
		cancel()
		if !source.Done() || !copied.Done() || calls.Load() != 1 {
			t.Fatalf("copies did not share completion: source done=%v, copy done=%v, calls=%d",
				source.Done(), copied.Done(), calls.Load())
		}
		if value, err := copied.Await().Unwrap(); value != wantValue || !errors.Is(err, wantError) {
			t.Errorf("cached result = (%v, %v), want %v", value, err, want)
		}
		if calls.Load() != 1 {
			t.Fatal("reading the cached result restarted the source")
		}
	}
}
