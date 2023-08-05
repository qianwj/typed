package aggregate

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Stage interface {
	Key() string
	Value() any
}

type Pipeline struct {
	stages []Stage
}

func (a *Pipeline) Append(stage ...Stage) *Pipeline {
	a.stages = append(a.stages, stage...)
	return a
}

func (a *Pipeline) Marshal() mongo.Pipeline {
	pipe := make([]bson.D, len(a.stages))

	for i, stage := range a.stages {
		pipe[i] = bson.D{
			{Key: stage.Key(), Value: stage.Value()},
		}
	}

	return pipe
}
