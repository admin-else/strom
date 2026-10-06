package phys

import (
	"math"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/util"
)

// AABB mirrors net.minecraft.world.phys.AABB.
//
// Deferred with the units that need them: the BlockPos constructors, move,
// intersects, unitCubeFromLowerCorner and encapsulatingFullBlocks
// (`net.minecraft.core.BlockPos`), of(BoundingBox)
// (`net.minecraft.world.level.levelgen.structure.BoundingBox`), min/max
// (Direction.Axis), the clip/getDirection/clipPoint ray-cast family and
// collidedAlongVector (`net.minecraft.core.Direction`), the JOML-backed Builder
// and move(Vector3f) overload, and Java's toString.
type AABB struct {
	MinX float64
	MinY float64
	MinZ float64
	MaxX float64
	MaxY float64
	MaxZ float64
}

// NewAABB mirrors AABB(double, double, double, double, double, double). The
// constructor normalises the corners with min/max, exactly like Java.
func NewAABB(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) (ret AABB) {
	return AABB{
		MinX: math.Min(minX, maxX),
		MinY: math.Min(minY, maxY),
		MinZ: math.Min(minZ, maxZ),
		MaxX: math.Max(minX, maxX),
		MaxY: math.Max(minY, maxY),
		MaxZ: math.Max(minZ, maxZ),
	}
}

// NewAABBFromVec3 mirrors AABB(Vec3, Vec3).
func NewAABBFromVec3(begin Vec3, end Vec3) (ret AABB) {
	return NewAABB(begin.X, begin.Y, begin.Z, end.X, end.Y, end.Z)
}

// NewAABBFromBlockPos mirrors AABB(BlockPos).
func NewAABBFromBlockPos(pos core.BlockPos) (ret AABB) {
	return NewAABB(float64(pos.X), float64(pos.Y), float64(pos.Z), float64(pos.X+1), float64(pos.Y+1), float64(pos.Z+1))
}

// UnitCubeFromLowerCorner mirrors AABB.unitCubeFromLowerCorner(Vec3).
func UnitCubeFromLowerCorner(pos Vec3) (ret AABB) {
	return NewAABB(pos.X, pos.Y, pos.Z, pos.X+1.0, pos.Y+1.0, pos.Z+1.0)
}

// EncapsulatingFullBlocks mirrors AABB.encapsulatingFullBlocks(BlockPos, BlockPos).
func EncapsulatingFullBlocks(pos0 core.BlockPos, pos1 core.BlockPos) (ret AABB) {
	return NewAABB(
		float64(minBlock(pos0.X, pos1.X)),
		float64(minBlock(pos0.Y, pos1.Y)),
		float64(minBlock(pos0.Z, pos1.Z)),
		float64(maxBlock(pos0.X, pos1.X)+1),
		float64(maxBlock(pos0.Y, pos1.Y)+1),
		float64(maxBlock(pos0.Z, pos1.Z)+1),
	)
}

// Min mirrors AABB.min(Direction.Axis).
func (b AABB) Min(axis core.DirectionAxis) (ret float64) {
	return axis.ChooseDouble(b.MinX, b.MinY, b.MinZ)
}

// Max mirrors AABB.max(Direction.Axis).
func (b AABB) Max(axis core.DirectionAxis) (ret float64) {
	return axis.ChooseDouble(b.MaxX, b.MaxY, b.MaxZ)
}

// SetMinX mirrors AABB.setMinX(double).
func (b AABB) SetMinX(minX float64) (ret AABB) {
	return NewAABB(minX, b.MinY, b.MinZ, b.MaxX, b.MaxY, b.MaxZ)
}

// SetMinY mirrors AABB.setMinY(double).
func (b AABB) SetMinY(minY float64) (ret AABB) {
	return NewAABB(b.MinX, minY, b.MinZ, b.MaxX, b.MaxY, b.MaxZ)
}

// SetMinZ mirrors AABB.setMinZ(double).
func (b AABB) SetMinZ(minZ float64) (ret AABB) {
	return NewAABB(b.MinX, b.MinY, minZ, b.MaxX, b.MaxY, b.MaxZ)
}

