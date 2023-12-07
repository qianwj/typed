package updates

import "go.mongodb.org/mongo-driver/mongo"

type InsertOneModel struct {
	TypedWriteModel
	internal *mongo.InsertOneModel
}

func NewInsertOne() *InsertOneModel {
	return &InsertOneModel{
		internal: mongo.NewInsertOneModel(),
	}
}

// Document specifies the document to be inserted. The document cannot be nil. If it does not have an _id field when
// transformed into BSON, one will be added automatically to the marshalled document. The original document will not be
// modified.
func (iom *InsertOneModel) Document(doc any) *InsertOneModel {
	iom.internal.SetDocument(doc)
	return iom
}

func (iom *InsertOneModel) WriteModel() mongo.WriteModel {
	return iom.internal
}
