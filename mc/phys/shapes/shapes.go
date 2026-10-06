package shapes

import (
	"math"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/phys"
)

// Shapes mirrors net.minecraft.world.phys.shapes.Shapes.
//
// Deferred: rotate/rotateHorizontal/rotateAll/rotateAttachFace and
// flipAxisIfNeeded (OctahedralGroup + AttachFace), which are block-model
// concerns. forAllEdges/forAllFaces are deferred with the debug paths.

// Shapes constants mirror the static fields on Shapes.
const (
	ShapesEPSILON     = 1.0e-7
	ShapesBIG_EPSILON = 1.0e-6
)

var (
	blockShape = makeBlockShape()
	infinity   = Box(math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1), math.Inf(1))
	emptyShape = newArrayVoxelShape(NewBitSetDiscreteVoxelShape(0, 0, 0), []float64{0.0}, []float64{0.0}, []float64{0.0})
)

func makeBlockShape() (ret *VoxelShape) {
	shape := NewBitSetDiscreteVoxelShape(1, 1, 1)
	shape.fill(0, 0, 0)
	return newCubeVoxelShape(shape)
}

// Empty mirrors Shapes.empty().
func Empty() (ret *VoxelShape) { return emptyShape }

// Block mirrors Shapes.block().
func Block() (ret *VoxelShape) { return blockShape }

// Infinity mirrors Shapes.INFINITY.
func Infinity() (ret *VoxelShape) { return infinity }

// Box mirrors Shapes.box(double x6).
func Box(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) (ret *VoxelShape) {
	if !(minX > maxX) && !(minY > maxY) && !(minZ > maxZ) {
		return Create(minX, minY, minZ, maxX, maxY, maxZ)
	}
	panic("The min values need to be smaller or equals to the max values")
}

// Create mirrors Shapes.create(double x6).
func Create(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) (ret *VoxelShape) {
	if !(maxX-minX < 1.0e-7) && !(maxY-minY < 1.0e-7) && !(maxZ-minZ < 1.0e-7) {
		xBits := FindBits(minX, maxX)
		yBits := FindBits(minY, maxY)
		zBits := FindBits(minZ, maxZ)
		if xBits < 0 || yBits < 0 || zBits < 0 {
			return newArrayVoxelShape(
				blockShape.shape,
				[]float64{minX, maxX},
				[]float64{minY, maxY},
				[]float64{minZ, maxZ},
			)
		}
		if xBits == 0 && yBits == 0 && zBits == 0 {
			return blockShape
		}
		xSize := 1 << xBits
		ySize := 1 << yBits
		zSize := 1 << zBits
		voxelShape := WithFilledBounds(
			xSize, ySize, zSize,
			int(math.Floor(minX*float64(xSize)+0.5)),
			int(math.Floor(minY*float64(ySize)+0.5)),
			int(math.Floor(minZ*float64(zSize)+0.5)),
			int(math.Floor(maxX*float64(xSize)+0.5)),
			int(math.Floor(maxY*float64(ySize)+0.5)),
			int(math.Floor(maxZ*float64(zSize)+0.5)),
		)
		return newCubeVoxelShape(voxelShape)
	}
	return emptyShape
}

// CreateFromAABB mirrors Shapes.create(AABB).
func CreateFromAABB(aabb phys.AABB) (ret *VoxelShape) {
	return Create(aabb.MinX, aabb.MinY, aabb.MinZ, aabb.MaxX, aabb.MaxY, aabb.MaxZ)
}

// FindBits mirrors Shapes.findBits(double, double). Exported for tests.
func FindBits(min float64, max float64) (ret int) {
	if !(min < -1.0e-7) && !(max > 1.0000001) {
		for bitsVal := 0; bitsVal <= 3; bitsVal++ {
			intervals := 1 << bitsVal
			shMin := min * float64(intervals)
			shMax := max * float64(intervals)
			foundMin := math.Abs(shMin-math.Round(shMin)) < 1.0e-7*float64(intervals)
			foundMax := math.Abs(shMax-math.Round(shMax)) < 1.0e-7*float64(intervals)
			if foundMin && foundMax {
				return bitsVal
			}
		}
		return -1
	}
	return -1
}

// Lcm mirrors Shapes.lcm(int, int).
func Lcm(first int, second int) (ret int64) { return lcmInt(first, second) }

// Or mirrors Shapes.or(VoxelShape, VoxelShape...).
func Or(first *VoxelShape, tail ...*VoxelShape) (ret *VoxelShape) {
	result := first
	for _, shape := range tail {
		result = Join(result, shape, BooleanOpOR)
	}
	return result
}

