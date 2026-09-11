// Package control holds generic control-flow helpers that do not
// fit the pattern-matching subpackage.
//
// The two primitives here cover the "do this N times" shape:
//
//   - Repeat is the simple loop helper. It runs f exactly
//     times times and returns nothing.
//   - RepeatE is the error-aware variant. It runs f up to
//     times times, stopping at the first non-nil error, and
//     returns (count, err) where count is the number of
//     successful iterations completed before the error.
//
// Both functions are zero-cost for times <= 0: Go's for-range
// over a non-positive integer skips the body without an explicit
// guard. The return value of RepeatE on times <= 0 is (0, nil).
//
// # Relationship to collections.Range
//
// collections.Range(0, n).ForEach(f) builds a Stream and then
// iterates it; control.Repeat(n, f) is the imperative equivalent
// without the Stream allocation. Pick Repeat when the iteration
// is a one-off side effect and Stream when the iteration is the
// start of a larger fluent pipeline (Map, Filter, Take, ...).
//
// The matching primitives (Pattern[T], Match, Between, etc.) live
// in the sibling package github.com/qianwj/typed/control/match.
package control

// Repeat calls f exactly times times and discards the result.
// times <= 0 is a no-op: the body is never entered and f is not
// called at all. Repeat returns nothing; callers that need
// progress reporting or error propagation use RepeatE.
//
// Repeat is the imperative equivalent of a for-i loop. The
// function is the natural form for "do this N times" where the
// body has no per-iteration value, no early break, and no error
// to surface. Examples include warming a cache, ticking a
// progress bar, or running a side effect a fixed number of times
// in a test fixture.
func Repeat(times int, f func()) {
	for range times {
		f()
	}
}

// RepeatE calls f up to times times, stopping at the first
// non-nil error. The return value is the count of successful
// iterations followed by the error that stopped the loop:
//
//   - On full success (f was called exactly times times and
//     returned nil every time) the return is (times, nil).
//     The count mirrors the requested times rather than the
//     actual iteration index, so the caller can use it as
//     "we did everything you asked for".
//   - On a non-nil error from the i-th call (0-indexed), the
//     return is (i, err): exactly i successful iterations
//     happened, then the (i+1)-th call returned err and the
//     loop stopped. f is not called again after the error.
//   - On times == 0 the body is never entered and the return
//     is (0, nil).
//   - On times < 0 the body is also never entered; the return
//     is (times, nil) — the same value as the input, which
//     is not a meaningful count. Negative times are a caller
//     error. The function does not panic so it stays safe to
//     call with computed counts, but callers should treat a
//     negative return as a signal to fix the input, not as a
//     real count.
//
// The (count, error) shape matches Go's standard (T, error)
// convention. On the error path the count is the number of
// *successful* iterations, so an error from the very first
// call returns (0, err) and an error from the last call
// returns (times - 1, err). This lets callers report partial
// progress ("12 of 100 attempts failed") without re-counting
// on the caller side.
//
// RepeatE is the right shape for retry loops, batch operations
// that should stop at the first failure, and any "do this N
// times unless something goes wrong" workflow. The nil error
// from f is treated as success; only a non-nil error short-
// circuits the loop.
func RepeatE(times int, f func() error) (int, error) {
	for i := range times {
		if err := f(); err != nil {
			return i, err
		}
	}
	return times, nil
}
