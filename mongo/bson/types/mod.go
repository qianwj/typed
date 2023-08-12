package types

import (
	"github.com/qianwj/typed/mongo/bson"
	"reflect"
)

var (
	Array    = reflect.TypeOf((*bson.Array)(nil)).Elem()
	Document = reflect.TypeOf((*bson.Document)(nil)).Elem()
	Map      = reflect.TypeOf((*bson.Map)(nil)).Elem()
)
