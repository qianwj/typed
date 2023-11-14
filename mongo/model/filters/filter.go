package filters

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/util"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter struct {
	data *bson.OrderedMap
}

func From(m rawbson.M) *Filter {
	return &Filter{data: bson.FromM(m)}
}

func Cast(d rawbson.D) *Filter {
	return &Filter{data: bson.FromD(d)}
}

func New() *Filter {
	return &Filter{data: bson.NewMap()}
}

func (f *Filter) Append(other *Filter) *Filter {
	for _, entry := range other.data.Entries() {
		f.data.Put(entry.Key, entry.Value)
	}
	return f
}

func (f *Filter) Get(key string) (any, bool) {
	return f.data.Get(key)
}

func (f *Filter) ToMap() map[string]any {
	return f.data.ToMap()
}

func (f *Filter) Primitive() primitive.D {
	return f.data.Primitive()
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	return f.data.MarshalJSON()
}

func (f *Filter) UnmarshalJSON(bytes []byte) error {
	if f.data == nil {
		f.data = bson.NewMap()
	}
	return f.data.UnmarshalJSON(bytes)
}

func (f *Filter) MarshalBSON() ([]byte, error) {
	return f.data.MarshalBSON()
}

func (f *Filter) UnmarshalBSON(bytes []byte) error {
	if f.data == nil {
		f.data = bson.NewMap()
	}
	return f.data.UnmarshalBSON(bytes)
}

func (f *Filter) putAsArray(key string, others ...*Filter) *Filter {
	f.data.PutAsArray(key, util.OrderedMap(others, func(it *Filter) *bson.OrderedMap {
		return it.data
	})...)
	return f
}
