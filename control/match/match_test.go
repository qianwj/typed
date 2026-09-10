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
