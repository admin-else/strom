package phys

import (
	"math"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/util"
)

// Vec3 mirrors net.minecraft.world.phys.Vec3. Java declares public final double
// x/y/z plus x()/y()/z(); Go keeps the exported X/Y/Z fields and omits the
// redundant accessors (see the same choice in the neon renderer port).
//
// Deferred with the physics units that need them: the Mojang Codec and
// StreamCodec, the JOML constructors and toVector3f, offsetRandom/offsetRandomXZ
// (`RandomSource`), and Java's toString.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// Vec3 constants mirror the static fields on Vec3.
var (
	Vec3ZERO   = NewVec3(0.0, 0.0, 0.0)
	Vec3X_AXIS = NewVec3(1.0, 0.0, 0.0)
	Vec3Y_AXIS = NewVec3(0.0, 1.0, 0.0)
	Vec3Z_AXIS = NewVec3(0.0, 0.0, 1.0)
)

// NewVec3 mirrors Vec3(double, double, double).
func NewVec3(x float64, y float64, z float64) (ret Vec3) {
	return Vec3{X: x, Y: y, Z: z}
}

// NewVec3FromVec3i mirrors Vec3(Vec3i).
func NewVec3FromVec3i(vec core.Vec3i) (ret Vec3) {
	return NewVec3(float64(vec.X), float64(vec.Y), float64(vec.Z))
}

// AtLowerCornerOf mirrors Vec3.atLowerCornerOf(Vec3i).
func AtLowerCornerOf(pos core.Vec3i) (ret Vec3) {
	return NewVec3(float64(pos.GetX()), float64(pos.GetY()), float64(pos.GetZ()))
}

// AtLowerCornerWithOffset mirrors Vec3.atLowerCornerWithOffset(Vec3i, double, double, double).
func AtLowerCornerWithOffset(pos core.Vec3i, x float64, y float64, z float64) (ret Vec3) {
	return NewVec3(float64(pos.GetX())+x, float64(pos.GetY())+y, float64(pos.GetZ())+z)
}

// AtCenterOf mirrors Vec3.atCenterOf(Vec3i).
func AtCenterOf(pos core.Vec3i) (ret Vec3) { return AtLowerCornerWithOffset(pos, 0.5, 0.5, 0.5) }

// AtCenterOfWithY mirrors Vec3.atCenterOfWithY(Vec3i, double).
func AtCenterOfWithY(pos core.Vec3i, y float64) (ret Vec3) {
	return NewVec3(float64(pos.GetX())+0.5, y, float64(pos.GetZ())+0.5)
}

// AtBottomCenterOf mirrors Vec3.atBottomCenterOf(Vec3i).
func AtBottomCenterOf(pos core.Vec3i) (ret Vec3) { return AtLowerCornerWithOffset(pos, 0.5, 0.0, 0.5) }

// UpFromBottomCenterOf mirrors Vec3.upFromBottomCenterOf(Vec3i, double).
func UpFromBottomCenterOf(pos core.Vec3i, yOffset float64) (ret Vec3) {
	return AtLowerCornerWithOffset(pos, 0.5, yOffset, 0.5)
}

// VectorTo mirrors Vec3.vectorTo(Vec3).
func (v Vec3) VectorTo(other Vec3) (ret Vec3) {
	return NewVec3(other.X-v.X, other.Y-v.Y, other.Z-v.Z)
}

// Normalize mirrors Vec3.normalize().
func (v Vec3) Normalize() (ret Vec3) {
	dist := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if dist < float64(util.MthEpsilon) {
		return Vec3ZERO
	}
	return NewVec3(v.X/dist, v.Y/dist, v.Z/dist)
}

