package executor

func toAny[T any](arr []T) []any {
	res := make([]any, len(arr))
	for i, t := range arr {
		res[i] = t
	}
	return res
}

func mapTo[T, U any](arr []T, convert func(t T) U) []U {
	res := make([]U, len(arr))
	for i, t := range arr {
		res[i] = convert(t)
	}
	return res
}
