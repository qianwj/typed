package collection

import (
	"github.com/qianwj/typed/mongo/executor"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregate"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/update"
	"github.com/qianwj/typed/mongo/util"
	raw "go.mongodb.org/mongo-driver/mongo"
)

type Collection[D model.Document[I], I model.DocumentId] struct {
	primary         *raw.Collection
	defaultReadpref *raw.Collection
}

func NewCollection[D model.Document[I], I model.DocumentId](primary, defaultReadpref *raw.Collection) *Collection[D, I] {
	return &Collection[D, I]{
		primary:         primary,
		defaultReadpref: defaultReadpref,
	}
}

func (c *Collection[D, I]) InsertOne(doc D) *executor.InsertOneExecutor[D, I] {
	return executor.NewInsertOneExecutor[D, I](c.primary, doc)
}

func (c *Collection[D, I]) InsertMany(docs ...D) *executor.InsertManyExecutor[D, I] {
	return executor.NewInsertManyExecutor[D, I](c.primary, docs...)
}

func (c *Collection[D, I]) FindOne(filter *filters.Filter) *executor.FindOneExecutor[D, I] {
	return executor.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindOneById(id I) *executor.FindOneExecutor[D, I] {
	return executor.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filters.Eq("_id", id))
}

func (c *Collection[D, I]) Find(filter *filters.Filter) *executor.FindExecutor[D, I] {
	return executor.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindByIds(ids []I) *executor.FindExecutor[D, I] {
	return executor.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filters.In("_id", util.ToAny(ids)))
}

func (c *Collection[D, I]) CountDocuments(filter *filters.Filter) *executor.CountExecutor[D, I] {
	return executor.NewCountExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindOneAndUpdate(filter *filters.Filter, update *update.Update) *executor.FindOneAndUpdateExecutor[D, I] {
	return executor.NewFindOneAndUpdateExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateOne(filter *filters.Filter, update *update.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateOneExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateMany(filter *filters.Filter, update *update.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateManyExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateById(id I, update *update.Update) *executor.UpdateExecutor[D, I] {
	return executor.NewUpdateByIdExecutor[D, I](c.primary, id, update)
}

func (c *Collection[D, I]) DeleteOne(filter *filters.Filter) *executor.DeleteExecutor[D, I] {
	return executor.NewDeleteOneExecutor[D, I](c.primary, filter)
}

func (c *Collection[D, I]) DeleteMany(filter *filters.Filter) *executor.DeleteExecutor[D, I] {
	return executor.NewDeleteManyExecutor[D, I](c.primary, filter)
}

func (c *Collection[D, I]) BulkWrite() *executor.BulkWriteExecutor[D, I] {
	return executor.NewBulkWriteExecutor[D, I](c.primary)
}

func (c *Collection[D, I]) Aggregate(pipe aggregate.Pipeline) *executor.AggregateExecutor[D, I] {
	return executor.NewAggregateExecutor[D, I](c.primary, c.defaultReadpref, pipe)
}

func (c *Collection[D, I]) Indexes() *IndexViewer {
	return FromIndexView(c.primary.Indexes())
}
