package requests

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type SearchRequestBuilder struct {
	internal *search.Request
}

func Search() *SearchRequestBuilder {
	return &SearchRequestBuilder{internal: search.NewRequest()}
}

// Aggregations Defines the aggregations that are run as part of the search request.
func (s *SearchRequestBuilder) Aggregations(agg map[string]types.Aggregations) *SearchRequestBuilder {
	s.internal.Aggregations = agg
	return s
}

// Collapse Collapses search results the values of the specified field.
func (s *SearchRequestBuilder) Collapse(collapse *types.FieldCollapse) *SearchRequestBuilder {
	s.internal.Collapse = collapse
	return s
}

// DocValueFields Array of wildcard (`*`) patterns.
// The request returns doc values for field names matching these patterns in the
// `hits.fields` property of the response.
func (s *SearchRequestBuilder) DocValueFields(docValueFields []types.FieldAndFormat) *SearchRequestBuilder {
	s.internal.DocvalueFields = docValueFields
	return s
}

// Explain If true, returns detailed information about score computation as part of a
// hit.
func (s *SearchRequestBuilder) Explain(explain bool) *SearchRequestBuilder {
	s.internal.Explain = &explain
	return s
}

// Ext Configuration of search extensions defined by Elasticsearch plugins.
func (s *SearchRequestBuilder) Ext(ext map[string]json.RawMessage) *SearchRequestBuilder {
	s.internal.Ext = ext
	return s
}

// Fields Array of wildcard (`*`) patterns.
// The request returns values for field names matching these patterns in the
// `hits.fields` property of the response.
func (s *SearchRequestBuilder) Fields(fields []types.FieldAndFormat) *SearchRequestBuilder {
	s.internal.Fields = fields
	return s
}

// From Starting document offset.
// Needs to be non-negative.
// By default, you cannot page through more than 10,000 hits using the `from`
// and `size` parameters.
// To page through more hits, use the `search_after` parameter.
func (s *SearchRequestBuilder) From(from int) *SearchRequestBuilder {
	s.internal.From = &from
	return s
}

// Highlight Specifies the highlighter to use for retrieving highlighted snippets from one
// or more fields in your search results.
func (s *SearchRequestBuilder) Highlight(highlight *types.Highlight) *SearchRequestBuilder {
	s.internal.Highlight = highlight
	return s
}

// Profile Set to `true` to return detailed timing information about the execution of
// individual components in a search request.
// NOTE: This is a debugging tool and adds significant overhead to search
// execution.
func (s *SearchRequestBuilder) Profile(profile bool) *SearchRequestBuilder {
	s.internal.Profile = &profile
	return s
}

// Query Defines the search definition using the Query DSL.
func (s *SearchRequestBuilder) Query(query *types.Query) *SearchRequestBuilder {
	s.internal.Query = query
	return s
}

// Size The number of hits to return.
// By default, you cannot page through more than 10,000 hits using the `from`
// and `size` parameters.
// To page through more hits, use the `search_after` parameter.
func (s *SearchRequestBuilder) Size(size int) *SearchRequestBuilder {
	s.internal.Size = &size
	return s
}

// Sort A comma-separated list of <field>:<direction> pairs.
func (s *SearchRequestBuilder) Sort(sort []types.SortCombinations) *SearchRequestBuilder {
	s.internal.Sort = sort
	return s
}

// Timeout Specifies the period of time to wait for a response from each shard.
// If no response is received before the timeout expires, the request fails and
// returns an error.
// Defaults to no timeout.
func (s *SearchRequestBuilder) Timeout(timeout string) *SearchRequestBuilder {
	s.internal.Timeout = &timeout
	return s
}

// Version If true, returns document version as part of a hit.
func (s *SearchRequestBuilder) Version(version bool) *SearchRequestBuilder {
	s.internal.Version = &version
	return s
}

//
//// EventCategoryField Field containing the event classification, such as process, file, or network.
//func (s *SearchRequestBuilder) EventCategoryField(eventCategoryField string) *SearchRequestBuilder {
//	s.internal.EventCategoryField = &eventCategoryField
//	return s
//}
//

//
//// Fields Array of wildcard (*) patterns. The response returns values for field names
//// matching these patterns in the fields property of each hit.
//func (s *SearchRequestBuilder) Fields(fields []types.FieldAndFormat) *SearchRequestBuilder {
//	s.internal.Fields = fields
//	return s
//}
//
//// Filter Query, written in Query DSL, used to filter the events on which the EQL query
//// runs.
//func (s *SearchRequestBuilder) Filter(filters []types.Query) *SearchRequestBuilder {
//	s.internal.Filter = filters
//	return s
//}
//
//func (s *SearchRequestBuilder) KeepAlive(keepAlive types.Duration) *SearchRequestBuilder {
//	s.internal.KeepAlive = keepAlive
//	return s
//}
//
//func (s *SearchRequestBuilder) KeepOnCompletion(keepOnCompletion *bool) *SearchRequestBuilder {
//	s.internal.KeepOnCompletion = keepOnCompletion
//	return s
//}
//
//func (s *SearchRequestBuilder) ResultPosition(pos *resultposition.ResultPosition) *SearchRequestBuilder {
//	s.internal.ResultPosition = pos
//	return s
//}
//
//func (s *SearchRequestBuilder) RuntimeMappings(fields types.RuntimeFields) *SearchRequestBuilder {
//	s.internal.RuntimeMappings = fields
//	return s
//}
//
//// Size For basic queries, the maximum number of matching events to return. Defaults
//// to 10
//func (s *SearchRequestBuilder) Size(size uint) *SearchRequestBuilder {
//	s.internal.Size = &size
//	return s
//}
//
//// TiebreakerField Field used to sort hits with the same timestamp in ascending order
//func (s *SearchRequestBuilder) TiebreakerField(field string) *SearchRequestBuilder {
//	s.internal.TiebreakerField = &field
//	return s
//}
//
//// TimestampField Field containing event timestamp. Default "@timestamp"
//func (s *SearchRequestBuilder) TimestampField(field string) *SearchRequestBuilder {
//	s.internal.TimestampField = &field
//	return s
//}
//
//func (s *SearchRequestBuilder) WaitForCompletionTimeout(duration types.Duration) *SearchRequestBuilder {
//	s.internal.WaitForCompletionTimeout = duration
//	return s
//}

func (s *SearchRequestBuilder) Build() *search.Request {
	return s.internal
}
