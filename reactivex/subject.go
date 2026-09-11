package reactivex

import (
	"context"
	"sync"
)

// Subject is a hot multicast publisher and subscriber. Each subscriber has
// independent demand and a buffer configured at construction time.
//
// Values are not replayed: OnNext with no subscribers discards the value, and
// a new subscriber only participates in later publications. Terminal state is
// retained, so late subscribers receive completion or the stored error.
//
// The default unbuffered blocking mode invokes callbacks synchronously when
// demand is available. A buffer or a non-blocking overflow strategy dispatches
// notifications asynchronously, still serializing callbacks per subscription.
// Across different subscriptions callbacks may run concurrently. Construct a
// Subject with NewSubject and do not copy it after use.
type Subject[T any] struct {
	mu          sync.Mutex
	subscribers map[*bufferedSubscription[T]]struct{}
	closed      bool
	err         error
	upstream    Subscription
	config      backpressureConfig
}

// NewSubject creates a hot multicast subject with optional backpressure settings.
// Options apply to every subscription, but queue contents and demand are never
// shared. Without options, publication waits for each subscriber's demand.
//
// WithBuffer and WithOverflow configure pending values and overflow handling.
// Invalid combinations, such as DropOldest without a positive buffer, panic
// during construction. A dropped value or overflow error affects the individual
// subscription, not the Subject's terminal state.
func NewSubject[T any](options ...BackpressureOption) *Subject[T] {
	return &Subject[T]{subscribers: make(map[*bufferedSubscription[T]]struct{}), config: applyBackpressureOptions(options)}
}

var _ Publisher[any] = (*Subject[any])(nil)
var _ Subscriber[any] = (*Subject[any])(nil)

// Subscribe registers one subscriber and returns its independent subscription.
// out must be non-nil and can request demand in OnSubscribe. A nil ctx becomes
// context.Background. Notifications wait until OnSubscribe returns.
//
// Cancelling removes this subscriber without closing the Subject. A subscription
// made after Subject termination receives the terminal notification without
// replaying prior values. In buffered mode this notification may be asynchronous.
func (s *Subject[T]) Subscribe(ctx context.Context, out Subscriber[T]) Subscription {
	if ctx == nil {
		ctx = context.Background()
	}
	var sub *bufferedSubscription[T]
	sub = newBufferedSubscription(out, s.config, s.config.overflow != OverflowBlock || s.config.buffer > 0, func() { s.remove(sub) })
	watchContext(ctx, sub)
	s.mu.Lock()
	if s.closed {
		err := s.err
		s.mu.Unlock()
		out.OnSubscribe(sub)
		sub.activate()
		sub.terminate(err)
		return sub
	}
	s.subscribers[sub] = struct{}{}
	s.mu.Unlock()
	out.OnSubscribe(sub)
	sub.activate()
	return sub
}

// ForEach registers callback functions and requests unbounded demand.
// Any callback may be nil. It returns without waiting for termination, and
// callback scheduling follows the Subject's configured synchronous or buffered
// mode. An omitted onError silently ignores overflow and upstream errors.
func (s *Subject[T]) ForEach(ctx context.Context, onNext func(T), onError func(error), onComplete func()) Subscription {
	out := &callbackSubscriber[T]{onNext: onNext, onError: onError, onComplete: onComplete}
	sub := s.Subscribe(ctx, out)
	sub.Request(^uint64(0))
	return sub
}

func (s *Subject[T]) remove(sub *bufferedSubscription[T]) {
	s.mu.Lock()
	delete(s.subscribers, sub)
	s.mu.Unlock()
}

// OnSubscribe records an upstream subscription connected to this subject.
// It cancels the handle immediately if the Subject has already terminated.
// Otherwise the handle is retained and cancelled when the Subject terminates.
//
// This method currently does not request upstream demand, and downstream
// requests are not automatically aggregated into upstream requests. When using
// the Subject as a Subscriber, explicitly request from the returned upstream
// subscription. Only the most recently supplied handle is retained.
func (s *Subject[T]) OnSubscribe(upstream Subscription) {
	s.mu.Lock()
	s.upstream = upstream
	closed := s.closed
	s.mu.Unlock()
	if closed {
		upstream.Cancel()
	}
}

// OnNext broadcasts value to all active subscribers.
// It snapshots the subscriber set before delivery; a subscription registered
// afterward does not receive this value. No subscribers means the value is
// discarded, and calls after Subject termination are ignored.
//
// Each subscription applies its own demand and overflow policy. With Block,
// a slow subscriber can delay this call and delivery to other subscribers.
// With dropping policies, OnNext returning does not guarantee every subscriber
// received the value. Concurrent publishers have no global ordering guarantee
// across subscribers; serialize calls when a common publication order matters.
// Do not synchronously republish into a full blocking Subject from its callback.
func (s *Subject[T]) OnNext(value T) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	subs := make([]*bufferedSubscription[T], 0, len(s.subscribers))
	for sub := range s.subscribers {
		subs = append(subs, sub)
	}
	s.mu.Unlock()
	for _, sub := range subs {
		sub.offer(value)
	}
}

// OnError terminates the subject and notifies all active subscribers.
// Supply a non-nil error. It cancels the recorded upstream handle and discards
// buffered values before scheduling the error; an in-flight callback may still
// finish first. The error notification itself does not need demand.
//
// The first terminal call wins. Later OnNext, OnError and OnComplete calls are
// ignored, and late subscribers observe the stored error. OverflowError on an
// individual subscription does not call this method or terminate its siblings.
func (s *Subject[T]) OnError(err error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed, s.err = true, err
	subs := make([]*bufferedSubscription[T], 0, len(s.subscribers))
	for sub := range s.subscribers {
		subs = append(subs, sub)
	}
	s.subscribers = make(map[*bufferedSubscription[T]]struct{})
	upstream := s.upstream
	s.mu.Unlock()
	if upstream != nil {
		upstream.Cancel()
	}
	for _, sub := range subs {
		sub.terminate(err)
	}
}

// OnComplete closes the Subject to new values and schedules normal completion
// for its current subscribers. It cancels the recorded upstream handle.
//
// Each subscriber drains its own pending queue before OnComplete is delivered.
// That drain still needs demand, so this method may return while subscribers
// remain unfinished. Call Cancel on an individual subscription to discard its
// queue. Late subscribers receive completion with no replay; later terminal
// calls and publications have no effect.
func (s *Subject[T]) OnComplete() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	subs := make([]*bufferedSubscription[T], 0, len(s.subscribers))
	for sub := range s.subscribers {
		subs = append(subs, sub)
	}
	s.subscribers = make(map[*bufferedSubscription[T]]struct{})
	upstream := s.upstream
	s.mu.Unlock()
	if upstream != nil {
		upstream.Cancel()
	}
	for _, sub := range subs {
		sub.terminate(nil)
	}
}
