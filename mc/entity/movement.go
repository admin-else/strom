package entity

import (
	"math"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/util"
)

// MovementContext is the optional level slice the LivingEntity movement reads:
// fluid and climbable queries plus a no-collision test for the sneak edge
// protection. A CollisionGetter that does not implement it makes every query
// return the empty/default answer.
type MovementContext interface {
	IsInWater(pos core.BlockPos) bool
	IsInLava(pos core.BlockPos) bool
	OnClimbable(pos core.BlockPos) bool
	IsInShallowFluid(pos core.BlockPos) bool
	BlockFriction(pos core.BlockPos) float32
	NoCollision(source *Entity, box phys.AABB) bool
}

// MovementContext returns the level's movement context, or nil.
func (e *Entity) MovementContext() (ret MovementContext) {
	if context, ok := e.level.(MovementContext); ok {
		return context
	}
	return nil
}

// BlockPosition mirrors Entity.blockPosition().
func (e *Entity) BlockPosition() (ret core.BlockPos) {
	return core.BlockPosContaining(e.X, e.Y, e.Z)
}

// GetBlockPosBelowThatAffectsMyMovement mirrors Entity.getBlockPosBelowThatAffectsMyMovement().
func (e *Entity) GetBlockPosBelowThatAffectsMyMovement() (ret core.BlockPos) {
	return core.BlockPosContaining(e.X, e.BoundingBox.MinY-0.500001, e.Z)
}

// IsInWater mirrors Entity.isInWater().
func (e *Entity) IsInWater() (ret bool) {
	if ctx := e.MovementContext(); ctx != nil {
		return ctx.IsInWater(e.BlockPosition())
	}
	return false
}

// IsInLava mirrors Entity.isInLava().
func (e *Entity) IsInLava() (ret bool) {
	if ctx := e.MovementContext(); ctx != nil {
		return ctx.IsInLava(e.BlockPosition())
	}
	return false
}

// IsInLiquid mirrors Entity.isInLiquid().
func (e *Entity) IsInLiquid() (ret bool) { return e.IsInWater() || e.IsInLava() }

// OnClimbable mirrors LivingEntity.onClimbable().
func (e *Entity) OnClimbable() (ret bool) {
	if ctx := e.MovementContext(); ctx != nil {
		return ctx.OnClimbable(e.BlockPosition())
	}
	return false
}

// IsFallFlying mirrors LivingEntity.isFallFlying().
func (e *Entity) IsFallFlying() (ret bool) { return e.fallFlying }

// SetFallFlying sets the fall-flying flag.
func (e *Entity) SetFallFlying(fallFlying bool) { e.fallFlying = fallFlying }

// IsAffectedByFluids mirrors LivingEntity.isAffectedByFluids().
func (e *Entity) IsAffectedByFluids() (ret bool) { return true }

// ShouldDiscardFriction mirrors LivingEntity.shouldDiscardFriction().
func (e *Entity) ShouldDiscardFriction() (ret bool) { return e.fallFlying }

// GetDefaultGravity mirrors LivingEntity.getDefaultGravity().
func (e *Entity) GetDefaultGravity() (ret float64) { return e.AttributeValue(AttributeGRAVITY) }

// GetGravity mirrors Entity.getGravity().
func (e *Entity) GetGravity() (ret float64) {
	if e.noGravity {
		return 0.0
	}
	return e.GetDefaultGravity()
}

// GetEffectiveGravity mirrors LivingEntity.getEffectiveGravity().
func (e *Entity) GetEffectiveGravity() (ret float64) {
	isFalling := e.GetDeltaMovement().Y <= 0.0
	if isFalling && e.HasEffect(MobEffectSLOW_FALLING) {
		return minFloat64(e.GetGravity(), 0.01)
	}
	return e.GetGravity()
}

// GetJumpBoostPower mirrors LivingEntity.getJumpBoostPower().
func (e *Entity) GetJumpBoostPower() (ret float32) {
	if amplifier, ok := e.GetEffectAmplifier(MobEffectJUMP_BOOST); ok {
		return 0.1 * float32(amplifier+1)
	}
	return 0.0
}

