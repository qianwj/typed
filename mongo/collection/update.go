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
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/options"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateExecutor[D model.Doc[I], I model.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	update *updates.Update
	multi  bool
	docId  *I
	opts   *rawopts.UpdateOptions
}

func newUpdateOneExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		opts:   rawopts.Update(),
	}
}

func newUpdateManyExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		multi:  true,
		opts:   rawopts.Update(),
	}
}

func newUpdateByIdExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, id I, update *updates.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		docId:  &id,
		update: update,
		opts:   rawopts.Update(),
	}
}

func (u *UpdateExecutor[D, I]) ArrayFilters(af options.ArrayFilters) *UpdateExecutor[D, I] {
	u.opts.SetArrayFilters(af.Raw())
	return u
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (u *UpdateExecutor[D, I]) BypassDocumentValidation() *UpdateExecutor[D, I] {
	u.opts.SetBypassDocumentValidation(true)
	return u
}

// Collation sets the value for the Collation field.
func (u *UpdateExecutor[D, I]) Collation(c *options.Collation) *UpdateExecutor[D, I] {
	u.opts.SetCollation(c.Raw())
	return u
}

// Hint sets the value for the Hint field.
func (u *UpdateExecutor[D, I]) Hint(index string) *UpdateExecutor[D, I] {
	u.opts.SetHint(index)
	return u
}

// Upsert sets the value for the Upsert field.
func (u *UpdateExecutor[D, I]) Upsert() *UpdateExecutor[D, I] {
	u.opts.SetUpsert(true)
	return u
}

// Let sets the value for the Let field.
func (u *UpdateExecutor[D, I]) Let(l rawbson.M) *UpdateExecutor[D, I] {
	u.opts.SetLet(l)
	return u
}

func (u *UpdateExecutor[D, I]) Execute(ctx context.Context) (*updates.UpdateResult[I], error) {
	var (
		err error
		res *mongo.UpdateResult
	)
	if u.docId != nil {
		res, err = u.coll.UpdateByID(ctx, u.docId, u.update, u.opts)
	} else if u.multi {
		res, err = u.coll.UpdateMany(ctx, u.filter, u.update, u.opts)
	} else {
		res, err = u.coll.UpdateOne(ctx, u.filter, u.update, u.opts)
	}
	return updates.FromUpdateResult[I](res), err
}
