// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package bson

import (
	"encoding/json"
	"github.com/qianwj/typed/mongo/util"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderedMap based on `bson.E` is an ordered map, unlike `bson.M` which is unordered.
// Additionally, OrderedMap can be converted to and from `bson.D` and `bson.M` as well.
type OrderedMap struct {
	entries []Entry
	dict    map[string]int
}

func D(pairs ...Entry) *OrderedMap {
	return NewMap(pairs...)
}

func NewMap(pairs ...Entry) *OrderedMap {
	m := &OrderedMap{
		entries: make([]Entry, 0),
		dict:    make(map[string]int),
	}
	for _, pair := range pairs {
		m.Put(pair.Key, pair.Value)
	}
	return m
}

func FromM(m rawbson.M) *OrderedMap {
	f := NewMap()
	for k, v := range m {
		f.Put(k, v)
	}
	return f
}

func FromD(d rawbson.D) *OrderedMap {
	f := NewMap()
	for _, e := range d {
		f.Put(e.Key, e.Value)
	}
	return f
}

func (m *OrderedMap) Put(key string, value any) *OrderedMap {
	if index, exists := m.dict[key]; exists {
		m.entries[index].Value = value
	} else {
		entry := Entry{Key: key, Value: value}
		m.entries = append(m.entries, entry)
		m.dict[key] = len(m.entries) - 1
	}
	return m
}

func (m *OrderedMap) Merge(other *OrderedMap) *OrderedMap {
	for _, entry := range other.entries {
		m.Put(entry.Key, entry.Value)
	}
	return m
}

func (m *OrderedMap) PutAsArray(key string, others ...*OrderedMap) *OrderedMap {
	validOthers := util.Filter(others, func(it *OrderedMap) bool {
		return !it.IsEmpty()
	})
	if len(validOthers) == 0 {
		return m
	}
	val, exist := m.Get(key)
	if !exist {
		m.Put(key, validOthers)
	} else {
		switch val.(type) {
		case rawbson.M:
			val = append([]*OrderedMap{FromM(val.(rawbson.M))}, validOthers...)
		case rawbson.D:
			val = append([]*OrderedMap{FromD(val.(rawbson.D))}, validOthers...)
		case rawbson.A:
			arr := make([]any, len(val.(rawbson.A)))
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
			for _, other := range validOthers {
				arr = append(arr, other)
			}
			val = arr
		case []*OrderedMap:
			val = append(val.([]*OrderedMap), validOthers...)
		}
		m.Put(key, val)
	}
	return m
}

func (m *OrderedMap) PutAsHash(key, hashKey string, val any) *OrderedMap {
	prev, exist := m.Get(key)
	if !exist {
		m.Put(key, NewMap(E(hashKey, val)))
	} else {
		prev.(*OrderedMap).Put(hashKey, val)
		m.Put(key, prev)
	}
	return m
}

func (m *OrderedMap) Get(key string) (any, bool) {
	if index, exists := m.dict[key]; exists {
		val := m.entries[index].Value
		if _, isNull := val.(primitive.Null); isNull {
			return nil, true
		}
		return val, true
	}
	return nil, false
}

func (m *OrderedMap) IsEmpty() bool {
	return len(m.entries) == 0
}

func (m *OrderedMap) Entries() []Entry {
	entries := make([]Entry, len(m.entries))
	for i, entry := range m.entries {
		entries[i] = entry
	}
	return entries
}

func (m *OrderedMap) Primitive() primitive.D {
	return util.OrderedMap(m.entries, func(e Entry) primitive.E {
		return e.Primitive()
	})
}

func (m *OrderedMap) ToMap() map[string]any {
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

func (m *OrderedMap) UnmarshalJSON(bytes []byte) error {
	var bMap rawbson.M
	if err := json.Unmarshal(bytes, &bMap); err != nil {
		return err
	}
	newMap := FromM(bMap)
	m.entries = newMap.entries
	m.dict = newMap.dict
	return nil
}

func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToMap())
}

func (m *OrderedMap) UnmarshalBSON(bytes []byte) error {
	var doc rawbson.D
	if err := rawbson.Unmarshal(bytes, &doc); err != nil {
		return err
	}
	newMap := FromD(doc)
	m.entries = newMap.entries
	m.dict = newMap.dict
	return nil
}

func (m *OrderedMap) MarshalBSON() ([]byte, error) {
	return rawbson.Marshal(m.Primitive())
}

func (m *OrderedMap) d2m(d rawbson.D) map[string]any {
	res := make(map[string]any)
	for _, e := range d {
		res[e.Key] = e
		switch e.Value.(type) {
		case rawbson.A:
			res[e.Key] = m.a2m(e.Value.(rawbson.A))
		case rawbson.M:
			res[e.Key] = FromM(e.Value.(rawbson.M)).ToMap()
		case rawbson.D:
			res[e.Key] = FromD(e.Value.(rawbson.D)).ToMap()
		case primitive.Null:
			res[e.Key] = nil
		case primitive.Regex:
			regex := e.Value.(primitive.Regex)
			res[e.Key] = "/" + regex.Pattern + "/" + regex.Options
		case *OrderedMap:
			res[e.Key] = e.Value.(*OrderedMap).ToMap()
		default:
			res[e.Key] = e.Value
		}
	}
	return res
}

func (m *OrderedMap) a2m(a rawbson.A) []any {
	arr := make([]any, len(a))
	for i, d := range a {
		switch d.(type) {
		case rawbson.D:
			arr[i] = m.d2m(d.(rawbson.D))
		case *OrderedMap:
			arr[i] = m.d2m(d.(*OrderedMap).Primitive())
		default:
			arr[i] = d
		}
	}
	return arr
}
