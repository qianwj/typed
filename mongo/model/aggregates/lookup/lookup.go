package lookup

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinCondition struct {
	From         string        `bson:"from,omitempty"`         // <collection to join>
	LocalField   string        `bson:"localField,omitempty"`   // <field from the input documents>,
	ForeignField string        `bson:"foreignField,omitempty"` // <field from the documents of the "from" collection>,
	Let          primitive.M   `bson:"let,omitempty"`          // { <var_1>: <expression>, …, <var_n>: <expression> }
	Pipeline     []primitive.D `bson:"pipeline,omitempty"`     // [ <pipeline to run on joined collection> ]
	As           string        `bson:"as,omitempty"`           // <output array field>
}

func New(from, as string) *JoinCondition {
	return &JoinCondition{
		From: from,
		As:   as,
	}
}

func (j *JoinCondition) Join(localField, foreignField string) *JoinCondition {
	j.LocalField = localField
	j.ForeignField = foreignField
	return j
}

func (j *JoinCondition) Pipe(pipeline []primitive.D, let primitive.M) *JoinCondition {
	j.Let = let
	j.Pipeline = pipeline
	return j
}

type GraphJoinCondition struct {
	From                    string          `bson:"from,omitempty"`                    // <collection>,
	StartWith               any             `bson:"startWith,omitempty"`               // <expression>,
	ConnectFromField        string          `bson:"connectFromField,omitempty"`        // <string>,
	ConnectToField          string          `bson:"connectToField,omitempty"`          // <string>,
	As                      string          `bson:"as,omitempty"`                      // <string>,
	MaxDepth                int             `bson:"maxDepth,omitempty"`                // <number>,
	DepthField              string          `bson:"depthField,omitempty"`              // <string>,
	RestrictSearchWithMatch *filters.Filter `bson:"restrictSearchWithMatch,omitempty"` // <document>
}

func NewGraph(from, connectFromField, connectToField, as string, startsWith any) *GraphJoinCondition {
	return &GraphJoinCondition{
		From:             from,
		StartWith:        startsWith,
		ConnectFromField: connectFromField,
		ConnectToField:   connectToField,
		As:               as,
	}
}

func (g *GraphJoinCondition) SetMaxDepth(max int) *GraphJoinCondition {
	g.MaxDepth = max
	return g
}

func (g *GraphJoinCondition) SetDepthField(field string) *GraphJoinCondition {
	g.DepthField = field
	return g
}

func (g *GraphJoinCondition) SetRestrictSearchWithMatch(filter *filters.Filter) *GraphJoinCondition {
	g.RestrictSearchWithMatch = filter
	return g
}
