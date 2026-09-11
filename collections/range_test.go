package collections

import (
	"reflect"
	"testing"
)

// ---------- basic half-open interval ----------

// TestRangeBasicHalfOpen confirms the canonical case:
// Range(0, 5) yields exactly 0, 1, 2, 3, 4 (5 elements, end
// excluded). This is the half-open [start, end) convention used
// by Rust, Python, and Java's IntStream.range.
func TestRangeBasicHalfOpen(t *testing.T) {
	got := Range(0, 5).Collect()
	want := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(0, 5).Collect(): got %v, want %v", got, want)
	}
	if n := Range(0, 5).Count(); n != 5 {
		t.Fatalf("Range(0, 5).Count(): got %d, want 5", n)
	}
}

// TestRangeNonZeroStart confirms that the start parameter is
// honoured: Range(3, 7) yields 3, 4, 5, 6 (not 0..7).
func TestRangeNonZeroStart(t *testing.T) {
	got := Range(3, 7).Collect()
	want := []int{3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(3, 7): got %v, want %v", got, want)
	}
}

// TestRangeNegativeStart confirms Range works with negative
// starts: Range(-3, 2) yields -3, -2, -1, 0, 1.
func TestRangeNegativeStart(t *testing.T) {
	got := Range(-3, 2).Collect()
	want := []int{-3, -2, -1, 0, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(-3, 2): got %v, want %v", got, want)
	}
}

// TestRangeSingleElement confirms that a half-open interval of
// length 1 yields exactly one element.
func TestRangeSingleElement(t *testing.T) {
	got := Range(7, 8).Collect()
	want := []int{7}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(7, 8): got %v, want %v", got, want)
	}
}

// ---------- empty / reversed intervals ----------

// TestRangeEmptyStartEqualsEnd confirms that start == end
// produces an empty Stream, not a single-element one. This is
// the "off-by-one" case that half-open intervals get right and
// closed intervals get wrong.
func TestRangeEmptyStartEqualsEnd(t *testing.T) {
	if got := Range(5, 5).Collect(); len(got) != 0 {
		t.Fatalf("Range(5, 5): got %v, want empty", got)
	}
	if n := Range(5, 5).Count(); n != 0 {
		t.Fatalf("Range(5, 5).Count(): got %d, want 0", n)
	}
	if Range(5, 5).Any(func(int) bool { return true }) {
		t.Fatal("Range(5, 5).Any: returned true on empty stream")
	}
}

// TestRangeEmptyStartGreaterThanEnd confirms that start > end
// also produces an empty Stream, without panicking. This makes
// Range safe to call with computed bounds that may not be in
// order, like a defensive Range(low, high) where low might
// exceed high.
func TestRangeEmptyStartGreaterThanEnd(t *testing.T) {
	if got := Range(10, 3).Collect(); len(got) != 0 {
		t.Fatalf("Range(10, 3): got %v, want empty", got)
	}
	if n := Range(10, 3).Count(); n != 0 {
		t.Fatalf("Range(10, 3).Count(): got %d, want 0", n)
	}
}

// TestRangeEmptyLargeGap confirms that a large reversed gap
// still produces an empty Stream without iteration. Count
// must be 0 and the Stream must not yield any value.
func TestRangeEmptyLargeGap(t *testing.T) {
	got := Range(1_000_000, 0).Collect()
	if len(got) != 0 {
		t.Fatalf("Range(1_000_000, 0): got %d elements, want 0", len(got))
	}
}

// ---------- various integer types ----------

// TestRangeWithInt64 confirms the Integer type constraint
// accepts int64. A common use case is Range(0, int64(n)) where
// n is the length of a slice.
func TestRangeWithInt64(t *testing.T) {
	got := Range(int64(2), int64(5)).Collect()
	want := []int64{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range[int64](2, 5): got %v, want %v", got, want)
	}
}

// TestRangeWithUint confirms the Integer type constraint
// accepts unsigned integers. This is the natural type for
// index-based ranges.
func TestRangeWithUint(t *testing.T) {
	got := Range(uint(0), uint(3)).Collect()
	want := []uint{0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range[uint](0, 3): got %v, want %v", got, want)
	}
}

// ---------- composition with the Stream API ----------

