package main

import (
	"cmp"
	"maps"
	"slices"
	"unicode"
)

var boundaries = []rune{' ', '_', '-', '.', ',', ';', ':', '(', ')', '[', ']', '{', '}'}

func CamelCase(s string) (ret string) {
	atWordBoundary := true
	for _, c := range s {
		if slices.Contains(boundaries, c) {
			atWordBoundary = true
			continue
		}
		if atWordBoundary && unicode.IsLower(c) {
			c = unicode.ToUpper(c)
		}
		atWordBoundary = false
		ret += string(c)
	}
	return
}

func OrderedKeys[Map map[K]V, K cmp.Ordered, V any](m Map) []K {
	keysIter := maps.Keys(m)
	keys := slices.Collect(keysIter)
	slices.Sort(keys)
	return keys
}

func ReverseMap[M map[K]V, K cmp.Ordered, V cmp.Ordered](m M) (ret map[V]K) {
	ret = map[V]K{}
	for k, v := range m {
		ret[v] = k
	}
	return
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func PopFront[T any](s []T) (front T, rest []T, ok bool) {
	if len(s) == 0 {
		return front, s, false
	}
	return s[0], s[1:], true
}

func CombineNamAndData(name string, data any) []any {
	return []any{name, data}
}
