package collection

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/executor"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/util"
	raw "go.mongodb.org/mongo-driver/mongo"
)

type TypedCollection[D bson.Doc[I], I bson.ID] struct {
	primary         *raw.Collection
	defaultReadpref *raw.Collection
}

func NewTypedCollection[D bson.Doc[I], I bson.ID](primary, defaultReadpref *raw.Collection) *TypedCollection[D, I] {
	return &TypedCollection[D, I]{
		primary:         primary,
		defaultReadpref: defaultReadpref,
	}
}

func (c *TypedCollection[D, I]) InsertOne(doc D) *executor.InsertOneExecutor[D, I] {
	return executor.NewInsertOneExecutor[D, I](c.primary, doc)
}

func (c *TypedCollection[D, I]) InsertMany(docs ...D) *executor.InsertManyExecutor[D, I] {
	return executor.NewInsertManyExecutor[D, I](c.primary, docs...)
}

func (c *TypedCollection[D, I]) FindOne(filter *filters.Filter) *executor.FindOneExecutor[D, I] {
	return executor.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindOneById(id I) *executor.FindOneExecutor[D, I] {
	return executor.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filters.Eq("_id", id))
}

func (c *TypedCollection[D, I]) Find(filter *filters.Filter) *executor.FindExecutor[D, I] {
	return executor.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindByIds(ids []I) *executor.FindExecutor[D, I] {
	return executor.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filters.In("_id", util.ToAny(ids)))
}

func (c *TypedCollection[D, I]) CountDocuments(filter *filters.Filter) *executor.CountExecutor[D, I] {
	return executor.NewCountExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindOneAndUpdate(filter *filters.Filter, update *updates.Update) *executor.FindOneAndUpdateExecutor[D, I] {
	return executor.NewFindOneAndUpdateExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) UpdateOne(filter *filters.Filter, update *updates.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateOneExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) UpdateMany(filter *filters.Filter, update *updates.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateManyExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) UpdateById(id I, update *updates.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateByIdExecutor[D, I](c.primary, id, update)
}

func (c *TypedCollection[D, I]) DeleteOne(filter *filters.Filter) *executor.DeleteExecutor[D, I] {
	return executor.NewDeleteOneExecutor[D, I](c.primary, filter)
}

func (c *TypedCollection[D, I]) DeleteMany(filter *filters.Filter) *executor.DeleteExecutor[D, I] {
	return executor.NewDeleteManyExecutor[D, I](c.primary, filter)
}

func (c *TypedCollection[D, I]) BulkWrite() *executor.TypedBulkWriteExecutor[D, I] {
	return executor.NewBulkWriteExecutor[D, I](c.primary)
}

func (c *TypedCollection[D, I]) Aggregate(pipe raw.Pipeline) *executor.AggregateExecutor[D, I] {
	return executor.NewAggregateExecutor[D, I](c.primary, c.defaultReadpref, pipe)
}

func (c *TypedCollection[D, I]) Distinct(field string) *executor.DistinctExecutor {
	return executor.NewDistinctExecutor(c.primary, c.defaultReadpref, field)
}

func (c *TypedCollection[D, I]) Indexes() *IndexViewer {
	return FromIndexView(c.primary.Indexes())
}
