package util

import "testing"

func TestMthIsMultipleOf(t *testing.T) {
	if !IsMultipleOf(8, 4) {
		t.Error("8 is a multiple of 4")
	}
	if IsMultipleOf(9, 4) {
		t.Error("9 is not a multiple of 4")
	}
}

func TestMthLog2(t *testing.T) {
	cases := map[int]int{1: 0, 2: 1, 4: 2, 8: 3, 3: 1, 5: 2, 16: 4}
	for input, want := range cases {
		if got := Log2(input); got != want {
			t.Errorf("Log2(%d) = %d, want %d", input, got, want)
		}
	}
}

func TestMthCeilLog2(t *testing.T) {
	cases := map[int]int{1: 0, 2: 1, 4: 2, 3: 2, 5: 3, 16: 4}
	for input, want := range cases {
		if got := CeilLog2(input); got != want {
			t.Errorf("CeilLog2(%d) = %d, want %d", input, got, want)
		}
	}
}

func TestMthRoundToward(t *testing.T) {
	if got := RoundToward(5, 4); got != 8 {
		t.Errorf("RoundToward(5,4) = %d, want 8", got)
	}
	if got := RoundToward(8, 4); got != 8 {
		t.Errorf("RoundToward(8,4) = %d, want 8", got)
	}
}

func TestMthClamp(t *testing.T) {
	if got := ClampFloat(2.0, 0.0, 1.0); got != 1.0 {
		t.Errorf("ClampFloat = %v, want 1.0", got)
	}
	if got := ClampInt(0, 1, 3); got != 1 {
		t.Errorf("ClampInt = %d, want 1", got)
	}
}
