package typed_mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne[T any](ctx context.Context, c *mongo.Collection, doc T) (primitive.ObjectID, error) {
	res, err := c.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func InsertMany[T any](ctx context.Context, c *mongo.Collection, docs []*T) ([]primitive.ObjectID, error) {
	res, err := c.InsertMany(ctx, util.ToInterfaceSlice[*T](docs))
	if err != nil {
		return nil, err
	}
	return util.ToAnySlice[primitive.ObjectID](res.InsertedIDs), nil
}
