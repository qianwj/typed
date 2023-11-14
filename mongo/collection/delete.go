package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type DeleteExecutor[D model.Doc[I], I model.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	multi  bool
	opts   *rawopts.DeleteOptions
}

func newDeleteOneExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, filter *filters.Filter) *DeleteExecutor[D, I] {
	return &DeleteExecutor[D, I]{
		coll:   primary,
		filter: filter,
		opts:   rawopts.Delete(),
	}
}

func newDeleteManyExecutor[D model.Doc[I], I model.ID](primary *mongo.Collection, filter *filters.Filter) *DeleteExecutor[D, I] {
	return &DeleteExecutor[D, I]{
		coll:   primary,
		filter: filter,
		multi:  true,
		opts:   rawopts.Delete(),
	}
}

// Collation sets the value for the Collation field.
func (d *DeleteExecutor[D, I]) Collation(c *options.Collation) *DeleteExecutor[D, I] {
	d.opts.SetCollation((*rawopts.Collation)(c))
	return d
}

// Comment sets the value for the Comment field.
func (d *DeleteExecutor[D, I]) Comment(comment string) *DeleteExecutor[D, I] {
	d.opts.SetComment(comment)
	return d
}

// Hint sets the value for the Hint field.
func (d *DeleteExecutor[D, I]) Hint(index string) *DeleteExecutor[D, I] {
	d.opts.SetHint(index)
	return d
}

// Let sets the value for the Let field.
func (d *DeleteExecutor[D, I]) Let(let rawbson.M) *DeleteExecutor[D, I] {
	d.opts.SetLet(let)
	return d
}

func (d *DeleteExecutor[D, I]) Execute(ctx context.Context) (int64, error) {
	var (
		err error
		res *mongo.DeleteResult
	)
	if d.multi {
		res, err = d.coll.DeleteMany(ctx, d.filter, d.opts)
	} else {
		res, err = d.coll.DeleteOne(ctx, d.filter, d.opts)
	}
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}
