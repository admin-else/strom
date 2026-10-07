package proto_base

// OptionalDeref returns the value pointed to by p, or the zero value of T when
// p is nil. Generated switch comparisons over optional fields use it so an
// absent optional falls through to the switch default instead of panicking.
func OptionalDeref[T any](p *T) T {
	var zero T
	if p != nil {
		return *p
	}
	return zero
}
