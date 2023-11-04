package bson

import "go.mongodb.org/mongo-driver/bson"

type Array bson.A

func A(elements ...any) Array {
	a := make(Array, len(elements))
	for i, element := range elements {
		a[i] = element
	}
	return a
}
