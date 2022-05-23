package options

import "go.mongodb.org/mongo-driver/mongo/options"

type DeleteOptions struct {
	internal *options.DeleteOptions
}

// Delete creates a new DeleteOptions instance.
func Delete() *DeleteOptions {
	return &DeleteOptions{
		internal: options.Delete(),
	}
}

// SetCollation sets the value for the Collation field.
func (do *DeleteOptions) SetCollation(c *options.Collation) *DeleteOptions {
	do.internal.Collation = c
	return do
}

// SetHint sets the value for the Hint field.
func (do *DeleteOptions) SetHint(hint interface{}) *DeleteOptions {
	do.internal.Hint = hint
	return do
}

// SetLet sets the value for the Let field.
func (do *DeleteOptions) SetLet(let interface{}) *DeleteOptions {
	do.internal.Let = let
	return do
}

// MergeDeleteOptions combines the given DeleteOptions instances into a single DeleteOptions in a last-one-wins fashion.
func MergeDeleteOptions(opts ...*DeleteOptions) *options.DeleteOptions {
	dOpts := options.Delete()
	for _, opt := range opts {
		do := opt.internal
		if do == nil {
			continue
		}
		if do.Collation != nil {
			dOpts.Collation = do.Collation
		}
		if do.Hint != nil {
			dOpts.Hint = do.Hint
		}
		if do.Let != nil {
			dOpts.Let = do.Let
		}
	}

	return dOpts
}
