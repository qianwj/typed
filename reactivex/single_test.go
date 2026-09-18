package reactivex

import (
	"context"
	"errors"
	"github.com/qianwj/typed/adt/result"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var errSentinelSingle = errors.New("sentinel-single")

func TestSingle_Subscribe_Success(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[int] { return result.Success[int](42) })

	var got int
	var err error
	var called atomic.Bool
	sub := s.Subscribe(func(v int) {
		called.Store(true)
		got = v
	}, func(e error) {
		err = e
	})

	<-sub.Done()

	if !called.Load() {
		t.Fatal("onSuccess not called")
	}
	if got != 42 {
		t.Errorf("got = %d, want 42", got)
	}
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
}

func TestSingle_Subscribe_Error(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })

	var successCalled atomic.Bool
	var gotErr error
	sub := s.Subscribe(func(int) {
		successCalled.Store(true)
	}, func(e error) {
		gotErr = e
	})

	<-sub.Done()

	if successCalled.Load() {
		t.Error("onSuccess called for an error Single")
	}
	if !errors.Is(gotErr, errSentinelSingle) {
		t.Errorf("gotErr = %v, want sentinel", gotErr)
	}
}

func TestSingle_Subscribe_NilCallbacks(t *testing.T) {
	t.Parallel()
	// Subscribing with nil callbacks should not panic on either
	// branch.
	s := NewSingle(func() result.Result[int] { return result.Success[int](7) })
	sub := s.Subscribe(nil, nil)
	<-sub.Done()

	s2 := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })
	sub2 := s2.Subscribe(nil, nil)
	<-sub2.Done()
}

func TestSingle_Await(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[string] { return result.Success[string]("hello") })

	v, err := s.Await().Unwrap()
	if err != nil || v != "hello" {
		t.Errorf("Await = (%q, %v), want (\"hello\", nil)", v, err)
	}

	// Subsequent Await returns the same cached result without
	// re-running fn.
	v2, err2 := s.Await().Unwrap()
	if err2 != nil || v2 != "hello" {
		t.Errorf("Await #2 = (%q, %v), want (\"hello\", nil)", v2, err2)
	}
	if !s.Done() {
		t.Error("Done() = false after Await, want true")
	}
}

func TestSingle_AwaitWithContext_Cancel(t *testing.T) {
	t.Parallel()

	// Producer never finishes on its own; AwaitWithContext must
	// return ctx.Err() without invoking fn.
	started := make(chan struct{})
	s := NewSingle(func() result.Result[int] {
		close(started)
		return result.Success[int](0)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	v, err := s.AwaitWithContext(ctx).Unwrap()
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if v != 0 {
		t.Errorf("v = %d, want 0", v)
	}

	// Producer should not have started because we returned early.
	select {
	case <-started:
		t.Error("fn was invoked even though AwaitWithContext returned ctx.Err()")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestSingle_AwaitWithContext_Done(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[int] { return result.Success[int](99) })

	v, err := s.AwaitWithContext(context.Background()).Unwrap()
	if err != nil || v != 99 {
		t.Errorf("AwaitWithContext = (%d, %v), want (99, nil)", v, err)
	}
}

func TestSingle_MultipleSubscribers_ShareResult(t *testing.T) {
	t.Parallel()

	// fn runs exactly once even with multiple subscribers.
	var calls atomic.Int32
	s := NewSingle(func() result.Result[string] {
		calls.Add(1)
		return result.Success[string]("shared")
	})

	var wg sync.WaitGroup
	wg.Add(3)
	for range 3 {
		go func() {
			defer wg.Done()
			sub := s.Subscribe(func(v string) {}, func(error) {})
			<-sub.Done()
		}()
	}
	wg.Wait()

	if got := calls.Load(); got != 1 {
		t.Errorf("fn calls = %d, want 1", got)
	}
}

func TestSingle_Map_Success(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[int] { return result.Success[int](21) }).Map(func(n int) string {
		return "x" + string(rune('0'+n%10))
	})

	v, err := s.Await().Unwrap()
	if err != nil || v != "x1" {
		t.Errorf("Map Await = (%q, %v), want (\"x1\", nil)", v, err)
	}
}

func TestSingle_Map_ErrorPassesThrough(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })
	mapped := s.Map(func(n int) string { return "should not run" })

	v, err := mapped.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
	if v != "" {
		t.Errorf("v = %q, want \"\"", v)
	}
}

