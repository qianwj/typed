package executor

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregate/pipe"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/update"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *Collection[D, I]) InsertOne(doc D) *InsertOneExecutor[D, I] {
	return &InsertOneExecutor[D, I]{
		coll: c,
		data: doc,
		opts: options.InsertOne(),
	}
}

func (c *Collection[D, I]) InsertMany(docs ...D) *InsertManyExecutor[D, I] {
	return &InsertManyExecutor[D, I]{
		coll: c,
		data: toAny(docs),
		opts: options.InsertMany(),
	}
}

func (c *Collection[D, I]) FindOne(filter *filter.Filter) *FindOneExecutor[D, I] {
	return &FindOneExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.FindOne(),
	}
}

func (c *Collection[D, I]) FindOneById(id I, primary bool) *FindOneExecutor[D, I] {
	return &FindOneExecutor[D, I]{
		coll:    c,
		filter:  filter.Eq("_id", id),
		primary: primary,
		opts:    options.FindOne(),
	}
}

func (c *Collection[D, I]) Find(filter *filter.Filter) *FindExecutor[D, I] {
	return &FindExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Find(),
	}
}

func (c *Collection[D, I]) FindByIds(ids []I) *FindExecutor[D, I] {
	return &FindExecutor[D, I]{
		coll:   c,
		filter: filter.In("_id", toAny(ids)),
		opts:   options.Find(),
	}
}

func (c *Collection[D, I]) CountDocuments(filter *filter.Filter) *CountExecutor[D, I] {
	return &CountExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Count(),
	}
}

func (c *Collection[D, I]) FindOneAndUpdate(filter *filter.Filter, update *update.Update) *FindOneAndUpdateExecutor[D, I] {
	return &FindOneAndUpdateExecutor[D, I]{
		coll:   c,
		filter: filter,
		update: update,
		opts:   options.FindOneAndUpdate(),
	}
}

func (c *Collection[D, I]) UpdateOne(filter *filter.Filter, update *update.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   c,
		filter: filter,
		update: update,
		opts:   options.Update(),
	}
}

func (c *Collection[D, I]) UpdateMany(filter *filter.Filter, update *update.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   c,
		filter: filter,
		update: update,
		multi:  true,
		opts:   options.Update(),
	}
}

func (c *Collection[D, I]) UpdateById(id I, update *update.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   c,
		docId:  &id,
		update: update,
		opts:   options.Update(),
	}
}

func (c *Collection[D, I]) DeleteOne(filter *filter.Filter) *DeleteExecutor[D, I] {
	return &DeleteExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Delete(),
	}
}

func (c *Collection[D, I]) DeleteMany(filter *filter.Filter) *DeleteExecutor[D, I] {
	return &DeleteExecutor[D, I]{
		coll:   c,
		filter: filter,
		multi:  true,
		opts:   options.Delete(),
	}
}

func (c *Collection[D, I]) BulkWrite() *BulkWriteExecutor[D, I] {
	return &BulkWriteExecutor[D, I]{
		coll: c,
		opts: options.BulkWrite(),
	}
}

func (c *Collection[D, I]) Aggregate(pipe *pipe.Pipeline) *AggregateExecutor[D, I] {
	return &AggregateExecutor[D, I]{
		coll: c,
		pipe: pipe,
		opts: options.Aggregate(),
	}
}

func (c *Collection[D, I]) Indexes() *IndexViewer {
	return &IndexViewer{view: c.primary.Indexes()}
}
