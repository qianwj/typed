package filter

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/regex"
)

func Eq(key string, val any) *Filter {
	return New().Eq(key, val)
}

func Gt(key string, val any) *Filter {
	return New().Gt(key, val)
}

func Gte(key string, val any) *Filter {
	return New().Gte(key, val)
}

func In(key string, items []any) *Filter {
	return New().In(key, items)
}

func Lt(key string, val any) *Filter {
	return New().Lt(key, val)
}

func Lte(key string, val any) *Filter {
	return New().Lte(key, val)
}

func Nin(key string, items []any) *Filter {
	return New().Nin(key, items)
}

func Ne(key string, val any) *Filter {
	return New().Ne(key, val)
}

func Exists(key string, val bool) *Filter {
	return New().Exists(key, val)
}

func Type(key string, val *model.DataType) *Filter {
	return New().Type(key, val)
}

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

func Expr(expression any) *Filter {
	return New().Expr(expression)
}

func Mod(key string, divisor, remainder float64) *Filter {
	return New().Mod(key, divisor, remainder)
}

func Like(key string, matcher *regex.Matcher) *Filter {
	return New().Like(key, matcher)
}

func Where(key, expression string) *Filter {
	return New().Where(key, expression)
}

func All(key string, items []any) *Filter {
	return New().All(key, items)
}

func Size(key string, size int64) *Filter {
	return New().Size(key, size)
}

func ElemMatch(sub *Filter) *Filter {
	return New().ElemMatch(sub)
}
