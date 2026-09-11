package control

import (
	"errors"
	"testing"
)

// ---------- Repeat ----------

// TestRepeatZeroTimesIsNoOp covers the boundary case where
// times == 0: the body must not run at all. This is the
// natural Go behaviour of 'for range 0' and is what makes
// Repeat safe to call with computed counts.
func TestRepeatZeroTimesIsNoOp(t *testing.T) {
	calls := 0
	Repeat(0, func() { calls++ })
	if calls != 0 {
		t.Fatalf("Repeat(0): calls = %d, want 0", calls)
	}
}

// TestRepeatNegativeTimesIsNoOp covers the times < 0 case.
// Go's 'for range n' iterates zero times for any non-positive
// n, so Repeat(-1, ...) and Repeat(-100, ...) must not invoke f.
// Repeat has no return value, so the only thing to assert is
// that f is not called and the call does not panic.
func TestRepeatNegativeTimesIsNoOp(t *testing.T) {
	calls := 0
	Repeat(-1, func() { calls++ })
	Repeat(-100, func() { calls++ })
	if calls != 0 {
		t.Fatalf("Repeat(-1) and Repeat(-100): calls = %d, want 0", calls)
	}
}

// TestRepeatExactlyNTimes confirms that Repeat(n, f) calls f
// exactly n times when n is positive. The calls happen in
// order 0, 1, ..., n-1, which is verified by appending to a
// slice.
func TestRepeatExactlyNTimes(t *testing.T) {
	calls := 0
	Repeat(5, func() { calls++ })
	if calls != 5 {
		t.Fatalf("Repeat(5): calls = %d, want 5", calls)
	}
}

// TestRepeatOrderIsSequential confirms that f is called in
// the natural order 0, 1, ..., times-1. A test that records
// the call index in a slice must see [0, 1, 2, ..., times-1].
func TestRepeatOrderIsSequential(t *testing.T) {
	var got []int
	Repeat(4, func() {
		got = append(got, len(got))
	})
	want := []int{0, 1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("Repeat(4): got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Repeat(4) at %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

// TestRepeatSideEffectsVisible confirms that mutations made
// by f are visible to the caller after Repeat returns. This
// is the property that makes Repeat useful for "do this N
// times" side effects.
func TestRepeatSideEffectsVisible(t *testing.T) {
	counter := 0
	Repeat(3, func() { counter++ })
	if counter != 3 {
		t.Fatalf("after Repeat(3): counter = %d, want 3", counter)
	}
}

// TestRepeatClosureCapture confirms that the body can read
// and write closure variables safely across iterations. The
// captured variable grows by one each call, and the final
// value must be exactly n.
func TestRepeatClosureCapture(t *testing.T) {
	sum := 0
	Repeat(10, func() { sum++ })
	if sum != 10 {
		t.Fatalf("after Repeat(10): sum = %d, want 10", sum)
	}
}

// TestRepeatWithNilFunc confirms that a nil func does not
// panic. 'for range n' with a nil f would panic on the first
// call, so Repeat is not safe with a nil func; the test
// documents the failure mode so a future change that
// accidentally makes this safe does not silently change
// behaviour.
//
// The expected behaviour is a panic, recovered and asserted
// here. If you intend to call f zero times, use Repeat(0, f).
func TestRepeatWithNilFuncPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Repeat(1, nil): expected panic, got none")
		}
	}()
	Repeat(1, nil)
}

// ---------- RepeatE ----------

// TestRepeatEZeroTimesIsNoOp covers the times == 0 case. The
// function body is never entered and the return is (0, nil).
func TestRepeatEZeroTimesIsNoOp(t *testing.T) {
	calls := 0
	gotN, err := RepeatE(0, func() error { calls++; return nil })
	if gotN != 0 || err != nil {
		t.Fatalf("RepeatE(0): got (%d, %v), want (0, nil)", gotN, err)
	}
	if calls != 0 {
		t.Fatalf("RepeatE(0): calls = %d, want 0", calls)
	}
}

// TestRepeatENegativeTimesDoesNotPanic covers the times < 0
// case. The body is never entered and f is not called, but
// the success path returns the input value back, so
// RepeatE(-1, ...) returns (-1, nil) and RepeatE(-100, ...)
// returns (-100, nil). Negative times are a caller error;
// the function does not panic so it stays safe to call with
// computed counts, but the returned count is not a meaningful
// iteration count and callers should not rely on it.
func TestRepeatENegativeTimesDoesNotPanic(t *testing.T) {
	calls := 0
	gotN, err := RepeatE(-1, func() error { calls++; return nil })
	if gotN != -1 || err != nil {
		t.Fatalf("RepeatE(-1): got (%d, %v), want (-1, nil)", gotN, err)
	}
	if calls != 0 {
		t.Fatalf("RepeatE(-1): calls = %d, want 0", calls)
	}

	calls = 0
	gotN, err = RepeatE(-100, func() error { calls++; return nil })
	if gotN != -100 || err != nil {
		t.Fatalf("RepeatE(-100): got (%d, %v), want (-100, nil)", gotN, err)
	}
	if calls != 0 {
		t.Fatalf("RepeatE(-100): calls = %d, want 0", calls)
	}
}

