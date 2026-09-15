package adt

import (
	"errors"
	"testing"
)

var (
	errSample = errors.New("sample error")
)

// leftOrRight returns (value, present) — same shape as (T, bool) but
// via Option, which is what Either's safe accessors return.
func leftOrRight[L, R any](e Either[L, R]) (L, bool) {
	opt := e.Left()
	if !opt.IsPresent() {
		var zero L
		return zero, false
	}
	return opt.Get(), true
}

func rightOrPresent[L, R any](e Either[L, R]) (R, bool) {
	opt := e.Right()
	if !opt.IsPresent() {
		var zero R
		return zero, false
	}
	return opt.Get(), true
}

func TestLeftAndRight(t *testing.T) {
	t.Parallel()

	l := Left[string, int]("left-val")
	if !l.IsLeft() || l.IsRight() {
		t.Fatalf("Left: IsLeft=%v IsRight=%v, want true/false", l.IsLeft(), l.IsRight())
	}

	r := Right[string, int](42)
	if !r.IsRight() || r.IsLeft() {
		t.Fatalf("Right: IsLeft=%v IsRight=%v, want false/true", r.IsLeft(), r.IsRight())
	}
}

func TestAccessors(t *testing.T) {
	t.Parallel()

	l := Left[string, int]("L")
	if got, ok := leftOrRight(l); !ok || got != "L" {
		t.Errorf("Left().Get() = (%q, %v), want (\"L\", true)", got, ok)
	}
	if _, ok := rightOrPresent[string](l); ok {
		t.Errorf("Left.Right() should be absent")
	}

	r := Right[string, int](7)
	if got, ok := rightOrPresent(r); !ok || got != 7 {
		t.Errorf("Right().Get() = (%d, %v), want (7, true)", got, ok)
	}
	if _, ok := leftOrRight[string](r); ok {
		t.Errorf("Right.Left() should be absent")
	}
}

func TestOrZero(t *testing.T) {
	t.Parallel()

	l := Left[string, int]("hello")
	if got := l.LeftOrZero(); got != "hello" {
		t.Errorf("Left.LeftOrZero() = %q, want \"hello\"", got)
	}
	if got := l.RightOrZero(); got != 0 {
		t.Errorf("Left.RightOrZero() = %d, want 0", got)
	}

	r := Right[string, int](99)
	if got := r.RightOrZero(); got != 99 {
		t.Errorf("Right.RightOrZero() = %d, want 99", got)
	}
	if got := r.LeftOrZero(); got != "" {
		t.Errorf("Right.LeftOrZero() = %q, want \"\"", got)
	}
}

func TestMapLeft(t *testing.T) {
	t.Parallel()

	l := Left[error, int](errSample)
	mapped := l.MapLeft(func(e error) string { return e.Error() })
	if !mapped.IsLeft() {
		t.Fatalf("MapLeft: IsLeft=%v, want true", mapped.IsLeft())
	}
	if got, ok := leftOrRight(mapped); !ok || got != "sample error" {
		t.Errorf("mapped.Left() = (%q, %v), want (\"sample error\", true)", got, ok)
	}

	r := Right[error, int](42)
	mappedR := r.MapLeft(func(e error) string { return "should not run" })
	if !mappedR.IsRight() {
		t.Fatalf("MapLeft on Right: IsRight=%v, want true", mappedR.IsRight())
	}
	if got, ok := rightOrPresent(mappedR); !ok || got != 42 {
		t.Errorf("mappedR.Right() = (%d, %v), want (42, true)", got, ok)
	}
}

func TestMapRight(t *testing.T) {
	t.Parallel()

	r := Right[error, int](21)
	mapped := r.MapRight(func(n int) int { return n * 2 })
	if !mapped.IsRight() {
		t.Fatalf("MapRight: IsRight=%v, want true", mapped.IsRight())
	}
	if got, ok := rightOrPresent(mapped); !ok || got != 42 {
		t.Errorf("mapped.Right() = (%d, %v), want (42, true)", got, ok)
	}

	l := Left[error, int](errSample)
	mappedL := l.MapRight(func(n int) int { return -1 })
	if !mappedL.IsLeft() {
		t.Fatalf("MapRight on Left: IsLeft=%v, want true", mappedL.IsLeft())
	}
	if got, ok := leftOrRight(mappedL); !ok || got != errSample {
		t.Errorf("mappedL.Left() = (%v, %v), want (errSample, true)", got, ok)
	}
}

func TestFold(t *testing.T) {
	t.Parallel()

	l := Left[string, int]("left")
	got := l.Fold(
		func(s string) int { return -len(s) },
		func(n int) int { return n },
	)
	if got != -4 {
		t.Errorf("Fold Left = %d, want -4", got)
	}

	r := Right[string, int](7)
	got = r.Fold(
		func(s string) int { return -len(s) },
		func(n int) int { return n },
	)
	if got != 7 {
		t.Errorf("Fold Right = %d, want 7", got)
	}
}

func TestFoldOnlyOneCallbackRuns(t *testing.T) {
	t.Parallel()

	var leftCalls, rightCalls int
	r := Right[error, int](1)
	r.Fold(
		func(e error) int { leftCalls++; return 0 },
		func(n int) int { rightCalls++; return n },
	)
	if leftCalls != 0 || rightCalls != 1 {
		t.Errorf("Fold Right: leftCalls=%d rightCalls=%d, want 0/1", leftCalls, rightCalls)
	}

	leftCalls, rightCalls = 0, 0
	l := Left[error, int](errSample)
	l.Fold(
		func(e error) int { leftCalls++; return 0 },
		func(n int) int { rightCalls++; return n },
	)
	if leftCalls != 1 || rightCalls != 0 {
		t.Errorf("Fold Left: leftCalls=%d rightCalls=%d, want 1/0", leftCalls, rightCalls)
	}
}

func TestRoundTripThroughOption(t *testing.T) {
	t.Parallel()
	// Verify the Option combinators compose with Either's safe
	// accessors.
	r := Right[error, int](123)
	opt := r.Right()
	if !opt.IsPresent() {
		t.Fatal("Right.Right() not present")
	}
	if opt.Map(func(n int) int { return n + 1 }).Get() != 124 {
		t.Errorf("Option Map on Right.Right() failed")
	}

	l := Left[error, int](errSample)
	if l.Right().IsPresent() {
		t.Error("Left.Right() should be absent")
	}
	if !l.Left().IsPresent() {
		t.Error("Left.Left() should be present")
	}
}