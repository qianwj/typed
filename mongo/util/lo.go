package util

func ToAny[T any](arr []T) []any {
	res := make([]any, len(arr))
	for i, t := range arr {
		res[i] = t
	}
	return res
}

func Map[T, U any](arr []T, convert func(t T) U) []U {
	res := make([]U, len(arr))
	for i, t := range arr {
		res[i] = convert(t)
	}
	return res
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
