package util

func MapElementsBy[K comparable, V any](arr []V, fn func(V) K) map[K][]V {
	result := map[K][]V{}

	for _, elem := range arr {
		result[fn(elem)] = append(result[fn(elem)], elem)
	}

	return result
}

func MapUniqueElementsBy[K comparable, V any](arr []V, fn func(V) K) map[K]V {
	result := map[K]V{}

	for _, elem := range arr {
		result[fn(elem)] = elem
	}

	return result
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))

	for key, _ := range m {
		result = append(result, key)
	}

	return result
}
