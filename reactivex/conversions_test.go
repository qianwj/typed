package reactivex

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/adt"
)

func TestToFlowableSharesCachedResult(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"Single", "Maybe"} {
		for _, tc := range []struct {
			name    string
			value   any
			present bool
			err     error
		}{
			{"value", 42, true, nil},
			{"zero", 0, true, nil},
			{"nil", nil, true, nil},
			{"typed nil", (*int)(nil), true, nil},
			{"empty", nil, false, nil},
			{"error", nil, false, errSentinelSingle},
		} {
			if kind == "Single" && tc.name == "empty" {
				continue
			}
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				var calls atomic.Int32
				var flow Flowable[any]
				if kind == "Single" {
					flow = NewSingle(func() adt.Result[any] {
						calls.Add(1)
						return adt.Wrap(tc.value, tc.err)
					}).ToFlowable()
				} else {
					flow = NewMaybe(func() adt.Result[adt.Option[any]] {
						calls.Add(1)
						value := adt.Empty[any]()
						if tc.present {
							value = adt.Of(tc.value)
						}
						return adt.Wrap(value, tc.err)
					}).ToFlowable()
				}
				if calls.Load() != 0 {
					t.Fatal("conversion started the source")
				}
				for range 2 {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					got, err := flow.ToSlice(ctx)
					cancel()
					want := []any{}
					if tc.err == nil && tc.present {
						want = append(want, tc.value)
					}
					if !errors.Is(err, tc.err) || !reflect.DeepEqual(got, want) {
						t.Fatalf("got (%v, %v), want (%v, %v)", got, err, want, tc.err)
					}
				}
				if calls.Load() != 1 {
					t.Fatalf("source ran %d times, want 1", calls.Load())
				}
			})
		}
	}
}

func TestToFlowableDemandIsPerSubscription(t *testing.T) {
	t.Parallel()
	flow := NewSingle(func() adt.Result[int] { return adt.Success(42) }).ToFlowable()
	values := make(chan int, 2)
	var terminals atomic.Int32
	slow := flow.Subscribe(nil, subscriberFuncs[int]{
		onNext:     func(value int) { values <- value },
		onComplete: func() { terminals.Add(1) },
	})
	defer slow.Cancel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	got, err := flow.ToSlice(ctx)
	if err != nil || !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("subscriber without demand blocked another subscriber: %v, %v", got, err)
	}
	slow.Request(0)
	select {
	case <-values:
		t.Fatal("value delivered without demand")
	case <-slow.Done():
		t.Fatal("completed before delivering the pending value")
	case <-time.After(20 * time.Millisecond):
	}
	slow.Request(1)
	waitCompletion(t, slow.Done())
	if len(values) != 1 || <-values != 42 || terminals.Load() != 1 {
		t.Fatal("expected one value followed by one completion")
	}
}

func TestToFlowableTerminalsDoNotRequireDemand(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		flow Flowable[int]
		err  error
	}{
		{NewSingle(func() adt.Result[int] { return adt.Failure[int](errSentinelSingle) }).ToFlowable(), errSentinelSingle},
		{NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Success(adt.Empty[int]()) }).ToFlowable(), nil},
		{NewMaybe(func() adt.Result[adt.Option[int]] { return adt.Failure[adt.Option[int]](errSentinelMaybe) }).ToFlowable(), errSentinelMaybe},
	} {
		var received int
		var got error
		sub := tc.flow.Subscribe(context.Background(), subscriberFuncs[int]{
			onNext:     func(int) { t.Error("unexpected value") },
			onError:    func(err error) { received++; got = err },
			onComplete: func() { received++ },
		})
		waitCompletion(t, sub.Done())
		if received != 1 || !errors.Is(got, tc.err) {
			t.Fatalf("terminal count=%d, error=%v, want %v", received, got, tc.err)
		}
	}
}

func TestToFlowableCancelDuringSetup(t *testing.T) {
	t.Parallel()
	for _, preCanceled := range []bool{false, true} {
		var calls atomic.Int32
		source := NewSingle(func() adt.Result[int] { calls.Add(1); return adt.Success(42) })
		ctx, cancel := context.WithCancel(context.Background())
		if preCanceled {
			cancel()
		}
		sub := source.ToFlowable().Subscribe(ctx, subscriberFuncs[int]{
			onSubscribe: func(sub Subscription) {
				if !preCanceled {
					sub.Cancel()
				}
			},
			onNext: func(int) { t.Error("value delivered after setup cancellation") },
		})
		waitCompletion(t, sub.Done())
		cancel()
		// Await must still be able to start the untouched shared source.
		if calls.Load() != 0 || source.Done() {
			t.Fatal("canceled setup started the shared source")
		}
		if source.Await().Value() != 42 || calls.Load() != 1 {
			t.Fatal("canceled setup affected the shared computation")
		}
	}
}

