package reactivex

import "errors"

// ErrBackpressureOverflow is sent to the affected subscriber when an
// OverflowError boundary receives a value after its pending queue is full.
// It is comparable with errors.Is and is never sent to sibling subscriptions.
var ErrBackpressureOverflow = errors.New("reactivex: backpressure buffer overflow")

// OverflowStrategy controls what happens when an asynchronous boundary has no
// free pending slot. It does not change Subscription demand; Request still
// controls whether a queued value may be delivered.
type OverflowStrategy uint8

const (
	// OverflowBlock waits for demand or buffer space and preserves every value.
	OverflowBlock OverflowStrategy = iota
	// OverflowDropLatest accepts the producer call but discards the incoming value
	// when the buffer is full. It preserves older pending values.
	OverflowDropLatest
	// OverflowDropOldest discards the oldest pending value and stores the incoming
	// value when the buffer is full. It requires WithBuffer with a positive size.
	OverflowDropOldest
	// OverflowKeepLatest keeps exactly one pending value and replaces it whenever a
	// newer value arrives. WithBuffer is ignored for this strategy.
	OverflowKeepLatest
	// OverflowError discards pending values, terminates that subscription and sends
	// ErrBackpressureOverflow. The Subject and sibling subscriptions continue.
	OverflowError
)

type backpressureConfig struct {
	// buffer counts pending values and excludes a callback currently in flight.
	buffer int
	// overflow determines behavior after buffer reaches capacity.
	overflow OverflowStrategy
}

// BackpressureOption configures one asynchronous boundary. Options are applied
// in order; later options replace earlier values. Passing nil is ignored.
type BackpressureOption func(*backpressureConfig)

// WithBuffer sets the number of pending values, excluding a callback currently
// in flight. WithBuffer(0) means that delivery must have available demand (or
// the selected drop/error policy handles the value immediately). Negative sizes
// panic. The default is zero.
func WithBuffer(size int) BackpressureOption {
	if size < 0 {
		panic("reactivex: negative buffer size")
	}
	return func(c *backpressureConfig) { c.buffer = size }
}

// WithOverflow selects an overflow strategy. Invalid strategy values panic.
// KeepLatest always uses one pending slot, regardless of WithBuffer. DropOldest
// with a zero buffer is rejected when the complete option set is applied.
func WithOverflow(strategy OverflowStrategy) BackpressureOption {
	if strategy > OverflowError {
		panic("reactivex: invalid overflow strategy")
	}
	return func(c *backpressureConfig) { c.overflow = strategy }
}

// applyBackpressureOptions resolves options and validates their combinations.
func applyBackpressureOptions(options []BackpressureOption) backpressureConfig {
	c := backpressureConfig{overflow: OverflowBlock}
	for _, option := range options {
		if option != nil {
			option(&c)
		}
	}
	if c.overflow == OverflowKeepLatest {
		c.buffer = 1
	}
	if c.overflow == OverflowDropOldest && c.buffer == 0 {
		panic("reactivex: DropOldest requires a positive buffer size")
	}
	return c
}
