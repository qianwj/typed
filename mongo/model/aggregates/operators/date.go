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

// DateFromParts constructs and returns a Date object given the date's constituent properties.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateFromParts/
func DateFromParts(parts *FromPartsOptions) bson.Entry {
	return bson.E(operator.DateFromParts, parts)
}

// DateFromString converts a date/time string to a date object.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateFromString/
func DateFromString(dateStr *FromStringOptions) bson.Entry {
	return bson.E(operator.DateFromString, dateStr)
}

// DateSubtract Decrements a `Date()` object by a specified number of time units.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateSubtract/
func DateSubtract(subtracter *DateSubtracter) bson.Entry {
	return bson.E(operator.DateSubtract, subtracter)
}

// DateToParts returns a document that contains the constituent parts of a given BSON Date value as individual
// properties. The properties returned are year, month, day, hour, minute, second and millisecond.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateToParts/
func DateToParts(parts *ToPartsOptions) bson.Entry {
	return bson.E(operator.DateToParts, parts)
}

// DateToString converts a date object to a string according to a user-specified format.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateToString/
func DateToString(toString *ToStringOptions) bson.Entry {
	return bson.E(operator.DateToString, toString)
}

// Hour returns the hour portion of a date as a number between 0 and 23.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/hour/
func Hour(date time.Time, timezone ...any) bson.Entry {
	if len(timezone) == 0 {
		return bson.E(operator.Hour, date)
	}
	return bson.E(operator.Hour, bson.M(
		bson.E("date", date),
		bson.E("timezone", timezone[0]),
	))
}

type DateAdder struct {
	data     bson.UnorderedMap
	timezone *string
}

func NewDateAdder(startDate, amount any, unit timeunit.DateTime) *DateAdder {
	return &DateAdder{
		data: bson.M(
			bson.E("startDate", startDate),
			bson.E("unit", unit),
			bson.E("amount", amount),
		),
	}
}

func (d *DateAdder) Timezone(zone any) *DateAdder {
	d.data["timezone"] = zone
	return d
}

func (d *DateAdder) MarshalBSON() ([]byte, error) {
	return d.data.Marshal()
}

type DateDiffer struct {
	data bson.UnorderedMap
}

func NewDateDiffer(startDate, endDate any, unit timeunit.DateTime) *DateDiffer {
	return &DateDiffer{
		data: bson.M(
			bson.E("startDate", startDate),
			bson.E("endDate", endDate),
			bson.E("unit", unit),
		),
	}
}

func (d *DateDiffer) StartOfWeek(weekday timeunit.Weekday) *DateDiffer {
	d.data["startOfWeek"] = weekday
	return d
}

func (d *DateDiffer) Timezone(zone any) *DateDiffer {
	d.data["timezone"] = zone
	return d
}

func (d *DateDiffer) MarshalBSON() ([]byte, error) {
	return d.data.Marshal()
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

type FromPartsOptions struct {
	Year         any `bson:"year,omitempty"`
	ISOWeekYear  any `bson:"isoWeekYear,omitempty"`
	Month        any `bson:"month,omitempty"`
	ISOWeek      any `bson:"isoWeek,omitempty"`
	Day          any `bson:"day,omitempty"`
	ISODayOfWeek any `bson:"isoDayOfWeek,omitempty"`
	Hour         any `bson:"hour,omitempty"`
	Minute       any `bson:"minute,omitempty"`
	Second       any `bson:"second,omitempty"`
	Millisecond  any `bson:"millisecond,omitempty"`
	Timezone     any `bson:"timezone,omitempty"`
}

type FromStringOptions struct {
	data bson.UnorderedMap
}

func NewFromString(dateStrExpr any) *FromStringOptions {
	return &FromStringOptions{
		data: bson.M(bson.E("dateString", dateStrExpr)),
	}
}

func (d *FromStringOptions) Format(format any) *FromStringOptions {
	d.data["format"] = format
	return d
}

func (d *FromStringOptions) Timezone(timezone any) *FromStringOptions {
	d.data["timezone"] = timezone
	return d
}

func (d *FromStringOptions) OnError(onErr any) *FromStringOptions {
	d.data["onError"] = onErr
	return d
}

func (d *FromStringOptions) OnNull(onNull any) *FromStringOptions {
	d.data["onNull"] = onNull
	return d
}

func (d *FromStringOptions) MarshalBSON() ([]byte, error) {
	return d.data.Marshal()
}

type ToPartsOptions struct {
	data bson.UnorderedMap
}

func NewToParts(date any) *ToPartsOptions {
	return &ToPartsOptions{data: bson.M(bson.E("date", date))}
}

func (t *ToPartsOptions) Timezone(zone any) *ToPartsOptions {
	t.data["timezone"] = zone
	return t
}

func (t *ToPartsOptions) ISO8601(iso8601 bool) *ToPartsOptions {
	t.data["iso8601"] = iso8601
	return t
}

func (t *ToPartsOptions) MarshalBSON() ([]byte, error) {
	return t.data.Marshal()
}

type ToStringOptions struct {
	data bson.UnorderedMap
}

func NewToString(dateExpr any) *ToStringOptions {
	return &ToStringOptions{
		data: bson.M(bson.E("date", dateExpr)),
	}
}

func (d *ToStringOptions) Format(format any) *ToStringOptions {
	d.data["format"] = format
	return d
}

func (d *ToStringOptions) Timezone(timezone any) *ToStringOptions {
	d.data["timezone"] = timezone
	return d
}

func (d *ToStringOptions) OnNull(onNull any) *ToStringOptions {
	d.data["onNull"] = onNull
	return d
}
