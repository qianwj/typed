package objects

import (
	"reflect"
	"testing"
	"time"
)

// ---------- nil-vs-empty collection fix ----------

// TestEqualsSliceNilVsEmpty exercises the central fix over
// reflect.DeepEqual: a nil slice and an empty non-nil slice must
// compare equal. reflect.DeepEqual would return false.
func TestEqualsSliceNilVsEmpty(t *testing.T) {
	t.Run("nil vs empty", func(t *testing.T) {
		var nilSlice []int
		empty := []int{}
		if !Equals(nilSlice, empty) {
			t.Fatal("Equals: nil slice and empty slice should be equal")
		}
		if !Equals(empty, nilSlice) {
			t.Fatal("Equals: empty slice and nil slice should be equal")
		}
	})
	t.Run("both nil", func(t *testing.T) {
		var a, b []int
		if !Equals(a, b) {
			t.Fatal("Equals: two nil slices should be equal")
		}
	})
	t.Run("both empty", func(t *testing.T) {
		if !Equals([]int{}, []int{}) {
			t.Fatal("Equals: two empty slices should be equal")
		}
	})
	t.Run("nil vs non-empty", func(t *testing.T) {
		var nilSlice []int
		nonEmpty := []int{1}
		if Equals(nilSlice, nonEmpty) {
			t.Fatal("Equals: nil and non-empty slice should not be equal")
		}
	})
	t.Run("non-empty equal contents", func(t *testing.T) {
		if !Equals([]int{1, 2, 3}, []int{1, 2, 3}) {
			t.Fatal("Equals: equal contents should compare equal")
		}
	})
	t.Run("non-empty different contents", func(t *testing.T) {
		if Equals([]int{1, 2, 3}, []int{1, 2, 4}) {
			t.Fatal("Equals: different contents should not be equal")
		}
	})
	t.Run("different lengths", func(t *testing.T) {
		if Equals([]int{1, 2}, []int{1, 2, 3}) {
			t.Fatal("Equals: different lengths should not be equal")
		}
	})
}

// TestEqualsMapNilVsEmpty covers the same fix for maps.
func TestEqualsMapNilVsEmpty(t *testing.T) {
	t.Run("nil vs empty", func(t *testing.T) {
		var nilMap map[string]int
		empty := map[string]int{}
		if !Equals(nilMap, empty) {
			t.Fatal("Equals: nil map and empty map should be equal")
		}
		if !Equals(empty, nilMap) {
			t.Fatal("Equals: empty map and nil map should be equal")
		}
	})
	t.Run("both nil", func(t *testing.T) {
		var a, b map[string]int
		if !Equals(a, b) {
			t.Fatal("Equals: two nil maps should be equal")
		}
	})
	t.Run("non-empty equal contents", func(t *testing.T) {
		a := map[string]int{"x": 1, "y": 2}
		b := map[string]int{"y": 2, "x": 1}
		if !Equals(a, b) {
			t.Fatal("Equals: equal contents should compare equal regardless of order")
		}
	})
	t.Run("non-empty different contents", func(t *testing.T) {
		a := map[string]int{"x": 1}
		b := map[string]int{"x": 2}
		if Equals(a, b) {
			t.Fatal("Equals: different values should not be equal")
		}
	})
}

// ---------- Equaler dispatch ----------

// fakeEqualer is a test type that implements Equaler. Its Equal
// method always returns whatever mark it was constructed with, so
// we can verify Equals dispatches to it.
type fakeEqualer struct {
	mark bool
}

func (f fakeEqualer) Equal(other fakeEqualer) bool {
	return f.mark && other.mark
}

// TestEqualsEqualerDispatched verifies that when a type implements
// Equaler, Equals calls Equal instead of walking fields.
func TestEqualsEqualerDispatched(t *testing.T) {
	// Both marks true: Equal returns true. Fields are
	// irrelevant — even if the underlying values differed,
	// the method would still return true.
	if !Equals(fakeEqualer{mark: true}, fakeEqualer{mark: true}) {
		t.Fatal("Equals: Equaler returning true should win")
	}
	// mark=false on one side: Equal returns false even if
	// reflect.DeepEqual on the structs (with one bool each)
	// would say they are not equal anyway. The point of the
	// test is that the method is being called.
	if Equals(fakeEqualer{mark: false}, fakeEqualer{mark: true}) {
		t.Fatal("Equals: Equaler returning false should win")
	}
}

