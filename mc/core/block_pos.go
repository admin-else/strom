package core

import "github.com/admin-else/strom/mc/util"

// BlockPos mirrors net.minecraft.core.BlockPos. Java extends Vec3i; Go keeps a
// separate struct with the same field order. The Mojang codecs, the Position
// overloads, MutableBlockPos, rotate(Rotation) and the traversal iterators
// (randomInCube/randomBetweenClosed/withinManhattan/spiralAround/
// breadthFirstTraversal/betweenCornersInDirection) are deferred until a
// consumer needs them.
type BlockPos struct {
	X int
	Y int
	Z int
}

// BlockPosZERO mirrors BlockPos.ZERO.
var BlockPosZERO = NewBlockPos(0, 0, 0)

// Packing constants mirror BlockPos's static fields.
var (
	packedHorizontalLength = 1 + util.Log2(util.SmallestEncompassingPowerOfTwo(30000000))
	packedYLength          = 64 - 2*packedHorizontalLength
	packedXMask            = (int64(1) << packedHorizontalLength) - 1
	packedYMask            = (int64(1) << packedYLength) - 1
	packedZMask            = packedXMask
	zOffset                = packedYLength
	xOffset                = packedYLength + packedHorizontalLength
	// MaxHorizontalCoordinate mirrors BlockPos.MAX_HORIZONTAL_COORDINATE.
	MaxHorizontalCoordinate = (1<<packedHorizontalLength)/2 - 1
)

// NewBlockPos mirrors BlockPos(int, int, int).
func NewBlockPos(x int, y int, z int) (ret BlockPos) { return BlockPos{X: x, Y: y, Z: z} }

// GetX mirrors BlockPos.getX().
func (p BlockPos) GetX() (ret int) { return p.X }

// GetY mirrors BlockPos.getY().
func (p BlockPos) GetY() (ret int) { return p.Y }

// GetZ mirrors BlockPos.getZ().
func (p BlockPos) GetZ() (ret int) { return p.Z }

// BlockPosGetX mirrors BlockPos.getX(long).
func BlockPosGetX(blockNode int64) (ret int) {
	return int(blockNode << (64 - xOffset - packedHorizontalLength) >> (64 - packedHorizontalLength))
}

// BlockPosGetY mirrors BlockPos.getY(long).
func BlockPosGetY(blockNode int64) (ret int) {
	return int(blockNode << (64 - packedYLength) >> (64 - packedYLength))
}

// BlockPosGetZ mirrors BlockPos.getZ(long).
func BlockPosGetZ(blockNode int64) (ret int) {
	return int(blockNode << (64 - zOffset - packedHorizontalLength) >> (64 - packedHorizontalLength))
}

// BlockPosOfPackaged mirrors BlockPos.of(long).
func BlockPosOfPackaged(blockNode int64) (ret BlockPos) {
	return NewBlockPos(BlockPosGetX(blockNode), BlockPosGetY(blockNode), BlockPosGetZ(blockNode))
}

// BlockPosContaining mirrors BlockPos.containing(double, double, double).
func BlockPosContaining(x float64, y float64, z float64) (ret BlockPos) {
	return NewBlockPos(util.FloorDouble(x), util.FloorDouble(y), util.FloorDouble(z))
}

// BlockPosMin mirrors BlockPos.min(BlockPos, BlockPos).
func BlockPosMin(a BlockPos, b BlockPos) (ret BlockPos) {
	return NewBlockPos(minInt(a.X, b.X), minInt(a.Y, b.Y), minInt(a.Z, b.Z))
}

// BlockPosMax mirrors BlockPos.max(BlockPos, BlockPos).
func BlockPosMax(a BlockPos, b BlockPos) (ret BlockPos) {
	return NewBlockPos(maxInt(a.X, b.X), maxInt(a.Y, b.Y), maxInt(a.Z, b.Z))
}

// AsLong mirrors BlockPos.asLong().
func (p BlockPos) AsLong() (ret int64) { return BlockPosAsLong(p.X, p.Y, p.Z) }

// BlockPosAsLong mirrors BlockPos.asLong(int, int, int).
func BlockPosAsLong(x int, y int, z int) (ret int64) {
	var node int64
	node |= (int64(x) & packedXMask) << xOffset
	node |= int64(y) & packedYMask
	return node | (int64(z)&packedZMask)<<zOffset
}

