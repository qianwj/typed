package operator

const (
	// Comparison
	Eq  = "$eq"
	Gt  = "$gt"
	Gte = "$gte"
	In  = "$in"
	Lt  = "$lt"
	Lte = "$lte"
	Ne  = "$ne"
	Nin = "$nin"

	// Logical
	And = "$and"
	Not = "$not"
	Nor = "$nor"
	Or  = "$or"

	// Element
	Exists = "$exists"
	Type   = "$type"

	// Evaluation
	Expr               = "$expr"
	JsonSchema         = "$jsonSchema"
	Mod                = "$mod"
	Regex              = "$regex"
	Text               = "$text"
	Where              = "$where"
	Search             = "$search"
	Language           = "$language"
	CaseSensitive      = "$caseSensitive"
	DiacriticSensitive = "$diacriticSensitive"

	// Geo spatial
	GeoIntersects = "$geoIntersects"
	GeoWithin     = "$geoWithin"
	Near          = "$near"
	NearSphere    = "$nearSphere"

	// Array
	All       = "$all"
	ElemMatch = "$elemMatch"
	Size      = "$size"

	// Bitwise
	BitsAllClear = "$bitsAllClear"
	BitsAllSet   = "$bitsAllSet"
	BitsAnyClear = "$bitsAnyClear"
	BitsAnySet   = "$bitsAnySet"

	// Comments
	Comment = "$comment"

	// Projection operators
	Dollar = "$"
	Meta   = "$meta"
	Slice  = "$slice"
)
