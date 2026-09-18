package reactivex

import (
	"context"
	"errors"
	"sync"

	"github.com/qianwj/typed/adt/option"
	r "github.com/qianwj/typed/adt/result"
)

// ErrNoElements is returned by FirstOrError when the source completes empty.
var ErrNoElements = errors.New("reactivex: flowable completed without elements")

// ToFlowable exposes this Single as a stream containing one value or an error.
// Conversion is lazy. All subscriptions share the Single's cached computation,
// while each subscription has independent demand and cancellation. A value
// waits for Request; an error does not require demand. Canceling the stream
// subscription does not cancel the Single or affect other subscribers.
func (s Single[T]) ToFlowable() Flowable[T] {
	return completionToFlowable(s.state, option.Of[T])
}

// ToFlowable exposes this Maybe as a stream containing zero or one value, or
// an error. Conversion is lazy and subscriptions share the cached computation.
// Present values, including nil, wait for Request; empty completion and errors
// do not require demand. Cancellation affects only the stream subscription.
func (m Maybe[T]) ToFlowable() Flowable[T] {
	return completionToFlowable(m.state, func(value option.Option[T]) option.Option[T] { return value })
}

// completionToFlowable buffers at most one value for each subscriber. The
// completion callback never waits for demand or runs downstream user code.
func completionToFlowable[V, T any](state *completion[V], valueOf func(V) option.Option[T]) Flowable[T] {
	return Flowable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		var mu sync.Mutex
		var upstream Subscription
		stopped := false
		sub := newBufferedSubscription(out, backpressureConfig{buffer: 1}, true, func() {
			mu.Lock()
			stopped = true
			current := upstream
			upstream = nil
			mu.Unlock()
			if current != nil {
				current.Cancel()
			}
		})
		out.OnSubscribe(sub)
		if ctx.Err() != nil {
			sub.Cancel()
		}
		select {
		case <-sub.Done():
			return sub
		default:
		}
		if ctx.Done() != nil {
			watchContext(ctx, sub)
		}
		sub.activate()
		current := state.subscribe(func(result r.Result[V]) {
			if result.IsFailure() {
				sub.terminate(result.Error())
				return
			}
			value := valueOf(result.Value())
			if value.IsPresent() {
				sub.offer(value.Get())
			}
			sub.terminate(nil)
		})
		mu.Lock()
		if stopped {
			mu.Unlock()
			current.Cancel()
		} else {
			upstream = current
			mu.Unlock()
		}
		return sub
	}}
}

// FirstElement returns a lazy Maybe for the first value emitted by this stream.
// It requests one value and cancels upstream as soon as that value arrives.
// Empty completion produces Success(Empty[T]()); errors before the first value
// produce Failure. Present zero and nil values are preserved.
//
// The returned Maybe subscribes once and caches its result. ctx controls that
// shared upstream subscription: cancellation produces Failure(ctx.Err()).
// Canceling an individual Maybe wait or subscription does not cancel upstream.
// A nil ctx means context.Background.
func (o Flowable[T]) FirstElement(ctx context.Context) Maybe[T] {
	return newMaybe(func(complete func(r.Result[option.Option[T]])) {
		subscribeFirst(ctx, o, complete)
	})
}

// FirstOrError returns a lazy Single for the first value emitted by this stream.
// An empty source fails with ErrNoElements. Further values are not inspected:
// upstream is canceled on the first value, without waiting for completion.
// Subscription, caching and context handling follow FirstElement.
func (o Flowable[T]) FirstOrError(ctx context.Context) Single[T] {
	return newSingle(func(complete func(r.Result[T])) {
		subscribeFirst(ctx, o, func(result r.Result[option.Option[T]]) {
			complete(result.FlatMap(func(value option.Option[T]) r.Result[T] {
				if value.IsEmpty() {
					return r.Failure[T](ErrNoElements)
				}
				return r.Success(value.Get())
			}))
		})
	})
}

func subscribeFirst[T any](ctx context.Context, source Flowable[T], complete func(r.Result[option.Option[T]])) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		complete(r.Failure[option.Option[T]](err))
		return
	}
	sub := &firstSubscriber[T]{complete: complete, done: make(chan struct{})}
	if ctx.Done() != nil {
		go func() {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
			case <-sub.done:
			}
		}()
	}
	source.Subscribe(ctx, sub)
}

// firstSubscriber arbitrates the first value, empty completion, error and
// context cancellation. It releases its lock before canceling upstream or
// notifying continuations, allowing synchronous and reentrant sources.
type firstSubscriber[T any] struct {
	mu       sync.Mutex
	upstream Subscription
	complete func(r.Result[option.Option[T]])
	done     chan struct{}
	settled  bool
}

func (s *firstSubscriber[T]) OnSubscribe(upstream Subscription) {
	s.mu.Lock()
	if s.settled {
		s.mu.Unlock()
		upstream.Cancel()
		return
	}
	s.upstream = upstream
	s.mu.Unlock()
	upstream.Request(1)
}

func (s *firstSubscriber[T]) OnNext(value T) {
	s.finish(r.Success(option.Of(value)))
}

func (s *firstSubscriber[T]) OnComplete() {
	s.finish(r.Success(option.Empty[T]()))
}

func (s *firstSubscriber[T]) OnError(err error) {
	s.finish(r.Failure[option.Option[T]](err))
}

func (s *firstSubscriber[T]) finish(result r.Result[option.Option[T]]) {
	s.mu.Lock()
	if s.settled {
		s.mu.Unlock()
		return
	}
	s.settled = true
	upstream, complete := s.upstream, s.complete
	s.upstream, s.complete = nil, nil
	close(s.done)
	s.mu.Unlock()
	if upstream != nil {
		upstream.Cancel()
	}
	complete(result)
}
