package player

import "github.com/admin-else/strom/mc/util"

// Vec3 is a position or delta in the world.
type Vec3 struct {
	X, Y, Z float64
}

// Rotation is a yaw/pitch pair in degrees.
type Rotation struct {
	Yaw, Pitch float32
}

// PositionMoveRotation mirrors net.minecraft.world.entity.PositionMoveRotation.
// The headless client does not track velocity, so Delta is carried but not
// simulated.
type PositionMoveRotation struct {
	Position Vec3
	Delta    Vec3
	Rotation Rotation
}

// Relative movement flags of Relative.SET_STREAM_CODEC.
const (
	relativeX      = 1 << 0
	relativeY      = 1 << 1
	relativeZ      = 1 << 2
	relativeYRot   = 1 << 3
	relativeXRot   = 1 << 4
	relativeDeltaX = 1 << 5
	relativeDeltaY = 1 << 6
	relativeDeltaZ = 1 << 7
)

// CalculateAbsolute mirrors PositionMoveRotation.calculateAbsolute.
func CalculateAbsolute(source, change PositionMoveRotation, relatives uint32) (ret PositionMoveRotation) {
	offsetX := 0.0
	if relatives&relativeX != 0 {
		offsetX = source.Position.X
	}
	offsetY := 0.0
	if relatives&relativeY != 0 {
		offsetY = source.Position.Y
	}
	offsetZ := 0.0
	if relatives&relativeZ != 0 {
		offsetZ = source.Position.Z
	}
	offsetYRot := float32(0)
	if relatives&relativeYRot != 0 {
		offsetYRot = source.Rotation.Yaw
	}
	offsetXRot := float32(0)
	if relatives&relativeXRot != 0 {
		offsetXRot = source.Rotation.Pitch
	}

	ret.Position = Vec3{offsetX + change.Position.X, offsetY + change.Position.Y, offsetZ + change.Position.Z}
	ret.Rotation = Rotation{offsetYRot + change.Rotation.Yaw, clamp(offsetXRot+change.Rotation.Pitch, -90, 90)}

	ret.Delta = Vec3{
		calculateDelta(source.Delta.X, change.Delta.X, relatives, relativeDeltaX),
		calculateDelta(source.Delta.Y, change.Delta.Y, relatives, relativeDeltaY),
		calculateDelta(source.Delta.Z, change.Delta.Z, relatives, relativeDeltaZ),
	}
	return
}

func calculateDelta(currentDelta, deltaChange float64, relatives uint32, relative uint32) (ret float64) {
	if relatives&relative != 0 {
		return currentDelta + deltaChange
	}
	return deltaChange
}

func clamp(v, min, max float32) (ret float32) {
	return util.ClampFloat(v, min, max)
}
