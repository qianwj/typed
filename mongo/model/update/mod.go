package update

import (
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type PopType int8

const (
	PopFirst PopType = -1
	PopLast  PopType = 1
)

func CurrentDate(fields ...string) *Update {
	return New().CurrentDate(fields...)
}

func CurrentTimestamp(fields ...string) *Update {
	return New().CurrentTimestamp(fields...)
}

func Inc(field string, delta int) *Update {
	return New().Inc(field, delta)
}

func Dec(field string, delta int) *Update {
	return New().Dec(field, delta)
}

func MinInt32(field string, cur int32) *Update {
	return New().MinInt32(field, cur)
}

func MinInt(field string, cur int) *Update {
	return New().MinInt(field, cur)
}

func MinInt64(field string, cur int64) *Update {
	return New().MinInt64(field, cur)
}

func MinFloat32(field string, cur float32) *Update {
	return New().MinFloat32(field, cur)
}

func MinFloat64(field string, cur float64) *Update {
	return New().MinFloat64(field, cur)
}

func MinTime(field string, time time.Time) *Update {
	return New().MinTime(field, time)
}

func MaxInt32(field string, cur int32) *Update {
	return New().MaxInt32(field, cur)
}

func MaxInt(field string, cur int) *Update {
	return New().MaxInt(field, cur)
}

func MaxInt64(field string, cur int64) *Update {
	return New().MaxInt64(field, cur)
}

func MaxFloat32(field string, cur float32) *Update {
	return New().MaxFloat32(field, cur)
}

func MaxFloat64(field string, cur float64) *Update {
	return New().MaxFloat64(field, cur)
}

func MaxTime(field string, time time.Time) *Update {
	return New().MaxTime(field, time)
}

func Rename(field, newname string) *Update {
	return New().Rename(field, newname)
}

func Set(field string, value any) *Update {
	return New().Set(field, value)
}

func SetAll(pairs bson.M) *Update {
	return New().SetAll(pairs)
}

func SetOnInsert(field string, value any) *Update {
	return New().SetOnInsert(field, value)
}

func SetOnInsertAll(pairs bson.M) *Update {
	return New().SetOnInsertAll(pairs)
}

func Unset(fields ...string) *Update {
	return New().Unset(fields...)
}

func AddToSet(field string, value any) *Update {
	return New().AddToSet(field, value)
}

func AddEachToSet(field string, value []any) *Update {
	return New().AddEachToSet(field, value)
}

func Push(field string, value any) *Update {
	return New().Push(field, value)
}

func PushEach(field string, value []any) *Update {
	return New().PushEach(field, value)
}

func InsertAll(field string, value []any, index int) *Update {
	return New().InsertAll(field, value, index)
}

func PushLimited(field string, value []string, limit int) *Update {
	return New().PushLimited(field, value, limit)
}

func PushSorted(field string, value []string, sort options.SortOptions) *Update {
	return New().PushSorted(field, value, sort)
}

func Pop(field string, pType PopType) *Update {
	return New().Pop(field, pType)
}

func Pull(field string, value any) *Update {
	return New().Pull(field, value)
}

func PullConditioned(condition *filter.Filter) *Update {
	return New().PullConditioned(condition)
}

func PullAll(field string, value []any) *Update {
	return New().PullAll(field, value)
}
