package model

import (
	"fmt"
	"github.com/qianwj/typed/mongo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ID is a constraint that defines the type of the `_id` field in mongo documents.
type ID interface {
	~string | bson.Number | primitive.ObjectID
}

// Doc is an interface that defines the type of the mongo document.
// If you use `TypedCollection`, your document type must implement this interface.
type Doc[I ID] interface {
	GetID() I
}

type Pair[V any] struct {
	Key   string
	Value V
}

type Addr struct {
	Host string
	Port int
}

func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
