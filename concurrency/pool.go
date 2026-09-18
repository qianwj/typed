package concurrency

import (
	"sync"

	"github.com/qianwj/typed/adt/option"
	"github.com/qianwj/typed/adt"
)

// Pool caches temporary values for reuse, reducing allocation and GC pressure.
// It is a typed wrapper around sync.Pool. Get and Put are safe for concurrent
// use; this does not make the values themselves safe for concurrent mutation.
// Prefer pointer types such as *bytes.Buffer to avoid boxing value types.
//
// Cached values may be discarded at any time without notification. Pool has
// no capacity limit and does not wait for values to be returned. Use it for
// disposable temporary objects, not resources requiring explicit cleanup.
//
// The zero value is usable and returns an empty Option on a cache miss.
// Use NewPool to supply a creator instead. Do not copy a Pool after first use.
type Pool[T any] struct {
	container sync.Pool
}

// NewPool returns an empty Pool that calls creator on a cache miss.
// A nil creator panics. Construction does not call creator; concurrent Get
// calls may invoke it concurrently. Each call should create an independent
// value when T holds mutable state. Nil results are allowed and produce an
// empty Option, including typed nil results.
func NewPool[T any](creator func() T) *Pool[T] {
	if creator == nil {
		panic("concurrency: NewPool requires a non-nil creator")
	}
	return &Pool[T]{
		container: sync.Pool{
			New: func() any { return creator() },
		},
	}
}

// Get returns an Option containing a cached value, or calls the creator on a
// cache miss. A miss without a creator, or a nil result (including a typed nil
// pointer, slice, map, channel, or function), returns an empty Option.
// Non-nil zero values such as 0, false, and "" produce a present Option.
// A cached typed nil yields an empty Option without retrying the creator.
// There is no guarantee that Get returns a previously stored value, even
// immediately after Put.
// Get does not reset values; the caller must prepare them for reuse.
func (p *Pool[T]) Get() option.Option[T] {
	value := p.container.Get()
	return option.OfNullable(value).
		FlatMap(func(v any) option.Option[T] {
			return adt.Cast[T](v).Option()
		})
}

// Put makes value available for reuse without resetting it. After returning a
// mutable value, stop using it and any aliases to its mutable storage, and do
// not return it again without first acquiring it. A nil interface is ignored;
// typed nil values follow sync.Pool's normal storage semantics, but Get
// represents them as an empty Option.
//
// Put(value) synchronizes before a Get that retrieves that same value.
func (p *Pool[T]) Put(value T) {
	p.container.Put(value)
}
