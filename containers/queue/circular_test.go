package queue

import (
	"testing"
)

func TestCircularBufferEnqueueDequeue(t *testing.T) {
	// Create a circular buffer with size 3
	cb := NewCircularBuffer[*int](3)

	// Enqueue three values
	value1 := 1
	value2 := 2
	value3 := 3
	if !cb.Enqueue(&value1) {
		t.Error("Enqueue failed for value1")
	}
	if !cb.Enqueue(&value2) {
		t.Error("Enqueue failed for value2")
	}
	if !cb.Enqueue(&value3) {
		t.Error("Enqueue failed for value3")
	}

	// Dequeue three values and verify their correctness
	dequeueValue1, ok := cb.Dequeue()
	if !ok || dequeueValue1 != &value1 {
		t.Error("Dequeue failed for value1")
	}
	dequeueValue2, ok := cb.Dequeue()
	if !ok || dequeueValue2 != &value2 {
		t.Error("Dequeue failed for value2")
	}
	dequeueValue3, ok := cb.Dequeue()
	if !ok || dequeueValue3 != &value3 {
		t.Error("Dequeue failed for value3")
	}

	// Attempt to dequeue from an empty buffer
	_, ok = cb.Dequeue()
	if ok {
		t.Error("Dequeue should fail for an empty buffer")
	}
}

func TestCircularBufferEnqueueFullBuffer(t *testing.T) {
	// Create a circular buffer with size 2
	cb := NewCircularBuffer[*int](2)

	// Enqueue two values successfully
	value1 := 1
	value2 := 2
	if !cb.Enqueue(&value1) {
		t.Error("Enqueue failed for value1")
	}
	if !cb.Enqueue(&value2) {
		t.Error("Enqueue failed for value2")
	}

	// Attempt to enqueue a third value (buffer is full)
	value3 := 3
	if cb.Enqueue(&value3) {
		t.Error("Enqueue should fail for value3 when buffer is full")
	}
}

func TestCircularBufferDequeueEmptyBuffer(t *testing.T) {
	// Create an empty circular buffer with size 2
	cb := NewCircularBuffer[any](2)

	// Attempt to dequeue from an empty buffer
	_, ok := cb.Dequeue()
	if ok {
		t.Error("Dequeue should fail for an empty buffer")
	}
}
