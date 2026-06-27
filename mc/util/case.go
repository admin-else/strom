package util

import (
	"slices"
	"unicode"
)

var boundaries = []rune{' ', '_', '-', '.', ',', ';', ':', '(', ')', '[', ']', '{', '}'}

func CamelCase(s string) string {
	ret := ""
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
	return ret
}

func SnakeCase(s string) string {
	ret := ""
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 {
				ret += "_"
			}
			c = unicode.ToLower(c)
		}
		ret += string(c)
	}
	return ret
}

func FirstLetterLower(s string) string {
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}
