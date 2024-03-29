package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Projector struct {
	def bson.M
}

func Project(def bson.M) *Projector {
	return &Projector{def: def}
}

func (p *Projector) Tag() {}

func (p *Projector) Marshal() primitive.D {
	return primitive.D{
		{Key: operator.Project, Value: p.def.Marshal()},
	}
}

func (p *Projector) ToMap() primitive.M {
	return primitive.M{
		operator.Project: p.def.ToMap(),
	}
}

type RootReplacer struct {
	expr Expression
}

func ReplaceRoot(expr Expression) *RootReplacer {
	return &RootReplacer{expr: expr}
}

func (r *RootReplacer) Tag() {}

func (r *RootReplacer) Marshal() primitive.D {
	return primitive.D{
		primitive.E{
			Key: operator.ReplaceRoot,
			Value: primitive.M{
				"newRoot": r.expr,
			},
		},
	}
}

func (r *RootReplacer) ToMap() primitive.M {
	return primitive.M{
		operator.ReplaceRoot: bson.M{
			"newRoot": r.expr,
		},
	}
}