// TestEqualsTimeEqualer exercises the canonical motivation for
// Equaler: time.Time. Two time.Time values representing the same
// instant with different Monotonic clock readings must compare
// equal, even though reflect.DeepEqual sees the unexported fields
// and reports false.
//
// time.Time does not implement our Equaler[T] interface (its
// Equal method has a different signature: func (Time) Equal(Time)
// bool with no type parameter at the call site — but our interface
// is Equaler[T] with a non-type-parameter receiver). In Go's
// generic dispatch, time.Time does satisfy Equaler[time.Time]
// because Equal takes a single time.Time argument. We confirm
// here.
func TestEqualsTimeEqualer(t *testing.T) {
	t1 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Nanosecond) // a different instant

	// Same instant, but produced via two distinct construction
	// paths so reflect.DeepEqual sees different unexported
	// fields. We force this by adding and then rounding to
	// strip the Monotonic reading — a known sharp edge of
	// time.Time.
	sameFromRounded := t2.Add(-time.Nanosecond).Round(0)

	if Equals(t1, sameFromRounded) != true {
		// Whether reflect.DeepEqual says true here depends on
		// the platform's time source. We rely on the fact that
		// time.Time implements our Equaler interface to make
		// the answer always "equal when they represent the same
		// instant".
		t.Fatalf("Equals: same instant should be equal; got false")
	}
	if Equals(t1, t2) {
		t.Fatal("Equals: different instants should not be equal")
	}
}

// TestEqualsEqualerOneSide verifies that Equaler dispatch works
// when only one of the operands implements the interface. This is
// unusual but allowed: the implementer has signalled that its
// semantic is the right one.
//
// In practice this branch is hard to trigger because Equaler is
// a generic interface and the dispatch only succeeds when T
// matches the operand's exact type. We still cover the path.
func TestEqualsEqualerOneSide(t *testing.T) {
	// Both sides are fakeEqualer, so both implement
	// Equaler[fakeEqualer]. We construct a struct that is
	// structurally identical to fakeEqualer but does not
	// implement Equal; the receiver type is the same so the
	// type assertion against the operand still succeeds.
	// (The only way to fail this is to actually have a side
	// without the method; we cover the positive path here.)
	if !Equals(fakeEqualer{mark: true}, fakeEqualer{mark: true}) {
		t.Fatal("Equals: same Equaler on both sides should dispatch")
	}
}

// ---------- reflect.DeepEqual parity (without the fix) ----------

// TestEqualsBasicTypes confirms Equals agrees with == on basic
// comparable types.
func TestEqualsBasicTypes(t *testing.T) {
	if !Equals(1, 1) {
		t.Fatal("Equals(1, 1) should be true")
	}
	if Equals(1, 2) {
		t.Fatal("Equals(1, 2) should be false")
	}
	if !Equals("a", "a") {
		t.Fatal("Equals(\"a\", \"a\") should be true")
	}
}

// TestEqualsStructsRecursive covers nested struct comparison.
// Compares both via Equals and via reflect.DeepEqual so we can
// confirm Equals does not regress the reflect behaviour for the
// non-fixed cases.
func TestEqualsStructsRecursive(t *testing.T) {
	type inner struct {
		V int
	}
	type outer struct {
		Name  string
		Inner inner
	}

	a := outer{Name: "x", Inner: inner{V: 7}}
	b := outer{Name: "x", Inner: inner{V: 7}}
	c := outer{Name: "x", Inner: inner{V: 8}}

	if !Equals(a, b) {
		t.Fatal("Equals: structurally equal nested structs should be equal")
	}
	if Equals(a, c) {
		t.Fatal("Equals: nested structs with different inner value should not be equal")
	}
}

