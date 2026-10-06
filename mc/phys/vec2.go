package phys

import (
	"math"

	"github.com/admin-else/strom/mc/util"
)

// Vec2 mirrors net.minecraft.world.phys.Vec2. Java declares public final float
// x/y; Go keeps the exported X/Y fields. The Mojang Codec is omitted (it only
// matters for JSON/NBT serialization, which the headless world does not use yet).
type Vec2 struct {
	X float32
	Y float32
}

// Vec2 constants mirror the static fields on Vec2.
var (
	Vec2ZERO       = NewVec2(0.0, 0.0)
	Vec2ONE        = NewVec2(1.0, 1.0)
	Vec2UNIT_X     = NewVec2(1.0, 0.0)
	Vec2NEG_UNIT_X = NewVec2(-1.0, 0.0)
	Vec2UNIT_Y     = NewVec2(0.0, 1.0)
	Vec2NEG_UNIT_Y = NewVec2(0.0, -1.0)
	Vec2MAX        = NewVec2(math.MaxFloat32, math.MaxFloat32)
	Vec2MIN        = NewVec2(math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32)
)

// NewVec2 mirrors Vec2(float, float).
func NewVec2(x float32, y float32) (ret Vec2) {
	return Vec2{X: x, Y: y}
}

// Scale mirrors Vec2.scale(float).
func (v Vec2) Scale(s float32) (ret Vec2) { return NewVec2(v.X*s, v.Y*s) }

// Dot mirrors Vec2.dot(Vec2).
func (v Vec2) Dot(other Vec2) (ret float32) { return v.X*other.X + v.Y*other.Y }

// Add mirrors Vec2.add(Vec2).
func (v Vec2) Add(other Vec2) (ret Vec2) { return NewVec2(v.X+other.X, v.Y+other.Y) }

// AddScalar mirrors Vec2.add(float).
func (v Vec2) AddScalar(s float32) (ret Vec2) { return NewVec2(v.X+s, v.Y+s) }

// Normalized mirrors Vec2.normalized().
func (v Vec2) Normalized() (ret Vec2) {
	dist := util.Sqrt(v.X*v.X + v.Y*v.Y)
	if dist < float32(1.0e-4) {
		return Vec2ZERO
	}
	return NewVec2(v.X/dist, v.Y/dist)
}

// Length mirrors Vec2.length().
func (v Vec2) Length() (ret float32) { return util.Sqrt(v.X*v.X + v.Y*v.Y) }

// LengthSquared mirrors Vec2.lengthSquared().
func (v Vec2) LengthSquared() (ret float32) { return v.X*v.X + v.Y*v.Y }

// DistanceToSqr mirrors Vec2.distanceToSqr(Vec2).
func (v Vec2) DistanceToSqr(p Vec2) (ret float32) {
	xd := p.X - v.X
	yd := p.Y - v.Y
	return xd*xd + yd*yd
}

// Negated mirrors Vec2.negated().
func (v Vec2) Negated() (ret Vec2) { return NewVec2(-v.X, -v.Y) }

// Rotate mirrors Vec2.rotate(double).
func (v Vec2) Rotate(angleRadians float64) (ret Vec2) {
	cosine := util.Cos(angleRadians)
	sine := util.Sin(angleRadians)
	return NewVec2(v.X*cosine-v.Y*sine, v.Y*cosine+v.X*sine)
}

// Equal mirrors Vec2.equals(Object); Java compares the floats with ==.
func (v Vec2) Equal(other Vec2) (ret bool) { return v.X == other.X && v.Y == other.Y }