func TestToFlowableContextCancelDoesNotCancelSharedSource(t *testing.T) {
	t.Parallel()
	started, release := make(chan struct{}), make(chan struct{})
	source := NewMaybe(func() adt.Result[adt.Option[int]] {
		close(started)
		<-release
		return adt.Success(adt.Of(42))
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var received atomic.Int32
	sub := source.ToFlowable().Subscribe(ctx, subscriberFuncs[int]{
		onSubscribe: func(sub Subscription) { sub.Request(1) },
		onNext:      func(int) { received.Add(1) },
		onComplete:  func() { received.Add(1) },
		onError:     func(error) { received.Add(1) },
	})
	waitCompletion(t, started)
	cancel()
	waitCompletion(t, sub.Done())
	close(release)
	if source.Await().Value().Get() != 42 || received.Load() != 0 {
		t.Fatal("cancellation affected the source or delivered a terminal signal")
	}
}

func TestFirstConversions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		values []any
	}{
		{"empty", nil},
		{"zero", []any{0}},
		{"nil", []any{nil}},
		{"typed nil", []any{(*int)(nil)}},
		{"multiple", []any{42, 99}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			flow := FromSlice(tc.values)
			maybe := flow.FirstElement(nil).AwaitWithContext(ctx)
			single := flow.FirstOrError(nil).AwaitWithContext(ctx)
			if len(tc.values) == 0 {
				if maybe.IsFailure() || maybe.Value().IsPresent() || !errors.Is(single.Error(), ErrNoElements) {
					t.Fatalf("empty: Maybe=%v, Single=%v", maybe, single)
				}
			} else if maybe.IsFailure() || maybe.Value().IsEmpty() || !reflect.DeepEqual(maybe.Value().Get(), tc.values[0]) || single.IsFailure() || !reflect.DeepEqual(single.Value(), tc.values[0]) {
				t.Fatalf("Maybe=%v, Single=%v, want first %v", maybe, single, tc.values[0])
			}
		})
	}
}

// conversionSubscription supports synchronous notifications during Request
// and a reentrant terminal notification during Cancel.
type conversionSubscription struct {
	request func(uint64)
	cancel  func()
	done    chan struct{}
	once    sync.Once
}

func (s *conversionSubscription) Request(n uint64) { s.request(n) }
func (s *conversionSubscription) Cancel() {
	// Release Once before invoking user code, which may cancel reentrantly.
	var cancel func()
	s.once.Do(func() { close(s.done); cancel = s.cancel })
	if cancel != nil {
		cancel()
	}
}
func (s *conversionSubscription) Done() <-chan struct{} { return s.done }

