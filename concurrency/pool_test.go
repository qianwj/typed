package concurrency

import (
	"bytes"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPoolCreator(t *testing.T) {
	calls := 0
	p := NewPool(func() *int {
		calls++
		value := 42
		return &value
	})
	if calls != 0 {
		t.Fatal("construction called creator")
	}
	a, b := p.Get().Get(), p.Get().Get()
	if calls != 2 || a == b || *a != 42 || *b != 42 {
		t.Fatalf("creator calls=%d, values=%v, %v", calls, a, b)
	}
}

func TestPoolNilCreatorPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected nil creator to panic")
		}
	}()
	NewPool[int](nil)
}

func TestPoolZeroValue(t *testing.T) {
	var p Pool[int]
	if got := p.Get(); got.IsPresent() {
		t.Fatalf("zero Pool.Get() = %v, want empty Option", got)
	}
	p.Put(42)
	if got := p.Get(); got.IsPresent() && got.Get() != 42 {
		t.Fatalf("Get() = %v, want cached value or empty Option on miss", got)
	}
}

func TestPoolNilValues(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		p := NewPool(func() any { return nil })
		p.Put(nil)
		if got := p.Get(); got.IsPresent() {
			t.Fatalf("Get() = %v, want empty Option", got)
		}
	})
	t.Run("nil pointer", func(t *testing.T) {
		p := NewPool(func() *int { return nil })
		p.Put(nil)
		if got := p.Get(); got.IsPresent() {
			t.Fatalf("Get() = %v, want empty Option for nil pointer", got)
		}
	})
	t.Run("typed nil interface", func(t *testing.T) {
		p := NewPool(func() any { return (*int)(nil) })
		p.Put((*int)(nil))
		if got := p.Get(); got.IsPresent() {
			t.Fatalf("Get() = %v, want empty Option for typed nil", got)
		}
	})
	for _, tc := range []struct {
		name  string
		value any
	}{
		{"nil slice", []int(nil)},
		{"nil map", map[string]int(nil)},
		{"nil channel", (chan int)(nil)},
		{"nil function", (func())(nil)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPool(func() any { return tc.value })
			if got := p.Get(); got.IsPresent() {
				t.Fatalf("creator result %T should produce empty Option", tc.value)
			}
			p.Put(tc.value)
			if got := p.Get(); got.IsPresent() {
				t.Fatalf("Get after Put(%T) should produce empty Option", tc.value)
			}
		})
	}
}

func TestPoolZeroValuesArePresent(t *testing.T) {
	for _, value := range []any{0, false, ""} {
		p := NewPool(func() any { return value })
		if got := p.Get(); got.IsEmpty() || got.Get() != value {
			t.Fatalf("Get() = %v, want present %v", got, value)
		}
	}
	p := NewPool(func() int { return 0 })
	if got := p.Get(); got.IsEmpty() || got.Get() != 0 {
		t.Fatalf("Get() = %v, want present zero", got)
	}
}

func TestPoolEmptySliceIsPresent(t *testing.T) {
	p := NewPool(func() []int { return []int{} })
	if got := p.Get(); got.IsEmpty() || got.Get() == nil || len(got.Get()) != 0 {
		t.Fatalf("Get() = %v, want present empty slice", got)
	}
}

func TestPoolSliceValues(t *testing.T) {
	p := NewPool(func() []int { return []int{1, 2} })
	p.Put([]int{3, 4})
	got := p.Get().Get()
	if !slices.Equal(got, []int{1, 2}) && !slices.Equal(got, []int{3, 4}) {
		t.Fatalf("Get() = %v, want cached or newly created slice", got)
	}
}

func TestPoolAfterGC(t *testing.T) {
	p := NewPool(func() *bytes.Buffer { return new(bytes.Buffer) })
	p.Put(p.Get().Get())
	runtime.GC()
	runtime.GC()
	got := p.Get()
	if got.IsEmpty() {
		t.Fatal("Get returned empty Option after GC despite non-nil creator")
	}
	buf := got.Get()
	buf.Reset()
	buf.WriteString("reusable")
	if buf.String() != "reusable" {
		t.Fatal("buffer unusable after GC")
	}
	p.Put(buf)
}

func TestPoolConcurrent(t *testing.T) {
	type item struct {
		owned atomic.Bool
		value int
	}
	p := NewPool(func() *item { return new(item) })
	var wg sync.WaitGroup
	for worker := 0; worker < 16; worker++ {
		wg.Go(func() {
			for i := 0; i < 200; i++ {
				got := p.Get()
				if got.IsEmpty() {
					t.Error("Get returned empty Option despite non-nil creator")
					return
				}
				value := got.Get()
				if !value.owned.CompareAndSwap(false, true) {
					t.Error("same object acquired by multiple goroutines")
					return
				}
				value.value = worker
				runtime.Gosched()
				if value.value != worker {
					t.Error("object changed while acquired")
				}
				value.owned.Store(false)
				p.Put(value)
			}
		})
	}
	wg.Wait()
}
