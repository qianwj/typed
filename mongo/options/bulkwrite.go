package options

import "go.mongodb.org/mongo-driver/mongo/options"

type BulkWriteOptions struct {
	internal *options.BulkWriteOptions
}

// BulkWrite creates a new *BulkWriteOptions instance.
func BulkWrite() *BulkWriteOptions {
	return &BulkWriteOptions{
		internal: options.BulkWrite(),
	}
}

// SetOrdered sets the value for the Ordered field.
func (b *BulkWriteOptions) SetOrdered(ordered bool) *BulkWriteOptions {
	b.internal.Ordered = &ordered
	return b
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (b *BulkWriteOptions) SetBypassDocumentValidation(bypass bool) *BulkWriteOptions {
	b.internal.BypassDocumentValidation = &bypass
	return b
}

// SetLet sets the value for the Let field. Let specifies parameters for all update and delete commands in the BulkWrite.
// This option is only valid for MongoDB versions >= 5.0. Older servers will report an error for using this option.
// This must be a document mapping parameter names to values. Values must be constant or closed expressions that do not
// reference document fields. Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
func (b *BulkWriteOptions) SetLet(let interface{}) *BulkWriteOptions {
	b.internal.Let = &let
	return b
}

// MergeBulkWriteOptions combines the given BulkWriteOptions instances into a single BulkWriteOptions in a last-one-wins
// fashion.
func MergeBulkWriteOptions(opts ...*BulkWriteOptions) *options.BulkWriteOptions {
	b := options.BulkWrite()
	for _, o := range opts {
		opt := o.internal
		if opt == nil {
			continue
		}
		if opt.Ordered != nil {
			b.Ordered = opt.Ordered
		}
		if opt.BypassDocumentValidation != nil {
			b.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Let != nil {
			b.Let = opt.Let
		}
	}

	return b
}
