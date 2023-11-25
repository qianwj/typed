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
	"errors"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneExecutor[D model.Doc[I], I model.ID] struct {
	coll *mongo.Collection
	data D
	opts *rawopts.InsertOneOptions
}

func newInsertOneExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, doc D) *InsertOneExecutor[D, I] {
	return &InsertOneExecutor[D, I]{
		coll: primary,
		data: doc,
		opts: rawopts.InsertOne(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (i *InsertOneExecutor[D, I]) BypassDocumentValidation() *InsertOneExecutor[D, I] {
	i.opts.SetBypassDocumentValidation(true)
	return i
}

// Comment sets the value for the Comment field.
func (i *InsertOneExecutor[D, I]) Comment(comment bson.M) *InsertOneExecutor[D, I] {
	i.opts.SetComment(comment)
	return i
}

func (i *InsertOneExecutor[D, I]) Execute(ctx context.Context) (I, error) {
	res, err := i.coll.InsertOne(ctx, i.data, i.opts)
	var id I
	if err != nil {
		return id, err
	}
	return res.InsertedID.(I), nil
}

type InsertManyExecutor[D model.Doc[I], I model.ID] struct {
	coll *mongo.Collection
	data []any
	opts *rawopts.InsertManyOptions
}

func newInsertManyExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, docs ...D) *InsertManyExecutor[D, I] {
	return &InsertManyExecutor[D, I]{
		coll: primary,
		data: util.ToAny(docs),
		opts: rawopts.InsertMany(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (i *InsertManyExecutor[D, I]) BypassDocumentValidation() *InsertManyExecutor[D, I] {
	i.opts.SetBypassDocumentValidation(true)
	return i
}

// Comment sets the value for the Comment field.
func (i *InsertManyExecutor[D, I]) Comment(comment bson.M) *InsertManyExecutor[D, I] {
	i.opts.SetComment(comment)
	return i
}

// Ordered sets the value for the Ordered field.
func (i *InsertManyExecutor[D, I]) Ordered() *InsertManyExecutor[D, I] {
	i.opts.SetOrdered(true)
	return i
}

func (i *InsertManyExecutor[D, I]) Add(docs ...D) *InsertManyExecutor[D, I] {
	i.data = append(i.data, util.ToAny(docs))
	return i
}

func (i *InsertManyExecutor[D, I]) Execute(ctx context.Context) ([]I, error) {
	if len(i.data) == 0 {
		return make([]I, 0), nil
	}
	res, err := i.coll.InsertMany(ctx, i.data, i.opts)
	return util.OrderedMap(res.InsertedIDs, func(i any) I {
		return i.(I)
	}), err
}

type InsertExecutor struct {
	coll                     *mongo.Collection
	bypassDocumentValidation *bool
	comment                  *bson.M
	ordered                  *bool
	docs                     []any
}

func newInsertExecutor(coll *mongo.Collection, docs ...any) *InsertExecutor {
	return &InsertExecutor{
		coll: coll,
		docs: docs,
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (i *InsertExecutor) BypassDocumentValidation() *InsertExecutor {
	i.bypassDocumentValidation = util.ToPtr(true)
	return i
}

// Comment sets the value for the Comment field.
func (i *InsertExecutor) Comment(comment bson.M) *InsertExecutor {
	i.comment = util.ToPtr(comment)
	return i
}

// Ordered sets the value for the Ordered field.
func (i *InsertExecutor) Ordered() *InsertExecutor {
	i.ordered = util.ToPtr(true)
	return i
}

// Add append insert data.
func (i *InsertExecutor) Add(docs ...any) *InsertExecutor {
	i.docs = append(i.docs, docs...)
	return i
}

func (i *InsertExecutor) One(ctx context.Context) (any, error) {
	if len(i.docs) < 1 {
		return nil, errors.New("insert empty data")
	}
	return i.coll.InsertOne(ctx, i.docs[0], i.insertOneOptions())
}

func (i *InsertExecutor) Many(ctx context.Context) ([]any, error) {
	if len(i.docs) < 1 {
		return nil, errors.New("insert empty data")
	}
	res, err := i.coll.InsertMany(ctx, i.docs, i.insertManyOptions())
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

func (i *InsertExecutor) insertOneOptions() *rawopts.InsertOneOptions {
	opts := rawopts.InsertOne()
	if i.bypassDocumentValidation != nil {
		opts = opts.SetBypassDocumentValidation(true)
	}
	if i.comment != nil {
		opts = opts.SetComment(true)
	}
	return opts
}

func (i *InsertExecutor) insertManyOptions() *rawopts.InsertManyOptions {
	opts := rawopts.InsertMany()
	if i.bypassDocumentValidation != nil {
		opts = opts.SetBypassDocumentValidation(true)
	}
	if i.comment != nil {
		opts = opts.SetComment(true)
	}
	if i.ordered != nil {
		opts = opts.SetOrdered(true)
	}
	return opts
}
