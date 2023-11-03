package filters

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
)

type interval int8

const (
	open             interval = iota + 1 // (a, b)
	leftHalfOpen                         // [a, b)
	rightHalfOpen                        // (a, b]
	closed                               // [a, b]
	leftUnbound                          // (a, +)
	leftHalfUnbound                      // [a, +)
	rightUnbound                         // (-, b)
	rightHalfUnbound                     // (-, b]
)

type Interval struct {
	left  any
	right any
	mode  interval
}

func OpenInterval(left, right any) *Interval {
	return &Interval{
		left:  left,
		right: right,
		mode:  open,
	}
}

func LeftHalfOpen(left, right any) *Interval {
	return &Interval{
		left:  left,
		right: right,
		mode:  leftHalfOpen,
	}
}

func RightHalfOpen(left, right any) *Interval {
	return &Interval{
		left:  left,
		right: right,
		mode:  rightHalfOpen,
	}
}

func ClosedInterval(left, right any) *Interval {
	return &Interval{
		left:  left,
		right: right,
		mode:  closed,
	}
}

func LeftUnboundInterval(left any) *Interval {
	return &Interval{
		left: left,
		mode: leftUnbound,
	}
}

func LeftHalfUnboundInterval(left any) *Interval {
	return &Interval{
		left: left,
		mode: leftHalfUnbound,
	}
}

func RightUnboundInterval(right any) *Interval {
	return &Interval{
		right: right,
		mode:  rightUnbound,
	}
}

func RightHalfUnboundInterval(right any) *Interval {
	return &Interval{
		right: right,
		mode:  rightHalfUnbound,
	}
}

func (i *Interval) query() *bson.Map {
	switch i.mode {
	case open:
		return bson.NewMap(bson.Entry{Key: operator.Gt, Value: i.left}, bson.Entry{Key: operator.Lt, Value: i.right})
	case leftHalfOpen:
		return bson.NewMap(bson.Entry{Key: operator.Gte, Value: i.left}, bson.Entry{Key: operator.Lt, Value: i.right})
	case rightHalfOpen:
		return bson.NewMap(bson.Entry{Key: operator.Gt, Value: i.left}, bson.Entry{Key: operator.Lte, Value: i.right})
	case closed:
		return bson.NewMap(bson.Entry{Key: operator.Gte, Value: i.left}, bson.Entry{Key: operator.Lte, Value: i.right})
	case leftUnbound:
		return bson.NewMap(bson.Entry{Key: operator.Gt, Value: i.left})
	case leftHalfUnbound:
		return bson.NewMap(bson.Entry{Key: operator.Gte, Value: i.left})
	case rightUnbound:
		return bson.NewMap(bson.Entry{Key: operator.Lt, Value: i.right})
	case rightHalfUnbound:
		return bson.NewMap(bson.Entry{Key: operator.Lte, Value: i.right})
	}
	return bson.NewMap()
}
