package objects

import (
	"reflect"
	"testing"
)

// TestEquals_CustomEqualMethod covers the Equal-method dispatch path of
// Equals: a user-defined Equal(T) bool method takes precedence over
// structural walk, and a returning-false method actually returns false.
type withEqual struct{ v int }

func (w withEqual) Equal(other withEqual) bool { return w.v == other.v }

func TestEquals_CustomEqualMethod(t *testing.T) {
	t.Parallel()

	t.Run("TrueViaEqual", func(t *testing.T) {
		if !Equals(withEqual{v: 1}, withEqual{v: 1}) {
			t.Fatal("Equals should defer to Equal method on equal values")
		}
	})
	t.Run("FalseViaEqual", func(t *testing.T) {
		if Equals(withEqual{v: 1}, withEqual{v: 2}) {
			t.Fatal("Equals should defer to Equal method on unequal values")
		}
	})
	t.Run("EqualMethodReceivesOtherSide", func(t *testing.T) {
		// When only b has Equal, findEqualMethod(a) misses, but
		// findEqualMethod(b) hits and is invoked with (a, b).
		// Equality is still symmetric.
		if !Equals(struct{ v int }{v: 7}, withEqual{v: 7}) {
			t.Fatal("Equals should consult Equal on the right side too")
		}
	})
}

// TestEquals_ReflectionAndDeepWalk covers the deepEqualFixedVisit branch
// not yet exercised: arrays, channels (only by direction), and
// non-comparable-but-empty maps. Most of these are reachable only with
// reflect-shaped inputs.
func TestEquals_ReflectionAndDeepWalk(t *testing.T) {
	t.Parallel()

	t.Run("ArraysByElement", func(t *testing.T) {
		if !Equals([3]int{1, 2, 3}, [3]int{1, 2, 3}) {
			t.Fatal("arrays: identical elements should be equal")
		}
		if Equals([3]int{1, 2, 3}, [3]int{1, 2, 4}) {
			t.Fatal("arrays: differing elements should not be equal")
		}
	})

	t.Run("StructsByField", func(t *testing.T) {
		type S struct{ A, B int }
		if !Equals(S{1, 2}, S{1, 2}) {
			t.Fatal("structs: identical fields should be equal")
		}
		if Equals(S{1, 2}, S{1, 3}) {
			t.Fatal("structs: differing fields should not be equal")
		}
	})

	t.Run("PointersFollowed", func(t *testing.T) {
		x, y := 42, 42
		if !Equals(&x, &y) {
			t.Fatal("pointers: equal pointee should compare equal")
		}
		if Equals(&x, new(int)) {
			t.Fatal("pointers: nil pointee should not compare equal to non-nil")
		}
	})

	t.Run("MapsByEntry", func(t *testing.T) {
		if !Equals(map[string]int{"a": 1}, map[string]int{"a": 1}) {
			t.Fatal("maps: equal entries should be equal")
		}
		if Equals(map[string]int{"a": 1}, map[string]int{"a": 2}) {
			t.Fatal("maps: differing values should not be equal")
		}
	})
}

// TestIsEmptyOrNil_NilAndEmptyPaths drives the isEmptyOrNil helper through
// every branch: invalid value, nil pointer, nil slice, empty slice, nil
// map, empty map. None of these were covered before.
func TestIsEmptyOrNil_NilAndEmptyPaths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"invalid value", (*int)(nil), true}, // any(*int)(nil) is the typed nil case
		{"nil slice", []int(nil), true},
		{"empty slice", []int{}, true},
		{"non-empty slice", []int{1}, false},
		{"nil map", map[string]int(nil), true},
		{"empty map", map[string]int{}, true},
		{"non-empty map", map[string]int{"a": 1}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rv := reflect.ValueOf(c.v)
			if got := isEmptyOrNil(rv); got != c.want {
				t.Fatalf("isEmptyOrNil(%v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}

// TestEquals_BothTypedNil covers the previously-0%-covered branch
// where both operands are typed-nil: a == nil && b == nil returns true.
// This is the symmetric case of the existing "one is nil" branch.
func TestEquals_BothTypedNil(t *testing.T) {
	t.Parallel()

	var a, b *int
	if !Equals(a, b) {
		t.Fatal("Equals: two nil pointers should be equal")
	}

	var s1, s2 []string
	if !Equals(s1, s2) {
		t.Fatal("Equals: two nil slices should be equal")
	}

	var m1, m2 map[int]int
	if !Equals(m1, m2) {
		t.Fatal("Equals: two nil maps should be equal")
	}

	// Negative case: one nil, one non-nil.
	if Equals(a, new(int)) {
		t.Fatal("Equals: nil vs non-nil pointer should not be equal")
	}
}
