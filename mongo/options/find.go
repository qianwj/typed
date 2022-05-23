package options

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	desc SortOrder = -1
	asc  SortOrder = 1
)

type (
	SortOrder   int
	SortOptions bson.D
	SortOption  bson.E
	Projection  bson.M
)

func (opts SortOptions) Append(next SortOptions) SortOptions {
	return append(opts, next...)
}

func Ascending(field string) SortOptions {
	return SortOptions{bson.E{Key: field, Value: asc}}
}

func Descending(field string) SortOptions {
	return SortOptions{bson.E{Key: field, Value: desc}}
}

func Need(field string) Projection {
	return Projection{
		field: 1,
	}
}

func Ignore(field string) Projection {
	return Projection{
		field: -1,
	}
}

func (p Projection) And(projection Projection) Projection {
	for k, v := range projection {
		p[k] = v
	}
	return p
}

type FindOptions struct {
	internal *options.FindOptions
}

func Find() *FindOptions {
	return &FindOptions{
		internal: options.Find(),
	}
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (f *FindOptions) SetAllowDiskUse(b bool) *FindOptions {
	f.internal.AllowDiskUse = &b
	return f
}

// SetAllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindOptions) SetAllowPartialResults(b bool) *FindOptions {
	f.internal.AllowPartialResults = &b
	return f
}

// SetBatchSize sets the value for the BatchSize field.
func (f *FindOptions) SetBatchSize(i int32) *FindOptions {
	f.internal.BatchSize = &i
	return f
}

// SetCollation sets the value for the Collation field.
func (f *FindOptions) SetCollation(collation *options.Collation) *FindOptions {
	f.internal.Collation = collation
	return f
}

// SetComment sets the value for the Comment field.
func (f *FindOptions) SetComment(comment string) *FindOptions {
	f.internal.Comment = &comment
	return f
}

// SetCursorType sets the value for the CursorType field.
func (f *FindOptions) SetCursorType(ct options.CursorType) *FindOptions {
	f.internal.CursorType = &ct
	return f
}

// SetHint sets the value for the Hint field.
func (f *FindOptions) SetHint(hint interface{}) *FindOptions {
	f.internal.Hint = hint
	return f
}

// SetLet sets the value for the Let field.
func (f *FindOptions) SetLet(let interface{}) *FindOptions {
	f.internal.Let = let
	return f
}

// SetLimit sets the value for the Limit field.
func (f *FindOptions) SetLimit(i int64) *FindOptions {
	f.internal.Limit = &i
	return f
}

// SetMax sets the value for the Max field.
func (f *FindOptions) SetMax(max interface{}) *FindOptions {
	f.internal.Max = max
	return f
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (f *FindOptions) SetMaxAwaitTime(d time.Duration) *FindOptions {
	f.internal.MaxAwaitTime = &d
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOptions) SetMaxTime(d time.Duration) *FindOptions {
	f.internal.MaxTime = &d
	return f
}

// SetMin sets the value for the Min field.
func (f *FindOptions) SetMin(min interface{}) *FindOptions {
	f.internal.Min = min
	return f
}

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
func (f *FindOptions) SetNoCursorTimeout(b bool) *FindOptions {
	f.internal.NoCursorTimeout = &b
	return f
}

// SetOplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is set.
func (f *FindOptions) SetOplogReplay(b bool) *FindOptions {
	f.internal.OplogReplay = &b
	return f
}

// SetProjection sets the value for the Projection field.
func (f *FindOptions) SetProjection(projection Projection) *FindOptions {
	f.internal.Projection = projection
	return f
}

// SetReturnKey sets the value for the ReturnKey field.
func (f *FindOptions) SetReturnKey(b bool) *FindOptions {
	f.internal.ReturnKey = &b
	return f
}

// SetShowRecordID sets the value for the ShowRecordID field.
func (f *FindOptions) SetShowRecordID(b bool) *FindOptions {
	f.internal.ShowRecordID = &b
	return f
}

// SetSkip sets the value for the Skip field.
func (f *FindOptions) SetSkip(i int64) *FindOptions {
	f.internal.Skip = &i
	return f
}

// SetSnapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindOptions) SetSnapshot(b bool) *FindOptions {
	f.internal.Snapshot = &b
	return f
}

// SetSort sets the value for the Sort field.
func (f *FindOptions) SetSort(sort SortOptions) *FindOptions {
	f.internal.Sort = sort
	return f
}

