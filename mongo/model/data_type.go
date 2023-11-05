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

package model

type DataType struct {
	number int8
	alias  string
}

var (
	DataTypeDouble              = &DataType{number: 1, alias: "double"}
	DataTypeString              = &DataType{number: 2, alias: "string"}
	DataTypeObject              = &DataType{number: 3, alias: "object"}
	DataTypeArray               = &DataType{number: 4, alias: "array"}
	DataTypeBinaryData          = &DataType{number: 5, alias: "binData"}
	DataTypeUndefined           = &DataType{number: 6, alias: "undefined"}
	DataTypeObjectId            = &DataType{number: 7, alias: "objectId"}
	DataTypeBool                = &DataType{number: 8, alias: "bool"}
	DataTypeDate                = &DataType{number: 9, alias: "date"}
	DataTypeNull                = &DataType{number: 10, alias: "null"}
	DataTypeRegex               = &DataType{number: 11, alias: "regex"}
	DataTypeDBPointer           = &DataType{number: 12, alias: "dbPointer"}
	DataTypeJavascript          = &DataType{number: 13, alias: "javascript"}
	DataTypeSymbol              = &DataType{number: 14, alias: "symbol"}
	DataTypeJavascriptWithScope = &DataType{number: 15, alias: "javascriptWithScope"}
	DataTypeInt32               = &DataType{number: 16, alias: "int"}
	DataTypeTimestamp           = &DataType{number: 17, alias: "timestamp"}
	DataTypeInt64               = &DataType{number: 18, alias: "long"}
	DataTypeDecimal128          = &DataType{number: 19, alias: "decimal"}
	DataTypeMinKey              = &DataType{number: -1, alias: "minKey"}
	DataTypeMaxKey              = &DataType{number: 127, alias: "maxKey"}
)

func (d *DataType) Order() int8 {
	return d.number
}

func (d *DataType) Alias() string {
	return d.alias
}
