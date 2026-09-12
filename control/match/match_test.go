package match

import (
	"fmt"
	"testing"
)

func TestValueMatcherFirstMatchWins(t *testing.T) {
	calls := 0

	got := Match(12).
		When(func(n int) bool {
			calls++
			return n > 10
		}, func(n int) string {
			calls++
			return fmt.Sprintf("large:%d", n)
		}).
		When(func(int) bool {
			t.Fatal("later predicate should not run after a match")
			return true
		}, func(int) string {
			t.Fatal("later handler should not run after a match")
			return ""
		}).
		Default(func(int) string {
			t.Fatal("default should not run after a match")
			return ""
		})

	if got != "large:12" {
		t.Fatalf("got %q, want %q", got, "large:12")
	}
	if calls != 2 {
		t.Fatalf("calls got %d, want 2", calls)
	}
}

func TestValueMatcherPatternHelpers(t *testing.T) {
	grade := Value(86).
		Case(Between(90, 100), func(int) string { return "A" }).
		Case(Between(80, 89), func(int) string { return "B" }).
		Case(Between(70, 79), func(int) string { return "C" }).
		Default(func(int) string { return "F" })

	if grade != "B" {
		t.Fatalf("grade got %q, want B", grade)
	}

	status := Value("paused").
		Case(In("new", "queued"), func(string) string { return "pending" }).
		Case(Or(Eq("paused"), Eq("suspended")), func(string) string { return "stopped" }).
		Case(Not(Any[string]()), func(string) string { return "impossible" }).
		OrElse("unknown")

	if status != "stopped" {
		t.Fatalf("status got %q, want stopped", status)
	}
}

func TestValueMatcherSupportsNonComparableValues(t *testing.T) {
	got := Value([]int{1, 2, 3}).
		When(func(values []int) bool {
			return len(values) > 2
		}, func(values []int) int {
			return len(values)
		}).
		OrElse(0)

	if got != 3 {
		t.Fatalf("got %d, want 3", got)
	}
}

func TestEqMatchesStructs(t *testing.T) {
	type profile struct {
		ID   int
		Tags []string
	}

	want := profile{ID: 7, Tags: []string{"go", "generics"}}
	got := Value(profile{ID: 7, Tags: []string{"go", "generics"}}).
		Case(Eq(want), func(profile) string {
			return "same"
		}).
		Default(func(profile) string {
			return "different"
		})

	if got != "same" {
		t.Fatalf("got %q, want same", got)
	}

	different := Value(profile{ID: 7, Tags: []string{"go"}}).
		Case(Eq(want), func(profile) string {
			return "same"
		}).
		OrElse("different")

	if different != "different" {
		t.Fatalf("got %q, want different", different)
	}
}

func TestValueMatcherUnwrap(t *testing.T) {
	matched := Value(3).
		Case(Eq(3), func(int) string { return "three" })

	if got, ok := matched.Unwrap(); !ok || got != "three" {
		t.Fatalf("matched unwrap got (%q, %v), want (three, true)", got, ok)
	}
	if !matched.Matched() {
		t.Fatalf("Matched got false, want true")
	}

	unmatched := Value(3).
		Case(Eq(4), func(int) string { return "four" })

	if got, ok := unmatched.Unwrap(); ok || got != "" {
		t.Fatalf("unmatched unwrap got (%q, %v), want zero false", got, ok)
	}
	if unmatched.Matched() {
		t.Fatalf("Matched got true, want false")
	}
}

func TestValueMatcherLazyFallback(t *testing.T) {
	calls := 0

	got := Value(3).
		Case(Eq(4), func(int) string { return "four" }).
		OrElseGet(func(n int) string {
			calls++
			return fmt.Sprintf("fallback:%d", n)
		})

	if got != "fallback:3" {
		t.Fatalf("got %q, want fallback:3", got)
	}
	if calls != 1 {
		t.Fatalf("fallback calls got %d, want 1", calls)
	}
}

