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
	"github.com/qianwj/typed/mongo/util"
)

type DateAdder struct {
	data bson.UnorderedMap
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

func (t *ToStringOptions) Format(format any) *ToStringOptions {
	t.data["format"] = format
	return t
}

func (t *ToStringOptions) Timezone(timezone any) *ToStringOptions {
	t.data["timezone"] = timezone
	return t
}

func (t *ToStringOptions) OnNull(onNull any) *ToStringOptions {
	t.data["onNull"] = onNull
	return t
}

func (t *ToStringOptions) MarshalBSON() ([]byte, error) {
	return t.data.Marshal()
}

type TruncOptions struct {
	data bson.UnorderedMap
}

func NewTrunc(date any, unit timeunit.DateTime) *TruncOptions {
	return &TruncOptions{
		data: bson.M(
			bson.E("date", date),
			bson.E("unit", unit),
		),
	}
}

func (t *TruncOptions) BinSize(binSize any) *TruncOptions {
	t.data["binSize"] = binSize
	return t
}

func (t *TruncOptions) Timezone(zone any) *TruncOptions {
	t.data["timezone"] = zone
	return t
}

func (t *TruncOptions) StartOfWeek(weekday timeunit.Weekday) *TruncOptions {
	t.data["startOfWeek"] = weekday
	return t
}

func (t *TruncOptions) MarshalBSON() ([]byte, error) {
	return t.data.Marshal()
}
