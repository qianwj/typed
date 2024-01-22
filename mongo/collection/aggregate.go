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
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/options"
	"github.com/qianwj/typed/streams"
	rawbson "go.mongodb.org/mongo-driver/bson"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AggregateExecutor[D model.Doc[I], I model.ID] struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	pipe            *aggregates.Pipeline
	primary         bool
	opts            *rawopts.AggregateOptions
}

func NewAggregator[D model.Doc[I], I model.ID](coll *raw.Collection, pipe *aggregates.Pipeline) *AggregateExecutor[D, I] {
	return &AggregateExecutor[D, I]{
		readprefPrimary: coll,
		readprefDefault: coll,
		pipe:            pipe,
		opts:            rawopts.Aggregate(),
	}
}

func newAggregateExecutor[D model.Doc[I], I model.ID](
	readprefPrimary, readprefDefault *raw.Collection,
	pipe *aggregates.Pipeline,
) *AggregateExecutor[D, I] {
	return &AggregateExecutor[D, I]{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		pipe:            pipe,
		opts:            rawopts.Aggregate(),
	}
}

func (a *AggregateExecutor[D, I]) Primary() *AggregateExecutor[D, I] {
	a.primary = true
	return a
}

// AllowDiskUse sets the value for the AllowDiskUse field.
func (a *AggregateExecutor[D, I]) AllowDiskUse() *AggregateExecutor[D, I] {
	a.opts.SetAllowDiskUse(true)
	return a
}

// BatchSize sets the value for the BatchSize field.
func (a *AggregateExecutor[D, I]) BatchSize(i int32) *AggregateExecutor[D, I] {
	a.opts.SetBatchSize(i)
	return a
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *AggregateExecutor[D, I]) BypassDocumentValidation() *AggregateExecutor[D, I] {
	a.opts.SetBypassDocumentValidation(true)
	return a
}

// Collation sets the value for the Collation field.
func (a *AggregateExecutor[D, I]) Collation(c *options.Collation) *AggregateExecutor[D, I] {
	a.opts.SetCollation(c.Raw())
	return a
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (a *AggregateExecutor[D, I]) MaxTime(d time.Duration) *AggregateExecutor[D, I] {
	a.opts.SetMaxTime(d)
	return a
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *AggregateExecutor[D, I]) MaxAwaitTime(d time.Duration) *AggregateExecutor[D, I] {
	a.opts.SetMaxAwaitTime(d)
	return a
}

// Comment sets the value for the Comment field.
func (a *AggregateExecutor[D, I]) Comment(s string) *AggregateExecutor[D, I] {
	a.opts.SetComment(s)
	return a
}

// Hint sets the value for the Hint field.
func (a *AggregateExecutor[D, I]) Hint(index string) *AggregateExecutor[D, I] {
	a.opts.SetHint(index)
	return a
}

// Let sets the value for the Let field.
func (a *AggregateExecutor[D, I]) Let(let any) *AggregateExecutor[D, I] {
	a.opts.SetLet(let)
	return a
}

// Custom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (a *AggregateExecutor[D, I]) Custom(c rawbson.M) *AggregateExecutor[D, I] {
	a.opts.SetCustom(c)
	return a
}

func (a *AggregateExecutor[D, I]) Collect(ctx context.Context, result any) error {
	cursor, err := a.cursor(ctx)
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}

func (a *AggregateExecutor[D, I]) Stream(ctx context.Context) (streams.Publisher[D], error) {
	cursor, err := a.cursor(ctx)
	if err != nil {
		return nil, err
	}
	return fromCursor[D](ctx, cursor), nil
}

func (a *AggregateExecutor[D, I]) cursor(ctx context.Context) (*raw.Cursor, error) {
	var (
		err    error
		cursor *raw.Cursor
	)
	if a.primary {
		cursor, err = a.readprefPrimary.Aggregate(ctx, a.pipe.Stages(), a.opts)
	} else {
		cursor, err = a.readprefDefault.Aggregate(ctx, a.pipe.Stages(), a.opts)
	}
	return cursor, err
}
