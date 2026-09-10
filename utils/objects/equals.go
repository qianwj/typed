package objects

import "reflect"

// Equaler is an optional interface a type can implement to provide
// its own notion of equality. The Equal method receives another
// value of the same type and reports whether it should be considered
// equal to the receiver.
//
// objects.Equals consults Equaler before falling back to a
// reflect-based deep comparison. This is the recommended path for
// types whose default reflect.DeepEqual behaviour is wrong:
//
//   - time.Time: its Monotonic clock reading makes the unexported
//     fields differ between two values that represent the same
//     instant. time.Time.Equal implements the correct semantic.
//   - any type that wants to ignore certain unexported fields or
//     do tolerance-based comparison (floating-point epsilon,
//     case-insensitive strings, …) can give Equals a precise rule
//     instead of relying on reflect.DeepEqual's structural walk.
//
// Equaler is a generic interface, so a type's Equal method must
// accept the exact same type. Type X implements Equaler[X], and
// Equals[X] picks it up automatically.
type Equaler[T any] interface {
	Equal(T) bool
}

// Equals reports whether a and b are deeply equal, with two
// intentional fixes over reflect.DeepEqual:
//
//  1. Equaler dispatch. If a's dynamic type implements
//     Equaler[T] for the function's type parameter T, that Equal
//     method is called instead of a recursive walk. This is the
//     override hook for types whose reflect.DeepEqual behaviour is
//     wrong (time.Time being the canonical example).
//
//  2. nil-vs-empty collection normalisation. A nil slice and an
//     empty non-nil slice are considered equal. The same rule
//     applies to maps. reflect.DeepEqual returns false for these
//     pairs, which is rarely what callers want.
//
// All other aspects mirror reflect.DeepEqual: pointers are
// followed, structs are compared field by field (including
// unexported fields), arrays element by element, cyclic
// structures are handled without infinite recursion.
func Equals[T any](a, b T) bool {
	// Equaler dispatch. Because the function signature requires
	// a and b to share the static type T, both either implement
	// Equaler[T] or neither does, so a single type assertion is
	// enough to detect the override.
	if eq, ok := any(a).(Equaler[T]); ok {
		return eq.Equal(b)
	}
	return deepEqualFixed(any(a), any(b))
}

// deepEqualFixed is the Equaler-bypassing path. It is
// reflect.DeepEqual with a single correction: nil and empty
// slices / maps are treated as equal.
func deepEqualFixed(x, y any) bool {
	// Fast path: both sides are untyped nil (or both typed nil
	// holding a nil interface). reflect.DeepEqual already
	// handles this correctly, so no fix is needed.
	if x == nil || y == nil {
		return x == nil && y == nil
	}

	rx := reflect.ValueOf(x)
	ry := reflect.ValueOf(y)

	// Mismatched Kinds: not equal. reflect.DeepEqual also
	// reports false here, but the check is cheap and lets us
	// skip the slice/map normalisation for non-collection types.
	if rx.Kind() != ry.Kind() {
		return false
	}

	// Mismatched Types: not equal, even when the Kind matches.
	// This catches the boxed-any case: Equals[any](MyInt(7), 7)
	// has Kind Int on both sides, but Type is MyInt vs int, and
	// reflect.DeepEqual would also report false.
	if rx.Type() != ry.Type() {
		return false
	}

	switch rx.Kind() {
	case reflect.Slice:
		if isEmptyOrNil(rx) && isEmptyOrNil(ry) {
			return true
		}
	case reflect.Map:
		if isEmptyOrNil(rx) && isEmptyOrNil(ry) {
			return true
		}
	}

	return reflect.DeepEqual(x, y)
}

// isEmptyOrNil reports whether v is a nil reference, or a
// zero-length slice / map (including nil). It only makes sense
// for slice and map kinds; the caller is expected to check.
func isEmptyOrNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	if v.IsNil() {
		return true
	}
	return v.Len() == 0
}
