package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewFilter() Filter {
	return Filter{}
}

func TypedFilter(raw bson.M) Filter {
	return Filter(raw)
}

func (f Filter) Eq(field string, value interface{}) Filter {
	f[field] = value
	return f
}

func (f Filter) In(field string, value []interface{}) Filter {
	f[field] = bson.M{
		"$in": value,
	}
	return f
}

func (f Filter) Nin(field string, value []interface{}) Filter {
	f[field] = bson.M{
		"$nin": value,
	}
	return f
}

func (f Filter) Or(conditions bson.A) Filter {
	f["$or"] = conditions
	return f
}

func (f Filter) IgnoreCaseRegex(field string, pattern string) {
	regex := primitive.Regex{
		Pattern: pattern,
		Options: "i",
	}
	f[field] = bson.M{"$regex": regex}
}

func (f Filter) Regex(field string, pattern string) {
	regex := primitive.Regex{Pattern: pattern}
	f[field] = bson.M{"$regex": regex}
}
