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

//func NewAggregatePipeline() AggregatePipeline {
//	return AggregatePipeline{}
//}
//
//func NewAggregateGroup() AggregateGroup {
//	return AggregateGroup{}
//}
//
//func (a AggregateGroup) Field(field, from string) AggregateGroup {
//	a[field] = from
//	return a
//}
//
//func (a AggregateGroup) SumFixed(field string, value int) AggregateGroup {
//	a[field] = bson.M{
//		"$sum": value,
//	}
//	return a
//}
//
//func (a AggregatePipeline) Match(filter Filter) AggregatePipeline {
//	return append(a, bson.D{
//		{Key: "$match", Value: filter},
//	})
//}
//
//func (a AggregatePipeline) Group(group AggregateGroup) AggregatePipeline {
//	return append(a, bson.D{
//		{Key: "$group", Value: group},
//	})
//}
//
//func (a AggregatePipeline) Page(page AggregatePage) AggregatePipeline {
//	return append(a,
//		bson.D{{Key: "$skip", Value: page.skip}},
//		bson.D{{Key: "$limit", Value: page.limit}},
//	)
//}