// BlockPosGetFlatIndex mirrors BlockPos.getFlatIndex(long).
func BlockPosGetFlatIndex(neighborBlockNode int64) (ret int64) { return neighborBlockNode & -16 }

// Offset mirrors BlockPos.offset(int, int, int).
func (p BlockPos) Offset(x int, y int, z int) (ret BlockPos) {
	if x == 0 && y == 0 && z == 0 {
		return p
	}
	return NewBlockPos(p.X+x, p.Y+y, p.Z+z)
}

// OffsetVec mirrors BlockPos.offset(Vec3i).
func (p BlockPos) OffsetVec(vec Vec3i) (ret BlockPos) { return p.Offset(vec.X, vec.Y, vec.Z) }

// Subtract mirrors BlockPos.subtract(Vec3i).
func (p BlockPos) Subtract(vec Vec3i) (ret BlockPos) { return p.Offset(-vec.X, -vec.Y, -vec.Z) }

// Multiply mirrors BlockPos.multiply(int).
func (p BlockPos) Multiply(scale int) (ret BlockPos) {
	if scale == 1 {
		return p
	}
	if scale == 0 {
		return BlockPosZERO
	}
	return NewBlockPos(p.X*scale, p.Y*scale, p.Z*scale)
}

// Above mirrors BlockPos.above().
func (p BlockPos) Above() (ret BlockPos) { return p.Relative(DirectionUP) }

// AboveSteps mirrors BlockPos.above(int).
func (p BlockPos) AboveSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionUP, steps) }

// Below mirrors BlockPos.below().
func (p BlockPos) Below() (ret BlockPos) { return p.Relative(DirectionDOWN) }

// BelowSteps mirrors BlockPos.below(int).
func (p BlockPos) BelowSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionDOWN, steps) }

// North mirrors BlockPos.north().
func (p BlockPos) North() (ret BlockPos) { return p.Relative(DirectionNORTH) }

// NorthSteps mirrors BlockPos.north(int).
func (p BlockPos) NorthSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionNORTH, steps) }

// South mirrors BlockPos.south().
func (p BlockPos) South() (ret BlockPos) { return p.Relative(DirectionSOUTH) }

// SouthSteps mirrors BlockPos.south(int).
func (p BlockPos) SouthSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionSOUTH, steps) }

// West mirrors BlockPos.west().
func (p BlockPos) West() (ret BlockPos) { return p.Relative(DirectionWEST) }

// WestSteps mirrors BlockPos.west(int).
func (p BlockPos) WestSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionWEST, steps) }

// East mirrors BlockPos.east().
func (p BlockPos) East() (ret BlockPos) { return p.Relative(DirectionEAST) }

// EastSteps mirrors BlockPos.east(int).
func (p BlockPos) EastSteps(steps int) (ret BlockPos) { return p.RelativeSteps(DirectionEAST, steps) }

// Relative mirrors BlockPos.relative(Direction).
func (p BlockPos) Relative(direction Direction) (ret BlockPos) {
	return NewBlockPos(p.X+direction.GetStepX(), p.Y+direction.GetStepY(), p.Z+direction.GetStepZ())
}

// RelativeSteps mirrors BlockPos.relative(Direction, int).
func (p BlockPos) RelativeSteps(direction Direction, steps int) (ret BlockPos) {
	if steps == 0 {
		return p
	}
	return NewBlockPos(
		p.X+direction.GetStepX()*steps,
		p.Y+direction.GetStepY()*steps,
		p.Z+direction.GetStepZ()*steps,
	)
}

// RelativeAxis mirrors BlockPos.relative(Direction.Axis, int).
func (p BlockPos) RelativeAxis(axis DirectionAxis, steps int) (ret BlockPos) {
	if steps == 0 {
		return p
	}
	xStep, yStep, zStep := 0, 0, 0
	switch axis {
	case AxisX:
		xStep = steps
	case AxisY:
		yStep = steps
	default:
		zStep = steps
	}
	return NewBlockPos(p.X+xStep, p.Y+yStep, p.Z+zStep)
}

// Cross mirrors BlockPos.cross(Vec3i).
func (p BlockPos) Cross(upVector Vec3i) (ret BlockPos) {
	return NewBlockPos(
		p.Y*upVector.Z-p.Z*upVector.Y,
		p.Z*upVector.X-p.X*upVector.Z,
		p.X*upVector.Y-p.Y*upVector.X,
	)
}

