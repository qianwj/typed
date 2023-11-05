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

import (
	"github.com/qianwj/typed/mongo/model/sorts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	keys []Pair[sorts.Order]
	opts *rawopts.IndexOptions
}

func NewIndex(keys ...Pair[sorts.Order]) *Index {
	return &Index{
		keys: keys,
		opts: rawopts.Index(),
	}
}

func From(model *mongo.IndexModel) *Index {
	keys, idx := model.Keys, Index{keys: make([]Pair[sorts.Order], 0)}
	switch keys.(type) {
	case bson.D:
		for _, e := range keys.(bson.D) {
			idx.keys = append(idx.keys, Pair[sorts.Order]{
				Key:   e.Key,
				Value: e.Value.(sorts.Order),
			})
		}
	}
	idx.opts = model.Options
	return &idx
}

// Background sets value for the Background field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.2.
func (i *Index) Background() *Index {
	i.opts.SetBackground(true)
	return i
}

// ExpireAfterSeconds sets value for the ExpireAfterSeconds field.
func (i *Index) ExpireAfterSeconds(seconds int32) *Index {
	i.opts.SetExpireAfterSeconds(seconds)
	return i
}

// Name sets the value for the Name field.
func (i *Index) Name(name string) *Index {
	i.opts.SetName(name)
	return i
}

// Sparse sets the value of the Sparse field.
func (i *Index) Sparse() *Index {
	i.opts.SetSparse(true)
	return i
}

// StorageEngine sets the value for the StorageEngine field.
func (i *Index) StorageEngine(engine interface{}) *Index {
	i.opts.SetStorageEngine(engine)
	return i
}

// Unique sets the value for the Unique field.
func (i *Index) Unique() *Index {
	i.opts.SetUnique(true)
	return i
}

// Version sets the value for the Version field.
func (i *Index) Version(version int32) *Index {
	i.opts.SetVersion(version)
	return i
}

// DefaultLanguage sets the value for the DefaultLanguage field.
func (i *Index) DefaultLanguage(language string) *Index {
	i.opts.SetDefaultLanguage(language)
	return i
}

// LanguageOverride sets the value of the LanguageOverride field.
func (i *Index) LanguageOverride(override string) *Index {
	i.opts.SetLanguageOverride(override)
	return i
}

// TextVersion sets the value for the TextVersion field.
func (i *Index) TextVersion(version int32) *Index {
	i.opts.SetTextVersion(version)
	return i
}

// Weights sets the value for the Weights field.
func (i *Index) Weights(weights interface{}) *Index {
	i.opts.SetWeights(weights)
	return i
}

// SphereVersion sets the value for the SphereVersion field.
func (i *Index) SphereVersion(version int32) *Index {
	i.opts.SetSphereVersion(version)
	return i
}

// Bits sets the value for the Bits field.
func (i *Index) Bits(bits int32) *Index {
	i.opts.SetBits(bits)
	return i
}

// Max sets the value for the Max field.
func (i *Index) Max(max float64) *Index {
	i.opts.SetMax(max)
	return i
}

// Min sets the value for the Min field.
func (i *Index) Min(min float64) *Index {
	i.opts.SetMin(min)
	return i
}

// BucketSize sets the value for the BucketSize field
func (i *Index) BucketSize(bucketSize int32) *Index {
	i.opts.SetBucketSize(bucketSize)
	return i
}

// PartialFilterExpression sets the value for the PartialFilterExpression field.
func (i *Index) PartialFilterExpression(expression interface{}) *Index {
	i.opts.SetPartialFilterExpression(expression)
	return i
}

// Collation sets the value for the Collation field.
func (i *Index) Collation(collation *rawopts.Collation) *Index {
	i.opts.SetCollation(collation)
	return i
}

// WildcardProjection sets the value for the WildcardProjection field.
func (i *Index) WildcardProjection(wildcardProjection interface{}) *Index {
	i.opts.SetWildcardProjection(wildcardProjection)
	return i
}

// Hidden sets the value for the Hidden field.
func (i *Index) Hidden() *Index {
	i.opts.SetHidden(true)
	return i
}

func (i *Index) Build() mongo.IndexModel {
	keys := bson.D{}
	for _, key := range i.keys {
		keys = append(keys, bson.E{Key: key.Key, Value: key.Value})
	}
	return mongo.IndexModel{
		Keys:    keys,
		Options: i.opts,
	}
}
