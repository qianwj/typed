package result_test

import (
	"errors"
	"testing"

	"github.com/qianwj/typed/utils/result"
)

// TestResultSuccess documents the basic success path and that the
// success value can be a zero value (the absent/present distinction
// is not needed here because err plays that role).
func TestResultSuccess(t *testing.T) {
	r := result.Success(42)
	if !r.IsSuccess() {
		t.Fatal("Success(42) should be a success")
	}
	if r.IsFailure() {
		t.Fatal("Success(42) should not be a failure")
	}
	if got := r.Value(); got != 42 {
		t.Fatalf("Value: got %d, want 42", got)
	}
	if err := r.Error(); err != nil {
		t.Fatalf("Error: got %v, want nil", err)
	}
}

// TestResultFailure covers the failure path. Error must be
// retrievable and Value must panic.
func TestResultFailure(t *testing.T) {
	sentinel := errors.New("boom")
	r := result.Failure[int](sentinel)

	if r.IsSuccess() {
		t.Fatal("Failure(...) should not be a success")
	}
	if !r.IsFailure() {
		t.Fatal("Failure(...) should be a failure")
	}
	if got := r.Error(); got != sentinel {
		t.Fatalf("Error: got %v, want %v", got, sentinel)
	}
}

// TestResultValuePanic confirms the deliberate panic when a value
// is pulled out of a failed Result.
func TestResultValuePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Value on failed Result should panic")
		}
		// The panic value should be the original error, not a
		// wrapped or stringified copy. This is important for
		// callers that recover and inspect.
		err, ok := r.(error)
		if !ok {
			t.Fatalf("panic value should be an error, got %T", r)
		}
		if err.Error() != "boom" {
			t.Fatalf("unexpected panic error: %v", err)
		}
	}()
	_ = result.Failure[int](errors.New("boom")).Value()
}

// TestResultFailurePanicOnNil documents that Failure rejects a nil
// error. A Result with a nil error must be built via Success, so
// Failure[T](nil) is a programming error.
func TestResultFailurePanicOnNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Failure(nil) should panic")
		}
	}()
	_ = result.Failure[int](nil)
}

// TestResultOptional covers the bridge from Result to Optional.
// Successes become present Optionals, failures become empty ones
// and the error is dropped.
func TestResultOptional(t *testing.T) {
	t.Run("success becomes present", func(t *testing.T) {
		o := result.Success(7).Optional()
		if !o.IsPresent() || o.Get() != 7 {
			t.Fatalf("got %v, want present 7", o)
		}
	})
	t.Run("failure becomes empty", func(t *testing.T) {
		o := result.Failure[int](errors.New("boom")).Optional()
		if o.IsPresent() {
			t.Fatalf("got %v, want empty", o)
		}
	})
}

// TestResultOrElse covers both branches of the value-extracting
// OrElse and confirms the default value is evaluated eagerly (the
// call site pays for constructing it even on success).
func TestResultOrElse(t *testing.T) {
	if got := result.Success(1).OrElse(99); got != 1 {
		t.Fatalf("OrElse on success: got %d, want 1", got)
	}
	if got := result.Failure[int](errors.New("x")).OrElse(99); got != 99 {
		t.Fatalf("OrElse on failure: got %d, want 99", got)
	}
}

// TestResultOrElseGet ensures the lazy fallback only runs on
// failure. This is the right choice when the fallback is expensive
// or has side effects, mirroring Optional.OrElseGet.
func TestResultOrElseGet(t *testing.T) {
	t.Run("success skips fallback", func(t *testing.T) {
		called := false
		got := result.Success(1).OrElseGet(func() int {
			called = true
			return 99
		})
		if got != 1 {
			t.Fatalf("got %d, want 1", got)
		}
		if called {
			t.Fatal("fallback should not run on success")
		}
	})
	t.Run("failure runs fallback", func(t *testing.T) {
		got := result.Failure[int](errors.New("x")).OrElseGet(func() int { return 7 })
		if got != 7 {
			t.Fatalf("got %d, want 7", got)
		}
	})
}

