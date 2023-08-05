package aggregate

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/operator"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type AddFieldsStage struct {
	value bson.M
}

func AddFields(value bson.M) Stage {
	return &AddFieldsStage{value: value}
}

func (a *AddFieldsStage) Key() string {
	return operator.AddFields
}

func (a *AddFieldsStage) Value() any {
	return a.value
}

type BucketStage struct {
	value bson.M
}

func Bucket(value bson.M) Stage {
	return &BucketStage{value: value}
}

func (b *BucketStage) Key() string {
	return operator.Bucket
}

func (b *BucketStage) Value() any {
	return b.value
}

type CountStage struct {
	value string
}

func Count(field string) Stage {
	return &CountStage{value: field}
}

func (c *CountStage) Key() string {
	return operator.Count
}

func (c *CountStage) Value() any {
	return c.value
}

type MatchStage struct {
	value bson.D
}

func Match(filter *filter.Filter) Stage {
	return &MatchStage{value: filter.Marshal()}
}

func (m *MatchStage) Key() string {
	return operator.Match
}

func (m *MatchStage) Value() any {
	return m.value
}

type SortStage struct {
	value bson.D
}

func Sort(fields ...model.Pair[options.SortOrder]) Stage {
	val := bson.D{}
	for _, field := range fields {
		val = append(val, bson.E{
			Key: field.Key, Value: field.Value,
		})
	}
	return &SortStage{value: val}
}

func (s *SortStage) Key() string {
	return operator.Sort
}

func (s *SortStage) Value() any {
	return s.value
}

type UnsetStage struct {
	value []string
}

func Unset(fields ...string) Stage {
	return &UnsetStage{value: fields}
}

func (u *UnsetStage) Key() string {
	return operator.Unset
}

func (u *UnsetStage) Value() any {
	return u.value
}
