package mongo

import (
	"context"
	"errors"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/options"
)

func Aggregate[D model.Document, U any](ctx context.Context, c TypedCollection[D], pipeline model.AggregatePipeline, opts ...*options.AggregateOptions) ([]*U, error) {
	if len(pipeline) == 0 {
		return nil, errors.New("pipeline must not empty")
	}
	cursor, err := c.collection().Aggregate(ctx, pipeline, options.MergeAggregateOptions(opts...))
	if err != nil {
		return nil, err
	}
	var data []*U
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}
