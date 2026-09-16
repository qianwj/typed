package reactivex

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var errSentinelSingle = errors.New("sentinel-single")

func TestSingle_Subscribe_Success(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() (int, error) { return 42, nil })

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

	s := NewSingle(func() (int, error) { return 0, errSentinelSingle })

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
	s := NewSingle(func() (int, error) { return 7, nil })
	sub := s.Subscribe(nil, nil)
	<-sub.Done()

	s2 := NewSingle(func() (int, error) { return 0, errSentinelSingle })
	sub2 := s2.Subscribe(nil, nil)
	<-sub2.Done()
}

func TestSingle_Await(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() (string, error) { return "hello", nil })

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
	s := NewSingle(func() (int, error) {
		close(started)
		return 0, nil
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

	s := NewSingle(func() (int, error) { return 99, nil })

	v, err := s.AwaitWithContext(context.Background()).Unwrap()
	if err != nil || v != 99 {
		t.Errorf("AwaitWithContext = (%d, %v), want (99, nil)", v, err)
	}
}

func TestSingle_MultipleSubscribers_ShareResult(t *testing.T) {
	t.Parallel()

	// fn runs exactly once even with multiple subscribers.
	var calls atomic.Int32
	s := NewSingle(func() (string, error) {
		calls.Add(1)
		return "shared", nil
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

	s := NewSingle(func() (int, error) { return 21, nil }).Map(func(n int) string {
		return "x" + string(rune('0'+n%10))
	})

	v, err := s.Await().Unwrap()
	if err != nil || v != "x1" {
		t.Errorf("Map Await = (%q, %v), want (\"x1\", nil)", v, err)
	}
}

func TestSingle_Map_ErrorPassesThrough(t *testing.T) {
	t.Parallel()

	s := NewSingle(func() (int, error) { return 0, errSentinelSingle })
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

	inner := NewSingle(func() (string, error) { return "inner", nil })
	outer := NewSingle(func() (int, error) { return 5, nil }).
		FlatMap(func(n int) *Single[string] {
			return inner
		})

	v, err := outer.Await().Unwrap()
	if err != nil || v != "inner" {
		t.Errorf("FlatMap Await = (%q, %v), want (\"inner\", nil)", v, err)
	}
}

func TestSingle_FlatMap_OuterError(t *testing.T) {
	t.Parallel()

	outer := NewSingle(func() (int, error) { return 0, errSentinelSingle }).
		FlatMap(func(int) *Single[string] {
			t.Error("FlatMap f should not run on outer error")
			return nil
		})

	_, err := outer.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_Zip_BothSuccess(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() (int, error) { return 3, nil })
	b := NewSingle(func() (int, error) { return 4, nil })
	zipped := a.Zip(b, func(x, y int) int { return x*x + y*y })

	v, err := zipped.Await().Unwrap()
	if err != nil || v != 25 {
		t.Errorf("Zip Await = (%d, %v), want (25, nil)", v, err)
	}
}

func TestSingle_Zip_LeftError(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() (int, error) { return 0, errSentinelSingle })
	b := NewSingle(func() (int, error) { return 4, nil })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	_, err := zipped.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_Zip_RightError(t *testing.T) {
	t.Parallel()

	a := NewSingle(func() (int, error) { return 3, nil })
	b := NewSingle(func() (int, error) { return 0, errSentinelSingle })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	_, err := zipped.Await().Unwrap()
	if !errors.Is(err, errSentinelSingle) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestSingle_AndThen_RunsNext(t *testing.T) {
	t.Parallel()

	first := NewSingle(func() (int, error) { return 1, nil })
	second := NewSingle(func() (string, error) { return "next", nil })
	seq := first.AndThen(second)

	v, err := seq.Await().Unwrap()
	if err != nil || v != "next" {
		t.Errorf("AndThen Await = (%q, %v), want (\"next\", nil)", v, err)
	}
}

func TestSingle_AndThen_SkipsNextOnError(t *testing.T) {
	t.Parallel()

	first := NewSingle(func() (int, error) { return 0, errSentinelSingle })
	second := NewSingle(func() (string, error) {
		t.Error("second should not run when first errors")
		return "", nil
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
	s := NewSingle(func() (int, error) { return 7, nil })

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
	s := NewSingle(func() (int, error) {
		close(started)
		<-gate
		return 11, nil
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
