package maps

import (
	"errors"

	"github.com/qianwj/typed/adt"
	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/stream"
	"github.com/qianwj/typed/utils/json"
)

// TreeMap is an ordered map backed by an AVL tree. Lookup, insertion,
// removal, and neighbor queries take O(log n) time. Traversal follows the
// comparator supplied to NewTreeMap.
//
// The comparator must define a consistent total ordering: negative means
// before, zero means the same key, and positive means after. Replacing a
// comparator-equal key updates its value while retaining the original key.
// Keys must remain stable under the comparator while stored in the tree.
// K is comparable because Collect and MapStream collect into plain Go maps.
//
// Construct with NewTreeMap or TreeMapOf; the zero value is not usable.
// Do not copy a TreeMap after use. Concurrent mutation requires external
// synchronization, and traversal callbacks must not mutate the tree.
type TreeMap[K comparable, V any] struct {
	root    *treeNode[K, V]
	compare func(K, K) int
	size    int
}

type treeNode[K comparable, V any] struct {
	entry       Entry[K, V]
	left, right *treeNode[K, V]
	height      int
}

// NewTreeMap returns an empty ordered map. A nil comparator panics.
// For naturally ordered keys, pass cmp.Compare[K].
func NewTreeMap[K comparable, V any](compare func(K, K) int) *TreeMap[K, V] {
	if compare == nil {
		panic("maps.NewTreeMap: nil comparator")
	}
	return &TreeMap[K, V]{compare: compare}
}

// TreeMapOf copies entries into a new ordered map. Later values replace
// earlier ones for comparator-equal keys.
func TreeMapOf[K comparable, V any](compare func(K, K) int, entries ...Entry[K, V]) *TreeMap[K, V] {
	m := NewTreeMap[K, V](compare)
	for _, e := range entries {
		m.Put(e.Key, e.Value)
	}
	return m
}

// Put inserts or replaces a value and returns the previous value, or its
// zero value when absent. Use Contains to distinguish absence from zero.
func (m *TreeMap[K, V]) Put(key K, value V) V {
	return m.put(key, value, true)
}

// PutIfAbsent inserts value only if no comparator-equal key exists.
func (m *TreeMap[K, V]) PutIfAbsent(key K, value V) {
	m.put(key, value, false)
}

func (m *TreeMap[K, V]) put(key K, value V, replace bool) V {
	if m.compare == nil {
		panic("maps.TreeMap: construct with NewTreeMap")
	}
	var old V
	var added bool
	m.root, old, added = m.insert(m.root, key, value, replace)
	if added {
		m.size++
	}
	return old
}

func (m *TreeMap[K, V]) insert(n *treeNode[K, V], key K, value V, replace bool) (*treeNode[K, V], V, bool) {
	var old V
	if n == nil {
		return &treeNode[K, V]{entry: Entry[K, V]{Key: key, Value: value}, height: 1}, old, true
	}
	var added bool
	switch c := m.compare(key, n.entry.Key); {
	case c < 0:
		n.left, old, added = m.insert(n.left, key, value, replace)
	case c > 0:
		n.right, old, added = m.insert(n.right, key, value, replace)
	default:
		old = n.entry.Value
		if replace {
			n.entry.Value = value
		}
		return n, old, false
	}
	return balanceTree(n), old, added
}

// Get returns the value and whether a comparator-equal key exists.
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	for n := m.root; n != nil; {
		switch c := m.compare(key, n.entry.Key); {
		case c < 0:
			n = n.left
		case c > 0:
			n = n.right
		default:
			return n.entry.Value, true
		}
	}
	var zero V
	return zero, false
}

// GetOrDefault returns the value for key, or fallback when absent.
func (m *TreeMap[K, V]) GetOrDefault(key K, fallback V) V {
	if v, ok := m.Get(key); ok {
		return v
	}
	return fallback
}

// Contains reports whether a comparator-equal key exists.
func (m *TreeMap[K, V]) Contains(key K) bool {
	_, ok := m.Get(key)
	return ok
}

// Remove deletes key and returns its previous value and whether it existed.
func (m *TreeMap[K, V]) Remove(key K) (V, bool) {
	var old V
	var removed bool
	m.root, old, removed = m.remove(m.root, key)
	if removed {
		m.size--
	}
	return old, removed
}

