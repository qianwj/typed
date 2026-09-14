package reactivex

import (
	"context"
	"errors"
	"sync"
	"testing"
)

// captureSubscriber records the next error / complete notification and
// discards values. It implements Subscriber[T] for any T so the operator
// forwarding tests can reuse it across type parameters.
type captureSubscriber[T any] struct {
	err   error
	done  bool
	count int
}

func (c *captureSubscriber[T]) OnSubscribe(s Subscription) { s.Request(1 << 30) }
func (c *captureSubscriber[T]) OnNext(T)                    { c.count++ }
func (c *captureSubscriber[T]) OnError(e error)             { c.err = e }
func (c *captureSubscriber[T]) OnComplete()                 { c.done = true }

// nullSubscription is a no-op Subscription for adapter tests that exercise
// only the forwarding methods, not the upstream-cancellation path.
type nullSubscription struct{}

func (nullSubscription) Request(uint64)            {}
func (nullSubscription) Cancel()                   {}
func (nullSubscription) Done() <-chan struct{}     { return nil }
func (nullSubscription) IsClosed() bool            { return true }
func (nullSubscription) Err() error                 { return nil }

// TestAdapter_OnErrorCompleteForwarding covers the trivial OnError /
// OnComplete forwarding methods on every operator's subscriber adapter.
// These methods were 0% covered before because no test exercised the
// upstream-emitted-error path through a fluent pipeline; the tests below
// construct each adapter directly and call the methods on it.
func TestAdapter_OnErrorCompleteForwarding(t *testing.T) {
	t.Parallel()

	want := errors.New("upstream failure")

	t.Run("Map", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &mapSubscriber[int, int]{
			out:   cap,
			sub:   nullSubscription{},
			ctx:   context.Background(),
			f:     func(context.Context, int) (int, error) { return 0, nil },
			once:  sync.Once{},
		}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
		cap.err = nil
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
	})

	t.Run("Filter", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &filterSubscriber[int]{
			out:       cap,
			sub:       nullSubscription{},
			predicate: func(int) bool { return true },
		}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
		cap.err = nil
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
	})

	t.Run("Take_ForwardsErrorWhenLeftPositive", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &takeSubscriber[int]{out: cap, sub: nullSubscription{}, left: 5}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
	})

	t.Run("Take_DropsErrorWhenLeftZero", func(t *testing.T) {
		// After the limit has been reached the Take operator emits its own
		// OnComplete and discards a late upstream error. This branch was
		// previously uncovered.
		cap := &captureSubscriber[int]{}
		s := &takeSubscriber[int]{out: cap, sub: nullSubscription{}, left: 0}
		s.OnError(want)
		if cap.err != nil {
			t.Fatalf("expected error to be dropped, got %v", cap.err)
		}
	})

	t.Run("Take_OnCompleteOnlyWhenLeftPositive", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &takeSubscriber[int]{out: cap, sub: nullSubscription{}, left: 0}
		s.OnComplete()
		if cap.done {
			t.Fatal("expected OnComplete to be dropped when left == 0")
		}
	})

	t.Run("Take_OnCompleteForwardedWhenLeftPositive", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &takeSubscriber[int]{out: cap, sub: nullSubscription{}, left: 2}
		s.OnComplete()
		if !cap.done {
			t.Fatal("expected OnComplete forwarded when left > 0")
		}
	})

	t.Run("Skip", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &skipSubscriber[int]{out: cap, sub: nullSubscription{}, left: 0}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
		cap.err = nil
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
	})

	t.Run("Scan", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &scanSubscriber[int, int]{
			out:   cap,
			value: 0,
			f:     func(acc, v int) int { return acc + v },
		}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
		cap.err = nil
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
	})

	t.Run("Reduce_OnError", func(t *testing.T) {
		cap := &captureSubscriber[int]{}
		s := &reduceSubscriber[int]{
			out: cap,
			f:   func(a, b int) int { return a + b },
		}
		s.OnError(want)
		if !errors.Is(cap.err, want) {
			t.Fatalf("OnError not forwarded: got %v", cap.err)
		}
	})

	t.Run("Reduce_OnCompleteEmitsPartial", func(t *testing.T) {
		// OnComplete with hasValue=true must emit the accumulator before
		// completion (it is the whole point of Reduce).
		cap := &captureSubscriber[int]{}
		s := &reduceSubscriber[int]{
			out: cap, f: func(a, b int) int { return a + b },
		}
		s.OnNext(7)
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
		if cap.count != 1 {
			t.Fatalf("Reduce should have emitted the accumulator once; got %d", cap.count)
		}
	})

	t.Run("Reduce_OnCompleteEmptyInput", func(t *testing.T) {
		// OnComplete with no input must not emit a value but still
		// forward completion.
		cap := &captureSubscriber[int]{}
		s := &reduceSubscriber[int]{
			out: cap, f: func(a, b int) int { return a + b },
		}
		s.OnComplete()
		if !cap.done {
			t.Fatal("OnComplete not forwarded")
		}
		if cap.count != 0 {
			t.Fatalf("Reduce must not emit a value on empty input; got %d", cap.count)
		}
	})
}
