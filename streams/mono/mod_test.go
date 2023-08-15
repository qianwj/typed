package mono

import (
	"fmt"
	"streams"
	"testing"
)

type testSub struct {
	streams.Publisher
}

func (t *testSub) OnNext(item any) {
	fmt.Printf("data: %+v\r\n", item)
}

func (t *testSub) OnError(err error) {
	fmt.Printf("error: %+v\r\n", err)
}

func (t *testSub) OnComplete() {
	println("complete")
}

func (t *testSub) OnSubscribe(streams.Subscription) {

}

func TestJust(t *testing.T) {
	observable := Just[int](1231)()
	observable.Subscribe(&testSub{})
}
