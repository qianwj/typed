package aggregates

import "github.com/qianwj/typed/mongo/bson"

type GroupId interface {
	isGroupId()
}

type StringGroupId string

func (s StringGroupId) isGroupId() {}

type ComplexGroupId bson.D

func (c ComplexGroupId) isGroupId() {}
