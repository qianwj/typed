package utils_test

import (
	"errors"
	"testing"

	"github.com/qianwj/typed/utils"
)

// TestOptionalEmpty exercises the basic absence story for value
// types, which is the case that breaks a nil-sentinel design.
func TestOptionalEmpty(t *testing.T) {
	t.Run("value type int", func(t *testing.T) {
		o := utils.Empty[int]()
		if o.IsPresent() {
			t.Fatal("Empty[int] should be absent")
		}
		if !o.IsEmpty() {
			t.Fatal("Empty[int] should report IsEmpty")
		}
	})
	t.Run("value type string", func(t *testing.T) {
		o := utils.Empty[string]()
		if o.IsPresent() || !o.IsEmpty() {
			t.Fatal("Empty[string] should be absent")
		}
	})
	t.Run("value type struct", func(t *testing.T) {
		o := utils.Empty[user]()
		if o.IsPresent() || !o.IsEmpty() {
			t.Fatal("Empty[user] should be absent")
		}
	})
}

// TestOptionalOf ensures a present Optional of a value type can
// carry a real (non-zero) value and is also able to carry a zero
// value while staying present.
func TestOptionalOf(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		o := utils.Of(42)
		if !o.IsPresent() {
			t.Fatal("Of(42) should be present")
		}
		if got := o.Get(); got != 42 {
			t.Fatalf("Get: got %d, want 42", got)
		}
	})
	t.Run("zero value is still present", func(t *testing.T) {
		o := utils.Of(0)
		if !o.IsPresent() {
			t.Fatal("Of(0) should be present — that is the whole point of the present flag")
		}
		if got := o.Get(); got != 0 {
			t.Fatalf("Get: got %d, want 0", got)
		}
	})
	t.Run("empty string is still present", func(t *testing.T) {
		o := utils.Of("")
		if !o.IsPresent() {
			t.Fatal(`Of("") should be present`)
		}
		if got := o.Get(); got != "" {
			t.Fatalf(`Get: got %q, want ""`, got)
		}
	})
}

// TestOptionalOfNullable covers the pointer/reference type branch.
// A typed nil must become an absent Optional; a non-nil value must
// become a present one.
func TestOptionalOfNullable(t *testing.T) {
	t.Run("nil pointer is absent", func(t *testing.T) {
		var p *user
		o := utils.OfNullable(p)
		if o.IsPresent() {
			t.Fatal("OfNullable(nil) should be absent")
		}
	})
	t.Run("non-nil pointer is present", func(t *testing.T) {
		u := &user{Name: "bob", Age: 30}
		o := utils.OfNullable(u)
		if !o.IsPresent() {
			t.Fatal("OfNullable(&u) should be present")
		}
		if got := o.Get(); got != u {
			t.Fatalf("Get: got %p, want %p", got, u)
		}
	})
	t.Run("nil slice is absent", func(t *testing.T) {
		var s []int
		if utils.OfNullable(s).IsPresent() {
			t.Fatal("OfNullable(nil slice) should be absent")
		}
	})
	t.Run("nil map is absent", func(t *testing.T) {
		var m map[string]int
		if utils.OfNullable(m).IsPresent() {
			t.Fatal("OfNullable(nil map) should be absent")
		}
	})
	t.Run("non-nil slice is present", func(t *testing.T) {
		if !utils.OfNullable([]int{1, 2, 3}).IsPresent() {
			t.Fatal("OfNullable([]int{1,2,3}) should be present")
		}
	})
	t.Run("nil channel is absent", func(t *testing.T) {
		var ch chan int
		if utils.OfNullable(ch).IsPresent() {
			t.Fatal("OfNullable(nil chan) should be absent")
		}
	})
	t.Run("nil function is absent", func(t *testing.T) {
		var f func()
		if utils.OfNullable(f).IsPresent() {
			t.Fatal("OfNullable(nil func) should be absent")
		}
	})
	t.Run("nil interface is absent", func(t *testing.T) {
		var i any
		if utils.OfNullable(i).IsPresent() {
			t.Fatal("OfNullable(nil interface) should be absent")
		}
	})
	t.Run("value type never nil", func(t *testing.T) {
		// Value types (int, string, struct, …) cannot be nil.
		// OfNullable must return a present Optional regardless
		// of the value, including the zero value. This is the
		// reason OfNullable cannot be used as a universal
		// constructor: prefer Of for value types.
		if !utils.OfNullable(0).IsPresent() {
			t.Fatal("OfNullable(0) should be present")
		}
		if !utils.OfNullable("").IsPresent() {
			t.Fatal(`OfNullable("") should be present`)
		}
		if !utils.OfNullable(user{}).IsPresent() {
			t.Fatal("OfNullable(user{}) should be present")
		}
	})
}

// TestOptionalGetPanicOnEmpty documents the deliberate panic when
// the caller extracts from an absent Optional without checking.
func TestOptionalGetPanicOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Get on empty Optional should panic")
		}
	}()
	_ = utils.Empty[int]().Get()
}

// TestOptionalOrElse covers both branches of OrElse.
func TestOptionalOrElse(t *testing.T) {
	if got := utils.Of(7).OrElse(99); got != 7 {
		t.Fatalf("OrElse on present: got %d, want 7", got)
	}
	if got := utils.Empty[int]().OrElse(99); got != 99 {
		t.Fatalf("OrElse on empty: got %d, want 99", got)
	}
}

