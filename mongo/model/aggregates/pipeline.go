package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	stages mongo.Pipeline
	dict   map[string]int
}

func New() *Pipeline {
	return &Pipeline{
		stages: make(mongo.Pipeline, 0),
		dict:   make(map[string]int),
	}
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

func (p *Pipeline) get(key string) (any, bool) {
	idx, exist := p.dict[key]
	if !exist {
		return nil, false
	}
	return p.stages[idx], true
}

func (p *Pipeline) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.stages)
}
