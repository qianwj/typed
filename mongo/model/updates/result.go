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

package updates

import (
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateResult[I any] struct {
	MatchedCount  int64 // The number of documents matched by the filter.
	ModifiedCount int64 // The number of documents modified by the operation.
	UpsertedCount int64 // The number of documents upserted by the operation.
	UpsertedID    I     // The _id field of the upserted document, or nil if no upsert was done.
}

func FromUpdateResult[I any](r *mongo.UpdateResult) *UpdateResult[I] {
	if r == nil {
		return nil
	}
	return &UpdateResult[I]{
		MatchedCount:  r.MatchedCount,
		ModifiedCount: r.ModifiedCount,
		UpsertedCount: r.UpsertedCount,
		UpsertedID:    r.UpsertedID.(I),
	}
}

// BulkWriteResult is the result type returned by a BulkWrite operation.
type BulkWriteResult[I model.ID] struct {
	// The number of documents inserted.
	InsertedCount int64

	// The number of documents matched by filters in update and replace operations.
	MatchedCount int64

	// The number of documents modified by update and replace operations.
	ModifiedCount int64

	// The number of documents deleted.
	DeletedCount int64

	// The number of documents upserted by update and replace operations.
	UpsertedCount int64

	// A map of operation index to the _id of each upserted document.
	UpsertedIDs map[int64]I
}

func FromBulkWriteResult[I model.ID](r *mongo.BulkWriteResult) *BulkWriteResult[I] {
	if r == nil {
		return nil
	}
	upsertIds := make(map[int64]I)
	for k, v := range r.UpsertedIDs {
		upsertIds[k] = v.(I)
	}
	return &BulkWriteResult[I]{
		MatchedCount:  r.MatchedCount,
		ModifiedCount: r.ModifiedCount,
		UpsertedCount: r.UpsertedCount,
		UpsertedIDs:   upsertIds,
	}
}
