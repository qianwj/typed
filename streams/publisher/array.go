package publisher

import (
	"go.uber.org/atomic"
	"math"
	"streams"
	"streams/item"
)

type ArraySubscription struct {
	actual    streams.Subscriber
	val       []item.Item
	idx       int
	cancelled atomic.Bool
	requested atomic.Int64
}

func NewArraySubscription(actual streams.Subscriber, val []item.Item) *ArraySubscription {
	return &ArraySubscription{
		actual: actual,
		val:    val,
	}
}

func (a *ArraySubscription) Request(n int64) {
	if n == math.MaxInt64 {
		a.fast()
	} else {
		a.slow(n)
	}
}

func (a *ArraySubscription) Cancel() {
	a.cancelled.CompareAndSwap(false, true)
}

func (a *ArraySubscription) Poll() *item.Item {
	i, val := a.idx, a.val
	if i != len(val) {
		it := val[i]
		a.idx = i + 1
		return &it
	}
	return nil
}

func (a *ArraySubscription) slow(n int64) {
	arr, size := a.val, len(a.val)
	subscriber := a.actual
	i, e := a.idx, int64(0)
	for {
		if a.cancelled.Load() {
			return
		}
		for i != size && e != n {
			it := arr[i]
			if it.IsErr() {
				subscriber.OnError(it.Err())
			} else {
				subscriber.OnNext(it)
			}
			if a.cancelled.Load() {
				return
			}
			i++
			e++
		}
		if i == size {
			subscriber.OnComplete()
			return
		}
		n = a.requested.Load()
		if n == e {
			a.idx = i
			n = a.requested.Sub(e)
			if n == 0 {
				return
			}
			e = int64(0)
		}
	}
}

func (a *ArraySubscription) fast() {
	arr, size := a.val, len(a.val)
	subscriber := a.actual
	for i := a.idx; i != size; i++ {
		if a.cancelled.Load() {
			return
		}
		it := arr[i]
		if it.IsErr() {
			subscriber.OnError(it.Err())
		} else {
			subscriber.OnNext(it)
		}
	}
	if a.cancelled.Load() {
		return
	}
	subscriber.OnComplete()
}
