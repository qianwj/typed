package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindExecutor_ToArray(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test1"))
		actual, err := exec.ToArray(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*TestDoc{prepared[0]}
		assert.Equal(t, expected, actual)
	})
	t.Run("test more", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test3"))
		actual, err := exec.ToArray(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*TestDoc{prepared[2], prepared[3], prepared[4]}
		assert.Equal(t, expected, actual)
	})
}

func TestFindExecutor_Collect(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test1")).
			Projection(projections.Includes("name").ExcludeId())
		var actual []*testDocProj
		err := exec.Collect(context.TODO(), &actual)
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*testDocProj{{Name: "test1"}}
		assert.Equal(t, expected, actual)
	})
	t.Run("test more", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test3")).
			Projection(projections.Includes("name").ExcludeId())
		var actual []*testDocProj
		err := exec.Collect(context.TODO(), &actual)
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*testDocProj{{Name: "test3"}, {Name: "test3"}, {Name: "test3"}}
		assert.Equal(t, expected, actual)
	})
}

func TestFindExecutor_Cursor(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test1"))
		cursor, err := exec.Cursor(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		if cursor.HasNext(context.TODO()) {
			actual, err := cursor.Next()
			if err != nil {
				t.Errorf("deserialized doc error. %+v", err)
				t.FailNow()
			}
			assert.Equal(t, prepared[0], actual)
		} else {
			t.Errorf("expect at least one element, but not found.")
			t.FailNow()
		}
	})
	t.Run("test more", func(t *testing.T) {
		exec := NewFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test3"))
		cursor, err := exec.Cursor(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		i, expected := 0, []*TestDoc{prepared[2], prepared[3], prepared[4]}
		for cursor.HasNext(context.TODO()) || i < 3 {
			actual, err := cursor.Next()
			if err != nil {
				t.Errorf("deserialized doc error. %+v", err)
				t.FailNow()
			}
			assert.Equal(t, expected[i], actual)
			i++
		}
	})
}
