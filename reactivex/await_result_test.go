package reactivex

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qianwj/typed/adt"
)

func TestMaybeAwaitResultStates(t *testing.T) {
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
		{"empty", 42, false, nil},
		{"error", 42, true, errSentinelMaybe},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			m := NewMaybe(func() (any, bool, error) {
				calls.Add(1)
				return tc.value, tc.present, tc.err
			})
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			// Check both the initial wait and a cached read with a canceled context.
			for _, got := range []adt.Result[adt.Option[any]]{m.Await(), m.AwaitWithContext(ctx)} {
				if tc.err != nil {
					if !got.IsFailure() || !errors.Is(got.Error(), tc.err) {
						t.Fatalf("got %v, want failure %v", got, tc.err)
					}
					continue
				}
				if got.IsFailure() || got.Value().IsPresent() != tc.present {
					t.Fatalf("got %v, want success with present=%v", got, tc.present)
				}
				if tc.present && got.Value().Get() != tc.value {
					t.Fatalf("got value %v, want %v", got.Value().Get(), tc.value)
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("source called %d times", calls.Load())
			}
		})
	}
}

func TestSingleAwaitResultStates(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value *int
		err   error
	}{
		{"nil success", nil, nil},
		{"failure", new(int), errSentinelSingle},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSingle(func() (*int, error) { return tc.value, tc.err })
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			for _, got := range []adt.Result[*int]{s.Await(), s.AwaitWithContext(ctx)} {
				value, err := got.Unwrap()
				if tc.err != nil {
					if !got.IsFailure() || !errors.Is(err, tc.err) || value != nil {
						t.Fatalf("got %v, want failure with zero value", got)
					}
				} else if !got.IsSuccess() || value != tc.value {
					t.Fatalf("got %v, want successful nil", got)
				}
			}
		})
	}
}

func TestMaybeAwaitCanceledWaitRetainsOutcome(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release := make(chan struct{})
	m := NewMaybe(func() (int, bool, error) {
		cancel()
		<-release
		return 42, true, nil
	})
	got := m.AwaitWithContext(ctx)
	close(release)
	if !got.IsFailure() || !errors.Is(got.Error(), context.Canceled) {
		t.Fatalf("got %v, want cancellation failure", got)
	}
	if final := m.Await(); final.IsFailure() || final.Value().IsEmpty() || final.Value().Get() != 42 {
		t.Fatalf("canceled wait lost eventual value: %v", final)
	}
}

func TestSingleAwaitCanceledWaitRetainsOutcome(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release := make(chan struct{})
	s := NewSingle(func() (int, error) {
		cancel()
		<-release
		return 42, nil
	})
	got := s.AwaitWithContext(ctx)
	close(release)
	if !got.IsFailure() || !errors.Is(got.Error(), context.Canceled) {
		t.Fatalf("got %v, want cancellation failure", got)
	}
	if final := s.Await(); final.IsFailure() || final.Value() != 42 {
		t.Fatalf("canceled wait lost eventual value: %v", final)
	}
}

func TestAwaitExpiredDeadlineReturnsFailure(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	s := NewSingle(func() (int, error) {
		t.Error("expired deadline started Single")
		return 0, nil
	})
	m := NewMaybe(func() (int, bool, error) {
		t.Error("expired deadline started Maybe")
		return 0, false, nil
	})
	if got := s.AwaitWithContext(ctx); !got.IsFailure() || !errors.Is(got.Error(), context.DeadlineExceeded) {
		t.Fatalf("Single: got %v, want deadline failure", got)
	}
	if got := m.AwaitWithContext(ctx); !got.IsFailure() || !errors.Is(got.Error(), context.DeadlineExceeded) {
		t.Fatalf("Maybe: got %v, want deadline failure", got)
	}
}
