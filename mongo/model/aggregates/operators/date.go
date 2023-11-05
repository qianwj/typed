// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package operators

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates/timeunit"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"time"
)

// DateAdd increments a `Date()` object by a specified number of time units.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateAdd/
func DateAdd(adder *DateAdder) bson.Entry {
	return bson.E(operator.DateAdd, adder)
}

// DateDiff returns the difference between two dates.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateDiff/
func DateDiff(differ *DateDiffer) bson.Entry {
	return bson.E(operator.DateDiff, differ)
}

// DateSubtract Decrements a `Date()` object by a specified number of time units.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateSubtract/
func DateSubtract(subtracter *DateSubtracter) bson.Entry {
	return bson.E(operator.DateSubtract, subtracter)
}

// Hour returns the hour portion of a date as a number between 0 and 23.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/hour/
func Hour(date time.Time, timezone ...string) bson.Entry {
	if len(timezone) == 0 {
		return bson.E(operator.Hour, date)
	}
	return bson.E(operator.Hour, bson.M(
		bson.E("date", date),
		bson.E("timezone", timezone[0]),
	))
}

type DateAdder struct {
	startDate any
	unit      timeunit.DateTime
	amount    any
	timezone  *string
}

func NewDateAdder(startDate, amount any, unit timeunit.DateTime) *DateAdder {
	return &DateAdder{
		startDate: startDate,
		amount:    amount,
		unit:      unit,
	}
}

func (d *DateAdder) Timezone(zone string) *DateAdder {
	d.timezone = util.ToPtr(zone)
	return d
}

func (d *DateAdder) MarshalBSON() ([]byte, error) {
	m := bson.M(
		bson.E("startDate", d.startDate),
		bson.E("unit", d.unit),
		bson.E("amount", d.amount),
	)
	if util.IsNonNil(d.timezone) {
		m["timezone"] = d.timezone
	}
	return m.Marshal()
}

type DateDiffer struct {
	startDate   any
	endDate     any
	unit        timeunit.DateTime
	startOfWeek *timeunit.Weekday
	timezone    *string
}

func NewDateDiffer(startDate, endDate any, unit timeunit.DateTime) *DateDiffer {
	return &DateDiffer{
		startDate: startDate,
		endDate:   endDate,
		unit:      unit,
	}
}

func (d *DateDiffer) StartOfWeek(weekday timeunit.Weekday) *DateDiffer {
	d.startOfWeek = util.ToPtr(weekday)
	return d
}

func (d *DateDiffer) Timezone(zone string) *DateDiffer {
	d.timezone = util.ToPtr(zone)
	return d
}

func (d *DateDiffer) MarshalBSON() ([]byte, error) {
	m := bson.M(
		bson.E("startDate", d.startDate),
		bson.E("endDate", d.endDate),
		bson.E("unit", d.unit),
	)
	if util.IsNonNil(d.timezone) {
		m["timezone"] = d.timezone
	}
	if util.IsNonNil(d.startOfWeek) {
		m["startOfWeek"] = d.startOfWeek
	}
	return m.Marshal()
}

type DateSubtracter struct {
	startDate any
	unit      timeunit.DateTime
	amount    any
	timezone  *string
}

func NewDateSubtracter(startDate, amount any, unit timeunit.DateTime) *DateSubtracter {
	return &DateSubtracter{
		startDate: startDate,
		amount:    amount,
		unit:      unit,
	}
}

func (d *DateSubtracter) Timezone(zone string) *DateSubtracter {
	d.timezone = util.ToPtr(zone)
	return d
}

func (d *DateSubtracter) MarshalBSON() ([]byte, error) {
	m := bson.M(
		bson.E("startDate", d.startDate),
		bson.E("unit", d.unit),
		bson.E("amount", d.amount),
	)
	if util.IsNonNil(d.timezone) {
		m["timezone"] = d.timezone
	}
	return m.Marshal()
}
