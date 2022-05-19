package core

type Object interface {
	Equals(o Object) bool
}

type Any interface {
	Object | comparable
}

type Comparable interface {
	Compare(other Any) int
}
