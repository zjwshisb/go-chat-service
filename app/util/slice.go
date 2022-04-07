package util

func SliceMap[T any, K any](s []T, f func(s T) K) []K {
	result := make([]K, len(s), len(s))
	for index, item := range s {
		result[index] = f(item)
	}
	return result
}
