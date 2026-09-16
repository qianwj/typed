package maps

import (
	"cmp"
	"encoding/json"
	stdmaps "maps"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/qianwj/typed/adt"
)

func checkTree(t *testing.T, m *TreeMap[int, int], want map[int]int) {
	t.Helper()
	seen := make(map[*treeNode[int, int]]bool)
	var walk func(*treeNode[int, int], *int, *int) int
	walk = func(n *treeNode[int, int], lower, upper *int) int {
		if n == nil {
			return 0
		}
		if seen[n] {
			t.Fatal("cycle or shared node in tree")
		}
		seen[n] = true
		if lower != nil && m.compare(n.entry.Key, *lower) <= 0 || upper != nil && m.compare(n.entry.Key, *upper) >= 0 {
			t.Fatalf("key %d violates tree ordering", n.entry.Key)
		}
		left := walk(n.left, lower, &n.entry.Key)
		right := walk(n.right, &n.entry.Key, upper)
		if left-right > 1 || right-left > 1 || n.height != 1+max(left, right) {
			t.Fatalf("key %d: left=%d right=%d height=%d", n.entry.Key, left, right, n.height)
		}
		return 1 + max(left, right)
	}
	walk(m.root, nil, nil)
	if len(seen) != len(want) || m.Size() != len(want) || m.IsEmpty() != (len(want) == 0) {
		t.Fatalf("node count=%d Size=%d, want %d", len(seen), m.Size(), len(want))
	}
	if got := m.Collect(); !stdmaps.Equal(got, want) {
		t.Fatalf("Collect = %v, want %v", got, want)
	}
	keys := m.Keys().Collect()
	wantKeys := slices.Collect(stdmaps.Keys(want))
	slices.SortFunc(wantKeys, m.compare)
	if !slices.Equal(keys, wantKeys) {
		t.Fatalf("Keys = %v, want %v", keys, wantKeys)
	}
}

func TestTreeMapRandomizedOperations(t *testing.T) {
	for _, descending := range []bool{false, true} {
		t.Run(strconv.FormatBool(descending), func(t *testing.T) {
			compare := cmp.Compare[int]
			if descending {
				compare = func(a, b int) int { return cmp.Compare(b, a) }
			}
			m := NewTreeMap[int, int](compare)
			want := make(map[int]int)
			rng := rand.New(rand.NewPCG(17, 29))
			for step := 0; step < 3000; step++ {
				key, value := rng.IntN(128), rng.IntN(1000)
				switch rng.IntN(4) {
				case 0:
					if old := m.Put(key, value); old != want[key] {
						t.Fatalf("Put old = %d, want %d", old, want[key])
					}
					want[key] = value
				case 1:
					old, existed := want[key]
					if got, ok := m.Remove(key); got != old || ok != existed {
						t.Fatalf("Remove(%d) = (%d,%v), want (%d,%v)", key, got, ok, old, existed)
					}
					delete(want, key)
				case 2:
					m.PutIfAbsent(key, value)
					if _, ok := want[key]; !ok {
						want[key] = value
					}
				case 3:
					value, exists := want[key]
					if got, ok := m.Get(key); got != value || ok != exists || m.Contains(key) != exists {
						t.Fatalf("Get(%d) = (%d,%v), want (%d,%v)", key, got, ok, value, exists)
					}
				}
				checkTree(t, m, want)
			}
			// Remove every remaining key, checking all intermediate AVL states.
			for _, key := range rng.Perm(128) {
				m.Remove(key)
				delete(want, key)
				checkTree(t, m, want)
			}
		})
	}
}

func TestTreeMapOrderedInsertAndDelete(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		m := NewTreeMap[int, int](cmp.Compare[int])
		want := make(map[int]int)
		for i := 0; i < 256; i++ {
			key := i
			if reverse {
				key = 255 - i
			}
			m.Put(key, key)
			want[key] = key
			checkTree(t, m, want)
		}
		for i := 0; i < 256; i++ {
			// Repeatedly deleting the root exercises two-child replacement.
			key := m.root.entry.Key
			if old, ok := m.Remove(key); !ok || old != key {
				t.Fatalf("Remove root %d = (%d,%v)", key, old, ok)
			}
			delete(want, key)
			checkTree(t, m, want)
		}
	}
}

