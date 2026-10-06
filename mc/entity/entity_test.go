package entity

import (
	"testing"

	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
)

type testLevel struct {
	blockColliders []*shapes.VoxelShape
}

func (l testLevel) GetEntityCollisions(source *Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	return nil
}

func (l testLevel) GetBlockCollisions(source *Entity, box phys.AABB) (ret []*shapes.VoxelShape) {
	for _, shape := range l.blockColliders {
		if bounds, ok := shape.Bounds(); ok && box.Intersects(bounds) {
			ret = append(ret, shape)
		}
	}
	return
}

func TestEntityFallsAndRestsOnFloor(t *testing.T) {
	level := testLevel{blockColliders: []*shapes.VoxelShape{
		shapes.Box(-10, -1, -10, 10, 0, 10),
	}}
	e := NewEntity(level, 0.5, 5, 0.5, 0.6, 1.8)
	for i := 0; i < 40; i++ {
		delta := phys.NewVec3(0, -0.5, 0)
		e.Move(MoverTypeSELF, delta)
	}
	if e.GetY() != 0.0 {
		t.Fatalf("expected to rest at y=0, got %v", e.GetY())
	}
	if !e.OnGround() {
		t.Fatalf("expected onGround after landing")
	}
}

func TestEntityStopsAtWall(t *testing.T) {
	level := testLevel{blockColliders: []*shapes.VoxelShape{
		shapes.Box(1, 0, -10, 2, 2, 10),
	}}
	e := NewEntity(level, 0, 0, 0, 0.6, 1.8)
	e.Move(MoverTypeSELF, phys.NewVec3(1, 0, 0))
	if !almost(e.GetX(), 0.7, 1e-9) {
		t.Fatalf("expected x=0.7 after wall, got %v", e.GetX())
	}
	if !e.HorizontalCollision() {
		t.Fatalf("expected horizontal collision")
	}
}

func TestEntityStepsUp(t *testing.T) {
	level := testLevel{blockColliders: []*shapes.VoxelShape{
		shapes.Box(-10, -1, -10, 10, 0, 10),
		shapes.Box(1, 0, -1, 3, 0.5, 1),
	}}
	e := NewEntity(level, 0, 0, 0, 0.6, 1.8)
	e.Move(MoverTypeSELF, phys.NewVec3(0, -0.1, 0)) // settle on floor
	e.Move(MoverTypeSELF, phys.NewVec3(1, 0, 0))
	if !almost(e.GetY(), 0.5, 1e-9) {
		t.Fatalf("expected step up to y=0.5, got %v", e.GetY())
	}
}

func TestInputForward(t *testing.T) {
	xxa, zza := ModifyInput(Input{Forward: true}.GetMoveVector(), false).X, ModifyInput(Input{Forward: true}.GetMoveVector(), false).Y
	if xxa != 0.0 || !almost(float64(zza), 0.98, 1e-6) {
		t.Fatalf("forward input: got (%v,%v) want (0,0.98)", xxa, zza)
	}
}

func TestInputDiagonalClampsToUnit(t *testing.T) {
	got := ModifyInput(Input{Forward: true, Right: true}.GetMoveVector(), false)
	if !almost(float64(got.Length()), 1.0, 1e-6) {
		t.Fatalf("diagonal input length: got %v want 1", got.Length())
	}
}

func TestInputSneakScales(t *testing.T) {
	full := ModifyInput(Input{Forward: true}.GetMoveVector(), false)
	sneak := ModifyInput(Input{Forward: true}.GetMoveVector(), true)
	if !(sneak.Length() < full.Length()) {
		t.Fatalf("sneaking should be slower: full=%v sneak=%v", full.Length(), sneak.Length())
	}
}

func TestInputSneakDiagonalLength(t *testing.T) {
	got := ModifyInput(Input{Forward: true, Right: true}.GetMoveVector(), true)
	if !almost(float64(got.Length()), 0.41578, 1e-4) {
		t.Fatalf("sneak diagonal length: got %v want ~0.41578", got.Length())
	}
}

func almost(a, b, eps float64) bool {
	if a > b {
		return a-b <= eps
	}
	return b-a <= eps
}

func TestFallDistanceAccumulatesAndResets(t *testing.T) {
	level := testLevel{blockColliders: []*shapes.VoxelShape{
		shapes.Box(-10, -1, -10, 10, 0, 10),
	}}
	e := NewEntity(level, 0.5, 5, 0.5, 0.6, 1.8)
	for i := 0; i < 3; i++ {
		e.Move(MoverTypeSELF, phys.NewVec3(0, -1, 0))
	}
	if e.GetFallDistance() <= 0 {
		t.Fatalf("expected accumulating fall distance, got %v", e.GetFallDistance())
	}
	for i := 0; i < 10; i++ {
		e.Move(MoverTypeSELF, phys.NewVec3(0, -1, 0))
	}
	if !e.OnGround() || e.GetFallDistance() != 0 {
		t.Fatalf("expected fall distance reset on landing, got %v onGround=%v", e.GetFallDistance(), e.OnGround())
	}
}

func TestJumpBoostIncreasesJumpPower(t *testing.T) {
	level := testLevel{}
	e := NewEntity(level, 0.5, 0, 0.5, 0.6, 1.8)
	base := e.GetJumpPower(1.0, 1.0)
	e.AddEffect(MobEffectJUMP_BOOST, 1, -1)
	boosted := e.GetJumpPower(1.0, 1.0)
	if !almost(float64(base), 0.42, 1e-6) {
		t.Fatalf("base jump power: got %v want 0.42", base)
	}
	if !almost(float64(boosted), 0.62, 1e-6) {
		t.Fatalf("boosted jump power: got %v want 0.62", boosted)
	}
}

func TestSlowFallingReducesGravity(t *testing.T) {
	level := testLevel{}
	e := NewEntity(level, 0.5, 0, 0.5, 0.6, 1.8)
	e.SetDeltaMovement(phys.NewVec3(0, -1, 0))
	e.AddEffect(MobEffectSLOW_FALLING, 0, -1)
	if got := e.GetEffectiveGravity(); !almost(got, 0.01, 1e-9) {
		t.Fatalf("slow falling gravity: got %v want 0.01", got)
	}
}

func TestHandleOnClimbableClampsMovement(t *testing.T) {
	level := testLevel{}
	e := NewEntity(level, 0.5, 0, 0.5, 0.6, 1.8)
	// OnClimbable is false for a plain CollisionGetter, so the delta is unchanged.
	if got := e.HandleOnClimbable(phys.NewVec3(1, -2, 3)); got != phys.NewVec3(1, -2, 3) {
		t.Fatalf("non-climbable delta must pass through, got %+v", got)
	}
}

func TestTickEffectsExpires(t *testing.T) {
	e := NewEntity(testLevel{}, 0, 0, 0, 0.6, 1.8)
	e.AddEffect(MobEffectJUMP_BOOST, 0, 2)
	e.TickEffects()
	e.TickEffects()
	if e.HasEffect(MobEffectJUMP_BOOST) {
		t.Fatalf("timed effect should expire")
	}
	e.AddEffect(MobEffectSLOW_FALLING, 0, -1)
	for i := 0; i < 5; i++ {
		e.TickEffects()
	}
	if !e.HasEffect(MobEffectSLOW_FALLING) {
		t.Fatalf("infinite effect should persist")
	}
}
