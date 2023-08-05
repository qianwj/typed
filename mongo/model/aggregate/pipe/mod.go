package pipe

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	stages mongo.Pipeline
	dict   map[string]int
}

func (p *Pipeline) put(key string, value any) {
	if index, exists := p.dict[key]; exists {
		p.stages[index] = bson.D{
			{Key: key, Value: value},
		}
	} else {
		p.stages = append(p.stages, bson.D{
			{Key: key, Value: value},
		})
		p.dict[key] = len(p.stages) - 1
	}
}

func (p *Pipeline) Marshal() mongo.Pipeline {
	return p.stages
}
