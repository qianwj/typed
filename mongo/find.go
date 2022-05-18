package typed_mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T any](ctx context.Context, c *mongo.Collection, filter model.Filter, opts ...*options.FindOneOptions) (*T, error) {
	singleResult := c.FindOne(ctx, filter, opts...)
	var doc T
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	if err := singleResult.Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func FindByDocIds[T any](ctx context.Context, c *mongo.Collection, ids []primitive.ObjectID, opts ...*options.FindOptions) ([]*T, error) {
	filter := model.NewFilter().In("_id", util.ToInterfaceSlice[primitive.ObjectID](ids))
	return Find[T](ctx, c, filter, opts...)
}

func Find[T any](ctx context.Context, c *mongo.Collection, filter model.Filter, opts ...*options.FindOptions) ([]*T, error) {
	cursor, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	var data []*T
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}