// SetMaxX mirrors AABB.setMaxX(double).
func (b AABB) SetMaxX(maxX float64) (ret AABB) {
	return NewAABB(b.MinX, b.MinY, b.MinZ, maxX, b.MaxY, b.MaxZ)
}

// SetMaxY mirrors AABB.setMaxY(double).
func (b AABB) SetMaxY(maxY float64) (ret AABB) {
	return NewAABB(b.MinX, b.MinY, b.MinZ, b.MaxX, maxY, b.MaxZ)
}

// SetMaxZ mirrors AABB.setMaxZ(double).
func (b AABB) SetMaxZ(maxZ float64) (ret AABB) {
	return NewAABB(b.MinX, b.MinY, b.MinZ, b.MaxX, b.MaxY, maxZ)
}

// Contract mirrors AABB.contract(double, double, double).
func (b AABB) Contract(xa float64, ya float64, za float64) (ret AABB) {
	minX, minY, minZ := b.MinX, b.MinY, b.MinZ
	maxX, maxY, maxZ := b.MaxX, b.MaxY, b.MaxZ
	if xa < 0.0 {
		minX -= xa
	} else if xa > 0.0 {
		maxX -= xa
	}
	if ya < 0.0 {
		minY -= ya
	} else if ya > 0.0 {
		maxY -= ya
	}
	if za < 0.0 {
		minZ -= za
	} else if za > 0.0 {
		maxZ -= za
	}
	return NewAABB(minX, minY, minZ, maxX, maxY, maxZ)
}

// ExpandTowards mirrors AABB.expandTowards(Vec3).
func (b AABB) ExpandTowards(delta Vec3) (ret AABB) {
	return b.ExpandTowardsXYZ(delta.X, delta.Y, delta.Z)
}

// ExpandTowardsXYZ mirrors AABB.expandTowards(double, double, double).
func (b AABB) ExpandTowardsXYZ(xa float64, ya float64, za float64) (ret AABB) {
	minX, minY, minZ := b.MinX, b.MinY, b.MinZ
	maxX, maxY, maxZ := b.MaxX, b.MaxY, b.MaxZ
	if xa < 0.0 {
		minX += xa
	} else if xa > 0.0 {
		maxX += xa
	}
	if ya < 0.0 {
		minY += ya
	} else if ya > 0.0 {
		maxY += ya
	}
	if za < 0.0 {
		minZ += za
	} else if za > 0.0 {
		maxZ += za
	}
	return NewAABB(minX, minY, minZ, maxX, maxY, maxZ)
}

// InflateXYZ mirrors AABB.inflate(double, double, double).
func (b AABB) InflateXYZ(xAdd float64, yAdd float64, zAdd float64) (ret AABB) {
	return NewAABB(b.MinX-xAdd, b.MinY-yAdd, b.MinZ-zAdd, b.MaxX+xAdd, b.MaxY+yAdd, b.MaxZ+zAdd)
}

// Inflate mirrors AABB.inflate(double).
func (b AABB) Inflate(amount float64) (ret AABB) {
	return b.InflateXYZ(amount, amount, amount)
}

// Intersect mirrors AABB.intersect(AABB).
func (b AABB) Intersect(other AABB) (ret AABB) {
	return NewAABB(
		math.Max(b.MinX, other.MinX),
		math.Max(b.MinY, other.MinY),
		math.Max(b.MinZ, other.MinZ),
		math.Min(b.MaxX, other.MaxX),
		math.Min(b.MaxY, other.MaxY),
		math.Min(b.MaxZ, other.MaxZ),
	)
}

// Minmax mirrors AABB.minmax(AABB).
func (b AABB) Minmax(other AABB) (ret AABB) {
	return NewAABB(
		math.Min(b.MinX, other.MinX),
		math.Min(b.MinY, other.MinY),
		math.Min(b.MinZ, other.MinZ),
		math.Max(b.MaxX, other.MaxX),
		math.Max(b.MaxY, other.MaxY),
		math.Max(b.MaxZ, other.MaxZ),
	)
}

