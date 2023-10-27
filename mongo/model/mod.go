package model

import (
	"fmt"
)

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