func TestTreeMapComparatorEquality(t *testing.T) {
	compare := func(a, b string) int { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) }
	m := TreeMapOf(compare, Entry[string, int]{Key: "Alpha", Value: 0})
	m.PutIfAbsent("ALPHA", 7)
	if got, ok := m.Get("alpha"); !ok || got != 0 {
		t.Fatalf("PutIfAbsent overwrote zero value: (%d,%v)", got, ok)
	}
	if old := m.Put("aLpHa", 9); old != 0 || m.Size() != 1 {
		t.Fatalf("Put comparator-equal key: old=%d size=%d", old, m.Size())
	}
	if e := m.First().Get(); e.Key != "Alpha" || e.Value != 9 {
		t.Fatalf("replacement should retain stored key: %+v", e)
	}
	if got := m.GetOrDefault("missing", 42); got != 42 {
		t.Fatalf("fallback = %d", got)
	}
	if got := m.GetOrDefault("alpha", 42); got != 9 {
		t.Fatalf("existing = %d", got)
	}
	if got, ok := m.Remove("ALPHA"); !ok || got != 9 || !m.IsEmpty() {
		t.Fatalf("Remove comparator-equal key = (%d,%v)", got, ok)
	}
}

func TestTreeMapNeighbors(t *testing.T) {
	for _, compare := range []func(int, int) int{cmp.Compare[int], func(a, b int) int { return cmp.Compare(b, a) }} {
		m := NewTreeMap[int, int](compare)
		for _, query := range []func(int) adt.Option[Entry[int, int]]{m.Floor, m.Ceiling, m.Lower, m.Higher} {
			if !query(0).IsEmpty() {
				t.Fatal("neighbor on empty tree must be absent")
			}
		}
		if !m.First().IsEmpty() || !m.Last().IsEmpty() {
			t.Fatal("empty extrema must be absent")
		}
		for _, k := range []int{0, 2, 4, 6, 8} {
			m.Put(k, k*10)
		}
		keys := m.Keys().Collect()
		if m.First().Get().Key != keys[0] || m.Last().Get().Key != keys[len(keys)-1] {
			t.Fatal("First/Last do not respect comparator")
		}
		for key := -1; key <= 9; key++ {
			for _, tc := range []struct {
				name              string
				query             func(int) adt.Option[Entry[int, int]]
				before, inclusive bool
			}{{"Floor", m.Floor, true, true}, {"Lower", m.Lower, true, false}, {"Ceiling", m.Ceiling, false, true}, {"Higher", m.Higher, false, false}} {
				var candidates []int
				for _, k := range keys {
					c := compare(k, key)
					if tc.inclusive && c == 0 || tc.before && c < 0 || !tc.before && c > 0 {
						candidates = append(candidates, k)
					}
				}
				got := tc.query(key)
				if len(candidates) == 0 {
					if !got.IsEmpty() {
						t.Fatalf("%s(%d) should be absent", tc.name, key)
					}
					continue
				}
				want := candidates[0]
				if tc.before {
					want = candidates[len(candidates)-1]
				}
				if got.IsEmpty() || got.Get() != (Entry[int, int]{Key: want, Value: want * 10}) {
					t.Fatalf("%s(%d) = %v, want key %d", tc.name, key, got, want)
				}
			}
		}
	}
}

func TestTreeMapRangeAndSnapshots(t *testing.T) {
	for _, compare := range []func(int, int) int{cmp.Compare[int], func(a, b int) int { return cmp.Compare(b, a) }} {
		m := NewTreeMap[int, int](compare)
		for i := 0; i < 10; i++ {
			m.Put(i, i*10)
		}
		all := m.Entries().Collect()
		for from := -1; from <= 10; from++ {
			for to := -1; to <= 10; to++ {
				var want []Entry[int, int]
				for _, e := range all {
					if compare(e.Key, from) >= 0 && compare(e.Key, to) < 0 {
						want = append(want, e)
					}
				}
				if got := m.Range(from, to).Collect(); !slices.Equal(got, want) {
					t.Fatalf("Range(%d,%d) = %v, want %v", from, to, got, want)
				}
			}
		}
	}
	m := TreeMapOf(cmp.Compare[int], Entry[int, string]{Key: 1, Value: "one"}, Entry[int, string]{Key: 2, Value: "two"})
	snapshot, ranged := m.Stream(), m.Range(1, 3)
	keys, values, entries, plain := m.Keys(), m.Values(), m.Entries(), m.Collect()
	m.Clear()
	m.Put(3, "three")
	want := []Entry[int, string]{{Key: 1, Value: "one"}, {Key: 2, Value: "two"}}
	if !stdmaps.Equal(snapshot.Collect(), plain) || !slices.Equal(ranged.Collect(), want) || !slices.Equal(entries.Collect(), want) {
		t.Fatal("mutations affected entry snapshots")
	}
	if !slices.Equal(keys.Collect(), []int{1, 2}) || !slices.Equal(values.Collect(), []string{"one", "two"}) || plain[1] != "one" {
		t.Fatal("mutations affected key/value/map snapshots")
	}
	plain[3] = "changed"
	if m.GetOrDefault(3, "") != "three" {
		t.Fatal("Collect exposes tree storage")
	}
}

