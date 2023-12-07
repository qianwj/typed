package updates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReplaceOneModel struct {
	TypedWriteModel
	internal *mongo.ReplaceOneModel
}

func NewReplaceOneModel() *ReplaceOneModel {
	return &ReplaceOneModel{
		internal: mongo.NewReplaceOneModel(),
	}
}

// Hint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (rom *ReplaceOneModel) Hint(hint string) *ReplaceOneModel {
	rom.internal.SetHint(hint)
	return rom
}

// Filter specifies a filter to use to select the document to replace. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (rom *ReplaceOneModel) Filter(filter *filters.Filter) *ReplaceOneModel {
	rom.internal.SetFilter(filter)
	return rom
}

// Replacement specifies a document that will be used to replace the selected document. It cannot be nil and cannot
// contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
func (rom *ReplaceOneModel) Replacement(rep any) *ReplaceOneModel {
	rom.internal.SetReplacement(rep)
	return rom
}

// Collation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (rom *ReplaceOneModel) Collation(collation *options.Collation) *ReplaceOneModel {
	rom.internal.SetCollation(collation.Raw())
	return rom
}

// Upsert specifies whether or not the replacement document should be inserted if no document matching the filter is
// found. If an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (rom *ReplaceOneModel) Upsert() *ReplaceOneModel {
	rom.internal.SetUpsert(true)
	return rom
}

func (rom *ReplaceOneModel) WriteModel() mongo.WriteModel {
	return rom.internal
}
