package concurrency_test

import (
	"fmt"

	"github.com/qianwj/typed/concurrency"
	"github.com/qianwj/typed/utils/option"
)

// ExampleBoundedBlockingQueue shows the optional-returning TryTake
// paired with Push: the non-blocking probe avoids the (T, bool) shape
// by returning option.Optional.
func ExampleBoundedBlockingQueue() {
	q := concurrency.NewBoundedBlockingQueue[int](4)
	q.Push(1)
	q.Push(2)

	var first option.Optional[int] = q.TryTake()
	var second option.Optional[int] = q.TryTake()
	var empty option.Optional[int] = q.TryTake()

	fmt.Println(first.OrElse(-1), second.OrElse(-1), empty.OrElse(-1))
	// Output: 1 2 -1
}
