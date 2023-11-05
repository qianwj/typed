package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	rawbson "go.mongodb.org/mongo-driver/bson"
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
	p.stages = append(p.stages, bson.D(
		bson.E(key, value),
	).Primitive())
}

func (p *Pipeline) MarshalBSON() ([]byte, error) {
	return rawbson.Marshal(p.stages)
}

func (p *Pipeline) Stages() []primitive.D {
	return p.stages
}
