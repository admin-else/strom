package util

import "fmt"

// Comparable mirrors java.lang.Comparable.
type Comparable[T any] interface {
	CompareTo(T) int
}

// Integer mirrors java.lang.Integer when it is used through the java.lang
// Comparable bound (Go's int cannot carry a CompareTo method).
type Integer int

// CompareTo mirrors Integer.compareTo(Integer).
func (i Integer) CompareTo(other Integer) int {
	switch {
	case i < other:
		return -1
	case i > other:
		return 1
	default:
		return 0
	}
}

// InclusiveRange mirrors net.minecraft.util.InclusiveRange.
type InclusiveRange[T Comparable[T]] struct {
	MinInclusive T
	MaxInclusive T
}

// NewInclusiveRange mirrors the canonical InclusiveRange(T, T) record
// constructor, including its monotonicity check.
func NewInclusiveRange[T Comparable[T]](minInclusive T, maxInclusive T) InclusiveRange[T] {
	if minInclusive.CompareTo(maxInclusive) > 0 {
		panic("min_inclusive must be less than or equal to max_inclusive")
	}
	return InclusiveRange[T]{MinInclusive: minInclusive, MaxInclusive: maxInclusive}
}

// NewInclusiveRangeOf mirrors the InclusiveRange(T) single-value constructor.
func NewInclusiveRangeOf[T Comparable[T]](value T) InclusiveRange[T] {
	return NewInclusiveRange(value, value)
}

// TryCreateInclusiveRange mirrors InclusiveRange.create(T, T), which returns a
// DataResult instead of throwing on a non-monotonic range.
func TryCreateInclusiveRange[T Comparable[T]](minInclusive T, maxInclusive T) (result InclusiveRange[T], err error) {
	if minInclusive.CompareTo(maxInclusive) <= 0 {
		return InclusiveRange[T]{MinInclusive: minInclusive, MaxInclusive: maxInclusive}, nil
	}
	return result, fmt.Errorf("min_inclusive must be less than or equal to max_inclusive")
}

// MapInclusiveRange mirrors InclusiveRange.map(Function). Java declares a
// generic method (<S> InclusiveRange<S> map(...)); Go methods cannot introduce
// type parameters, so the mapped element type is a free-function parameter.
func MapInclusiveRange[T Comparable[T], S Comparable[S]](r InclusiveRange[T], mapper func(T) S) InclusiveRange[S] {
	return NewInclusiveRange(mapper(r.MinInclusive), mapper(r.MaxInclusive))
}

// IsValueInRange mirrors InclusiveRange.isValueInRange(T).
func (r InclusiveRange[T]) IsValueInRange(value T) bool {
	return value.CompareTo(r.MinInclusive) >= 0 && value.CompareTo(r.MaxInclusive) <= 0
}

// Contains mirrors InclusiveRange.contains(InclusiveRange<T>).
func (r InclusiveRange[T]) Contains(subRange InclusiveRange[T]) bool {
	return subRange.MinInclusive.CompareTo(r.MinInclusive) >= 0 && subRange.MaxInclusive.CompareTo(r.MaxInclusive) <= 0
}

// String mirrors InclusiveRange.toString().
func (r InclusiveRange[T]) String() string {
	return "[" + fmt.Sprint(r.MinInclusive) + ", " + fmt.Sprint(r.MaxInclusive) + "]"
}
