package model

import "go.mongodb.org/mongo-driver/bson"

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
