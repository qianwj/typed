package bson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ObjectId primitive.ObjectID

func NewObjectId() ObjectId {
	return ObjectId(primitive.NewObjectID())
}

func (o ObjectId) Cast() any {
	return o
}

func (ObjectId) IsVariable() {}

func (o ObjectId) String() string {
	return "ObjectId(" + primitive.ObjectID(o).Hex() + ")"
}

func (o ObjectId) MarshalJSON() ([]byte, error) {
	return primitive.ObjectID(o).MarshalJSON()
}

type Null primitive.Null

func (n Null) Cast() any {
	return n
}

func (Null) IsVariable() {}

type String string

func (s String) Cast() any {
	return s
}

func (String) IsVariable() {}

func (String) IsExpression() {}

type Int int32

func (i Int) Cast() any {
	return i
}

func (Int) IsVariable() {}

type Long int32

func (l Long) Cast() any {
	return l
}

func (Long) IsVariable() {}

type Bool bool

func (b Bool) Cast() any {
	return b
}

func (Bool) IsVariable() {}

type Float float32

func (f Float) Cast() any {
	return f
}

func (Float) IsVariable() {}

type Double float64

func (d Double) Cast() any {
	return d
}

func (Double) IsVariable() {}

type DateTime time.Time

func Now() DateTime {
	return DateTime(time.Now())
}

func (d DateTime) Cast() any {
	return d
}

func (DateTime) IsVariable() {}

func (d DateTime) String() string {
	return "ISODate(" + time.Time(d).Format(time.RFC3339Nano) + ")"
}
