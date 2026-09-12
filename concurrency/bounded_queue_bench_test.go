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
	const cap = 1024
	q := NewBoundedBlockingQueue[int](cap)

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
	// Wait for the consumer to drain. The TryTake loop above is bounded by
	// both `done` and `q.Size()`, so it will exit as soon as it observes
	// both, but we give it a moment in case it is mid-iteration.
	for {
		if q.Size() == 0 && consumedCount.Load() >= int64(b.N) {
			break
		}
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
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			if w%2 == 0 {
				for i := 0; i < perWorker; i++ {
					q.Push(i)
				}
			} else {
				for i := 0; i < perWorker; i++ {
					q.Take()
				}
			}
		}()
	}
	wg.Wait()
}

// BenchmarkChannel_1P1C is the reference benchmark for the equivalent
// channel-based implementation. Use it as a sanity check that our queue
// is in the right ballpark.
func BenchmarkChannel_1P1C(b *testing.B) {
	const cap = 1024
	ch := make(chan int, cap)

	var done atomic.Bool
	var consumedCount atomic.Int64
	consumerDone := make(chan struct{})

	go func() {
		for !done.Load() || len(ch) > 0 {
			select {
			case <-ch:
				consumedCount.Add(1)
			default:
				runtime.Gosched()
			}
		}
		close(consumerDone)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	b.StopTimer()

	done.Store(true)
	for {
		if len(ch) == 0 && consumedCount.Load() >= int64(b.N) {
			break
		}
		runtime.Gosched()
	}
	<-consumerDone
}

// BenchmarkChannel_MPMC is the equivalent for the channel implementation.
func BenchmarkChannel_MPMC(b *testing.B) {
	const (
		cap     = 64
		workers = 8
	)
	ch := make(chan int, cap)

	var wg sync.WaitGroup
	wg.Add(workers)
	perWorker := b.N / workers
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			if w%2 == 0 {
				for i := 0; i < perWorker; i++ {
					ch <- i
				}
			} else {
				for i := 0; i < perWorker; i++ {
					<-ch
				}
			}
		}()
	}
	wg.Wait()
}

// BenchmarkDrainTo measures the cost of pulling a batch of 64 elements out of
// a 1024-slot queue in one shot. Useful to see the constant factor compared to
// a tight loop of single-element TryTake calls.
func BenchmarkDrainTo(b *testing.B) {
	const cap = 1024
	q := NewBoundedBlockingQueue[int](cap)
	// Pre-fill so DrainTo always has work to do.
	for i := 0; i < cap; i++ {
		q.TryPush(i)
	}
	dst := make([]int, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.DrainTo(dst)
		// Replenish the drained batch so the next iteration has work.
		for j := 0; j < 64; j++ {
			q.TryPush(i*64 + j)
		}
	}
}
