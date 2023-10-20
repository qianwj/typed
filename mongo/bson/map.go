package bson

import (
	"encoding/json"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Map struct {
	entries []Entry
	dict    map[string]int
}

func NewMap() *Map {
	return &Map{
		entries: make([]Entry, 0),
		dict:    make(map[string]int),
	}
}

func FromM(m bson.M) *Map {
	f := NewMap()
	for k, v := range m {
		f.Put(k, v)
	}
	return f
}

func FromD(d bson.D) *Map {
	f := NewMap()
	for _, e := range d {
		f.Put(e.Key, e.Value)
	}
	return f
}

func (m *Map) Put(key string, value any) *Map {
	if index, exists := m.dict[key]; exists {
		// 键已存在，更新值
		m.entries[index].Value = value
	} else {
		// 键不存在，添加新的键值对
		entry := Entry{Key: key, Value: value}
		m.entries = append(m.entries, entry)
		m.dict[key] = len(m.entries) - 1
	}
	return m
}

func (m *Map) Raw() bson.D {
	return util.Map(m.entries, func(e Entry) bson.E {
		return bson.E(e)
	})
}

func (m *Map) ToMap() map[string]any {
	dst := make(map[string]any)
	for _, entry := range m.entries {
		switch entry.Value.(type) {
		case bson.D:
			dst[entry.Key] = m.d2m(entry.Value.(bson.D))
		case bson.A:
			dst[entry.Key] = m.a2m(entry.Value.(bson.A))
		default:
			dst[entry.Key] = entry.Value
		}
	}
	return dst
}

func (m *Map) UnmarshalJSON(bytes []byte) error {
	var bMap bson.M
	if err := json.Unmarshal(bytes, &bMap); err != nil {
		return err
	}
	newMap := FromM(bMap)
	m.entries = newMap.entries
	m.dict = newMap.dict
	return nil
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToMap())
}

func (m *Map) UnmarshalBSON(bytes []byte) error {
	var doc bson.D
	if err := bson.Unmarshal(bytes, doc); err != nil {
		return err
	}
	newMap := FromD(doc)
	m.entries = newMap.entries
	m.dict = newMap.dict
	return nil
}

func (m *Map) MarshalBSON() ([]byte, error) {
	return bson.Marshal(m.Raw())
}

func (m *Map) d2m(d bson.D) bson.M {
	res := bson.M{}
	for _, e := range d {
		res[e.Key] = e
		switch e.Value.(type) {
		case bson.A:
			res[e.Key] = m.a2m(e.Value.(bson.A))
		case bson.D:
			res[e.Key] = m.d2m(e.Value.(bson.D))
		case primitive.Null:
			res[e.Key] = nil
		case primitive.Regex:
			regex := e.Value.(primitive.Regex)
			res[e.Key] = "/" + regex.Pattern + "/" + regex.Options
		case primitive.ObjectID:
			res[e.Key] = e.Value.(primitive.ObjectID)
		case *Map:
			res[e.Key] = e.Value.(*Map).ToMap()
		default:
			res[e.Key] = e.Value
		}
	}
	return res
}

func (m *Map) a2m(a bson.A) []any {
	arr := make([]any, len(a))
	for i, d := range a {
		switch d.(type) {
		case bson.D:
			arr[i] = m.d2m(d.(bson.D))
		case *Map:
			arr[i] = m.d2m(d.(*Map).Raw())
		default:
			arr[i] = d
		}
	}
	return arr
}
