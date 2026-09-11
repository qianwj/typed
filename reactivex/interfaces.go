package reactivex

import (
	"context"
	"sync"
)

// Subscriber receives subscription setup, values and terminal notifications
// from a Publisher. Errors are notifications, not values of T.
//
// Implementations should request their initial demand in OnSubscribe and keep
// callbacks short. A publisher may invoke callbacks on a producer goroutine or
// directly on a caller's goroutine; subscribers must not depend on a particular
// execution thread. State shared by several subscriptions needs its own
// synchronization. Callback panics are not converted into OnError notifications.
type Subscriber[T any] interface {
	// OnSubscribe supplies the handle for this subscription before values are
	// delivered. It may call Request or Cancel before returning. Retain the
	// handle when demand will be replenished as processing finishes.
	OnSubscribe(Subscription)
	// OnNext processes a delivered value. Its receiver controls any mutable
	// state it shares with other subscriptions or application goroutines.
	OnNext(T)
	// OnError reports abnormal termination. It does not consume demand and
	// must not be interpreted as an ordinary data item to request again.
	OnError(error)
	// OnComplete reports normal termination and does not consume demand.
	// A buffered publisher may defer this notification until pending values
	// have been delivered, which still requires sufficient demand.
	OnComplete()
}

// Publisher is the subscription boundary implemented by Observable and Subject.
// Accept a Publisher when a component only consumes notifications; use the
// concrete Observable type when building a fluent operator pipeline.
//
// The interface does not imply replay, independent inputs or multicast. Those
// properties come from the source: a slice can be iterated again, whereas two
// subscriptions to the same channel compete to receive its values.
type Publisher[T any] interface {
	// Subscribe connects a subscriber to the publisher and returns its handle.
	// ctx controls the subscription lifetime, not ownership of input resources.
	Subscribe(context.Context, Subscriber[T]) Subscription
}

// Subscription controls demand and the lifetime of one subscription.
// Request and Cancel may be called from different goroutines. A subscription
// belongs to one subscriber; cancelling it does not close the original channel
// or cancel sibling subscribers of a Subject.
type Subscription interface {
	// Request adds n to outstanding demand; it is not an absolute batch size.
	// Zero adds no demand, and addition saturates at the maximum uint64 value
	// instead of wrapping. Operators may compensate for filtered values.
	Request(n uint64)
	// Cancel signals cancellation and is idempotent. It does not wait for a
	// callback already executing, and does not forcibly interrupt user code.
	Cancel()
	// Done closes on cancellation or termination. It is a lifecycle signal,
	// not a join operation for all work started by a producer or callback.
	Done() <-chan struct{}
}

// subscription supplies the demand gate for sources without an overflow queue.
// changed wakes producers waiting for demand; done is also used to release the
// context watcher. It must not be copied after its condition variable is bound.
type subscription struct {
	done    chan struct{}
	mu      sync.Mutex
	demand  uint64
	changed *sync.Cond
	once    sync.Once
}

func newSubscription() *subscription {
	s := &subscription{done: make(chan struct{})}
	s.changed = sync.NewCond(&s.mu)
	return s
}

type cancellable interface {
	Cancel()
	Done() <-chan struct{}
}

// watchContext translates context cancellation into subscription cancellation.
// It also exits on Done so a completed stream does not retain a watcher until
// a long-lived parent context is cancelled.
func watchContext(ctx context.Context, s cancellable) {
	go func() {
		select {
		case <-ctx.Done():
			s.Cancel()
		case <-s.Done():
		}
	}()
}

func (s *subscription) Request(n uint64) {
	if n == 0 {
		return
	}
	s.mu.Lock()
	if ^uint64(0)-s.demand < n {
		s.demand = ^uint64(0)
	} else {
		s.demand += n
	}
	s.changed.Broadcast()
	s.mu.Unlock()
}

func (s *subscription) Cancel() {
	s.once.Do(func() { close(s.done); s.mu.Lock(); s.changed.Broadcast(); s.mu.Unlock() })
}

func (s *subscription) Done() <-chan struct{} { return s.done }

func (s *subscription) complete() { s.Cancel() }

// acquire waits for demand and consumes one unit before a source delivers a
// value. The source's context watcher wakes waiters when the context ends.
func (s *subscription) acquire(ctx context.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for s.demand == 0 {
		select {
		case <-s.done:
			return false
		case <-ctx.Done():
			s.Cancel()
			return false
		default:
		}
		s.changed.Wait()
	}
	s.demand--
	return true
}

// Observable describes a typed source and its operator chain. Constructing or
// copying an Observable does not start consumption; subscribing does.
//
// Each subscription creates new operator state. Source resources are not
// necessarily duplicated: FromChannel subscriptions share their input, and
// FromSlice subscriptions share the original backing array. The zero value
// has no source and cannot be subscribed to; use a constructor such as Just.
type Observable[T any] struct {
	// subscribe is private so sources control lifecycle setup.
	subscribe func(context.Context, Subscriber[T]) Subscription
}

var _ Publisher[any] = Observable[any]{}
