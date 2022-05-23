package options

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AggregateOptions struct {
	internal *options.AggregateOptions
}

// Aggregate creates a new AggregateOptions instance.
func Aggregate() *AggregateOptions {
	return &AggregateOptions{}
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (ao *AggregateOptions) SetAllowDiskUse(b bool) *AggregateOptions {
	ao.internal.AllowDiskUse = &b
	return ao
}

// SetBatchSize sets the value for the BatchSize field.
func (ao *AggregateOptions) SetBatchSize(i int32) *AggregateOptions {
	ao.internal.BatchSize = &i
	return ao
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ao *AggregateOptions) SetBypassDocumentValidation(b bool) *AggregateOptions {
	ao.internal.BypassDocumentValidation = &b
	return ao
}

// SetCollation sets the value for the Collation field.
func (ao *AggregateOptions) SetCollation(c *options.Collation) *AggregateOptions {
	ao.internal.Collation = c
	return ao
}

// SetMaxTime sets the value for the MaxTime field.
func (ao *AggregateOptions) SetMaxTime(d time.Duration) *AggregateOptions {
	ao.internal.MaxTime = &d
	return ao
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (ao *AggregateOptions) SetMaxAwaitTime(d time.Duration) *AggregateOptions {
	ao.internal.MaxAwaitTime = &d
	return ao
}

// SetComment sets the value for the Comment field.
func (ao *AggregateOptions) SetComment(s string) *AggregateOptions {
	ao.internal.Comment = &s
	return ao
}

// SetHint sets the value for the Hint field.
func (ao *AggregateOptions) SetHint(h interface{}) *AggregateOptions {
	ao.internal.Hint = h
	return ao
}

// SetLet sets the value for the Let field.
func (ao *AggregateOptions) SetLet(let interface{}) *AggregateOptions {
	ao.internal.Let = let
	return ao
}

// SetCustom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (ao *AggregateOptions) SetCustom(c bson.M) *AggregateOptions {
	ao.internal.Custom = c
	return ao
}

// MergeAggregateOptions combines the given AggregateOptions instances into a single AggregateOptions in a last-one-wins
// fashion.
func MergeAggregateOptions(opts ...*AggregateOptions) *options.AggregateOptions {
	aggOpts := options.Aggregate()
	for _, opt := range opts {
		ao := opt.internal
		if ao == nil {
			continue
		}
		if ao.AllowDiskUse != nil {
			aggOpts.AllowDiskUse = ao.AllowDiskUse
		}
		if ao.BatchSize != nil {
			aggOpts.BatchSize = ao.BatchSize
		}
		if ao.BypassDocumentValidation != nil {
			aggOpts.BypassDocumentValidation = ao.BypassDocumentValidation
		}
		if ao.Collation != nil {
			aggOpts.Collation = ao.Collation
		}
		if ao.MaxTime != nil {
			aggOpts.MaxTime = ao.MaxTime
		}
		if ao.MaxAwaitTime != nil {
			aggOpts.MaxAwaitTime = ao.MaxAwaitTime
		}
		if ao.Comment != nil {
			aggOpts.Comment = ao.Comment
		}
		if ao.Hint != nil {
			aggOpts.Hint = ao.Hint
		}
		if ao.Let != nil {
			aggOpts.Let = ao.Let
		}
		if ao.Custom != nil {
			aggOpts.Custom = ao.Custom
		}
	}

	return aggOpts
}