// TestOptionalOrElseGet ensures the fallback function only runs
// when the receiver is absent. This is the property the eager
// OrElse does not provide.
func TestOptionalOrElseGet(t *testing.T) {
	t.Run("present skips the fallback", func(t *testing.T) {
		called := false
		got := utils.Of(5).OrElseGet(func() int {
			called = true
			return 99
		})
		if got != 5 {
			t.Fatalf("got %d, want 5", got)
		}
		if called {
			t.Fatal("fallback should not be called when present")
		}
	})
	t.Run("empty calls the fallback", func(t *testing.T) {
		got := utils.Empty[int]().OrElseGet(func() int { return 42 })
		if got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})
}

// TestOptionalOrElseThrow documents the (T, error) shape. Note that
// the function does not panic despite the name: the Go convention
// is to return an error, and the name is borrowed from Java's
// Optional.orElseThrow.
func TestOptionalOrElseThrow(t *testing.T) {
	t.Run("present returns value with nil error", func(t *testing.T) {
		v, err := utils.Of(7).OrElseThrow("missing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 7 {
			t.Fatalf("got %d, want 7", v)
		}
	})
	t.Run("empty returns the message", func(t *testing.T) {
		_, err := utils.Empty[int]().OrElseThrow("missing")
		if err == nil || err.Error() != "missing" {
			t.Fatalf("got %v, want error 'missing'", err)
		}
	})
}

// TestOptionalIfPresent confirms the side-effecting callback is
// only invoked when the Optional is present.
func TestOptionalIfPresent(t *testing.T) {
	called := 0
	utils.Of(10).IfPresent(func(int) { called++ })
	utils.Empty[int]().IfPresent(func(int) { called++ })
	if called != 1 {
		t.Fatalf("IfPresent: callback ran %d times, want 1", called)
	}
}

// TestOptionalIfPresentOrElse checks that exactly one of the two
// callbacks runs, depending on presence.
func TestOptionalIfPresentOrElse(t *testing.T) {
	present, absent := 0, 0
	track := func(int) { present++ }
	trackAbsent := func() { absent++ }

	utils.Of(1).IfPresentOrElse(track, trackAbsent)
	utils.Empty[int]().IfPresentOrElse(track, trackAbsent)

	if present != 1 {
		t.Fatalf("present callback ran %d times, want 1", present)
	}
	if absent != 1 {
		t.Fatalf("absent callback ran %d times, want 1", absent)
	}
}

// TestOptionalFilter covers all three branches: present+keep,
// present+drop, and absent.
func TestOptionalFilter(t *testing.T) {
	t.Run("present and kept", func(t *testing.T) {
		got := utils.Of(10).Filter(func(n int) bool { return n > 5 })
		if !got.IsPresent() || got.Get() != 10 {
			t.Fatalf("got %v, want present 10", got)
		}
	})
	t.Run("present and dropped", func(t *testing.T) {
		got := utils.Of(3).Filter(func(n int) bool { return n > 5 })
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("absent stays absent", func(t *testing.T) {
		called := false
		got := utils.Empty[int]().Filter(func(int) bool {
			called = true
			return true
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
		if called {
			t.Fatal("predicate should not run on empty receiver")
		}
	})
}

// TestOptionalMap exercises a type-changing chain and verifies the
// absent branch is propagated without invoking f.
func TestOptionalMap(t *testing.T) {
	t.Run("present applies f and changes type", func(t *testing.T) {
		got := utils.Of(7).Map(func(n int) string {
			return "n=" + itoa(n)
		})
		if !got.IsPresent() || got.Get() != "n=7" {
			t.Fatalf("got %v, want present n=7", got)
		}
	})
	t.Run("empty propagates without calling f", func(t *testing.T) {
		called := false
		got := utils.Empty[int]().Map(func(int) string {
			called = true
			return "x"
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
		if called {
			t.Fatal("f should not run on empty receiver")
		}
	})
}

// TestOptionalFlatMap checks that FlatMap delegates to the inner
// Optional returned by f when present, and propagates absence
// without calling f.
func TestOptionalFlatMap(t *testing.T) {
	t.Run("present delegates to f", func(t *testing.T) {
		got := utils.Of(5).FlatMap(func(n int) utils.Optional[string] {
			if n > 0 {
				return utils.Of(itoa(n))
			}
			return utils.Empty[string]()
		})
		if !got.IsPresent() || got.Get() != "5" {
			t.Fatalf("got %v, want present 5", got)
		}
	})
	t.Run("f returning empty stays empty", func(t *testing.T) {
		got := utils.Of(-1).FlatMap(func(int) utils.Optional[string] {
			return utils.Empty[string]()
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("empty propagates without calling f", func(t *testing.T) {
		called := false
		got := utils.Empty[int]().FlatMap(func(int) utils.Optional[string] {
			called = true
			return utils.Of("x")
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
		if called {
			t.Fatal("f should not run on empty receiver")
		}
	})
}

// TestOptionalChained combinators read more naturally than a manual
// sequence of conditionals. This is a smoke test for the
// documentation example.
func TestOptionalChained(t *testing.T) {
	parse := func(s string) utils.Optional[int] {
		if s == "" {
			return utils.Empty[int]()
		}
		v, err := atoi(s)
		if err != nil {
			return utils.Empty[int]()
		}
		return utils.Of(v)
	}

	got := parse("123").
		Filter(func(n int) bool { return n > 0 }).
		Map(func(n int) int { return n * 2 })

	if !got.IsPresent() || got.Get() != 246 {
		t.Fatalf("chained: got %v, want present 246", got)
	}
}

// ---------- helpers ----------

type user struct {
	Name string
	Age  int
}

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
