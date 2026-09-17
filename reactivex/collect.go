package reactivex

import (
	"context"
	"sync"
)

// ToSlice starts a subscription and blocks until completion, an upstream error
// or cancellation of ctx. It requests maximum demand and appends received
// values in delivery order. Empty input returns a non-nil empty slice.
//
// The returned slice is always a snapshot — the library never mutates a
// slice it has handed back. On normal completion the returned error is nil
// and the slice holds every value delivered before completion. An upstream
// error returns the values collected before that error together with the
// error. Context cancellation calls Cancel and returns a snapshot of the
// values received so far together with ctx.Err(); if termination and context
// cancellation race, either outcome may be selected.
//
// Collection retains every value in memory and cannot complete normally for
// an infinite source. Bound such sources before collecting. A nil ctx means
// context.Background.
func (o Flowable[T]) ToSlice(ctx context.Context) ([]T, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		mu     sync.Mutex
		values = make([]T, 0)
	)
	result := make(chan error, 1)
	var once sync.Once
	sub := o.ForEach(ctx, func(value T) {
		mu.Lock()
		values = append(values, value)
		mu.Unlock()
	}, func(err error) {
		once.Do(func() { result <- err })
	}, func() {
		once.Do(func() { result <- nil })
	})
	select {
	case err := <-result:
		// Success path: the source has signalled terminal, so no OnNext
		// can be in flight. The live slice is stable and safe to return
		// without copying — the source's terminal callback is a sequence
		// point that comes after the last OnNext.
		return values, err
	case <-ctx.Done():
		// Cancel path: the producer's in-flight OnNext may still be
		// appending. Cancel signals the source to stop, then take a
		// snapshot under the same mutex so the copy is atomic with
		// respect to any concurrent append. An in-flight append either
		// completes before the snapshot (its value is in the snapshot)
		// or after (its value is in the internal buffer but not in the
		// returned slice); either is acceptable.
		sub.Cancel()
		mu.Lock()
		snapshot := append([]T(nil), values...)
		mu.Unlock()
		return snapshot, ctx.Err()
	}
}

// Collect collects o with ToSlice, then folds the resulting values into initial
// by calling f in delivery order. R can differ from T, allowing an aggregate,
// map or caller-owned collection to be built without depending on collections.
//
// On empty input, Collect returns initial with a nil error and does not call f.
// If collection fails, it returns initial and the error without calling f for
// any partial values. It does not accumulate incrementally while data arrives:
// the current implementation retains the entire input slice first.
//
// The accumulator is passed by value; maps, slices and pointers inside it are
// not deep-copied. f runs on the calling goroutine after collection completes.
//
//	total, err := reactivex.Collect(ctx, reactivex.Just(1, 2, 3), 0,
//	    func(sum, value int) int { return sum + value })
//	// On success, total is 6.
func Collect[T, R any](ctx context.Context, o Flowable[T], initial R, f func(R, T) R) (R, error) {
	result := initial
	values, err := o.ToSlice(ctx)
	if err != nil {
		return result, err
	}
	for _, value := range values {
		result = f(result, value)
	}
	return result, nil
}
