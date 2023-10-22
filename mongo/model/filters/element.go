package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func Exists(key string, val bool) *Filter {
	return New().Exists(key, val)
}

func Type(key string, val *DataType) *Filter {
	return New().Type(key, val)
}

func (f *Filter) Exists(key string, val bool) *Filter {
	f.data.Put(key, bson.M{operator.Exists: val})
	return f
}

func (f *Filter) Type(key string, val *DataType) *Filter {
	f.data.Put(key, bson.M{operator.Type: val.Order()})
	return f
}