// GetJumpPower mirrors LivingEntity.getJumpPower() with a block jump factor.
func (e *Entity) GetJumpPower(multiplier float32, blockJumpFactor float32) (ret float32) {
	return float32(e.AttributeValue(AttributeJUMP_STRENGTH))*multiplier*blockJumpFactor + e.GetJumpBoostPower()
}

// JumpFromGround mirrors LivingEntity.jumpFromGround().
func (e *Entity) JumpFromGround() {
	jumpPower := e.GetJumpPower(1.0, 1.0)
	if jumpPower <= 1.0e-5 {
		return
	}
	movement := e.GetDeltaMovement()
	e.SetDeltaMovement(phys.NewVec3(movement.X, mathMax64(float64(jumpPower), movement.Y), movement.Z))
	if e.sprinting {
		angle := float32(e.yRot) * util.MthDegToRad
		e.AddDeltaMovement(phys.NewVec3(-float64(util.Sin(float64(angle)))*0.2, 0.0, float64(util.Cos(float64(angle)))*0.2))
	}
}

// GetWaterSlowDown mirrors LivingEntity.getWaterSlowDown().
func (e *Entity) GetWaterSlowDown() (ret float32) { return 0.8 }

// ComputeModifiedFriction mirrors LivingEntity.computeModifiedFriction(float, float).
func ComputeModifiedFriction(friction float32, modifier float32) (ret float32) {
	return util.ClampFloat(1.0-(1.0-friction)*modifier, 0.0, 1.0)
}

// ShouldTravelInFluid mirrors LivingEntity.shouldTravelInFluid(FluidState).
func (e *Entity) ShouldTravelInFluid() (ret bool) {
	return e.IsInLiquid() && e.IsAffectedByFluids()
}

// GetLookAngle mirrors Entity.getLookAngle() (calculateViewVector).
func (e *Entity) GetLookAngle() (ret phys.Vec3) {
	return phys.DirectionFromRotation(e.xRot, e.yRot)
}

// Travel mirrors LivingEntity.travel(Vec3): fluid, fall-flying or air.
func (e *Entity) Travel(input phys.Vec3) {
	if e.ShouldTravelInFluid() {
		e.TravelInFluid(input)
	} else if e.IsFallFlying() {
		e.TravelFallFlying(input)
	} else {
		e.TravelInAir(input)
	}
}

// TravelInAir mirrors LivingEntity.travelInAir(Vec3).
func (e *Entity) TravelInAir(input phys.Vec3) {
	posBelow := e.GetBlockPosBelowThatAffectsMyMovement()
	blockFriction := float32(1.0)
	if e.OnGround() {
		friction := float32(0.6)
		if ctx := e.MovementContext(); ctx != nil {
			friction = ctx.BlockFriction(posBelow)
		}
		blockFriction = ComputeModifiedFriction(friction, float32(e.AttributeValue(AttributeFRICTION_MODIFIER)))
	}
	movement := e.HandleRelativeFrictionAndCalculateMovement(input, blockFriction)
	movementY := movement.Y
	if amplifier, ok := e.GetEffectAmplifier(MobEffectLEVITATION); ok {
		movementY += (0.05*float64(amplifier+1) - movement.Y) * 0.2
	} else {
		movementY -= e.GetEffectiveGravity()
	}

	if e.ShouldDiscardFriction() {
		e.SetDeltaMovement(phys.NewVec3(movement.X, movementY, movement.Z))
		return
	}
	airDragModifier := float32(e.AttributeValue(AttributeAIR_DRAG_MODIFIER))
	airDrag := ComputeModifiedFriction(0.91, airDragModifier)
	friction := blockFriction * airDrag
	verticalFriction := ComputeModifiedFriction(0.98, airDragModifier)
	e.SetDeltaMovement(phys.NewVec3(movement.X*float64(friction), movementY*float64(verticalFriction), movement.Z*float64(friction)))
}

