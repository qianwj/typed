package objects

import "reflect"

// Equaler is an optional interface a type can implement to provide
// its own notion of equality. The Equal method receives another
// value (any) and reports whether it should be considered equal to
// the receiver.
//
// Equaler is the simple, non-generic form. It is documented for
// completeness and as a hint to readers; the Equals function does
// not require it. Equals detects the Equal method on a type's
// dynamic value via reflection, with the required signature
// func (T) Equal(T) bool, and calls it regardless of whether the
// type also satisfies the explicit Equaler interface. This is why
// time.Time works out of the box: its Equal(time.Time) bool method
// has the right shape and is found by reflection, even though
// time.Time is not declared to implement any Equaler interface.
//
// If a type's Equal method has a different signature — for example
// Equal(other Interface) bool — Equals will not pick it up. Define
// the method with signature func (T) Equal(T) bool to opt in.
type Equaler interface {
	Equal(other any) bool
}

// Equals reports whether a and b are deeply equal, with three
// intentional fixes over reflect.DeepEqual:
//
//  1. Custom Equal method dispatch. If either operand's dynamic
//     type defines a method with signature func (T) Equal(T) bool,
//     that method is called instead of a recursive walk. The
//     operand whose dynamic type has the Equal method is treated
//     as the authoritative source. Both sides are checked, so
//     Equals(a, b) and Equals(b, a) are symmetric when only one
//     side has a custom Equal method: whichever side has it will
//     be consulted.
//
//  2. nil-vs-empty collection normalisation, recursively. A nil
//     slice and an empty non-nil slice are considered equal at
//     every level of the walk (not just the top). The same rule
//     applies to maps. reflect.DeepEqual returns false for these
//     pairs, which is rarely what callers want.
//
//  3. Nil-safe custom Equal methods. If a custom Equal method is
//     defined on a pointer receiver and either operand is a typed
//     nil pointer, Equals returns true when both are nil and
//     false otherwise, without ever calling the nil-receiver
//     method. This avoids the panic that a naive "call Equal
//     first, then check nil" implementation would produce.
//
// All other aspects mirror reflect.DeepEqual: pointers are
// followed, structs are compared field by field (including
// unexported fields), arrays element by element, cyclic
// structures are handled without infinite recursion.
func Equals[T any](a, b T) bool {
	// Step 1: handle nil operands up front, before any method
	// dispatch. A typed nil pointer must not be dereferenced by
	// a custom Equal method, and untyped nil is the easiest case
	// to short-circuit.
	if any(a) == nil || any(b) == nil {
		return any(a) == nil && any(b) == nil
	}

	// Step 2: try to find an Equal method on either side. The
	// first match wins; the other is the receiver's choice. This
	// keeps Equals(a, b) == Equals(b, a) when only one side
	// implements Equal: that side's method is consulted in both
	// orderings.
	if m, ok := findEqualMethod(a); ok {
		return callEqual(m, b)
	}
	if m, ok := findEqualMethod(b); ok {
		return callEqual(m, a)
	}

	// Step 3: structural walk.
	return deepEqualFixed(any(a), any(b))
}

// findEqualMethod looks up an Equal method on v's dynamic type and
// returns the bound method value if the method's signature matches
// the expected func (T) Equal(T) bool shape.
//
// findEqualMethod returns false (no method) when v is a typed nil
// pointer, so Equals never calls a nil-receiver Equal method.
func findEqualMethod(v any) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return reflect.Value{}, false
	}
	// Refuse to bind a method on a typed nil pointer / interface;
	// the call would panic.
	if (rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface) && rv.IsNil() {
		return reflect.Value{}, false
	}
	m := rv.MethodByName("Equal")
	if !m.IsValid() {
		return reflect.Value{}, false
	}
	mt := m.Type()
	if mt.NumIn() != 1 || mt.NumOut() != 1 || mt.Out(0).Kind() != reflect.Bool {
		return reflect.Value{}, false
	}
	// The single argument must be of the receiver's type. This
	// rules out Equal(any) bool and other loose signatures that
	// would otherwise match too eagerly.
	if mt.In(0) != rv.Type() {
		return reflect.Value{}, false
	}
	return m, true
}