// AtY mirrors BlockPos.atY(int).
func (p BlockPos) AtY(y int) (ret BlockPos) { return NewBlockPos(p.X, y, p.Z) }

// Immutable mirrors BlockPos.immutable().
func (p BlockPos) Immutable() (ret BlockPos) { return p }

// AsVec3i mirrors the upcast BlockPos -> Vec3i.
func (p BlockPos) AsVec3i() (ret Vec3i) { return NewVec3i(p.X, p.Y, p.Z) }

// Equal mirrors BlockPos.equals(Object).
func (p BlockPos) Equal(other BlockPos) (ret bool) {
	return p.X == other.X && p.Y == other.Y && p.Z == other.Z
}

// CompareTo mirrors Vec3i.compareTo inherited by BlockPos.
func (p BlockPos) CompareTo(other BlockPos) (ret int) { return p.AsVec3i().CompareTo(other.AsVec3i()) }

// DistSqr mirrors Vec3i.distSqr inherited by BlockPos.
func (p BlockPos) DistSqr(other BlockPos) (ret float64) { return p.AsVec3i().DistSqr(other.AsVec3i()) }

// DistToCenterSqr mirrors Vec3i.distToCenterSqr(double, double, double).
func (p BlockPos) DistToCenterSqr(x float64, y float64, z float64) (ret float64) {
	return p.AsVec3i().DistToCenterSqr(x, y, z)
}

// DistToLowCornerSqr mirrors Vec3i.distToLowCornerSqr(double, double, double).
func (p BlockPos) DistToLowCornerSqr(x float64, y float64, z float64) (ret float64) {
	return p.AsVec3i().DistToLowCornerSqr(x, y, z)
}

// DistManhattan mirrors Vec3i.distManhattan inherited by BlockPos.
func (p BlockPos) DistManhattan(other BlockPos) (ret int) {
	return p.AsVec3i().DistManhattan(other.AsVec3i())
}

// DistChessboard mirrors Vec3i.distChessboard inherited by BlockPos.
func (p BlockPos) DistChessboard(other BlockPos) (ret int) {
	return p.AsVec3i().DistChessboard(other.AsVec3i())
}

// DiffersHorizontally mirrors Vec3i.differsHorizontally inherited by BlockPos.
func (p BlockPos) DiffersHorizontally(other BlockPos) (ret bool) {
	return p.X != other.X || p.Z != other.Z
}

// ToShortString mirrors Vec3i.toShortString inherited by BlockPos.
func (p BlockPos) ToShortString() (ret string) { return p.AsVec3i().ToShortString() }

// `p.Get(axis)` mirrors Vec3i.get(Direction.Axis).
func (p BlockPos) Get(axis DirectionAxis) (ret int) { return axis.ChooseInt(p.X, p.Y, p.Z) }

// BetweenClosedAABB mirrors BlockPos.betweenClosed(AABB). It materialises the
// iteration in the same order as Java (x fastest, then y, then z).
func BetweenClosedAABB(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) (ret []BlockPos) {
	a := BlockPosContaining(minX, minY, minZ)
	b := BlockPosContaining(maxX, maxY, maxZ)
	return BetweenClosed(a, b)
}

// BetweenClosed mirrors BlockPos.betweenClosed(BlockPos, BlockPos).
func BetweenClosed(a BlockPos, b BlockPos) (ret []BlockPos) {
	return BetweenClosedXYZ(
		minInt(a.X, b.X), minInt(a.Y, b.Y), minInt(a.Z, b.Z),
		maxInt(a.X, b.X), maxInt(a.Y, b.Y), maxInt(a.Z, b.Z),
	)
}

// BetweenClosedXYZ mirrors BlockPos.betweenClosed(int, int, int, int, int, int).
func BetweenClosedXYZ(minX int, minY int, minZ int, maxX int, maxY int, maxZ int) (ret []BlockPos) {
	width := maxX - minX + 1
	height := maxY - minY + 1
	depth := maxZ - minZ + 1
	end := width * height * depth
	if end <= 0 {
		return nil
	}
	ret = make([]BlockPos, 0, end)
	for index := 0; index < end; index++ {
		x := index % width
		slice := index / width
		y := slice % height
		z := slice / height
		ret = append(ret, NewBlockPos(minX+x, minY+y, minZ+z))
	}
	return
}

func minInt(a int, b int) (ret int) {
	if a < b {
		return a
	}
	return b
}
