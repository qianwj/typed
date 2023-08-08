package collection

import (
	"github.com/qianwj/typed/mongo/builder"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregate"
	"github.com/qianwj/typed/mongo/model/filter"
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

func (c *Collection[D, I]) InsertOne(doc D) *builder.InsertOneExecutor[D, I] {
	return builder.NewInsertOneExecutor[D, I](c.primary, doc)
}

func (c *Collection[D, I]) InsertMany(docs ...D) *builder.InsertManyExecutor[D, I] {
	return builder.NewInsertManyExecutor[D, I](c.primary, docs...)
}

func (c *Collection[D, I]) FindOne(filter *filter.Filter) *builder.FindOneExecutor[D, I] {
	return builder.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindOneById(id I) *builder.FindOneExecutor[D, I] {
	return builder.NewFindOneExecutor[D, I](c.primary, c.defaultReadpref, filter.Eq("_id", id))
}

func (c *Collection[D, I]) Find(filter *filter.Filter) *builder.FindExecutor[D, I] {
	return builder.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindByIds(ids []I) *builder.FindExecutor[D, I] {
	return builder.NewFindExecutor[D, I](c.primary, c.defaultReadpref, filter.In("_id", util.ToAny(ids)))
}

func (c *Collection[D, I]) CountDocuments(filter *filter.Filter) *builder.CountExecutor[D, I] {
	return builder.NewCountExecutor[D, I](c.primary, c.defaultReadpref, filter)
}

func (c *Collection[D, I]) FindOneAndUpdate(filter *filter.Filter, update *update.Update) *builder.FindOneAndUpdateExecutor[D, I] {
	return builder.NewFindOneAndUpdateExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateOne(filter *filter.Filter, update *update.Update) *builder.UpdateExecutor[D, I] {
	return builder.NewUpdateOneExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateMany(filter *filter.Filter, update *update.Update) *builder.UpdateExecutor[D, I] {
	return builder.NewUpdateManyExecutor[D, I](c.primary, filter, update)
}

func (c *Collection[D, I]) UpdateById(id I, update *update.Update) *builder.UpdateExecutor[D, I] {
	return builder.NewUpdateByIdExecutor[D, I](c.primary, id, update)
}

func (c *Collection[D, I]) DeleteOne(filter *filter.Filter) *builder.DeleteExecutor[D, I] {
	return builder.NewDeleteOneExecutor[D, I](c.primary, filter)
}

func (c *Collection[D, I]) DeleteMany(filter *filter.Filter) *builder.DeleteExecutor[D, I] {
	return builder.NewDeleteManyExecutor[D, I](c.primary, filter)
}

func (c *Collection[D, I]) BulkWrite() *builder.BulkWriteExecutor[D, I] {
	return builder.NewBulkWriteExecutor[D, I](c.primary)
}

func (c *Collection[D, I]) Aggregate(pipe aggregate.Pipeline) *builder.AggregateExecutor[D, I] {
	return builder.NewAggregateExecutor[D, I](c.primary, c.defaultReadpref, pipe)
}

func (c *Collection[D, I]) Indexes() *IndexViewer {
	return FromIndexView(c.primary.Indexes())
}
