package mongo

import (
	rawmongo "go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoDocuments        = rawmongo.ErrNoDocuments
	ErrClientDisconnected = rawmongo.ErrClientDisconnected
)
