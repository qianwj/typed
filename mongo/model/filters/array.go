package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func All(key string, items []any) *Filter {
	return New().All(key, items)
}

func Size(key string, size int64) *Filter {
	return New().Size(key, size)
}

func ElemMatch(key string, sub *Filter) *Filter {
	return New().ElemMatch(key, sub)
}

func ElemMatchWithInterval(key string, interval *Interval) *Filter {
	return New().ElemMatchWithInterval(key, interval)
}

func (f *Filter) All(key string, items []any) *Filter {
	f.data.Put(key, bson.M{operator.All: items})
	return f
}

func (f *Filter) Size(key string, size int64) *Filter {
	f.data.Put(key, bson.M{operator.Size: size})
	return f
}

func (f *Filter) ElemMatch(key string, sub *Filter) *Filter {
	f.data.Put(key, bson.M{operator.ElemMatch: sub.data.Raw()})
	return f
}

func (f *Filter) ElemMatchWithInterval(key string, interval *Interval) *Filter {
	f.data.Put(key, bson.M{operator.ElemMatch: interval.query()})
	return f
}
