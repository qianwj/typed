package reactivex

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestSubject_OnSubscribeCancelsUpstreamWhenTerminated covers the closed
// branch of Subject.OnSubscribe: when the subject is already terminated,
// a fresh upstream subscription must be cancelled instead of being held.
func TestSubject_OnSubscribeCancelsUpstreamWhenTerminated(t *testing.T) {
	t.Parallel()
	s := NewSubject[int]()
	s.OnError(errors.New("seed"))

	var cancelled atomic.Bool
	upstream := &cancelHookSubscription{cancel: func() { cancelled.Store(true) }}
	s.OnSubscribe(upstream)
	if !cancelled.Load() {
		t.Fatal("expected upstream.Cancel() to be called when subject is already terminated")
	}
}

// TestSubject_OnErrorNoSubscribers covers the no-subscribers branch of
// OnError: it should still record the terminal state so late subscribers
// see the stored error.
func TestSubject_OnErrorNoSubscribers(t *testing.T) {
	t.Parallel()
	s := NewSubject[int]()
	want := errors.New("boom")
	s.OnError(want)
	// Late subscriber: should see the error immediately on Subscribe.
	got := make(chan error, 1)
	s.Subscribe(context.Background(), &captureSubscriber[int]{
		// not used; we'll provide a custom subscriber below
	})
	// Replace with a dedicated subscriber to capture the late error.
	done := make(chan struct{})
	s.Subscribe(context.Background(), &errCaptureSubscriber{err: got, done: done})
	select {
	case err := <-got:
		if !errors.Is(err, want) {
			t.Fatalf("late subscriber got %v, want %v", err, want)
		}
	case <-time.After(time.Second):
		t.Fatal("late subscriber did not receive terminal error")
	}
	<-done
}

// errCaptureSubscriber captures the next OnError and signals completion
// without retaining the value, used to assert that the Subject delivers
// its stored error to subscribers that arrive after termination.
type errCaptureSubscriber struct {
	err  chan<- error
	done chan struct{}
}

func (e *errCaptureSubscriber) OnSubscribe(s Subscription) { s.Request(1) }
func (e *errCaptureSubscriber) OnNext(int)                 {}
func (e *errCaptureSubscriber) OnError(err error)          { e.err <- err; close(e.done) }
func (e *errCaptureSubscriber) OnComplete()                 { close(e.done) }

// TestOptions_ApplyBackpressure covers applyBackpressureOptions: nil
// entries are skipped, DropOldest without a positive buffer is rejected
// by the deferred panic at NewSubject time, and KeepLatest forces buffer
// to 1.
func TestOptions_ApplyBackpressure(t *testing.T) {
	t.Parallel()

	t.Run("NilOptionSkipped", func(t *testing.T) {
		got := applyBackpressureOptions([]BackpressureOption{nil, WithBuffer(8)})
		if got.buffer != 8 {
			t.Fatalf("buffer = %d, want 8", got.buffer)
		}
		if got.overflow != OverflowBlock {
			t.Fatalf("overflow = %d, want OverflowBlock", got.overflow)
		}
	})

	t.Run("KeepLatestForcesBufferOne", func(t *testing.T) {
		got := applyBackpressureOptions([]BackpressureOption{WithBuffer(64), WithOverflow(OverflowKeepLatest)})
		if got.buffer != 1 {
			t.Fatalf("KeepLatest should force buffer to 1; got %d", got.buffer)
		}
	})

	t.Run("DropOldestRejectsZeroBuffer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for DropOldest with zero buffer")
			}
		}()
		applyBackpressureOptions([]BackpressureOption{WithOverflow(OverflowDropOldest)})
	})

	t.Run("DefaultOverflowIsBlock", func(t *testing.T) {
		got := applyBackpressureOptions(nil)
		if got.overflow != OverflowBlock {
			t.Fatalf("default overflow = %d, want OverflowBlock", got.overflow)
		}
	})

	t.Run("InvalidOverflowPanics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid overflow strategy")
			}
		}()
		WithOverflow(OverflowStrategy(99))
	})

	t.Run("NegativeBufferPanics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for negative buffer size")
			}
		}()
		WithBuffer(-1)
	})
}

// TestCollect_ErrorAndCancellation covers the error and ctx-cancel
// branches of ToSlice / Collect that the smoke test does not exercise.
func TestCollect_ErrorAndCancellation(t *testing.T) {
	t.Parallel()

	t.Run("ToSliceSurfacesUpstreamError", func(t *testing.T) {
		want := errors.New("upstream fail")
		obs := Create[int](func(ctx context.Context, emit func(int) bool, complete func()) {
			// Trigger an upstream error via Subject instead, since Create
			// has no direct error-emit hook. The simpler path: subscribe
			// a Subject and ask ToSlice to collect; closing the Subject
			// with an error before subscribing is the cleanest way.
			_ = ctx
			_ = emit
			_ = complete
		})
		_ = obs // not used; the Subject path below is the actual test
		s := NewSubject[int]()
		s.OnError(want)
		values, err := subjectAsObservable(s).ToSlice(context.Background())
		if err == nil || !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
		if len(values) != 0 {
			t.Fatalf("values = %v, want empty", values)
		}
	})

	t.Run("CollectReturnsErrorAndInitial", func(t *testing.T) {
		want := errors.New("boom")
		s := NewSubject[int]()
		s.OnError(want)
		sum, err := Collect(context.Background(), subjectAsObservable(s), 100,
			func(acc, v int) int { return acc + v })
		if err == nil || !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
		if sum != 100 {
			t.Fatalf("sum = %d, want 100 (initial unchanged on error)", sum)
		}
	})

	t.Run("SubscribeWithNilCtx", func(t *testing.T) {
		// Subscribe must treat a nil ctx as context.Background.
		var s Subscription
		s = Just(1, 2, 3).Subscribe(nil, &captureSubscriber[int]{})
		if s == nil {
			t.Fatal("expected non-nil subscription")
		}
	})

	t.Run("ToSliceWithNilCtx", func(t *testing.T) {
		values, err := Just(7, 8, 9).ToSlice(nil)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(values) != 3 || values[0] != 7 || values[1] != 8 || values[2] != 9 {
			t.Fatalf("values = %v", values)
		}
	})
}

// subjectAsObservable wraps a hot Subject in an Observable so ToSlice /
// Collect can subscribe to it through the standard Observable API.
func subjectAsObservable[T any](s *Subject[T]) Observable[T] {
	return Observable[T]{
		subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
			return s.Subscribe(ctx, out)
		},
	}
}

// withCancelHook returns a no-op Subscription whose Cancel records
// that it was called; used to assert that Subject.OnSubscribe cancels
// the upstream when the subject is already terminated.
func (n nullSubscription) withCancelHook(fn func()) cancelHookSubscription {
	return cancelHookSubscription{cancel: fn}
}

type cancelHookSubscription struct {
	cancel       func()
	requestCount atomic.Int64
}

func (c *cancelHookSubscription) Request(n uint64) { c.requestCount.Add(int64(n)) }
func (c *cancelHookSubscription) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}
func (c *cancelHookSubscription) Done() <-chan struct{} { return nil }
func (c *cancelHookSubscription) IsClosed() bool        { return false }
func (c *cancelHookSubscription) Err() error             { return nil }
