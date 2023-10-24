package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestFindOneAndUpdateExecutor_Execute(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test2", Age: 10},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test upsert", func(t *testing.T) {
		doc, err := NewFindOneAndUpdateExecutor[*TestDoc, primitive.ObjectID](
			testColl, filters.Eq("name", "test3"), updates.Set("age", 18),
		).Upsert().ReturnAfter().Execute(context.TODO())
		if err != nil {
			t.Errorf("upsert doc error. %+v", err)
			t.FailNow()
		}
		assert.Equal(t, "test3", doc.Name)
		assert.Equal(t, 18, doc.Age)
	})
	t.Run("test update only exist", func(t *testing.T) {
		_, err := NewFindOneAndUpdateExecutor[*TestDoc, primitive.ObjectID](
			testColl, filters.Eq("name", "test4"), updates.Set("age", 18),
		).ReturnAfter().Execute(context.TODO())
		assert.Equal(t, mongo.ErrNoDocuments, err)
		doc, err := NewFindOneAndUpdateExecutor[*TestDoc, primitive.ObjectID](
			testColl, filters.Eq("name", "test2"), updates.Set("age", 22),
		).ReturnAfter().Execute(context.TODO())
		if err != nil {
			t.Errorf("upsert doc error. %+v", err)
			t.FailNow()
		}
		assert.Equal(t, 22, doc.Age)
	})
}

func TestFindOneAndUpdateExecutor_Collect(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test2", Age: 10},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test upsert", func(t *testing.T) {
		var doc testDocProj
		err := NewFindOneAndUpdateExecutor[*TestDoc, primitive.ObjectID](
			testColl, filters.Eq("name", "test3"), updates.Set("age", 18),
		).Upsert().ReturnAfter().Collect(context.TODO(), &doc)
		if err != nil {
			t.Errorf("upsert doc error. %+v", err)
			t.FailNow()
		}
		assert.Equal(t, "test3", doc.Name)
	})
}
