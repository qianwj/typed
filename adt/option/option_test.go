package option

import (
	"errors"
	"testing"
)

// TestOptionEmpty exercises the basic absence story for value
// types, which is the case that breaks a nil-sentinel design.
func TestOptionEmpty(t *testing.T) {
	t.Run("value type int", func(t *testing.T) {
		o := Empty[int]()
		if o.IsPresent() {
			t.Fatal("Empty[int] should be absent")
		}
		if !o.IsEmpty() {
			t.Fatal("Empty[int] should report IsEmpty")
		}
	})
	t.Run("value type string", func(t *testing.T) {
		o := Empty[string]()
		if o.IsPresent() || !o.IsEmpty() {
			t.Fatal("Empty[string] should be absent")
		}
	})
	t.Run("value type struct", func(t *testing.T) {
		o := Empty[user]()
		if o.IsPresent() || !o.IsEmpty() {
			t.Fatal("Empty[user] should be absent")
		}
	})
}

// TestOptionOf ensures a present Option of a value type can
// carry a real (non-zero) value and is also able to carry a zero
// value while staying present.
func TestOptionOf(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		o := Of(42)
		if !o.IsPresent() {
			t.Fatal("Of(42) should be present")
		}
		if got := o.Get(); got != 42 {
			t.Fatalf("Get: got %d, want 42", got)
		}
	})
	t.Run("zero value is still present", func(t *testing.T) {
		o := Of(0)
		if !o.IsPresent() {
			t.Fatal("Of(0) should be present — that is the whole point of the present flag")
		}
		if got := o.Get(); got != 0 {
			t.Fatalf("Get: got %d, want 0", got)
		}
	})
	t.Run("empty string is still present", func(t *testing.T) {
		o := Of("")
		if !o.IsPresent() {
			t.Fatal(`Of("") should be present`)
		}
		if got := o.Get(); got != "" {
			t.Fatalf(`Get: got %q, want ""`, got)
		}
	})
}

// TestOptionOfNullable covers the pointer/reference type branch.
// A typed nil must become an absent Option; a non-nil value must
// become a present one.
func TestOptionOfNullable(t *testing.T) {
	t.Run("nil pointer is absent", func(t *testing.T) {
		var p *user
		o := OfNullable(p)
		if o.IsPresent() {
			t.Fatal("OfNullable(nil) should be absent")
		}
	})
	t.Run("non-nil pointer is present", func(t *testing.T) {
		u := &user{Name: "bob", Age: 30}
		o := OfNullable(u)
		if !o.IsPresent() {
			t.Fatal("OfNullable(&u) should be present")
		}
		if got := o.Get(); got != u {
			t.Fatalf("Get: got %p, want %p", got, u)
		}
	})
	t.Run("nil slice is absent", func(t *testing.T) {
		var s []int
		if OfNullable(s).IsPresent() {
			t.Fatal("OfNullable(nil slice) should be absent")
		}
	})
	t.Run("nil map is absent", func(t *testing.T) {
		var m map[string]int
		if OfNullable(m).IsPresent() {
			t.Fatal("OfNullable(nil map) should be absent")
		}
	})
	t.Run("non-nil slice is present", func(t *testing.T) {
		if !OfNullable([]int{1, 2, 3}).IsPresent() {
			t.Fatal("OfNullable([]int{1,2,3}) should be present")
		}
	})
	t.Run("nil channel is absent", func(t *testing.T) {
		var ch chan int
		if OfNullable(ch).IsPresent() {
			t.Fatal("OfNullable(nil chan) should be absent")
		}
	})
	t.Run("nil function is absent", func(t *testing.T) {
		var f func()
		if OfNullable(f).IsPresent() {
			t.Fatal("OfNullable(nil func) should be absent")
		}
	})
	t.Run("nil interface is absent", func(t *testing.T) {
		var i any
		if OfNullable(i).IsPresent() {
			t.Fatal("OfNullable(nil interface) should be absent")
		}
	})
	t.Run("typed nil inside an interface is absent", func(t *testing.T) {
		for _, value := range []any{(*user)(nil), []int(nil), map[string]int(nil), (chan int)(nil), (func())(nil)} {
			if OfNullable(value).IsPresent() {
				t.Errorf("OfNullable(%T(nil)) should be absent", value)
			}
		}
	})
	t.Run("non-nil reference values are present", func(t *testing.T) {
		for _, value := range []any{&user{}, []int{}, map[string]int{}, make(chan int), func() {}} {
			o := OfNullable(value)
			if o.IsEmpty() {
				t.Errorf("OfNullable(%T) should be present", value)
			}
		}
	})
	t.Run("value type never nil", func(t *testing.T) {
		// Value types (int, string, struct, …) cannot be nil.
		// OfNullable must return a present Option regardless
		// of the value, including the zero value. This is the
		// reason OfNullable cannot be used as a universal
		// constructor: prefer Of for value types.
		if !OfNullable(0).IsPresent() {
			t.Fatal("OfNullable(0) should be present")
		}
		if !OfNullable("").IsPresent() {
			t.Fatal(`OfNullable("") should be present`)
		}
		if !OfNullable(user{}).IsPresent() {
			t.Fatal("OfNullable(user{}) should be present")
		}
	})
}

