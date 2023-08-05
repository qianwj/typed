package pipe

import "go.mongodb.org/mongo-driver/bson"

type groupId interface {
	tag()
}
type StringGroupId string

func (s *StringGroupId) tag() {}

type DocGroupId bson.D

func (d *DocGroupId) tag() {}
