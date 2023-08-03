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
)

func (d *DataType) Order() int {
	return d.number
}
