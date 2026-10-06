// Package core mirrors the net.minecraft.core subset the physics and world
// layers need: Direction (+Axis/AxisDirection), Vec3i and BlockPos.
//
// Ported 1:1 from net/minecraft/core/Direction.java, Vec3i.java and
// BlockPos.java. JOML-backed members (Direction.rotate, Direction.getRotation,
// orderedByNearest, toMutable, Vector3f accessors), the Mojang codecs, and the
// exotic BlockPos traversal iterators are deferred until a consumer needs them.
package core

import "github.com/admin-else/strom/mc/util"

// Direction mirrors net.minecraft.core.Direction. The enum ordinal equals the
// 3D data value for every constant.
type Direction int

const (
	DirectionDOWN Direction = iota
	DirectionUP
	DirectionNORTH
	DirectionSOUTH
	DirectionWEST
	DirectionEAST
)

// DirectionVALUES mirrors Direction.VALUES.
var DirectionVALUES = []Direction{
	DirectionDOWN, DirectionUP, DirectionNORTH, DirectionSOUTH, DirectionWEST, DirectionEAST,
}

var (
	directionNames = [6]string{"down", "up", "north", "south", "west", "east"}
	// directionOppositeIndex mirrors the constructor's oppositeIndex argument.
	directionOppositeIndex = [6]int{1, 0, 3, 2, 5, 4}
	// directionData2D mirrors the constructor's data2d argument: vertical
	// directions use -1, horizontals are SOUTH=0, WEST=1, NORTH=2, EAST=3.
	directionData2D = [6]int{-1, -1, 2, 0, 1, 3}
	// directionNormal mirrors the constructor's Vec3i normal argument.
	directionNormal = [6][3]int{
		{0, -1, 0},
		{0, 1, 0},
		{0, 0, -1},
		{0, 0, 1},
		{-1, 0, 0},
		{1, 0, 0},
	}
)

// directionBY_2D_DATA mirrors Direction.BY_2D_DATA: the horizontal directions
// sorted by their data2d value (SOUTH=0, WEST=1, NORTH=2, EAST=3).
var directionBY_2D_DATA = [4]Direction{DirectionSOUTH, DirectionWEST, DirectionNORTH, DirectionEAST}

// DirectionAxis mirrors Direction.Axis.
type DirectionAxis int

const (
	AxisX DirectionAxis = iota
	AxisY
	AxisZ
)

// DirectionAxisVALUES mirrors Direction.Axis.VALUES.
var DirectionAxisVALUES = []DirectionAxis{AxisX, AxisY, AxisZ}

// DirectionAxisDirection mirrors Direction.AxisDirection.
type DirectionAxisDirection int

const (
	DirectionAxisDirectionPOSITIVE DirectionAxisDirection = iota
	DirectionAxisDirectionNEGATIVE
)

// Name mirrors Direction.AxisDirection.getName()/toString().
func (a DirectionAxisDirection) Name() (ret string) {
	if a == DirectionAxisDirectionPOSITIVE {
		return "Towards positive"
	}
	return "Towards negative"
}

// GetStep mirrors Direction.AxisDirection.getStep().
func (a DirectionAxisDirection) GetStep() (ret int) {
	if a == DirectionAxisDirectionPOSITIVE {
		return 1
	}
	return -1
}

// Opposite mirrors Direction.AxisDirection.opposite().
func (a DirectionAxisDirection) Opposite() (ret DirectionAxisDirection) {
	if a == DirectionAxisDirectionPOSITIVE {
		return DirectionAxisDirectionNEGATIVE
	}
	return DirectionAxisDirectionPOSITIVE
}

// Name mirrors Direction.getName()/getSerializedName()/toString().
func (d Direction) Name() (ret string) { return directionNames[d] }

// GetSerializedName mirrors Direction.getSerializedName().
func (d Direction) GetSerializedName() (ret string) { return d.Name() }

// Get3DDataValue mirrors Direction.get3DDataValue().
func (d Direction) Get3DDataValue() (ret int) { return int(d) }

// Get2DDataValue mirrors Direction.get2DDataValue().
func (d Direction) Get2DDataValue() (ret int) { return directionData2D[d] }

