package filter

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

func ElemMatch(sub *Filter) *Filter {
	return New().ElemMatch(sub)
}

func (f *Filter) All(key string, items []any) *Filter {
	f.put(key, bson.M{operator.All: items})
	return f
}

func (f *Filter) Size(key string, size int64) *Filter {
	f.put(key, bson.M{operator.Size: size})
	return f
}

func (f *Filter) ElemMatch(sub *Filter) *Filter {
	f.put(operator.ElemMatch, bson.D(sub.entries))
	return f
}