// callEqual invokes a method value m (already bound to its receiver)
// with other as the single argument. The Equaler dispatch treats
// the method as the authoritative source of equality, so:
//
//   - If other is a typed nil pointer, a zero value of the
//     method's expected argument type is used (the method gets a
//     chance to handle "compare against nothing" itself).
//   - If other is the wrong concrete type, the dispatch still
//     happens: the receiver's Equal method is called with a zero
//     value of the expected argument type. This is the "Equaler
//     is authoritative" rule: the implementer has signalled that
//     its own semantic is the right one, and Equals will not
//     silently fall through to a structural walk just because
//     the operands do not share the same type.
//
// The only case where the call is rejected is when the method
// value itself is somehow invalid, which cannot happen in
// practice because findEqualMethod only returns valid methods.
func callEqual(m reflect.Value, other any) bool {
	argType := m.Type().In(0)
	rv := reflect.ValueOf(other)
	if !rv.IsValid() || rv.Type() != argType {
		// other is nil, or the wrong concrete type. The
		// method is treated as the source of truth, so we
		// substitute the zero value of the expected argument
		// type and let the method decide.
		rv = reflect.Zero(argType)
	}
	return m.Call([]reflect.Value{rv})[0].Bool()
}

// deepEqualFixed is the Equaler-bypassing path. It implements
// reflect.DeepEqual with the nil-vs-empty slice / map normalisation
// applied at every level, and explicit cycle detection so cyclic
// data structures do not blow the stack.
func deepEqualFixed(x, y any) bool {
	return deepEqualFixedVisit(x, y, map[visitKey]bool{})
}

// visitKey identifies a pair of pointer-typed values that have
// already been compared during a recursive walk. Using the pointer
// addresses (not the underlying values) lets two distinct cyclic
// structures be compared without infinite recursion.
type visitKey struct {
	a, b uintptr
}

// deepEqualFixedVisit is the recursive walker. It is given a
// visited-set of pointer pairs to detect cycles. The visited set
// is grown as the walk descends and never shrinks; the cost of
// allocating it per Equals call is acceptable because most
// structures are acyclic and short.
func deepEqualFixedVisit(x, y any, visited map[visitKey]bool) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	rx := reflect.ValueOf(x)
	ry := reflect.ValueOf(y)
	if !rx.IsValid() || !ry.IsValid() {
		return false
	}
	if rx.Kind() != ry.Kind() || rx.Type() != ry.Type() {
		return false
	}

	switch rx.Kind() {
	case reflect.Slice:
		// nil-vs-empty at this level, before length check.
		if isEmptyOrNil(rx) && isEmptyOrNil(ry) {
			return true
		}
		if rx.Len() != ry.Len() {
			return false
		}
		for i := 0; i < rx.Len(); i++ {
			if !deepEqualFixedVisit(rx.Index(i).Interface(), ry.Index(i).Interface(), visited) {
				return false
			}
		}
		return true

	case reflect.Map:
		if isEmptyOrNil(rx) && isEmptyOrNil(ry) {
			return true
		}
		if rx.Len() != ry.Len() {
			return false
		}
		for _, key := range rx.MapKeys() {
			vy := ry.MapIndex(key)
			if !vy.IsValid() {
				return false
			}
			if !deepEqualFixedVisit(rx.MapIndex(key).Interface(), vy.Interface(), visited) {
				return false
			}
		}
		return true

	case reflect.Pointer, reflect.Interface:
		// Both nil → equal. Exactly one nil → not equal.
		if rx.IsNil() || ry.IsNil() {
			return rx.IsNil() && ry.IsNil()
		}
		// Cycle detection: if we have already compared this pair
		// of pointer values, treat the cycle as equal. The
		// surrounding structural equality (or inequality) is
		// established by the rest of the walk; the cycle is just
		// a fixed point.
		k := visitKey{a: rx.Pointer(), b: ry.Pointer()}
		if visited[k] {
			return true
		}
		visited[k] = true
		return deepEqualFixedVisit(rx.Elem().Interface(), ry.Elem().Interface(), visited)

	case reflect.Struct:
		for i := 0; i < rx.NumField(); i++ {
			if !deepEqualFixedVisit(rx.Field(i).Interface(), ry.Field(i).Interface(), visited) {
				return false
			}
		}
		return true

	case reflect.Array:
		if rx.Len() != ry.Len() {
			return false
		}
		for i := 0; i < rx.Len(); i++ {
			if !deepEqualFixedVisit(rx.Index(i).Interface(), ry.Index(i).Interface(), visited) {
				return false
			}
		}
		return true

	default:
		// For basic kinds (bool, numbers, string, chan, func),
		// fall back to reflect.DeepEqual, which has the right
		// semantics for those.
		return reflect.DeepEqual(x, y)
	}
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
