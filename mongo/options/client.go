package options

import (
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	CompressorType          string
	CredentialAuthMechanism string
)

const (
	Snappy      CompressorType          = "snappy"
	ZLIB        CompressorType          = "zlib"
	ZSTD        CompressorType          = "zstd"
	MongoDBCR   CredentialAuthMechanism = "MONGODB-CR"
	Plain       CredentialAuthMechanism = "PLAIN"
	GSSAPI      CredentialAuthMechanism = "GSSAPI"
	MongoDBX509 CredentialAuthMechanism = "MONGODB-X509"
	MongoDBAWS  CredentialAuthMechanism = "MONGODB-AWS"
)

const (
	AuthMechanismServiceNameKey          = "SERVICE_NAME"
	AuthMechanismCanonicalizeHostNameKey = "CANONICALIZE_HOST_NAME"
	AuthMechanismServiceRealmKey         = "SERVICE_REALM"
	AuthMechanismServiceHostKey          = "SERVICE_HOST"
	AuthMechanismAWSSessionTokenKey      = "AWS_SESSION_TOKEN"
)

func (c CompressorType) String() string {
	return string(c)
}

type Credential struct {
	authMechanism           CredentialAuthMechanism
	authMechanismProperties map[string]string
	authSource              string
	username                string
	password                *string
}

func NewCredential() *Credential {
	return &Credential{}
}

func (c *Credential) AuthMechanism(mechanism CredentialAuthMechanism) *Credential {
	c.authMechanism = mechanism
	return c
}

func (c *Credential) AuthMechanismServiceName(name string) *Credential {
	c.authMechanismProperties[AuthMechanismServiceNameKey] = name
	return c
}

func (c *Credential) AuthMechanismCanonicalizeHostName(name string) *Credential {
	c.authMechanismProperties[AuthMechanismCanonicalizeHostNameKey] = name
	return c
}

func (c *Credential) AuthMechanismServiceRealm(realm string) *Credential {
	c.authMechanismProperties[AuthMechanismServiceRealmKey] = realm
	return c
}

func (c *Credential) AuthMechanismServiceHost(host string) *Credential {
	c.authMechanismProperties[AuthMechanismServiceHostKey] = host
	return c
}

func (c *Credential) AuthMechanismAWSSessionToken(token string) *Credential {
	c.authMechanismProperties[AuthMechanismAWSSessionTokenKey] = token
	return c
}

func (c *Credential) AuthSource(authSource string) *Credential {
	c.authSource = authSource
	return c
}

func (c *Credential) Username(name string) *Credential {
	c.username = name
	return c
}

func (c *Credential) Password(pwd string) *Credential {
	c.password = util.ToPtr(pwd)
	return c
}

func (c *Credential) Auth() options.Credential {
	if c == nil {
		return options.Credential{}
	}
	auth := options.Credential{
		AuthMechanism:           string(c.authMechanism),
		AuthMechanismProperties: c.authMechanismProperties,
		AuthSource:              c.authSource,
		Username:                c.username,
	}
	if util.NonNil(c.password) {
		auth.Password = util.FromPtr(c.password)
		auth.PasswordSet = true
	}
	return auth
}

type BSONOptions struct {
	useJSONStructTags       bool
	errorOnInlineDuplicates bool
	intMinSize              bool
	nilMapAsEmpty           bool
	nilSliceAsEmpty         bool
	nilByteSliceAsEmpty     bool
	omitZeroStruct          bool
	stringifyMapKeysWithFmt bool
	allowTruncatingDoubles  bool
	binaryAsSlice           bool
	defaultDocumentD        bool
	defaultDocumentM        bool
	useLocalTimeZone        bool
	zeroMaps                bool
	zeroStructs             bool
}

func NewBSONOptions() *BSONOptions {
	return &BSONOptions{}
}

// UseJSONStructTags causes the driver to fall back to using the "json"
// struct tag if a "bson" struct tag is not specified.
func (b *BSONOptions) UseJSONStructTags() *BSONOptions {
	b.useJSONStructTags = true
	return b
}

// ErrorOnInlineDuplicates causes the driver to return an error if there is
// a duplicate field in the marshaled BSON when the "inline" struct tag
// option is set.
func (b *BSONOptions) ErrorOnInlineDuplicates() *BSONOptions {
	b.errorOnInlineDuplicates = true
	return b
}

// IntMinSize causes the driver to marshal Go integer values (int, int8,
// int16, int32, int64, uint, uint8, uint16, uint32, or uint64) as the
// minimum BSON int size (either 32 or 64 bits) that can represent the
// integer value.
func (b *BSONOptions) IntMinSize() *BSONOptions {
	b.intMinSize = true
	return b
}

