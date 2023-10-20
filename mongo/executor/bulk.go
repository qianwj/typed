package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/update"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type BulkWriteExecutor[D bson.Doc[I], I bson.ID] struct {
	coll   *mongo.Collection
	models []mongo.WriteModel
	opts   *rawopts.BulkWriteOptions
}

func NewBulkWriteExecutor[D bson.Doc[I], I bson.ID](primary *mongo.Collection) *BulkWriteExecutor[D, I] {
	return &BulkWriteExecutor[D, I]{
		coll: primary,
		opts: rawopts.BulkWrite(),
	}
}

// Comment sets the value for the Comment field.
func (b *BulkWriteExecutor[D, I]) Comment(comment string) *BulkWriteExecutor[D, I] {
	b.opts.SetComment(comment)
	return b
}

// Ordered sets the value for the Ordered field.
func (b *BulkWriteExecutor[D, I]) Ordered() *BulkWriteExecutor[D, I] {
	b.opts.SetOrdered(true)
	return b
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (b *BulkWriteExecutor[D, I]) BypassDocumentValidation() *BulkWriteExecutor[D, I] {
	b.opts.SetBypassDocumentValidation(true)
	return b
}

// Let sets the value for the Let field. Let specifies parameters for all update and delete commands in the BulkWrite.
// This option is only valid for MongoDB versions >= 5.0. Older servers will report an error for using this option.
// This must be a document mapping parameter names to values. Values must be constant or closed expressions that do not
// reference document fields. Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
func (b *BulkWriteExecutor[D, I]) Let(let rawbson.M) *BulkWriteExecutor[D, I] {
	b.opts.SetLet(let)
	return b
}

func (b *BulkWriteExecutor[D, I]) UpdateOne(update *update.TypedUpdateOneModel) *BulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *BulkWriteExecutor[D, I]) UpdateMany(update *update.TypedUpdateManyModel) *BulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *BulkWriteExecutor[D, I]) Execute(ctx context.Context) (*model.BulkWriteResult[I], error) {
	if len(b.models) == 0 {
		return &model.BulkWriteResult[I]{}, nil
	}
	res, err := b.coll.BulkWrite(ctx, b.models)
	return model.FromBulkWriteResult[I](res), err
}
