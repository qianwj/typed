package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/model/update"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneAndUpdateExecutor[D model.Document[I], I model.DocumentId] struct {
	coll   *mongo.Collection
	filter *filter.Filter
	update *update.Update
	opts   *rawopts.FindOneAndUpdateOptions
}

func NewFindOneAndUpdateExecutor[D model.Document[I], I model.DocumentId](
	primary *mongo.Collection,
	filter *filter.Filter,
	update *update.Update,
) *FindOneAndUpdateExecutor[D, I] {
	return &FindOneAndUpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		opts:   rawopts.FindOneAndUpdate(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndUpdateExecutor[D, I]) BypassDocumentValidation() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetBypassDocumentValidation(true)
	return f
}

// ArrayFilters sets the value for the ArrayFilters field.
func (f *FindOneAndUpdateExecutor[D, I]) ArrayFilters(af options.ArrayFilters) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetArrayFilters(af.Raw())
	return f
}

// Collation sets the value for the Collation field.
func (f *FindOneAndUpdateExecutor[D, I]) Collation(collation *options.Collation) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetCollation((*rawopts.Collation)(collation))
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneAndUpdateExecutor[D, I]) Comment(comment interface{}) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndUpdateExecutor[D, I]) MaxTime(d time.Duration) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Projection sets the value for the Projection field.
func (f *FindOneAndUpdateExecutor[D, I]) Projection(projection *projections.ProjectionOptions) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetProjection(projection.Marshal())
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) ReturnBefore() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.Before)
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) ReturnAfter() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.After)
	return f
}

// Sort sets the value for the Sort field.
func (f *FindOneAndUpdateExecutor[D, I]) Sort(sort *sorts.SortOptions) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetSort(sort.Marshal())
	return f
}

// Upsert sets the value for the Upsert field.
func (f *FindOneAndUpdateExecutor[D, I]) Upsert() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetUpsert(true)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneAndUpdateExecutor[D, I]) Hint(index string) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindOneAndUpdateExecutor[D, I]) Let(l bson.M) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetLet(l)
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	res := f.coll.FindOneAndUpdate(ctx, f.filter.Marshal(), f.update.Marshal(), f.opts)
	var data D
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}
