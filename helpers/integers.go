package helpers

import "math"

// Uint64ToInt64 if value is bigger than math.MaxInt64 return math.MaxInt64
// otherwise returns original value casted to int64
func Uint64ToInt64(i uint64) int64 {
	if i > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(i)
}

func Uint64ToInt(i uint64) int {
	if i > math.MaxInt {
		return math.MaxInt
	}
	return int(i)
}
