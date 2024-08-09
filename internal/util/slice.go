package util

import "golang.org/x/exp/constraints"

func SumFunc[T any, V constraints.Integer](arr []T, fn func(T) V) int {
	var result int

	for _, elem := range arr {
		result += int(fn(elem))
	}

	return result
}

func SliceMap[T any, P any](arr []T, fn func(T) P) []P {
	result := make([]P, 0, len(arr))

	for _, elem := range arr {
		result = append(result, fn(elem))
	}

	return result
}
