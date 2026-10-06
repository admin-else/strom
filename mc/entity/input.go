package entity

import (
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/util"
)

// Input mirrors the movement subset of the client key state
// (net.minecraft.client.player.ClientInput / KeyboardInput).
type Input struct {
	Forward  bool
	Backward bool
	Left     bool
	Right    bool
	Jump     bool
	Shift    bool
	Sprint   bool
}

// AttributeDefaults.SNEAKING_SPEED.
const sneakingSpeedDefault float32 = 0.3

// GetMoveVector mirrors ClientInput.getMoveVector(): x is left(+)/right(-),
// y is forward(+)/backward(-), normalized to unit length.
func (in Input) GetMoveVector() (ret phys.Vec2) {
	x := float32(0.0)
	y := float32(0.0)
	if in.Left {
		x += 1.0
	}
	if in.Right {
		x -= 1.0
	}
	if in.Forward {
		y += 1.0
	}
	if in.Backward {
		y -= 1.0
	}
	return phys.NewVec2(x, y).Normalized()
}

// ModifyInput mirrors LocalPlayer.modifyInput(Vec2) with the default sneaking
// speed. The item-use multiplier is deferred (no items yet).
func ModifyInput(input phys.Vec2, movingSlowly bool) (ret phys.Vec2) {
	if input.LengthSquared() == 0.0 {
		return input
	}
	newInput := input.Scale(0.98)
	if movingSlowly {
		newInput = newInput.Scale(sneakingSpeedDefault)
	}
	return ModifyInputSpeedForSquareMovement(newInput)
}

// ModifyInputSpeedForSquareMovement mirrors LocalPlayer.modifyInputSpeedForSquareMovement(Vec2).
func ModifyInputSpeedForSquareMovement(input phys.Vec2) (ret phys.Vec2) {
	length := input.Length()
	if length <= 0.0 {
		return input
	}
	direction := input.Scale(1.0 / length)
	distanceToUnitSquare := DistanceToUnitSquare(direction)
	modifiedLength := minFloat32(length*distanceToUnitSquare, 1.0)
	return direction.Scale(modifiedLength)
}

// DistanceToUnitSquare mirrors LocalPlayer.distanceToUnitSquare(Vec2).
func DistanceToUnitSquare(direction phys.Vec2) (ret float32) {
	directionX := util.AbsFloat(direction.X)
	directionY := util.AbsFloat(direction.Y)
	var tan float32
	if directionY > directionX {
		tan = directionX / directionY
	} else {
		tan = directionY / directionX
	}
	return util.Sqrt(1.0 + util.SquareFloat(tan))
}

// ApplyInput mirrors LocalPlayer.applyInput() for a controlled camera: it fills
// the (xxa, zza) movement vector from the key input.
func (e *Entity) ApplyInput(in Input, movingSlowly bool) (xxa float32, zza float32) {
	modifiedInput := ModifyInput(in.GetMoveVector(), movingSlowly)
	return modifiedInput.X, modifiedInput.Y
}

// TravelWithInput is the end-to-end input -> travel wiring: it turns the key
// state into the movement vector and runs the LivingEntity travel step. Speed
// and gravity come from the entity attributes.
func (e *Entity) TravelWithInput(in Input, movingSlowly bool) {
	xxa, zza := e.ApplyInput(in, movingSlowly)
	e.jumping = in.Jump
	e.SetSprinting(in.Sprint)
	e.ShiftKeyDown = in.Shift
	if in.Jump && e.OnGround() {
		e.JumpFromGround()
	}
	e.Travel(phys.NewVec3(float64(xxa), 0.0, float64(zza)))
}

func minFloat32(a float32, b float32) (ret float32) {
	if a < b {
		return a
	}
	return b
}