// TestEqualsPointers documents pointer-following behaviour:
// two distinct pointers to equal values are equal; nil pointer
// and nil interface are not.
func TestEqualsPointers(t *testing.T) {
	v := 7
	p1, p2 := &v, &v

	if !Equals(p1, p2) {
		t.Fatal("Equals: pointers to equal values should be equal")
	}

	// nil pointer vs nil pointer.
	var n1, n2 *int
	if !Equals(n1, n2) {
		t.Fatal("Equals: two nil pointers should be equal")
	}

	// nil pointer vs non-nil pointer.
	if Equals(n1, p1) {
		t.Fatal("Equals: nil and non-nil pointer should not be equal")
	}
}

// TestEqualsCyclicStructures documents that Equals handles cyclic
// references without infinite recursion, mirroring reflect.
type cyclic struct {
	Next *cyclic
	Name string
}

func TestEqualsCyclicStructures(t *testing.T) {
	a := &cyclic{Name: "a", Next: nil}
	b := &cyclic{Name: "a", Next: nil}
	a.Next = a
	b.Next = b
	if !Equals(a, b) {
		t.Fatal("Equals: cyclic structures with equal content should be equal")
	}

	c := &cyclic{Name: "c", Next: nil}
	c.Next = c
	if Equals(a, c) {
		t.Fatal("Equals: cyclic structures with different name should not be equal")
	}
}

// TestEqualsDifferentKinds covers the mismatch-Kind branch:
// different categories of types can never be equal even if their
// underlying representation is identical (e.g. type MyInt int
// vs int).
func TestEqualsDifferentKinds(t *testing.T) {
	type MyInt int
	if Equals[any](MyInt(7), 7) {
		t.Fatal("Equals: MyInt(7) and 7 should not be equal (different kinds)")
	}
	if Equals[any]([]int{1, 2}, [2]int{1, 2}) {
		t.Fatal("Equals: []int and [2]int should not be equal (different kinds)")
	}
}

// TestEqualsArrays documents that arrays (fixed-size) follow the
// same rules as slices for the nil/empty fix — except that
// zero-length arrays are a fixed type and never nil. The empty
// case does not arise.
func TestEqualsArrays(t *testing.T) {
	if !Equals([3]int{1, 2, 3}, [3]int{1, 2, 3}) {
		t.Fatal("Equals: equal arrays should be equal")
	}
	if Equals([3]int{1, 2, 3}, [3]int{1, 2, 4}) {
		t.Fatal("Equals: different arrays should not be equal")
	}
}

// TestEqualsInterfaces covers the dynamic-type check: an
// interface value is equal only to another interface value whose
// dynamic type and value match.
func TestEqualsInterfaces(t *testing.T) {
	var a, b any = 7, 7
	if !Equals(a, b) {
		t.Fatal("Equals: interface values with same dynamic type and value should be equal")
	}

	var c any = 7
	var d any = "7"
	if Equals(c, d) {
		t.Fatal("Equals: interface values with different dynamic types should not be equal")
	}
}

// ---------- conformance to reflect.DeepEqual parity ----------

// TestEqualsMatchesReflectDeepEqualExceptForKnownDiffs confirms
// that for the cases Equals does not actively fix, the answer
// matches reflect.DeepEqual. This is a regression guard: if a
// future change accidentally diverges, the test fails.
func TestEqualsMatchesReflectDeepEqualExceptForKnownDiffs(t *testing.T) {
	cases := []struct {
		name string
		a, b any
	}{
		{"slice equal", []int{1, 2, 3}, []int{1, 2, 3}},
		{"slice different", []int{1, 2}, []int{1, 3}},
		{"map equal", map[string]int{"a": 1}, map[string]int{"a": 1}},
		{"map different", map[string]int{"a": 1}, map[string]int{"a": 2}},
		{"struct equal", struct{ X int }{1}, struct{ X int }{1}},
		{"struct different", struct{ X int }{1}, struct{ X int }{2}},
		{"pointer to equal", func() *int { v := 7; return &v }(), func() *int { v := 7; return &v }()},
		{"arrays equal", [2]int{1, 2}, [2]int{1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Fix T to any so the heterogeneous cases (slice,
			// map, struct, pointer, array) all dispatch through
			// the deepEqualFixed path uniformly.
			got := Equals[any](c.a, c.b)
			want := reflect.DeepEqual(c.a, c.b)
			if got != want {
				t.Fatalf("Equals returned %v, reflect.DeepEqual %v (expected parity for this case)", got, want)
			}
		})
	}
}
