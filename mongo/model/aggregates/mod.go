/** MIT License
 *
 * Copyright (c) 2023 qianwj
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates/lookup"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
)

func Count(field string) *Pipeline {
	return New().Count(field)
}

func GraphLookup(cond *lookup.GraphJoinCondition) *Pipeline {
	return New().GraphLookup(cond)
}

func Group(id any, fields ...bson.Entry) *Pipeline {
	return New().Group(id, fields...)
}

func Limit(limit int64) *Pipeline {
	return New().Limit(limit)
}

// Lookup Performs a left outer join to a collection in the same database to filter in documents from the "joined"
// collection for processing. The `$lookup` stage adds a new array field to each input document. The new array field
// contains the matching documents from the "joined" collection. The `$lookup` stage passes these reshaped documents to
// the next stage.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/
func Lookup(cond *lookup.JoinCondition) *Pipeline {
	return New().Lookup(cond)
}

// Match filters the documents to pass only the documents that match the specified condition(s) to the next pipeline
// stage.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/match/
func Match(filter *filters.Filter) *Pipeline {
	return New().Match(filter)
}

func Project(projection *projections.Options) *Pipeline {
	return New().Project(projection)
}

func Set(fields *bson.OrderedMap) *Pipeline {
	return New().Set(fields)
}

// ShardedDataDistribution returns information on the distribution of data in sharded collections.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/shardedDataDistribution/
func ShardedDataDistribution() *Pipeline {
	return New().ShardedDataDistribution()
}

// Skip skips over the specified number of documents that pass into the stage and passes the remaining documents to the
// next stage in the pipeline.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/skip/
func Skip(skip int64) *Pipeline {
	return New().Skip(skip)
}

// Sort sorts all input documents and returns them to the pipeline in sorted order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sort/
func Sort(opts *sorts.Options) *Pipeline {
	return New().Sort(opts)
}

// SortByCount Groups incoming documents based on the value of a specified expression, then computes the count of documents in
// each distinct group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sortByCount/
func SortByCount(expression any) *Pipeline {
	return New().SortByCount(expression)
}

// UnionWith Performs a union of two collections. `$unionWith` combines pipeline results from two collections into a
// single result set. The stage outputs the combined result set (including duplicates) to the next stage.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/unionWith/
func UnionWith(coll string, pipeline *Pipeline) *Pipeline {
	return New().UnionWith(coll, pipeline)
}

// Unset removes/excludes fields from documents.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/unset/
func Unset(fields ...string) *Pipeline {
	return New().Unset(fields...)
}

// Unwind deconstructs an array field from the input documents to output a document for each element. Each output
// document is the input document with the value of the array field replaced by the element.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/unwind/
func Unwind(opts *UnwindOptions) *Pipeline {
	return New().Unwind(opts)
}
