package util

import (
	"fmt"
	"math"
)

const (
	KB = 1000
	MB = KB * 1000
	GB = MB * 1000
)

func ToHumanReadableBytes(v int) string {
	switch {
	case v < KB:
		return fmt.Sprintf("%d Bytes", v)
	case v > KB && v < MB:
		res := float64(v) / KB
		if res == math.Trunc(res) {
			return fmt.Sprintf("%d KB", int(res))
		}

		return fmt.Sprintf("%.3f KB", res)
	case v > MB && v < GB:
		res := float64(v) / MB
		if res == math.Trunc(res) {
			return fmt.Sprintf("%d MB", int(res))
		}

		return fmt.Sprintf("%.3f MB", res)
	default:
		res := float64(v) / GB
		if res == math.Trunc(res) {
			return fmt.Sprintf("%d GB", int(res))
		}

		return fmt.Sprintf("%.3f GB", res)
	}
}