// TravelInFluid mirrors LivingEntity.travelInFluid(Vec3).
func (e *Entity) TravelInFluid(input phys.Vec3) {
	isFalling := e.GetDeltaMovement().Y <= 0.0
	oldY := e.Y
	baseGravity := e.GetEffectiveGravity()
	if e.IsInWater() {
		e.TravelInWater(input, baseGravity, isFalling, oldY)
	} else {
		e.TravelInLava(input, baseGravity, isFalling, oldY)
	}
}

// TravelInWater mirrors LivingEntity.travelInWater(...).
func (e *Entity) TravelInWater(input phys.Vec3, baseGravity float64, isFalling bool, oldY float64) {
	slowDown := e.GetWaterSlowDown()
	if e.sprinting {
		slowDown = 0.9
	}
	speed := float32(0.02)
	waterWalker := float32(e.AttributeValue(AttributeWATER_MOVEMENT_EFFICIENCY))
	if !e.OnGround() {
		waterWalker *= 0.5
	}
	if waterWalker > 0.0 {
		slowDown += (0.54600006 - slowDown) * waterWalker
		speed += (e.GetSpeed() - speed) * waterWalker
	}
	if e.HasEffect(MobEffectDOLPHINS_GRACE) {
		slowDown = 0.96
	}
	e.MoveRelative(speed, input)
	e.Move(MoverTypeSELF, e.GetDeltaMovement())
	movement := e.GetDeltaMovement()
	if e.horizontalCollision && e.OnClimbable() {
		movement = phys.NewVec3(movement.X, 0.2, movement.Z)
	}
	movement = movement.MultiplyXYZ(float64(slowDown), 0.8, float64(slowDown))
	e.SetDeltaMovement(e.GetFluidFallingAdjustedMovement(baseGravity, isFalling, movement))
	e.JumpOutOfFluid(oldY)
}

// TravelInLava mirrors LivingEntity.travelInLava(...).
func (e *Entity) TravelInLava(input phys.Vec3, baseGravity float64, isFalling bool, oldY float64) {
	e.MoveRelative(0.02, input)
	e.Move(MoverTypeSELF, e.GetDeltaMovement())
	if e.IsInShallowFluid() {
		e.SetDeltaMovement(e.GetDeltaMovement().MultiplyXYZ(0.5, 0.8, 0.5))
		e.SetDeltaMovement(e.GetFluidFallingAdjustedMovement(baseGravity, isFalling, e.GetDeltaMovement()))
	} else {
		e.SetDeltaMovement(e.GetDeltaMovement().Scale(0.5))
	}
	if baseGravity != 0.0 {
		e.AddDeltaMovement(phys.NewVec3(0.0, -baseGravity/4.0, 0.0))
	}
	e.JumpOutOfFluid(oldY)
}

// IsInShallowFluid mirrors LivingEntity.isInShallowFluid(TagKey<Fluid>).
func (e *Entity) IsInShallowFluid() (ret bool) {
	if ctx := e.MovementContext(); ctx != nil {
		return ctx.IsInShallowFluid(e.BlockPosition())
	}
	return true
}

// JumpOutOfFluid mirrors LivingEntity.jumpOutOfFluid(double).
func (e *Entity) JumpOutOfFluid(oldY float64) {
	if !e.jumping && !e.sprinting {
		return
	}
	movement := e.GetDeltaMovement()
	if e.horizontalCollision && e.isFree(movement.X, movement.Y+0.6-float64(e.Y)+oldY, movement.Z) {
		e.SetDeltaMovement(phys.NewVec3(movement.X, 0.3, movement.Z))
	}
}

