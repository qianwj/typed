package update

import (
	"github.com/qianwj/typed/mongo/model"
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

func MinNumber[N model.Number](field string, cur N) *Update {
	switch cur.(type) {
	case int32:
		return New().MinInt32(field, int32(cur))
	case int:
		return New().MinInt(field, int(cur))
	case int64:
		return New().MinInt64(field, int64(cur))
	case float32:
		return New().MinFloat32(field, float32(cur))
	case float64:
		return New().MinFloat64(field, float64(cur))
	}
	return New()
}

func MinTime(field string, time time.Time) *Update {
	return New().MinTime(field, time)
}

func MaxNumber[N model.Number](field string, cur N) *Update {
	switch cur.(type) {
	case int32:
		return New().MaxInt32(field, int32(cur))
	case int:
		return New().MaxInt(field, int(cur))
	case int64:
		return New().MaxInt64(field, int64(cur))
	case float32:
		return New().MaxFloat32(field, float32(cur))
	case float64:
		return New().MaxFloat64(field, float64(cur))
	}
	return New()
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
