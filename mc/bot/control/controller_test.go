package control

import (
	"math"
	"testing"
	"time"

	"github.com/admin-else/strom/mc/entity"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
)

type fakeTarget struct {
	x, y, z    float64
	yaw, pitch float32
}

func (f *fakeTarget) Position() (x float64, y float64, z float64) { return f.x, f.y, f.z }
func (f *fakeTarget) Rotation() (yaw float32, pitch float32)      { return f.yaw, f.pitch }
func (f *fakeTarget) SetPosition(x float64, y float64, z float64) { f.x, f.y, f.z = x, y, z }
func (f *fakeTarget) SetRotation(yaw float32, pitch float32)      { f.yaw, f.pitch = yaw, pitch }

// fakeLevel is an entity.CollisionGetter with a flat floor at y=0.
type fakeLevel struct {
	shapes []*shapes.VoxelShape
}

func (l fakeLevel) GetEntityCollisions(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	return nil
}

func (l fakeLevel) GetBlockCollisions(source *entity.Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	for _, shape := range l.shapes {
		if bounds, ok := shape.Bounds(); ok && box.Intersects(bounds) {
			ret = append(ret, shape)
		}
	}
	return
}

func newTestController(t *testing.T) (*Controller, *fakeTarget) {
	t.Helper()
	level := fakeLevel{shapes: []*shapes.VoxelShape{shapes.Box(-10, -1, -10, 10, 0, 10)}}
	target := &fakeTarget{}
	ctrl := New(level, target)
	for i := 0; i < 5; i++ { // settle onto the floor
		if err := ctrl.Tick(); err != nil {
			t.Fatalf("settle tick: %v", err)
		}
	}
	return ctrl, target
}

func TestControllerWalkForwardMovesAlongYaw(t *testing.T) {
	ctrl, target := newTestController(t)
	if !ctrl.Entity().OnGround() {
		t.Fatalf("expected to be on the floor before walking")
	}
	startZ := target.z
	ctrl.Move("w", time.Second)
	for i := 0; i < 20; i++ {
		if err := ctrl.Tick(); err != nil {
			t.Fatalf("tick: %v", err)
		}
	}
	if target.z <= startZ+0.5 {
		t.Fatalf("expected forward (+Z) movement, startZ=%v z=%v", startZ, target.z)
	}
}

func TestControllerLookAt(t *testing.T) {
	ctrl, _ := newTestController(t)
	ctrl.LookAt(0, 0, 10, 0)
	if err := ctrl.Tick(); err != nil {
		t.Fatalf("tick: %v", err)
	}
	if yaw, _ := ctrl.Rotation(); math.Abs(float64(yaw)) > 1e-3 {
		t.Fatalf("looking +Z should give yaw 0, got %v", yaw)
	}

	ctrl.LookAt(0, 0, -10, 0)
	if err := ctrl.Tick(); err != nil {
		t.Fatalf("tick: %v", err)
	}
	if yaw, _ := ctrl.Rotation(); math.Abs(math.Abs(float64(yaw))-180) > 1e-3 {
		t.Fatalf("looking -Z should give |yaw| 180, got %v", yaw)
	}
}

func TestControllerLookAtDelay(t *testing.T) {
	ctrl, _ := newTestController(t)
	ctrl.LookAt(0, 0, -10, 200*time.Millisecond) // 4 ticks
	ctrl.Tick()
	if yaw, _ := ctrl.Rotation(); yaw != 0 {
		t.Fatalf("rotation must not change before the delay, got %v", yaw)
	}
	for i := 0; i < 4; i++ {
		ctrl.Tick()
	}
	if yaw, _ := ctrl.Rotation(); math.Abs(math.Abs(float64(yaw))-180) > 1e-3 {
		t.Fatalf("rotation should apply after the delay, got %v", yaw)
	}
}

func TestControllerJumpSetsUpwardVelocity(t *testing.T) {
	ctrl, _ := newTestController(t)
	ctrl.Jump(time.Second)
	if err := ctrl.Tick(); err != nil {
		t.Fatalf("tick: %v", err)
	}
	if ctrl.Entity().GetDeltaMovement().Y <= 0 {
		t.Fatalf("jump should give upward velocity, got %v", ctrl.Entity().GetDeltaMovement().Y)
	}
}
