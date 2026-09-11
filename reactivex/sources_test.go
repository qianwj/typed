package reactivex

import (
	"context"
	"testing"
	"time"
)

func TestDemandAndCancel(t *testing.T) {
	var sub Subscription
	values := make(chan int, 2)
	sub = FromSlice([]int{1, 2, 3}).Subscribe(context.Background(), subscriberFuncs[int]{onSubscribe: func(s Subscription) { sub = s }, onNext: func(v int) { values <- v }})
	sub.Request(2)
	if a, b := <-values, <-values; a != 1 || b != 2 {
		t.Fatalf("got %d, %d", a, b)
	}
	sub.Cancel()
	select {
	case <-sub.Done():
	case <-time.After(time.Second):
		t.Fatal("cancel did not complete")
	}
}
