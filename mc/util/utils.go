package util

import (
	"cmp"
	"maps"
	"slices"
)

func OrderedKeys[Map map[K]V, K cmp.Ordered, V any](m Map) []K {
	keysIter := maps.Keys(m)
	keys := slices.Collect(keysIter)
	slices.Sort(keys)
	return keys
}

func AssertAndConvertMapValues[T any, M map[K]any, K cmp.Ordered](m M) map[K]T {
	ret := map[K]T{}
	for k, v := range m {
		ret[k] = v.(T)
	}
	return ret
}

func ReverseMap[M map[K]V, K cmp.Ordered, V cmp.Ordered](m M) map[V]K {
	ret := map[V]K{}
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

func PopFront[T any](s []T) (T, []T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, s, false
	}
	return s[0], s[1:], true
}

func CombineNamAndData(name string, data any) []any {
	return []any{name, data}
}