// Join mirrors Shapes.join(VoxelShape, VoxelShape, BooleanOp).
func Join(first *VoxelShape, second *VoxelShape, op BooleanOp) (ret *VoxelShape) {
	return JoinUnoptimized(first, second, op).Optimize()
}

// JoinUnoptimized mirrors Shapes.joinUnoptimized(VoxelShape, VoxelShape, BooleanOp).
func JoinUnoptimized(first *VoxelShape, second *VoxelShape, op BooleanOp) (ret *VoxelShape) {
	if op(false, false) {
		panic("BooleanOp must not return true for (false, false)")
	}
	if first == second {
		if op(true, true) {
			return first
		}
		return emptyShape
	}
	firstOnlyMatters := op(true, false)
	secondOnlyMatters := op(false, true)
	if first.IsEmpty() {
		if secondOnlyMatters {
			return second
		}
		return emptyShape
	}
	if second.IsEmpty() {
		if firstOnlyMatters {
			return first
		}
		return emptyShape
	}
	xMerger := createIndexMerger(1, first.GetCoords(core.AxisX), second.GetCoords(core.AxisX), firstOnlyMatters, secondOnlyMatters)
	yMerger := createIndexMerger(xMerger.Size()-1, first.GetCoords(core.AxisY), second.GetCoords(core.AxisY), firstOnlyMatters, secondOnlyMatters)
	zMerger := createIndexMerger((xMerger.Size()-1)*(yMerger.Size()-1), first.GetCoords(core.AxisZ), second.GetCoords(core.AxisZ), firstOnlyMatters, secondOnlyMatters)
	voxelShape := BitSetJoin(first.shape, second.shape, xMerger, yMerger, zMerger, op)
	if isDiscreteCubeMerger(xMerger) && isDiscreteCubeMerger(yMerger) && isDiscreteCubeMerger(zMerger) {
		return newCubeVoxelShape(voxelShape)
	}
	return newArrayVoxelShape(voxelShape, xMerger.GetList(), yMerger.GetList(), zMerger.GetList())
}

func isDiscreteCubeMerger(merger IndexMerger) (ret bool) {
	_, ok := merger.(DiscreteCubeMerger)
	return ok
}

// JoinIsNotEmpty mirrors Shapes.joinIsNotEmpty(VoxelShape, VoxelShape, BooleanOp).
func JoinIsNotEmpty(first *VoxelShape, second *VoxelShape, op BooleanOp) (ret bool) {
	if op(false, false) {
		panic("BooleanOp must not return true for (false, false)")
	}
	firstEmpty := first.IsEmpty()
	secondEmpty := second.IsEmpty()
	if !firstEmpty && !secondEmpty {
		if first == second {
			return op(true, true)
		}
		firstOnlyMatters := op(true, false)
		secondOnlyMatters := op(false, true)
		for _, axis := range core.DirectionAxisVALUES {
			if first.Max(axis) < second.Min(axis)-1.0e-7 {
				return firstOnlyMatters || secondOnlyMatters
			}
			if second.Max(axis) < first.Min(axis)-1.0e-7 {
				return firstOnlyMatters || secondOnlyMatters
			}
		}
		xMerger := createIndexMerger(1, first.GetCoords(core.AxisX), second.GetCoords(core.AxisX), firstOnlyMatters, secondOnlyMatters)
		yMerger := createIndexMerger(xMerger.Size()-1, first.GetCoords(core.AxisY), second.GetCoords(core.AxisY), firstOnlyMatters, secondOnlyMatters)
		zMerger := createIndexMerger((xMerger.Size()-1)*(yMerger.Size()-1), first.GetCoords(core.AxisZ), second.GetCoords(core.AxisZ), firstOnlyMatters, secondOnlyMatters)
		return joinIsNotEmptyMergers(xMerger, yMerger, zMerger, first.shape, second.shape, op)
	}
	return op(!firstEmpty, !secondEmpty)
}

func joinIsNotEmptyMergers(xMerger IndexMerger, yMerger IndexMerger, zMerger IndexMerger, first DiscreteVoxelShape, second DiscreteVoxelShape, op BooleanOp) (ret bool) {
	inverted := xMerger.ForMergedIndexes(func(x1 int, x2 int, xr int) bool {
		return yMerger.ForMergedIndexes(func(y1 int, y2 int, yr int) bool {
			return zMerger.ForMergedIndexes(func(z1 int, z2 int, zr int) bool {
				return !op(DvsIsFullWide(first, x1, y1, z1), DvsIsFullWide(second, x2, y2, z2))
			})
		})
	})
	return !inverted
}

