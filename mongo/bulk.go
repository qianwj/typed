package mongo

import (
	"context"
	"github.com/qianwj/typed/mongo/model/modify"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type BulkWriteOperation struct {
	coll   *mongo.Collection
	models []mongo.WriteModel
	opts   []*options.BulkWriteOptions
}

func newBulkWriteOperation(coll *mongo.Collection, opts ...*options.BulkWriteOptions) *BulkWriteOperation {
	return &BulkWriteOperation{coll: coll, opts: opts}
}

func (b *BulkWriteOperation) UpdateOne(update *modify.TypedUpdateOneModel) *BulkWriteOperation {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *BulkWriteOperation) Execute(ctx context.Context) (*mongo.BulkWriteResult, error) {
	return b.coll.BulkWrite(ctx, b.models, options.MergeBulkWriteOptions(b.opts...))
}
