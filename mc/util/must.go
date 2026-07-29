package util

import "fmt"

// MustOk returns v if ok is true, otherwise panics.
func MustOk[T any](v T, ok bool) T {
	if !ok {
		panic("not ok")
	}
	return v
}

// MustT returns v if err is nil, otherwise panics.
func MustT[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Must panics if err is non-nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// AssertTypeError asserts that v is of type T and returns it, or returns an error.
func AssertTypeError[T any](v any) (T, error) {
	var zero T
	r, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed: %T cannot be asserted to %T", v, zero)
	}
	return r, nil
}

// MapGetError returns the value for key k from map m, or an error if not present.
func MapGetError[K comparable, V any, M map[K]V](m M, k K) (V, error) {
	v, ok := m[k]
	if !ok {
		return v, fmt.Errorf("key %v not found", k)
	}
	return v, nil
}
