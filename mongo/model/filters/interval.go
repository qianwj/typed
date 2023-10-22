package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
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

func (i *Interval) query() bson.M {
	switch i.mode {
	case open:
		return bson.M{operator.Gt: i.left, operator.Lt: i.right}
	case leftHalfOpen:
		return bson.M{operator.Gte: i.left, operator.Lt: i.right}
	case rightHalfOpen:
		return bson.M{operator.Gt: i.left, operator.Lte: i.right}
	case closed:
		return bson.M{operator.Gte: i.left, operator.Lte: i.right}
	case leftUnbound:
		return bson.M{operator.Gt: i.left}
	case leftHalfUnbound:
		return bson.M{operator.Gt: i.left}
	case rightUnbound:
		return bson.M{operator.Lt: i.right}
	case rightHalfUnbound:
		return bson.M{operator.Lte: i.right}
	}
	return bson.M{}
}
