package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type typedWriteModel interface {
	WriteModel() mongo.WriteModel
}

type TypedUpdateOneModel struct {
	typedWriteModel
	internal *mongo.UpdateOneModel
}

func NewTypeUpdateOneModel() *TypedUpdateOneModel {
	return &TypedUpdateOneModel{internal: mongo.NewUpdateOneModel()}
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (uom *TypedUpdateOneModel) Hint(hint interface{}) *TypedUpdateOneModel {
	uom.internal.Hint = hint
	return uom
}

// Filter specifies a filter to use to select the document to update. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (uom *TypedUpdateOneModel) Filter(filter Filter) *TypedUpdateOneModel {
	uom.internal.Filter = filter
	return uom
}

// Update specifies the modifications to be made to the selected document. The value must be a document containing
// update operators (https://docs.mongodb.com/manual/reference/operator/update/). It cannot be nil or empty.
func (uom *TypedUpdateOneModel) Update(update Update) *TypedUpdateOneModel {
	uom.internal.Update = update
	return uom
}

// ArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (uom *TypedUpdateOneModel) ArrayFilters(filters options.ArrayFilters) *TypedUpdateOneModel {
	uom.internal.ArrayFilters = &filters
	return uom
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (uom *TypedUpdateOneModel) Collation(collation *options.Collation) *TypedUpdateOneModel {
	uom.internal.Collation = collation
	return uom
}

// Upsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (uom *TypedUpdateOneModel) Upsert(upsert bool) *TypedUpdateOneModel {
	uom.internal.Upsert = &upsert
	return uom
}

func (uom *TypedUpdateOneModel) WriteModel() mongo.WriteModel {
	return uom.internal
}
