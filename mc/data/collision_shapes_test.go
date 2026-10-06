package data

import "testing"

func TestCollisionShapeStoneIsFullCube(t *testing.T) {
	stone, ok := LookupBlockByName("26.4-snapshot-2", "stone")
	if !ok {
		t.Fatalf("stone not found")
	}
	boxes, ok := CollisionShapeBoxesAt("26.4-snapshot-2", stone.DefaultState)
	if !ok {
		t.Fatalf("stone has no collision shape")
	}
	if len(boxes) != 1 {
		t.Fatalf("stone boxes: got %d want 1", len(boxes))
	}
	want := [6]float64{0, 0, 0, 1, 1, 1}
	if boxes[0] != want {
		t.Fatalf("stone box: got %v want %v", boxes[0], want)
	}
}

func TestCollisionShapeAirIsEmpty(t *testing.T) {
	air, ok := LookupBlockByName("26.4-snapshot-2", "air")
	if !ok {
		t.Fatalf("air not found")
	}
	if _, ok := CollisionShapeBoxesAt("26.4-snapshot-2", air.DefaultState); ok {
		t.Fatalf("air should have no collision shape")
	}
}

func TestCollisionShapeSlabIsPerState(t *testing.T) {
	slab, ok := LookupBlockByName("26.4-snapshot-2", "oak_slab")
	if !ok {
		t.Fatalf("oak_slab not found")
	}
	boxes, ok := CollisionShapeBoxesAt("26.4-snapshot-2", slab.DefaultState)
	if !ok {
		t.Fatalf("oak_slab default state has no collision shape")
	}
	if len(boxes) == 0 || boxes[0][4] > 1.0 {
		t.Fatalf("oak_slab default state should be a partial shape, got %v", boxes)
	}
}
