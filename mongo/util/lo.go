package util

func ToAny[T any](arr []T) []any {
	res := make([]any, len(arr))
	for i, t := range arr {
		res[i] = t
	}
	return res
}

func OrderedMap[T, U any](arr []T, convert func(T) U) []U {
	res := make([]U, len(arr))
	for i, t := range arr {
		res[i] = convert(t)
	}
	return res
}

func Filter[T any](arr []T, predict func(T) bool) []T {
	res := make([]T, 0)
	for _, t := range arr {
		if predict(t) {
			res = append(res, t)
		}
	}
	return res
}

func Reduce[T, R any](arr []T, acc func(prev R, cur T) R, initial R) R {
	for _, item := range arr {
		initial = acc(initial, item)
	}
	return initial
}

func ErrMap[T, U any](arr []T, convert func(t T) (U, error)) ([]U, error) {
	res := make([]U, len(arr))
	var err error
	for i, t := range arr {
		res[i], err = convert(t)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