func TestFirstRequestsOneAndCancelsBeforePublishing(t *testing.T) {
	t.Parallel()
	var subscriptions, requested atomic.Int32
	var upstream Subscription
	returned := make(chan struct{})
	flow := Flowable[int]{subscribe: func(_ context.Context, out Subscriber[int]) Subscription {
		subscriptions.Add(1)
		sub := &conversionSubscription{done: make(chan struct{})}
		upstream = sub
		sub.cancel = out.OnComplete
		sub.request = func(n uint64) {
			requested.Add(int32(n))
			out.OnNext(42)
			// A late notification from an uncooperative source cannot overwrite first.
			out.OnError(errSentinelSingle)
		}
		out.OnSubscribe(sub)
		close(returned)
		return sub
	}}
	first := flow.FirstOrError(nil)
	if subscriptions.Load() != 0 {
		t.Fatal("conversion subscribed eagerly")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for range 2 {
		if got := first.AwaitWithContext(ctx); got.IsFailure() || got.Value() != 42 {
			t.Fatalf("first = %v, want 42", got)
		}
	}
	waitCompletion(t, returned)
	waitCompletion(t, upstream.Done())
	if subscriptions.Load() != 1 || requested.Load() != 1 {
		t.Fatalf("subscriptions=%d, requested=%d, want 1 each", subscriptions.Load(), requested.Load())
	}
}

func TestFirstPropagatesUpstreamError(t *testing.T) {
	t.Parallel()
	flow := NewSingle(func() adt.Result[int] { return adt.Failure[int](errSentinelSingle) }).ToFlowable()
	if got := flow.FirstElement(nil).Await(); !errors.Is(got.Error(), errSentinelSingle) {
		t.Fatalf("Maybe error = %v", got.Error())
	}
	if got := flow.FirstOrError(nil).Await(); !errors.Is(got.Error(), errSentinelSingle) {
		t.Fatalf("Single error = %v", got.Error())
	}
}

func TestFirstOrErrorAfterReduce(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result := Just(1, 2, 3).Reduce(func(a, b int) int { return a + b }).FirstOrError(ctx).Await()
	if result.IsFailure() || result.Value() != 6 {
		t.Fatalf("reduced first result = %v, want 6", result)
	}
}

func TestReduceTranslatesDemandOnce(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	upstream := &conversionSubscription{
		done: make(chan struct{}),
		request: func(n uint64) {
			requests.Add(1)
			if n != ^uint64(0) {
				t.Errorf("upstream request = %d, want maximum demand", n)
			}
		},
	}
	source := Flowable[int]{subscribe: func(_ context.Context, out Subscriber[int]) Subscription {
		out.OnSubscribe(upstream)
		return upstream
	}}
	var supplied Subscription
	returned := source.Reduce(func(a, b int) int { return a + b }).Subscribe(nil, subscriberFuncs[int]{
		onSubscribe: func(sub Subscription) { supplied = sub },
	})
	if supplied != returned {
		t.Fatal("Subscribe and OnSubscribe returned different demand handles")
	}
	returned.Request(0)
	if requests.Load() != 0 {
		t.Fatal("zero demand requested upstream input")
	}
	returned.Request(1)
	returned.Request(10)
	if requests.Load() != 1 {
		t.Fatalf("upstream requested %d times, want 1", requests.Load())
	}
	returned.Cancel()
	waitCompletion(t, returned.Done())
}

func TestFirstContextCancellation(t *testing.T) {
	t.Parallel()
	for _, preCanceled := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		// Multiple paths below early-return via t.Fatal before they
		// reach the bottom of the loop. context.CancelFunc is safe
		// to call more than once, so a single defer covers them all
		// without affecting the preCanceled == true branch, which
		// still cancels at the top.
		defer cancel()
		subscribed := make(chan Subscription, 1)
		flow := Flowable[int]{subscribe: func(_ context.Context, out Subscriber[int]) Subscription {
			sub := newSubscription()
			out.OnSubscribe(sub)
			subscribed <- sub
			return sub
		}}
		if preCanceled {
			cancel()
		}
		first := flow.FirstElement(ctx)
		result := make(chan adt.Result[adt.Option[int]], 1)
		go func() { result <- first.Await() }()
		if !preCanceled {
			var upstream Subscription
			select {
			case upstream = <-subscribed:
			case <-time.After(5 * time.Second):
				cancel()
				t.Fatal("source did not subscribe")
			}
			cancel()
			waitCompletion(t, upstream.Done())
		}
		select {
		case got := <-result:
			if !errors.Is(got.Error(), context.Canceled) {
				t.Fatalf("result = %v, want cancellation", got)
			}
		case <-time.After(5 * time.Second):
			cancel()
			t.Fatal("cancellation did not settle the result")
		}
		select {
		case <-subscribed:
			t.Fatal("pre-canceled conversion subscribed")
		default:
		}
	}
}

func TestFirstCancelWaitKeepsUpstreamAlive(t *testing.T) {
	t.Parallel()
	connected := make(chan Subscriber[int], 1)
	upstream := newSubscription()
	flow := Flowable[int]{subscribe: func(_ context.Context, out Subscriber[int]) Subscription {
		out.OnSubscribe(upstream)
		connected <- out
		return upstream
	}}
	first := flow.FirstElement(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan adt.Result[adt.Option[int]], 1)
	go func() { result <- first.AwaitWithContext(ctx) }()
	var out Subscriber[int]
	select {
	case out = <-connected:
	case <-time.After(5 * time.Second):
		t.Fatal("source did not subscribe")
	}
	cancel()
	select {
	case got := <-result:
		if !errors.Is(got.Error(), context.Canceled) {
			t.Fatalf("wait result = %v", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait did not cancel")
	}
	select {
	case <-upstream.Done():
		t.Fatal("canceling a wait canceled the shared subscription")
	default:
	}
	out.OnNext(42)
	if got := first.Await(); got.IsFailure() || got.Value().IsEmpty() || got.Value().Get() != 42 {
		t.Fatalf("eventual result = %v, want 42", got)
	}
	waitCompletion(t, upstream.Done())
}