// GetOpposite mirrors Direction.getOpposite().
func (d Direction) GetOpposite() (ret Direction) { return Direction(directionOppositeIndex[d]) }

// GetAxis mirrors Direction.getAxis().
func (d Direction) GetAxis() (ret DirectionAxis) {
	switch d {
	case DirectionDOWN, DirectionUP:
		return AxisY
	case DirectionNORTH, DirectionSOUTH:
		return AxisZ
	default:
		return AxisX
	}
}

// GetAxisDirection mirrors Direction.getAxisDirection().
func (d Direction) GetAxisDirection() (ret DirectionAxisDirection) {
	switch d {
	case DirectionDOWN, DirectionNORTH, DirectionWEST:
		return DirectionAxisDirectionNEGATIVE
	default:
		return DirectionAxisDirectionPOSITIVE
	}
}

// GetStepX mirrors Direction.getStepX().
func (d Direction) GetStepX() (ret int) { return directionNormal[d][0] }

// GetStepY mirrors Direction.getStepY().
func (d Direction) GetStepY() (ret int) { return directionNormal[d][1] }

// GetStepZ mirrors Direction.getStepZ().
func (d Direction) GetStepZ() (ret int) { return directionNormal[d][2] }

// GetUnitVec3i mirrors Direction.getUnitVec3i().
func (d Direction) GetUnitVec3i() (ret Vec3i) {
	return NewVec3i(directionNormal[d][0], directionNormal[d][1], directionNormal[d][2])
}

// GetStep mirrors Direction.getStep().
func (d Direction) GetStep() (ret Vec3i) { return d.GetUnitVec3i() }

// ToYRot mirrors Direction.toYRot().
func (d Direction) ToYRot() (ret float32) { return float32(d.Get2DDataValue()&3) * 90.0 }

// IsFacingAngle mirrors Direction.isFacingAngle(float).
func (d Direction) IsFacingAngle(yAngle float32) (ret bool) {
	radians := yAngle * util.MthDegToRad
	dx := -util.Sin(float64(radians))
	dz := util.Cos(float64(radians))
	return float32(directionNormal[d][0])*dx+float32(directionNormal[d][2])*dz > 0.0
}

// GetYRot mirrors Direction.getYRot(Direction).
func GetYRot(direction Direction) (ret float32) {
	switch direction {
	case DirectionNORTH:
		return 180.0
	case DirectionSOUTH:
		return 0.0
	case DirectionWEST:
		return 90.0
	case DirectionEAST:
		return -90.0
	default:
		panic("No y-Rot for vertical axis: " + direction.Name())
	}
}

// DirectionByName mirrors Direction.byName(String); ok is false for unknown names.
func DirectionByName(name string) (result Direction, ok bool) {
	for _, d := range DirectionVALUES {
		if d.Name() == name {
			return d, true
		}
	}
	return DirectionDOWN, false
}

// DirectionFrom3DDataValue mirrors Direction.from3DDataValue(int).
func DirectionFrom3DDataValue(data int) (ret Direction) {
	return Direction(util.AbsInt(data) % len(DirectionVALUES))
}

// DirectionFrom2DDataValue mirrors Direction.from2DDataValue(int).
func DirectionFrom2DDataValue(data int) (ret Direction) {
	return directionBY_2D_DATA[util.AbsInt(data)%len(directionBY_2D_DATA)]
}

// DirectionFromYRot mirrors Direction.fromYRot(double).
func DirectionFromYRot(yRot float64) (ret Direction) {
	return DirectionFrom2DDataValue(util.FloorDouble(yRot/90.0+0.5) & 3)
}

// DirectionFromAxisAndDirection mirrors Direction.fromAxisAndDirection(Axis, AxisDirection).
func DirectionFromAxisAndDirection(axis DirectionAxis, direction DirectionAxisDirection) (ret Direction) {
	switch axis {
	case AxisX:
		if direction == DirectionAxisDirectionPOSITIVE {
			return DirectionEAST
		}
		return DirectionWEST
	case AxisY:
		if direction == DirectionAxisDirectionPOSITIVE {
			return DirectionUP
		}
		return DirectionDOWN
	default:
		if direction == DirectionAxisDirectionPOSITIVE {
			return DirectionSOUTH
		}
		return DirectionNORTH
	}
}

