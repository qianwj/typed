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

// Arithmetic Expression Operators.
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#arithmetic-expression-operators
const (
	Abs      = "$abs"
	Add      = "$add"
	Ceil     = "$ceil"
	Divide   = "$divide"
	Exp      = "$exp"
	Floor    = "$floor"
	Ln       = "$ln"
	Log      = "$log"
	Log10    = "$log10"
	Multiply = "$multiply"
	Pow      = "$pow"
	Round    = "$round"
	Sqrt     = "$sqrt"
	Subtract = "$subtract"
	Trunc    = "$trunc"
)

// Array Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#array-expression-operators
const (
	ArrayElemAt   = "$arrayElemAt"
	ArrayToObject = "$arrayToObject"
	ConcatArrays  = "$concatArrays"
	Filter        = "$filter"
	FirstN        = "$firstN"
	IndexOfArray  = "$indexOfArray"
	IsArray       = "$isArray"
	LastN         = "$lastN"
	Map           = "$map"
	MaxN          = "$maxN"
	MinN          = "$minN"
	ObjectToArray = "$objectToArray"
	Range         = "$range"
	Reduce        = "$reduce"
	ReverseArray  = "$reverseArray"
	SortArray     = "$sortArray"
	Zip           = "$zip"
)

// Boolean Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#boolean-expression-operators
// $and, $not, $or already defined in other files.

// Comparison Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#comparison-expression-operators
const (
	Cmp = "$cmp"
	Eq  = "$eq"
)

// Conditional Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#conditional-expression-operators
const (
	Cond   = "$cond"
	IfNull = "$ifNull"
	Switch = "$switch"
)

// Custom Aggregation Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#custom-aggregation-expression-operators
const (
	Accumulator = "$accumulator"
	Function    = "$function"
)

// Date Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#date-expression-operators
const (
	DateAdd      = "$dateAdd"
	DateDiff     = "$dateDiff"
	DateSubtract = "$dateSubtract"
	Hour         = "$hour"
)

// Object Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#object-expression-operators
const ()

// String Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#string-expression-operators
const (
	Concat = "$concat"
	Split  = "$split"
)

// Type Expression Operators
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#type-expression-operators
const (
	Convert    = "$convert"
	IsNumber   = "$isNumber"
	ToBool     = "$toBool"
	ToDate     = "$toDate"
	ToDecimal  = "$toDecimal"
	ToDouble   = "$toDouble"
	ToInt      = "$toInt"
	ToLong     = "$toLong"
	ToObjectID = "$toObjectId"
	ToString   = "$toString"
)

// Accumulators ($group, $bucket, $bucketAuto, $setWindowFields)
// See https://www.mongodb.com/docs/manual/reference/aggregation-quick-reference/#accumulators---group---bucket---bucketauto---setwindowfields-
const (
	Avg     = "$avg"
	Bottom  = "$bottom"
	BottomN = "$bottomN"
	First   = "$first"
	Last    = "$last"
	Median  = "$median"
	Sum     = "$sum"
	Top     = "$top"
	TopN    = "$topN"
)

// Aggregation Pipeline Operators
// refer: https://docs.mongodb.com/manual/reference/operator/aggregation/
const (
	DocumentNumber = "$documentNumber"

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
	StrLenBytes  = "$strLenBytes"
	StrLenCP     = "$strLenCP"
	Strcasecmp   = "$strcasecmp"
	Substr       = "$substr"
	SubstrBytes  = "$substrBytes"
	SubstrCP     = "$substrCP"
	ToLower      = "$toLower"
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

	StdDevPop  = "$stdDevPop"
	StdDevSamp = "$stdDevSamp"

	// Variable Expression Operators
	Let = "$let"
)
