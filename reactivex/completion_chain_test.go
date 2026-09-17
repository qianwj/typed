package reactivex

import (
	"errors"
	"runtime"
	"testing"

	"github.com/qianwj/typed/adt"
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
				chain := NewSingle(func() adt.Result[int] { return adt.Success(produce()) })
				for range pairs {
					chain = chain.Map(func(v int) int { return v + 1 }).FlatMap(func(v int) Single[int] {
						return NewSingle(func() adt.Result[int] { return adt.Success(v + 1) })
					})
				}
				subscribe = func() Subscription { return chain.Subscribe(nil, nil) }
				value = func() int { return chain.Await().Value() }
			case "Maybe":
				chain := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of(produce())) })
				for range pairs {
					chain = chain.Map(func(v int) int { return v + 1 }).FlatMap(func(v int) Maybe[int] {
						return NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of(v + 1)) })
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
		for _, want := range []adt.Result[int]{adt.Success(42), adt.Failure[int](errSentinelSingle)} {
			inner := NewSingle(func() adt.Result[int] { return want })
			inner.Await()
			outer := NewSingle(func() adt.Result[int] { return adt.Success(1) })
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
		for _, want := range []adt.Result[adt.Option[int]]{
			adt.Success(adt.Of(42)), adt.Success(adt.Empty[int]()), adt.Failure[adt.Option[int]](errSentinelMaybe),
		} {
			inner := NewMaybe(func() adt.Result[adt.Option[int]] { return want })
			inner.Await()
			outer := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of(1)) })
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
	source := NewSingle(func() adt.Result[int] {
		close(started)
		<-release
		return adt.Success(1)
	})
	inner := NewSingle(func() adt.Result[int] {
		close(innerStarted)
		<-innerRelease
		return adt.Success(2)
	})
	chain := source.FlatMap(func(int) Single[int] { return inner })
	chain.Subscribe(nil, nil)
	waitCompletion(t, started)
	sub := source.Subscribe(func(int) {}, nil)
	close(release)
	waitCompletion(t, innerStarted)
	waitCompletion(t, sub.Done())
}