// GetFluidFallingAdjustedMovement mirrors LivingEntity.getFluidFallingAdjustedMovement(...).
func (e *Entity) GetFluidFallingAdjustedMovement(baseGravity float64, isFalling bool, movement phys.Vec3) (ret phys.Vec3) {
	if baseGravity != 0.0 && !e.sprinting {
		var yd float64
		if isFalling && util.AbsFloat(float32(movement.Y-0.005)) >= 0.003 && util.AbsFloat(float32(movement.Y-baseGravity/16.0)) < 0.003 {
			yd = -0.003
		} else {
			yd = movement.Y - baseGravity/16.0
		}
		return phys.NewVec3(movement.X, yd, movement.Z)
	}
	return movement
}

// TravelFlying mirrors LivingEntity.travelFlying(Vec3, float).
func (e *Entity) TravelFlying(input phys.Vec3, speed float32) {
	e.TravelFlyingFull(input, 0.02, 0.02, speed)
}

// TravelFlyingFull mirrors LivingEntity.travelFlying(Vec3, float, float, float).
func (e *Entity) TravelFlyingFull(input phys.Vec3, waterSpeed float32, lavaSpeed float32, airSpeed float32) {
	switch {
	case e.IsInWater():
		e.MoveRelative(waterSpeed, input)
		e.Move(MoverTypeSELF, e.GetDeltaMovement())
		e.SetDeltaMovement(e.GetDeltaMovement().Scale(0.8))
	case e.IsInLava():
		e.MoveRelative(lavaSpeed, input)
		e.Move(MoverTypeSELF, e.GetDeltaMovement())
		e.SetDeltaMovement(e.GetDeltaMovement().Scale(0.5))
	default:
		e.MoveRelative(airSpeed, input)
		e.Move(MoverTypeSELF, e.GetDeltaMovement())
		e.SetDeltaMovement(e.GetDeltaMovement().Scale(0.91))
	}
}

// TravelFallFlying mirrors LivingEntity.travelFallFlying(Vec3).
func (e *Entity) TravelFallFlying(input phys.Vec3) {
	if e.OnClimbable() {
		e.TravelInAir(input)
		e.SetFallFlying(false)
		return
	}
	lastMovement := e.GetDeltaMovement()
	e.SetDeltaMovement(e.UpdateFallFlyingMovement(lastMovement))
	e.Move(MoverTypeSELF, e.GetDeltaMovement())
}

// UpdateFallFlyingMovement mirrors LivingEntity.updateFallFlyingMovement(Vec3).
func (e *Entity) UpdateFallFlyingMovement(movement phys.Vec3) (ret phys.Vec3) {
	lookAngle := e.GetLookAngle()
	leanAngle := float64(e.xRot * util.MthDegToRad)
	lookHorLength := math.Sqrt(lookAngle.X*lookAngle.X + lookAngle.Z*lookAngle.Z)
	moveHorLength := movement.HorizontalDistance()
	gravity := e.GetEffectiveGravity()
	liftForce := util.SquareDouble(math.Cos(leanAngle))
	movement = movement.Add(phys.NewVec3(0.0, gravity*(-1.0+liftForce*0.75), 0.0))
	if movement.Y < 0.0 && lookHorLength > 0.0 {
		convert := movement.Y * -0.1 * liftForce
		movement = movement.Add(phys.NewVec3(lookAngle.X*convert/lookHorLength, convert, lookAngle.Z*convert/lookHorLength))
	}
	if leanAngle < 0.0 && lookHorLength > 0.0 {
		convert := moveHorLength * -float64(util.Sin(leanAngle)) * 0.04
		movement = movement.Add(phys.NewVec3(-lookAngle.X*convert/lookHorLength, convert*3.2, -lookAngle.Z*convert/lookHorLength))
	}
	if lookHorLength > 0.0 {
		movement = movement.Add(phys.NewVec3(
			(lookAngle.X/lookHorLength*moveHorLength-movement.X)*0.1,
			0.0,
			(lookAngle.Z/lookHorLength*moveHorLength-movement.Z)*0.1,
		))
	}
	return movement.MultiplyXYZ(0.99, 0.98, 0.99)
}