// TestRepeatEAllSuccess confirms that RepeatE(n, f) returns
// (n, nil) when f returns nil every time. The (n, nil) shape
// is the success signal: the count matches the requested
// times.
func TestRepeatEAllSuccess(t *testing.T) {
	calls := 0
	gotN, err := RepeatE(7, func() error { calls++; return nil })
	if gotN != 7 || err != nil {
		t.Fatalf("RepeatE(7): got (%d, %v), want (7, nil)", gotN, err)
	}
	if calls != 7 {
		t.Fatalf("RepeatE(7): calls = %d, want 7", calls)
	}
}

// TestRepeatEErrorOnFirstCall confirms the (0, err) return
// when the very first call to f returns a non-nil error.
// f must be called exactly once.
func TestRepeatEErrorOnFirstCall(t *testing.T) {
	boom := errors.New("boom")
	calls := 0
	gotN, err := RepeatE(5, func() error {
		calls++
		return boom
	})
	if gotN != 0 || err != boom {
		t.Fatalf("RepeatE(5) error-on-first: got (%d, %v), want (0, %v)", gotN, err, boom)
	}
	if calls != 1 {
		t.Fatalf("RepeatE(5) error-on-first: calls = %d, want 1", calls)
	}
}

// TestRepeatEErrorOnMiddleCall confirms the (i, err) return
// when f returns nil for the first i calls and a non-nil
// error on the (i+1)-th call. f must be called exactly i+1
// times.
func TestRepeatEErrorOnMiddleCall(t *testing.T) {
	boom := errors.New("middle")
	calls := 0
	gotN, err := RepeatE(10, func() error {
		calls++
		if calls == 3 {
			return boom
		}
		return nil
	})
	if gotN != 2 || err != boom {
		t.Fatalf("RepeatE(10) error-on-third: got (%d, %v), want (2, %v)", gotN, err, boom)
	}
	if calls != 3 {
		t.Fatalf("RepeatE(10) error-on-third: calls = %d, want 3", calls)
	}
}

// TestRepeatEErrorOnLastCall confirms the (n-1, err) return
// when the error happens on the final call. f must be called
// exactly n times.
func TestRepeatEErrorOnLastCall(t *testing.T) {
	boom := errors.New("last")
	calls := 0
	gotN, err := RepeatE(5, func() error {
		calls++
		if calls == 5 {
			return boom
		}
		return nil
	})
	if gotN != 4 || err != boom {
		t.Fatalf("RepeatE(5) error-on-last: got (%d, %v), want (4, %v)", gotN, err, boom)
	}
	if calls != 5 {
		t.Fatalf("RepeatE(5) error-on-last: calls = %d, want 5", calls)
	}
}

// TestRepeatEStopsAfterError confirms that f is not called
// again after the first non-nil error. This is the
// "short-circuit on error" property that makes RepeatE
// useful for retry / batch loops.
func TestRepeatEStopsAfterError(t *testing.T) {
	boom := errors.New("stop here")
	calls := 0
	_, err := RepeatE(100, func() error {
		calls++
		if calls == 4 {
			return boom
		}
		return nil
	})
	if err != boom {
		t.Fatalf("RepeatE(100) stop-after-error: err = %v, want %v", err, boom)
	}
	if calls != 4 {
		t.Fatalf("RepeatE(100) stop-after-error: calls = %d, want 4 (must not call f past the first error)", calls)
	}
}

// TestRepeatENilErrorTreatedAsSuccess confirms that a
// function that returns nil does not short-circuit the loop.
// A common bug in similar primitives is to break on err == nil
// (the inverse of the intended check); this test guards
// against that.
func TestRepeatENilErrorTreatedAsSuccess(t *testing.T) {
	calls := 0
	gotN, err := RepeatE(3, func() error {
		calls++
		return nil
	})
	if gotN != 3 || err != nil {
		t.Fatalf("RepeatE(3) all-nil: got (%d, %v), want (3, nil)", gotN, err)
	}
	if calls != 3 {
		t.Fatalf("RepeatE(3) all-nil: calls = %d, want 3", calls)
	}
}

// TestRepeatERecordsCallOrder confirms that f is called in the
// natural order 0, 1, ..., n-1, and that the recorded call
// index matches the iteration that returned. This is the
// foundation of "report which iteration failed" reporting.
func TestRepeatERecordsCallOrder(t *testing.T) {
	var called []int
	_, _ = RepeatE(5, func() error {
		called = append(called, len(called))
		return nil
	})
	want := []int{0, 1, 2, 3, 4}
	if len(called) != len(want) {
		t.Fatalf("RepeatE(5) order: got %v, want %v", called, want)
	}
	for i := range called {
		if called[i] != want[i] {
			t.Fatalf("RepeatE(5) order at %d: got %d, want %d", i, called[i], want[i])
		}
	}
}

// TestRepeatEMixedNilAndErrors exercises the realistic
// "sometimes-error" pattern: f returns nil for several
// iterations, then an error, then would return another error
// if called again. RepeatE must stop at the first error and
// return the right (count, err) shape.
func TestRepeatEMixedNilAndErrors(t *testing.T) {
	first := errors.New("first failure")
	calls := 0
	gotN, err := RepeatE(10, func() error {
		calls++
		switch calls {
		case 1, 2, 3:
			return nil
		case 4:
			return first
		default:
			return errors.New("should not reach")
		}
	})
	if gotN != 3 || err != first {
		t.Fatalf("RepeatE(10) mixed: got (%d, %v), want (3, %v)", gotN, err, first)
	}
	if calls != 4 {
		t.Fatalf("RepeatE(10) mixed: calls = %d, want 4", calls)
	}
}
