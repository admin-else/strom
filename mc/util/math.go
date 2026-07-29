package util

import "math"

// FloorDiv returns floor division (rounds toward negative infinity).
func FloorDiv(a, b int32) int32 {
	return (a - ((a%b + b) % b)) / b
}

// BpeByNum returns the minimum bits per element needed to represent x distinct values.
func BpeByNum(x float64) uint8 {
	return uint8(math.Ceil(math.Log2(x)))
}
