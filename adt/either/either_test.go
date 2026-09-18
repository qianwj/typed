package either

import (
	"errors"
	"strconv"
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

func TestEitherMapRightTypeChange(t *testing.T) {
	t.Parallel()
	calls := 0
	stringify := func(n int) string {
		calls++
		return strconv.Itoa(n)
	}
	got := Right[error, int](42).MapRight(stringify)
	if !got.IsRight() || got.Right().Get() != "42" || calls != 1 {
		t.Fatalf("MapRight = %v, calls = %d; want Right(42), 1 call", got, calls)
	}
	left := Left[error, int](errSample).MapRight(stringify)
	if !left.IsLeft() || left.Left().Get() != errSample || calls != 1 {
		t.Fatalf("MapRight on Left = %v, calls = %d; want original Left, no new calls", left, calls)
	}
}

func TestEitherFlatMapRight(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		input Either[string, int]
		want  Either[string, string]
		calls int
	}{
		{"right to right", Right[string, int](42), Right[string, string]("42"), 1},
		{"right to left", Right[string, int](-1), Left[string, string]("negative"), 1},
		{"left skips callback", Left[string, int]("original"), Left[string, string]("original"), 0},
		{"zero value is left", Either[string, int]{}, Left[string, string](""), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			got := tc.input.FlatMapRight(func(n int) Either[string, string] {
				calls++
				if n < 0 {
					return Left[string, string]("negative")
				}
				return Right[string, string](strconv.Itoa(n))
			})
			if got != tc.want || calls != tc.calls {
				t.Fatalf("FlatMapRight = %v, calls = %d; want %v, %d", got, calls, tc.want, tc.calls)
			}
		})
	}
}

func TestEitherFlatMapLeft(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		input Either[int, string]
		want  Either[string, string]
		calls int
	}{
		{"left to left", Left[int, string](42), Left[string, string]("42"), 1},
		{"left to right", Left[int, string](-1), Right[string, string]("recovered"), 1},
		{"right skips callback", Right[int, string]("original"), Right[string, string]("original"), 0},
		{"zero value invokes callback", Either[int, string]{}, Left[string, string]("0"), 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			got := tc.input.FlatMapLeft(func(n int) Either[string, string] {
				calls++
				if n < 0 {
					return Right[string, string]("recovered")
				}
				return Left[string, string](strconv.Itoa(n))
			})
			if got != tc.want || calls != tc.calls {
				t.Fatalf("FlatMapLeft = %v, calls = %d; want %v, %d", got, calls, tc.want, tc.calls)
			}
		})
	}
}

func TestEitherFlatMapLeftRecoveryChain(t *testing.T) {
	t.Parallel()
	got := Left[string, int]("missing").
		FlatMapLeft(func(problem string) Either[error, int] {
			if problem != "missing" {
				t.Fatalf("callback received %q, want missing", problem)
			}
			return Right[error, int](21)
		}).
		FlatMapLeft(func(error) Either[string, int] {
			t.Fatal("FlatMapLeft callback ran after recovery to Right")
			return Left[string, int]("unexpected")
		}).
		FlatMapRight(func(n int) Either[string, int] {
			return Right[string, int](n * 2)
		})
	if !got.IsRight() || got.Right().Get() != 42 {
		t.Fatalf("chain = %v, want Right(42)", got)
	}
}

func TestEitherFlatMapRightChainShortCircuits(t *testing.T) {
	t.Parallel()
	got := Right[error, int](42).
		FlatMapRight(func(int) Either[error, string] {
			return Left[error, string](errSample)
		}).
		FlatMapRight(func(string) Either[error, bool] {
			t.Fatal("FlatMapRight callback ran after Left")
			return Right[error, bool](true)
		}).
		MapRight(func(bool) int {
			t.Fatal("MapRight callback ran after Left")
			return 1
		})
	if !got.IsLeft() || got.Left().Get() != errSample {
		t.Fatalf("chain = %v, want original Left", got)
	}
}

func TestEitherFlatMapRightPreservesNilBranches(t *testing.T) {
	t.Parallel()
	left := Left[*int, int](nil).FlatMapRight(func(int) Either[*int, string] {
		t.Fatal("FlatMapRight callback ran on Left(nil)")
		return Right[*int, string]("")
	})
	if !left.IsLeft() || !left.Left().IsPresent() || left.Left().Get() != nil {
		t.Fatalf("got %v, want Left(nil)", left)
	}
	right := Right[string, *int](nil).FlatMapRight(func(p *int) Either[string, *int] {
		if p != nil {
			t.Fatal("callback did not receive nil")
		}
		return Right[string, *int](p)
	})
	if !right.IsRight() || !right.Right().IsPresent() || right.Right().Get() != nil {
		t.Fatalf("got %v, want Right(nil)", right)
	}
}

func TestEitherFlatMapLeftPreservesNilBranches(t *testing.T) {
	t.Parallel()
	right := Right[int, *int](nil).FlatMapLeft(func(int) Either[string, *int] {
		t.Fatal("FlatMapLeft callback ran on Right(nil)")
		return Left[string, *int]("")
	})
	if !right.IsRight() || !right.Right().IsPresent() || right.Right().Get() != nil {
		t.Fatalf("got %v, want Right(nil)", right)
	}
	left := Left[*int, string](nil).FlatMapLeft(func(p *int) Either[*int, string] {
		if p != nil {
			t.Fatal("callback did not receive nil")
		}
		return Left[*int, string](p)
	})
	if !left.IsLeft() || !left.Left().IsPresent() || left.Left().Get() != nil {
		t.Fatalf("got %v, want Left(nil)", left)
	}
}