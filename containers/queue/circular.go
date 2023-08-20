package queue

import (
	"github.com/qianwj/typed/containers"
	"sync/atomic"
)

type CircularBuffer[T any] struct {
	containers.Queue[T]
	buffer []T
	mask   uint64
	read   uint64
	write  uint64
}

func NewCircularBuffer[T any](size uint64) *CircularBuffer[T] {
	if size == 0 || (size&(size-1)) != 0 {
		panic("Buffer size must be a power of 2")
	}
	return &CircularBuffer[T]{
		buffer: make([]T, size),
		mask:   size - 1,
		read:   0,
		write:  0,
	}
}

func (cb *CircularBuffer[T]) Enqueue(value T) bool {
	write := atomic.LoadUint64(&cb.write)
	read := cb.read
	if write-read >= uint64(len(cb.buffer)) {
		return false // Buffer is full
	}
	cb.buffer[write&cb.mask] = value
	atomic.StoreUint64(&cb.write, write+1)
	return true
}

func (cb *CircularBuffer[T]) Dequeue() (T, bool) {
	read := atomic.LoadUint64(&cb.read)
	write := cb.write
	var t T
	if read >= write {
		return t, false // Buffer is empty
	}
	value := cb.buffer[read&cb.mask]
	atomic.StoreUint64(&cb.read, read+1)
	return value, true
}
