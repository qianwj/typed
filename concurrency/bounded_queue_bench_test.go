package concurrency

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkBoundedQueue_1P1C measures single-producer / single-consumer
// throughput using an asynchronous consumer so the producer rarely blocks.
// The "items per second" number is dominated by the queue's per-operation
// cost.
func BenchmarkBoundedQueue_1P1C(b *testing.B) {
	const capacity = 1024
	q := NewBoundedBlockingQueue[int](capacity)

	var done atomic.Bool
	var consumedCount atomic.Int64
	consumerDone := make(chan struct{})

	go func() {
		for !done.Load() || q.Size() > 0 {
			if opt := q.TryPoll(); !opt.IsEmpty() {
				_ = opt.Get()
			}
			consumedCount.Add(1)
			runtime.Gosched()
		}
		close(consumerDone)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	b.StopTimer()

	done.Store(true)
	// Wait for the consumer to drain. The TryPoll loop above is bounded by
	// both `done` and `q.Size()`, so it will exit as soon as it observes
	// both, but we give it a moment in case it is mid-iteration.
	for q.Size() > 0 || consumedCount.Load() < int64(b.N) {
		runtime.Gosched()
	}
	// One more nudge: signal done again and let the consumer exit on its own.
	<-consumerDone
}

// BenchmarkBoundedQueue_MPMC measures throughput with many producers and
// many consumers — the worst case for a single-mutex design.
func BenchmarkBoundedQueue_MPMC(b *testing.B) {
	const (
		capacity = 64
		workers  = 8
	)
	q := NewBoundedBlockingQueue[int](capacity)

	var wg sync.WaitGroup
	wg.Add(workers)
	perWorker := b.N / workers
	for w := range workers {
		go func() {
			defer wg.Done()
			if w%2 == 0 {
				for i := range perWorker {
					q.Push(i)
				}
			} else {
				for range perWorker {
					q.Poll()
				}
			}
		}()
	}
	wg.Wait()
}
