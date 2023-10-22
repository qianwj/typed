/*
 Copyright 2023 The typed Authors.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package filters

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
