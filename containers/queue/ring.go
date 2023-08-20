package queue

import (
	"github.com/qianwj/typed/containers"
	"sync/atomic"
)

type RingBuffer[T any] struct {
	containers.Queue[T]
	buffer []T
	size   uint64
	mask   uint64
	read   uint64
	write  uint64
}

func NewRingBuffer[T any](size uint64) *RingBuffer[T] {
	if size == 0 || (size&(size-1)) != 0 {
		panic("Buffer size must be a power of 2")
	}
	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
		mask:   size - 1,
		read:   0,
		write:  0,
	}
}

func (rb *RingBuffer[T]) Enqueue(value T) bool {
	read := atomic.LoadUint64(&rb.read)
	write := atomic.LoadUint64(&rb.write)

	if write-read >= rb.size {
		return false // Buffer is full
	}

	index := write & rb.mask
	rb.buffer[index] = value

	for !atomic.CompareAndSwapUint64(&rb.write, write, write+1) {
		write = atomic.LoadUint64(&rb.write)
		if write-read >= rb.size {
			return false // Buffer is full
		}
		index = write & rb.mask
		rb.buffer[index] = value
	}

	return true
}

func (rb *RingBuffer[T]) Dequeue() (T, bool) {
	read := atomic.LoadUint64(&rb.read)
	write := atomic.LoadUint64(&rb.write)

	var t T
	if read >= write {
		return t, false // Buffer is empty
	}

	index := read & rb.mask
	value := rb.buffer[index]

	for !atomic.CompareAndSwapUint64(&rb.read, read, read+1) {
		read = atomic.LoadUint64(&rb.read)
		if read >= write {
			return t, false // Buffer is empty
		}
		index = read & rb.mask
		value = rb.buffer[index]
	}

	return value, true
}
