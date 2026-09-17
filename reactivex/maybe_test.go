package reactivex

import (
	"context"
	"errors"
	"github.com/qianwj/typed/adt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var errSentinelMaybe = errors.New("sentinel-maybe")

func TestMaybe_Subscribe_Success(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](7)) })

	var successCalled, completeCalled, errorCalled atomic.Bool
	var got int
	sub := m.Subscribe(
		func(v int) { successCalled.Store(true); got = v },
		func() { completeCalled.Store(true) },
		func(error) { errorCalled.Store(true) },
	)
	<-sub.Done()

	if !successCalled.Load() || completeCalled.Load() || errorCalled.Load() {
		t.Errorf("success=%v complete=%v error=%v, want true/false/false",
			successCalled.Load(), completeCalled.Load(), errorCalled.Load())
	}
	if got != 7 {
		t.Errorf("got = %d, want 7", got)
	}
}

func TestMaybe_Subscribe_Complete(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) })

	var successCalled, completeCalled, errorCalled atomic.Bool
	sub := m.Subscribe(
		func(int) { successCalled.Store(true) },
		func() { completeCalled.Store(true) },
		func(error) { errorCalled.Store(true) },
	)
	<-sub.Done()

	if successCalled.Load() || !completeCalled.Load() || errorCalled.Load() {
		t.Errorf("success=%v complete=%v error=%v, want false/true/false",
			successCalled.Load(), completeCalled.Load(), errorCalled.Load())
	}
}

func TestMaybe_Subscribe_Error(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) })

	var successCalled, completeCalled atomic.Bool
	var gotErr error
	sub := m.Subscribe(
		func(int) { successCalled.Store(true) },
		func() { completeCalled.Store(true) },
		func(e error) { gotErr = e },
	)
	<-sub.Done()

	if successCalled.Load() || completeCalled.Load() {
		t.Error("success or complete called for an error Maybe")
	}
	if !errors.Is(gotErr, errSentinelMaybe) {
		t.Errorf("gotErr = %v, want sentinel", gotErr)
	}
}

func TestMaybe_Subscribe_IgnoresValueOnFailure(t *testing.T) {
	t.Parallel()
	// Wrap can retain a value internally on failure; only OnError may fire.
	m := NewMaybe(func() adt.Result[adt.Option[int]] {
		return adt.Wrap(adt.Of(42), errSentinelMaybe)
	})

	var successCalled atomic.Bool
	var gotErr error
	sub := m.Subscribe(
		func(int) { successCalled.Store(true) },
		nil,
		func(e error) { gotErr = e },
	)
	<-sub.Done()

	if successCalled.Load() {
		t.Error("OnSuccess called despite err != nil")
	}
	if !errors.Is(gotErr, errSentinelMaybe) {
		t.Errorf("gotErr = %v, want sentinel", gotErr)
	}
}

func TestMaybe_Subscribe_NilCallbacks(t *testing.T) {
	t.Parallel()
	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](7)) })
	sub := m.Subscribe(nil, nil, nil)
	<-sub.Done()

	m2 := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) })
	sub2 := m2.Subscribe(nil, nil, nil)
	<-sub2.Done()

	m3 := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) })
	sub3 := m3.Subscribe(nil, nil, nil)
	<-sub3.Done()
}

func TestMaybe_Await_Success(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[string]] { return adt.Success(adt.Of[string]("x")) })

	v, err := m.Await().Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != "x" {
		t.Errorf("Await = (%v, %v), want (\"x\", true, nil)", v, err)
	}
	if !m.Done() {
		t.Error("Done() = false after Await")
	}
}

func TestMaybe_Await_Complete(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[string]] { return adt.Success(adt.Empty[string]()) })

	v, err := m.Await().Unwrap()
	if err != nil || v.IsPresent() {
		t.Errorf("Await = (%v, %v), want (\"\", false, nil)", v, err)
	}
}

func TestMaybe_Await_Error(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[string]] { return adt.Failure[adt.Option[string]](errSentinelMaybe) })

	v, err := m.Await().Unwrap()
	if !errors.Is(err, errSentinelMaybe) {
		t.Errorf("err = %v, want sentinel", err)
	}
	if v.IsPresent() {
		t.Error("present = true on error, want false")
	}
}

func TestMaybe_AwaitWithContext_CancelBeforeStart(t *testing.T) {
	t.Parallel()
	// ctx already done — fn should not be invoked.
	started := make(chan struct{})
	m := NewMaybe(func() adt.Result[adt.Option[int]] {
		close(started)
		return adt.Success(adt.Empty[int]())
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := m.AwaitWithContext(ctx).Unwrap()
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	select {
	case <-started:
		t.Error("fn invoked despite ctx already canceled")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestMaybe_AwaitWithContext_Done(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](99)) })
	v, err := m.AwaitWithContext(context.Background()).Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != 99 {
		t.Errorf("AwaitWithContext = (%v, %v), want (99, true, nil)", v, err)
	}
}

