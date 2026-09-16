package stream_test

import (
	"fmt"
	"strings"

	"github.com/qianwj/typed/collections/lists"
	"github.com/qianwj/typed/collections/stream"
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

func ExampleStream_Associate() {
	result := stream.Of("a", "bb", "ccc", "dd").
		Associate(func(s string) (int, string) { return len(s), s }).
		Filter(func(k int, _ string) bool { return k == 2 }).
		MapValues(func(_ int, v string) string { return strings.ToUpper(v) }).
		Collect()
	fmt.Println(result[2])
	// Output: DD
}