func TestTreeMapStream(t *testing.T) {
	m := TreeMapOf(func(a, b int) int { return cmp.Compare(b, a) },
		Entry[int, string]{Key: 1, Value: "one"},
		Entry[int, string]{Key: 3, Value: "three"},
		Entry[int, string]{Key: 2, Value: "two"})
	snapshot := m.Stream()
	var visited []int
	out := snapshot.Filter(func(k int, _ string) bool {
		visited = append(visited, k)
		return k > 1
	}).MapValues(func(k int, v string) int { return k + len(v) })
	if len(visited) != 0 {
		t.Fatal("stream evaluated before collection")
	}
	m.Put(3, "changed")
	m.Remove(2)
	got := out.Collect()
	if !stdmaps.Equal(got, map[int]int{3: 8, 2: 5}) || !slices.Equal(visited, []int{3, 2, 1}) {
		t.Fatalf("stream lost snapshot contents or comparator order: %v, visited %v", got, visited)
	}
	if keys := m.Stream().Map(func(k int, _ string) int { return k }).Collect(); !slices.Equal(keys, []int{3, 1}) {
		t.Fatalf("Map lost comparator order: %v", keys)
	}
	if empty := NewTreeMap[int, string](cmp.Compare[int]).Stream().Collect(); empty == nil || len(empty) != 0 {
		t.Fatalf("empty stream collected to %v", empty)
	}
}

func TestTreeMapTransforms(t *testing.T) {
	compare := func(a, b int) int { return cmp.Compare(b, a) }
	m := TreeMapOf(compare, Entry[int, int]{Key: 1, Value: 10}, Entry[int, int]{Key: 2, Value: 20}, Entry[int, int]{Key: 3, Value: 30})
	out := m.Filter(func(k, v int) bool { return k > 1 && v >= 20 }).
		MapValues(func(k, v int) []int { return []int{k, v} })
	if !slices.Equal(out.Keys().Collect(), []int{3, 2}) || !slices.Equal(out.First().Get().Value, []int{3, 30}) {
		t.Fatalf("transforms changed comparator or values: %v", out.Collect())
	}
	out.Remove(3)
	if m.Size() != 3 || !m.Contains(3) {
		t.Fatal("transforms share mutable tree nodes")
	}
}

func TestTreeMapJSONAndConstruction(t *testing.T) {
	t.Run("nil comparator", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("nil comparator must panic")
			}
		}()
		NewTreeMap[int, int](nil)
	})
	t.Run("zero value write", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("zero value Put must panic")
			}
		}()
		var m TreeMap[int, int]
		m.Put(1, 1)
	})
	m := NewTreeMap[int, string](func(a, b int) int { return cmp.Compare(b, a) })
	if err := json.Unmarshal([]byte(`{"1":"one","2":"two"}`), m); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(m.Keys().Collect(), []int{2, 1}) {
		t.Fatal("decoding lost comparator")
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var plain map[int]string
	if err := json.Unmarshal(data, &plain); err != nil {
		t.Fatal(err)
	}
	if !stdmaps.Equal(plain, map[int]string{1: "one", 2: "two"}) {
		t.Fatalf("JSON = %s", data)
	}
	if err := m.UnmarshalJSON([]byte(`{"3":123}`)); err == nil || m.Size() != 2 || !m.Contains(1) {
		t.Fatal("invalid decode must preserve the tree")
	}
	if err := m.UnmarshalJSON([]byte(`null`)); err != nil || !m.IsEmpty() {
		t.Fatal("null must clear tree")
	}
	m.Put(1, "one")
	m.Put(2, "two")
	if m.First().Get().Key != 2 {
		t.Fatal("null decode lost comparator")
	}
	var zero TreeMap[int, string]
	if err := zero.UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatal("uninitialized decode must fail")
	}
}
