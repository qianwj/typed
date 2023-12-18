package requests

import (
	indicescreate "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type CreateIndexRequestBuilder struct {
	internal *indicescreate.Request
}

func CreateIndex() *CreateIndexRequestBuilder {
	return &CreateIndexRequestBuilder{internal: indicescreate.NewRequest()}
}

// Aliases Aliases for the index.
func (c *CreateIndexRequestBuilder) Aliases(aliases map[string]types.Alias) *CreateIndexRequestBuilder {
	c.internal.Aliases = aliases
	return c
}

// Mappings Mapping for fields in the index. If specified, this mapping can include:
// - Field names
// - Field data types
// - Mapping parameters
func (c *CreateIndexRequestBuilder) Mappings(mappings *types.TypeMapping) *CreateIndexRequestBuilder {
	c.internal.Mappings = mappings
	return c
}

// Settings Configuration options for the index.
func (c *CreateIndexRequestBuilder) Settings(settings *types.IndexSettings) *CreateIndexRequestBuilder {
	c.internal.Settings = settings
	return c
}

func (c *CreateIndexRequestBuilder) Build() *indicescreate.Request {
	return c.internal
}