func TestTypeMatcherFirstMatchWins(t *testing.T) {
	var value any = "typed"
	calls := 0

	got := Type(value).
		Case(func(int) string {
			t.Fatal("int case should not match")
			return ""
		}).
		Case(func(s string) string {
			calls++
			return "string:" + s
		}).
		Case(func(any) string {
			t.Fatal("any case should not run after a match")
			return ""
		}).
		Default(func(any) string {
			t.Fatal("default should not run after a match")
			return ""
		})

	if got != "string:typed" {
		t.Fatalf("got %q, want string:typed", got)
	}
	if calls != 1 {
		t.Fatalf("calls got %d, want 1", calls)
	}
}

func TestTypeMatcherCaseWhenAndInterface(t *testing.T) {
	var value any = timeoutError{}

	got := Type(value).
		CaseWhen(func(err interface{ Timeout() bool }) bool {
			return err.Timeout()
		}, func(errorWithTimeout interface{ Timeout() bool }) string {
			if !errorWithTimeout.Timeout() {
				t.Fatal("unexpected non-timeout error")
			}
			return "timeout"
		}).
		Default(func(any) string { return "other" })

	if got != "timeout" {
		t.Fatalf("got %q, want timeout", got)
	}
}

func TestTypeMatcherNil(t *testing.T) {
	var pointer *int

	got := Type(pointer).
		Nil(func() string { return "nil" }).
		Case(func(*int) string { return "pointer" }).
		Default(func(any) string { return "other" })

	if got != "nil" {
		t.Fatalf("got %q, want nil", got)
	}

	typedNilCase := Type(pointer).
		Case(func(value *int) string {
			if value != nil {
				t.Fatal("got non-nil pointer")
			}
			return "typed nil pointer"
		}).
		Nil(func() string { return "nil" }).
		Default(func(any) string { return "other" })

	if typedNilCase != "typed nil pointer" {
		t.Fatalf("got %q, want typed nil pointer", typedNilCase)
	}
}

func TestTypeMatcherUnwrap(t *testing.T) {
	matched := Type(10).
		Case(func(n int) string { return fmt.Sprintf("%d", n) })

	if got, ok := matched.Unwrap(); !ok || got != "10" {
		t.Fatalf("matched unwrap got (%q, %v), want (10, true)", got, ok)
	}

	unmatched := Type(10).
		Case(func(string) string { return "string" })

	if got, ok := unmatched.Unwrap(); ok || got != "" {
		t.Fatalf("unmatched unwrap got (%q, %v), want zero false", got, ok)
	}
}

type timeoutError struct{}

func (timeoutError) Timeout() bool {
	return true
}

// ---------- gap-fill coverage for the match package ----------
//
// The tests in this block cover the function paths that the original
// match_test.go suite did not exercise: every 0% function in the
// package coverage report, plus the no-match fallback branches of
// the Chain methods. They are written in the same style as the rest
// of the file (no subtests with t.Run, focused assertions).

// TestMatcherDefaultWithoutCases covers the standalone Default /
// OrElse / OrElseGet methods on Matcher, which run when no prior
// Case has been declared.
func TestMatcherDefaultWithoutCases(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		got := Value(42).Default(func(n int) string { return fmt.Sprintf("v=%d", n) })
		if got != "v=42" {
			t.Fatalf("got %q, want v=42", got)
		}
	})
	t.Run("OrElse", func(t *testing.T) {
		got := Value(42).OrElse("missing")
		if got != "missing" {
			t.Fatalf("got %q, want missing", got)
		}
	})
	t.Run("OrElseGet", func(t *testing.T) {
		calls := 0
		got := Value(42).OrElseGet(func(n int) string {
			calls++
			return fmt.Sprintf("v=%d", n)
		})
		if got != "v=42" {
			t.Fatalf("got %q, want v=42", got)
		}
		if calls != 1 {
			t.Fatalf("OrElseGet calls got %d, want 1", calls)
		}
	})
}

// TestMatcherDefaultPanics covers the nil-handler defensive panic
// for the standalone Default / OrElseGet on Matcher.
func TestMatcherDefaultPanics(t *testing.T) {
	t.Run("Default nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil Default handler")
			}
		}()
		Value(1).Default[string](nil)
	})
	t.Run("OrElseGet nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil OrElseGet handler")
			}
		}()
		Value(1).OrElseGet[string](nil)
	})
}