func (m *TreeMap[K, V]) remove(n *treeNode[K, V], key K) (*treeNode[K, V], V, bool) {
	var old V
	if n == nil {
		return nil, old, false
	}
	var removed bool
	switch c := m.compare(key, n.entry.Key); {
	case c < 0:
		n.left, old, removed = m.remove(n.left, key)
	case c > 0:
		n.right, old, removed = m.remove(n.right, key)
	default:
		old, removed = n.entry.Value, true
		if n.left == nil {
			return n.right, old, true
		}
		if n.right == nil {
			return n.left, old, true
		}
		successor := n.right
		for successor.left != nil {
			successor = successor.left
		}
		n.entry = successor.entry
		n.right, _, _ = m.remove(n.right, successor.entry.Key)
	}
	return balanceTree(n), old, removed
}

// Size returns the number of entries.
func (m *TreeMap[K, V]) Size() int { return m.size }

// IsEmpty reports whether the tree has no entries.
func (m *TreeMap[K, V]) IsEmpty() bool { return m.size == 0 }

// Clear removes all entries while preserving the comparator.
func (m *TreeMap[K, V]) Clear() { m.root, m.size = nil, 0 }

// First returns the least entry under the comparator, or an empty Option.
func (m *TreeMap[K, V]) First() adt.Option[Entry[K, V]] {
	n := m.root
	if n == nil {
		return adt.Empty[Entry[K, V]]()
	}
	for n.left != nil {
		n = n.left
	}
	return adt.Of(n.entry)
}

// Last returns the greatest entry under the comparator, or an empty Option.
func (m *TreeMap[K, V]) Last() adt.Option[Entry[K, V]] {
	n := m.root
	if n == nil {
		return adt.Empty[Entry[K, V]]()
	}
	for n.right != nil {
		n = n.right
	}
	return adt.Of(n.entry)
}

// Floor returns the greatest entry with key <= key under the comparator.
func (m *TreeMap[K, V]) Floor(key K) adt.Option[Entry[K, V]] {
	return m.neighbor(key, true, true)
}

// Ceiling returns the least entry with key >= key under the comparator.
func (m *TreeMap[K, V]) Ceiling(key K) adt.Option[Entry[K, V]] {
	return m.neighbor(key, false, true)
}

// Lower returns the greatest entry strictly before key under the comparator.
func (m *TreeMap[K, V]) Lower(key K) adt.Option[Entry[K, V]] {
	return m.neighbor(key, true, false)
}

// Higher returns the least entry strictly after key under the comparator.
func (m *TreeMap[K, V]) Higher(key K) adt.Option[Entry[K, V]] {
	return m.neighbor(key, false, false)
}

func (m *TreeMap[K, V]) neighbor(key K, before, inclusive bool) adt.Option[Entry[K, V]] {
	var candidate *treeNode[K, V]
	for n := m.root; n != nil; {
		c := m.compare(n.entry.Key, key)
		if c == 0 && inclusive {
			return adt.Of(n.entry)
		}
		if before {
			if c < 0 {
				candidate, n = n, n.right
			} else {
				n = n.left
			}
		} else if c > 0 {
			candidate, n = n, n.left
		} else {
			n = n.right
		}
	}
	if candidate == nil {
		return adt.Empty[Entry[K, V]]()
	}
	return adt.Of(candidate.entry)
}

// ForEach visits entries in comparator order. visit must not mutate the tree.
func (m *TreeMap[K, V]) ForEach(visit func(K, V)) {
	var walk func(*treeNode[K, V])
	walk = func(n *treeNode[K, V]) {
		if n == nil {
			return
		}
		walk(n.left)
		visit(n.entry.Key, n.entry.Value)
		walk(n.right)
	}
	walk(m.root)
}

// Keys returns an independent list of keys in comparator order.
func (m *TreeMap[K, V]) Keys() *lists.ArrayList[K] {
	out := make([]K, 0, m.size)
	m.ForEach(func(k K, _ V) { out = append(out, k) })
	return lists.ArrayListOf(out...)
}

// Values returns an independent list of values in key order.
func (m *TreeMap[K, V]) Values() *lists.ArrayList[V] {
	out := make([]V, 0, m.size)
	m.ForEach(func(_ K, v V) { out = append(out, v) })
	return lists.ArrayListOf(out...)
}

// Entries returns an independent list of entries in comparator order.
func (m *TreeMap[K, V]) Entries() *lists.ArrayList[Entry[K, V]] {
	return lists.ArrayListOf(m.entries()...)
}

func (m *TreeMap[K, V]) entries() []Entry[K, V] {
	out := make([]Entry[K, V], 0, m.size)
	m.ForEach(func(k K, v V) { out = append(out, Entry[K, V]{Key: k, Value: v}) })
	return out
}

