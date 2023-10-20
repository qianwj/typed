package bson

import (
	"encoding/json"
	"github.com/qianwj/typed/mongo/util"
	rawbson "go.mongodb.org/mongo-driver/bson"
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

func FromM(m rawbson.M) *Map {
	f := NewMap()
	for k, v := range m {
		f.Put(k, v)
	}
	return f
}

func FromD(d rawbson.D) *Map {
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

func (m *Map) PutAsArray(key string, others ...*Map) {
	val, exist := m.Get(key)
	if !exist {
		m.Put(key, others)
	} else {
		switch val.(type) {
		case rawbson.M:
			val = append([]*Map{FromM(val.(rawbson.M))}, others...)
		case rawbson.D:
			val = append([]*Map{FromD(val.(rawbson.D))}, others...)
		case rawbson.A:
			arr := make([]any, 0)
			for i, elem := range val.(rawbson.A) {
				switch elem.(type) {
				case rawbson.M:
					arr[i] = FromM(elem.(rawbson.M))
				case rawbson.D:
					arr[i] = FromD(elem.(rawbson.D))
				default:
					arr[i] = elem
				}
			}
			for _, other := range others {
				arr = append(arr, other)
			}
			val = arr
		}
		m.Put(key, val)
	}
}

func (m *Map) Get(key string) (any, bool) {
	if index, exists := m.dict[key]; exists {
		return m.entries[index].Value, true
	}
	return nil, false
}

func (m *Map) Entries() []Entry {
	entries := make([]Entry, len(m.entries))
	for i, entry := range m.entries {
		entries[i] = entry
	}
	return entries
}

func (m *Map) Raw() rawbson.D {
	return util.Map(m.entries, func(e Entry) rawbson.E {
		return rawbson.E(e)
	})
}

func (m *Map) ToMap() map[string]any {
	dst := make(map[string]any)
	for _, entry := range m.entries {
		switch entry.Value.(type) {
		case primitive.Null:
			dst[entry.Key] = nil
		case rawbson.D:
			dst[entry.Key] = m.d2m(entry.Value.(rawbson.D))
		case rawbson.A:
			dst[entry.Key] = m.a2m(entry.Value.(rawbson.A))
		default:
			dst[entry.Key] = entry.Value
		}
	}
	return dst
}

func (m *Map) UnmarshalJSON(bytes []byte) error {
	var bMap rawbson.M
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
	var doc rawbson.D
	if err := rawbson.Unmarshal(bytes, doc); err != nil {
		return err
	}
	newMap := FromD(doc)
	m.entries = newMap.entries
	m.dict = newMap.dict
	return nil
}

func (m *Map) MarshalBSON() ([]byte, error) {
	return rawbson.Marshal(m.Raw())
}

func (m *Map) d2m(d rawbson.D) map[string]any {
	res := make(map[string]any)
	for _, e := range d {
		res[e.Key] = e
		switch e.Value.(type) {
		case rawbson.A:
			res[e.Key] = m.a2m(e.Value.(rawbson.A))
		case rawbson.D:
			res[e.Key] = m.d2m(e.Value.(rawbson.D))
		case primitive.Null:
			res[e.Key] = nil
		case primitive.Regex:
			regex := e.Value.(primitive.Regex)
			res[e.Key] = "/" + regex.Pattern + "/" + regex.Options
		case *Map:
			res[e.Key] = e.Value.(*Map).ToMap()
		default:
			res[e.Key] = e.Value
		}
	}
	return res
}

func (m *Map) a2m(a rawbson.A) []any {
	arr := make([]any, len(a))
	for i, d := range a {
		switch d.(type) {
		case rawbson.D:
			arr[i] = m.d2m(d.(rawbson.D))
		case *Map:
			arr[i] = m.d2m(d.(*Map).Raw())
		default:
			arr[i] = d
		}
	}
	return arr
}
