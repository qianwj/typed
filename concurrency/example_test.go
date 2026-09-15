package concurrency_test

import (
	"fmt"

	"github.com/qianwj/typed/concurrency"
)

// ExampleBoundedBlockingQueue shows the optional-returning TryPoll
// paired with Push: the non-blocking probe avoids the (T, bool) shape
// by returning option.Optional.
//
// The explicit type annotations on each TryPoll call are the
// assertions under test (compile-time check that TryPoll returns
// option.Optional[int] exactly).
func ExampleBoundedBlockingQueue() {
	q := concurrency.NewBoundedBlockingQueue[int](4)
	q.Push(1)
	q.Push(2)

	var first, second, empty = q.TryPoll(), q.TryPoll(), q.TryPoll()

	fmt.Println(first.OrElse(-1), second.OrElse(-1), empty.OrElse(-1))
	// Output: 1 2 -1
}