// MergeFindOptions combines the given FindOptions instances into a single FindOptions in a last-one-wins fashion.
func MergeFindOptions(opts ...*FindOptions) *options.FindOptions {
	fo := options.Find()
	for _, o := range opts {
		opt := o.internal
		if opt == nil {
			continue
		}
		if opt.AllowDiskUse != nil {
			fo.AllowDiskUse = opt.AllowDiskUse
		}
		if opt.AllowPartialResults != nil {
			fo.AllowPartialResults = opt.AllowPartialResults
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.CursorType != nil {
			fo.CursorType = opt.CursorType
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
		if opt.Limit != nil {
			fo.Limit = opt.Limit
		}
		if opt.Max != nil {
			fo.Max = opt.Max
		}
		if opt.MaxAwaitTime != nil {
			fo.MaxAwaitTime = opt.MaxAwaitTime
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Min != nil {
			fo.Min = opt.Min
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.OplogReplay != nil {
			fo.OplogReplay = opt.OplogReplay
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnKey != nil {
			fo.ReturnKey = opt.ReturnKey
		}
		if opt.ShowRecordID != nil {
			fo.ShowRecordID = opt.ShowRecordID
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Snapshot != nil {
			fo.Snapshot = opt.Snapshot
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}
	return fo
}

type FindOneAndUpdateOptions struct {
	internal *options.FindOneAndUpdateOptions
}

// FindOneAndUpdate creates a new FindOneAndUpdateOptions instance.
func FindOneAndUpdate() *FindOneAndUpdateOptions {
	return &FindOneAndUpdateOptions{
		internal: options.FindOneAndUpdate(),
	}
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndUpdateOptions) SetBypassDocumentValidation(b bool) *FindOneAndUpdateOptions {
	f.internal.BypassDocumentValidation = &b
	return f
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (f *FindOneAndUpdateOptions) SetArrayFilters(filters options.ArrayFilters) *FindOneAndUpdateOptions {
	f.internal.ArrayFilters = &filters
	return f
}

// SetCollation sets the value for the Collation field.
func (f *FindOneAndUpdateOptions) SetCollation(collation *options.Collation) *FindOneAndUpdateOptions {
	f.internal.Collation = collation
	return f
}

// SetMaxTime sets the value for the MaxTime field.
func (f *FindOneAndUpdateOptions) SetMaxTime(d time.Duration) *FindOneAndUpdateOptions {
	f.internal.MaxTime = &d
	return f
}

// SetProjection sets the value for the Projection field.
func (f *FindOneAndUpdateOptions) SetProjection(projection Projection) *FindOneAndUpdateOptions {
	f.internal.Projection = projection
	return f
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (f *FindOneAndUpdateOptions) SetReturnDocument(rd options.ReturnDocument) *FindOneAndUpdateOptions {
	f.internal.ReturnDocument = &rd
	return f
}

// SetSort sets the value for the Sort field.
func (f *FindOneAndUpdateOptions) SetSort(sort SortOptions) *FindOneAndUpdateOptions {
	f.internal.Sort = sort
	return f
}

// SetUpsert sets the value for the Upsert field.
func (f *FindOneAndUpdateOptions) SetUpsert(b bool) *FindOneAndUpdateOptions {
	f.internal.Upsert = &b
	return f
}

// SetHint sets the value for the Hint field.
func (f *FindOneAndUpdateOptions) SetHint(hint interface{}) *FindOneAndUpdateOptions {
	f.internal.Hint = hint
	return f
}

// SetLet sets the value for the Let field.
func (f *FindOneAndUpdateOptions) SetLet(let interface{}) *FindOneAndUpdateOptions {
	f.internal.Let = let
	return f
}

// MergeFindOneAndUpdateOptions combines the given FindOneAndUpdateOptions instances into a single
// FindOneAndUpdateOptions in a last-one-wins fashion.
func MergeFindOneAndUpdateOptions(opts ...*FindOneAndUpdateOptions) *options.FindOneAndUpdateOptions {
	fo := options.FindOneAndUpdate()
	for _, o := range opts {
		opt := o.internal
		if opt == nil {
			continue
		}
		if opt.ArrayFilters != nil {
			fo.ArrayFilters = opt.ArrayFilters
		}
		if opt.BypassDocumentValidation != nil {
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnDocument != nil {
			fo.ReturnDocument = opt.ReturnDocument
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Upsert != nil {
			fo.Upsert = opt.Upsert
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
	}

	return fo
}
