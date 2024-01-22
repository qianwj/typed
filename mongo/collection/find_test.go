package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/util"
	"github.com/qianwj/typed/streams"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindExecutor_ToArray(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test_to_array1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_to_array2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_to_array3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_to_array3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test_to_array3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test_to_array4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_to_array1"))
		actual, err := exec.ToArray(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*TestDoc{prepared[0]}
		assert.Equal(t, expected, actual)
	})
	t.Run("test more", func(t *testing.T) {
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_to_array3"))
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
		{Id: primitive.NewObjectID(), Name: "test_collect1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_collect2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_collect3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_collect3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test_collect3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test_collect4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_collect1")).
			Projection(projections.Includes("name").ExcludeId())
		var actual []*testDocProj
		err := exec.Collect(context.TODO(), &actual)
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*testDocProj{{Name: "test_collect1"}}
		assert.Equal(t, expected, actual)
	})
	t.Run("test more", func(t *testing.T) {
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_collect3")).
			Projection(projections.Includes("name").ExcludeId())
		var actual []*testDocProj
		err := exec.Collect(context.TODO(), &actual)
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		expected := []*testDocProj{{Name: "test_collect3"}, {Name: "test_collect3"}, {Name: "test_collect3"}}
		assert.Equal(t, expected, actual)
	})
}

func TestFindExecutor_Cursor(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test_cursor1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test_cursor4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))
	t.Run("test one", func(t *testing.T) {
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_cursor1"))
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
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_cursor3"))
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

func TestFindExecutor_Stream(t *testing.T) {
	prepared := []*TestDoc{
		{Id: primitive.NewObjectID(), Name: "test_cursor1", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor2", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 10},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 20},
		{Id: primitive.NewObjectID(), Name: "test_cursor3", Age: 30},
		{Id: primitive.NewObjectID(), Name: "test_cursor4", Age: 40},
	}
	_, _ = testColl.InsertMany(context.TODO(), util.ToAny(prepared))

	t.Run("test one", func(t *testing.T) {
		sub := streams.NewFixedSubscriber[*TestDoc](
			streams.DoOnNext[*TestDoc](func(doc *TestDoc) {
				t.Logf("on next: %+v", doc)
			}),
			streams.DoOnError[*TestDoc](func(err error) {
				t.Errorf("on error: %+v", err)
			}),
			streams.DoOnComplete[*TestDoc](func() {
				t.Logf("complete.")
			}),
		)
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_cursor1"))
		publisher, err := exec.Stream(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		publisher.Subscribe(sub)
	})
	t.Run("test more", func(t *testing.T) {
		sub := &testSubscriber{t: t}
		exec := newFindExecutor[*TestDoc, primitive.ObjectID](testColl, testColl, filters.Eq("name", "test_cursor3"))
		publisher, err := exec.Stream(context.TODO())
		if err != nil {
			t.Errorf("find error: %+v", err)
			t.FailNow()
		}
		publisher.Subscribe(sub)
	})
}

type testSubscriber struct {
	streams.Subscriber[*TestDoc]
	t            *testing.T
	complete     bool
	subscription streams.Subscription[*TestDoc]
}

func (s *testSubscriber) OnSubscribe(subscription streams.Subscription[*TestDoc]) {
	s.subscription = subscription
	for !s.complete {
		subscription.Request(1)
	}
}

func (s *testSubscriber) OnNext(val *TestDoc) {
	s.t.Logf("on next: %+v", val)
}

func (s *testSubscriber) OnError(err error) {
	s.t.Errorf("on error: %+v", err)
	if s.subscription != nil {
		s.subscription.Cancel()
	}
}

func (s *testSubscriber) OnComplete() {
	s.complete = true
}
