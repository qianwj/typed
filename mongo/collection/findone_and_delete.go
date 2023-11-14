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

type FindOneAndDeleteExecutor[D bson.Doc[I], I bson.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	opts   *rawopts.FindOneAndDeleteOptions
}

func newFindOneAndDeleteExecutor[D bson.Doc[I], I bson.ID](
	primary *mongo.Collection,
	filter *filters.Filter,
) *FindOneAndDeleteExecutor[D, I] {
	return &FindOneAndDeleteExecutor[D, I]{
		coll:   primary,
		filter: filter,
		opts:   rawopts.FindOneAndDelete(),
	}
}

// Collation sets the value for the Collation field.
func (f *FindOneAndDeleteExecutor[D, I]) Collation(collation *options.Collation) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetCollation((*rawopts.Collation)(collation))
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneAndDeleteExecutor[D, I]) Comment(comment interface{}) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndDeleteExecutor[D, I]) MaxTime(d time.Duration) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Projection sets the value for the Projection field.
func (f *FindOneAndDeleteExecutor[D, I]) Projection(projection *projections.Options) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetProjection(projection)
	return f
}

// Sort sets the value for the Sort field.
func (f *FindOneAndDeleteExecutor[D, I]) Sort(sort *sorts.Options) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneAndDeleteExecutor[D, I]) Hint(index string) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindOneAndDeleteExecutor[D, I]) Let(l bson.UnorderedMap) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetLet(l)
	return f
}

func (f *FindOneAndDeleteExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	res := f.coll.FindOneAndDelete(ctx, f.filter, f.opts)
	var data D
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneAndDeleteExecutor[D, I]) Collect(ctx context.Context, data any) error {
	res := f.coll.FindOneAndDelete(ctx, f.filter, f.opts)
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}
