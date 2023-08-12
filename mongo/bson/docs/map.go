package docs

import (
	"encoding/json"
	"errors"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/bson/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

func From(m M) bson.Map {
	f := New()
	for k, v := range m {
		f.Put(k, v.(bson.Bson))
	}
	return f
}

func Cast(d D) bson.Map {
	f := New()
	for _, e := range d {
		f.Put(e.Key, e.Value.(bson.Bson))
	}
	return f
}

type bsonMap struct {
	entries []bson.E
	dict    map[string]int
}

func New() bson.Map {
	return &bsonMap{
		entries: make([]bson.E, 0),
		dict:    make(map[string]int),
	}
}

func (m *bsonMap) Cast() any {
	return m
}

func (m *bsonMap) Put(key string, value bson.Bson) bson.Map {
	if index, exists := m.dict[key]; exists {
		// 键已存在，更新值
		m.entries[index].Value = value
	} else {
		// 键不存在，添加新的键值对
		entry := bson.E{Key: key, Value: value}
		m.entries = append(m.entries, entry)
		m.dict[key] = len(m.entries) - 1
	}
	return m
}

func (m *bsonMap) Entries() []bson.Entry {
	entries := make([]bson.Entry, len(m.entries))
	for i, entry := range m.entries {
		entries[i] = entry
	}
	return entries
}

func (m *bsonMap) PutAll(other bson.Map) {
	for _, entry := range other.Entries() {
		m.Put(entry.GetKey(), entry.GetVal())
	}
}

func (m *bsonMap) Get(key string) (bson.Bson, bool) {
	if index, exists := m.dict[key]; exists {
		// 键存在，返回对应值
		return m.entries[index].GetVal(), true
	}
	// 键不存在
	return nil, false
}

func (m *bsonMap) ToMap() bson.Map {
	return m
}

func (m *bsonMap) Marshal() primitive.D {
	res := make(primitive.D, len(m.entries))
	for i, entry := range m.entries {
		res[i] = primitive.E{
			Key:   entry.GetKey(),
			Value: entry.GetVal(),
		}
	}
	return res
}

func (m *bsonMap) MarshalJSON() ([]byte, error) {
	res := map[string]any{}
	var err error
	for _, entry := range m.entries {
		res[entry.GetKey()], err = m.marshalBJSON(entry.GetVal())
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(res)
}

func (m *bsonMap) marshalBJSON(b any) (any, error) {
	switch b.(type) {
	case nil:
		return nil, nil
	case string:
		return b, nil
	case int, int8, int16, int32, int64, float32, float64, bool:
		return b, nil
	case bson.String, bson.Int, bson.Long, bson.Bool, bson.Float, bson.Double:
		return b, nil
	case bson.DateTime:
		return b.(bson.DateTime).String(), nil
	case bson.Null:
		return nil, nil
	case bson.ObjectId:
		return b.(bson.ObjectId).String(), nil
	default:
		bType := reflect.TypeOf(b)
		if bType.Implements(types.Array) {
			docs := make([]any, 0)
			iter := b.(bson.Array).Iter()
			if iter == nil {
				return nil, nil
			}
			for iter.HasNext() {
				next, ok := iter.Next()
				if !ok {
					return nil, errors.New("unknown type data")
				}
				doc, err := m.marshalBJSON(next)
				if err != nil {
					return nil, err
				}
				docs = append(docs, doc)
			}
			return docs, nil
		} else if bType.Implements(types.Document) {
			doc := b.(bson.Document).ToMap()
			return doc, nil
		}
	}
	return nil, errors.New("unknown type")
}

//func (m *bsonMap) d2m(d bson.D) bson.M {
//	res := bson.M{}
//	for _, e := range d {
//		res[e.Key] = e
//		switch e.Value.(type) {
//		case bson.A:
//			res[e.Key] = m.a2m(e.Value.(bson.A))
//		case bson.D:
//			res[e.Key] = m.d2m(e.Value.(bson.D))
//		case primitive.Null:
//			res[e.Key] = "null"
//		case primitive.Regex:
//			regex := e.Value.(primitive.Regex)
//			res[e.Key] = "/" + regex.Pattern + "/" + regex.Options
//		case primitive.ObjectID:
//			res[e.Key] = e.Value.(primitive.ObjectID).String()
//		default:
//			res[e.Key] = e.Value
//		}
//	}
//	return res
//}
//
//func (m *bsonMap) a2m(a primitive.A) []Document {
//	arr := make([]bson.M, len(a))
//	for i, d := range a {
//		arr[i] = m.d2m(d.(bson.D))
//	}
//	return arr
//}
