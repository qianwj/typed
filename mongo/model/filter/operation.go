package filter

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *Filter) Eq(key string, val any) *Filter {
	f.put(key, val)
	return f
}

func (f *Filter) Gt(key string, val any) *Filter {
	f.put(key, bson.M{operator.Gt: val})
	return f
}

func (f *Filter) Gte(key string, val any) *Filter {
	f.put(key, bson.M{operator.Gte: val})
	return f
}

func (f *Filter) In(key string, items []any) *Filter {
	f.put(key, bson.M{operator.In: items})
	return f
}

func (f *Filter) Lt(key string, val any) *Filter {
	f.put(key, bson.M{operator.Lt: val})
	return f
}

func (f *Filter) Lte(key string, val any) *Filter {
	f.put(key, bson.M{operator.Lte: val})
	return f
}

func (f *Filter) Nin(key string, items []any) *Filter {
	f.put(key, bson.M{operator.Nin: items})
	return f
}

func (f *Filter) Ne(key string, val any) *Filter {
	f.put(key, bson.M{operator.Ne: val})
	return f
}

func (f *Filter) Exists(key string, val bool) *Filter {
	f.put(key, bson.M{operator.Exists: val})
	return f
}

func (f *Filter) Type(key string, val *model.DataType) *Filter {
	f.put(key, bson.M{operator.Type: val.Order()})
	return f
}

func (f *Filter) And(others ...*Filter) *Filter {
	f.putAsArray(operator.And, others...)
	return f
}

func (f *Filter) Not(key string, sub *Filter) *Filter {
	f.put(key, bson.M{operator.Not: bson.D(sub.entries)})
	return f
}

func (f *Filter) Nor(others ...*Filter) *Filter {
	f.putAsArray(operator.Nor, others...)
	return f
}

func (f *Filter) Or(others ...*Filter) *Filter {
	f.putAsArray(operator.Or, others...)
	return f
}