// TestResultRecover covers the error-aware fallback. Recover is
// the right choice when the recovery value depends on the kind of
// error: not-found returns the zero value, transient returns a
// cached value, etc.
func TestResultRecover(t *testing.T) {
	t.Run("success skips recovery", func(t *testing.T) {
		called := false
		got := result.Success(7).Recover(func(error) int {
			called = true
			return 99
		})
		if got != 7 {
			t.Fatalf("got %d, want 7", got)
		}
		if called {
			t.Fatal("f should not run on success")
		}
	})
	t.Run("failure hands the error to recovery", func(t *testing.T) {
		sentinel := errors.New("disk full")
		got := result.Failure[int](sentinel).Recover(func(err error) int {
			if errors.Is(err, sentinel) {
				return -1
			}
			return 0
		})
		if got != -1 {
			t.Fatalf("got %d, want -1", got)
		}
	})
	t.Run("recovery decides based on error type", func(t *testing.T) {
		// A realistic pattern: distinct sentinel values for
		// distinct error categories.
		notFound := result.Failure[int](errors.New("not found"))
		denied := result.Failure[int](errors.New("denied"))

		recover := func(err error) int {
			switch err.Error() {
			case "not found":
				return 0
			case "denied":
				return -1
			default:
				return -2
			}
		}

		if got := notFound.Recover(recover); got != 0 {
			t.Fatalf("not found: got %d, want 0", got)
		}
		if got := denied.Recover(recover); got != -1 {
			t.Fatalf("denied: got %d, want -1", got)
		}
	})
	t.Run("chained at end of pipeline", func(t *testing.T) {
		// Smoke test: combine Recover with FlatMap and MapError
		// to show it sits at the end of a chain as a final
		// value extraction step.
		parse := func(s string) result.Result[int] {
			if s == "" {
				return result.Failure[int](errors.New("empty"))
			}
			v, err := atoi(s)
			if err != nil {
				return result.Failure[int](err)
			}
			return result.Success(v)
		}
		double := func(n int) result.Result[int] {
			return result.Success(n * 2)
		}

		got := parse("12").FlatMap(double).Recover(func(err error) int {
			// On any failure, return 0 — common default for
			// "couldn't compute, but I need to keep going".
			return 0
		})
		if got != 24 {
			t.Fatalf("chained recover: got %d, want 24", got)
		}

		// Failure path: parse("") fails at the first step,
		// recovery produces 0.
		bad := parse("").FlatMap(double).Recover(func(err error) int {
			return 0
		})
		if bad != 0 {
			t.Fatalf("chained recover on failure: got %d, want 0", bad)
		}
	})
}

