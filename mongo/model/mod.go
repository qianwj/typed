package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type Bson interface {
	bson.M | bson.D | bson.A | bson.E
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
