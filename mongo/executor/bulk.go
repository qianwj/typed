package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/updates"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type TypedBulkWriteExecutor[D bson.Doc[I], I bson.ID] struct {
	coll   *mongo.Collection
	models []mongo.WriteModel
	opts   *rawopts.BulkWriteOptions
}

func NewBulkWriteExecutor[D bson.Doc[I], I bson.ID](primary *mongo.Collection) *TypedBulkWriteExecutor[D, I] {
	return &TypedBulkWriteExecutor[D, I]{
		coll: primary,
		opts: rawopts.BulkWrite(),
	}
}

// Comment sets the value for the Comment field.
func (b *TypedBulkWriteExecutor[D, I]) Comment(comment string) *TypedBulkWriteExecutor[D, I] {
	b.opts.SetComment(comment)
	return b
}

// Ordered sets the value for the Ordered field.
func (b *TypedBulkWriteExecutor[D, I]) Ordered() *TypedBulkWriteExecutor[D, I] {
	b.opts.SetOrdered(true)
	return b
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (b *TypedBulkWriteExecutor[D, I]) BypassDocumentValidation() *TypedBulkWriteExecutor[D, I] {
	b.opts.SetBypassDocumentValidation(true)
	return b
}

// Let sets the value for the Let field. Let specifies parameters for all update and delete commands in the BulkWrite.
// This option is only valid for MongoDB versions >= 5.0. Older servers will report an error for using this option.
// This must be a document mapping parameter names to values. Values must be constant or closed expressions that do not
// reference document fields. Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
func (b *TypedBulkWriteExecutor[D, I]) Let(let rawbson.M) *TypedBulkWriteExecutor[D, I] {
	b.opts.SetLet(let)
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) UpdateOne(update *updates.TypedUpdateOneModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) UpdateMany(update *updates.TypedUpdateManyModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) Execute(ctx context.Context) (*updates.BulkWriteResult[I], error) {
	if len(b.models) == 0 {
		return &updates.BulkWriteResult[I]{}, nil
	}
	res, err := b.coll.BulkWrite(ctx, b.models)
	return updates.FromBulkWriteResult[I](res), err
}
