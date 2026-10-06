package util

import (
	"math"
	"testing"
)

func TestMthSinCosTable(t *testing.T) {
	cases := []struct {
		in   float64
		sin  float32
		cos  float32
		desc string
	}{
		{0, 0, 1, "0"},
		{math.Pi / 2, 1, 0, "pi/2"},
		{math.Pi, 0, -1, "pi"},
		{-math.Pi / 2, -1, 0, "-pi/2"},
	}
	for _, c := range cases {
		if got := Sin(c.in); got != c.sin {
			t.Fatalf("Sin(%s): got %v want %v", c.desc, got, c.sin)
		}
		if got := Cos(c.in); got != c.cos {
			t.Fatalf("Cos(%s): got %v want %v", c.desc, got, c.cos)
		}
	}
}

func TestMthSinMatchesStdlibWithinQuantization(t *testing.T) {
	for i := -1000; i <= 1000; i++ {
		x := float64(i) * 0.01
		if diff := math.Abs(float64(Sin(x)) - math.Sin(x)); diff > 1.0/10430.0 {
			t.Fatalf("Sin(%v) diverges from math.Sin by %v", x, diff)
		}
	}
}

func TestMthLerpAndSquare(t *testing.T) {
	if got := LerpFloat(0.25, 0, 4); got != 1 {
		t.Fatalf("LerpFloat: got %v want 1", got)
	}
	if got := LerpDouble(0.25, 0, 4); got != 1 {
		t.Fatalf("LerpDouble: got %v want 1", got)
	}
	if got := SquareDouble(3); got != 9 {
		t.Fatalf("SquareDouble: got %v want 9", got)
	}
	if got := LengthSquared3Double(1, 2, 2); got != 9 {
		t.Fatalf("LengthSquared3Double: got %v want 9", got)
	}
	if got := LengthSquared2Double(3, 4); got != 25 {
		t.Fatalf("LengthSquared2Double: got %v want 25", got)
	}
}

func TestMthClampedLerp(t *testing.T) {
	if got := ClampedLerpFloat(-1, 2, 8); got != 2 {
		t.Fatalf("clamped low: got %v want 2", got)
	}
	if got := ClampedLerpFloat(2, 2, 8); got != 8 {
		t.Fatalf("clamped high: got %v want 8", got)
	}
	if got := ClampedLerpDouble(0.5, 2, 8); got != 5 {
		t.Fatalf("clamped mid: got %v want 5", got)
	}
}

func TestMthAbsIntNoFloatRoundtrip(t *testing.T) {
	if got := AbsInt(-5); got != 5 {
		t.Fatalf("AbsInt(-5): got %v want 5", got)
	}
	if got := AbsInt(7); got != 7 {
		t.Fatalf("AbsInt(7): got %v want 7", got)
	}
}

func TestMthNarrowingSaturates(t *testing.T) {
	if got := FloorDouble(3e9); got != math.MaxInt32 {
		t.Fatalf("FloorDouble(3e9): got %v want MaxInt32", got)
	}
	if got := CeilDouble(-3e9); got != math.MinInt32 {
		t.Fatalf("CeilDouble(-3e9): got %v want MinInt32", got)
	}
	if got := FloorDouble(math.NaN()); got != 0 {
		t.Fatalf("FloorDouble(NaN): got %v want 0", got)
	}
	if got := FloorDouble(math.Inf(1)); got != math.MaxInt32 {
		t.Fatalf("FloorDouble(+Inf): got %v want MaxInt32", got)
	}
}

func TestMthRadToDegIsFloatDivision(t *testing.T) {
	if got, want := MthRadToDeg, float32(180.0)/float32(math.Pi); got != want {
		t.Fatalf("MthRadToDeg: got %v want %v", got, want)
	}
}
