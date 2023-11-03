package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindOneExecutor_Execute(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test_one_exec1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_one_exec2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_one_exec3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_one_exec3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test_one_exec3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test_one_exec4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := newFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_one_exec1"))
		actual, err := exec.Execute(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		assert.Equal(t, prepared[0], actual)
	})
	t.Run("test multi", func(t *testing.T) {
		exec := newFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_one_exec3").Eq("age", 30))
		actual, err := exec.Execute(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		assert.Equal(t, prepared[4], actual)
	})
	t.Run("test by id", func(t *testing.T) {
		exec := newFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("_id", prepared[3].Id))
		actual, err := exec.Execute(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		assert.Equal(t, prepared[3], actual)
	})
}

func TestFindOneExecutor_Collect(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test_one_collect1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_one_collect2", Age: 10},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := newFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_one_collect1")).
			Projection(projections.Includes("name").ExcludeId())
		var actual testDocProj
		err := exec.Collect(context.TODO(), &actual)
		if err != nil {
			t.Errorf("find one error: %+v", err)
			t.FailNow()
		}
		expected := &testDocProj{Name: "test_one_collect1"}
		assert.Equal(t, expected, &actual)
	})
}
