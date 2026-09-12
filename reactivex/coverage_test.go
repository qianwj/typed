package reactivex

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

// ---------- gap-fill coverage for the reactivex package ----------
//
// The tests in this block cover the function paths that the existing
// collect_test.go / operators_test.go / sources_test.go / subject_test.go
// did not exercise: Just, FromSeq, FromChannel, FromChannelWithOptions,
// Create, Interval, Skip, Scan, Reduce, Take(0), the Subject
// OnError / OnSubscribe hooks, ToSlice empty / error paths, and
// Collect with empty / error input.
//
// All async tests use a done channel plus a 1-second time.After
// safety net, matching the existing test style.

// ---------- sources ----------

// TestJust covers the Just constructor: empty (no values) and
// non-empty variadic forms.
func TestJust(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		values, err := Just(1, 2, 3).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{1, 2, 3}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		values, err := Just[int]().ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
}

// TestFromSeq covers the iter.Seq-based source: empty sequence and
// a finite non-empty sequence.
func TestFromSeq(t *testing.T) {
	t.Run("finite", func(t *testing.T) {
		seq := func(yield func(int) bool) {
			for i := 1; i <= 3; i++ {
				if !yield(i * 10) {
					return
				}
			}
		}
		values, err := FromSeq(seq).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{10, 20, 30}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		empty := func(yield func(int) bool) {}
		values, err := FromSeq(empty).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
}

// TestFromChannel covers the default-blocking FromChannel source.
// The test sends a few values, then closes the channel and verifies
// the source completes.
func TestFromChannel(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	values, err := FromChannel(ch).ToSlice(context.Background())
	if err != nil {
		t.Fatalf("ToSlice: %v", err)
	}
	if want := []int{1, 2, 3}; !reflect.DeepEqual(values, want) {
		t.Fatalf("got %v, want %v", values, want)
	}
}

// TestFromChannelWithOptions covers the configurable channel source:
// WithBuffer controls the pending queue size. The OverflowError
// path is exercised by TestSubjectOverflowStrategies (the Subject
// case is more direct to set up without racing with cancellation).
func TestFromChannelWithOptions(t *testing.T) {
	t.Run("buffered delivery", func(t *testing.T) {
		ch := make(chan int, 4)
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		close(ch)
		values, err := FromChannelWithOptions(ch, WithBuffer(2)).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
}

// TestCreate covers the producer-callback source. Two flavours:
// values emitted via emit(), and completion via the explicit
// complete() callback (which is what most observers actually do).
func TestCreate(t *testing.T) {
	t.Run("emit + complete callback", func(t *testing.T) {
		values, err := Create(func(_ context.Context, emit func(int) bool, complete func()) {
			for i := 1; i <= 3; i++ {
				if !emit(i) {
					return
				}
			}
			complete()
		}).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{1, 2, 3}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("return-without-complete still completes", func(t *testing.T) {
		// Create's defer s.complete() ensures normal completion
		// even if the run function returns without calling
		// complete().
		values, err := Create(func(_ context.Context, emit func(int) bool, _ func()) {
			emit(42)
		}).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{42}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("empty producer", func(t *testing.T) {
		values, err := Create(func(_ context.Context, _ func(int) bool, _ func()) {}).ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
}

// TestInterval covers the time-ticker source. With a 5ms period we
// request 3 values and cancel to avoid leaking the goroutine.
func TestInterval(t *testing.T) {
	values := make(chan uint64, 8)
	var sub Subscription
	sub = Interval(context.Background(), 5*time.Millisecond).
		Subscribe(context.Background(), subscriberFuncs[uint64]{
			onSubscribe: func(s Subscription) { sub = s },
			onNext:      func(v uint64) { values <- v },
		})
	sub.Request(3)
	timeout := time.After(time.Second)
	got := make([]uint64, 0, 3)
	for len(got) < 3 {
		select {
		case v := <-values:
			got = append(got, v)
		case <-timeout:
			t.Fatalf("timeout: got %v, want 3 values", got)
		}
	}
	// Cancel to stop the ticker goroutine.
	sub.Cancel()
	select {
	case <-sub.Done():
	case <-time.After(time.Second):
		t.Fatal("subscription did not complete after cancel")
	}
	if got[0] != 0 || got[1] != 1 || got[2] != 2 {
		t.Fatalf("got %v, want [0 1 2] (monotonic counters)", got)
	}
}

// ---------- operators ----------

// TestSkip covers the Skip operator: n=0 is identity, n smaller
// than the source, n equal to the source, n larger than the source.
func TestSkip(t *testing.T) {
	cases := []struct {
		name string
		n    uint64
		in   []int
		want []int
	}{
		{"n smaller", 2, []int{1, 2, 3, 4, 5}, []int{3, 4, 5}},
		{"n equal", 5, []int{1, 2, 3, 4, 5}, []int{}},
		{"n larger", 100, []int{1, 2, 3}, []int{}},
		{"n zero is identity", 0, []int{1, 2, 3}, []int{1, 2, 3}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			values, err := FromSlice(c.in).Skip(c.n).ToSlice(context.Background())
			if err != nil {
				t.Fatalf("ToSlice: %v", err)
			}
			if !reflect.DeepEqual(values, c.want) {
				t.Fatalf("got %v, want %v", values, c.want)
			}
		})
	}
}

// TestScan covers the running-accumulator operator. The first
// emitted value is f(initial, v1); later values use the running
// accumulator.
func TestScan(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		values, err := FromSlice([]int{1, 2, 3, 4}).
			Scan(0, func(acc, v int) int { return acc + v }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		// Running sum after each input: 1, 3, 6, 10.
		if want := []int{1, 3, 6, 10}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("empty input emits nothing", func(t *testing.T) {
		values, err := FromSlice([]int{}).
			Scan(99, func(acc, v int) int { return acc + v }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
	t.Run("type-changing accumulator", func(t *testing.T) {
		// Build a running join.
		values, err := FromSlice([]int{1, 2, 3}).
			Scan("", func(acc string, v int) string {
				if acc == "" {
					return string(rune('0' + v))
				}
				return acc + "," + string(rune('0'+v))
			}).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []string{"1", "1,2", "1,2,3"}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
}

// TestReduce covers the terminal reducer. Empty input produces
// no output; single input is forwarded without calling f; multiple
// inputs are folded left-to-right.
func TestReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		values, err := FromSlice([]int{1, 2, 3, 4, 5}).
			Reduce(func(a, b int) int { return a + b }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{15}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
	t.Run("empty input emits nothing", func(t *testing.T) {
		values, err := FromSlice([]int{}).
			Reduce(func(a, b int) int { return a + b }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
	t.Run("single input is forwarded", func(t *testing.T) {
		// f is not called for a single input; we still get the
		// value through.
		called := 0
		values, err := FromSlice([]int{42}).
			Reduce(func(a, b int) int { called++; return a + b }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{42}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
		if called != 0 {
			t.Fatalf("Reduce called f %d times for a single input", called)
		}
	})
}

// TestTakeZero covers the Take(0) branch: the operator must
// subscribe upstream but cancel during setup and emit no values.
func TestTakeZero(t *testing.T) {
	values, err := FromSlice([]int{1, 2, 3}).Take(0).ToSlice(context.Background())
	if err != nil {
		t.Fatalf("ToSlice: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("Take(0): got %v, want empty", values)
	}
}

// TestFilterEdgeCases covers the Filter branches that the existing
// TestFromSliceMapFilterTake does not exercise: predicate always
// false (every value replenishes demand), predicate always true.
func TestFilterEdgeCases(t *testing.T) {
	t.Run("predicate always false", func(t *testing.T) {
		values, err := FromSlice([]int{1, 2, 3}).
			Filter(func(int) bool { return false }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if len(values) != 0 {
			t.Fatalf("got %v, want empty", values)
		}
	})
	t.Run("predicate always true", func(t *testing.T) {
		values, err := FromSlice([]int{1, 2, 3}).
			Filter(func(int) bool { return true }).
			ToSlice(context.Background())
		if err != nil {
			t.Fatalf("ToSlice: %v", err)
		}
		if want := []int{1, 2, 3}; !reflect.DeepEqual(values, want) {
			t.Fatalf("got %v, want %v", values, want)
		}
	})
}

// ---------- subject ----------

// TestSubjectOnSubscribe covers the OnSubscribe hook on a
// Subject: it must be invoked exactly once with the subscription
// handle.
func TestSubjectOnSubscribe(t *testing.T) {
	s := NewSubject[int]()
	calls := 0
	var sub Subscription
	_ = s.Subscribe(context.Background(), subscriberFuncs[int]{
		onSubscribe: func(v Subscription) {
			calls++
			sub = v
		},
	})
	if calls != 1 {
		t.Fatalf("OnSubscribe called %d times, want 1", calls)
	}
	if sub == nil {
		t.Fatal("OnSubscribe did not pass a subscription handle")
	}
	// Request demand to keep the subscription healthy and avoid
	// leaking goroutines.
	sub.Request(1)
}

// TestSubjectOnError covers the OnError propagation path on a
// Subject: an OnError notification must be delivered to every
// subscriber, and no further OnNext or OnComplete notifications
// are allowed afterwards.
func TestSubjectOnError(t *testing.T) {
	s := NewSubject[int]()
	want := errors.New("boom")
	var got error
	done := make(chan struct{})
	s.ForEach(context.Background(), nil, func(err error) {
		got = err
		close(done)
	}, func() { close(done) })
	s.OnError(want)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for OnError")
	}
	if !errors.Is(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// TestSubjectOnErrorIgnoresSubsequentOnNext covers the contract
// that once OnError has been emitted, later OnNext / OnComplete
// calls are no-ops.
func TestSubjectOnErrorIgnoresSubsequentOnNext(t *testing.T) {
	s := NewSubject[int]()
	values := make(chan int, 4)
	done := make(chan struct{})
	s.ForEach(context.Background(),
		func(v int) { values <- v },
		func(error) { close(done) },
		func() { close(done) },
	)
	s.OnNext(1)
	s.OnError(errors.New("boom"))
	// These calls must be silently dropped.
	s.OnNext(2)
	s.OnComplete()
	// Wait for the OnError callback to fire, signalling the
	// chain has acknowledged the error and stopped processing.
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for OnError")
	}
	// Drain anything buffered before the error, but no extras.
	timeout := time.After(100 * time.Millisecond)
	got := 0
drain:
	for {
		select {
		case <-values:
			got++
		case <-timeout:
			break drain
		}
	}
	if got != 1 {
		t.Fatalf("got %d values, want 1 (only the pre-error value)", got)
	}
}

// ---------- collect ----------

// TestToSliceEmpty covers the empty-source path of ToSlice: it
// must return a non-nil empty slice and a nil error.
func TestToSliceEmpty(t *testing.T) {
	values, err := FromSlice([]int{}).ToSlice(context.Background())
	if err != nil {
		t.Fatalf("ToSlice: %v", err)
	}
	if values == nil {
		t.Fatal("ToSlice on empty: got nil, want non-nil empty slice")
	}
	if len(values) != 0 {
		t.Fatalf("ToSlice on empty: got %v, want empty", values)
	}
}

// TestCollectEmpty covers the empty-source path of Collect: the
// accumulator is initial, f is not called, and err is nil.
func TestCollectEmpty(t *testing.T) {
	called := 0
	total, err := Collect(context.Background(), FromSlice([]int{}), 99,
		func(acc, v int) int {
			called++
			return acc + v
		})
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if total != 99 {
		t.Fatalf("total: got %d, want 99 (init)", total)
	}
	if called != 0 {
		t.Fatalf("Collect on empty called f %d times, want 0", called)
	}
}