// NilMapAsEmpty causes the driver to marshal nil Go maps as empty BSON
// documents instead of BSON null.
//
// Empty BSON documents take up slightly more space than BSON null, but
// preserve the ability to use document update operations like "$set" that
// do not work on BSON null.
func (b *BSONOptions) NilMapAsEmpty() *BSONOptions {
	b.nilMapAsEmpty = true
	return b
}

// NilSliceAsEmpty causes the driver to marshal nil Go slices as empty BSON
// arrays instead of BSON null.
//
// Empty BSON arrays take up slightly more space than BSON null, but
// preserve the ability to use array update operations like "$push" or
// "$addToSet" that do not work on BSON null.
func (b *BSONOptions) NilSliceAsEmpty() *BSONOptions {
	b.nilSliceAsEmpty = true
	return b
}

// NilByteSliceAsEmpty causes the driver to marshal nil Go byte slices as
// empty BSON binary values instead of BSON null.
func (b *BSONOptions) NilByteSliceAsEmpty() *BSONOptions {
	b.nilByteSliceAsEmpty = true
	return b
}

// OmitZeroStruct causes the driver to consider the zero value for a struct
// (e.g. MyStruct{}) as empty and omit it from the marshaled BSON when the
// "omitempty" struct tag option is set.
func (b *BSONOptions) OmitZeroStruct() *BSONOptions {
	b.omitZeroStruct = true
	return b
}

// StringifyMapKeysWithFmt causes the driver to convert Go map keys to BSON
// document field name strings using fmt.Sprint instead of the default
// string conversion logic.
func (b *BSONOptions) StringifyMapKeysWithFmt() *BSONOptions {
	b.stringifyMapKeysWithFmt = true
	return b
}

// AllowTruncatingDoubles causes the driver to truncate the fractional part
// of BSON "double" values when attempting to unmarshal them into a Go
// integer (int, int8, int16, int32, or int64) struct field. The truncation
// logic does not apply to BSON "decimal128" values.
func (b *BSONOptions) AllowTruncatingDoubles() *BSONOptions {
	b.allowTruncatingDoubles = true
	return b
}

// BinaryAsSlice causes the driver to unmarshal BSON binary field values
// that are the "Generic" or "Old" BSON binary subtype as a Go byte slice
// instead of a primitive.Binary.
func (b *BSONOptions) BinaryAsSlice() *BSONOptions {
	b.binaryAsSlice = true
	return b
}

// DefaultDocumentD causes the driver to always unmarshal documents into the
// primitive.D type. This behavior is restricted to data typed as
// "interface{}" or "map[string]interface{}".
func (b *BSONOptions) DefaultDocumentD() *BSONOptions {
	b.defaultDocumentD = true
	return b
}

// DefaultDocumentM causes the driver to always unmarshal documents into the
// primitive.M type. This behavior is restricted to data typed as
// "interface{}" or "map[string]interface{}".
func (b *BSONOptions) DefaultDocumentM() *BSONOptions {
	b.defaultDocumentM = true
	return b
}

// UseLocalTimeZone causes the driver to unmarshal time.Time values in the
// local timezone instead of the UTC timezone.
func (b *BSONOptions) UseLocalTimeZone() *BSONOptions {
	b.useLocalTimeZone = true
	return b
}

// ZeroMaps causes the driver to delete any existing values from Go maps in
// the destination value before unmarshaling BSON documents into them.
func (b *BSONOptions) ZeroMaps() *BSONOptions {
	b.zeroMaps = true
	return b
}

// ZeroStructs causes the driver to delete any existing values from Go
// structs in the destination value before unmarshaling BSON documents into
// them.
func (b *BSONOptions) ZeroStructs() *BSONOptions {
	b.zeroStructs = true
	return b
}

func (b *BSONOptions) Raw() *options.BSONOptions {
	if b == nil {
		return nil
	}
	return &options.BSONOptions{
		UseJSONStructTags:       b.useJSONStructTags,
		ErrorOnInlineDuplicates: b.errorOnInlineDuplicates,
		IntMinSize:              b.intMinSize,
		NilMapAsEmpty:           b.nilMapAsEmpty,
		NilSliceAsEmpty:         b.nilSliceAsEmpty,
		NilByteSliceAsEmpty:     b.nilByteSliceAsEmpty,
		OmitZeroStruct:          b.omitZeroStruct,
		StringifyMapKeysWithFmt: b.stringifyMapKeysWithFmt,
		AllowTruncatingDoubles:  b.allowTruncatingDoubles,
		BinaryAsSlice:           b.binaryAsSlice,
		DefaultDocumentD:        b.defaultDocumentD,
		DefaultDocumentM:        b.defaultDocumentM,
		UseLocalTimeZone:        b.useLocalTimeZone,
		ZeroMaps:                b.zeroMaps,
		ZeroStructs:             b.zeroStructs,
	}
}
