package util

import "math"

func BpeByNum(x float64) uint8 {
	return uint8(math.Ceil(math.Log2(x)))
}
