package typed_mongo

import (
	"context"
	"errors"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOneAndUpdate[T any](ctx context.Context, c *mongo.Collection, m *model.FindOneAndUpdate, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	singleResult := c.FindOneAndUpdate(ctx, m.Filter, m.Update, opts...)
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	var t T
	if err := singleResult.Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateMany[T any](ctx context.Context, c *mongo.Collection, m *model.UpdateMany, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	res, err := c.UpdateMany(ctx, m.Filter, m.Update, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateById[T any](ctx context.Context, c *mongo.Collection, m *model.UpdateById, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if len(m.Update) == 0 {
		return nil, errors.New("update is empty")
	}
	res, err := c.UpdateByID(ctx, m.Id, m.Update, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