// DirectionGetApproximateNearest mirrors Direction.getApproximateNearest(float, float, float).
func DirectionGetApproximateNearest(dx float32, dy float32, dz float32) (ret Direction) {
	result := DirectionNORTH
	highestDot := smallestNonzeroFloat32
	for _, direction := range DirectionVALUES {
		normal := directionNormal[direction]
		dot := dx*float32(normal[0]) + dy*float32(normal[1]) + dz*float32(normal[2])
		if dot > highestDot {
			highestDot = dot
			result = direction
		}
	}
	return result
}

// DirectionGetNearest mirrors Direction.getNearest(int, int, int, *Direction).
// orElse is returned when no axis dominates; ok reports whether a direction was
// chosen from the input (matching the Java @Nullable return).
func DirectionGetNearest(x int, y int, z int, orElse Direction) (ret Direction, ok bool) {
	absX, absY, absZ := util.AbsInt(x), util.AbsInt(y), util.AbsInt(z)
	switch {
	case absX > absZ && absX > absY:
		if x < 0 {
			return DirectionWEST, true
		}
		return DirectionEAST, true
	case absZ > absX && absZ > absY:
		if z < 0 {
			return DirectionNORTH, true
		}
		return DirectionSOUTH, true
	case absY > absX && absY > absZ:
		if y < 0 {
			return DirectionDOWN, true
		}
		return DirectionUP, true
	default:
		return orElse, false
	}
}

// DirectionGet mirrors Direction.get(AxisDirection, Axis).
func DirectionGet(axisDirection DirectionAxisDirection, axis DirectionAxis) (ret Direction) {
	for _, d := range DirectionVALUES {
		if d.GetAxisDirection() == axisDirection && d.GetAxis() == axis {
			return d
		}
	}
	panic("No such direction: " + axis.Name() + " " + axisDirection.Name())
}

// GetClockWise mirrors Direction.getClockWise() (the Y-axis rotation).
func (d Direction) GetClockWise() (ret Direction) {
	switch d {
	case DirectionNORTH:
		return DirectionEAST
	case DirectionSOUTH:
		return DirectionWEST
	case DirectionWEST:
		return DirectionNORTH
	case DirectionEAST:
		return DirectionSOUTH
	default:
		panic("Unable to get Y-rotated facing of " + d.Name())
	}
}

// GetCounterClockWise mirrors Direction.getCounterClockWise() (Y-axis).
func (d Direction) GetCounterClockWise() (ret Direction) {
	switch d {
	case DirectionNORTH:
		return DirectionWEST
	case DirectionSOUTH:
		return DirectionEAST
	case DirectionWEST:
		return DirectionSOUTH
	case DirectionEAST:
		return DirectionNORTH
	default:
		panic("Unable to get CCW facing of " + d.Name())
	}
}

func (d Direction) getClockWiseX() (ret Direction) {
	switch d {
	case DirectionDOWN:
		return DirectionSOUTH
	case DirectionUP:
		return DirectionNORTH
	case DirectionNORTH:
		return DirectionDOWN
	case DirectionSOUTH:
		return DirectionUP
	default:
		panic("Unable to get X-rotated facing of " + d.Name())
	}
}

func (d Direction) getCounterClockWiseX() (ret Direction) {
	switch d {
	case DirectionDOWN:
		return DirectionNORTH
	case DirectionUP:
		return DirectionSOUTH
	case DirectionNORTH:
		return DirectionUP
	case DirectionSOUTH:
		return DirectionDOWN
	default:
		panic("Unable to get X-rotated facing of " + d.Name())
	}
}

func (d Direction) getClockWiseZ() (ret Direction) {
	switch d {
	case DirectionDOWN:
		return DirectionWEST
	case DirectionUP:
		return DirectionEAST
	case DirectionWEST:
		return DirectionUP
	case DirectionEAST:
		return DirectionDOWN
	default:
		panic("Unable to get Z-rotated facing of " + d.Name())
	}
}

