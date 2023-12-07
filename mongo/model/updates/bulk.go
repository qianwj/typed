package updates

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type TypedWriteModel interface {
	WriteModel() mongo.WriteModel
}