func TestSingle_FlatMap_Success(t *testing.T) {
	t.Parallel()

	inner := NewSingle(func() result.Result[string] { return result.Success[string]("inner") })
	outer := NewSingle(func() result.Result[int] { return result.Success[int](5) }).
		FlatMap(func(n int) Single[string] {
			return inner
		})

	v, err := outer.Await().Unwrap()
	if err != nil || v != "inner" {
		t.Errorf("FlatMap Await = (%q, %v), want (\"inner\", nil)", v, err)
	}
}

func TestSingle_FlatMap_OuterError(t *testing.T) {
	t.Parallel()

	outer := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) }).
		FlatMap(func(int) Single[string] {
			t.Error("FlatMap f should not run on outer error")
			return Single[string]{}
		})

	_, err := outer.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_Zip_BothSuccess(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() result.Result[int] { return result.Success[int](3) })
	b := NewSingle(func() result.Result[int] { return result.Success[int](4) })
	zipped := a.Zip(b, func(x, y int) int { return x*x + y*y })

	v, err := zipped.Await().Unwrap()
	if err != nil || v != 25 {
		t.Errorf("Zip Await = (%d, %v), want (25, nil)", v, err)
	}
}

func TestSingle_Zip_LeftError(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })
	b := NewSingle(func() result.Result[int] { return result.Success[int](4) })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	_, err := zipped.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_Zip_RightError(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() result.Result[int] { return result.Success[int](3) })
	b := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	_, err := zipped.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_AndThen_RunsNext(t *testing.T) {
	t.Parallel()

	first := NewSingle(func() result.Result[int] { return result.Success[int](1) })
	second := NewSingle(func() result.Result[string] { return result.Success[string]("next") })
	seq := first.AndThen(second)

	v, err := seq.Await().Unwrap()
	if err != nil || v != "next" {
		t.Errorf("AndThen Await = (%q, %v), want (\"next\", nil)", v, err)
	}
}

func TestSingle_AndThen_SkipsNextOnError(t *testing.T) {
	t.Parallel()

	first := NewSingle(func() result.Result[int] { return result.Failure[int](errSentinelSingle) })
	second := NewSingle(func() result.Result[string] {
		t.Error("second should not run when first errors")
		return result.Success[string]("")
	})
	seq := first.AndThen(second)

	_, err := seq.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_SubscribeThenAwait(t *testing.T) {
	t.Parallel()

	// Mix reactive + blocking consumers on the same Single.
	s := NewSingle(func() result.Result[int] { return result.Success[int](7) })

	var reactive int
	sub := s.Subscribe(func(v int) { reactive = v }, nil)
	<-sub.Done()

	v, err := s.Await().Unwrap()
	if err != nil || v != 7 || reactive != 7 {
		t.Errorf("got reactive=%d blocking=(%d, %v), want 7 / (7, nil)", reactive, v, err)
	}
}

func TestSingle_Done_InitialFalseThenTrue(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	gate := make(chan struct{})
	s := NewSingle(func() result.Result[int] {
		close(started)
		<-gate
		return result.Success[int](11)
	})

	if s.Done() {
		t.Error("Done() = true before any Await, want false")
	}

	done := make(chan struct{})
	go func() {
		v, err := s.Await().Unwrap()
		if v != 11 || err != nil {
			t.Errorf("Await = (%d, %v), want (11, nil)", v, err)
		}
		close(done)
	}()

	<-started
	if s.Done() {
		t.Error("Done() = true while fn is still running, want false")
	}

	close(gate)
	<-done
	if !s.Done() {
		t.Error("Done() = false after Await returned, want true")
	}
}
