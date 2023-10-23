package update

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Update struct {
	currentDate bson.M
	increment   bson.M
	min         bson.M
	max         bson.M
	mul         bson.M
	rename      bson.M
	set         bson.M
	setOnInsert bson.M
	unset       bson.M
	push        bson.M
	addToSet    bson.M
	pop         bson.M
	pull        bson.M
	pullAll     bson.M
	bit         bson.M
}

func New() *Update {
	return &Update{
		currentDate: bson.M{},
		increment:   bson.M{},
		min:         bson.M{},
		max:         bson.M{},
		mul:         bson.M{},
		rename:      bson.M{},
		set:         bson.M{},
		setOnInsert: bson.M{},
		unset:       bson.M{},
		push:        bson.M{},
		addToSet:    bson.M{},
		pop:         bson.M{},
		pull:        bson.M{},
		pullAll:     bson.M{},
		bit:         bson.M{},
	}
}

func (u *Update) Document() bson.M {
	res := bson.M{}
	if len(u.currentDate) > 0 {
		res[operator.CurrentDate] = u.currentDate
	}
	if len(u.increment) > 0 {
		res[operator.Inc] = u.increment
	}
	if len(u.min) > 0 {
		res[operator.Min] = u.min
	}
	if len(u.max) > 0 {
		res[operator.Max] = u.max
	}
	if len(u.mul) > 0 {
		res[operator.Mul] = u.mul
	}
	if len(u.rename) > 0 {
		res[operator.Rename] = u.rename
	}
	if len(u.set) > 0 {
		res[operator.Set] = u.set
	}
	if len(u.setOnInsert) > 0 {
		res[operator.SetOnInsert] = u.setOnInsert
	}
	if len(u.unset) > 0 {
		res[operator.Unset] = u.unset
	}
	if len(u.push) > 0 {
		res[operator.Push] = u.push
	}
	if len(u.addToSet) > 0 {
		res[operator.AddToSet] = u.addToSet
	}
	if len(u.pop) > 0 {
		res[operator.Pop] = u.pop
	}
	if len(u.pull) > 0 {
		res[operator.Pull] = u.pull
	}
	if len(u.pullAll) > 0 {
		res[operator.PullAll] = u.pullAll
	}
	if len(u.bit) > 0 {
		res[operator.Bit] = u.bit
	}
	return res
}

func (u *Update) MarshalBSON() ([]byte, error) {
	return bson.Marshal(u.Document())
}
