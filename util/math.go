package util

import "math"

func FloorDiv(a, b int32) int32 {
	return (a - ((a%b + b) % b)) / b
}

func BpeByNum(x float64) uint8 {
	return uint8(math.Ceil(math.Log2(x)))
}