// HandleRelativeFrictionAndCalculateMovement mirrors LivingEntity.handleRelativeFrictionAndCalculateMovement(...).
func (e *Entity) HandleRelativeFrictionAndCalculateMovement(input phys.Vec3, blockFriction float32) (ret phys.Vec3) {
	e.MoveRelative(e.GetFrictionInfluencedSpeed(blockFriction), input)
	e.SetDeltaMovement(e.HandleOnClimbable(e.GetDeltaMovement()))
	e.Move(MoverTypeSELF, e.GetDeltaMovement())
	movement := e.GetDeltaMovement()
	if (e.horizontalCollision || e.jumping) && e.OnClimbable() {
		movement = phys.NewVec3(movement.X, 0.2, movement.Z)
	}
	return movement
}

// HandleOnClimbable mirrors LivingEntity.handleOnClimbable(Vec3).
func (e *Entity) HandleOnClimbable(delta phys.Vec3) (ret phys.Vec3) {
	if !e.OnClimbable() {
		return delta
	}
	xd := util.ClampDouble(delta.X, -0.15, 0.15)
	zd := util.ClampDouble(delta.Z, -0.15, 0.15)
	yd := mathMax64(delta.Y, -0.15)
	if yd < 0.0 && e.suppressSlidingDownLadder {
		yd = 0.0
	}
	return phys.NewVec3(xd, yd, zd)
}

// GetFrictionInfluencedSpeed mirrors LivingEntity.getFrictionInfluencedSpeed(float).
func (e *Entity) GetFrictionInfluencedSpeed(blockFriction float32) (ret float32) {
	if e.OnGround() {
		if blockFriction > 0.6 {
			return e.GetSpeed() * (0.21600002 / (blockFriction * blockFriction * blockFriction))
		}
		return e.GetSpeed()
	}
	return e.GetFlyingSpeed()
}

// GetSpeed mirrors LivingEntity.getSpeed().
func (e *Entity) GetSpeed() (ret float32) {
	return float32(e.AttributeValue(AttributeMOVEMENT_SPEED)) * e.speedMultiplierOrDefault()
}

// GetFlyingSpeed mirrors LivingEntity.getFlyingSpeed().
func (e *Entity) GetFlyingSpeed() (ret float32) { return 0.02 }

func (e *Entity) speedMultiplierOrDefault() (ret float32) {
	if e.speedMultiplier == 0.0 {
		return 1.0
	}
	return e.speedMultiplier
}

// SetSpeedMultiplier sets the walk-speed multiplier (LivingEntity.setSpeed).
func (e *Entity) SetSpeedMultiplier(multiplier float32) { e.speedMultiplier = multiplier }

// SetSprinting mirrors LivingEntity.setSprinting for the movement slice.
func (e *Entity) SetSprinting(sprinting bool) { e.sprinting = sprinting }

// IsSprinting mirrors Entity.isSprinting().
func (e *Entity) IsSprinting() (ret bool) { return e.sprinting }

// IsJumping mirrors LivingEntity.isJumping().
func (e *Entity) IsJumping() (ret bool) { return e.jumping }

// IsFree mirrors Entity.isFree(double, double, double).
func (e *Entity) IsFree(dx float64, dy float64, dz float64) (ret bool) { return e.isFree(dx, dy, dz) }

func (e *Entity) isFree(dx float64, dy float64, dz float64) (ret bool) {
	return e.canNoCollision(e.BoundingBox.MoveXYZ(dx, dy, dz).Deflate(1.0e-7))
}

func (e *Entity) canNoCollision(box phys.AABB) (ret bool) {
	if ctx := e.MovementContext(); ctx != nil {
		return ctx.NoCollision(e, box)
	}
	return len(e.level.GetBlockCollisions(e, box)) == 0
}

func minFloat64(a float64, b float64) (ret float64) {
	if a < b {
		return a
	}
	return b
}

func mathMax64(a float64, b float64) (ret float64) {
	if a > b {
		return a
	}
	return b
}
