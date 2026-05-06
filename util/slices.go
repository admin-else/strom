package util

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func IsLastElement[T comparable](s []T, e T) bool {
	return len(s) > 0 && s[len(s)-1] == e
}

func GetNextLargestNumberSlice[T Numeric](a T, s []T) (ret T, found bool) {
	for _, v := range s {
		if v > a {
			return v, true
		}
	}
	return
}
