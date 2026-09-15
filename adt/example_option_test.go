package adt_test

import (
	"fmt"
	"strings"

	"github.com/qianwj/typed/adt"
)

// ExampleOption_OrElse shows the simplest way to turn an absent
// Option into a fallback value without unpacking the (T, bool) shape.
func ExampleOption_OrElse() {
	present := adt.Of(42)
	absent := adt.Empty[int]()
	fmt.Println(present.OrElse(-1))
	fmt.Println(absent.OrElse(-1))
	// Output:
	// 42
	// -1
}

// ExampleOption_Map demonstrates chaining an Option through a
// transformation, with the result collapsing to absent on an absent
// input.
func ExampleOption_Map() {
	mapped := adt.Of("  hello  ").
		Map(strings.TrimSpace).
		Map(strings.ToUpper)
	fmt.Println(mapped.Get())
	// Output: HELLO
}

// ExampleOption_FlatMap composes two Option-returning operations
// into one chain without manual IsPresent / Get dance.
func ExampleOption_FlatMap() {
	parse := func(s string) adt.Option[int] {
		if s == "" {
			return adt.Empty[int]()
		}
		return adt.Of(len(s))
	}
	chain := adt.Of("typed").FlatMap(parse)
	fmt.Println(chain.Get(), chain.OrElse(0))
	// Output: 5 5
}
