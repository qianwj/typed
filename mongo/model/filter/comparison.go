package filter

import (
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
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
