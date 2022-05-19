package core

type Object interface {
	Equals(o Object) bool
}

type Any interface {
	Object | comparable
}
