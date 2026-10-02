package util

import "testing"

func TestInclusiveRangeBasics(t *testing.T) {
	r := NewInclusiveRange(Integer(2), Integer(5))
	if !r.IsValueInRange(Integer(2)) || !r.IsValueInRange(Integer(5)) || !r.IsValueInRange(Integer(3)) {
		t.Fatal("expected values inside range")
	}
	if r.IsValueInRange(Integer(1)) || r.IsValueInRange(Integer(6)) {
		t.Fatal("expected values outside range")
	}
	if got := r.String(); got != "[2, 5]" {
		t.Fatalf("String() = %q", got)
	}
	single := NewInclusiveRangeOf(Integer(7))
	if !single.Contains(single) {
		t.Fatal("expected single-value range to contain itself")
	}
	if !r.Contains(NewInclusiveRange(Integer(3), Integer(4))) {
		t.Fatal("expected range to contain subrange")
	}
	if r.Contains(NewInclusiveRange(Integer(0), Integer(4))) {
		t.Fatal("expected range not to contain lower subrange")
	}
}

func TestInclusiveRangePanicsOnInverted(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for inverted range")
		}
	}()
	NewInclusiveRange(Integer(5), Integer(2))
}

func TestMapInclusiveRange(t *testing.T) {
	r := NewInclusiveRange(Integer(1), Integer(3))
	mapped := MapInclusiveRange(r, func(i Integer) Integer { return i * 10 })
	if mapped.MinInclusive != 10 || mapped.MaxInclusive != 30 {
		t.Fatalf("mapped = %v", mapped)
	}
}