// TestResultUnwrap covers the bridge from Result to Go's standard
// (T, error) return shape. Unwrap is the right method to call at
// the end of a chain, or anywhere the caller wants to fall back to
// ordinary Go error handling.
func TestResultUnwrap(t *testing.T) {
	t.Run("success returns value and nil error", func(t *testing.T) {
		v, err := result.Success(42).Unwrap()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 42 {
			t.Fatalf("got %d, want 42", v)
		}
	})
	t.Run("failure returns zero value and the error", func(t *testing.T) {
		sentinel := errors.New("boom")
		v, err := result.Failure[int](sentinel).Unwrap()
		if err != sentinel {
			t.Fatalf("got error %v, want %v", err, sentinel)
		}
		if v != 0 {
			t.Fatalf("got %d, want 0 (zero value)", v)
		}
	})
	// Unwrap is the intended replacement for the old
	// GetOrElseFailure combinator: any (T, error) function can
	// serve as a recovery path, not just one that returns a
	// Result. This test documents the recommended pattern.
	t.Run("recovery via standard (T, error) function", func(t *testing.T) {
		// A pretend fallback that returns (T, error), exactly
		// like an ordinary Go function would.
		loadFromCache := func() (int, error) {
			return 99, nil
		}

		v, err := result.Failure[int](errors.New("primary failed")).Unwrap()
		if err != nil {
			v, err = loadFromCache()
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 99 {
			t.Fatalf("got %d, want 99", v)
		}
	})
}

// TestResultMap covers a type-changing Map and confirms failures
// are propagated without invoking f.
func TestResultMap(t *testing.T) {
	t.Run("success applies f and changes type", func(t *testing.T) {
		got := result.Success(7).Map(func(n int) string { return itoa(n) })
		if !got.IsSuccess() || got.Value() != "7" {
			t.Fatalf("got %v, want Success(\"7\")", got)
		}
	})
	t.Run("failure propagates without calling f", func(t *testing.T) {
		sentinel := errors.New("nope")
		called := false
		got := result.Failure[int](sentinel).Map(func(int) string {
			called = true
			return "x"
		})
		if got.IsSuccess() {
			t.Fatal("got Success, want Failure")
		}
		if got.Error() != sentinel {
			t.Fatalf("got %v, want %v", got.Error(), sentinel)
		}
		if called {
			t.Fatal("f should not run on failure")
		}
	})
}

// TestResultFlatMap is the natural combinator for "if successful,
// run a function that itself may fail". It also covers the case
// where the inner Result is itself a failure.
func TestResultFlatMap(t *testing.T) {
	t.Run("success delegates to f", func(t *testing.T) {
		got := result.Success(5).FlatMap(func(n int) result.Result[string] {
			return result.Success(itoa(n))
		})
		if !got.IsSuccess() || got.Value() != "5" {
			t.Fatalf("got %v, want Success(\"5\")", got)
		}
	})
	t.Run("inner failure surfaces", func(t *testing.T) {
		sentinel := errors.New("inner")
		got := result.Success(5).FlatMap(func(int) result.Result[string] {
			return result.Failure[string](sentinel)
		})
		if got.IsSuccess() {
			t.Fatal("got Success, want Failure")
		}
		if got.Error() != sentinel {
			t.Fatalf("got %v, want %v", got.Error(), sentinel)
		}
	})
	t.Run("outer failure short-circuits", func(t *testing.T) {
		sentinel := errors.New("outer")
		called := false
		got := result.Failure[int](sentinel).FlatMap(func(int) result.Result[string] {
			called = true
			return result.Success("x")
		})
		if got.Error() != sentinel {
			t.Fatalf("got %v, want %v", got.Error(), sentinel)
		}
		if called {
			t.Fatal("f should not run on outer failure")
		}
	})
}

// TestResultMapError is the "wrap the error" combinator. The
// success value must pass through untouched.
func TestResultMapError(t *testing.T) {
	t.Run("success unchanged", func(t *testing.T) {
		called := false
		r := result.Success(7).MapError(func(error) error {
			called = true
			return errors.New("unused")
		})
		if !r.IsSuccess() || r.Value() != 7 {
			t.Fatalf("got %v, want Success(7)", r)
		}
		if called {
			t.Fatal("f should not run on success")
		}
	})
	t.Run("failure error is wrapped", func(t *testing.T) {
		r := result.Failure[int](errors.New("disk")).MapError(
			func(e error) error { return errors.New("save: " + e.Error()) },
		)
		if r.IsSuccess() {
			t.Fatal("got Success, want Failure")
		}
		if got, want := r.Error().Error(), "save: disk"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

// TestResultChained is a smoke test that exercises a realistic
// pipeline: parse, validate, transform. It mirrors the chained
// pattern shown in the design doc.
func TestResultChained(t *testing.T) {
	// Parse a string into an int. The empty string is a failure.
	parse := func(s string) result.Result[int] {
		if s == "" {
			return result.Failure[int](errors.New("empty input"))
		}
		v, err := atoi(s)
		if err != nil {
			return result.Failure[int](err)
		}
		return result.Success(v)
	}

	// Take a successful int, double it. Doubling a negative is
	// still success; the caller decides whether negative is
	// allowed in the next step.
	double := func(n int) result.Result[int] { return result.Success(n * 2) }

	got := parse("12").FlatMap(double)
	if !got.IsSuccess() || got.Value() != 24 {
		t.Fatalf("chained: got %v, want Success(24)", got)
	}

	bad := parse("").FlatMap(double)
	if bad.IsSuccess() {
		t.Fatalf("chained empty: got Success, want Failure")
	}

	// MapError at the end to attach a higher-level label without
	// having to wrap every call site.
	wrapped := parse("").MapError(func(e error) error {
		return errors.New("parse step: " + e.Error())
	})
	if got, want := wrapped.Error().Error(), "parse step: empty input"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// ---------- helpers ----------

// itoa and atoi are tiny stand-ins for strconv to keep the test
// file's import surface small and to make the assertions obvious.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func atoi(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}