// TestTypeMatcherDefaultWithoutCases covers the standalone Default /
// OrElse / OrElseGet methods on TypeMatcher, which run when no
// prior Case has been declared.
func TestTypeMatcherDefaultWithoutCases(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		got := Type(123).Default(func(v any) string { return fmt.Sprintf("v=%v", v) })
		if got != "v=123" {
			t.Fatalf("got %q, want v=123", got)
		}
	})
	t.Run("OrElse", func(t *testing.T) {
		got := Type(123).OrElse("missing")
		if got != "missing" {
			t.Fatalf("got %q, want missing", got)
		}
	})
	t.Run("OrElseGet", func(t *testing.T) {
		calls := 0
		got := Type(123).OrElseGet(func(v any) string {
			calls++
			return fmt.Sprintf("v=%v", v)
		})
		if got != "v=123" {
			t.Fatalf("got %q, want v=123", got)
		}
		if calls != 1 {
			t.Fatalf("OrElseGet calls got %d, want 1", calls)
		}
	})
}

// TestTypeMatcherDefaultPanics covers the nil-handler defensive
// panic for the standalone Default / OrElseGet on TypeMatcher.
func TestTypeMatcherDefaultPanics(t *testing.T) {
	t.Run("Default nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil Default handler")
			}
		}()
		Type(1).Default[string](nil)
	})
	t.Run("OrElseGet nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil OrElseGet handler")
			}
		}()
		Type(1).OrElseGet[string](nil)
	})
}

// TestTypeMatcherWhen covers the predicate-based When method on
// TypeMatcher / TypeChain, which matches on an arbitrary predicate
// over the raw input value rather than on a type assertion.
func TestTypeMatcherWhen(t *testing.T) {
	t.Run("predicate true", func(t *testing.T) {
		got := Type(15).
			When(func(v any) bool { return v.(int) > 10 }, func(v any) string {
				return fmt.Sprintf("big:%d", v.(int))
			}).
			OrElse("small")
		if got != "big:15" {
			t.Fatalf("got %q, want big:15", got)
		}
	})
	t.Run("predicate false", func(t *testing.T) {
		got := Type(3).
			When(func(v any) bool { return v.(int) > 10 }, func(v any) string {
				return "big"
			}).
			OrElse("small")
		if got != "small" {
			t.Fatalf("got %q, want small", got)
		}
	})
	t.Run("first match wins", func(t *testing.T) {
		calls := 0
		got := Type(15).
			When(func(v any) bool { return v.(int) > 0 }, func(v any) string {
				calls++
				return "first"
			}).
			When(func(v any) bool { return v.(int) > 0 }, func(v any) string {
				calls++
				return "second"
			}).
			OrElse("none")
		if got != "first" {
			t.Fatalf("got %q, want first", got)
		}
		if calls != 1 {
			t.Fatalf("When calls got %d, want 1", calls)
		}
	})
	t.Run("nil predicate panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil When predicate")
			}
		}()
		Type(1).When[string](nil, func(any) string { return "" })
	})
	t.Run("nil handler panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil When handler")
			}
		}()
		Type(1).When[string](func(any) bool { return true }, nil)
	})
}

