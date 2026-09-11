package reactivex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSubjectMulticast(t *testing.T) {
	subject := NewSubject[int]()
	var a, b []int
	subject.ForEach(context.Background(), func(v int) { a = append(a, v) }, nil, nil)
	subject.ForEach(context.Background(), func(v int) { b = append(b, v) }, nil, nil)
	subject.OnNext(1)
	subject.OnNext(2)
	subject.OnComplete()
	if len(a) != 2 || len(b) != 2 || a[0] != 1 || a[1] != 2 || b[0] != 1 || b[1] != 2 {
		t.Fatalf("a=%v b=%v", a, b)
	}
}

func TestSubjectOverflowStrategies(t *testing.T) {
	check := func(strategy OverflowStrategy, want int) {
		s := NewSubject[int](WithBuffer(1), WithOverflow(strategy))
		values := make(chan int, 1)
		var sub Subscription
		sub = s.Subscribe(context.Background(), subscriberFuncs[int]{onSubscribe: func(v Subscription) { sub = v }, onNext: func(v int) { values <- v }})
		s.OnNext(1)
		s.OnNext(2)
		sub.Request(1)
		select {
		case got := <-values:
			if got != want {
				t.Fatalf("strategy %d: got %d, want %d", strategy, got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
		sub.Cancel()
	}
	check(OverflowDropLatest, 1)
	check(OverflowDropOldest, 2)
	s := NewSubject[int](WithBuffer(1), WithOverflow(OverflowError))
	errs := make(chan error, 1)
	s.Subscribe(context.Background(), subscriberFuncs[int]{onError: func(err error) { errs <- err }})
	s.OnNext(1)
	s.OnNext(2)
	select {
	case err := <-errs:
		if !errors.Is(err, ErrBackpressureOverflow) {
			t.Fatalf("got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
