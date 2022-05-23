package options

import "go.mongodb.org/mongo-driver/mongo/options"

type UpdateOptions struct {
	internal *options.UpdateOptions
}

// Update creates a new UpdateOptions instance.
func Update() *UpdateOptions {
	return &UpdateOptions{
		internal: options.Update(),
	}
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (uo *UpdateOptions) SetArrayFilters(af options.ArrayFilters) *UpdateOptions {
	uo.internal.ArrayFilters = &af
	return uo
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions {
	uo.internal.BypassDocumentValidation = &b
	return uo
}

// SetCollation sets the value for the Collation field.
func (uo *UpdateOptions) SetCollation(c *options.Collation) *UpdateOptions {
	uo.internal.Collation = c
	return uo
}

// SetHint sets the value for the Hint field.
func (uo *UpdateOptions) SetHint(h interface{}) *UpdateOptions {
	uo.internal.Hint = h
	return uo
}

// SetUpsert sets the value for the Upsert field.
func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions {
	uo.internal.Upsert = &b
	return uo
}

// SetLet sets the value for the Let field.
func (uo *UpdateOptions) SetLet(l interface{}) *UpdateOptions {
	uo.internal.Let = l
	return uo
}

// MergeUpdateOptions combines the given UpdateOptions instances into a single UpdateOptions in a last-one-wins fashion.
func MergeUpdateOptions(opts ...*UpdateOptions) *options.UpdateOptions {
	uOpts := options.Update()
	for _, opt := range opts {
		uo := opt.internal
		if uo == nil {
			continue
		}
		if uo.ArrayFilters != nil {
			uOpts.ArrayFilters = uo.ArrayFilters
		}
		if uo.BypassDocumentValidation != nil {
			uOpts.BypassDocumentValidation = uo.BypassDocumentValidation
		}
		if uo.Collation != nil {
			uOpts.Collation = uo.Collation
		}
		if uo.Hint != nil {
			uOpts.Hint = uo.Hint
		}
		if uo.Upsert != nil {
			uOpts.Upsert = uo.Upsert
		}
		if uo.Let != nil {
			uOpts.Let = uo.Let
		}
	}

	return uOpts
}
