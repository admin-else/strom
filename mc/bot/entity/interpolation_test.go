package entity

import (
	"math"
	"testing"
)

func TestSteppedInterpolationMovesTowardTargetAndLerps(t *testing.T) {
	e := newEntity()
	e.SnapTo(Vec3{}, 0, 0)
	e.MoveOrInterpolateToLinear(Vec3{X: 1}, 0, 0)

	e.Tick()
	if x := e.Position(1.0).X; math.Abs(x-1.0/3.0) > 1e-6 {
		t.Fatalf("after tick 1 x=%v, want 1/3", x)
	}
	e.Tick()
	if x := e.Position(1.0).X; math.Abs(x-2.0/3.0) > 1e-6 {
		t.Fatalf("after tick 2 x=%v, want 2/3", x)
	}
	e.Tick()
	if x := e.Position(1.0).X; math.Abs(x-1.0) > 1e-6 {
		t.Fatalf("after tick 3 x=%v, want 1", x)
	}
	if e.interpolation.HasActiveInterpolation() {
		t.Fatal("interpolation should be done after reaching the target")
	}
	// partial tick lerps from the previous-tick snapshot toward the live pose
	mid := e.Position(0.5).X
	if mid <= 2.0/3.0 || mid >= 1.0 {
		t.Fatalf("partial-tick x=%v, want between 2/3 and 1", mid)
	}
}

func TestLinearInterpolationHandlerConverges(t *testing.T) {
	e := newEntity()
	e.SnapTo(Vec3{}, 0, 0)
	e.interpolation = NewLinearInterpolationHandler(e, 3)
	e.MoveOrInterpolateToLinear(Vec3{X: 9}, 0, 0)
	want := []float64{3, 6, 9}
	for i, w := range want {
		e.Tick()
		if x := e.Position(1.0).X; math.Abs(x-w) > 1e-6 {
			t.Fatalf("step %d x=%v, want %v", i, x, w)
		}
	}
}

func TestInterpolatedRotationUsesRotLerp(t *testing.T) {
	e := newEntity()
	e.SnapTo(Vec3{}, 350, 0)
	e.MoveOrInterpolateToLinear(Vec3{}, 10, 0) // +20 degrees the short way
	e.Tick()
	yaw := e.Yaw(1.0)
	// rotLerp takes the -340 -> +20 wrapped delta, so one step is 350 + 20/3
	if math.Abs(float64(yaw)-356.66666) > 1e-3 {
		t.Fatalf("yaw=%v, want ~356.67 (wrapped the short way)", yaw)
	}
}

func TestAgeInTicksAddsPartialTick(t *testing.T) {
	e := newEntity()
	e.Tick()
	e.Tick()
	if got := e.AgeInTicks(0.25); got != 2.25 {
		t.Fatalf("AgeInTicks=%v, want 2.25", got)
	}
}
