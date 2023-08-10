package filters

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Null = primitive.Null{}
	Nil  = Null
)

type Filter struct {
	entries []bson.E
	dict    map[string]int
}

func New() *Filter {
	return &Filter{
		entries: make([]bson.E, 0),
		dict:    make(map[string]int),
	}
}

func From(m bson.M) *Filter {
	f := New()
	for k, v := range m {
		f.put(k, v)
	}
	return f
}

func Cast(d bson.D) *Filter {
	f := New()
	for _, e := range d {
		f.put(e.Key, e.Value)
	}
	return f
}

func (f *Filter) Append(other *Filter) *Filter {
	for _, entry := range other.entries {
		f.put(entry.Key, entry.Value)
	}
	return f
}

func (f *Filter) put(key string, value any) {
	if index, exists := f.dict[key]; exists {
		// 键已存在，更新值
		f.entries[index].Value = value
	} else {
		// 键不存在，添加新的键值对
		entry := bson.E{Key: key, Value: value}
		f.entries = append(f.entries, entry)
		f.dict[key] = len(f.entries) - 1
	}
}

func (f *Filter) putAsArray(key string, others ...*Filter) {
	val, exist := f.get(key)
	if !exist {
		f.put(key, _map(others, func(it *Filter) bson.D {
			return it.entries
		}))
	} else {
		f.put(
			key,
			append(
				val.(bson.A),
				_map(others, func(it *Filter) bson.D {
					return it.entries
				})...,
			),
		)
	}
}

func (f *Filter) get(key string) (any, bool) {
	if index, exists := f.dict[key]; exists {
		// 键存在，返回对应值
		return f.entries[index].Value, true
	}
	// 键不存在
	return nil, false
}

func (f *Filter) Document() bson.M {
	return f.d2m(f.entries)
}

func (f *Filter) Marshal() bson.D {
	return f.entries
}

func (f *Filter) d2m(d bson.D) bson.M {
	res := bson.M{}
	for _, e := range d {
		res[e.Key] = e
		switch e.Value.(type) {
		case bson.A:
			res[e.Key] = f.a2m(e.Value.(bson.A))
		case bson.D:
			res[e.Key] = f.d2m(e.Value.(bson.D))
		case primitive.Null:
			res[e.Key] = "null"
		case primitive.Regex:
			regex := e.Value.(primitive.Regex)
			res[e.Key] = "/" + regex.Pattern + "/" + regex.Options
		case primitive.ObjectID:
			res[e.Key] = e.Value.(primitive.ObjectID).String()
		default:
			res[e.Key] = e.Value
		}
	}
	return res
}

func (f *Filter) a2m(a bson.A) []bson.M {
	arr := make([]bson.M, len(a))
	for i, d := range a {
		arr[i] = f.d2m(d.(bson.D))
	}
	return arr
}