// Stream returns a snapshot of key-value pairs in comparator order. Later
// changes to the tree do not affect it. Collect returns a Go map without order.
func (m *TreeMap[K, V]) Stream() stream.MapStream[K, V] {
	entries := m.entries()
	return stream.FromPairs(func(yield func(K, V) bool) {
		for _, e := range entries {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	})
}

// Range returns an ordered snapshot of entries in [from, to) under the
// comparator. It takes O(log n + k) time and O(k) output space for k entries.
// Equal or reversed bounds return an empty stream.
func (m *TreeMap[K, V]) Range(from, to K) stream.Stream[Entry[K, V]] {
	if m.compare(from, to) >= 0 {
		return stream.Empty[Entry[K, V]]()
	}
	out := make([]Entry[K, V], 0)
	var walk func(*treeNode[K, V])
	walk = func(n *treeNode[K, V]) {
		if n == nil {
			return
		}
		lo, hi := m.compare(n.entry.Key, from), m.compare(n.entry.Key, to)
		if lo > 0 {
			walk(n.left)
		}
		if lo >= 0 && hi < 0 {
			out = append(out, n.entry)
		}
		if hi < 0 {
			walk(n.right)
		}
	}
	walk(m.root)
	return stream.FromSlice(out)
}

// Collect returns an independent Go map. Go maps do not preserve tree order;
// use Entries for ordered output.
func (m *TreeMap[K, V]) Collect() map[K]V {
	out := make(map[K]V, m.size)
	m.ForEach(func(k K, v V) { out[k] = v })
	return out
}

// Filter returns a new TreeMap with the same comparator and matching entries.
func (m *TreeMap[K, V]) Filter(p func(K, V) bool) *TreeMap[K, V] {
	out := NewTreeMap[K, V](m.compare)
	m.ForEach(func(k K, v V) {
		if p(k, v) {
			out.Put(k, v)
		}
	})
	return out
}

// MapValues transforms values into a new TreeMap with the same keys and comparator.
func (m *TreeMap[K, V]) MapValues[R any](f func(K, V) R) *TreeMap[K, R] {
	out := NewTreeMap[K, R](m.compare)
	m.ForEach(func(k K, v V) { out.Put(k, f(k, v)) })
	return out
}

// MarshalJSON encodes entries as a JSON object, like HashMap. JSON object
// member order is not guaranteed to follow the tree comparator.
func (m *TreeMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Encode(m.Collect()).Unwrap()
}

// UnmarshalJSON replaces entries, preserving the comparator. The receiver
// must be constructed with NewTreeMap first. Invalid input leaves it unchanged.
// JSON null clears the tree. Comparator-equal keys in an object are collapsed;
// when their Go keys differ, which value wins is unspecified.
func (m *TreeMap[K, V]) UnmarshalJSON(data []byte) error {
	if m.compare == nil {
		return errors.New("maps.TreeMap: construct with NewTreeMap before decoding JSON")
	}
	result := json.Decode[map[K]V](data)
	if err := result.Error(); err != nil {
		return err
	}
	out := NewTreeMap[K, V](m.compare)
	for k, v := range result.Value() {
		out.Put(k, v)
	}
	m.root, m.size = out.root, out.size
	return nil
}

func treeHeight[K comparable, V any](n *treeNode[K, V]) int {
	if n == nil {
		return 0
	}
	return n.height
}

func updateTreeHeight[K comparable, V any](n *treeNode[K, V]) {
	n.height = 1 + max(treeHeight(n.left), treeHeight(n.right))
}

func rotateTreeLeft[K comparable, V any](n *treeNode[K, V]) *treeNode[K, V] {
	top := n.right
	n.right, top.left = top.left, n
	updateTreeHeight(n)
	updateTreeHeight(top)
	return top
}

func rotateTreeRight[K comparable, V any](n *treeNode[K, V]) *treeNode[K, V] {
	top := n.left
	n.left, top.right = top.right, n
	updateTreeHeight(n)
	updateTreeHeight(top)
	return top
}

func balanceTree[K comparable, V any](n *treeNode[K, V]) *treeNode[K, V] {
	updateTreeHeight(n)
	balance := treeHeight(n.left) - treeHeight(n.right)
	if balance > 1 {
		if treeHeight(n.left.left) < treeHeight(n.left.right) {
			n.left = rotateTreeLeft(n.left)
		}
		return rotateTreeRight(n)
	}
	if balance < -1 {
		if treeHeight(n.right.right) < treeHeight(n.right.left) {
			n.right = rotateTreeRight(n.right)
		}
		return rotateTreeLeft(n)
	}
	return n
}
