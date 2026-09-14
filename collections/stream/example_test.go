package stream_test

import (
	"fmt"

	"github.com/qianwj/typed/collections/lists"
)

// ExampleArrayList_Stream shows the explicit lazy layer: a Stream[T]
// over the same ArrayList values, with early-terminating terminals.
func ExampleArrayList_Stream() {
	got := lists.ArrayListOf(1, 2, 3, 4, 5).
		Stream().
		Filter(func(v int) bool { return v%2 == 1 }).
		Map(func(v int) int { return v * 10 }).
		Take(2).
		Collect()
	fmt.Println(got)
	// Output: [10 30]
}
