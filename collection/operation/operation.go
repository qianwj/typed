package operation

func Map[T, U any](input []T, convert func(t T) U) []U {
	res := make([]U, len(input))
	for i, t := range input {
		res[i] = convert(t)
	}
	return res
}
