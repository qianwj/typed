package publisher

import (
	"streams"
	"sync"
)

type ScalarSubscription struct {
	actual streams.Subscriber
	once   sync.Once
	val    any
}

func newScalarSubscription(actual streams.Subscriber, val any) *ScalarSubscription {
	return &ScalarSubscription{
		actual: actual,
		val:    val,
	}
}

func (s *ScalarSubscription) Request(n int64) {

}

func (s *ScalarSubscription) Cancel() {

}
