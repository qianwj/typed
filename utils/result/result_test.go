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

// ---------- Wrap ----------

// TestResultWrapSuccessNilErr confirms that Wrap(value, nil)
// produces a success carrying value, the same as
// Success(value).
func TestResultWrapSuccessNilErr(t *testing.T) {
	r := result.Wrap(42, nil)
	if !r.IsSuccess() {
		t.Fatal("Wrap(42, nil): IsSuccess = false, want true")
	}
	if r.IsFailure() {
		t.Fatal("Wrap(42, nil): IsFailure = true, want false")
	}
	if got := r.Value(); got != 42 {
		t.Fatalf("Wrap(42, nil).Value: got %d, want 42", got)
	}
	if err := r.Error(); err != nil {
		t.Fatalf("Wrap(42, nil).Error: got %v, want nil", err)
	}
}

// TestResultWrapFailureNonNilErr confirms that Wrap(_, err)
// with a non-nil err produces a failure carrying err. The
// value parameter is stored internally but not visible
// through the Result API.
func TestResultWrapFailureNonNilErr(t *testing.T) {
	sentinel := errors.New("wrap failure")
	r := result.Wrap(99, sentinel)
	if r.IsSuccess() {
		t.Fatal("Wrap(99, err): IsSuccess = true, want false")
	}
	if !r.IsFailure() {
		t.Fatal("Wrap(99, err): IsFailure = false, want true")
	}
	if err := r.Error(); err != sentinel {
		t.Fatalf("Wrap(99, err).Error: got %v, want %v", err, sentinel)
	}
}

// TestResultWrapZeroValueSuccess confirms that Wrap(0, nil)
// produces a success carrying the zero value, not a failure.
// This is the case where Wrap differs from Failure: a zero
// value with a nil err is unambiguously a success.
func TestResultWrapZeroValueSuccess(t *testing.T) {
	r := result.Wrap(0, nil)
	if !r.IsSuccess() {
		t.Fatal("Wrap(0, nil): IsSuccess = false, want true (zero value with nil err is success)")
	}
	if got := r.Value(); got != 0 {
		t.Fatalf("Wrap(0, nil).Value: got %d, want 0", got)
	}
}

// TestResultWrapObservationalIgnoresValueOnFailure confirms
// that the value parameter is unreachable through the Result
// API when err is non-nil. This is the property that lets
// callers write `return result.Wrap(r, err)` after a (T, error)
// call without first checking which branch they are in.
func TestResultWrapObservationalIgnoresValueOnFailure(t *testing.T) {
	err := errors.New("ignored-value test")
	withIgnored := result.Wrap(42, err)
	withZero := result.Wrap(0, err)
	// Both must be observationally identical: same error,
	// same Unwrap output, same Optional, same Map / OrElse.
	if withIgnored.Error() != withZero.Error() {
		t.Fatal("ignored vs zero: errors differ")
	}
	if vIgnored, eIgnored := withIgnored.Unwrap(); eIgnored != err || vIgnored != 0 {
		t.Fatalf("withIgnored.Unwrap: got (%v, %v), want (0, err)", vIgnored, eIgnored)
	}
	if vZero, eZero := withZero.Unwrap(); eZero != err || vZero != 0 {
		t.Fatalf("withZero.Unwrap: got (%v, %v), want (0, err)", vZero, eZero)
	}
	if withIgnored.Optional().IsPresent() || withZero.Optional().IsPresent() {
		t.Fatal("Optional of a failure must be absent, regardless of value")
	}
	if withIgnored.OrElse(7) != withZero.OrElse(7) {
		t.Fatal("OrElse of failures must agree, regardless of value")
	}
	if withIgnored.OrElseGet(func() int { return 7 }) != withZero.OrElseGet(func() int { return 7 }) {
		t.Fatal("OrElseGet of failures must agree, regardless of value")
	}
}

// TestResultWrapMatchesSuccessAndFailure confirms that Wrap
// is exactly equivalent to Success on the success path and
// exactly equivalent to Failure on the failure path. The
// comparison is made through the public methods, so this
// also serves as a regression test for the internal storage
// invariants.
func TestResultWrapMatchesSuccessAndFailure(t *testing.T) {
	cases := []struct {
		name    string
		wrapped result.Result[int]
		plain   result.Result[int]
	}{
		{"nil err, value 42", result.Wrap(42, nil), result.Success(42)},
		{"nil err, value 0", result.Wrap(0, nil), result.Success(0)},
		{"err, value 99", result.Wrap(99, errors.New("x")), result.Failure[int](errors.New("x"))},
		{"err, value 0", result.Wrap(0, errors.New("x")), result.Failure[int](errors.New("x"))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.wrapped.IsSuccess() != c.plain.IsSuccess() {
				t.Fatalf("IsSuccess differs: %v vs %v", c.wrapped.IsSuccess(), c.plain.IsSuccess())
			}
			if c.wrapped.IsFailure() != c.plain.IsFailure() {
				t.Fatalf("IsFailure differs")
			}
			// Error must agree (or both be nil).
			wErr, pErr := c.wrapped.Error(), c.plain.Error()
			if (wErr == nil) != (pErr == nil) {
				t.Fatalf("Error nilness differs: %v vs %v", wErr, pErr)
			}
		})
	}
}

