package executor

import (
	"github.com/qianwj/typed/mongo"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/options"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[D model.Document[I], I model.DocumentId] struct {
	primary   *raw.Collection
	secondary *raw.Collection
	//FindOne(ctx context.Context, filter *filter.Filter, opts ...*options.FindOneOptions) (D, error)
	//Find(ctx context.Context, filter model.Filter, opts ...*options.FindOptions) ([]D, error)
	//FindByDocIds(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) ([]D, error)
	//CountDocuments(ctx context.Context, filter model.Filter, opts ...*options.CountOptions) (int64, error)
	//InsertOne(ctx context.Context, doc D) (primitive.ObjectID, error)
	//InsertMany(ctx context.Context, docs []D) ([]primitive.ObjectID, error)
	//FindOneAndUpdate(ctx context.Context, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (D, error)
	//UpdateOne(ctx context.Context, m *model.UpdateOne, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//UpdateMany(ctx context.Context, m *model.UpdateMany, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//UpdateById(ctx context.Context, m *model.UpdateById, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//DeleteOne(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//DeleteMany(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//InitializeBulkWriteOp(opts ...*options.BulkWriteOptions) *BulkWriteOperation
}

func FromDatabase[D model.Document[I], I model.DocumentId](db *mongo.Database, name string, opts ...*rawopts.CollectionOptions) *Collection[D, I] {
	return &Collection[D, I]{
		primary:   db.Primary(name, opts...),
		secondary: db.Secondary(name, opts...),
	}
}

//func (c *Collection[D, I]) FindOne(ctx context.Context, filter *filter.Filter) (D, error) {
//	singleResult := c.internal.FindOne(ctx, filter.Marshal(), options.MergeFindOneOptions(opts...))
//	var doc D
//	if singleResult.Err() != nil {
//		return doc, singleResult.Err()
//	}
//	if err := singleResult.Decode(doc); err != nil {
//		return doc, err
//	}
//	return doc, nil
//}

func (c *Collection[D, I]) Find(filter *filter.Filter) *FindExecutor[D, I] {
	return &FindExecutor[D, I]{
		coll:   c,
		filter: filter,
		opts:   options.Find(),
	}
}

//
//func (c typedCollectionImpl[D]) FindByDocIds(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) ([]D, error) {
//	filter := model.NewFilter().In("_id", lo.ToAnySlice[primitive.ObjectID](ids))
//	return c.Find(ctx, filter, opts...)
//}
//
//func (c typedCollectionImpl[D]) CountDocuments(ctx context.Context, filter model.Filter, opts ...*options.CountOptions) (int64, error) {
//	return c.internal.CountDocuments(ctx, filter, options.MergeCountOptions(opts...))
//}
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
//func (c typedCollectionImpl[D]) UpdateOne(ctx context.Context, m *model.UpdateOne, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	if len(m.Update) == 0 {
//		return nil, errors.New("update is empty")
//	}
//	return c.internal.UpdateOne(ctx, m.Filter, m.Update, options.MergeUpdateOptions(opts...))
//}
//
//func (c typedCollectionImpl[D]) UpdateMany(ctx context.Context, m *model.UpdateMany, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	if len(m.Update) == 0 {
//		return nil, errors.New("update is empty")
//	}
//	res, err := c.internal.UpdateMany(ctx, m.Filter, m.Update, options.MergeUpdateOptions(opts...))
//	if err != nil {
//		return nil, err
//	}
//	return res, nil
//}
//
//func (c typedCollectionImpl[D]) UpdateById(ctx context.Context, m *model.UpdateById, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	if len(m.Update) == 0 {
//		return nil, errors.New("update is empty")
//	}
//	res, err := c.internal.UpdateByID(ctx, m.Id, m.Update, options.MergeUpdateOptions(opts...))
//	if err != nil {
//		return nil, err
//	}
//	return res, nil
//}
//
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
