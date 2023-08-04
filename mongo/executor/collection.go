package executor

import (
	"context"
	"github.com/qianwj/typed/mongo"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	raw "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[D model.Document[I], I model.DocumentId] struct {
	primary   *raw.Collection
	secondary *raw.Collection
	//FindOneAndUpdate(ctx context.Context, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (D, error)
	//DeleteOne(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//DeleteMany(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//InitializeBulkWriteOp(opts ...*options.BulkWriteOptions) *BulkWriteOperation
}

func FromDatabase[D model.Document[I], I model.DocumentId](db *mongo.Database, name string, opts ...*options.CollectionOptions) *Collection[D, I] {
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

//
//func (c typedCollectionImpl[D]) InsertOne(ctx context.Context, doc D) (primitive.ObjectID, error) {
//	res, err := c.internal.InsertOne(ctx, doc)
//	if err != nil {
//		return primitive.NilObjectID, err
//	}
//	return res.InsertedID.(primitive.ObjectID), nil
//}
//
//func (c typedCollectionImpl[D]) InsertMany(ctx context.Context, docs []D) ([]primitive.ObjectID, error) {
//	res, err := c.internal.InsertMany(ctx, lo.ToAnySlice[D](docs))
//	if err != nil {
//		return nil, err
//	}
//	data, _ := lo.FromAnySlice[primitive.ObjectID](res.InsertedIDs)
//	return data, nil
//}
//
//func (c typedCollectionImpl[D]) FindOneAndUpdate(ctx context.Context, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (D, error) {
//	var t D
//	if len(m.Update) == 0 {
//		return t, errors.New("update is empty")
//	}
//	singleResult := c.internal.FindOneAndUpdate(ctx, m.Filter, m.Update, options.MergeFindOneAndUpdateOptions(opts...))
//	if singleResult.Err() != nil {
//		return t, singleResult.Err()
//	}
//	if err := singleResult.Decode(t); err != nil {
//		return t, err
//	}
//	return t, nil
//}
//

func (c *Collection[D, I]) UpdateOne(filter *filter.Filter) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Update(),
	}
}

func (c *Collection[D, I]) UpdateMany(filter *filter.Filter) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   c,
		filter: filter,
		multi:  true,
		opts:   options.Update(),
	}
}

func (c *Collection[D, I]) UpdateById(id I) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:  c,
		docId: &id,
		opts:  options.Update(),
	}
}

//func (c typedCollectionImpl[D]) DeleteOne(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	return c.internal.DeleteOne(ctx, filter, options.MergeDeleteOptions(opts...))
//}
//
//func (c typedCollectionImpl[D]) DeleteMany(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	return c.internal.DeleteMany(ctx, filter, options.MergeDeleteOptions(opts...))
//}
//
//func (c typedCollectionImpl[D]) InitializeBulkWriteOp(opts ...*options.BulkWriteOptions) *BulkWriteOperation {
//	return newBulkWriteOperation(c.internal, opts...)
//}

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
