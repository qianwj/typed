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

package options

import "go.mongodb.org/mongo-driver/mongo/options"

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

func BSON() *BSONOptions {
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
