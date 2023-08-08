package update

import (
	operator2 "github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (u *Update) CurrentDate(fields ...string) *Update {
	for _, field := range fields {
		u.currentDate[field] = true
	}
	return u
}

func (u *Update) CurrentTimestamp(fields ...string) *Update {
	for _, field := range fields {
		u.currentDate[field] = bson.M{
			operator2.Type: operator2.Timestamp,
		}
	}
	return u
}

func (u *Update) Inc(field string, delta int) *Update {
	u.increment[field] = delta
	return u
}

func (u *Update) Dec(field string, delta int) *Update {
	u.increment[field] = -delta
	return u
}

func (u *Update) MinInt32(field string, cur int32) *Update {
	u.min[field] = cur
	return u
}

func (u *Update) MinInt(field string, cur int) *Update {
	u.min[field] = cur
	return u
}

func (u *Update) MinInt64(field string, cur int64) *Update {
	u.min[field] = cur
	return u
}

func (u *Update) MinFloat32(field string, cur float32) *Update {
	u.min[field] = cur
	return u
}

func (u *Update) MinFloat64(field string, cur float64) *Update {
	u.min[field] = cur
	return u
}

func (u *Update) MinTime(field string, time time.Time) *Update {
	u.min[field] = time
	return u
}

func (u *Update) MaxInt32(field string, cur int32) *Update {
	u.max[field] = cur
	return u
}

func (u *Update) MaxInt(field string, cur int) *Update {
	u.max[field] = cur
	return u
}

func (u *Update) MaxInt64(field string, cur int64) *Update {
	u.max[field] = cur
	return u
}

func (u *Update) MaxFloat32(field string, cur float32) *Update {
	u.max[field] = cur
	return u
}

func (u *Update) MaxFloat64(field string, cur float64) *Update {
	u.max[field] = cur
	return u
}

func (u *Update) MaxTime(field string, time time.Time) *Update {
	u.max[field] = time
	return u
}

func (u *Update) Rename(field, newname string) *Update {
	u.rename[field] = newname
	return u
}

func (u *Update) Set(field string, value any) *Update {
	u.set[field] = value
	return u
}

func (u *Update) SetAll(pairs bson.M) *Update {
	for k, v := range pairs {
		u.set[k] = v
	}
	return u
}

func (u *Update) SetOnInsert(field string, value any) *Update {
	u.setOnInsert[field] = value
	return u
}

func (u *Update) SetOnInsertAll(pairs bson.M) *Update {
	for k, v := range pairs {
		u.setOnInsert[k] = v
	}
	return u
}

func (u *Update) Unset(fields ...string) *Update {
	for _, field := range fields {
		u.unset[field] = ""
	}
	return u
}
