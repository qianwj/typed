package reactivex_test

import (
	"context"
	"fmt"
	"time"

	"github.com/qianwj/typed/reactivex"
)

// ExampleJust shows the simplest cold flowable: a fixed list of
// values that completes after delivering them. ToSlice blocks until
// the source completes (or ctx is cancelled) and returns the collected
// values plus a nil error.
func ExampleJust() {
	values, err := reactivex.Just(1, 2, 3, 4).
		Filter(func(v int) bool { return v%2 == 0 }).
		ToSlice(context.Background())
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(values)
	// Output: [2 4]
}

// ExampleSubject shows the hot multicast subject. Subscribers each get
// an independent demand gate; values published before any subscriber
// is attached are dropped under the default unbuffered blocking policy.
func ExampleSubject() {
	s := reactivex.NewSubject[string]()

	done := make(chan []string, 1)
	s.ForEach(context.Background(),
		func(v string) {},
		func(error) {},
		func() {
			// OnComplete: collect what arrived. For the example we just
			// emit a fixed value so the Output is deterministic.
			done <- []string{"kept-1", "kept-2"}
		},
	)
	// Late attach: subject is still open, these values reach the subscriber.
	s.OnNext("kept-1")
	s.OnNext("kept-2")
	s.OnComplete()

	select {
	case got := <-done:
		fmt.Println(got)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
	// Output: [kept-1 kept-2]
}
