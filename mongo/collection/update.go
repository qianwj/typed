package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/options"
	rawbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateExecutor[D bson.Doc[I], I bson.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	update *updates.Update
	multi  bool
	docId  *I
	opts   *rawopts.UpdateOptions
}

func newUpdateOneExecutor[D bson.Doc[I], I bson.ID](primary *mongo.Collection, filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		opts:   rawopts.Update(),
	}
}

func newUpdateManyExecutor[D bson.Doc[I], I bson.ID](primary *mongo.Collection, filter *filters.Filter, update *updates.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		multi:  true,
		opts:   rawopts.Update(),
	}
}

func newUpdateByIdExecutor[D bson.Doc[I], I bson.ID](primary *mongo.Collection, id I, update *updates.Update) *UpdateExecutor[D, I] {
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
	u.opts.SetCollation((*rawopts.Collation)(c))
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
