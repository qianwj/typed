package options

import (
	"github.com/qianwj/typed/mongo/model/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func (uo *UpdateOptions) SetArrayFilters(af ArrayFilters) *UpdateOptions {
	uo.internal.SetArrayFilters(af.Raw())
	return uo
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions {
	uo.internal.BypassDocumentValidation = &b
	return uo
}

// SetCollation sets the value for the Collation field.
func (uo *UpdateOptions) SetCollation(c *Collation) *UpdateOptions {
	uo.internal.Collation = (*options.Collation)(c)
	return uo
}

// SetHint sets the value for the Hint field.
func (uo *UpdateOptions) SetHint(index string) *UpdateOptions {
	uo.internal.Hint = index
	return uo
}

// SetUpsert sets the value for the Upsert field.
func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions {
	uo.internal.Upsert = &b
	return uo
}

// SetLet sets the value for the Let field.
func (uo *UpdateOptions) SetLet(l bson.M) *UpdateOptions {
	uo.internal.Let = l
	return uo
}

type ArrayFilters struct {
	Items    []*filter.Filter
	Registry *bsoncodec.Registry
}

func (af *ArrayFilters) Raw() options.ArrayFilters {
	// complete this function
	return options.ArrayFilters{}
}
