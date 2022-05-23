package options

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneOptions struct {
	internal *options.FindOneOptions
}

// FindOne creates a new FindOneOptions instance.
func FindOne() *FindOneOptions {
	return &FindOneOptions{}
}

// SetAllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindOneOptions) SetAllowPartialResults(b bool) *FindOneOptions {
	f.internal.AllowPartialResults = &b
	return f
}

// SetBatchSize sets the value for the BatchSize field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetBatchSize(i int32) *FindOneOptions {
	f.internal.BatchSize = &i
	return f
}

// SetCollation sets the value for the Collation field.
func (f *FindOneOptions) SetCollation(collation *options.Collation) *FindOneOptions {
	f.internal.Collation = collation
	return f
}

// SetComment sets the value for the Comment field.
func (f *FindOneOptions) SetComment(comment string) *FindOneOptions {
	f.internal.Comment = &comment
	return f
}

// SetCursorType sets the value for the CursorType field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetCursorType(ct options.CursorType) *FindOneOptions {
	f.internal.CursorType = &ct
	return f
}

// SetHint sets the value for the Hint field.
func (f *FindOneOptions) SetHint(hint interface{}) *FindOneOptions {
	f.internal.Hint = hint
	return f
}

// SetMax sets the value for the Max field.
func (f *FindOneOptions) SetMax(max interface{}) *FindOneOptions {
	f.internal.Max = max
	return f
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetMaxAwaitTime(d time.Duration) *FindOneOptions {
	f.internal.MaxAwaitTime = &d
	return f
}

// SetMaxTime sets the value for the MaxTime field.
func (f *FindOneOptions) SetMaxTime(d time.Duration) *FindOneOptions {
	f.internal.MaxTime = &d
	return f
}

// SetMin sets the value for the Min field.
func (f *FindOneOptions) SetMin(min interface{}) *FindOneOptions {
	f.internal.Min = min
	return f
}

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetNoCursorTimeout(b bool) *FindOneOptions {
	f.internal.NoCursorTimeout = &b
	return f
}

// SetOplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is
// set.
func (f *FindOneOptions) SetOplogReplay(b bool) *FindOneOptions {
	f.internal.OplogReplay = &b
	return f
}

// SetProjection sets the value for the Projection field.
func (f *FindOneOptions) SetProjection(projection Projection) *FindOneOptions {
	f.internal.Projection = projection
	return f
}

// SetReturnKey sets the value for the ReturnKey field.
func (f *FindOneOptions) SetReturnKey(b bool) *FindOneOptions {
	f.internal.ReturnKey = &b
	return f
}

// SetShowRecordID sets the value for the ShowRecordID field.
func (f *FindOneOptions) SetShowRecordID(b bool) *FindOneOptions {
	f.internal.ShowRecordID = &b
	return f
}

// SetSkip sets the value for the Skip field.
func (f *FindOneOptions) SetSkip(i int64) *FindOneOptions {
	f.internal.Skip = &i
	return f
}

// SetSnapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindOneOptions) SetSnapshot(b bool) *FindOneOptions {
	f.internal.Snapshot = &b
	return f
}

// SetSort sets the value for the Sort field.
func (f *FindOneOptions) SetSort(sort SortOptions) *FindOneOptions {
	f.internal.Sort = sort
	return f
}

// MergeFindOneOptions combines the given FindOneOptions instances into a single FindOneOptions in a last-one-wins
// fashion.
func MergeFindOneOptions(opts ...*FindOneOptions) *options.FindOneOptions {
	fo := options.FindOne()
	for _, o := range opts {
		opt := o.internal
		if opt == nil {
			continue
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
