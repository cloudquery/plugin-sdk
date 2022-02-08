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
