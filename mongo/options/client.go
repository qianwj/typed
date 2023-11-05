package options

import "go.mongodb.org/mongo-driver/mongo/options"

type ClientCompressorType string

const (
	Snappy ClientCompressorType = "snappy"
	ZLIB   ClientCompressorType = "zlib"
	ZSTD   ClientCompressorType = "zstd"
)

func (c ClientCompressorType) String() string {
	return string(c)
}

type (
	Credential  options.Credential
	BSONOptions options.BSONOptions
)
