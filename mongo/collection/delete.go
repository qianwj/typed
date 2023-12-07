package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type DeleteExecutor struct {
	coll   *mongo.Collection
	filter *filters.Filter
	opts   *rawopts.DeleteOptions
}

func newDeleteExecutor(coll *mongo.Collection, filter *filters.Filter) *DeleteExecutor {
	return &DeleteExecutor{
		coll:   coll,
		filter: filter,
		opts:   rawopts.Delete(),
	}
}

// Collation sets the value for the Collation field.
func (d *DeleteExecutor) Collation(c *options.Collation) *DeleteExecutor {
	d.opts.SetCollation((*rawopts.Collation)(c))
	return d
}

// Comment sets the value for the Comment field.
func (d *DeleteExecutor) Comment(comment bson.UnorderedMap) *DeleteExecutor {
	d.opts.SetComment(comment)
	return d
}

// Hint sets the value for the Hint field.
func (d *DeleteExecutor) Hint(index string) *DeleteExecutor {
	d.opts.SetHint(index)
	return d
}

// Let sets the value for the Let field.
func (d *DeleteExecutor) Let(let bson.UnorderedMap) *DeleteExecutor {
	d.opts.SetLet(let)
	return d
}

func (d *DeleteExecutor) One(ctx context.Context) (int64, error) {
	res, err := d.coll.DeleteOne(ctx, d.filter, d.opts)
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}

func (d *DeleteExecutor) Many(ctx context.Context) (int64, error) {
	res, err := d.coll.DeleteMany(ctx, d.filter, d.opts)
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}
