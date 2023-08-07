package aggregate

import "go.mongodb.org/mongo-driver/mongo"

type Pipeline interface {
	Marshal() mongo.Pipeline
}
