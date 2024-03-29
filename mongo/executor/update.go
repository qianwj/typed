package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/update"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateExecutor[D model.Document[I], I model.DocumentId] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	update *update.Update
	multi  bool
	docId  *I
	opts   *rawopts.UpdateOptions
}

func NewUpdateOneExecutor[D model.Document[I], I model.DocumentId](primary *mongo.Collection, filter *filters.Filter, update *update.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		opts:   rawopts.Update(),
	}
}

func NewUpdateManyExecutor[D model.Document[I], I model.DocumentId](primary *mongo.Collection, filter *filters.Filter, update *update.Update) *UpdateExecutor[D, I] {
	return &UpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		multi:  true,
		opts:   rawopts.Update(),
	}
}

func NewUpdateByIdExecutor[D model.Document[I], I model.DocumentId](primary *mongo.Collection, id I, update *update.Update) *UpdateExecutor[D, I] {
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
func (u *UpdateExecutor[D, I]) Let(l bson.M) *UpdateExecutor[D, I] {
	u.opts.SetLet(l)
	return u
}

func (u *UpdateExecutor[D, I]) Execute(ctx context.Context) (*model.UpdateResult[I], error) {
	var (
		err error
		res *mongo.UpdateResult
	)
	if u.docId != nil {
		res, err = u.coll.UpdateByID(ctx, u.docId, u.update.Marshal(), u.opts)
	} else if u.multi {
		res, err = u.coll.UpdateMany(ctx, u.filter.Marshal(), u.update.Marshal(), u.opts)
	} else {
		res, err = u.coll.UpdateOne(ctx, u.filter.Marshal(), u.update.Marshal(), u.opts)
	}
	return model.FromUpdateResult[I](res), err
}
