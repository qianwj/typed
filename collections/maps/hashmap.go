// Package maps provides HashMap[K, V], a type-safe unordered collection of
// key-value pairs built on top of the built-in map.
//
// As with the other concrete generic collections in this module,
// HashMap's methods (including type-changing ones such as MapValues[R])
// take advantage of Go 1.27's generic methods. An interface carrying
// MapValues would not be able to declare its own R parameter, so the
// type is exposed as a concrete generic type rather than an interface.
//
// All methods use pointer receivers (*HashMap[K, V]). Copying a HashMap
// value would share the underlying map and lead to surprising mutations;
// the pointer-receiver convention makes the "do not copy" rule uniform
// across the package.
//
// HashMap is eager: intermediate operations return new HashMaps
// immediately. For lazy pipelines over the entries, call Stream() to
// obtain a stream.Stream[Entry[K, V]].
package maps

import (
	"maps"

	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/stream"
	"github.com/qianwj/typed/utils/json"
)

// Entry is a single key-value pair as yielded by HashMap.Entries and
// HashMap.Stream. Key matches the comparable constraint of the parent
// HashMap, Value is unconstrained.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// HashMap is an unordered collection of K -> V pairs backed by a private
// map[K]V.
//
// The underlying map is not exported; callers cannot index into it or
// otherwise bypass the API. All reads and writes go through the methods
// defined on *HashMap. To obtain a plain map[K]V, call Collect (which
// returns a copy).
//
// As a concrete generic type, HashMap's methods (including type-changing
// ones such as MapValues[R]) can use Go 1.27's generic methods.
type HashMap[K comparable, V any] struct {
	items map[K]V
}

// ---------- Constructors ----------

// NewHashMap returns a new empty HashMap.
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{items: make(map[K]V)}
}

// HashMapOf returns a HashMap containing the given entries.
//
// Later entries with the same key overwrite earlier ones, matching the
// behavior of map[K]V assignment. The entries are copied into a fresh
// underlying map; mutating the input slice afterwards does not affect
// the returned HashMap.
func HashMapOf[K comparable, V any](entries ...Entry[K, V]) *HashMap[K, V] {
	m := &HashMap[K, V]{items: make(map[K]V, len(entries))}
	for _, e := range entries {
		m.items[e.Key] = e.Value
	}
	return m
}

// HashMapFromMap returns a HashMap containing the entries of src. The
// entries are copied; subsequent mutations of src do not affect the
// returned HashMap, and vice versa.
func HashMapFromMap[K comparable, V any](src map[K]V) *HashMap[K, V] {
	m := &HashMap[K, V]{items: make(map[K]V, len(src))}
	maps.Copy(m.items, src)
	return m
}

// ---------- Basic operations ----------

// Put inserts or updates the value for key and returns the previous value
// for that key, or the zero value of V if the key was not present.
//
// Because V may legitimately hold its zero value, the return cannot be
// used to test whether the key existed; use Contains for that.
func (m *HashMap[K, V]) Put(key K, value V) V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	old, _ := m.items[key]
	m.items[key] = value
	return old
}

// PutIfAbsent inserts value for key only if the key is not already
// present. It is a no-op when the key exists.
func (m *HashMap[K, V]) PutIfAbsent(key K, value V) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	_, ok := m.items[key]
	if !ok {
		m.items[key] = value
	}
}

// Get returns the value for key, or the zero value and false if the key
// is absent.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.items[key]
	return v, ok
}

// GetOrDefault returns the value for key, or defaultValue if the key is
// not present. The HashMap is not mutated.
func (m *HashMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	v, ok := m.items[key]
	if !ok {
		return defaultValue
	}
	return v
}

// Remove deletes key from the HashMap, returning the previous value and
// true if the key was present.
func (m *HashMap[K, V]) Remove(key K) (V, bool) {
	v, ok := m.items[key]
	if ok {
		delete(m.items, key)
	}
	return v, ok
}

// Contains reports whether key is present in the HashMap.
func (m *HashMap[K, V]) Contains(key K) bool {
	_, ok := m.items[key]
	return ok
}

// Size returns the number of entries.
func (m *HashMap[K, V]) Size() int {
	return len(m.items)
}

