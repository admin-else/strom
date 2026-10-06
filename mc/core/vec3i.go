package core

import (
	"strconv"

	"github.com/admin-else/strom/mc/util"
)

// Vec3i mirrors net.minecraft.core.Vec3i. Java declares private x/y/z with
// getX/getY/getZ; Go keeps exported X/Y/Z fields and the getters for call-site
// fidelity. The Mojang codecs, JOML toMutable, the MoreObjects toString and the
// Position-based overloads are deferred until a consumer needs them.
type Vec3i struct {
	X int
	Y int
	Z int
}

// Vec3iZERO mirrors Vec3i.ZERO.
var Vec3iZERO = NewVec3i(0, 0, 0)

// NewVec3i mirrors Vec3i(int, int, int).
func NewVec3i(x int, y int, z int) (ret Vec3i) { return Vec3i{X: x, Y: y, Z: z} }

// GetX mirrors Vec3i.getX().
func (v Vec3i) GetX() (ret int) { return v.X }

// GetY mirrors Vec3i.getY().
func (v Vec3i) GetY() (ret int) { return v.Y }

// GetZ mirrors Vec3i.getZ().
func (v Vec3i) GetZ() (ret int) { return v.Z }

// Equal mirrors Vec3i.equals(Object).
func (v Vec3i) Equal(other Vec3i) (ret bool) {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// HashCode mirrors Vec3i.hashCode().
func (v Vec3i) HashCode() (ret int) { return (v.Y+v.Z*31)*31 + v.X }

// CompareTo mirrors Vec3i.compareTo(Vec3i).
func (v Vec3i) CompareTo(other Vec3i) (ret int) {
	if v.Y == other.Y {
		if v.Z == other.Z {
			return v.X - other.X
		}
		return v.Z - other.Z
	}
	return v.Y - other.Y
}

// Offset mirrors Vec3i.offset(int, int, int).
func (v Vec3i) Offset(x int, y int, z int) (ret Vec3i) {
	if x == 0 && y == 0 && z == 0 {
		return v
	}
	return NewVec3i(v.X+x, v.Y+y, v.Z+z)
}

// OffsetVec mirrors Vec3i.offset(Vec3i).
func (v Vec3i) OffsetVec(vec Vec3i) (ret Vec3i) { return v.Offset(vec.X, vec.Y, vec.Z) }

// Subtract mirrors Vec3i.subtract(Vec3i).
func (v Vec3i) Subtract(vec Vec3i) (ret Vec3i) { return v.Offset(-vec.X, -vec.Y, -vec.Z) }

// Multiply mirrors Vec3i.multiply(int).
func (v Vec3i) Multiply(scale int) (ret Vec3i) {
	if scale == 1 {
		return v
	}
	if scale == 0 {
		return Vec3iZERO
	}
	return NewVec3i(v.X*scale, v.Y*scale, v.Z*scale)
}

// MultiplyXYZ mirrors Vec3i.multiply(int, int, int).
func (v Vec3i) MultiplyXYZ(xScale int, yScale int, zScale int) (ret Vec3i) {
	return NewVec3i(v.X*xScale, v.Y*yScale, v.Z*zScale)
}

// Above mirrors Vec3i.above().
func (v Vec3i) Above() (ret Vec3i) { return v.AboveSteps(1) }

// AboveSteps mirrors Vec3i.above(int).
func (v Vec3i) AboveSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionUP, steps) }

// Below mirrors Vec3i.below().
func (v Vec3i) Below() (ret Vec3i) { return v.BelowSteps(1) }

// BelowSteps mirrors Vec3i.below(int).
func (v Vec3i) BelowSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionDOWN, steps) }

// North mirrors Vec3i.north().
func (v Vec3i) North() (ret Vec3i) { return v.NorthSteps(1) }

// NorthSteps mirrors Vec3i.north(int).
func (v Vec3i) NorthSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionNORTH, steps) }

// South mirrors Vec3i.south().
func (v Vec3i) South() (ret Vec3i) { return v.SouthSteps(1) }

// SouthSteps mirrors Vec3i.south(int).
func (v Vec3i) SouthSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionSOUTH, steps) }

// West mirrors Vec3i.west().
func (v Vec3i) West() (ret Vec3i) { return v.WestSteps(1) }

// WestSteps mirrors Vec3i.west(int).
func (v Vec3i) WestSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionWEST, steps) }

