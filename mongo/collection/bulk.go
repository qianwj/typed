// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/updates"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type TypedBulkWriteExecutor[D model.Doc[I], I model.ID] struct {
	coll   *mongo.Collection
	models []mongo.WriteModel
	opts   *rawopts.BulkWriteOptions
}

func newBulkWriteExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection) *TypedBulkWriteExecutor[D, I] {
	return &TypedBulkWriteExecutor[D, I]{
		coll: primary,
		opts: rawopts.BulkWrite(),
	}
}

// Comment sets the value for the Comment field.
func (b *TypedBulkWriteExecutor[D, I]) Comment(comment bson.UnorderedMap) *TypedBulkWriteExecutor[D, I] {
	b.opts.SetComment(comment)
	return b
}

// Unordered sets the value for the Ordered field.
func (b *TypedBulkWriteExecutor[D, I]) Unordered() *TypedBulkWriteExecutor[D, I] {
	b.opts.SetOrdered(false)
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
func (b *TypedBulkWriteExecutor[D, I]) Let(let bson.UnorderedMap) *TypedBulkWriteExecutor[D, I] {
	b.opts.SetLet(let)
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) UpdateOne(update *updates.UpdateOneModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) UpdateMany(update *updates.UpdateManyModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, update.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) ReplaceOne(replace *updates.ReplaceOneModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, replace.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) DeleteOne(delete *updates.DeleteOneModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, delete.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) DeleteMany(delete *updates.DeleteManyModel) *TypedBulkWriteExecutor[D, I] {
	b.models = append(b.models, delete.WriteModel())
	return b
}

func (b *TypedBulkWriteExecutor[D, I]) Execute(ctx context.Context) (*updates.BulkWriteResult[I], error) {
	if len(b.models) == 0 {
		return &updates.BulkWriteResult[I]{}, nil
	}
	res, err := b.coll.BulkWrite(ctx, b.models)
	return updates.FromBulkWriteResult[I](res), err
}
