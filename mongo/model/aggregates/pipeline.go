package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
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
	p.stages = append(p.stages, bson.D(bson.E(key, value)).Primitive())
}

func (p *Pipeline) Stages() mongo.Pipeline {
	return p.stages
}