// Dot mirrors Vec3.dot(Vec3).
func (v Vec3) Dot(other Vec3) (ret float64) {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross mirrors Vec3.cross(Vec3).
func (v Vec3) Cross(other Vec3) (ret Vec3) {
	return NewVec3(v.Y*other.Z-v.Z*other.Y, v.Z*other.X-v.X*other.Z, v.X*other.Y-v.Y*other.X)
}

// Subtract mirrors Vec3.subtract(Vec3).
func (v Vec3) Subtract(other Vec3) (ret Vec3) {
	return v.SubtractXYZ(other.X, other.Y, other.Z)
}

// SubtractScalar mirrors Vec3.subtract(double).
func (v Vec3) SubtractScalar(value float64) (ret Vec3) {
	return v.SubtractXYZ(value, value, value)
}

// SubtractXYZ mirrors Vec3.subtract(double, double, double).
func (v Vec3) SubtractXYZ(x float64, y float64, z float64) (ret Vec3) {
	return v.AddXYZ(-x, -y, -z)
}

// AddScalar mirrors Vec3.add(double).
func (v Vec3) AddScalar(value float64) (ret Vec3) {
	return v.AddXYZ(value, value, value)
}

// Add mirrors Vec3.add(Vec3).
func (v Vec3) Add(other Vec3) (ret Vec3) {
	return v.AddXYZ(other.X, other.Y, other.Z)
}

// AddXYZ mirrors Vec3.add(double, double, double).
func (v Vec3) AddXYZ(x float64, y float64, z float64) (ret Vec3) {
	return NewVec3(v.X+x, v.Y+y, v.Z+z)
}

// CloserThan mirrors Vec3.closerThan(Position, double).
func (v Vec3) CloserThan(pos Vec3, distance float64) (ret bool) {
	return v.DistanceToSqr(pos) < distance*distance
}

// CloserThanXZ mirrors Vec3.closerThan(Vec3, double, double).
func (v Vec3) CloserThanXZ(other Vec3, distanceXZ float64, distanceY float64) (ret bool) {
	dx := other.X - v.X
	dy := other.Y - v.Y
	dz := other.Z - v.Z
	return util.LengthSquared2Double(dx, dz) < util.SquareDouble(distanceXZ) && math.Abs(dy) < distanceY
}

// DistanceTo mirrors Vec3.distanceTo(Vec3).
func (v Vec3) DistanceTo(other Vec3) (ret float64) {
	xd := other.X - v.X
	yd := other.Y - v.Y
	zd := other.Z - v.Z
	return math.Sqrt(xd*xd + yd*yd + zd*zd)
}

// DistanceToSqr mirrors Vec3.distanceToSqr(Vec3).
func (v Vec3) DistanceToSqr(other Vec3) (ret float64) {
	return v.DistanceToSqrXYZ(other.X, other.Y, other.Z)
}

// DistanceToSqrXYZ mirrors Vec3.distanceToSqr(double, double, double).
func (v Vec3) DistanceToSqrXYZ(x float64, y float64, z float64) (ret float64) {
	xd := x - v.X
	yd := y - v.Y
	zd := z - v.Z
	return xd*xd + yd*yd + zd*zd
}

// Scale mirrors Vec3.scale(double).
func (v Vec3) Scale(scale float64) (ret Vec3) {
	return v.MultiplyXYZ(scale, scale, scale)
}

// Reverse mirrors Vec3.reverse().
func (v Vec3) Reverse() (ret Vec3) { return v.Scale(-1.0) }

// Multiply mirrors Vec3.multiply(Vec3).
func (v Vec3) Multiply(scale Vec3) (ret Vec3) {
	return v.MultiplyXYZ(scale.X, scale.Y, scale.Z)
}

// MultiplyXYZ mirrors Vec3.multiply(double, double, double).
func (v Vec3) MultiplyXYZ(xScale float64, yScale float64, zScale float64) (ret Vec3) {
	return NewVec3(v.X*xScale, v.Y*yScale, v.Z*zScale)
}

// Horizontal mirrors Vec3.horizontal().
func (v Vec3) Horizontal() (ret Vec3) { return NewVec3(v.X, 0.0, v.Z) }

// Length mirrors Vec3.length().
func (v Vec3) Length() (ret float64) {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// LengthSqr mirrors Vec3.lengthSqr().
func (v Vec3) LengthSqr() (ret float64) { return v.X*v.X + v.Y*v.Y + v.Z*v.Z }

// HorizontalDistance mirrors Vec3.horizontalDistance().
func (v Vec3) HorizontalDistance() (ret float64) {
	return math.Sqrt(v.X*v.X + v.Z*v.Z)
}

// HorizontalDistanceSqr mirrors Vec3.horizontalDistanceSqr().
func (v Vec3) HorizontalDistanceSqr() (ret float64) { return v.X*v.X + v.Z*v.Z }

// Lerp mirrors Vec3.lerp(Vec3, double).
func (v Vec3) Lerp(other Vec3, a float64) (ret Vec3) {
	return NewVec3(
		util.LerpDouble(a, v.X, other.X),
		util.LerpDouble(a, v.Y, other.Y),
		util.LerpDouble(a, v.Z, other.Z),
	)
}

// XRot mirrors Vec3.xRot(float).
func (v Vec3) XRot(radians float32) (ret Vec3) {
	cos := util.Cos(float64(radians))
	sin := util.Sin(float64(radians))
	return NewVec3(v.X, v.Y*float64(cos)+v.Z*float64(sin), v.Z*float64(cos)-v.Y*float64(sin))
}

// YRot mirrors Vec3.yRot(float).
func (v Vec3) YRot(radians float32) (ret Vec3) {
	cos := util.Cos(float64(radians))
	sin := util.Sin(float64(radians))
	return NewVec3(v.X*float64(cos)+v.Z*float64(sin), v.Y, v.Z*float64(cos)-v.X*float64(sin))
}

// ZRot mirrors Vec3.zRot(float).
func (v Vec3) ZRot(radians float32) (ret Vec3) {
	cos := util.Cos(float64(radians))
	sin := util.Sin(float64(radians))
	return NewVec3(v.X*float64(cos)+v.Y*float64(sin), v.Y*float64(cos)-v.X*float64(sin), v.Z)
}

// RotateClockwise90 mirrors Vec3.rotateClockwise90().
func (v Vec3) RotateClockwise90() (ret Vec3) { return NewVec3(-v.Z, v.Y, v.X) }

// DirectionFromRotationVec2 mirrors Vec3.directionFromRotation(Vec2).
func DirectionFromRotationVec2(rotation Vec2) (ret Vec3) {
	return DirectionFromRotation(rotation.X, rotation.Y)
}

// DirectionFromRotation mirrors Vec3.directionFromRotation(float, float).
func DirectionFromRotation(rotX float32, rotY float32) (ret Vec3) {
	yArg := float64(-rotY*util.MthDegToRad - util.MthPI)
	yCos := util.Cos(yArg)
	ySin := util.Sin(yArg)
	xArg := float64(-rotX * util.MthDegToRad)
	xCos := -util.Cos(xArg)
	xSin := util.Sin(xArg)
	return NewVec3(float64(ySin*xCos), float64(xSin), float64(yCos*xCos))
}

// Rotation mirrors Vec3.rotation().
func (v Vec3) Rotation() (ret Vec2) {
	yaw := float32(math.Atan2(-v.X, v.Z)) * (180.0 / float32(math.Pi))
	pitch := float32(math.Asin(-v.Y/math.Sqrt(v.X*v.X+v.Y*v.Y+v.Z*v.Z))) * (180.0 / float32(math.Pi))
	return NewVec2(pitch, yaw)
}

// ProjectedOn mirrors Vec3.projectedOn(Vec3).
func (v Vec3) ProjectedOn(onto Vec3) (ret Vec3) {
	if onto.LengthSqr() == 0.0 {
		return onto
	}
	return onto.Scale(v.Dot(onto)).Scale(1.0 / onto.LengthSqr())
}

// Align mirrors Vec3.align(EnumSet<Direction.Axis>).
func (v Vec3) Align(axes []core.DirectionAxis) (ret Vec3) {
	x, y, z := v.X, v.Y, v.Z
	for _, axis := range axes {
		switch axis {
		case core.AxisX:
			x = float64(util.FloorDouble(v.X))
		case core.AxisY:
			y = float64(util.FloorDouble(v.Y))
		case core.AxisZ:
			z = float64(util.FloorDouble(v.Z))
		}
	}
	return NewVec3(x, y, z)
}

// Get mirrors Vec3.get(Direction.Axis).
func (v Vec3) Get(axis core.DirectionAxis) (ret float64) {
	return axis.ChooseDouble(v.X, v.Y, v.Z)
}

// With mirrors Vec3.with(Direction.Axis, double).
func (v Vec3) With(axis core.DirectionAxis, value float64) (ret Vec3) {
	x, y, z := v.X, v.Y, v.Z
	switch axis {
	case core.AxisX:
		x = value
	case core.AxisY:
		y = value
	case core.AxisZ:
		z = value
	}
	return NewVec3(x, y, z)
}

// Relative mirrors Vec3.relative(Direction, double).
func (v Vec3) Relative(direction core.Direction, distance float64) (ret Vec3) {
	normal := direction.GetUnitVec3i()
	return NewVec3(
		v.X+distance*float64(normal.GetX()),
		v.Y+distance*float64(normal.GetY()),
		v.Z+distance*float64(normal.GetZ()),
	)
}

// IsFinite mirrors Vec3.isFinite().
func (v Vec3) IsFinite() (ret bool) {
	return !math.IsNaN(v.X) && !math.IsInf(v.X, 0) &&
		!math.IsNaN(v.Y) && !math.IsInf(v.Y, 0) &&
		!math.IsNaN(v.Z) && !math.IsInf(v.Z, 0)
}

// Equal mirrors Vec3.equals(Object) using Double.compare semantics: any two NaNs
// are equal and -0.0 is not equal to 0.0. Java's toString is deferred with the
// other formatting-only members.
func (v Vec3) Equal(other Vec3) (ret bool) {
	return doubleCompareEqual(v.X, other.X) && doubleCompareEqual(v.Y, other.Y) && doubleCompareEqual(v.Z, other.Z)
}