// TestTypeChainCaseWhen covers the type assertion + guard branch on
// TypeChain.CaseWhen.
func TestTypeChainCaseWhen(t *testing.T) {
	t.Run("type and guard both match", func(t *testing.T) {
		var value any = timeoutError{}
		got := Type(value).
			CaseWhen(func(err timeoutError) bool { return err.Timeout() },
				func(err timeoutError) string { return "timeout" }).
			OrElse("other")
		if got != "timeout" {
			t.Fatalf("got %q, want timeout", got)
		}
	})
	t.Run("type matches but guard fails", func(t *testing.T) {
		var value any = timeoutError{}
		got := Type(value).
			CaseWhen(func(err timeoutError) bool { return false },
				func(err timeoutError) string { return "timeout" }).
			OrElse("other")
		if got != "other" {
			t.Fatalf("got %q, want other", got)
		}
	})
	t.Run("type does not match", func(t *testing.T) {
		var value any = 42
		got := Type(value).
			CaseWhen(func(err timeoutError) bool { return true },
				func(err timeoutError) string { return "timeout" }).
			OrElse("other")
		if got != "other" {
			t.Fatalf("got %q, want other", got)
		}
	})
	t.Run("first match wins", func(t *testing.T) {
		var value any = timeoutError{}
		calls := 0
		got := Type(value).
			CaseWhen(func(err timeoutError) bool {
				calls++
				return true
			}, func(err timeoutError) string { return "first" }).
			CaseWhen(func(err timeoutError) bool {
				calls++
				return true
			}, func(err timeoutError) string { return "second" }).
			OrElse("other")
		if got != "first" {
			t.Fatalf("got %q, want first", got)
		}
		if calls != 1 {
			t.Fatalf("CaseWhen calls got %d, want 1", calls)
		}
	})
	t.Run("nil guard panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil CaseWhen guard")
			}
		}()
		Type(timeoutError{}).CaseWhen[timeoutError, string](nil,
			func(timeoutError) string { return "" })
	})
	t.Run("nil handler panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil CaseWhen handler")
			}
		}()
		Type(timeoutError{}).CaseWhen[timeoutError, string](
			func(timeoutError) bool { return true }, nil)
	})
}

// TestTypeChainFallbacks covers the no-match fallback branches of
// TypeChain.Default / OrElse / OrElseGet / Matched / Unwrap. The
// matched branches are exercised by other tests; only the fallback
// paths were missing.
func TestTypeChainFallbacks(t *testing.T) {
	t.Run("Default fallback", func(t *testing.T) {
		got := Type(123).
			Case(func(s string) string { return "string" }).
			Default(func(v any) string { return fmt.Sprintf("v=%v", v) })
		if got != "v=123" {
			t.Fatalf("got %q, want v=123", got)
		}
	})
	t.Run("OrElse fallback", func(t *testing.T) {
		got := Type(123).
			Case(func(s string) string { return "string" }).
			OrElse("missing")
		if got != "missing" {
			t.Fatalf("got %q, want missing", got)
		}
	})
	t.Run("OrElseGet fallback", func(t *testing.T) {
		calls := 0
		got := Type(123).
			Case(func(s string) string { return "string" }).
			OrElseGet(func(v any) string {
				calls++
				return fmt.Sprintf("v=%v", v)
			})
		if got != "v=123" {
			t.Fatalf("got %q, want v=123", got)
		}
		if calls != 1 {
			t.Fatalf("OrElseGet calls got %d, want 1", calls)
		}
	})
	t.Run("Matched reports false", func(t *testing.T) {
		chain := Type(123).Case(func(s string) string { return "string" })
		if chain.Matched() {
			t.Fatal("Matched: expected false on unmatched chain")
		}
	})
	t.Run("Unwrap on no match", func(t *testing.T) {
		chain := Type(123).Case(func(s string) string { return "string" })
		got, ok := chain.Unwrap()
		if ok {
			t.Fatal("Unwrap: expected ok=false on unmatched chain")
		}
		if got != "" {
			t.Fatalf("Unwrap: got %q, want zero value", got)
		}
	})
}

// TestTypeChainFallbacksPanics covers the nil-handler panic paths
// for the no-match fallback branches of Default / OrElseGet on
// TypeChain.
func TestTypeChainFallbacksPanics(t *testing.T) {
	t.Run("Default nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil Default handler")
			}
		}()
		// Bind to a typed var so the compiler infers R = string
		// without needing a non-nil handler to anchor the type.
		var f func(any) string
		var result string
		result = Type(123).
			Case(func(s string) string { return "" }).
			Default(f)
		_ = result
	})
	t.Run("OrElseGet nil", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil OrElseGet handler")
			}
		}()
		var f func(any) string
		var result string
		result = Type(123).
			Case(func(s string) string { return "" }).
			OrElseGet(f)
		_ = result
	})
}