// TestOptionGetPanicOnEmpty documents the deliberate panic when
// the caller extracts from an absent Option without checking.
func TestOptionGetPanicOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Get on empty Option should panic")
		}
	}()
	_ = Empty[int]().Get()
}

// TestOptionOrElse covers both branches of OrElse.
func TestOptionOrElse(t *testing.T) {
	if got := Of(7).OrElse(99); got != 7 {
		t.Fatalf("OrElse on present: got %d, want 7", got)
	}
	if got := Empty[int]().OrElse(99); got != 99 {
		t.Fatalf("OrElse on empty: got %d, want 99", got)
	}
}

// TestOptionOrElseGet ensures the fallback function only runs
// when the receiver is absent. This is the property the eager
// OrElse does not provide.
func TestOptionOrElseGet(t *testing.T) {
	t.Run("present skips the fallback", func(t *testing.T) {
		called := false
		got := Of(5).OrElseGet(func() int {
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
		got := Empty[int]().OrElseGet(func() int { return 42 })
		if got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})
}

// TestOptionOrElseThrow documents the (T, error) shape. Note that
// the function does not panic despite the name: the Go convention
// is to return an error, and the name is borrowed from Java's
// Option.orElseThrow.
func TestOptionOrElseThrow(t *testing.T) {
	t.Run("present returns value with nil error", func(t *testing.T) {
		v, err := Of(7).OrElseThrow("missing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 7 {
			t.Fatalf("got %d, want 7", v)
		}
	})
	t.Run("empty returns the message", func(t *testing.T) {
		_, err := Empty[int]().OrElseThrow("missing")
		if err == nil || err.Error() != "missing" {
			t.Fatalf("got %v, want error 'missing'", err)
		}
	})
}

// TestOptionIfPresent confirms the side-effecting callback is
// only invoked when the Option is present.
func TestOptionIfPresent(t *testing.T) {
	called := 0
	Of(10).IfPresent(func(int) { called++ })
	Empty[int]().IfPresent(func(int) { called++ })
	if called != 1 {
		t.Fatalf("IfPresent: callback ran %d times, want 1", called)
	}
}

// TestOptionIfPresentOrElse checks that exactly one of the two
// callbacks runs, depending on presence.
func TestOptionIfPresentOrElse(t *testing.T) {
	present, absent := 0, 0
	track := func(int) { present++ }
	trackAbsent := func() { absent++ }

	Of(1).IfPresentOrElse(track, trackAbsent)
	Empty[int]().IfPresentOrElse(track, trackAbsent)

	if present != 1 {
		t.Fatalf("present callback ran %d times, want 1", present)
	}
	if absent != 1 {
		t.Fatalf("absent callback ran %d times, want 1", absent)
	}
}

// TestOptionFilter covers all three branches: present+keep,
// present+drop, and absent.
func TestOptionFilter(t *testing.T) {
	t.Run("present and kept", func(t *testing.T) {
		got := Of(10).Filter(func(n int) bool { return n > 5 })
		if !got.IsPresent() || got.Get() != 10 {
			t.Fatalf("got %v, want present 10", got)
		}
	})
	t.Run("present and dropped", func(t *testing.T) {
		got := Of(3).Filter(func(n int) bool { return n > 5 })
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("absent stays absent", func(t *testing.T) {
		called := false
		got := Empty[int]().Filter(func(int) bool {
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

// TestOptionMap exercises a type-changing chain and verifies the
// absent branch is propagated without invoking f.
func TestOptionMap(t *testing.T) {
	t.Run("present applies f and changes type", func(t *testing.T) {
		got := Of(7).Map(func(n int) string {
			return "n=" + itoa(n)
		})
		if !got.IsPresent() || got.Get() != "n=7" {
			t.Fatalf("got %v, want present n=7", got)
		}
	})
	t.Run("empty propagates without calling f", func(t *testing.T) {
		called := false
		got := Empty[int]().Map(func(int) string {
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

// TestOptionFlatMap checks that FlatMap delegates to the inner
// Option returned by f when present, and propagates absence
// without calling f.
func TestOptionFlatMap(t *testing.T) {
	t.Run("present delegates to f", func(t *testing.T) {
		got := Of(5).FlatMap(func(n int) Option[string] {
			if n > 0 {
				return Of(itoa(n))
			}
			return Empty[string]()
		})
		if !got.IsPresent() || got.Get() != "5" {
			t.Fatalf("got %v, want present 5", got)
		}
	})
	t.Run("f returning empty stays empty", func(t *testing.T) {
		got := Of(-1).FlatMap(func(int) Option[string] {
			return Empty[string]()
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("empty propagates without calling f", func(t *testing.T) {
		called := false
		got := Empty[int]().FlatMap(func(int) Option[string] {
			called = true
			return Of("x")
		})
		if got.IsPresent() {
			t.Fatalf("got %v, want empty", got)
		}
		if called {
			t.Fatal("f should not run on empty receiver")
		}
	})
}

// TestOptionChained combinators read more naturally than a manual
// sequence of conditionals. This is a smoke test for the
// documentation example.
func TestOptionChained(t *testing.T) {
	parse := func(s string) Option[int] {
		if s == "" {
			return Empty[int]()
		}
		v, err := atoi(s)
		if err != nil {
			return Empty[int]()
		}
		return Of(v)
	}

	got := parse("123").
		Filter(func(n int) bool { return n > 0 }).
		Map(func(n int) int { return n * 2 })

	if !got.IsPresent() || got.Get() != 246 {
		t.Fatalf("chained: got %v, want present 246", got)
	}
}

// ---------- Wrap ----------

// TestOptionWrapPresentTrue confirms that Wrap(value, true)
// produces a present Option carrying value. The presence flag
// is the sole determinant of the result; value is stored verbatim.
func TestOptionWrapPresentTrue(t *testing.T) {
	o := Wrap(42, true)
	if !o.IsPresent() {
		t.Fatal("Wrap(42, true): IsPresent = false, want true")
	}
	if o.IsEmpty() {
		t.Fatal("Wrap(42, true): IsEmpty = true, want false")
	}
	if got := o.Get(); got != 42 {
		t.Fatalf("Wrap(42, true).Get: got %d, want 42", got)
	}
}

// TestOptionWrapPresentFalse confirms that Wrap(value, false)
// produces an absent Option. The value parameter is stored in
// the underlying struct but is not reachable through any Option
// method — Get must panic, IsPresent / IsEmpty must report absent.
func TestOptionWrapPresentFalse(t *testing.T) {
	o := Wrap(99, false)
	if o.IsPresent() {
		t.Fatal("Wrap(99, false): IsPresent = true, want false")
	}
	if !o.IsEmpty() {
		t.Fatal("Wrap(99, false): IsEmpty = false, want true")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Wrap(99, false).Get: should panic")
		}
	}()
	_ = o.Get()
}

// TestOptionWrapZeroValueStillPresent confirms that Wrap(0, true)
// produces a present Option carrying the zero value, not an absent
// one. This is the case where Wrap differs from "absent iff value
// is the zero value" — the present flag is the only signal that
// matters.
func TestOptionWrapZeroValueStillPresent(t *testing.T) {
	o := Wrap(0, true)
	if !o.IsPresent() {
		t.Fatal("Wrap(0, true): IsPresent = false, want true (zero value with present=true is present)")
	}
	if got := o.Get(); got != 0 {
		t.Fatalf("Wrap(0, true).Get: got %d, want 0", got)
	}

	emptyString := Wrap("", true)
	if !emptyString.IsPresent() {
		t.Fatal(`Wrap("", true): IsPresent = false, want true`)
	}
	if got := emptyString.Get(); got != "" {
		t.Fatalf(`Wrap("", true).Get: got %q, want ""`, got)
	}
}

// TestOptionWrapObservationalIgnoresValueOnAbsent confirms that
// the value parameter is unreachable through the Option API when
// present is false. This is the property that lets callers write
// `return Wrap(val, ok)` after a (T, bool) lookup without first
// checking which branch they are in.
//
// Wrap(42, false) and Wrap(0, false) must agree on every public
// observation: same IsPresent/IsEmpty, same panic on Get, same
// fallback from OrElse/OrElseGet, same absent from Filter/Map.
func TestOptionWrapObservationalIgnoresValueOnAbsent(t *testing.T) {
	withValue := Wrap(42, false)
	withZero := Wrap(0, false)

	if withValue.IsPresent() || withZero.IsPresent() {
		t.Fatal("absent Wrap: IsPresent should be false regardless of value")
	}
	if !withValue.IsEmpty() || !withZero.IsEmpty() {
		t.Fatal("absent Wrap: IsEmpty should be true regardless of value")
	}
	if withValue.OrElse(7) != withZero.OrElse(7) {
		t.Fatal("OrElse of absent Wrap must agree, regardless of value")
	}
	if withValue.OrElseGet(func() int { return 7 }) != withZero.OrElseGet(func() int { return 7 }) {
		t.Fatal("OrElseGet of absent Wrap must agree, regardless of value")
	}
	if withValue.Filter(func(int) bool { return true }).IsPresent() {
		t.Fatal("Filter on absent Wrap: should remain absent")
	}
	if withValue.Map(func(n int) int { return n + 1 }).IsPresent() {
		t.Fatal("Map on absent Wrap: should remain absent")
	}
}

// TestOptionWrapMatchesEmpty confirms that Wrap(zero, false) is
// exactly equivalent to Empty[T]() through the public methods.
// This also serves as a regression test for the internal storage
// invariant: the value field of an absent Wrap is not zeroed by
// the constructor, so the public API must not leak it.
func TestOptionWrapMatchesEmpty(t *testing.T) {
	cases := []struct {
		name    string
		wrapped Option[int]
		plain   Option[int]
	}{
		{"zero value, false", Wrap(0, false), Empty[int]()},
		{"non-zero value, false", Wrap(99, false), Empty[int]()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.wrapped.IsPresent() != c.plain.IsPresent() {
				t.Fatalf("IsPresent differs: %v vs %v", c.wrapped.IsPresent(), c.plain.IsPresent())
			}
			if c.wrapped.IsEmpty() != c.plain.IsEmpty() {
				t.Fatalf("IsEmpty differs")
			}
			if c.wrapped.OrElse(7) != c.plain.OrElse(7) {
				t.Fatalf("OrElse differs")
			}
		})
	}
}

// TestOptionWrapComposesWithMap confirms that a Wrap-built
// Option participates in the standard Map chain. Map on a
// present Wrap transforms the value; Map on an absent Wrap
// propagates the absence without calling f.
func TestOptionWrapComposesWithMap(t *testing.T) {
	transformed := Wrap(21, true).Map(func(n int) int { return n * 2 })
	if !transformed.IsPresent() || transformed.Get() != 42 {
		t.Fatalf("Map on present Wrap: got %v, want present 42", transformed)
	}

	called := false
	absent := Wrap(999, false).Map(func(n int) int {
		called = true
		return n * 2
	})
	if absent.IsPresent() {
		t.Fatal("Map on absent Wrap: should remain absent")
	}
	if called {
		t.Fatal("Map on absent Wrap: f should not run")
	}
}

// TestOptionWrapComposesWithOrElse confirms that Wrap-built
// Options participate in OrElse and OrElseGet. The present
// path returns the wrapped value; the absent path returns the
// fallback without touching the wrapped value.
func TestOptionWrapComposesWithOrElse(t *testing.T) {
	if v := Wrap(42, true).OrElse(0); v != 42 {
		t.Fatalf("OrElse on present Wrap: got %d, want 42", v)
	}
	if v := Wrap(0, false).OrElse(99); v != 99 {
		t.Fatalf("OrElse on absent Wrap: got %d, want 99 (fallback)", v)
	}
	if v := Wrap(0, false).OrElseGet(func() int { return 100 }); v != 100 {
		t.Fatalf("OrElseGet on absent Wrap: got %d, want 100", v)
	}
}

// TestOptionWrapPracticalAdapter documents the realistic shape
// that Wrap is for: a function that returns (T, bool) and is
// being plugged into an Option chain. The test confirms Wrap
// does what the doc comment claims — replaces the if-ok ladder
// with a single return.
func TestOptionWrapPracticalAdapter(t *testing.T) {
	// Source function: returns (int, bool), the standard
	// (T, bool) shape. We will call it twice and route each
	// result through an Option chain via Wrap.
	double := func(n int) (int, bool) {
		if n < 0 {
			return 0, false
		}
		return n * 2, true
	}

	// Wrap the hit into a chain.
	hit := Wrap(double(21)) // (42, true)
	if v := hit.Map(func(n int) int { return n + 1 }).Get(); v != 43 {
		t.Fatalf("Wrap(hit).Map: got %d, want 43", v)
	}

	// Wrap the miss into a chain that falls back.
	miss := Wrap(double(-1)) // (0, false)
	if v := miss.OrElse(99); v != 99 {
		t.Fatalf("Wrap(miss).OrElse: got %d, want 99 (fallback)", v)
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
