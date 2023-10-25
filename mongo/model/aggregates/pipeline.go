package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	stages mongo.Pipeline
}

func New() *Pipeline {
	return &Pipeline{
		stages: make(mongo.Pipeline, 0),
	}
}

func (p *Pipeline) append(key string, value any) {
	p.stages = append(p.stages, bson.D{
		{Key: key, Value: value},
	})
}

func (p *Pipeline) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.stages)
}

func (p *Pipeline) Stages() []primitive.D {
	return p.stages
}
