package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneExecutor[D model.Document[I], I model.DocumentId] struct {
	coll *Collection[D, I]
	data D
	opts *rawopts.InsertOneOptions
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (i *InsertOneExecutor[D, I]) BypassDocumentValidation() *InsertOneExecutor[D, I] {
	i.opts.SetBypassDocumentValidation(true)
	return i
}

// Comment sets the value for the Comment field.
func (i *InsertOneExecutor[D, I]) Comment(comment string) *InsertOneExecutor[D, I] {
	i.opts.SetComment(comment)
	return i
}

func (i *InsertOneExecutor[D, I]) Execute(ctx context.Context) (I, error) {
	res, err := i.coll.primary.InsertOne(ctx, i.data, i.opts)
	var id I
	if err != nil {
		return id, err
	}
	return res.InsertedID.(I), nil
}

type InsertManyExecutor[D model.Document[I], I model.DocumentId] struct {
	coll *Collection[D, I]
	data []any
	opts *rawopts.InsertManyOptions
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (i *InsertManyExecutor[D, I]) BypassDocumentValidation() *InsertManyExecutor[D, I] {
	i.opts.SetBypassDocumentValidation(true)
	return i
}

// Comment sets the value for the Comment field.
func (i *InsertManyExecutor[D, I]) Comment(comment string) *InsertManyExecutor[D, I] {
	i.opts.SetComment(comment)
	return i
}

// SetOrdered sets the value for the Ordered field.
func (i *InsertManyExecutor[D, I]) SetOrdered() *InsertManyExecutor[D, I] {
	i.opts.SetOrdered(true)
	return i
}

func (i *InsertManyExecutor[D, I]) Add(docs ...D) *InsertManyExecutor[D, I] {
	i.data = append(i.data, toAny(docs))
	return i
}

func (i *InsertManyExecutor[D, I]) Execute(ctx context.Context) ([]I, error) {
	if len(i.data) == 0 {
		return make([]I, 0), nil
	}
	res, err := i.coll.primary.InsertMany(ctx, i.data, i.opts)
	return mapTo(res.InsertedIDs, func(i any) I {
		return i.(I)
	}), err
}
