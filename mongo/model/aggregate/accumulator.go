package aggregate

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Accumulate struct {
	init           string  // <code>,
	initArgs       bson.A  // <array expression>,        // Optional
	accumulate     string  // <code>,
	accumulateArgs bson.A  // <array expression>,
	merge          string  // <code>,
	finalize       *string // <code>,                    // Optional
	lang           string  // <string>
}

func (a *Accumulate) Marshal() bson.D {
	// complete this function
	return bson.D{}
}
