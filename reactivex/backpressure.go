package reactivex

import "sync"

// bufferedSubscription serializes notifications and owns one bounded ring for
// one Subject or channel subscription. The ring stores pending values only;
// one callback may be executing outside the ring at the same time.
//
// Callbacks, Done closure and owner cleanup run outside the main mutex where
// possible. Request, Cancel, offer and terminal transitions are safe to call
// concurrently. It is an implementation detail rather than a public queue
// abstraction; its behavior is exposed through Subscription and options.
type bufferedSubscription[T any] struct {
	// mu protects demand, queue and terminal state. Callbacks run after mu is
	// released so subscribers can safely cancel or request reentrantly.
	mu                                             sync.Mutex
	changed                                        *sync.Cond
	out                                            Subscriber[T]
	config                                         backpressureConfig
	async                                          bool
	done                                           chan struct{}
	doneClosed                                     bool
	onStop                                         func()
	cleanup                                        sync.Once
	queue                                          []T
	head, count                                    int
	demand                                         uint64
	ready, draining, closed, cancelled, terminated bool
	err                                            error
}

func newBufferedSubscription[T any](out Subscriber[T], config backpressureConfig, async bool, onStop func()) *bufferedSubscription[T] {
	s := &bufferedSubscription[T]{out: out, config: config, async: async, done: make(chan struct{}), queue: make([]T, config.buffer), onStop: onStop}
	s.changed = sync.NewCond(&s.mu)
	return s
}

// Done returns the lifecycle signal for this buffered subscription.
func (s *bufferedSubscription[T]) Done() <-chan struct{} { return s.done }

func (s *bufferedSubscription[T]) stop() {
	s.cleanup.Do(func() {
		if s.onStop != nil {
			s.onStop()
		}
	})
}

func (s *bufferedSubscription[T]) closeDoneLocked() {
	if !s.doneClosed {
		s.doneClosed = true
		close(s.done)
	}
}

// Cancel drops pending values, wakes blocked producers and removes the owner
// subscription. It is idempotent and does not wait for an in-flight callback.
func (s *bufferedSubscription[T]) Cancel() {
	s.mu.Lock()
	s.cancelled = true
	clear(s.queue)
	s.count = 0
	s.closeDoneLocked()
	s.changed.Broadcast()
	s.mu.Unlock()
	s.stop()
}

func (s *bufferedSubscription[T]) dispatch(f func()) {
	// Buffered boundaries use a worker per drain attempt so a slow callback does
	// not prevent the input reader from applying the configured overflow policy.
	// Unbuffered blocking boundaries retain synchronous callback behavior.
	if s.async {
		go f()
	} else {
		f()
	}
}

func (s *bufferedSubscription[T]) startLocked() bool {
	// The caller holds mu. A single drain owner prevents concurrent callbacks for
	// one subscription while allowing different Subject subscribers to progress
	// independently.
	if !s.ready || s.draining || s.cancelled || s.terminated {
		return false
	}
	if !(s.count > 0 && s.demand > 0) && !(s.closed && (s.err != nil || s.count == 0)) {
		return false
	}
	s.draining = true
	return true
}

// activate ensures no notification precedes the return from OnSubscribe.
func (s *bufferedSubscription[T]) activate() {
	s.mu.Lock()
	s.ready = true
	start := s.startLocked()
	s.changed.Broadcast()
	s.mu.Unlock()
	if start {
		s.dispatch(s.drain)
	}
}

// Request increases demand and starts draining queued values when possible.
// Demand saturates instead of wrapping at uint64 maximum.
func (s *bufferedSubscription[T]) Request(n uint64) {
	s.mu.Lock()
	if s.cancelled || s.terminated {
		s.mu.Unlock()
		return
	}
	if ^uint64(0)-s.demand < n {
		s.demand = ^uint64(0)
	} else {
		s.demand += n
	}
	start := s.startLocked()
	s.changed.Broadcast()
	s.mu.Unlock()
	if start {
		s.dispatch(s.drain)
	}
}

func (s *bufferedSubscription[T]) consumeDemandLocked() {
	if s.demand != ^uint64(0) {
		s.demand--
	}
}

// offer accepts one value at an asynchronous boundary. It either reserves a
// demand slot for immediate delivery, appends to the ring, waits for capacity,
// applies a drop policy, or starts terminal error handling. The return value is
// false only when the producer should stop (cancellation or OverflowError).
func (s *bufferedSubscription[T]) offer(value T) bool {
	s.mu.Lock()
	for {
		if s.cancelled || s.closed || s.terminated {
			s.mu.Unlock()
			return false
		}
		if s.ready && !s.draining && s.demand > 0 && s.count == 0 {
			s.consumeDemandLocked()
			s.draining = true
			s.mu.Unlock()
			s.dispatch(func() {
				s.mu.Lock()
				cancelled := s.cancelled
				s.mu.Unlock()
				if !cancelled {
					s.out.OnNext(value)
				}
				s.drain()
			})
			return true
		}
		if s.count < len(s.queue) {
			s.queue[(s.head+s.count)%len(s.queue)] = value
			s.count++
			start := s.startLocked()
			s.mu.Unlock()
			if start {
				s.dispatch(s.drain)
			}
			return true
		}
		switch s.config.overflow {
		case OverflowBlock:
			s.changed.Wait()
		case OverflowDropLatest:
			s.mu.Unlock()
			return true
		case OverflowDropOldest, OverflowKeepLatest:
			s.queue[s.head] = value
			s.head = (s.head + 1) % len(s.queue)
			s.mu.Unlock()
			return true
		case OverflowError:
			s.closed, s.err = true, ErrBackpressureOverflow
			clear(s.queue)
			s.count = 0
			start := s.startLocked()
			s.changed.Broadcast()
			s.mu.Unlock()
			if start {
				s.dispatch(s.drain)
			}
			return false
		}
	}
}

// terminate marks the boundary complete. Normal completion preserves queued
// values and waits for demand; an error discards them and can be delivered
// immediately. Calling terminate more than once has no effect.
func (s *bufferedSubscription[T]) terminate(err error) {
	s.mu.Lock()
	if s.closed || s.cancelled || s.terminated {
		s.mu.Unlock()
		return
	}
	s.closed, s.err = true, err
	if err != nil {
		clear(s.queue)
		s.count = 0
	}
	start := s.startLocked()
	s.changed.Broadcast()
	s.mu.Unlock()
	if start {
		s.dispatch(s.drain)
	}
}

// drain is the sole callback owner for a subscription. It repeatedly consumes
// queued values while demand exists, then either waits for a later Request or
// sends the pending terminal signal. It must never hold mu while invoking user
// code.
func (s *bufferedSubscription[T]) drain() {
	for {
		s.mu.Lock()
		if s.cancelled {
			s.draining = false
			s.changed.Broadcast()
			s.mu.Unlock()
			return
		}
		if s.count > 0 && s.demand > 0 {
			value := s.queue[s.head]
			var zero T
			s.queue[s.head] = zero
			s.head = (s.head + 1) % len(s.queue)
			s.count--
			s.consumeDemandLocked()
			s.changed.Broadcast()
			s.mu.Unlock()
			s.out.OnNext(value)
			continue
		}
		s.draining = false
		s.changed.Broadcast()
		if s.closed && (s.err != nil || s.count == 0) {
			s.terminated = true
			err := s.err
			s.mu.Unlock()
			if err != nil {
				s.out.OnError(err)
			} else {
				s.out.OnComplete()
			}
			s.mu.Lock()
			s.closeDoneLocked()
			s.mu.Unlock()
			s.stop()
			return
		}
		s.mu.Unlock()
		return
	}
}
