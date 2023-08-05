package model

import "go.mongodb.org/mongo-driver/mongo"

type UpdateResult[I DocumentId] struct {
	MatchedCount  int64 // The number of documents matched by the filter.
	ModifiedCount int64 // The number of documents modified by the operation.
	UpsertedCount int64 // The number of documents upserted by the operation.
	UpsertedID    I     // The _id field of the upserted document, or nil if no upsert was done.
}

func FromUpdateResult[I DocumentId](r *mongo.UpdateResult) *UpdateResult[I] {
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
type BulkWriteResult[I DocumentId] struct {
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

func FromBulkWriteResult[I DocumentId](r *mongo.BulkWriteResult) *BulkWriteResult[I] {
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
