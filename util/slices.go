package util

import "math/rand/v2"

func ShuffleSlice[T any](s []T) (ret []T) {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}