// Collide mirrors Shapes.collide(Direction.Axis, AABB, Iterable<VoxelShape>, double).
func Collide(axis core.DirectionAxis, moving phys.AABB, shapeList []*VoxelShape, distance float64) (ret float64) {
	for _, shape := range shapeList {
		if math.Abs(distance) < 1.0e-7 {
			return 0.0
		}
		distance = shape.Collide(axis, moving, distance)
	}
	return distance
}

// BlockOccludes mirrors Shapes.blockOccludes(VoxelShape, VoxelShape, Direction).
func BlockOccludes(shape *VoxelShape, occluder *VoxelShape, direction core.Direction) (ret bool) {
	if shape == Block() && occluder == Block() {
		return true
	}
	if occluder.IsEmpty() {
		return false
	}
	axis := direction.GetAxis()
	sign := direction.GetAxisDirection()
	first := occluder
	second := shape
	op := BooleanOpONLY_SECOND
	if sign == core.DirectionAxisDirectionPOSITIVE {
		first = shape
		second = occluder
		op = BooleanOpONLY_FIRST
	}
	return fuzzyEquals(first.Max(axis), 1.0, 1.0e-7) &&
		fuzzyEquals(second.Min(axis), 0.0, 1.0e-7) &&
		!JoinIsNotEmpty(newSliceShape(first, axis, DvsGetSize(first.shape, axis)-1), newSliceShape(second, axis, 0), op)
}

// MergedFaceOccludes mirrors Shapes.mergedFaceOccludes(VoxelShape, VoxelShape, Direction).
func MergedFaceOccludes(shape *VoxelShape, occluder *VoxelShape, direction core.Direction) (ret bool) {
	if shape != Block() && occluder != Block() {
		axis := direction.GetAxis()
		sign := direction.GetAxisDirection()
		first := occluder
		second := shape
		if sign == core.DirectionAxisDirectionPOSITIVE {
			first = shape
			second = occluder
		}
		if !fuzzyEquals(first.Max(axis), 1.0, 1.0e-7) {
			first = emptyShape
		}
		if !fuzzyEquals(second.Min(axis), 0.0, 1.0e-7) {
			second = emptyShape
		}
		return !JoinIsNotEmpty(
			Block(),
			JoinUnoptimized(newSliceShape(first, axis, DvsGetSize(first.shape, axis)-1), newSliceShape(second, axis, 0), BooleanOpOR),
			BooleanOpONLY_FIRST,
		)
	}
	return true
}

// FaceShapeOccludes mirrors Shapes.faceShapeOccludes(VoxelShape, VoxelShape).
func FaceShapeOccludes(shape *VoxelShape, occluder *VoxelShape) (ret bool) {
	if shape == Block() || occluder == Block() {
		return true
	}
	if shape.IsEmpty() && occluder.IsEmpty() {
		return false
	}
	return !JoinIsNotEmpty(Block(), JoinUnoptimized(shape, occluder, BooleanOpOR), BooleanOpONLY_FIRST)
}

func createIndexMerger(cost int, first []float64, second []float64, firstOnlyMatters bool, secondOnlyMatters bool) (ret IndexMerger) {
	firstSize := len(first) - 1
	secondSize := len(second) - 1
	if p1, ok1 := isCubePointRange(first); ok1 {
		if p2, ok2 := isCubePointRange(second); ok2 {
			size := lcmInt(p1, p2)
			if int64(cost)*size <= 256 {
				return NewDiscreteCubeMerger(p1, p2)
			}
		}
	}
	if first[firstSize] < second[0]-1.0e-7 {
		return NewNonOverlappingMerger(first, second, false)
	}
	if second[secondSize] < first[0]-1.0e-7 {
		return NewNonOverlappingMerger(second, first, true)
	}
	if firstSize == secondSize && coordsListEqual(first, second) {
		return NewIdenticalMerger(first)
	}
	return NewIndirectMerger(first, second, firstOnlyMatters, secondOnlyMatters)
}

// Equal mirrors Shapes.equal(VoxelShape, VoxelShape).
func Equal(first *VoxelShape, second *VoxelShape) (ret bool) {
	return !JoinIsNotEmpty(first, second, BooleanOpNOT_SAME)
}

// coordsListEqual mirrors Objects.equals(DoubleList, DoubleList), whose element
// equality is Double.equals (NaN equals NaN, -0.0 is not 0.0).
func coordsListEqual(a []float64, b []float64) (ret bool) {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !coordsEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func coordsEqual(a float64, b float64) (ret bool) {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Float64bits(a) == math.Float64bits(b)
}
