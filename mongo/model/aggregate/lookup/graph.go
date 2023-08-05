package lookup

import (
	"github.com/qianwj/typed/mongo/model/filter"
	"go.mongodb.org/mongo-driver/bson"
)

type GraphLookup struct {
	from                    string         // <collection>,
	startWith               string         // <expression>,
	connectFromField        string         // <string>,
	connectToField          string         // <string>,
	as                      string         // <string>,
	maxDepth                int            // <number>,
	depthField              string         // <string>,
	restrictSearchWithMatch *filter.Filter // <document>
}

func (l *GraphLookup) Marshal() bson.D {
	// complete this function
	return bson.D{}
}
