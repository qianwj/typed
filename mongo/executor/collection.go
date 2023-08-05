package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/update"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[D model.Document[I], I model.DocumentId] struct {
	primary   *raw.Collection
	secondary *raw.Collection
}

func FromDatabase[D model.Document[I], I model.DocumentId](db *Database, name string, opts ...*options.CollectionOptions) *Collection[D, I] {
	return &Collection[D, I]{
		primary:   db.Primary(name, opts...),
		secondary: db.Secondary(name, opts...),
	}
}

func (c *Collection[D, I]) InsertOne(ctx context.Context, doc *D) (I, error) {
	res, err := c.primary.InsertOne(ctx, doc)
	var id I
	if err != nil {
		return id, err
	}
	return res.InsertedID.(I), nil
}

func (c *Collection[D, I]) InsertMany(ctx context.Context, docs []*D) ([]I, error) {
	res, err := c.primary.InsertMany(ctx, toAny(docs))
	if err != nil {
		return nil, err
	}

	return mapTo(res.InsertedIDs, func(i any) I {
		return i.(I)
	}), nil
}

func (c *Collection[D, I]) FindOne(filter *filter.Filter) *FindOneExecutor[D, I] {
	return &FindOneExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.FindOne(),
	}
}

func (c *Collection[D, I]) Find(filter *filter.Filter) *FindExecutor[D, I] {
	return &FindExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Find(),
	}
}

func (c *Collection[D, I]) FindByDocIds(ids []I) *FindExecutor[D, I] {
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

func toAny[T any](arr []T) []any {
	res := make([]any, len(arr))
	for i, t := range arr {
		res[i] = t
	}
	return res
}

func mapTo[T, U any](arr []T, convert func(t T) U) []U {
	res := make([]U, len(arr))
	for i, t := range arr {
		res[i] = convert(t)
	}
	return res
}