package bucket

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bucket struct {
	GroupBy    any         `bson:"groupBy,omitempty"`
	Boundaries []any       `bson:"boundaries,omitempty"`
	Default    string      `bson:"default,omitempty"`
	Out        primitive.M `bson:"out,omitempty"`
}

func New(groupBy any, boundaries []any) *Bucket {
	return &Bucket{GroupBy: groupBy, Boundaries: boundaries}
}

func (b *Bucket) SetDefault(def string) *Bucket {
	if b != nil {
		b.Default = def
	}
	return b
}

func (b *Bucket) SetOut(out primitive.M) *Bucket {
	if b != nil {
		b.Out = out
	}
	return b
}
