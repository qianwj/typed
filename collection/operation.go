package collection

import "github.com/qianwj/typed/core/result"

func Map[T, U any](input []T, convert func(t T) U) []U {
	res := make([]U, len(input))
	for i, t := range input {
		res[i] = convert(t)
	}
	return res
}

func MapResult[T, U any](input []T, convert func(t T) (U, error)) result.Result[[]U] {
	res := make([]U, len(input))
	var err error
	for i, t := range input {
		res[i], err = convert(t)
		if err != nil {
			return result.Err[[]U](err)
		}
	}
	return result.Ok[[]U](res)
}