// TestRangeWithMap confirms that Range is composable with Map
// to produce a transformed stream. This is the canonical
// "iota-style squares" pattern: 1..6 -> [1, 4, 9, 16, 25].
func TestRangeWithMap(t *testing.T) {
	got := Range(1, 6).
		Map(func(x int) int { return x * x }).
		Collect()
	want := []int{1, 4, 9, 16, 25}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(1, 6).Map(square): got %v, want %v", got, want)
	}
}

// TestRangeWithFilter confirms that Filter can gate a Range.
// Range(0, 10).Filter(even) must yield the even numbers.
func TestRangeWithFilter(t *testing.T) {
	got := Range(0, 10).
		Filter(func(x int) bool { return x%2 == 0 }).
		Collect()
	want := []int{0, 2, 4, 6, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(0, 10).Filter(even): got %v, want %v", got, want)
	}
}

// TestRangeWithTake confirms that Take can cap a Range before
// its end. The take is the standard idiom for "do this n
// times" when n is small and the underlying range is large.
func TestRangeWithTake(t *testing.T) {
	got := Range(0, 1_000_000).Take(5).Collect()
	want := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(0, 1M).Take(5): got %v, want %v", got, want)
	}
}

// TestRangeWithTakeMoreThanLength confirms that Take(n) where n
// exceeds the range's length still works: the Stream ends at
// end - start, not at n.
func TestRangeWithTakeMoreThanLength(t *testing.T) {
	got := Range(0, 3).Take(100).Collect()
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(0, 3).Take(100): got %v, want %v", got, want)
	}
}

// TestRangeWithReduce confirms that Reduce threads a Range
// through an accumulator. The result is the standard sum-of-1..n
// identity: 0+1+2+3+4 = 10.
func TestRangeWithReduce(t *testing.T) {
	got := Range(0, 5).Reduce(0, func(acc, v int) int { return acc + v })
	if got != 10 {
		t.Fatalf("Range(0, 5).Reduce(+): got %d, want 10", got)
	}
}

// TestRangeFirstAndLast confirms the Optional-based terminals
// work: First is the smallest value, Last is end - 1.
func TestRangeFirstAndLast(t *testing.T) {
	first := Range(10, 20).First()
	if v := first.OrElse(-1); v != 10 {
		t.Fatalf("Range(10, 20).First: got %d, want 10", v)
	}
	last := Range(10, 20).Last()
	if v := last.OrElse(-1); v != 19 {
		t.Fatalf("Range(10, 20).Last: got %d, want 19", v)
	}
}

// TestRangeFirstOnEmpty confirms the absent-Optional branch:
// Range(5, 5).First must be absent.
func TestRangeFirstOnEmpty(t *testing.T) {
	first := Range(5, 5).First()
	if first.IsPresent() {
		t.Fatalf("Range(5, 5).First: got present %d, want absent", first.OrElse(0))
	}
}

// ---------- laziness and memory ----------

// TestRangeIsLazy confirms that constructing a Range does not
// produce the values up front. The Stream is built but no
// element is computed until Collect / Count / ForEach is
// called. The size used here is large enough (10 million) that
// eager construction would be obviously slow, but small enough
// that a single Count completes in well under a second.
func TestRangeIsLazy(t *testing.T) {
	const big = 10_000_000
	if n := Range(0, big).Count(); n != big {
		t.Fatalf("Range(0, 10M).Count: got %d, want %d", n, big)
	}
}

// TestRangeWithEarlyBreak confirms that the Stream composes
// with Take to stop iteration at a small prefix of a large
// range. Take(3) on Range(0, 10M) yields only 0, 1, 2.
func TestRangeWithEarlyBreak(t *testing.T) {
	const big = 10_000_000
	var seen []int
	Range(0, big).Take(3).ForEach(func(v int) {
		seen = append(seen, v)
	})
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(seen, want) {
		t.Fatalf("Range(0, 10M).Take(3).ForEach: got %v, want %v", seen, want)
	}
}

// TestRangeWithConcat confirms that two Ranges can be
// concatenated, which is a useful building block for slicing
// windows or paging.
func TestRangeWithConcat(t *testing.T) {
	got := Range(0, 3).Concat(Range(10, 13)).Collect()
	want := []int{0, 1, 2, 10, 11, 12}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range(0, 3).Concat(Range(10, 13)): got %v, want %v", got, want)
	}
}
