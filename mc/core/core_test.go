package core

import "testing"

func TestDirection2DDataValues(t *testing.T) {
	cases := map[Direction]int{
		DirectionSOUTH: 0,
		DirectionWEST:  1,
		DirectionNORTH: 2,
		DirectionEAST:  3,
	}
	for d, want := range cases {
		if got := d.Get2DDataValue(); got != want {
			t.Fatalf("%s get2DDataValue: got %d want %d", d.Name(), got, want)
		}
	}
	if DirectionDOWN.Get2DDataValue() != -1 || DirectionUP.Get2DDataValue() != -1 {
		t.Fatalf("vertical directions must have data2d -1")
	}
}

func TestDirectionFrom2DDataValueOrder(t *testing.T) {
	want := []Direction{DirectionSOUTH, DirectionWEST, DirectionNORTH, DirectionEAST}
	for i, w := range want {
		if got := DirectionFrom2DDataValue(i); got != w {
			t.Fatalf("from2DDataValue(%d): got %s want %s", i, got.Name(), w.Name())
		}
	}
	if got := DirectionFrom2DDataValue(4); got != DirectionSOUTH {
		t.Fatalf("from2DDataValue wrap: got %s want south", got.Name())
	}
}

func TestDirectionNormalsAndAxis(t *testing.T) {
	if got := DirectionEAST.GetUnitVec3i(); got != NewVec3i(1, 0, 0) {
		t.Fatalf("east normal: got %+v", got)
	}
	if got := DirectionDOWN.GetOpposite(); got != DirectionUP {
		t.Fatalf("down opposite: got %s", got.Name())
	}
	if got := DirectionWEST.GetAxis(); got != AxisX {
		t.Fatalf("west axis: got %v", got)
	}
	if got := DirectionFromAxisAndDirection(AxisY, DirectionAxisDirectionNEGATIVE); got != DirectionDOWN {
		t.Fatalf("fromAxisAndDirection: got %s", got.Name())
	}
}

func TestDirectionClockWise(t *testing.T) {
	if got := DirectionNORTH.GetClockWise(); got != DirectionEAST {
		t.Fatalf("north cw: got %s", got.Name())
	}
	if got := DirectionNORTH.GetCounterClockWise(); got != DirectionWEST {
		t.Fatalf("north ccw: got %s", got.Name())
	}
}

func TestBlockPosPackingRoundTrip(t *testing.T) {
	coords := []BlockPos{
		NewBlockPos(0, 0, 0),
		NewBlockPos(1, 2, 3),
		NewBlockPos(-1, 0, -1),
		NewBlockPos(33554431, 255, -33554431),
		NewBlockPos(-33554432, -64, 33554431),
	}
	for _, p := range coords {
		node := p.AsLong()
		got := BlockPosOfPackaged(node)
		if !got.Equal(p) {
			t.Fatalf("packing round trip: %+v -> %+v", p, got)
		}
	}
}

func TestBlockPosContainingFloors(t *testing.T) {
	if got := BlockPosContaining(-0.5, 1.0, 2.9); got != NewBlockPos(-1, 1, 2) {
		t.Fatalf("containing: got %+v", got)
	}
}

func TestBetweenClosedOrderAndCount(t *testing.T) {
	got := BetweenClosedXYZ(0, 0, 0, 1, 1, 1)
	if len(got) != 8 {
		t.Fatalf("count: got %d want 8", len(got))
	}
	// x fastest, then y, then z.
	if got[0] != NewBlockPos(0, 0, 0) || got[1] != NewBlockPos(1, 0, 0) ||
		got[2] != NewBlockPos(0, 1, 0) || got[4] != NewBlockPos(0, 0, 1) {
		t.Fatalf("order wrong: %+v", got[:5])
	}
}

func TestVec3iAlgebra(t *testing.T) {
	v := NewVec3i(1, 2, 3)
	if got := v.East().AboveSteps(2); got != NewVec3i(2, 4, 3) {
		t.Fatalf("relative: got %+v", got)
	}
	if got := NewVec3i(1, 0, 0).Cross(NewVec3i(0, 0, 1)); got != NewVec3i(0, -1, 0) {
		t.Fatalf("cross: got %+v", got)
	}
	if got := v.Get(AxisY); got != 2 {
		t.Fatalf("get axis: got %d", got)
	}
}
