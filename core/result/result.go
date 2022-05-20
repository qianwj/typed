package result

type Result[T any] struct {
	result T
	err    error
}

func Ok[T any](result T) Result[T] {
	return Result[T]{result, nil}
}

func Err[T any](error error) Result[T] {
	var t T
	return Result[T]{t, error}
}

func (r Result[T]) Succeed() bool {
	return r.err == nil
}

func (r Result[T]) Get() T {
	return r.result
}

func (r Result[T]) GetOr(t T) T {
	if r.Succeed() {
		return r.Get()
	}
	return t
}

func (r Result[T]) Unwrap() (T, error) {
	return r.result, r.err
}

func (r Result[T]) Failed() bool {
	return r.err != nil
}

func (r Result[T]) Error() error {
	return r.err
}
