package reactivex

import (
	"errors"
	"runtime"
	"testing"

	"github.com/qianwj/typed/adt/option"
	"github.com/qianwj/typed/adt/result"
)

func TestCompositionDoesNotParkGoroutinePerOperator(t *testing.T) {
	// Keep this test serial so other parallel tests cannot inflate the count.
	for _, kind := range []string{"Single", "Maybe"} {
		t.Run(kind, func(t *testing.T) {
			started, release := make(chan struct{}), make(chan struct{})
			produce := func() int { close(started); <-release; return 0 }
			const pairs = 256
			var subscribe func() Subscription
			var value func() int
			switch kind {
			case "Single":
				chain := NewSingle(func() result.Result[int] { return result.Success(produce()) })
				for range pairs {
					chain = chain.Map(func(v int) int { return v + 1 }).FlatMap(func(v int) Single[int] {
						return NewSingle(func() result.Result[int] { return result.Success(v + 1) })
					})
				}
				subscribe = func() Subscription { return chain.Subscribe(nil, nil) }
				value = func() int { return chain.Await().Value() }
			case "Maybe":
				chain := NewMaybe(func() result.Result[option.Option[int]] { return result.Success(option.Of(produce())) })
				for range pairs {
					chain = chain.Map(func(v int) int { return v + 1 }).FlatMap(func(v int) Maybe[int] {
						return NewMaybe(func() result.Result[option.Option[int]] { return result.Success(option.Of(v + 1)) })
					})
				}
				subscribe = func() Subscription { return chain.Subscribe(nil, nil, nil) }
				value = func() int { return chain.Await().Value().Get() }
			}
			before := runtime.NumGoroutine()
			sub := subscribe()
			waitCompletion(t, started)
			growth := runtime.NumGoroutine() - before
			close(release)
			waitCompletion(t, sub.Done())
			if got := value(); got != 2*pairs {
				t.Fatalf("chain result = %d, want %d", got, 2*pairs)
			}
			// Leave room for incidental runtime work, but reject a waiter per stage.
			if growth > 16 {
				t.Fatalf("%d operators added %d goroutines while waiting on one source", 2*pairs, growth)
			}
			t.Logf("%d operators: %d additional goroutines while source is blocked", 2*pairs, growth)
		})
	}
}

func TestFlatMapCachedInner(t *testing.T) {
	t.Parallel()
	t.Run("Single", func(t *testing.T) {
		for _, want := range []result.Result[int]{result.Success(42), result.Failure[int](errSentinelSingle)} {
			inner := NewSingle(func() result.Result[int] { return want })
			inner.Await()
			outer := NewSingle(func() result.Result[int] { return result.Success(1) })
			outer.Await()
			chain := outer.FlatMap(func(int) Single[int] { return inner })
			waitCompletion(t, chain.Subscribe(nil, nil).Done())
			got, err := chain.Await().Unwrap()
			value, wantErr := want.Unwrap()
			if got != value || !errors.Is(err, wantErr) {
				t.Fatalf("got (%v, %v), want (%v, %v)", got, err, value, wantErr)
			}
		}
	})
	t.Run("Maybe", func(t *testing.T) {
		for _, want := range []result.Result[option.Option[int]]{
			result.Success(option.Of(42)), result.Success(option.Empty[int]()), result.Failure[option.Option[int]](errSentinelMaybe),
		} {
			inner := NewMaybe(func() result.Result[option.Option[int]] { return want })
			inner.Await()
			outer := NewMaybe(func() result.Result[option.Option[int]] { return result.Success(option.Of(1)) })
			outer.Await()
			chain := outer.FlatMap(func(int) Maybe[int] { return inner })
			waitCompletion(t, chain.Subscribe(nil, nil, nil).Done())
			got, err := chain.Await().Unwrap()
			value, wantErr := want.Unwrap()
			if got != value || !errors.Is(err, wantErr) {
				t.Fatalf("got (%v, %v), want (%v, %v)", got, err, value, wantErr)
			}
		}
	})
}

func TestFlatMapInnerSourceDoesNotBlockUpstreamSubscribers(t *testing.T) {
	t.Parallel()
	started, release := make(chan struct{}), make(chan struct{})
	innerStarted, innerRelease := make(chan struct{}), make(chan struct{})
	defer close(innerRelease)
	source := NewSingle(func() result.Result[int] {
		close(started)
		<-release
		return result.Success(1)
	})
	inner := NewSingle(func() result.Result[int] {
		close(innerStarted)
		<-innerRelease
		return result.Success(2)
	})
	chain := source.FlatMap(func(int) Single[int] { return inner })
	chain.Subscribe(nil, nil)
	waitCompletion(t, started)
	sub := source.Subscribe(func(int) {}, nil)
	close(release)
	waitCompletion(t, innerStarted)
	waitCompletion(t, sub.Done())
}
