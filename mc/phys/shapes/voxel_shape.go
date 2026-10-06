package shapes

import (
	"math"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/util"
)

// DoubleLineConsumer mirrors Shapes.DoubleLineConsumer.
type DoubleLineConsumer func(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64)

// VoxelShape mirrors net.minecraft.world.phys.shapes.VoxelShape. Java's abstract
// class with getCoords/findIndex overrides is modelled as a struct holding those
// two operations as function fields.
type VoxelShape struct {
	shape          DiscreteVoxelShape
	getCoordsFn    func(axis core.DirectionAxis) []float64
	findIndexFn    func(axis core.DirectionAxis, coord float64) int
	faces          [6]*VoxelShape
	facesAllocated bool
}

func newVoxelShape(shape DiscreteVoxelShape, getCoords func(core.DirectionAxis) []float64, findIndex func(core.DirectionAxis, float64) int) (ret *VoxelShape) {
	return &VoxelShape{shape: shape, getCoordsFn: getCoords, findIndexFn: findIndex}
}

// Min mirrors VoxelShape.min(Direction.Axis).
func (v *VoxelShape) Min(axis core.DirectionAxis) (ret float64) {
	i := v.shape.firstFull(axis)
	if i >= DvsGetSize(v.shape, axis) {
		return math.Inf(1)
	}
	return v.get(axis, i)
}

// Max mirrors VoxelShape.max(Direction.Axis).
func (v *VoxelShape) Max(axis core.DirectionAxis) (ret float64) {
	i := v.shape.lastFull(axis)
	if i <= 0 {
		return math.Inf(-1)
	}
	return v.get(axis, i)
}

// Bounds mirrors VoxelShape.bounds(). ok is false where Java throws.
func (v *VoxelShape) Bounds() (ret phys.AABB, ok bool) {
	if v.IsEmpty() {
		return phys.AABB{}, false
	}
	return phys.NewAABB(
		v.Min(core.AxisX), v.Min(core.AxisY), v.Min(core.AxisZ),
		v.Max(core.AxisX), v.Max(core.AxisY), v.Max(core.AxisZ),
	), true
}

// SingleEncompassing mirrors VoxelShape.singleEncompassing().
func (v *VoxelShape) SingleEncompassing() (ret *VoxelShape) {
	if v.IsEmpty() {
		return Empty()
	}
	return Box(v.Min(core.AxisX), v.Min(core.AxisY), v.Min(core.AxisZ), v.Max(core.AxisX), v.Max(core.AxisY), v.Max(core.AxisZ))
}

func (v *VoxelShape) get(axis core.DirectionAxis, i int) (ret float64) {
	return v.GetCoords(axis)[i]
}

// GetCoords mirrors VoxelShape.getCoords(Direction.Axis).
func (v *VoxelShape) GetCoords(axis core.DirectionAxis) (ret []float64) {
	return v.getCoordsFn(axis)
}

// IsEmpty mirrors VoxelShape.isEmpty().
func (v *VoxelShape) IsEmpty() (ret bool) { return DvsIsEmpty(v.shape) }

// Move mirrors VoxelShape.move(Vec3).
func (v *VoxelShape) Move(delta phys.Vec3) (ret *VoxelShape) {
	return v.MoveXYZ(delta.X, delta.Y, delta.Z)
}

// MoveVec3i mirrors VoxelShape.move(Vec3i).
func (v *VoxelShape) MoveVec3i(delta core.Vec3i) (ret *VoxelShape) {
	return v.MoveXYZ(float64(delta.GetX()), float64(delta.GetY()), float64(delta.GetZ()))
}

// MoveXYZ mirrors VoxelShape.move(double, double, double).
func (v *VoxelShape) MoveXYZ(dx float64, dy float64, dz float64) (ret *VoxelShape) {
	if v.IsEmpty() {
		return Empty()
	}
	return newArrayVoxelShape(
		v.shape,
		offsetList(v.GetCoords(core.AxisX), dx),
		offsetList(v.GetCoords(core.AxisY), dy),
		offsetList(v.GetCoords(core.AxisZ), dz),
	)
}

// Optimize mirrors VoxelShape.optimize().
func (v *VoxelShape) Optimize() (ret *VoxelShape) {
	result := Empty()
	v.ForAllBoxes(func(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) {
		result = JoinUnoptimized(result, Box(x1, y1, z1, x2, y2, z2), BooleanOpOR)
	})
	return result
}

