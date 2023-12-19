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

// IndicesBoost Boosts the _score of documents from specified indices.
func (s *SearchRequestBuilder) IndicesBoost(indicesBoost []map[string]types.Float64) *SearchRequestBuilder {
	s.internal.IndicesBoost = indicesBoost
	return s
}

// Knn Defines the approximate kNN search to run.
func (s *SearchRequestBuilder) Knn(knn []types.KnnQuery) *SearchRequestBuilder {
	s.internal.Knn = knn
	return s
}

// MinScore Minimum `_score` for matching documents.
// Documents with a lower `_score` are not included in the search results.
func (s *SearchRequestBuilder) MinScore(minScore *types.Float64) *SearchRequestBuilder {
	s.internal.MinScore = minScore
	return s
}

// Pit Limits the search to a point in time (PIT).
// If you provide a PIT, you cannot specify an `<index>` in the request path.
func (s *SearchRequestBuilder) Pit(pit *types.PointInTimeReference) *SearchRequestBuilder {
	s.internal.Pit = pit
	return s
}

// PostFilter Use the `post_filter` parameter to filter search results.
// The search hits are filtered after the aggregations are calculated.
// A post filter has no impact on the aggregation results.
func (s *SearchRequestBuilder) PostFilter(postFilter *types.Query) *SearchRequestBuilder {
	s.internal.PostFilter = postFilter
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

// Rank Defines the Reciprocal Rank Fusion (RRF) to use.
func (s *SearchRequestBuilder) Rank(rank *types.RankContainer) *SearchRequestBuilder {
	s.internal.Rank = rank
	return s
}

// Rescore Can be used to improve precision by reordering just the top (for example 100
// - 500) documents returned by the `query` and `post_filter` phases.
func (s *SearchRequestBuilder) Rescore(rescore []types.Rescore) *SearchRequestBuilder {
	s.internal.Rescore = rescore
	return s
}

// RuntimeMappings Defines one or more runtime fields in the search request.
// These fields take precedence over mapped fields with the same name.
func (s *SearchRequestBuilder) RuntimeMappings(mappings types.RuntimeFields) *SearchRequestBuilder {
	s.internal.RuntimeMappings = mappings
	return s
}

// ScriptFields Retrieve a script evaluation (based on different fields) for each hit.
func (s *SearchRequestBuilder) ScriptFields(fields map[string]types.ScriptField) *SearchRequestBuilder {
	s.internal.ScriptFields = fields
	return s
}

// SearchAfter Used to retrieve the next page of hits using a set of sort values from the
// previous page.
func (s *SearchRequestBuilder) SearchAfter(searchAfter []types.FieldValue) *SearchRequestBuilder {
	s.internal.SearchAfter = searchAfter
	return s
}

// SeqNoPrimaryTerm If `true`, returns sequence number and primary term of the last modification
// of each hit.
func (s *SearchRequestBuilder) SeqNoPrimaryTerm(seqNoPrimaryTerm bool) *SearchRequestBuilder {
	s.internal.SeqNoPrimaryTerm = &seqNoPrimaryTerm
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

// Slice Can be used to split a scrolled search into multiple slices that can be
// consumed independently.
func (s *SearchRequestBuilder) Slice(slice *types.SlicedScroll) *SearchRequestBuilder {
	s.internal.Slice = slice
	return s
}

// Sort A comma-separated list of <field>:<direction> pairs.
func (s *SearchRequestBuilder) Sort(sort []types.SortCombinations) *SearchRequestBuilder {
	s.internal.Sort = sort
	return s
}

// Source Indicates which source fields are returned for matching documents.
// These fields are returned in the hits._source property of the search
// response.
func (s *SearchRequestBuilder) Source(source types.SourceConfig) *SearchRequestBuilder {
	s.internal.Source_ = source
	return s
}

// Stats groups to associate with the search.
// Each group maintains a statistics aggregation for its associated searches.
// You can retrieve these stats using the indices stats API.
func (s *SearchRequestBuilder) Stats(stats []string) *SearchRequestBuilder {
	s.internal.Stats = stats
	return s
}

// StoredFields List of stored fields to return as part of a hit.
// If no fields are specified, no stored fields are included in the response.
// If this field is specified, the `_source` parameter defaults to `false`.
// You can pass `_source: true` to return both source fields and stored fields
// in the search response.
func (s *SearchRequestBuilder) StoredFields(fields []string) *SearchRequestBuilder {
	s.internal.StoredFields = fields
	return s
}

// Suggest Defines a suggester that provides similar looking terms based on a provided
// text.
func (s *SearchRequestBuilder) Suggest(suggest *types.Suggester) *SearchRequestBuilder {
	s.internal.Suggest = suggest
	return s
}

// TerminateAfter Maximum number of documents to collect for each shard.
// If a query reaches this limit, Elasticsearch terminates the query early.
// Elasticsearch collects documents before sorting.
// Use with caution.
// Elasticsearch applies this parameter to each shard handling the request.
// When possible, let Elasticsearch perform early termination automatically.
// Avoid specifying this parameter for requests that target data streams with
// backing indices across multiple data tiers.
// If set to `0` (default), the query does not terminate early.
func (s *SearchRequestBuilder) TerminateAfter(terminateAfter int64) *SearchRequestBuilder {
	s.internal.TerminateAfter = &terminateAfter
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

// TrackScores If true, calculate and return document scores, even if the scores are not
// used for sorting.
func (s *SearchRequestBuilder) TrackScores(scores bool) *SearchRequestBuilder {
	s.internal.TrackScores = &scores
	return s
}

// TrackTotalHits Number of hits matching the query to count accurately.
// If `true`, the exact number of hits is returned at the cost of some
// performance.
// If `false`, the  response does not include the total number of hits matching
// the query.
func (s *SearchRequestBuilder) TrackTotalHits(trackTotalHits types.TrackHits) *SearchRequestBuilder {
	s.internal.TrackTotalHits = trackTotalHits
	return s
}

// Version If true, returns document version as part of a hit.
func (s *SearchRequestBuilder) Version(version bool) *SearchRequestBuilder {
	s.internal.Version = &version
	return s
}

func (s *SearchRequestBuilder) Build() *search.Request {
	return s.internal
}
