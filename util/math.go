package util

func FloorDiv(a, b int32) int32 {
	return (a - ((a%b + b) % b)) / b
}