// TestAndPattern covers the And pattern helper, which matches only
// when every supplied pattern matches.
func TestAndPattern(t *testing.T) {
	t.Run("all match", func(t *testing.T) {
		p := And(Eq(15), Between(10, 20))
		if !p.Match(15) {
			t.Fatal("And(Eq(15), Between(10, 20)) should match 15")
		}
	})
	t.Run("one does not match", func(t *testing.T) {
		p := And(Eq(15), Between(10, 20))
		if p.Match(5) {
			t.Fatal("And(Eq(15), Between(10, 20)) should not match 5")
		}
	})
	t.Run("no patterns matches everything", func(t *testing.T) {
		// And() with no patterns is the universal pattern.
		p := And[int]()
		if !p.Match(0) || !p.Match(99) || !p.Match(-1) {
			t.Fatal("And() with no patterns should match every value")
		}
	})
	t.Run("nil pattern panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil pattern inside And")
			}
		}()
		p := And(Eq(15), nil)
		p.Match(15)
	})
}

// TestNotPattern covers the negation branch of Not: Match returns
// true exactly when the inner pattern returns false. The existing
// tests use Not with Any as a way to test "all other cases", which
// only exercises the false branch of the inner. Here we exercise
// the full match cycle.
func TestNotPattern(t *testing.T) {
	t.Run("inner false -> Not true", func(t *testing.T) {
		p := Not(Eq(5))
		if !p.Match(10) {
			t.Fatal("Not(Eq(5)) should match 10")
		}
	})
	t.Run("inner true -> Not false", func(t *testing.T) {
		p := Not(Eq(5))
		if p.Match(5) {
			t.Fatal("Not(Eq(5)) should not match 5")
		}
	})
	t.Run("nil pattern panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil pattern inside Not")
			}
		}()
		_ = Not[int](nil)
	})
}

// TestPatternFuncNilPanics covers the nil-defensive panic inside
// PatternFunc.Match: a nil PatternFunc value must panic rather than
// silently return false.
func TestPatternFuncNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil PatternFunc")
		}
	}()
	var p PatternFunc[int]
	_ = p.Match(1)
}

// TestPredicateNilPanics covers the nil-defensive panic inside
// Predicate: passing a nil predicate must panic, not silently
// return a PatternFunc that always returns false.
func TestPredicateNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil Predicate function")
		}
	}()
	_ = Predicate[int](nil)
}

// TestRequirePatternCoversAllNilPaths exercises the nil-pattern
// defensive panic in the Case and applyValueCase paths. The
// helpers Or / And / Not all call requirePattern; here we trigger
// the panic directly via a Case with a nil pattern.
func TestRequirePatternCoversAllNilPaths(t *testing.T) {
	t.Run("Case with nil pattern panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil Case pattern")
			}
		}()
		var p Pattern[int]
		_ = Value(1).Case(p, func(int) string { return "" })
	})
	t.Run("Chain Case with nil pattern panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil Chain Case pattern")
			}
		}()
		var p Pattern[int]
		_ = Value(1).Case(Eq(2), func(int) string { return "two" }).
			Case(p, func(int) string { return "" })
	})
	t.Run("Or with nil pattern panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for nil pattern inside Or")
			}
		}()
		// Or(Eq(1), nil): the first pattern matches 1 and short-circuits
		// before the nil is reached. Use a value that does not match
		// the first pattern, so the loop advances to the nil entry.
		p := Or(Eq(1), nil)
		p.Match(99)
	})
}

// TestChainOrElseGetMatchedLazy verifies that OrElseGet does not
// call the fallback function when a case has already matched.
// Without this, the lazy-fallback guarantee is not enforced.
func TestChainOrElseGetMatchedLazy(t *testing.T) {
	calls := 0
	got := Value(3).
		Case(Eq(3), func(int) string { return "three" }).
		OrElseGet(func(int) string {
			calls++
			return "fallback"
		})
	if got != "three" {
		t.Fatalf("got %q, want three", got)
	}
	if calls != 0 {
		t.Fatalf("OrElseGet calls got %d, want 0 (matched path is lazy)", calls)
	}
}

// TestChainMatchedReportsTrue covers the matched-true branch of
// Chain.Matched, which the existing TestValueMatcherUnwrap already
// covers for the value path. Here we also exercise the type
// chain.
func TestChainMatchedReportsTrue(t *testing.T) {
	t.Run("value chain", func(t *testing.T) {
		chain := Value(3).Case(Eq(3), func(int) string { return "three" })
		if !chain.Matched() {
			t.Fatal("value chain: expected Matched=true")
		}
	})
	t.Run("type chain", func(t *testing.T) {
		chain := Type(123).Case(func(n int) string { return "int" })
		if !chain.Matched() {
			t.Fatal("type chain: expected Matched=true")
		}
	})
}

