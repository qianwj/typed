package bson

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestMap_PutAsArray(t *testing.T) {

	testCases := []struct {
		Expected bson.D
		Actual   *Map
	}{
		{Expected: bson.D{}, Actual: NewMap().PutAsArray("aa", NewMap())},
		{
			Expected: bson.D{{Key: "aa", Value: bson.A{"a", "b", bson.D{{Key: "c", Value: "d"}}}}},
			Actual: FromM(bson.M{
				"aa": bson.A{
					"a", "b",
				},
			}).PutAsArray("aa", NewMap().Put("c", "d")),
		},
		{
			Expected: bson.D{{Key: "aa", Value: bson.A{bson.M{"66": "77"}, bson.M{"c": "d"}}}},
			Actual: FromM(bson.M{
				"aa": bson.A{
					bson.M{"66": "77"},
				},
			}).PutAsArray("aa", NewMap().Put("c", "d")),
		},
		{
			Expected: bson.D{{Key: "aa", Value: bson.A{bson.M{"66": "77"}, bson.M{"c": "d"}}}},
			Actual: FromM(bson.M{
				"aa": bson.A{
					bson.D{{Key: "66", Value: "77"}},
				},
			}).PutAsArray("aa", NewMap().Put("c", "d")),
		},
		{
			Expected: bson.D{{Key: "bb", Value: bson.A{bson.M{"66": "77"}, bson.M{"c": "d"}}}},
			Actual: FromM(bson.M{
				"bb": bson.M{"66": "77"},
			}).PutAsArray("bb", NewMap().Put("c", "d")),
		},
		{
			Expected: bson.D{{Key: "bb", Value: bson.A{bson.M{"66": "77"}, bson.M{"c": "d"}}}},
			Actual: FromM(bson.M{
				"bb": bson.D{{Key: "66", Value: "77"}},
			}).PutAsArray("bb", NewMap().Put("c", "d")),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_PutAsArray_%d", i), func(t *testing.T) {
			eb, _ := bson.Marshal(testCase.Expected)
			ab, _ := bson.Marshal(testCase.Actual)
			assert.Equal(t, eb, ab)
		})
	}
}

func TestMap_Entries(t *testing.T) {
	expected := []Entry{
		{Key: "a", Value: "b"},
		{Key: "c", Value: "d"},
	}
	actual := NewMap().Put("a", "b").Put("c", "d")
	assert.Equal(t, expected, actual.Entries())
}

func TestMap_ToMap(t *testing.T) {
	testCases := []struct {
		Expected *Map
		Actual   map[string]any
	}{
		{Expected: NewMap(), Actual: make(map[string]any)},
		{Expected: NewMap().Put("a", "b"), Actual: map[string]any{"a": "b"}},
		{Expected: NewMap().Put("a", Null), Actual: map[string]any{"a": nil}},
		{
			Expected: NewMap().Put("a", bson.D{{Key: "ab", Value: "c"}}),
			Actual:   map[string]any{"a": map[string]any{"ab": "c"}},
		},
		{
			Expected: NewMap().Put("a", bson.A{1, 2}),
			Actual:   map[string]any{"a": []any{1, 2}},
		},
		{
			Expected: NewMap().Put("a", bson.D{{Key: "ab", Value: bson.M{"cc": "dd"}}}),
			Actual:   map[string]any{"a": map[string]any{"ab": map[string]any{"cc": "dd"}}},
		},
		{
			Expected: NewMap().Put("a", bson.D{{Key: "ab", Value: bson.D{{Key: "cc", Value: "dd"}}}}),
			Actual:   map[string]any{"a": map[string]any{"ab": map[string]any{"cc": "dd"}}},
		},
		{
			Expected: NewMap().Put("a", bson.D{{Key: "ab", Value: bson.D{{Key: "cc", Value: nil}}}}),
			Actual:   map[string]any{"a": map[string]any{"ab": map[string]any{"cc": nil}}},
		},
		//{
		//	Expected: NewMap().Put("a", bson.D{{Key: "ab", Value: bson.D{{Key: "cc", Value: primitive.Regex{
		//		Pattern: "^*b66y!z&",
		//		Options: "im",
		//	}}}}}),
		//	Actual: map[string]any{"a": map[string]any{"ab": map[string]any{"cc": "/^*b66y!z&/im"}}},
		//},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_ToMap_%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.Expected.ToMap(), testCase.Actual)
		})
	}
}

func TestMap_MarshalJSON(t *testing.T) {
	oid := primitive.NewObjectID()
	now := time.Now()
	m := NewMap()
	m.Put("aa", "bb")
	m.Put("bb", 1)
	m.Put("_id", oid)
	m.Put("time", now)
	m.Put("cc", Null)
	m.Put("dd", bson.A{
		bson.D{
			{Key: "sdd", Value: "ddd"},
			{Key: "sdd2", Value: "ddd1"},
		},
		bson.M{
			"aaa": nil,
		},
		NewMap().Put("sub1", 222).Put("sub2", NewMap().Put("ddv", false)),
	})
	bytes, err := m.MarshalJSON()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("map: %+v", string(bytes))
	var reversed Map
	err = json.Unmarshal(bytes, &reversed)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testCases := []struct {
		Key      string
		Expected any
		Exist    bool
	}{
		{"_id", oid, true},
		{"aa", "bb", true},
		{"time", now, true},
		{"cc", nil, true},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_JSON_%d", i), func(t *testing.T) {
			actual, ok := reversed.Get(testCase.Key)
			assert.Equal(t, testCase.Exist, ok)
			if testCase.Exist {
				eb, _ := json.Marshal(testCase.Expected)
				ab, _ := json.Marshal(actual)
				assert.Equal(t, eb, ab)
			}
		})
	}
}

func TestMap_UnmarshalBSON(t *testing.T) {
	expected := D(
		E("a", "b"),
		E("666", int32(1)),
	)
	bytes, _ := expected.MarshalBSON()
	var actual Map
	if err := bson.Unmarshal(bytes, &actual); err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.Equal(t, expected, &actual)
}
