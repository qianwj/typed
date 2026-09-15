package concurrency

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkUnboundedQueue_1P1C measures single-producer / single-consumer
// throughput using an asynchronous consumer so the producer rarely blocks.
// The "items per second" number is dominated by the queue's per-operation
// cost.
func BenchmarkUnboundedQueue_1P1C(b *testing.B) {
	q := NewUnboundedBlockingQueue[int]()

	var done atomic.Bool
	var consumedCount atomic.Int64
	consumerDone := make(chan struct{})

	go func() {
		for !done.Load() || q.Size() > 0 {
			if opt := q.TryTake(); !opt.IsEmpty() {
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
	for {
		if q.Size() == 0 && consumedCount.Load() >= int64(b.N) {
			break
		}
		runtime.Gosched()
	}
	<-consumerDone
}

// BenchmarkUnboundedQueue_MPMC measures throughput with many producers and
// many consumers. Producers never block (queue is unbounded), but the cond
// Signal path and the mutex are exercised heavily.
func BenchmarkUnboundedQueue_MPMC(b *testing.B) {
	const workers = 8
	q := NewUnboundedBlockingQueue[int]()

	var wg sync.WaitGroup
	wg.Add(workers)
	perWorker := b.N / workers
	for w := range workers {
		w := w
		go func() {
			defer wg.Done()
			if w%2 == 0 {
				for i := range perWorker {
					q.Push(i)
				}
			} else {
				for range perWorker {
					q.Take()
				}
			}
		}()
	}
	wg.Wait()
}

// Note: we intentionally do not ship a "push into a fresh queue with no
// consumer" benchmark. It would measure how fast the queue grows without
// anyone draining it — which is mostly a property of `make([]T, n)` and
// GC pressure on the resulting large slice, not of the queue itself.
// Useful numbers come from 1P1C and MPMC, where the queue size stays
// bounded by the workload.