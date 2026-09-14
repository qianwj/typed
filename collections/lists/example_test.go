package lists_test

import (
	"fmt"

	"github.com/qianwj/typed/collections/lists"
)

// ExampleArrayListOf shows the left-to-right fluent style on
// ArrayList: build, filter, transform, collect — without a for loop.
func ExampleArrayListOf() {
	adults := lists.ArrayListOf(15, 22, 17, 30, 12).
		Filter(func(age int) bool { return age >= 18 }).
		Map(func(age int) int { return age * 2 }).
		Collect()
	fmt.Println(adults)
	// Output: [44 60]
}

// ExampleArrayList_First shows the Optional-based accessor on
// ArrayList, replacing the (T, bool) shape with a fluent chain.
func ExampleArrayList_First() {
	values := lists.ArrayListOf[int]()
	got := values.First().OrElse(-1)
	fmt.Println("empty:", got)

	values = lists.ArrayListOf(7, 8, 9)
	got = values.First().OrElse(-1)
	fmt.Println("first:", got)
	// Output:
	// empty: -1
	// first: 7
}
