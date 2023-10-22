package filters

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	rawbson "go.mongodb.org/mongo-driver/bson"
)

type Filter struct {
	data *bson.Map
}

func From(m rawbson.M) *Filter {
	return &Filter{data: bson.FromM(m)}
}

func Cast(d rawbson.D) *Filter {
	return &Filter{data: bson.FromD(d)}
}

func New() *Filter {
	return &Filter{data: bson.NewMap()}
}

func (f *Filter) Append(other *Filter) *Filter {
	for _, entry := range other.data.Entries() {
		f.data.Put(entry.Key, entry.Value)
	}
	return f
}

func (f *Filter) Get(key string) (any, bool) {
	return f.data.Get(key)
}

func (f *Filter) ToMap() map[string]any {
	return f.data.ToMap()
}

func (f *Filter) Raw() rawbson.D {
	return f.data.Raw()
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	return f.data.MarshalJSON()
}

func (f *Filter) UnmarshalJSON(bytes []byte) error {
	if f.data == nil {
		f.data = bson.NewMap()
	}
	return f.data.UnmarshalJSON(bytes)
}

func (f *Filter) MarshalBSON() ([]byte, error) {
	return f.data.MarshalBSON()
}

func (f *Filter) UnmarshalBSON(bytes []byte) error {
	if f.data == nil {
		f.data = bson.NewMap()
	}
	return f.data.UnmarshalBSON(bytes)
}

func (f *Filter) putAsArray(key string, others ...*Filter) *Filter {
	f.data.PutAsArray(key, util.Map(others, func(it *Filter) *bson.Map {
		return it.data
	})...)
	return f
}

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

func (i *Interval) query() rawbson.M {
	switch i.mode {
	case open:
		return rawbson.M{operator.Gt: i.left, operator.Lt: i.right}
	case leftHalfOpen:
		return rawbson.M{operator.Gte: i.left, operator.Lt: i.right}
	case rightHalfOpen:
		return rawbson.M{operator.Gt: i.left, operator.Lte: i.right}
	case closed:
		return rawbson.M{operator.Gte: i.left, operator.Lte: i.right}
	case leftUnbound:
		return rawbson.M{operator.Gt: i.left}
	case leftHalfUnbound:
		return rawbson.M{operator.Gt: i.left}
	case rightUnbound:
		return rawbson.M{operator.Lt: i.right}
	case rightHalfUnbound:
		return rawbson.M{operator.Lte: i.right}
	}
	return rawbson.M{}
}