func (d Direction) getCounterClockWiseZ() (ret Direction) {
	switch d {
	case DirectionDOWN:
		return DirectionEAST
	case DirectionUP:
		return DirectionWEST
	case DirectionWEST:
		return DirectionDOWN
	case DirectionEAST:
		return DirectionUP
	default:
		panic("Unable to get Z-rotated facing of " + d.Name())
	}
}

// GetClockWiseAxis mirrors Direction.getClockWise(Axis).
func (d Direction) GetClockWiseAxis(axis DirectionAxis) (ret Direction) {
	switch axis {
	case AxisX:
		if d != DirectionWEST && d != DirectionEAST {
			return d.getClockWiseX()
		}
		return d
	case AxisY:
		if d != DirectionUP && d != DirectionDOWN {
			return d.GetClockWise()
		}
		return d
	default:
		if d != DirectionNORTH && d != DirectionSOUTH {
			return d.getClockWiseZ()
		}
		return d
	}
}

// GetCounterClockWiseAxis mirrors Direction.getCounterClockWise(Axis).
func (d Direction) GetCounterClockWiseAxis(axis DirectionAxis) (ret Direction) {
	switch axis {
	case AxisX:
		if d != DirectionWEST && d != DirectionEAST {
			return d.getCounterClockWiseX()
		}
		return d
	case AxisY:
		if d != DirectionUP && d != DirectionDOWN {
			return d.GetCounterClockWise()
		}
		return d
	default:
		if d != DirectionNORTH && d != DirectionSOUTH {
			return d.getCounterClockWiseZ()
		}
		return d
	}
}

// Name mirrors Direction.Axis.getName().
func (a DirectionAxis) Name() (ret string) {
	switch a {
	case AxisX:
		return "x"
	case AxisY:
		return "y"
	default:
		return "z"
	}
}

// IsVertical mirrors Direction.Axis.isVertical().
func (a DirectionAxis) IsVertical() (ret bool) { return a == AxisY }

// IsHorizontal mirrors Direction.Axis.isHorizontal().
func (a DirectionAxis) IsHorizontal() (ret bool) { return a == AxisX || a == AxisZ }

// DirectionAxisByName mirrors Direction.Axis.byName(String).
func DirectionAxisByName(name string) (result DirectionAxis, ok bool) {
	for _, a := range DirectionAxisVALUES {
		if a.Name() == name {
			return a, true
		}
	}
	return AxisX, false
}

// GetPositive mirrors Direction.Axis.getPositive().
func (a DirectionAxis) GetPositive() (ret Direction) {
	switch a {
	case AxisX:
		return DirectionEAST
	case AxisY:
		return DirectionUP
	default:
		return DirectionSOUTH
	}
}

// GetNegative mirrors Direction.Axis.getNegative().
func (a DirectionAxis) GetNegative() (ret Direction) {
	switch a {
	case AxisX:
		return DirectionWEST
	case AxisY:
		return DirectionDOWN
	default:
		return DirectionNORTH
	}
}

// GetDirections mirrors Direction.Axis.getDirections().
func (a DirectionAxis) GetDirections() (ret []Direction) {
	return []Direction{a.GetPositive(), a.GetNegative()}
}

// ChooseInt mirrors Direction.Axis.choose(int, int, int).
func (a DirectionAxis) ChooseInt(x int, y int, z int) (ret int) {
	switch a {
	case AxisX:
		return x
	case AxisY:
		return y
	default:
		return z
	}
}

// ChooseDouble mirrors Direction.Axis.choose(double, double, double).
func (a DirectionAxis) ChooseDouble(x float64, y float64, z float64) (ret float64) {
	switch a {
	case AxisX:
		return x
	case AxisY:
		return y
	default:
		return z
	}
}

// ChooseBool mirrors Direction.Axis.choose(boolean, boolean, boolean).
func (a DirectionAxis) ChooseBool(x bool, y bool, z bool) (ret bool) {
	switch a {
	case AxisX:
		return x
	case AxisY:
		return y
	default:
		return z
	}
}

const smallestNonzeroFloat32 float32 = 1.401298464324817e-45
