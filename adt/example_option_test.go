package adt_test

import (
	"fmt"
	"strings"

	"github.com/qianwj/typed/adt"
)

// ExampleOptional_OrElse shows the simplest way to turn an absent
// Optional into a fallback value without unpacking the (T, bool) shape.
func ExampleOptional_OrElse() {
	present := adt.Of(42)
	absent := adt.Empty[int]()
	fmt.Println(present.OrElse(-1))
	fmt.Println(absent.OrElse(-1))
	// Output:
	// 42
	// -1
}

// ExampleOptional_Map demonstrates chaining an Optional through a
// transformation, with the result collapsing to absent on an absent
// input.
func ExampleOptional_Map() {
	mapped := adt.Of("  hello  ").
		Map(strings.TrimSpace).
		Map(strings.ToUpper)
	fmt.Println(mapped.Get())
	// Output: HELLO
}

// ExampleOptional_FlatMap composes two Optional-returning operations
// into one chain without manual IsPresent / Get dance.
func ExampleOptional_FlatMap() {
	parse := func(s string) adt.Optional[int] {
		if s == "" {
			return adt.Empty[int]()
		}
		return adt.Of(len(s))
	}
	chain := adt.Of("typed").FlatMap(parse)
	fmt.Println(chain.Get(), chain.OrElse(0))
	// Output: 5 5
}