func TestMaybe_MultipleSubscribers_ShareResult(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	m := NewMaybe(func() adt.Result[adt.Option[string]] {
		calls.Add(1)
		return adt.Success(adt.Of[string]("shared"))
	})

	var wg sync.WaitGroup
	wg.Add(3)
	for range 3 {
		go func() {
			defer wg.Done()
			sub := m.Subscribe(func(string) {}, func() {}, func(error) {})
			<-sub.Done()
		}()
	}
	wg.Wait()

	if got := calls.Load(); got != 1 {
		t.Errorf("fn calls = %d, want 1", got)
	}
}

func TestMaybe_Map_Success(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](5)) }).
		Map(func(n int) string { return "x" + string(rune('0'+n)) })

	v, err := m.Await().Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != "x5" {
		t.Errorf("Map Await = (%v, %v), want (\"x5\", true, nil)", v, err)
	}
}

func TestMaybe_Map_CompletePassesThrough(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) }).
		Map(func(int) string { t.Error("f should not run on complete"); return "" })

	v, err := m.Await().Unwrap()
	if err != nil || v.IsPresent() {
		t.Errorf("Await = (%v, %v), want (_, false, nil)", v, err)
	}
}

func TestMaybe_Map_ErrorPassesThrough(t *testing.T) {
	t.Parallel()

	m := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) }).
		Map(func(int) string { t.Error("f should not run on error"); return "" })

	v, err := m.Await().Unwrap()
	if !errors.Is(err, errSentinelMaybe) {
		t.Errorf("err = %v, want sentinel", err)
	}
	if v.IsPresent() {
		t.Error("present = true on error, want false")
	}
}

func TestMaybe_FlatMap_Success(t *testing.T) {
	t.Parallel()

	inner := NewMaybe(func() adt.Result[adt.Option[string]] { return adt.Success(adt.Of[string]("inner")) })
	outer := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](1)) }).
		FlatMap(func(int) *Maybe[string] { return inner })

	v, err := outer.Await().Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != "inner" {
		t.Errorf("FlatMap Await = (%v, %v), want (\"inner\", true, nil)", v, err)
	}
}

func TestMaybe_FlatMap_OuterComplete(t *testing.T) {
	t.Parallel()

	outer := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) }).
		FlatMap(func(int) *Maybe[string] {
			t.Error("FlatMap f should not run on outer complete")
			return nil
		})

	v, err := outer.Await().Unwrap()
	if err != nil || v.IsPresent() {
		t.Errorf("Await = (%v, %v), want (_, false, nil)", v, err)
	}
}

func TestMaybe_Zip_BothPresent(t *testing.T) {
	t.Parallel()

	a := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](3)) })
	b := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](4)) })
	zipped := a.Zip(b, func(x, y int) int { return x*x + y*y })

	v, err := zipped.Await().Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != 25 {
		t.Errorf("Zip Await = (%v, %v), want (25, true, nil)", v, err)
	}
}

func TestMaybe_Zip_LeftComplete(t *testing.T) {
	t.Parallel()

	a := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) })
	b := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](4)) })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	v, err := zipped.Await().Unwrap()
	if err != nil || v.IsPresent() {
		t.Errorf("Await = (%v, %v), want (_, false, nil)", v, err)
	}
}

func TestMaybe_Zip_LeftError(t *testing.T) {
	t.Parallel()

	a := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) })
	b := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](4)) })
	zipped := a.Zip(b, func(x, y int) int { return x + y })

	_, err := zipped.Await().Unwrap()
	if !errors.Is(err, errSentinelMaybe) {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestMaybe_AndThen_RunsNextOnSuccess(t *testing.T) {
	t.Parallel()

	first := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Of[int](1)) })
	second := NewMaybe(func() adt.Result[adt.Option[string]] { return adt.Success(adt.Of[string]("next")) })
	seq := first.AndThen(second)

	v, err := seq.Await().Unwrap()
	if err != nil || v.IsEmpty() || v.Get() != "next" {
		t.Errorf("AndThen Await = (%v, %v), want (\"next\", true, nil)", v, err)
	}
}

func TestMaybe_AndThen_SkipsNextOnComplete(t *testing.T) {
	t.Parallel()

	first := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) })
	second := NewMaybe(func() adt.Result[adt.Option[string]] {
		t.Error("second should not run when first completes empty")
		return adt.Success(adt.Empty[string]())
	})
	seq := first.AndThen(second)

	v, err := seq.Await().Unwrap()
	if err != nil || v.IsPresent() {
		t.Errorf("Await = (%v, %v), want (_, false, nil)", v, err)
	}
}

func TestMaybe_AndThen_SkipsNextOnError(t *testing.T) {
	t.Parallel()

	first := NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) })
	second := NewMaybe(func() adt.Result[adt.Option[string]] {
		t.Error("second should not run when first errors")
		return adt.Success(adt.Empty[string]())
	})
	seq := first.AndThen(second)

	_, err := seq.Await().Unwrap()
	if !errors.Is(err, errSentinelMaybe) {
		t.Errorf("err = %v, want sentinel", err)
	}
}
