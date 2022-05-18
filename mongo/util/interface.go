package util

func ToInterfaceSlice[T any](slice []T) []interface{} {
	res := make([]interface{}, len(slice))
	for i, t := range slice {
		res[i] = t
	}
	return res
}

func ToAnySlice[T any](slice []interface{}) []T {
	res := make([]T, len(slice))
	for i, t := range slice {
		res[i] = t.(T)
	}
	return res
}
