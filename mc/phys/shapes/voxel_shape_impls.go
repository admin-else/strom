package shapes

import (
	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/util"
)

// newCubeVoxelShape mirrors CubeVoxelShape(DiscreteVoxelShape).
func newCubeVoxelShape(shape DiscreteVoxelShape) (ret *VoxelShape) {
	return newVoxelShape(
		shape,
		func(axis core.DirectionAxis) []float64 {
			return cubePointRangeList(DvsGetSize(shape, axis))
		},
		func(axis core.DirectionAxis, coord float64) int {
			size := DvsGetSize(shape, axis)
			return util.FloorDouble(util.ClampDouble(coord*float64(size), -1.0, float64(size)))
		},
	)
}

// newArrayVoxelShape mirrors ArrayVoxelShape(DiscreteVoxelShape, DoubleList x3).
func newArrayVoxelShape(shape DiscreteVoxelShape, xs []float64, ys []float64, zs []float64) (ret *VoxelShape) {
	xSize := shape.xSize() + 1
	ySize := shape.ySize() + 1
	zSize := shape.zSize() + 1
	if xSize != len(xs) || ySize != len(ys) || zSize != len(zs) {
		panic("Lengths of point arrays must be consistent with the size of the VoxelShape.")
	}
	return newVoxelShape(
		shape,
		func(axis core.DirectionAxis) []float64 {
			switch axis {
			case core.AxisX:
				return xs
			case core.AxisY:
				return ys
			default:
				return zs
			}
		},
		nil,
	)
}

// sliceCoords mirrors SliceShape.SLICE_COORDS = new CubePointRange(1).
var sliceCoords = cubePointRangeList(1)

// newSliceShape mirrors SliceShape(VoxelShape, Direction.Axis, int).
func newSliceShape(delegate *VoxelShape, axis core.DirectionAxis, point int) (ret *VoxelShape) {
	d := delegate.shape
	shape := newSubShape(
		d,
		axis.ChooseInt(point, 0, 0),
		axis.ChooseInt(0, point, 0),
		axis.ChooseInt(0, 0, point),
		axis.ChooseInt(point+1, d.xSize(), d.xSize()),
		axis.ChooseInt(d.ySize(), point+1, d.ySize()),
		axis.ChooseInt(d.zSize(), d.zSize(), point+1),
	)
	return newVoxelShape(
		shape,
		func(a core.DirectionAxis) []float64 {
			if a == axis {
				return sliceCoords
			}
			return delegate.GetCoords(a)
		},
		nil,
	)
}
