package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/util"
	raw "go.mongodb.org/mongo-driver/mongo"
)

type TypedCollection[D model.Doc[I], I model.ID] struct {
	primary         *raw.Collection
	defaultReadpref *raw.Collection
}

func NewTypedCollection[D model.Doc[I], I model.ID](primary, defaultReadpref *raw.Collection) *TypedCollection[D, I] {
	return &TypedCollection[D, I]{
		primary:         primary,
		defaultReadpref: defaultReadpref,
	}
}

func (c *TypedCollection[D, I]) InsertOne(doc D) *InsertOneExecutor[D, I] {
	return newInsertOneExecutor[D, I](c.primary, doc)
}

func (c *TypedCollection[D, I]) InsertMany(docs ...D) *InsertManyExecutor[D, I] {
	return newInsertManyExecutor[D, I](c.primary, docs...)
}

func (c *TypedCollection[D, I]) FindOne(filter *filters.Filter) *FindOneExecutor[D, I] {
	return newFindOneExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindOneById(id I) *FindOneExecutor[D, I] {
	return newFindOneExecutor[D, I](c.primary, c.defaultReadpref, filters.Eq("_id", id))
}

func (c *TypedCollection[D, I]) Find(filter *filters.Filter) *FindExecutor[D, I] {
	return newFindExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindByIds(ids []I) *FindExecutor[D, I] {
	return newFindExecutor[D, I](c.primary, c.defaultReadpref, filters.In("_id", util.ToAny(ids)))
}

func (c *TypedCollection[D, I]) CountDocuments(filter *filters.Filter) *CountExecutor[D, I] {
	return newCountExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *TypedCollection[D, I]) FindOneAndUpdate(filter *filters.Filter, update *updates.Update) *FindOneAndUpdateExecutor[D, I] {
	return newFindOneAndUpdateExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) FindOneAndReplace(filter *filters.Filter, replacement any) *FindOneAndReplaceExecutor[D, I] {
	return newFindOneAndReplaceExecutor[D, I](c.primary, filter, replacement)
}

func (c *TypedCollection[D, I]) FindOneAndDelete(filter *filters.Filter) *FindOneAndDeleteExecutor[D, I] {
	return newFindOneAndDeleteExecutor[D, I](c.primary, filter)
}

func (c *TypedCollection[D, I]) UpdateOne(filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return newUpdateOneExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) UpdateMany(filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return newUpdateManyExecutor[D, I](c.primary, filter, update)
}

func (c *TypedCollection[D, I]) UpdateById(id I, update *updates.Update) *UpdateExecutor[D, I] {
	return newUpdateByIdExecutor[D, I](c.primary, id, update)
}

func (c *TypedCollection[D, I]) Delete(filter *filters.Filter) *DeleteExecutor {
	return newDeleteExecutor(c.primary, filter)
}

func (c *TypedCollection[D, I]) BulkWrite() *TypedBulkWriteExecutor[D, I] {
	return newBulkWriteExecutor[D, I](c.primary)
}

func (c *TypedCollection[D, I]) Aggregate(pipe *aggregates.Pipeline) *AggregateExecutor[D, I] {
	return newAggregateExecutor[D, I](c.primary, c.defaultReadpref, pipe)
}

func (c *TypedCollection[D, I]) Distinct(field string) *DistinctExecutor {
	return newDistinctExecutor(c.primary, c.defaultReadpref, field)
}

func (c *TypedCollection[D, I]) Indexes() *IndexViewer {
	return fromIndexView(c.primary.Indexes())
}

func (c *TypedCollection[D, I]) Drop(ctx context.Context) error {
	return c.primary.Drop(ctx)
}
