package updates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateOneModel struct {
	TypedWriteModel
	internal *mongo.UpdateOneModel
}

// NewUpdateOne creates a new UpdateOneModel.
func NewUpdateOne() *UpdateOneModel {
	return &UpdateOneModel{internal: mongo.NewUpdateOneModel()}
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (uom *UpdateOneModel) Hint(index string) *UpdateOneModel {
	uom.internal.SetHint(index)
	return uom
}

// Filter specifies a filter to use to select the document to updates. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (uom *UpdateOneModel) Filter(filter *filters.Filter) *UpdateOneModel {
	uom.internal.SetFilter(filter)
	return uom
}

// Update specifies the modifications to be made to the selected document. The value must be a document containing
// updates operators (https://docs.mongodb.com/manual/reference/operator/update/). It cannot be nil or empty.
func (uom *UpdateOneModel) Update(update *Update) *UpdateOneModel {
	uom.internal.SetUpdate(update)
	return uom
}

// ArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (uom *UpdateOneModel) ArrayFilters(filters options.ArrayFilters) *UpdateOneModel {
	uom.internal.SetArrayFilters(filters.Raw())
	return uom
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (uom *UpdateOneModel) Collation(collation *options.Collation) *UpdateOneModel {
	uom.internal.SetCollation(collation.Raw())
	return uom
}

// Upsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (uom *UpdateOneModel) Upsert() *UpdateOneModel {
	uom.internal.SetUpsert(true)
	return uom
}

func (uom *UpdateOneModel) WriteModel() mongo.WriteModel {
	return uom.internal
}

type UpdateManyModel struct {
	TypedWriteModel
	internal *mongo.UpdateManyModel
}

// NewUpdateMany creates a new UpdateManyModel.
func NewUpdateMany() *UpdateManyModel {
	return &UpdateManyModel{internal: mongo.NewUpdateManyModel()}
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (umm *UpdateManyModel) Hint(index string) *UpdateManyModel {
	umm.internal.SetHint(index)
	return umm
}

// Filter specifies a filter to use to select documents to updates. The filter must be a document containing query
// operators. It cannot be nil.
func (umm *UpdateManyModel) Filter(filter *filters.Filter) *UpdateManyModel {
	umm.internal.SetFilter(filter)
	return umm
}

// Update specifies the modifications to be made to the selected documents. The value must be a document containing
// updates operators (https://www.mongodb.com/docs/manual/reference/operator/update/). It cannot be nil or empty.
func (umm *UpdateManyModel) Update(update *Update) *UpdateManyModel {
	umm.internal.SetUpdate(update)
	return umm
}

// ArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (umm *UpdateManyModel) ArrayFilters(filters options.ArrayFilters) *UpdateManyModel {
	umm.internal.SetArrayFilters(filters.Raw())
	return umm
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (umm *UpdateManyModel) Collation(collation *options.Collation) *UpdateManyModel {
	umm.internal.SetCollation(collation.Raw())
	return umm
}

// Upsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (umm *UpdateManyModel) Upsert() *UpdateManyModel {
	umm.internal.SetUpsert(true)
	return umm
}

func (umm *UpdateManyModel) WriteModel() mongo.WriteModel {
	return umm.internal
}
