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

package database

import (
	"context"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/options"
	rawbson "go.mongodb.org/mongo-driver/bson"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AggregateExecutor struct {
	readprefPrimary *raw.Database
	readprefDefault *raw.Database
	pipe            *aggregates.Pipeline
	primary         bool
	opts            *rawopts.AggregateOptions
}

func newAggregateExecutor(readprefPrimary, readprefDefault *raw.Database, pipe *aggregates.Pipeline) *AggregateExecutor {
	return &AggregateExecutor{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		pipe:            pipe,
		opts:            rawopts.Aggregate(),
	}
}

func (a *AggregateExecutor) Primary() *AggregateExecutor {
	a.primary = true
	return a
}

// AllowDiskUse sets the value for the AllowDiskUse field.
func (a *AggregateExecutor) AllowDiskUse() *AggregateExecutor {
	a.opts.SetAllowDiskUse(true)
	return a
}

// BatchSize sets the value for the BatchSize field.
func (a *AggregateExecutor) BatchSize(i int32) *AggregateExecutor {
	a.opts.SetBatchSize(i)
	return a
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *AggregateExecutor) BypassDocumentValidation() *AggregateExecutor {
	a.opts.SetBypassDocumentValidation(true)
	return a
}

// Collation sets the value for the Collation field.
func (a *AggregateExecutor) Collation(c *options.Collation) *AggregateExecutor {
	a.opts.SetCollation(c.Raw())
	return a
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (a *AggregateExecutor) MaxTime(d time.Duration) *AggregateExecutor {
	a.opts.SetMaxTime(d)
	return a
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *AggregateExecutor) MaxAwaitTime(d time.Duration) *AggregateExecutor {
	a.opts.SetMaxAwaitTime(d)
	return a
}

// Comment sets the value for the Comment field.
func (a *AggregateExecutor) Comment(s string) *AggregateExecutor {
	a.opts.SetComment(s)
	return a
}

// Hint sets the value for the Hint field.
func (a *AggregateExecutor) Hint(index string) *AggregateExecutor {
	a.opts.SetHint(index)
	return a
}

// Let sets the value for the Let field.
func (a *AggregateExecutor) Let(let rawbson.M) *AggregateExecutor {
	a.opts.SetLet(let)
	return a
}

// Custom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (a *AggregateExecutor) Custom(c rawbson.M) *AggregateExecutor {
	a.opts.SetCustom(c)
	return a
}

func (a *AggregateExecutor) Collect(ctx context.Context, result any) error {
	var (
		err    error
		cursor *raw.Cursor
	)
	if a.primary {
		cursor, err = a.readprefPrimary.Aggregate(ctx, a.pipe.Stages(), a.opts)
	} else {
		cursor, err = a.readprefDefault.Aggregate(ctx, a.pipe.Stages(), a.opts)
	}
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}
