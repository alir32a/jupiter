package util

func GetStructsField[T any, V any](arr []T, fn func(T) V) []V {
	result := make([]V, 0, len(arr))

	for _, el := range arr {
		result = append(result, fn(el))
	}

	return result
}

func MapSlice[T any](arr []T, fn func(T) T) {
	for i, el := range arr {
		arr[i] = fn(el)
	}
}

func MapStructsByField[T any, K comparable](arr []T, fn func(T) K) map[K][]T {
	result := map[K][]T{}

	for _, el := range arr {
		result[fn(el)] = append(result[fn(el)], el)
	}

	return result
}

func MapStructsByUniqueField[T any, K comparable](arr []T, fn func(T) K) map[K]T {
	result := map[K]T{}

	for _, el := range arr {
		result[fn(el)] = el
	}

	return result
}
