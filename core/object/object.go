package object

type Object interface {
	Equals(o Object) bool
}

func Equals[T any](a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch a.(type) {
	case Object:
		return a.(Object).Equals(b.(Object))
	default:
		return a == b
	}
}
