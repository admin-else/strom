package shapes

import "github.com/admin-else/strom/mc/core"

// DiscreteVoxelShape mirrors net.minecraft.world.phys.shapes.DiscreteVoxelShape.
// Java's abstract class with two concrete overrides is modelled as an interface
// (unexported methods, since only this package implements or consumes it); the
// concrete methods become package-level helpers.
//
// Deferred: rotate(OctahedralGroup) (JOML/OctahedralGroup), forAllEdges and
// forAllFaces (only used by the debug/rendering paths, not by collision).
type DiscreteVoxelShape interface {
	xSize() int
	ySize() int
	zSize() int
	isFull(x int, y int, z int) bool
	fill(x int, y int, z int)
	firstFull(axis core.DirectionAxis) int
	lastFull(axis core.DirectionAxis) int
}

// DvsGetSize mirrors DiscreteVoxelShape.getSize(Direction.Axis).
func DvsGetSize(s DiscreteVoxelShape, axis core.DirectionAxis) (ret int) {
	return axis.ChooseInt(s.xSize(), s.ySize(), s.zSize())
}

// DvsGetXSize mirrors DiscreteVoxelShape.getXSize().
func DvsGetXSize(s DiscreteVoxelShape) (ret int) { return DvsGetSize(s, core.AxisX) }

// DvsGetYSize mirrors DiscreteVoxelShape.getYSize().
func DvsGetYSize(s DiscreteVoxelShape) (ret int) { return DvsGetSize(s, core.AxisY) }

// DvsGetZSize mirrors DiscreteVoxelShape.getZSize().
func DvsGetZSize(s DiscreteVoxelShape) (ret int) { return DvsGetSize(s, core.AxisZ) }

// DvsIsFullWide mirrors DiscreteVoxelShape.isFullWide(int, int, int).
func DvsIsFullWide(s DiscreteVoxelShape, x int, y int, z int) (ret bool) {
	if x < 0 || y < 0 || z < 0 {
		return false
	}
	if x < s.xSize() && y < s.ySize() && z < s.zSize() {
		return s.isFull(x, y, z)
	}
	return false
}

// DvsIsFullWideTransform mirrors DiscreteVoxelShape.isFullWide(AxisCycle, int, int, int).
func DvsIsFullWideTransform(s DiscreteVoxelShape, transform core.AxisCycle, x int, y int, z int) (ret bool) {
	return DvsIsFullWide(s,
		transform.CycleInt(x, y, z, core.AxisX),
		transform.CycleInt(x, y, z, core.AxisY),
		transform.CycleInt(x, y, z, core.AxisZ),
	)
}

// DvsIsFullTransform mirrors DiscreteVoxelShape.isFull(AxisCycle, int, int, int).
func DvsIsFullTransform(s DiscreteVoxelShape, transform core.AxisCycle, x int, y int, z int) (ret bool) {
	return s.isFull(
		transform.CycleInt(x, y, z, core.AxisX),
		transform.CycleInt(x, y, z, core.AxisY),
		transform.CycleInt(x, y, z, core.AxisZ),
	)
}

// DvsIsEmpty mirrors DiscreteVoxelShape.isEmpty().
func DvsIsEmpty(s DiscreteVoxelShape) (ret bool) {
	for _, axis := range core.DirectionAxisVALUES {
		if s.firstFull(axis) >= s.lastFull(axis) {
			return true
		}
	}
	return false
}

// DvsFirstFullABC mirrors DiscreteVoxelShape.firstFull(Direction.Axis, int, int).
func DvsFirstFullABC(s DiscreteVoxelShape, aAxis core.DirectionAxis, b int, c int) (ret int) {
	aSize := DvsGetSize(s, aAxis)
	if b >= 0 && c >= 0 {
		bAxis := core.AxisCycleFORWARD.CycleAxis(aAxis)
		cAxis := core.AxisCycleBACKWARD.CycleAxis(aAxis)
		if b < DvsGetSize(s, bAxis) && c < DvsGetSize(s, cAxis) {
			transform := core.AxisCycleBetween(core.AxisX, aAxis)
			for a := 0; a < aSize; a++ {
				if DvsIsFullTransform(s, transform, a, b, c) {
					return a
				}
			}
			return aSize
		}
		return aSize
	}
	return aSize
}

// DvsLastFullABC mirrors DiscreteVoxelShape.lastFull(Direction.Axis, int, int).
func DvsLastFullABC(s DiscreteVoxelShape, aAxis core.DirectionAxis, b int, c int) (ret int) {
	if b >= 0 && c >= 0 {
		bAxis := core.AxisCycleFORWARD.CycleAxis(aAxis)
		cAxis := core.AxisCycleBACKWARD.CycleAxis(aAxis)
		if b < DvsGetSize(s, bAxis) && c < DvsGetSize(s, cAxis) {
			aSize := DvsGetSize(s, aAxis)
			transform := core.AxisCycleBetween(core.AxisX, aAxis)
			for a := aSize - 1; a >= 0; a-- {
				if DvsIsFullTransform(s, transform, a, b, c) {
					return a + 1
				}
			}
			return 0
		}
		return 0
	}
	return 0
}

// IntLineConsumer mirrors DiscreteVoxelShape.IntLineConsumer.
type IntLineConsumer func(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int)

// IntFaceConsumer mirrors DiscreteVoxelShape.IntFaceConsumer.
type IntFaceConsumer func(direction core.Direction, x int, y int, z int)

// DvsForAllBoxes mirrors DiscreteVoxelShape.forAllBoxes(IntLineConsumer, boolean).
func DvsForAllBoxes(s DiscreteVoxelShape, consumer IntLineConsumer, mergeNeighbors bool) {
	bitSetForAllBoxes(s, consumer, mergeNeighbors)
}