// ForAllBoxes mirrors VoxelShape.forAllBoxes(Shapes.DoubleLineConsumer).
func (v *VoxelShape) ForAllBoxes(consumer DoubleLineConsumer) {
	xCoords := v.GetCoords(core.AxisX)
	yCoords := v.GetCoords(core.AxisY)
	zCoords := v.GetCoords(core.AxisZ)
	DvsForAllBoxes(v.shape, func(xi1 int, yi1 int, zi1 int, xi2 int, yi2 int, zi2 int) {
		consumer(xCoords[xi1], yCoords[yi1], zCoords[zi1], xCoords[xi2], yCoords[yi2], zCoords[zi2])
	}, true)
}

// ToAabbs mirrors VoxelShape.toAabbs().
func (v *VoxelShape) ToAabbs() (ret []phys.AABB) {
	v.ForAllBoxes(func(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) {
		ret = append(ret, phys.NewAABB(x1, y1, z1, x2, y2, z2))
	})
	return
}

// MinABC mirrors VoxelShape.min(Direction.Axis, double, double).
func (v *VoxelShape) MinABC(aAxis core.DirectionAxis, b float64, c float64) (ret float64) {
	bAxis := core.AxisCycleFORWARD.CycleAxis(aAxis)
	cAxis := core.AxisCycleBACKWARD.CycleAxis(aAxis)
	bi := v.findIndex(bAxis, b)
	ci := v.findIndex(cAxis, c)
	i := DvsFirstFullABC(v.shape, aAxis, bi, ci)
	if i >= DvsGetSize(v.shape, aAxis) {
		return math.Inf(1)
	}
	return v.get(aAxis, i)
}

// MaxABC mirrors VoxelShape.max(Direction.Axis, double, double).
func (v *VoxelShape) MaxABC(aAxis core.DirectionAxis, b float64, c float64) (ret float64) {
	bAxis := core.AxisCycleFORWARD.CycleAxis(aAxis)
	cAxis := core.AxisCycleBACKWARD.CycleAxis(aAxis)
	bi := v.findIndex(bAxis, b)
	ci := v.findIndex(cAxis, c)
	i := DvsLastFullABC(v.shape, aAxis, bi, ci)
	if i <= 0 {
		return math.Inf(-1)
	}
	return v.get(aAxis, i)
}

func (v *VoxelShape) findIndex(axis core.DirectionAxis, coord float64) (ret int) {
	if v.findIndexFn != nil {
		return v.findIndexFn(axis, coord)
	}
	return defaultFindIndex(v, axis, coord)
}

// Clip mirrors VoxelShape.clip(Vec3, Vec3, BlockPos). A nil result means the
// Java method returned null.
func (v *VoxelShape) Clip(from phys.Vec3, to phys.Vec3, pos core.BlockPos) (ret *phys.BlockHitResult) {
	if v.IsEmpty() {
		return nil
	}
	diff := to.Subtract(from)
	if diff.LengthSqr() < 1.0e-7 {
		return nil
	}
	testPoint := from.Add(diff.Scale(0.001))
	if DvsIsFullWide(v.shape,
		v.findIndex(core.AxisX, testPoint.X-float64(pos.GetX())),
		v.findIndex(core.AxisY, testPoint.Y-float64(pos.GetY())),
		v.findIndex(core.AxisZ, testPoint.Z-float64(pos.GetZ())),
	) {
		result := phys.NewBlockHitResult(
			testPoint,
			core.DirectionGetApproximateNearest(float32(diff.X), float32(diff.Y), float32(diff.Z)).GetOpposite(),
			pos,
			true,
		)
		return &result
	}
	hit, ok := phys.ClipAABBs(v.ToAabbs(), from, to, pos)
	if !ok {
		return nil
	}
	return &hit
}

// ClosestPointTo mirrors VoxelShape.closestPointTo(Vec3). ok is false where Java
// returns Optional.empty().
func (v *VoxelShape) ClosestPointTo(point phys.Vec3) (ret phys.Vec3, ok bool) {
	if v.IsEmpty() {
		return phys.Vec3{}, false
	}
	var closest phys.Vec3
	haveClosest := false
	v.ForAllBoxes(func(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) {
		x := util.ClampDouble(point.X, x1, x2)
		y := util.ClampDouble(point.Y, y1, y2)
		z := util.ClampDouble(point.Z, z1, z2)
		if !haveClosest || point.DistanceToSqrXYZ(x, y, z) < point.DistanceToSqr(closest) {
			closest = phys.NewVec3(x, y, z)
			haveClosest = true
		}
	})
	return closest, true
}

