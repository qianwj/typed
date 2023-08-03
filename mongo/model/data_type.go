package model

type DataType struct {
	number int
	alias  string
}

var (
	DataTypeDouble = &DataType{
		number: 1,
		alias:  "double",
	}
	DataTypeString = &DataType{
		number: 2,
		alias:  "string",
	}
	DataTypeObject = &DataType{
		number: 3,
		alias:  "object",
	}
	DataTypeArray = &DataType{
		number: 4,
		alias:  "array",
	}
	DataTypeBinaryData = &DataType{
		number: 5,
		alias:  "binData",
	}
	DataTypeUndefined = &DataType{
		number: 6,
		alias:  "undefined",
	}
	DataTypeObjectId = &DataType{
		number: 7,
		alias:  "objectId",
	}
	DataTypeBool = &DataType{
		number: 8,
		alias:  "bool",
	}
	DataTypeDate = &DataType{
		number: 9,
		alias:  "date",
	}
	DataTypeNull = &DataType{
		number: 10,
		alias:  "null",
	}
	DataTypeRegex = &DataType{
		number: 11,
		alias:  "regex",
	}
)

func (d *DataType) Order() int {
	return d.number
}