// ---------- final gap-fill: panic paths and matched-or-fully-iterated branches ----------

// TestAnyPatternMatchesEverything covers the inner closure body of
// Any(), which is the `return true` line. The existing tests only
// use Any inside Not, and that wraps Any in a Not which never
// reaches the inner return when the chain is already matched.
func TestAnyPatternMatchesEverything(t *testing.T) {
	for _, v := range []int{0, 1, -1, 100, -100} {
		if !Any[int]().Match(v) {
			t.Fatalf("Any should match %d", v)
		}
	}
}

// TestOrPatternNoMatch covers the `return false` branch of Or: when
// no supplied pattern matches, the function returns false at the
// end of the loop.
func TestOrPatternNoMatch(t *testing.T) {
	p := Or(Eq(1), Eq(2), Eq(3))
	if p.Match(99) {
		t.Fatal("Or(Eq(1), Eq(2), Eq(3)) should not match 99")
	}
	if !p.Match(2) {
		t.Fatal("Or(Eq(1), Eq(2), Eq(3)) should match 2")
	}
}

// TestCaseNilHandler covers the panic in applyValueCase when the
// pattern matches but the handler is nil.
func TestCaseNilHandler(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil Case handler")
		}
	}()
	// Use a typed nil handler so the compiler can infer R = string.
	var h func(int) string
	_ = Value(3).Case(Eq(3), h)
}

// TestChainDefaultNoMatch covers the no-match branch of
// Chain.Default, where then is non-nil and must be invoked with
// the original value. The existing tests do not exercise this
// branch: they check Matched() / Unwrap() but never call Default
// on an unmatched chain.
func TestChainDefaultNoMatch(t *testing.T) {
	got := Value(42).
		Case(Eq(99), func(int) string { return "ninety-nine" }).
		Default(func(n int) string { return fmt.Sprintf("default:%d", n) })
	if got != "default:42" {
		t.Fatalf("got %q, want default:42", got)
	}
}

// TestChainDefaultNilHandler covers the nil-handler panic on
// Chain.Default's no-match branch.
func TestChainDefaultNilHandler(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil Chain.Default handler")
		}
	}()
	var f func(int) string
	_ = Value(42).
		Case(Eq(99), func(int) string { return "ninety-nine" }).
		Default(f)
}

// TestChainOrElseGetNilFallback covers the nil-fallback panic in
// Chain.OrElseGet.
func TestChainOrElseGetNilFallback(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil OrElseGet fallback")
		}
	}()
	var f func(int) string
	_ = Value(42).
		Case(Eq(99), func(int) string { return "ninety-nine" }).
		OrElseGet(f)
}

// TestTypeCaseNilHandler covers the panic in applyTypeCase when the
// type matches but the handler is nil.
func TestTypeCaseNilHandler(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil Type.Case handler")
		}
	}()
	var h func(int) string
	_ = Type(123).Case(h)
}

// TestTypeNilNilHandler covers the panic in applyTypeNil when the
// value is nil but the handler is nil.
func TestTypeNilNilHandler(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil Type.Nil handler")
		}
	}()
	var pointer *int
	var h func() string
	_ = Type(pointer).Nil(h)
}

// TestTypeChainOrElseGetMatched covers the matched-true branch of
// TypeChain.OrElseGet: when a case has already matched, the
// matched result is returned without calling the fallback.
func TestTypeChainOrElseGetMatched(t *testing.T) {
	calls := 0
	got := Type(123).
		Case(func(n int) string { return fmt.Sprintf("%d", n) }).
		OrElseGet(func(any) string {
			calls++
			return "fallback"
		})
	if got != "123" {
		t.Fatalf("got %q, want 123", got)
	}
	if calls != 0 {
		t.Fatalf("OrElseGet calls got %d, want 0 (matched path is lazy)", calls)
	}
}