// MoveXYZ mirrors AABB.move(double, double, double).
func (b AABB) MoveXYZ(xa float64, ya float64, za float64) (ret AABB) {
	return NewAABB(b.MinX+xa, b.MinY+ya, b.MinZ+za, b.MaxX+xa, b.MaxY+ya, b.MaxZ+za)
}

// MoveVec3 mirrors AABB.move(Vec3).
func (b AABB) MoveVec3(pos Vec3) (ret AABB) { return b.MoveXYZ(pos.X, pos.Y, pos.Z) }

// MoveBlockPos mirrors AABB.move(BlockPos).
func (b AABB) MoveBlockPos(pos core.BlockPos) (ret AABB) {
	return NewAABB(
		b.MinX+float64(pos.X), b.MinY+float64(pos.Y), b.MinZ+float64(pos.Z),
		b.MaxX+float64(pos.X), b.MaxY+float64(pos.Y), b.MaxZ+float64(pos.Z),
	)
}

// Intersects mirrors AABB.intersects(AABB).
func (b AABB) Intersects(aabb AABB) (ret bool) {
	return b.IntersectsXYZ(aabb.MinX, aabb.MinY, aabb.MinZ, aabb.MaxX, aabb.MaxY, aabb.MaxZ)
}

// IntersectsXYZ mirrors AABB.intersects(double, double, double, double, double, double).
func (b AABB) IntersectsXYZ(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) (ret bool) {
	return b.MinX < maxX && b.MaxX > minX && b.MinY < maxY && b.MaxY > minY && b.MinZ < maxZ && b.MaxZ > minZ
}

// IntersectsVec3 mirrors AABB.intersects(Vec3, Vec3).
func (b AABB) IntersectsVec3(min Vec3, max Vec3) (ret bool) {
	return b.IntersectsXYZ(
		math.Min(min.X, max.X), math.Min(min.Y, max.Y), math.Min(min.Z, max.Z),
		math.Max(min.X, max.X), math.Max(min.Y, max.Y), math.Max(min.Z, max.Z),
	)
}

// IntersectsBlockPos mirrors AABB.intersects(BlockPos).
func (b AABB) IntersectsBlockPos(pos core.BlockPos) (ret bool) {
	return b.IntersectsXYZ(
		float64(pos.X), float64(pos.Y), float64(pos.Z),
		float64(pos.X+1), float64(pos.Y+1), float64(pos.Z+1),
	)
}

// Contains mirrors AABB.contains(Vec3).
func (b AABB) Contains(vec Vec3) (ret bool) {
	return b.ContainsXYZ(vec.X, vec.Y, vec.Z)
}

// ContainsXYZ mirrors AABB.contains(double, double, double).
func (b AABB) ContainsXYZ(x float64, y float64, z float64) (ret bool) {
	return x >= b.MinX && x < b.MaxX && y >= b.MinY && y < b.MaxY && z >= b.MinZ && z < b.MaxZ
}

// GetSize mirrors AABB.getSize().
func (b AABB) GetSize() (ret float64) {
	return (b.GetXsize() + b.GetYsize() + b.GetZsize()) / 3.0
}

// GetXsize mirrors AABB.getXsize().
func (b AABB) GetXsize() (ret float64) { return b.MaxX - b.MinX }

// GetYsize mirrors AABB.getYsize().
func (b AABB) GetYsize() (ret float64) { return b.MaxY - b.MinY }

// GetZsize mirrors AABB.getZsize().
func (b AABB) GetZsize() (ret float64) { return b.MaxZ - b.MinZ }

// DeflateXYZ mirrors AABB.deflate(double, double, double).
func (b AABB) DeflateXYZ(xSubtract float64, ySubtract float64, zSubtract float64) (ret AABB) {
	return b.InflateXYZ(-xSubtract, -ySubtract, -zSubtract)
}

// Deflate mirrors AABB.deflate(double).
func (b AABB) Deflate(amount float64) (ret AABB) { return b.Inflate(-amount) }

// DistanceToSqrVec3 mirrors AABB.distanceToSqr(Vec3).
func (b AABB) DistanceToSqrVec3(point Vec3) (ret float64) {
	dx := math.Max(math.Max(b.MinX-point.X, point.X-b.MaxX), 0.0)
	dy := math.Max(math.Max(b.MinY-point.Y, point.Y-b.MaxY), 0.0)
	dz := math.Max(math.Max(b.MinZ-point.Z, point.Z-b.MaxZ), 0.0)
	return util.LengthSquared3Double(dx, dy, dz)
}

