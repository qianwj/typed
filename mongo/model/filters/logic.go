package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func And(others ...*Filter) *Filter {
	return New().And(others...)
}

func Not(key string, sub *Filter) *Filter {
	return New().Not(key, sub)
}

func Nor(others ...*Filter) *Filter {
	return New().Nor(others...)
}

func Or(others ...*Filter) *Filter {
	return New().Or(others...)
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
