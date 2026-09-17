package reactivex

import (
	"context"
	"sync"

	"github.com/qianwj/typed/adt"
)

// completion shares one lazy computation and its terminal result. Public
// consumers start it asynchronously; internal continuations reuse the current
// execution goroutine instead of launching goroutines that wait on other nodes.
type completion[T any] struct {
	mu        sync.Mutex
	start     func()
	result    adt.Result[T]
	settled   bool
	done      chan struct{}
	callbacks []*completionCallback[T]
}

type completionCallback[T any] struct {
	fn func(adt.Result[T])
}

func newCompletion[T any](start func(func(adt.Result[T]))) *completion[T] {
	c := &completion[T]{done: make(chan struct{})}
	c.start = func() { start(c.complete) }
	return c
}

// takeStart claims the source once. The source and all callbacks execute
// outside mu, including when they reenter a completed node.
func (c *completion[T]) takeStart() func() {
	c.mu.Lock()
	start := c.start
	c.start = nil
	c.mu.Unlock()
	return start
}

func (c *completion[T]) startAsync() {
	if start := c.takeStart(); start != nil {
		go start()
	}
}

func (c *completion[T]) complete(result adt.Result[T]) {
	c.mu.Lock()
	if c.settled {
		c.mu.Unlock()
		return
	}
	c.result = result
	c.settled = true
	callbacks := c.callbacks
	c.callbacks = nil
	close(c.done)
	c.mu.Unlock()

	for _, callback := range callbacks {
		callback.fn(result)
	}
}

// register atomically chooses between observing future completion and reading
// the cached result. No completion can fall between the check and registration.
func (c *completion[T]) register(callback *completionCallback[T]) (adt.Result[T], bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.settled {
		return c.result, true
	}
	c.callbacks = append(c.callbacks, callback)
	return adt.Result[T]{}, false
}

// onComplete is for internal composition, where running a cached continuation
// synchronously is intentional. Public subscriptions use subscribe instead.
func (c *completion[T]) onComplete(fn func(adt.Result[T])) {
	if result, ready := c.register(&completionCallback[T]{fn: fn}); ready {
		fn(result)
		return
	}
	if start := c.takeStart(); start != nil {
		start()
	}
}

func (c *completion[T]) isDone() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *completion[T]) await(ctx context.Context) adt.Result[T] {
	if c.isDone() {
		return c.result
	}
	if err := ctx.Err(); err != nil {
		return adt.Failure[T](err)
	}
	c.startAsync()
	select {
	case <-c.done:
		return c.result
	case <-ctx.Done():
		return adt.Failure[T](ctx.Err())
	}
}

func (c *completion[T]) subscribe(fn func(adt.Result[T])) Subscription {
	sub := &completionSubscription[T]{fn: fn, done: make(chan struct{})}
	callback := &completionCallback[T]{fn: sub.deliver}
	sub.remove = func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		for i, current := range c.callbacks {
			if current == callback {
				copy(c.callbacks[i:], c.callbacks[i+1:])
				c.callbacks[len(c.callbacks)-1] = nil
				c.callbacks = c.callbacks[:len(c.callbacks)-1]
				return
			}
		}
	}
	if result, ready := c.register(callback); ready {
		// A cached user callback must not block the Subscribe caller.
		go sub.deliver(result)
	} else {
		c.startAsync()
	}
	return sub
}

// A subscription's Done closes after delivery, or immediately on cancellation.
// Cancel suppresses a callback that has not started; an in-flight callback and
// the shared computation continue independently.
type completionSubscription[T any] struct {
	mu     sync.Mutex
	fn     func(adt.Result[T])
	remove func()
	done   chan struct{}
	once   sync.Once
}

func (sub *completionSubscription[T]) deliver(result adt.Result[T]) {
	sub.mu.Lock()
	fn := sub.fn
	sub.fn = nil
	sub.remove = nil
	sub.mu.Unlock()
	if fn != nil {
		defer sub.finish()
		fn(result)
	}
}

func (sub *completionSubscription[T]) finish() {
	sub.once.Do(func() { close(sub.done) })
}

func (sub *completionSubscription[T]) Request(uint64) {}

func (sub *completionSubscription[T]) Cancel() {
	sub.mu.Lock()
	sub.fn = nil
	remove := sub.remove
	sub.remove = nil
	sub.mu.Unlock()
	if remove != nil {
		remove()
	}
	sub.finish()
}

func (sub *completionSubscription[T]) Done() <-chan struct{} {
	return sub.done
}
