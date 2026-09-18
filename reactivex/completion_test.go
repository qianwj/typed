package reactivex

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/adt/option"
	r "github.com/qianwj/typed/adt/result"
)

func waitCompletion(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for completion")
	}
}

func TestSingleLateSubscription(t *testing.T) {
	t.Parallel()
	for _, want := range []r.Result[*int]{r.Success[*int](nil), r.Failure[*int](errSentinelSingle)} {
		var calls atomic.Int32
		s := NewSingle(func() r.Result[*int] { calls.Add(1); return want })
		s.Await()
		for range 3 {
			var received int
			sub := s.Subscribe(func(value *int) {
				received++
				if want.IsFailure() || value != nil {
					t.Error("unexpected success")
				}
			}, func(err error) {
				received++
				if !errors.Is(err, want.Error()) {
					t.Errorf("error = %v, want %v", err, want.Error())
				}
			})
			waitCompletion(t, sub.Done())
			if received != 1 {
				t.Fatalf("received %d callbacks, want 1", received)
			}
		}
		if calls.Load() != 1 {
			t.Fatal("late subscription restarted the source")
		}
	}
}

func TestMaybeLateSubscription(t *testing.T) {
	t.Parallel()
	for _, want := range []r.Result[option.Option[*int]]{
		r.Success(option.Of[*int](nil)),
		r.Success(option.Empty[*int]()),
		r.Failure[option.Option[*int]](errSentinelMaybe),
	} {
		var calls atomic.Int32
		m := NewMaybe(func() r.Result[option.Option[*int]] { calls.Add(1); return want })
		m.Await()
		for range 3 {
			var received int
			sub := m.Subscribe(func(value *int) {
				received++
				if want.IsFailure() || want.Value().IsEmpty() || value != nil {
					t.Error("unexpected success")
				}
			}, func() {
				received++
				if want.IsFailure() || want.Value().IsPresent() {
					t.Error("unexpected empty completion")
				}
			}, func(err error) {
				received++
				if !errors.Is(err, want.Error()) {
					t.Errorf("error = %v, want %v", err, want.Error())
				}
			})
			waitCompletion(t, sub.Done())
			if received != 1 {
				t.Fatalf("received %d callbacks, want 1", received)
			}
		}
		if calls.Load() != 1 {
			t.Fatal("late subscription restarted the source")
		}
	}
}

func TestCompletionConcurrentRegistration(t *testing.T) {
	t.Parallel()
	gate := make(chan struct{})
	var sources atomic.Int32
	c := newCompletion(func(complete func(r.Result[int])) {
		sources.Add(1)
		<-gate
		complete(r.Success(42))
	})
	const consumers = 128
	counts := make([]atomic.Int32, consumers)
	var workers sync.WaitGroup
	for i := range consumers {
		workers.Go(func() {
			sub := c.subscribe(func(result r.Result[int]) {
				counts[i].Add(1)
				if result.IsFailure() || result.Value() != 42 {
					t.Errorf("unexpected result: %v", result)
				}
			})
			<-sub.Done()
		})
	}
	close(gate) // Completion races registration by the remaining consumers.
	done := make(chan struct{})
	go func() { workers.Wait(); close(done) }()
	waitCompletion(t, done)
	for i := range counts {
		if counts[i].Load() != 1 {
			t.Fatalf("consumer %d received %d callbacks", i, counts[i].Load())
		}
	}
	if sources.Load() != 1 {
		t.Fatalf("source ran %d times", sources.Load())
	}
}

func TestCompletionSettlesOnce(t *testing.T) {
	t.Parallel()
	returned := make(chan struct{})
	c := newCompletion(func(complete func(r.Result[int])) {
		complete(r.Success(42))
		complete(r.Failure[int](errSentinelSingle))
		close(returned)
	})
	var calls atomic.Int32
	sub := c.subscribe(func(r.Result[int]) { calls.Add(1) })
	waitCompletion(t, returned)
	waitCompletion(t, sub.Done())
	if result := c.await(context.Background()); result.IsFailure() || result.Value() != 42 || calls.Load() != 1 {
		t.Fatalf("duplicate completion changed the result or repeated delivery: %v, calls=%d", result, calls.Load())
	}
}

