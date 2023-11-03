package filters

import (
	"github.com/qianwj/typed/mongo/model/filters/not"
	"github.com/qianwj/typed/mongo/operator"
	rawbson "go.mongodb.org/mongo-driver/bson"
)

func And(others ...*Filter) *Filter {
	return New().And(others...)
}

func Not(key string, exp *not.OperatorExpression) *Filter {
	return New().Not(key, exp)
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

func (f *Filter) Not(key string, exp *not.OperatorExpression) *Filter {
	f.data.Put(key, rawbson.M{operator.Not: exp})
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
