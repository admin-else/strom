package util

import "fmt"

func MustOk[T any](v T, ok bool) T {
	if !ok {
		panic("not ok")
	}
	return v
}

func MustT[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertTypeError[T any](v any) (T, error) {
	var zero T
	r, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed: %T cannot be asserted to %T", v, zero)
	}
	return r, nil
}

func MapGetError[K comparable, V any, M map[K]V](m M, k K) (V, error) {
	v, ok := m[k]
	if !ok {
		return v, fmt.Errorf("key %v not found", k)
	}
	return v, nil
}
