package mongo

import (
	"github.com/qianwj/typed/mongo/executor"
	rawmongo "go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoDocuments        = rawmongo.ErrNoDocuments
	ErrClientDisconnected = rawmongo.ErrClientDisconnected
	ErrEmptyIndex         = executor.ErrEmptyIndex
)
