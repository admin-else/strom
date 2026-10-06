package phys

import (
	"math"
	"testing"
)

func almost(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestAABBConstructorNormalizes(t *testing.T) {
	got := NewAABB(3, 4, 5, 0, 1, 2)
	want := NewAABB(0, 1, 2, 3, 4, 5)
	if !got.Equal(want) {
		t.Fatalf("normalize: got %+v want %+v", got, want)
	}
}

func TestAABBContractExpand(t *testing.T) {
	box := NewAABB(0, 0, 0, 1, 1, 1)
	if got, want := box.Contract(0.5, 0, 0), NewAABB(0, 0, 0, 0.5, 1, 1); !got.Equal(want) {
		t.Fatalf("contract: got %+v want %+v", got, want)
	}
	if got, want := box.Contract(-0.5, 0, 0), NewAABB(0.5, 0, 0, 1, 1, 1); !got.Equal(want) {
		t.Fatalf("contract neg: got %+v want %+v", got, want)
	}
	if got, want := box.ExpandTowardsXYZ(0.5, 0, 0), NewAABB(0, 0, 0, 1.5, 1, 1); !got.Equal(want) {
		t.Fatalf("expand: got %+v want %+v", got, want)
	}
	if got, want := box.ExpandTowardsXYZ(-0.5, 0, 0), NewAABB(-0.5, 0, 0, 1, 1, 1); !got.Equal(want) {
		t.Fatalf("expand neg: got %+v want %+v", got, want)
	}
}

func TestAABBIntersectMinmax(t *testing.T) {
	a := NewAABB(0, 0, 0, 1, 1, 1)
	b := NewAABB(0.5, 0.5, 0.5, 2, 2, 2)
	if got, want := a.Intersect(b), NewAABB(0.5, 0.5, 0.5, 1, 1, 1); !got.Equal(want) {
		t.Fatalf("intersect: got %+v want %+v", got, want)
	}
	if got, want := a.Minmax(b), NewAABB(0, 0, 0, 2, 2, 2); !got.Equal(want) {
		t.Fatalf("minmax: got %+v want %+v", got, want)
	}
}

func TestAABBContainsIsHalfOpen(t *testing.T) {
	box := NewAABB(0, 0, 0, 1, 1, 1)
	if !box.ContainsXYZ(0, 0, 0) || !box.ContainsXYZ(0.999, 0.999, 0.999) {
		t.Fatalf("expected inside points to be contained")
	}
	if box.ContainsXYZ(1, 1, 1) {
		t.Fatalf("max corner must be exclusive")
	}
}

func TestAABBDistanceAndGeometry(t *testing.T) {
	box := NewAABB(0, 0, 0, 1, 1, 1)
	if got := box.DistanceToSqrVec3(NewVec3(2, 0, 0)); got != 1 {
		t.Fatalf("distanceToSqr: got %v want 1", got)
	}
	if got, want := OfSize(NewVec3(0, 0, 0), 2, 2, 2), NewAABB(-1, -1, -1, 1, 1, 1); !got.Equal(want) {
		t.Fatalf("ofSize: got %+v want %+v", got, want)
	}
	if got := box.GetCenter(); !almost(got.X, 0.5, 0) || !almost(got.Y, 0.5, 0) || !almost(got.Z, 0.5, 0) {
		t.Fatalf("center: got %+v", got)
	}
}

func TestAABBNextDeflated(t *testing.T) {
	box := NewAABB(0, 0, 0, 1, 1, 1).NextDeflated()
	if !(box.MinX > 0 && box.MaxX < 1) {
		t.Fatalf("nextDeflated must shrink: %+v", box)
	}
}

func TestVec3BasicAlgebra(t *testing.T) {
	a := NewVec3(1, 2, 3)
	b := NewVec3(4, 5, 6)
	if got, want := a.Add(b), NewVec3(5, 7, 9); got != want {
		t.Fatalf("add: got %+v want %+v", got, want)
	}
	if got, want := a.Subtract(b), NewVec3(-3, -3, -3); got != want {
		t.Fatalf("subtract: got %+v want %+v", got, want)
	}
	if got := a.Dot(b); got != 32 {
		t.Fatalf("dot: got %v want 32", got)
	}
	if got, want := a.Cross(b), NewVec3(-3, 6, -3); got != want {
		t.Fatalf("cross: got %+v want %+v", got, want)
	}
	if got, want := a.Scale(2), NewVec3(2, 4, 6); got != want {
		t.Fatalf("scale: got %+v want %+v", got, want)
	}
}

func TestVec3NormalizeTinyIsZero(t *testing.T) {
	if got := NewVec3(1e-9, 0, 0).Normalize(); got != Vec3ZERO {
		t.Fatalf("tiny normalize: got %+v want zero", got)
	}
	if got := NewVec3(3, 0, 4).Normalize(); !almost(got.X, 0.6, 1e-12) || !almost(got.Z, 0.8, 1e-12) {
		t.Fatalf("normalize: got %+v", got)
	}
}

func TestVec3RotationDegrees(t *testing.T) {
	cases := []struct {
		v    Vec3
		yaw  float32
		psh  float32
		desc string
	}{
		{NewVec3(0, 0, 1), 0, 0, "+Z"},
		{NewVec3(-1, 0, 0), 90, 0, "-X"},
		{NewVec3(0, 1, 0), 0, -90, "+Y"},
	}
	for _, c := range cases {
		got := c.v.Rotation()
		if !almost(float64(got.X), float64(c.psh), 1e-4) || !almost(float64(got.Y), float64(c.yaw), 1e-4) {
			t.Fatalf("%s: rotation got (pitch=%v yaw=%v) want (pitch=%v yaw=%v)", c.desc, got.X, got.Y, c.psh, c.yaw)
		}
	}
}

func TestVec2RotateQuarterTurn(t *testing.T) {
	got := NewVec2(1, 0).Rotate(math.Pi / 2)
	if !almost(float64(got.X), 0, 1e-4) || !almost(float64(got.Y), 1, 1e-4) {
		t.Fatalf("rotate: got %+v want (0,1)", got)
	}
}

func TestDirectionFromRotationStraightUp(t *testing.T) {
	got := DirectionFromRotation(-90, 0)
	if !almost(got.Y, 1, 1e-4) || !almost(got.X, 0, 1e-4) || !almost(got.Z, 0, 1e-4) {
		t.Fatalf("directionFromRotation(-90,0): got %+v want straight up", got)
	}
}

func TestVec3EqualDoubleCompareSemantics(t *testing.T) {
	if NewVec3(math.NaN(), 0, 0).Equal(NewVec3(math.NaN(), 0, 0)) != true {
		t.Fatalf("NaN must equal NaN")
	}
	if NewVec3(0, 0, 0).Equal(NewVec3(math.Copysign(0, -1), 0, 0)) != false {
		t.Fatalf("0.0 must not equal -0.0")
	}
}

func TestVec2NormalizedTinyUsesFloatEpsilon(t *testing.T) {
	if got := NewVec2(1e-6, 0).Normalized(); !got.Equal(Vec2ZERO) {
		t.Fatalf("tiny normalize: got %+v want zero", got)
	}
}
