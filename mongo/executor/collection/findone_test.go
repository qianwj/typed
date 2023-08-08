package collection

import (
	"context"
	"fmt"
	"github.com/qianwj/typed/mongo/builder"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/options"
	"github.com/qianwj/typed/mongo/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestCollection_FindOne(t *testing.T) {
	createTime := time.UnixMilli(time.Now().UTC().UnixMilli()).UTC()
	docs := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "Amy", Age: 18, CreateTime: createTime},
		{Id: primitive.NewObjectID(), Name: "Lily", Age: 19, CreateTime: createTime.Add(time.Minute)},
		{Id: primitive.NewObjectID(), Name: "Mike", Age: 21, CreateTime: createTime.Add(10 * time.Minute)},
		{Id: primitive.NewObjectID(), Name: "LiLei", Age: 30, CreateTime: createTime.Add(time.Hour)},
	}
	_, err := testColl.InsertMany(context.TODO(), util.ToAny(docs))
	if err != nil {
		t.FailNow()
		return
	}
	coll := NewCollection[*TestDoc, primitive.ObjectID](testColl, testColl)
	tests := []struct {
		Expect   *TestDoc
		Err      error
		Executor *builder.FindOneExecutor[*TestDoc, primitive.ObjectID]
	}{
		{Expect: docs[0], Executor: coll.FindOne(filter.Eq("name", "Amy"))},
		{Err: mongo.ErrNoDocuments, Executor: coll.FindOne(filter.Gt("createTime", createTime.Add(time.Hour)))},
		{Expect: docs[3], Executor: coll.FindOne(filter.Gte("createTime", createTime.Add(time.Hour)))},
		{Expect: docs[0], Executor: coll.FindOne(filter.Nin("age", []any{19, 21})).Sort(options.Ascending("age"))},
		{Expect: docs[3], Executor: coll.FindOne(filter.Nin("age", []any{19, 21})).Sort(options.Ascending("age")).Skip(1)},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			actual, err := tc.Executor.Execute(context.TODO())
			if tc.Err != nil {
				assert.Equal(t, tc.Err, err)
			} else {
				assert.Equal(t, tc.Expect, actual)
			}
		})
	}
	opts := mtest.NewOptions().
		ClientType(mtest.Mock).
		CreateClient(true).
		DatabaseName(testDBName).
		CollectionName(testCollName)
	dataMock := mtest.New(t, opts)
	defer dataMock.Close()
	dataMock.Run("mock", func(mt *mtest.T) {
		mockError := mtest.CommandError{
			Message: mongo.ErrNoDocuments.Error(),
		}
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mockError))
		coll := NewCollection[model.BsonMap[string], string](mt.Coll, mt.Coll)
		_, err := coll.FindOne(filter.Eq("a", "b")).Execute(context.TODO())
		if !assert.Equal(mt, mockError.Message, err.Error()) {
			mt.FailNow()
		}
	})
}
