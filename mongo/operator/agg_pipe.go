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

package operator

// Aggregation Pipeline Operators
// refer: https://docs.mongodb.com/manual/reference/operator/aggregation/
const (
	Abs            = "$abs"
	Accumulator    = "$accumulator"
	Add            = "$add"
	Avg            = "$avg"
	Bottom         = "$bottom"
	BottomN        = "$bottomN"
	Ceil           = "$ceil"
	Cmp            = "$cmp"
	Concat         = "$concat"
	ConcatArrays   = "$concatArrays"
	Cond           = "$cond"
	Convert        = "$convert"
	DateAdd        = "$dateAdd"
	DateDiff       = "$dateDiff"
	DateSubtract   = "$dateSubtract"
	Divide         = "$divide"
	DocumentNumber = "$documentNumber"
	Filter         = "$filter"
	First          = "$first"
	FirstN         = "$firstN"
	Floor          = "$floor"
	Function       = "$function"
	Hour           = "$hour"
	IfNull         = "$ifNull"
	Last           = "$last"
	LastN          = "$lastN"
	Map            = "$map"
	MaxN           = "$maxN"
	Median         = "$median"
	MinN           = "$minN"
	Multiply       = "$multiply"
	ObjectToArray  = "$objectToArray"
	Reduce         = "$reduce"
	Subtract       = "$subtract"
	Sum            = "$sum"
	Top            = "$top"
	TopN           = "$topN"

	Exp   = "$exp"
	Ln    = "$ln"
	Log   = "$log"
	Log10 = "$log10"
	Pow   = "$pow"
	Round = "$round"
	Sqrt  = "$sqrt"
	Trunc = "$trunc"

	// Array Expression Operators
	ArrayElemAt   = "$arrayElemAt"
	ArrayToObject = "$arrayToObject"
	IndexOfArray  = "$indexOfArray"
	IsArray       = "$isArray"
	Range         = "$range"
	ReverseArray  = "$reverseArray"
	Zip           = "$zip"

	Switch = "$switch"

	// Data Size Operators
	BinarySize = "$binarySize"
	BsonSize   = "$bsonSize"

	// Date Expression Operators
	DateFromParts  = "$dateFromParts"
	DateFromString = "$dateFromString"
	DateToParts    = "$dateToParts"
	DateToString   = "$dateToString"
	DayOfMonth     = "$dayOfMonth"
	DayOfWeek      = "$dayOfWeek"
	DayOfYear      = "$dayOfYear"
	IsoDayOfWeek   = "$isoDayOfWeek"
	IsoWeek        = "$isoWeek"
	IsoWeekYear    = "$isoWeekYear"
	Millisecond    = "$millisecond"
	Minute         = "$minute"
	Month          = "$month"
	Second         = "$second"
	ToDate         = "$toDate"
	Week           = "$week"
	Year           = "$year"

	// Literal Expression Operator
	Literal = "$literal"

	// Object Expression Operators
	MergeObjects = "$mergeObjects"

	// Set Expression Operators
	AllElementsTrue = "$allElementsTrue"
	AnyElementTrue  = "$anyElementTrue"
	SetDifference   = "$setDifference"
	SetEquals       = "$setEquals"
	SetIntersection = "$setIntersection"
	SetIsSubset     = "$setIsSubset"
	SetUnion        = "$setUnion"

	// String Expression Operators
	IndexOfBytes = "$indexOfBytes"
	IndexOfCP    = "$indexOfCP"
	Ltrim        = "$ltrim"
	RegexFind    = "$regexFind"
	RegexFindAll = "$regexFindAll"
	RegexMatch   = "$regexMatch"
	Rtrim        = "$rtrim"
	Split        = "$split"
	StrLenBytes  = "$strLenBytes"
	StrLenCP     = "$strLenCP"
	Strcasecmp   = "$strcasecmp"
	Substr       = "$substr"
	SubstrBytes  = "$substrBytes"
	SubstrCP     = "$substrCP"
	ToLower      = "$toLower"
	ToString     = "$toString"
	Trim         = "$trim"
	ToUpper      = "$toUpper"
	ReplaceOne   = "$replaceOne"
	ReplaceAll   = "$replaceAll"

	// Trigonometry Expression Operators
	Sin              = "$sin"
	Cos              = "$cos"
	Tan              = "$tan"
	Asin             = "$asin"
	Acos             = "$acos"
	Atan             = "$atan"
	Atan2            = "$atan2"
	Asinh            = "$asinh"
	Acosh            = "$acosh"
	Atanh            = "$atanh"
	DegreesToRadians = "$degreesToRadians"
	RadiansToDegrees = "$radiansToDegrees"

	// Type Expression Operators

	ToBool     = "$toBool"
	ToDecimal  = "$toDecimal"
	ToDouble   = "$toDouble"
	ToInt      = "$toInt"
	ToLong     = "$toLong"
	ToObjectID = "$toObjectId"
	IsNumber   = "$isNumber"

	StdDevPop  = "$stdDevPop"
	StdDevSamp = "$stdDevSamp"

	// Variable Expression Operators
	Let = "$let"
)