// IsEmpty reports whether the HashMap has no entries.
func (m *HashMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear removes all entries.
func (m *HashMap[K, V]) Clear() {
	clear(m.items)
}

// ---------- Iteration and projection ----------

// ForEach invokes visit on every (key, value) pair. Iteration order is
// unspecified, matching the underlying map's semantics.
func (m *HashMap[K, V]) ForEach(visit func(K, V)) {
	for k, v := range m.items {
		visit(k, v)
	}
}

// Keys returns an ArrayList of all keys. Iteration order is unspecified.
func (m *HashMap[K, V]) Keys() *lists.ArrayList[K] {
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return lists.ArrayListOf(keys...)
}

// Values returns an ArrayList of all values. Iteration order is
// unspecified.
func (m *HashMap[K, V]) Values() *lists.ArrayList[V] {
	vals := make([]V, 0, len(m.items))
	for _, v := range m.items {
		vals = append(vals, v)
	}
	return lists.ArrayListOf(vals...)
}

// Entries returns an ArrayList of all entries as Entry[K, V] values.
// Iteration order is unspecified.
func (m *HashMap[K, V]) Entries() *lists.ArrayList[Entry[K, V]] {
	entries := make([]Entry[K, V], 0, len(m.items))
	for k, v := range m.items {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return lists.ArrayListOf(entries...)
}

// ---------- Transformations ----------

// Filter returns a new HashMap keeping only the entries for which p
// returns true.
func (m *HashMap[K, V]) Filter(p func(K, V) bool) *HashMap[K, V] {
	out := &HashMap[K, V]{items: make(map[K]V, len(m.items))}
	for k, v := range m.items {
		if p(k, v) {
			out.items[k] = v
		}
	}
	return out
}

// FilterKeys returns a new HashMap keeping only the entries whose key
// satisfies p.
func (m *HashMap[K, V]) FilterKeys(p func(K) bool) *HashMap[K, V] {
	out := &HashMap[K, V]{items: make(map[K]V, len(m.items))}
	for k, v := range m.items {
		if p(k) {
			out.items[k] = v
		}
	}
	return out
}

// FilterValues returns a new HashMap keeping only the entries whose
// value satisfies p.
func (m *HashMap[K, V]) FilterValues(p func(V) bool) *HashMap[K, V] {
	out := &HashMap[K, V]{items: make(map[K]V, len(m.items))}
	for k, v := range m.items {
		if p(v) {
			out.items[k] = v
		}
	}
	return out
}

// MapValues applies f to every value and returns a new HashMap with the
// same keys but transformed values. This is a type-changing method, only
// possible because HashMap is a concrete generic type in Go 1.27.
func (m *HashMap[K, V]) MapValues[R any](f func(K, V) R) *HashMap[K, R] {
	out := &HashMap[K, R]{items: make(map[K]R, len(m.items))}
	for k, v := range m.items {
		out.items[k] = f(k, v)
	}
	return out
}

// Concat returns a new HashMap that merges other into m. Entries from
// other overwrite entries with the same key in m. Neither input is
// mutated.
func (m *HashMap[K, V]) Concat(other *HashMap[K, V]) *HashMap[K, V] {
	out := &HashMap[K, V]{items: make(map[K]V, len(m.items)+len(other.items))}
	maps.Copy(out.items, m.items)
	maps.Copy(out.items, other.items)
	return out
}

// ---------- Stream and collect ----------

// Stream returns a lazy stream.Stream[Entry[K, V]] that is a snapshot of
// the HashMap at the time Stream() is called. Mutations to the HashMap
// after Stream() do not affect the Stream.
//
// To make the snapshot contract hold, Stream materializes the entries
// into a fresh backing slice. The copy is O(n); iteration of the Stream
// is then the same cost as iterating a plain slice. This cost is paid
// once at the HashMap -> Stream boundary, after which the Stream itself
// remains lazy.
//
// The Stream is single-use.
func (m *HashMap[K, V]) Stream() stream.Stream[Entry[K, V]] {
	buf := make([]Entry[K, V], 0, len(m.items))
	for k, v := range m.items {
		buf = append(buf, Entry[K, V]{Key: k, Value: v})
	}
	return stream.FromSlice(buf)
}

// Collect returns a freshly allocated map[K]V containing the same
// entries as the HashMap. The returned map is decoupled from the
// HashMap's internal storage; mutating it does not affect the HashMap.
func (m *HashMap[K, V]) Collect() map[K]V {
	out := make(map[K]V, len(m.items))
	maps.Copy(out, m.items)
	return out
}

// MarshalJSON encodes the HashMap as a JSON object whose keys
// are the K values and whose values are the corresponding V
// values. The output is identical to marshaling Collect().
//
// Because Go's map iteration order is unspecified, the order
// of keys in the marshaled object is not stable across runs;
// this matches the v1 json package's behaviour for plain
// map[K]V values and is documented for users who need
// stable serialisation.
//
// MarshalJSON delegates to utils/json.Encode, the project-
// wide wrapper around encoding/json/v2. Keys and values that
// satisfy json.Marshaler (or v2's marshaler variant) are
// encoded by their respective MarshalJSON methods; otherwise
// the default encoding for K and V applies.
func (m *HashMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Encode(m.items).Unwrap()
}

// UnmarshalJSON decodes a JSON object into the HashMap,
// replacing any existing contents. Each key/value pair in
// the object becomes an entry in the HashMap; the previous
// entries are dropped.
//
// UnmarshalJSON delegates to utils/json.Decode and plumbs
// the resulting map[K]V directly into the HashMap's internal
// storage. JSON null is accepted and treated as an empty
// object: a nil-decoded HashMap is empty (items == nil).
// Decoding into a key type that does not support a particular
// JSON value (for example, decoding a JSON string into an int
// key) returns a type-mismatch error.
//
// The return value is the underlying v2 error if data is not
// a JSON object or if any key or value fails to decode into
// K or V.
func (m *HashMap[K, V]) UnmarshalJSON(data []byte) error {
	r := json.Decode[map[K]V](data)
	if err := r.Error(); err != nil {
		return err
	}
	m.items = r.Value()
	return nil
}
