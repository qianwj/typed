package options

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CountOptions struct {
	internal *options.CountOptions
}

// Count creates a new CountOptions instance.
func Count() *CountOptions {
	return &CountOptions{}
}

// SetCollation sets the value for the Collation field.
func (co *CountOptions) SetCollation(c *options.Collation) *CountOptions {
	co.internal.Collation = c
	return co
}

// SetHint sets the value for the Hint field.
func (co *CountOptions) SetHint(h interface{}) *CountOptions {
	co.internal.Hint = h
	return co
}

// SetLimit sets the value for the Limit field.
func (co *CountOptions) SetLimit(i int64) *CountOptions {
	co.internal.Limit = &i
	return co
}

// SetMaxTime sets the value for the MaxTime field.
func (co *CountOptions) SetMaxTime(d time.Duration) *CountOptions {
	co.internal.MaxTime = &d
	return co
}

// SetSkip sets the value for the Skip field.
func (co *CountOptions) SetSkip(i int64) *CountOptions {
	co.internal.Skip = &i
	return co
}

// MergeCountOptions combines the given CountOptions instances into a single CountOptions in a last-one-wins fashion.
func MergeCountOptions(opts ...*CountOptions) *options.CountOptions {
	countOpts := options.Count()
	for _, opt := range opts {
		co := opt.internal
		if co == nil {
			continue
		}
		if co.Collation != nil {
			countOpts.Collation = co.Collation
		}
		if co.Hint != nil {
			countOpts.Hint = co.Hint
		}
		if co.Limit != nil {
			countOpts.Limit = co.Limit
		}
		if co.MaxTime != nil {
			countOpts.MaxTime = co.MaxTime
		}
		if co.Skip != nil {
			countOpts.Skip = co.Skip
		}
	}

	return countOpts
}
