package updates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteOneModel struct {
	TypedWriteModel
	internal *mongo.DeleteOneModel
}

func NewDeleteOne() *DeleteOneModel {
	return &DeleteOneModel{
		internal: mongo.NewDeleteOneModel(),
	}
}

// Filter specifies a filter to use to select the document to delete. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (dom *DeleteOneModel) Filter(filter *filters.Filter) *DeleteOneModel {
	dom.internal.SetFilter(filter)
	return dom
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dom *DeleteOneModel) Collation(collation *options.Collation) *DeleteOneModel {
	dom.internal.SetCollation(collation.Raw())
	return dom
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dom *DeleteOneModel) Hint(hint string) *DeleteOneModel {
	dom.internal.SetHint(hint)
	return dom
}

func (dom *DeleteOneModel) WriteModel() mongo.WriteModel {
	return dom.internal
}

type DeleteManyModel struct {
	TypedWriteModel
	internal *mongo.DeleteManyModel
}

func NewDeleteMany() *DeleteManyModel {
	return &DeleteManyModel{
		internal: mongo.NewDeleteManyModel(),
	}
}

// Filter specifies a filter to use to select the document to delete. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (dmm *DeleteManyModel) Filter(filter *filters.Filter) *DeleteManyModel {
	dmm.internal.SetFilter(filter)
	return dmm
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dmm *DeleteManyModel) Collation(collation *options.Collation) *DeleteManyModel {
	dmm.internal.SetCollation(collation.Raw())
	return dmm
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dmm *DeleteManyModel) Hint(hint string) *DeleteManyModel {
	dmm.internal.SetHint(hint)
	return dmm
}

func (dmm *DeleteManyModel) WriteModel() mongo.WriteModel {
	return dmm.internal
}
