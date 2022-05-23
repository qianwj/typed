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

func (f Filter) Gt(field string, value interface{}) Filter {
	f[field] = bson.M{
		"$gt": value,
	}
	return f
}

func (f Filter) Gte(field string, value interface{}) Filter {
	f[field] = bson.M{
		"$gte": value,
	}
	return f
}

func (f Filter) Lt(field string, value interface{}) Filter {
	f[field] = bson.M{
		"$lt": value,
	}
	return f
}

func (f Filter) Lte(field string, value interface{}) Filter {
	f[field] = bson.M{
		"$lte": value,
	}
	return f
}

func (f Filter) Ne(field string, value interface{}) Filter {
	f[field] = bson.M{
		"$ne": value,
	}
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

func (f Filter) Or(conditions ...Filter) Filter {
	f["$or"] = conditions
	return f
}

func (f Filter) Not(field string, filter Filter) Filter {
	f[field] = bson.M{
		"$not": filter,
	}
	return f
}

func (f Filter) Exists(field string, exists bool) Filter {
	f[field] = bson.M{
		"$exists": exists,
	}
	return f
}

func (f Filter) Type(field string, fieldType FieldType) Filter {
	f[field] = bson.M{
		"$type": fieldType,
	}
	return f
}

func (f Filter) ElementMatch(filter Filter) Filter {
	f["$elemMatch"] = filter
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
