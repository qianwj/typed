package result

type Result[T any, E error] struct {
	result T
	err    E
}

func Ok[T any](result T) Result[T, error] {
	return Result[T, error]{result, nil}
}

func Err[T any, E error](error E) Result[T, E] {
	return Result[T, E]{nil, error}
}

func (r Result[T, E]) Succeed() bool {
	return r.err == nil
}

func (r Result[T, E]) Get() T {
	return r.result
}

func (r Result[T, E]) GetOr(t T) T {
	if r.Succeed() {
		return r.Get()
	}
	return t
}

func (r Result[T, E]) Unwrap() (T, E) {
	return r.result, r.err
}

func (r Result[T, E]) Failed() bool {
	return r.err != nil
}

func (r Result[T, E]) Error() E {
	return r.err
}