// East mirrors Vec3i.east().
func (v Vec3i) East() (ret Vec3i) { return v.EastSteps(1) }

// EastSteps mirrors Vec3i.east(int).
func (v Vec3i) EastSteps(steps int) (ret Vec3i) { return v.RelativeSteps(DirectionEAST, steps) }

// Relative mirrors Vec3i.relative(Direction).
func (v Vec3i) Relative(direction Direction) (ret Vec3i) { return v.RelativeSteps(direction, 1) }

// RelativeSteps mirrors Vec3i.relative(Direction, int).
func (v Vec3i) RelativeSteps(direction Direction, steps int) (ret Vec3i) {
	if steps == 0 {
		return v
	}
	return NewVec3i(v.X+direction.GetStepX()*steps, v.Y+direction.GetStepY()*steps, v.Z+direction.GetStepZ()*steps)
}

// RelativeAxis mirrors Vec3i.relative(Direction.Axis, int).
func (v Vec3i) RelativeAxis(axis DirectionAxis, steps int) (ret Vec3i) {
	if steps == 0 {
		return v
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
	return NewVec3i(v.X+xStep, v.Y+yStep, v.Z+zStep)
}

// Cross mirrors Vec3i.cross(Vec3i).
func (v Vec3i) Cross(upVector Vec3i) (ret Vec3i) {
	return NewVec3i(
		v.Y*upVector.Z-v.Z*upVector.Y,
		v.Z*upVector.X-v.X*upVector.Z,
		v.X*upVector.Y-v.Y*upVector.X,
	)
}

// CloserThan mirrors Vec3i.closerThan(Vec3i, double).
func (v Vec3i) CloserThan(pos Vec3i, distance float64) (ret bool) {
	return v.DistSqr(pos) < util.SquareDouble(distance)
}

// DistSqr mirrors Vec3i.distSqr(Vec3i).
func (v Vec3i) DistSqr(pos Vec3i) (ret float64) {
	return v.DistToLowCornerSqr(float64(pos.X), float64(pos.Y), float64(pos.Z))
}

// DistToCenterSqr mirrors Vec3i.distToCenterSqr(double, double, double).
func (v Vec3i) DistToCenterSqr(x float64, y float64, z float64) (ret float64) {
	dx := float64(v.X) + 0.5 - x
	dy := float64(v.Y) + 0.5 - y
	dz := float64(v.Z) + 0.5 - z
	return dx*dx + dy*dy + dz*dz
}

// DistToLowCornerSqr mirrors Vec3i.distToLowCornerSqr(double, double, double).
func (v Vec3i) DistToLowCornerSqr(x float64, y float64, z float64) (ret float64) {
	dx := float64(v.X) - x
	dy := float64(v.Y) - y
	dz := float64(v.Z) - z
	return dx*dx + dy*dy + dz*dz
}

// DistManhattan mirrors Vec3i.distManhattan(Vec3i).
func (v Vec3i) DistManhattan(pos Vec3i) (ret int) {
	xd := float32(util.AbsInt(pos.X - v.X))
	yd := float32(util.AbsInt(pos.Y - v.Y))
	zd := float32(util.AbsInt(pos.Z - v.Z))
	return int(xd + yd + zd)
}

// DistChessboard mirrors Vec3i.distChessboard(Vec3i).
func (v Vec3i) DistChessboard(pos Vec3i) (ret int) {
	xd := util.AbsInt(v.X - pos.X)
	yd := util.AbsInt(v.Y - pos.Y)
	zd := util.AbsInt(v.Z - pos.Z)
	return maxInt(maxInt(xd, yd), zd)
}

// DiffersHorizontally mirrors Vec3i.differsHorizontally(Vec3i).
func (v Vec3i) DiffersHorizontally(pos Vec3i) (ret bool) { return v.X != pos.X || v.Z != pos.Z }

// Get mirrors Vec3i.get(Direction.Axis).
func (v Vec3i) Get(axis DirectionAxis) (ret int) { return axis.ChooseInt(v.X, v.Y, v.Z) }

// ToShortString mirrors Vec3i.toShortString().
func (v Vec3i) ToShortString() (ret string) {
	return strconv.Itoa(v.X) + ", " + strconv.Itoa(v.Y) + ", " + strconv.Itoa(v.Z)
}

func maxInt(a int, b int) (ret int) {
	if a > b {
		return a
	}
	return b
}