// TestResultWrapWithUnwrap confirms that Wrap is the forward
// half of the (T, error) bridge and Unwrap is the reverse
// half: Wrap + Unwrap round-trips.
func TestResultWrapWithUnwrap(t *testing.T) {
	v, err := result.Wrap(42, nil).Unwrap()
	if err != nil {
		t.Fatalf("Unwrap on success: err = %v, want nil", err)
	}
	if v != 42 {
		t.Fatalf("Unwrap on success: v = %d, want 42", v)
	}

	sentinel := errors.New("bridge test")
	v, err = result.Wrap(0, sentinel).Unwrap()
	if err != sentinel {
		t.Fatalf("Unwrap on failure: err = %v, want %v", err, sentinel)
	}
	if v != 0 {
		t.Fatalf("Unwrap on failure: v = %d, want 0", v)
	}
}

// TestResultWrapComposesWithMap confirms that a Wrap-built
// result participates in the standard Map chain. Map on a
// success transforms the value; Map on a failure propagates
// the error.
func TestResultWrapComposesWithMap(t *testing.T) {
	success := result.Wrap(21, nil).Map(func(n int) int { return n * 2 })
	if v := success.Value(); v != 42 {
		t.Fatalf("Map on success: got %d, want 42", v)
	}

	failure := result.Wrap(0, errors.New("map test")).
		Map(func(n int) int { return n * 2 })
	if !failure.IsFailure() {
		t.Fatal("Map on failure: should remain a failure")
	}
}

// TestResultWrapComposesWithOrElse confirms that Wrap-built
// results participate in OrElse and OrElseGet.
func TestResultWrapComposesWithOrElse(t *testing.T) {
	if v := result.Wrap(42, nil).OrElse(0); v != 42 {
		t.Fatalf("OrElse on success: got %d, want 42", v)
	}
	if v := result.Wrap(0, errors.New("orelse test")).OrElse(99); v != 99 {
		t.Fatalf("OrElse on failure: got %d, want 99 (fallback)", v)
	}
	if v := result.Wrap(0, errors.New("orelseget test")).
		OrElseGet(func() int { return 100 }); v != 100 {
		t.Fatalf("OrElseGet on failure: got %d, want 100", v)
	}
}

// TestResultWrapComposesWithRecover confirms that the
// error-aware fallback chain works on a Wrap-built failure:
// Recover runs the fallback and returns the fallback T.
func TestResultWrapComposesWithRecover(t *testing.T) {
	got := result.Wrap(0, errors.New("recover test")).
		Recover(func(err error) int {
			if err == nil {
				t.Fatal("Recover: f called with nil error")
			}
			return -1
		})
	if got != -1 {
		t.Fatalf("Recover: got %d, want -1 (fallback)", got)
	}

	// Recover on a Wrap-built success returns the original
	// value, untouched.
	got = result.Wrap(42, nil).
		Recover(func(err error) int {
			t.Fatal("Recover: f called on success")
			return 0
		})
	if got != 42 {
		t.Fatalf("Recover on success: got %d, want 42", got)
	}
}

// TestResultWrapComposesWithOptional confirms that
// Optional() collapses a Wrap-built result the same way it
// collapses Success / Failure: success → present Optional,
// failure → absent Optional (error dropped).
func TestResultWrapComposesWithOptional(t *testing.T) {
	present := result.Wrap(42, nil).Optional()
	if !present.IsPresent() {
		t.Fatal("Optional on success: should be present")
	}
	if v := present.OrElse(0); v != 42 {
		t.Fatalf("Optional on success: got %d, want 42", v)
	}

	absent := result.Wrap(0, errors.New("optional test")).Optional()
	if absent.IsPresent() {
		t.Fatal("Optional on failure: should be absent")
	}
}

// TestResultWrapPracticalAdapter documents the realistic
// shape that Wrap is for: a function that returns (T, error)
// and is being plugged into a Result chain. The test confirms
// the adapter does what the doc comment claims — replaces
// the if-err ladder with a single return.
func TestResultWrapPracticalAdapter(t *testing.T) {
	// Source function: returns (int, error), the standard
	// (T, error) shape. We will call it twice and route
	// each result through a Result chain via Wrap.
	double := func(n int) (int, error) {
		if n < 0 {
			return 0, errors.New("negative input")
		}
		return n * 2, nil
	}

	// Wrap the success into a chain.
	okResult := result.Wrap(double(21)) // 42, nil
	if v := okResult.Map(func(n int) int { return n + 1 }).Value(); v != 43 {
		t.Fatalf("Wrap(success).Map: got %d, want 43", v)
	}

	// Wrap the failure into a chain that falls back.
	badResult := result.Wrap(double(-1)) // 0, error
	if v := badResult.OrElse(99); v != 99 {
		t.Fatalf("Wrap(failure).OrElse: got %d, want 99 (fallback)", v)
	}
}
