package shapes

import (
	"testing"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/phys"
)

func TestShapesBlockBounds(t *testing.T) {
	bounds, ok := Block().Bounds()
	if !ok {
		t.Fatalf("block bounds missing")
	}
	want := phys.NewAABB(0, 0, 0, 1, 1, 1)
	if !bounds.Equal(want) {
		t.Fatalf("block bounds: got %+v want %+v", bounds, want)
	}
}

func TestShapesBoxBounds(t *testing.T) {
	shape := Box(0.25, 0.0, 0.25, 0.75, 1.0, 0.75)
	bounds, _ := shape.Bounds()
	want := phys.NewAABB(0.25, 0.0, 0.25, 0.75, 1.0, 0.75)
	if !bounds.Equal(want) {
		t.Fatalf("box bounds: got %+v want %+v", bounds, want)
	}
}

func TestShapesOrOfTwoHalvesIsBlock(t *testing.T) {
	left := Box(0, 0, 0, 0.5, 1, 1)
	right := Box(0.5, 0, 0, 1, 1, 1)
	joined := Join(left, right, BooleanOpOR)
	if !Equal(joined, Block()) {
		t.Fatalf("join of halves should equal block; got bounds %+v", mustBounds(joined))
	}
}

func TestShapesCollideAxisX(t *testing.T) {
	shape := Block()
	moving := phys.NewAABB(2, 0, 0, 3, 1, 1)
	got := shape.Collide(core.AxisX, moving, -3.0)
	if got != -1.0 {
		t.Fatalf("collide: got %v want -1", got)
	}
	if got := shape.Collide(core.AxisX, moving, 2.0); got != 2.0 {
		t.Fatalf("collide away: got %v want 2", got)
	}
}

func TestShapesCollideMany(t *testing.T) {
	shapeList := []*VoxelShape{Box(0, 0, 0, 1, 1, 1)}
	moving := phys.NewAABB(2, 0, 0, 3, 1, 1)
	if got := Collide(core.AxisX, moving, shapeList, -3.0); got != -1.0 {
		t.Fatalf("shapes collide: got %v want -1", got)
	}
}

func TestShapesBlockOccludes(t *testing.T) {
	if !BlockOccludes(Block(), Block(), core.DirectionUP) {
		t.Fatalf("block should occlude block upward")
	}
	if BlockOccludes(Empty(), Block(), core.DirectionUP) {
		t.Fatalf("empty shape should not occlude")
	}
}

func TestShapesIsEmpty(t *testing.T) {
	if !Empty().IsEmpty() {
		t.Fatalf("Empty() must be empty")
	}
	if Block().IsEmpty() {
		t.Fatalf("Block() must not be empty")
	}
}

func TestVoxelShapeClip(t *testing.T) {
	shape := Block()
	hit := shape.Clip(phys.NewVec3(-1, 0.5, 0.5), phys.NewVec3(1, 0.5, 0.5), core.NewBlockPos(0, 0, 0))
	if hit == nil {
		t.Fatalf("expected a hit")
	}
	if hit.GetDirection() != core.DirectionWEST {
		t.Fatalf("hit face: got %s want west", hit.GetDirection().Name())
	}
	if got := hit.GetLocation().X; got != 0.0 {
		t.Fatalf("hit x: got %v want 0", got)
	}
}

func TestAABBClip(t *testing.T) {
	got, ok := phys.AABBClip(0, 0, 0, 1, 1, 1, phys.NewVec3(-1, 0.5, 0.5), phys.NewVec3(1, 0.5, 0.5))
	if !ok {
		t.Fatalf("expected clip hit")
	}
	if got.X != 0.0 || got.Y != 0.5 || got.Z != 0.5 {
		t.Fatalf("clip point: got %+v", got)
	}
}

func mustBounds(shape *VoxelShape) (ret phys.AABB) {
	ret, _ = shape.Bounds()
	return
}
