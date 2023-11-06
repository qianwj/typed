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
	"github.com/qianwj/typed/mongo/operator"
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

// DateTrunc truncates a date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateTrunc/
func DateTrunc(trunc *TruncOptions) bson.Entry {
	return bson.E(operator.DateTrunc, trunc)
}

// DayOfMonth returns the day of the month for a date as a number between 1 and 31.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dayOfMonth/
func DayOfMonth(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.DayOfMonth, date, timezone...)
}

// DayOfWeek returns the day of the week for a date as a number between 1 (Sunday) and 7 (Saturday).
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dayOfWeek/
func DayOfWeek(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.DayOfWeek, date, timezone...)
}

// DayOfYear returns the day of the year for a date as a number between 1 and 366.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dayOfYear/
func DayOfYear(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.DayOfYear, date, timezone...)
}

// Hour returns the hour portion of a date as a number between 0 and 23.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/hour/
func Hour(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Hour, date, timezone...)
}

// ISODayOfWeek returns the weekday number in ISO 8601 format, ranging from 1 (for Monday) to 7 (for Sunday).
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/isoDayOfWeek/
func ISODayOfWeek(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.ISODayOfWeek, date, timezone...)
}

// ISOWeek returns the week number in ISO 8601 format, ranging from 1 to 53. Week numbers start at 1 with the week
// (Monday through Sunday) that contains the year's first Thursday.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/isoWeek/
func ISOWeek(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.ISOWeek, date, timezone...)
}

// ISOWeekYear returns the year number in ISO 8601 format. The year starts with the Monday of week 1 and ends with the
// Sunday of the last week.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/isoWeekYear/
func ISOWeekYear(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.ISOWeekYear, date, timezone...)
}

// Millisecond returns the millisecond portion of a date as an integer between 0 and 999.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/millisecond/
func Millisecond(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Millisecond, date, timezone...)
}

// Minute returns the minute portion of a date as a number between 0 and 59.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/millisecond/
func Minute(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Minute, date, timezone...)
}

// Month returns the month of a date as a number between 1 and 12.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/month/
func Month(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Month, date, timezone...)
}

// Second returns the second portion of a date as a number between 0 and 59, but can be 60 to account for leap seconds.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/second/
func Second(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Second, date, timezone...)
}

// Week returns the week of the year for a date as a number between 0 and 53.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/week/
func Week(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Week, date, timezone...)
}

// Year returns the year portion of a date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/year/
func Year(date any, timezone ...any) bson.Entry {
	return computeDateWithZone(operator.Year, date, timezone...)
}