// DistanceToSqrAABB mirrors AABB.distanceToSqr(AABB).
func (b AABB) DistanceToSqrAABB(other AABB) (ret float64) {
	dx := math.Max(math.Max(b.MinX-other.MaxX, other.MinX-b.MaxX), 0.0)
	dy := math.Max(math.Max(b.MinY-other.MaxY, other.MinY-b.MaxY), 0.0)
	dz := math.Max(math.Max(b.MinZ-other.MaxZ, other.MinZ-b.MaxZ), 0.0)
	return util.LengthSquared3Double(dx, dy, dz)
}

// HasNaN mirrors AABB.hasNaN().
func (b AABB) HasNaN() (ret bool) {
	return math.IsNaN(b.MinX) || math.IsNaN(b.MinY) || math.IsNaN(b.MinZ) ||
		math.IsNaN(b.MaxX) || math.IsNaN(b.MaxY) || math.IsNaN(b.MaxZ)
}

// GetCenter mirrors AABB.getCenter().
func (b AABB) GetCenter() (ret Vec3) {
	return NewVec3(
		util.LerpDouble(0.5, b.MinX, b.MaxX),
		util.LerpDouble(0.5, b.MinY, b.MaxY),
		util.LerpDouble(0.5, b.MinZ, b.MaxZ),
	)
}

// GetBottomCenter mirrors AABB.getBottomCenter().
func (b AABB) GetBottomCenter() (ret Vec3) {
	return NewVec3(
		util.LerpDouble(0.5, b.MinX, b.MaxX),
		b.MinY,
		util.LerpDouble(0.5, b.MinZ, b.MaxZ),
	)
}

// GetMinPosition mirrors AABB.getMinPosition().
func (b AABB) GetMinPosition() (ret Vec3) { return NewVec3(b.MinX, b.MinY, b.MinZ) }

// GetMaxPosition mirrors AABB.getMaxPosition().
func (b AABB) GetMaxPosition() (ret Vec3) { return NewVec3(b.MaxX, b.MaxY, b.MaxZ) }

// OfSize mirrors AABB.ofSize(Vec3, double, double, double).
func OfSize(center Vec3, sizeX float64, sizeY float64, sizeZ float64) (ret AABB) {
	return NewAABB(
		center.X-sizeX/2.0, center.Y-sizeY/2.0, center.Z-sizeZ/2.0,
		center.X+sizeX/2.0, center.Y+sizeY/2.0, center.Z+sizeZ/2.0,
	)
}

// NextDeflated mirrors AABB.nextDeflated().
func (b AABB) NextDeflated() (ret AABB) {
	return NewAABB(
		math.Nextafter(b.MinX, math.Inf(1)),
		math.Nextafter(b.MinY, math.Inf(1)),
		math.Nextafter(b.MinZ, math.Inf(1)),
		math.Nextafter(b.MaxX, math.Inf(-1)),
		math.Nextafter(b.MaxY, math.Inf(-1)),
		math.Nextafter(b.MaxZ, math.Inf(-1)),
	)
}

// Clip mirrors AABB.clip(Vec3, Vec3). ok is false when the segment misses.
func (b AABB) Clip(from Vec3, to Vec3) (ret Vec3, ok bool) {
	return AABBClip(b.MinX, b.MinY, b.MinZ, b.MaxX, b.MaxY, b.MaxZ, from, to)
}

// AABBClip mirrors AABB.clip(double x6, Vec3, Vec3).
func AABBClip(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64, from Vec3, to Vec3) (ret Vec3, ok bool) {
	scaleReference := 1.0
	dx := to.X - from.X
	dy := to.Y - from.Y
	dz := to.Z - from.Z
	_, ok = getDirectionXYZ(minX, minY, minZ, maxX, maxY, maxZ, from, &scaleReference, core.DirectionDOWN, false, dx, dy, dz)
	if !ok {
		return Vec3{}, false
	}
	scale := scaleReference
	return from.AddXYZ(scale*dx, scale*dy, scale*dz), true
}

