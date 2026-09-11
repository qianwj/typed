package reactivex

import (
	"context"
	"iter"
	"sync/atomic"
	"time"
)

// Just creates a cold observable that emits values in order and completes.
// Every subscription starts at the first value. With no arguments, the source
// completes without emitting a value or waiting for demand.
//
// Just delegates to FromSlice. Passing an existing slice with values... retains
// that slice's backing array; it does not make a defensive copy.
func Just[T any](values ...T) Observable[T] { return FromSlice(values) }

// FromSlice creates a cold observable that iterates values in slice order.
// Each subscription starts its own producer goroutine and demand gate. Nil
// and empty slices both complete without producing a value.
//
// The source retains values without copying. Do not mutate the backing array
// concurrently with a subscription; copy it before construction when a stable
// snapshot is required. Subscriber callbacks run on the producer goroutine,
// so a slow callback also slows this source.
func FromSlice[T any](values []T) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		s := newSubscription()
		watchContext(ctx, s)
		out.OnSubscribe(s)
		go func() {
			defer s.complete()
			for _, v := range values {
				if !s.acquire(ctx) {
					return
				}
				out.OnNext(v)
			}
			out.OnComplete()
		}()
		return s
	}}
}

// FromSeq invokes seq in a producer goroutine for each subscription. Values
// retain the order supplied by seq, and delivery waits for subscriber demand.
//
// Subscription setup is repeatable, but seq itself may capture single-use or
// mutable state. The caller is responsible for making repeated or concurrent
// invocations valid. seq is not cloned, and a nil seq is not an empty source.
//
// Returning false from the emission callback stops iteration. This cannot
// interrupt a sequence blocked inside its own code before yielding a value;
// such sequences need their own cooperative cancellation.
func FromSeq[T any](seq iter.Seq[T]) Observable[T] {
	return Create(func(ctx context.Context, emit func(T) bool, complete func()) {
		for v := range seq {
			if !emit(v) {
				return
			}
		}
		complete()
	})
}

// FromChannel adapts a channel using the default blocking backpressure policy.
// It is equivalent to FromChannelWithOptions(ch): no pending queue and no
// overflow drops. The source does not own or close ch.
//
// Subscriptions read from the same channel and compete for values. For a
// broadcast feed, send values through a Subject instead. Channel closure
// completes the stream once any value already read has been delivered.
func FromChannel[T any](ch <-chan T) Observable[T] {
	return FromChannelWithOptions(ch)
}

// FromChannelWithOptions reads ch through a per-subscription demand gate and
// bounded queue. Options are resolved at construction; each subscription gets
// its own queue, but subscriptions still compete to receive from the same ch.
//
// WithBuffer bounds pending queue entries, excluding an in-flight callback.
// With OverflowBlock the reader can additionally hold one value already read
// from ch while waiting for space. Buffering or a non-blocking overflow policy
// dispatches notifications separately from the input reader.
//
// Closing ch schedules normal completion after queued values drain. Draining
// requires demand; errors discard queued values. Cancelling a subscription
// stops its reader without closing ch. Producers writing to ch must therefore
// arrange their own cancellation. A nil ch waits until the subscription ends.
//
// For example, to terminate a subscriber that falls behind:
//
//	source := reactivex.FromChannelWithOptions(input,
//	    reactivex.WithBuffer(64),
//	    reactivex.WithOverflow(reactivex.OverflowError),
//	)
func FromChannelWithOptions[T any](ch <-chan T, options ...BackpressureOption) Observable[T] {
	config := applyBackpressureOptions(options)
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		s := newBufferedSubscription(out, config, config.overflow != OverflowBlock || config.buffer > 0, nil)
		watchContext(ctx, s)
		out.OnSubscribe(s)
		s.activate()
		go func() {
			for {
				select {
				case <-s.Done():
					return
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						s.terminate(nil)
						return
					}
					if !s.offer(v) {
						return
					}
				}
			}
		}()
		return s
	}}
}

// Create builds a source from a producer function run once per subscription
// in its own goroutine. run receives the subscription context, an emit
// callback and a complete callback.
//
// emit waits for demand and delivers one value. The producer must stop when
// emit returns false, including after explicit completion. complete sends
// OnComplete at most once; returning from run also invokes completion unless
// it has already been sent. This API has no explicit error-emission callback.
//
// Invoke emit and complete serially, and do not retain them after run returns.
// The producer owns its resources and should release them with defer. Cancelling
// the Subscription does not cancel the supplied context itself or interrupt
// unrelated blocking work; context-aware I/O must use a caller-owned context.
func Create[T any](run func(context.Context, func(T) bool, func())) Observable[T] {
	return Observable[T]{subscribe: func(ctx context.Context, out Subscriber[T]) Subscription {
		s := newSubscription()
		watchContext(ctx, s)
		out.OnSubscribe(s)
		go func() {
			defer s.complete()
			done := atomic.Bool{}
			complete := func() {
				if done.CompareAndSwap(false, true) {
					out.OnComplete()
				}
			}
			emit := func(v T) bool {
				if done.Load() || !s.acquire(ctx) {
					return false
				}
				out.OnNext(v)
				return true
			}
			run(ctx, emit, complete)
			if !done.Load() {
				complete()
			}
		}()
		return s
	}}
}

// Interval creates a ticker for each subscription and emits counters starting
// at zero. The first value follows the first tick, not subscription setup.
// Demand gates delivery; slow consumers can miss ticker ticks rather than
// accumulating an unbounded backlog. Counters count delivered values, not ticks.
//
// period must be positive. Otherwise time.NewTicker panics in the producer
// goroutine when the source is subscribed to. A subscription owns its ticker
// and stops it when that producer exits.
//
// In the current implementation the constructor's ctx argument is not used;
// the context supplied to Subscribe or ForEach controls the ticker loop.
func Interval(ctx context.Context, period time.Duration) Observable[uint64] {
	return Create(func(ctx context.Context, emit func(uint64) bool, complete func()) {
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		var i uint64
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !emit(i) {
					return
				}
				i++
			}
		}
	})
}
