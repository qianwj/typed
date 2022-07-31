package mongo

import (
	raw "go.mongodb.org/mongo-driver/mongo"
)

var client *raw.Client

//func FindOne[D model.Document](ctx context.Context, filter model.Filter, opts ...*options.FindOneOptions) (*D, error) {
//	collection := NewTypedCollection()
//}