func TestCompletionCancelRacesDelivery(t *testing.T) {
	t.Parallel()
	gate, finished := make(chan struct{}), make(chan struct{})
	c := newCompletion(func(complete func(r.Result[int])) {
		<-gate
		complete(r.Success(42))
		close(finished)
	})
	var workers sync.WaitGroup
	counts := make([]atomic.Int32, 64)
	for i := range counts {
		sub := c.subscribe(func(r.Result[int]) { counts[i].Add(1) })
		workers.Go(func() {
			<-gate
			sub.Cancel()
			<-sub.Done()
		})
	}
	close(gate)
	waitCompletion(t, finished)
	workers.Wait()
	for i := range counts {
		if counts[i].Load() > 1 {
			t.Fatalf("consumer %d received duplicate callbacks", i)
		}
	}
	if c.await(context.Background()).Value() != 42 {
		t.Fatal("canceling subscribers changed the shared result")
	}
}

func TestCompletionCancelSubscription(t *testing.T) {
	t.Parallel()
	gate := make(chan struct{})
	s := NewSingle(func() r.Result[int] { <-gate; return r.Success(42) })
	var calls atomic.Int32
	canceled := s.Subscribe(func(int) { calls.Add(1) }, nil)
	canceled.Cancel()
	canceled.Cancel()
	waitCompletion(t, canceled.Done())
	var got int
	active := s.Subscribe(func(value int) { got = value }, nil)
	close(gate)
	waitCompletion(t, active.Done())
	if got != 42 || calls.Load() != 0 || s.Await().Value() != 42 {
		t.Fatal("cancellation either delivered a callback or affected the shared result")
	}
}

func TestCompletionCachedCallbackDoesNotBlockSubscribe(t *testing.T) {
	t.Parallel()
	s := NewSingle(func() r.Result[int] { return r.Success(42) })
	s.Await()
	release := make(chan struct{})
	defer close(release)
	started := make(chan struct{})
	returned := make(chan Subscription, 1)
	go func() {
		returned <- s.Subscribe(func(int) { close(started); <-release }, nil)
	}()
	waitCompletion(t, started)
	select {
	case sub := <-returned:
		sub.Cancel() // Must not wait for the callback to return.
		waitCompletion(t, sub.Done())
	case <-time.After(5 * time.Second):
		t.Fatal("Subscribe blocked on a cached callback")
	}
}

func TestCompletionCallbackCanReenter(t *testing.T) {
	t.Parallel()
	gate := make(chan struct{})
	s := NewSingle(func() r.Result[int] { <-gate; return r.Success(42) })
	var sub Subscription
	sub = s.Subscribe(func(value int) {
		if s.Await().Value() != value {
			t.Error("reentrant Await lost the result")
		}
		nested := s.Subscribe(func(int) {}, nil)
		<-nested.Done()
		sub.Cancel()
	}, nil)
	close(gate)
	waitCompletion(t, sub.Done())
}

func TestSingleCachedMapCanCancelWait(t *testing.T) {
	t.Parallel()
	s := NewSingle(func() r.Result[int] { return r.Success(42) })
	s.Await()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release := make(chan struct{})
	mapped := s.Map(func(v int) int { cancel(); <-release; return v + 1 })
	returned := make(chan r.Result[int], 1)
	go func() { returned <- mapped.AwaitWithContext(ctx) }()
	select {
	case result := <-returned:
		close(release)
		if !errors.Is(result.Error(), context.Canceled) {
			t.Errorf("result = %v, want cancellation", result)
		}
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("cached transform blocked AwaitWithContext")
	}
	if result := mapped.Await(); result.IsFailure() || result.Value() != 43 {
		t.Fatalf("eventual result = %v, want 43", result)
	}
}

func TestMaybeCachedMapCanCancelWait(t *testing.T) {
	t.Parallel()
	m := NewMaybe(func() r.Result[option.Option[int]] { return r.Success(option.Of(42)) })
	m.Await()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release := make(chan struct{})
	mapped := m.Map(func(v int) int { cancel(); <-release; return v + 1 })
	returned := make(chan r.Result[option.Option[int]], 1)
	go func() { returned <- mapped.AwaitWithContext(ctx) }()
	select {
	case result := <-returned:
		close(release)
		if !errors.Is(result.Error(), context.Canceled) {
			t.Errorf("result = %v, want cancellation", result)
		}
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("cached transform blocked AwaitWithContext")
	}
	if result := mapped.Await(); result.IsFailure() || result.Value().IsEmpty() || result.Value().Get() != 43 {
		t.Fatalf("eventual result = %v, want present 43", result)
	}
}
