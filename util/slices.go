package util

import (
	"cmp"
	"math/rand/v2"
	"slices"
)

func ShuffleSlice[T any](s []T) (ret []T) {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

func CountSlice[T comparable](s []T, o T) (ret int) {
	for _, v := range s {
		if v == o {
			ret++
		}
	}
	return
}

func GetUniqueSlice[T comparable](s []T) (unique map[T]struct{}) {
	unique = make(map[T]struct{})
	for _, v := range s {
		if _, ok := unique[v]; !ok {
			unique[v] = struct{}{}
		}
	}
	return
}

func CountUniqueSlice[T comparable](s []T) (ret int) {
	return len(GetUniqueSlice(s))
}

func MakeSingleValuedSlice[T any](v T, n int) (ret []T) {
	ret = make([]T, n)
	for i := range ret {
		ret[i] = v
	}
	return
}

func SetToSlice[T comparable](set map[T]struct{}) (slice []T) {
	slice = make([]T, 0, len(set))
	for k := range set {
		slice = append(slice, k)
	}
	return
}

func SetToSliceOrdered[T cmp.Ordered](set map[T]struct{}) (slice []T) {
	slice = SetToSlice(set)
	slices.Sort(slice)
	return
}

func MinSlice[T cmp.Ordered](s []T) (ret T) {
	ret = s[0]
	for _, v := range s[1:] {
		if v < ret {
			ret = v
		}
	}
	return
}

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func SliceRange[T Numeric](s, e T) (ret []T) {
	if s > e {
		return
	}
	for v := s; v < e; v++ {
		ret = append(ret, v)
	}
	return
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