// GetFaceShape mirrors VoxelShape.getFaceShape(Direction).
func (v *VoxelShape) GetFaceShape(direction core.Direction) (ret *VoxelShape) {
	if v.IsEmpty() || v == Block() {
		return v
	}
	if !v.facesAllocated {
		v.facesAllocated = true
	} else if face := v.faces[direction]; face != nil {
		return face
	}
	face := v.calculateFace(direction)
	v.faces[direction] = face
	return face
}

func (v *VoxelShape) calculateFace(direction core.Direction) (ret *VoxelShape) {
	axis := direction.GetAxis()
	if v.isCubeLikeAlong(axis) {
		return v
	}
	sign := direction.GetAxisDirection()
	coord := 1.0e-7
	if sign == core.DirectionAxisDirectionPOSITIVE {
		coord = 0.9999999
	}
	index := v.findIndex(axis, coord)
	slice := newSliceShape(v, axis, index)
	if slice.IsEmpty() {
		return Empty()
	}
	if slice.isCubeLike() {
		return Block()
	}
	return slice
}

func (v *VoxelShape) isCubeLike() (ret bool) {
	for _, axis := range core.DirectionAxisVALUES {
		if !v.isCubeLikeAlong(axis) {
			return false
		}
	}
	return true
}

func (v *VoxelShape) isCubeLikeAlong(axis core.DirectionAxis) (ret bool) {
	coords := v.GetCoords(axis)
	return len(coords) == 2 && fuzzyEquals(coords[0], 0.0, 1.0e-7) && fuzzyEquals(coords[1], 1.0, 1.0e-7)
}

// Collide mirrors VoxelShape.collide(Direction.Axis, AABB, double).
func (v *VoxelShape) Collide(axis core.DirectionAxis, moving phys.AABB, distance float64) (ret float64) {
	return v.collideX(core.AxisCycleBetween(axis, core.AxisX), moving, distance)
}

func (v *VoxelShape) collideX(transform core.AxisCycle, moving phys.AABB, distance float64) (ret float64) {
	if v.IsEmpty() {
		return distance
	}
	if math.Abs(distance) < 1.0e-7 {
		return 0.0
	}
	inverse := transform.Inverse()
	aAxis := inverse.CycleAxis(core.AxisX)
	bAxis := inverse.CycleAxis(core.AxisY)
	cAxis := inverse.CycleAxis(core.AxisZ)
	maxA := moving.Max(aAxis)
	minA := moving.Min(aAxis)
	aMin := v.findIndex(aAxis, minA+1.0e-7)
	aMax := v.findIndex(aAxis, maxA-1.0e-7)
	bMin := maxInt(0, v.findIndex(bAxis, moving.Min(bAxis)+1.0e-7))
	bMax := minInt(DvsGetSize(v.shape, bAxis), v.findIndex(bAxis, moving.Max(bAxis)-1.0e-7)+1)
	cMin := maxInt(0, v.findIndex(cAxis, moving.Min(cAxis)+1.0e-7))
	cMax := minInt(DvsGetSize(v.shape, cAxis), v.findIndex(cAxis, moving.Max(cAxis)-1.0e-7)+1)
	aSize := DvsGetSize(v.shape, aAxis)
	if distance > 0.0 {
		for a := aMax + 1; a < aSize; a++ {
			for b := bMin; b < bMax; b++ {
				for c := cMin; c < cMax; c++ {
					if DvsIsFullWideTransform(v.shape, inverse, a, b, c) {
						newDistance := v.get(aAxis, a) - maxA
						if newDistance >= -1.0e-7 {
							distance = math.Min(distance, newDistance)
						}
						return distance
					}
				}
			}
		}
	} else if distance < 0.0 {
		for a := aMin - 1; a >= 0; a-- {
			for b := bMin; b < bMax; b++ {
				for c := cMin; c < cMax; c++ {
					if DvsIsFullWideTransform(v.shape, inverse, a, b, c) {
						newDistance := v.get(aAxis, a+1) - minA
						if newDistance <= 1.0e-7 {
							distance = math.Max(distance, newDistance)
						}
						return distance
					}
				}
			}
		}
	}
	return distance
}

func defaultFindIndex(v *VoxelShape, axis core.DirectionAxis, coord float64) (ret int) {
	return util.BinarySearch(0, DvsGetSize(v.shape, axis)+1, func(index int) bool {
		return coord < v.get(axis, index)
	}) - 1
}

func fuzzyEquals(a float64, b float64, tolerance float64) (ret bool) {
	return math.Abs(a-b) <= tolerance
}

func offsetList(list []float64, offset float64) (ret []float64) {
	ret = make([]float64, len(list))
	for i, value := range list {
		ret[i] = value + offset
	}
	return
}