// ClipAABBs mirrors AABB.clip(Iterable<AABB>, Vec3, Vec3, BlockPos). ok is false
// when the Java method returns null.
func ClipAABBs(aabbs []AABB, from Vec3, to Vec3, pos core.BlockPos) (ret BlockHitResult, ok bool) {
	scaleReference := 1.0
	var direction core.Direction
	directionOK := false
	dx := to.X - from.X
	dy := to.Y - from.Y
	dz := to.Z - from.Z

	for _, aabb := range aabbs {
		moved := aabb.MoveBlockPos(pos)
		direction, directionOK = getDirectionXYZ(
			moved.MinX, moved.MinY, moved.MinZ, moved.MaxX, moved.MaxY, moved.MaxZ,
			from, &scaleReference, direction, directionOK, dx, dy, dz,
		)
	}

	if !directionOK {
		return BlockHitResult{}, false
	}

	scale := scaleReference
	return NewBlockHitResult(from.AddXYZ(scale*dx, scale*dy, scale*dz), direction, pos, false), true
}

func getDirectionXYZ(
	minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64,
	from Vec3, scaleReference *float64, direction core.Direction, directionOK bool,
	dx float64, dy float64, dz float64,
) (ret core.Direction, retOK bool) {
	ret, retOK = direction, directionOK
	if dx > 1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dx, dy, dz, minX, minY, maxY, minZ, maxZ, core.DirectionWEST, from.X, from.Y, from.Z)
	} else if dx < -1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dx, dy, dz, maxX, minY, maxY, minZ, maxZ, core.DirectionEAST, from.X, from.Y, from.Z)
	}
	if dy > 1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dy, dz, dx, minY, minZ, maxZ, minX, maxX, core.DirectionDOWN, from.Y, from.Z, from.X)
	} else if dy < -1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dy, dz, dx, maxY, minZ, maxZ, minX, maxX, core.DirectionUP, from.Y, from.Z, from.X)
	}
	if dz > 1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dz, dx, dy, minZ, minX, maxX, minY, maxY, core.DirectionNORTH, from.Z, from.X, from.Y)
	} else if dz < -1.0e-7 {
		ret, retOK = clipPoint(scaleReference, ret, retOK, dz, dx, dy, maxZ, minX, maxX, minY, maxY, core.DirectionSOUTH, from.Z, from.X, from.Y)
	}
	return ret, retOK
}

func clipPoint(
	scaleReference *float64, direction core.Direction, directionOK bool,
	da float64, db float64, dc float64,
	point float64, minB float64, maxB float64, minC float64, maxC float64,
	newDirection core.Direction, fromA float64, fromB float64, fromC float64,
) (ret core.Direction, retOK bool) {
	s := (point - fromA) / da
	pb := fromB + s*db
	pc := fromC + s*dc
	if 0.0 < s && s < *scaleReference && minB-1.0e-7 < pb && pb < maxB+1.0e-7 && minC-1.0e-7 < pc && pc < maxC+1.0e-7 {
		*scaleReference = s
		return newDirection, true
	}
	return direction, directionOK
}

func minBlock(a int, b int) (ret int) {
	if a < b {
		return a
	}
	return b
}

func maxBlock(a int, b int) (ret int) {
	if a > b {
		return a
	}
	return b
}

// Equal mirrors AABB.equals(Object) for two AABBs. Java uses
// Double.compare, so -0.0 and 0.0 are not equal and any two NaNs are.
func (b AABB) Equal(other AABB) (ret bool) {
	return doubleCompareEqual(b.MinX, other.MinX) && doubleCompareEqual(b.MinY, other.MinY) && doubleCompareEqual(b.MinZ, other.MinZ) &&
		doubleCompareEqual(b.MaxX, other.MaxX) && doubleCompareEqual(b.MaxY, other.MaxY) && doubleCompareEqual(b.MaxZ, other.MaxZ)
}

func doubleCompareEqual(a float64, b float64) (ret bool) {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Float64bits(a) == math.Float64bits(b)
}
