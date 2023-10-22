package executor

import (
	"testing"
)

func TestCollection_FindOne(t *testing.T) {
	//createTime := time.UnixMilli(time.Now().UTC().UnixMilli()).UTC()
	//docs := []*TestDoc{
	//	{Id: primitive.NewObjectID(), Name: "Amy", Age: 18, CreateTime: createTime},
	//	{Id: primitive.NewObjectID(), Name: "Lily", Age: 19, CreateTime: createTime.Add(time.Minute)},
	//	{Id: primitive.NewObjectID(), Name: "Mike", Age: 21, CreateTime: createTime.Add(10 * time.Minute)},
	//	{Id: primitive.NewObjectID(), Name: "LiLei", Age: 30, CreateTime: createTime.Add(time.Hour)},
	//}
	//_, err := testColl.InsertMany(context.TODO(), util.ToAny(docs))
	//if err != nil {
	//	t.FailNow()
	//	return
	//}
	//tests := []struct {
	//	Expect   *TestDoc
	//	Err      error
	//	Executor *FindOneExecutor[*TestDoc, primitive.ObjectID]
	//}{
	//	{Expect: docs[0], Executor: NewFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "Amy"))},
	//	{Err: mongo.ErrNoDocuments, Executor: NewFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Gt("createTime", createTime.Add(time.Hour)))},
	//	{Expect: docs[3], Executor: NewFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Gte("createTime", createTime.Add(time.Hour)))},
	//	{Expect: docs[0], Executor: NewFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Nin("age", []any{19, 21})).Sort(sorts.Ascending("age"))},
	//	{Expect: docs[3], Executor: NewFindOneExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Nin("age", []any{19, 21})).Sort(sorts.Ascending("age")).Skip(1)},
	//}
	//for i, tc := range tests {
	//	t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
	//		actual, err := tc.Executor.Execute(context.TODO())
	//		if tc.Err != nil {
	//			assert.Equal(t, tc.Err, err)
	//		} else {
	//			assert.Equal(t, tc.Expect, actual)
	//		}
	//	})
	//}
	//opts := mtest.NewOptions().
	//	ClientType(mtest.Mock).
	//	CreateClient(true).
	//	DatabaseName(testDBName).
	//	CollectionName(testCollName)
	//dataMock := mtest.New(t, opts)
	//defer dataMock.Close()
	//dataMock.Run("mock", func(mt *mtest.T) {
	//	mockError := mtest.CommandError{
	//		Message: mongo.ErrNoDocuments.Error(),
	//	}
	//	mt.AddMockResponses(mtest.CreateCommandErrorResponse(mockError))
	//	_, err := NewFindOneExecutor[model.BsonMap[string], string](mt.Coll, mt.Coll, filters.Eq("a", "b")).Execute(context.TODO())
	//	if !assert.Equal(mt, mockError.Message, err.Error()) {
	//		mt.FailNow()
	//	}
	//})
}
