package mongo

import (
	"context"
	"errors"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/options"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TypedCollection[D model.Document] interface {
	collection() *mongo.Collection
	FindOne(ctx context.Context, filter model.Filter, opts ...*options.FindOneOptions) (D, error)
	Find(ctx context.Context, filter model.Filter, opts ...*options.FindOptions) ([]D, error)
	FindByDocIds(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) ([]D, error)
	CountDocuments(ctx context.Context, filter model.Filter, opts ...*options.CountOptions) (int64, error)
	InsertOne(ctx context.Context, doc D) (primitive.ObjectID, error)
	InsertMany(ctx context.Context, docs []D) ([]primitive.ObjectID, error)
	FindOneAndUpdate(ctx context.Context, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (D, error)
	UpdateOne(ctx context.Context, m *model.UpdateOne, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, m *model.UpdateMany, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateById(ctx context.Context, m *model.UpdateById, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	InitializeBulkWriteOp(opts ...*options.BulkWriteOptions) *BulkWriteOperation
}

type typedCollectionImpl[D model.Document] struct {
	TypedCollection[D]
	internal *mongo.Collection
}

func NewTypedCollection[D model.Document](db *mongo.Database, name string, opts ...*options.CollectionOptions) TypedCollection[D] {
	return &typedCollectionImpl[D]{internal: db.Collection(name, options.MergeCollectionOptions(opts...))}
}

func (c typedCollectionImpl[D]) collection() *mongo.Collection {
	return c.internal
}

func (c typedCollectionImpl[D]) FindOne(ctx context.Context, filter model.Filter, opts ...*options.FindOneOptions) (D, error) {
	singleResult := c.internal.FindOne(ctx, filter, options.MergeFindOneOptions(opts...))
	var doc D
	if singleResult.Err() != nil {
		return doc, singleResult.Err()
	}
	if err := singleResult.Decode(doc); err != nil {
		return doc, err
	}
	return doc, nil
}

func (c typedCollectionImpl[D]) Find(ctx context.Context, filter model.Filter, opts ...*options.FindOptions) ([]D, error) {
	cursor, err := c.internal.Find(ctx, filter, options.MergeFindOptions(opts...))
	if err != nil {
		return nil, err
	}
	var data []D
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c typedCollectionImpl[D]) FindByDocIds(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) ([]D, error) {
	filter := model.NewFilter().In("_id", lo.ToAnySlice[primitive.ObjectID](ids))
	return c.Find(ctx, filter, opts...)
}

func (c typedCollectionImpl[D]) CountDocuments(ctx context.Context, filter model.Filter, opts ...*options.CountOptions) (int64, error) {
	return c.internal.CountDocuments(ctx, filter, options.MergeCountOptions(opts...))
}

func (c typedCollectionImpl[D]) InsertOne(ctx context.Context, doc D) (primitive.ObjectID, error) {
	res, err := c.internal.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (c typedCollectionImpl[D]) InsertMany(ctx context.Context, docs []D) ([]primitive.ObjectID, error) {
	res, err := c.internal.InsertMany(ctx, lo.ToAnySlice[D](docs))
	if err != nil {
		return nil, err
	}
	data, _ := lo.FromAnySlice[primitive.ObjectID](res.InsertedIDs)
	return data, nil
}

func (c typedCollectionImpl[D]) FindOneAndUpdate(ctx context.Context, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (D, error) {
	var t D
	if len(m.Update) == 0 {
		return t, errors.New("update is empty")
	}
	singleResult := c.internal.FindOneAndUpdate(ctx, m.Filter, m.Update, options.MergeFindOneAndUpdateOptions(opts...))
	if singleResult.Err() != nil {
		return t, singleResult.Err()
	}
	if err := singleResult.Decode(t); err != nil {
		return t, err
	}
	return t, nil
}

func (c typedCollectionImpl[D]) UpdateOne(ctx context.Context, m *model.UpdateOne, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	return c.internal.UpdateOne(ctx, m.Filter, m.Update, options.MergeUpdateOptions(opts...))
}

func (c typedCollectionImpl[D]) UpdateMany(ctx context.Context, m *model.UpdateMany, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	res, err := c.internal.UpdateMany(ctx, m.Filter, m.Update, options.MergeUpdateOptions(opts...))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c typedCollectionImpl[D]) UpdateById(ctx context.Context, m *model.UpdateById, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	res, err := c.internal.UpdateByID(ctx, m.Id, m.Update, options.MergeUpdateOptions(opts...))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c typedCollectionImpl[D]) DeleteOne(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.internal.DeleteOne(ctx, filter, options.MergeDeleteOptions(opts...))
}

func (c typedCollectionImpl[D]) DeleteMany(ctx context.Context, filter model.Filter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.internal.DeleteMany(ctx, filter, options.MergeDeleteOptions(opts...))
}

func (c typedCollectionImpl[D]) InitializeBulkWriteOp(opts ...*options.BulkWriteOptions) *BulkWriteOperation {
	return newBulkWriteOperation(c.internal, opts...)
}
