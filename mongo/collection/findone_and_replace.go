package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneAndReplaceExecutor[D bson.Doc[I], I bson.ID] struct {
	coll        *mongo.Collection
	filter      *filters.Filter
	replacement any
	opts        *rawopts.FindOneAndReplaceOptions
}

func newFindOneAndReplaceExecutor[D bson.Doc[I], I bson.ID](
	primary *mongo.Collection,
	filter *filters.Filter,
	replacement any,
) *FindOneAndReplaceExecutor[D, I] {
	return &FindOneAndReplaceExecutor[D, I]{
		coll:        primary,
		filter:      filter,
		replacement: replacement,
		opts:        rawopts.FindOneAndReplace(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndReplaceExecutor[D, I]) BypassDocumentValidation() *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetBypassDocumentValidation(true)
	return f
}

// Collation sets the value for the Collation field.
func (f *FindOneAndReplaceExecutor[D, I]) Collation(collation *options.Collation) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetCollation((*rawopts.Collation)(collation))
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneAndReplaceExecutor[D, I]) Comment(comment interface{}) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndReplaceExecutor[D, I]) MaxTime(d time.Duration) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Projection sets the value for the Projection field.
func (f *FindOneAndReplaceExecutor[D, I]) Projection(projection *projections.Options) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetProjection(projection)
	return f
}

func (f *FindOneAndReplaceExecutor[D, I]) ReturnBefore() *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.Before)
	return f
}

func (f *FindOneAndReplaceExecutor[D, I]) ReturnAfter() *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.After)
	return f
}

// Sort sets the value for the Sort field.
func (f *FindOneAndReplaceExecutor[D, I]) Sort(sort *sorts.Options) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

// Upsert sets the value for the Upsert field.
func (f *FindOneAndReplaceExecutor[D, I]) Upsert() *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetUpsert(true)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneAndReplaceExecutor[D, I]) Hint(index string) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindOneAndReplaceExecutor[D, I]) Let(l bson.UnorderedMap) *FindOneAndReplaceExecutor[D, I] {
	f.opts.SetLet(l)
	return f
}

func (f *FindOneAndReplaceExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	res := f.coll.FindOneAndReplace(ctx, f.filter, f.replacement, f.opts)
	var data D
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneAndReplaceExecutor[D, I]) Collect(ctx context.Context, data any) error {
	res := f.coll.FindOneAndReplace(ctx, f.filter, f.replacement, f.opts)
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}
