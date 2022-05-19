package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AggregatePipeline mongo.Pipeline
	AggregateGroup    bson.M
	AggregatePage     struct {
		skip  int64
		limit int64
	}
)

func (a AggregatePipeline) Match(filter Filter) AggregatePipeline {
	return append(a, bson.D{
		{Key: "$match", Value: filter},
	})
}

func (a AggregatePipeline) Group(group AggregateGroup) AggregatePipeline {
	return append(a, bson.D{
		{Key: "$group", Value: group},
	})
}

func (a AggregatePipeline) Page(page AggregatePage) AggregatePipeline {
	return append(a,
		bson.D{{Key: "$skip", Value: page.skip}},
		bson.D{{Key: "$limit", Value: page.limit}},
	)
}
