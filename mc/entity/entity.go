// Package entity mirrors the collision and movement slice of
// net.minecraft.world.entity.Entity / LivingEntity / Player that a headless
// client or bot needs. Rendering-only members live in neon.
//
// The world is abstracted behind CollisionGetter so the package stays headless:
// neon (or a test) supplies the block and entity collision shapes for a box.
// Deferred, not faked: fluids, climbing, elytra/fall-flying, attributes, mob
// effects, fall damage, world border, pistons, and the profiler/restitution
// bookkeeping in Entity.move.
package entity

import (
	"math"
	"sort"

	"github.com/admin-else/strom/mc/core"
	"github.com/admin-else/strom/mc/phys"
	"github.com/admin-else/strom/mc/phys/shapes"
	"github.com/admin-else/strom/mc/util"
)

// CollisionGetter is the slice of net.minecraft.world.level.Level that
// Entity.collide consumes: the entity collision shapes and the block collision
// shapes intersecting a box.
type CollisionGetter interface {
	GetEntityCollisions(source *Entity, box phys.AABB) (ret []*shapes.VoxelShape)
	GetBlockCollisions(source *Entity, box phys.AABB) (ret []*shapes.VoxelShape)
}

// Entity mirrors the movement state of net.minecraft.world.entity.Entity.
type Entity struct {
	X float64
	Y float64
	Z float64

	yRot float32
	xRot float32

	deltaMovement phys.Vec3

	BoundingBox phys.AABB

	onGround                 bool
	horizontalCollision      bool
	verticalCollision        bool
	verticalCollisionBelow   bool
	minorHorizontalCollision bool

	NoPhysics bool
	MaxUpStep float64

	Width  float64
	Height float64

	jumping bool
	level   CollisionGetter
}

// NewEntity mirrors Entity(EntityType, Level) for the collision slice.
func NewEntity(level CollisionGetter, x float64, y float64, z float64, width float64, height float64) (ret *Entity) {
	ret = &Entity{
		X: x, Y: y, Z: z,
		Width: width, Height: height,
		MaxUpStep: 0.6,
		level:     level,
	}
	ret.RecomputeBoundingBox()
	return
}

// GetX mirrors Entity.getX().
func (e *Entity) GetX() (ret float64) { return e.X }

// GetY mirrors Entity.getY().
func (e *Entity) GetY() (ret float64) { return e.Y }

// GetZ mirrors Entity.getZ().
func (e *Entity) GetZ() (ret float64) { return e.Z }

// Position mirrors Entity.position().
func (e *Entity) Position() (ret phys.Vec3) { return phys.NewVec3(e.X, e.Y, e.Z) }

// GetBoundingBox mirrors Entity.getBoundingBox().
func (e *Entity) GetBoundingBox() (ret phys.AABB) { return e.BoundingBox }

// RecomputeBoundingBox mirrors Entity.recomputeBoundingBox() for the default
// (non-posing) dimensions.
func (e *Entity) RecomputeBoundingBox() {
	half := e.Width / 2.0
	e.BoundingBox = phys.NewAABB(e.X-half, e.Y, e.Z-half, e.X+half, e.Y+e.Height, e.Z+half)
}

// SetPos mirrors Entity.setPos(double, double, double).
func (e *Entity) SetPos(x float64, y float64, z float64) {
	e.X, e.Y, e.Z = x, y, z
	e.RecomputeBoundingBox()
}

// GetDeltaMovement mirrors Entity.getDeltaMovement().
func (e *Entity) GetDeltaMovement() (ret phys.Vec3) { return e.deltaMovement }

// SetDeltaMovement mirrors Entity.setDeltaMovement(Vec3).
func (e *Entity) SetDeltaMovement(delta phys.Vec3) { e.deltaMovement = delta }

// AddDeltaMovement mirrors Entity.addDeltaMovement(Vec3).
func (e *Entity) AddDeltaMovement(movement phys.Vec3) {
	e.deltaMovement = e.deltaMovement.Add(movement)
}

// OnGround mirrors Entity.onGround().
func (e *Entity) OnGround() (ret bool) { return e.onGround }

// HorizontalCollision mirrors Entity.horizontalCollision.
func (e *Entity) HorizontalCollision() (ret bool) { return e.horizontalCollision }

// VerticalCollision mirrors Entity.verticalCollision.
func (e *Entity) VerticalCollision() (ret bool) { return e.verticalCollision }

// SetJumping mirrors Entity.setJumping for the movement slice.
func (e *Entity) SetJumping(jumping bool) { e.jumping = jumping }

// SetRotation sets the yaw/pitch used by moveRelative.
func (e *Entity) SetRotation(yRot float32, xRot float32) { e.yRot = yRot; e.xRot = xRot }

// GetYRot mirrors Entity.getYRot().
func (e *Entity) GetYRot() (ret float32) { return e.yRot }

// Move mirrors Entity.move(MoverType, Vec3). Fall damage, restitution, piston
// limits and the profiler are deferred.
func (e *Entity) Move(moverType MoverType, delta phys.Vec3) {
	if e.NoPhysics {
		e.SetPos(e.X+delta.X, e.Y+delta.Y, e.Z+delta.Z)
		e.horizontalCollision = false
		e.verticalCollision = false
		e.verticalCollisionBelow = false
		return
	}

	delta = e.MaybeBackOffFromEdge(delta, moverType)
	movement := e.Collide(delta)
	movementLength := movement.LengthSqr()
	if movementLength > 1.0e-7 || delta.LengthSqr()-movementLength < 1.0e-7 {
		pos := e.Position()
		e.SetPos(pos.X+movement.X, pos.Y+movement.Y, pos.Z+movement.Z)
	}

	xCollision := !mthEqual(delta.X, movement.X)
	zCollision := !mthEqual(delta.Z, movement.Z)
	e.horizontalCollision = xCollision || zCollision
	movedVertically := math.Abs(delta.Y) > 0.0
	if movedVertically || e.IsLocalInstanceAuthoritative() {
		e.verticalCollision = delta.Y != movement.Y
		e.verticalCollisionBelow = e.verticalCollision && delta.Y < 0.0
		e.onGround = e.verticalCollisionBelow
	}
	if e.horizontalCollision {
		e.minorHorizontalCollision = false
	} else {
		e.minorHorizontalCollision = false
	}
}

// IsLocalInstanceAuthoritative mirrors Entity.isLocalInstanceAuthoritative().
// A headless client drives its own player, so it is authoritative.
func (e *Entity) IsLocalInstanceAuthoritative() (ret bool) { return true }

// MaybeBackOffFromEdge mirrors Entity.maybeBackOffFromEdge(Vec3, MoverType). The
// base implementation returns the delta unchanged; sneaking edge protection is a
// LivingEntity/Player override.
func (e *Entity) MaybeBackOffFromEdge(delta phys.Vec3, moverType MoverType) (ret phys.Vec3) {
	return delta
}

// Collide mirrors Entity.collide(Vec3), including the maxUpStep step-up search.
func (e *Entity) Collide(movement phys.Vec3) (ret phys.Vec3) {
	aabb := e.GetBoundingBox()
	entityColliders := e.level.GetEntityCollisions(e, aabb.ExpandTowards(movement).ExpandTowardsXYZ(0.0, e.MaxUpStep, 0.0))
	var movementStep phys.Vec3
	if movement.LengthSqr() == 0.0 {
		movementStep = movement
	} else {
		movementStep = CollideBoundingBox(e, movement, aabb, e.level, entityColliders)
	}
	xCollision := movement.X != movementStep.X
	yCollision := movement.Y != movementStep.Y
	zCollision := movement.Z != movementStep.Z
	onGroundAfterCollision := yCollision && movement.Y < 0.0
	if e.MaxUpStep > 0.0 && (onGroundAfterCollision || e.onGround) && (xCollision || zCollision) {
		groundedAABB := aabb
		if onGroundAfterCollision {
			groundedAABB = aabb.MoveXYZ(0.0, movementStep.Y, 0.0)
		}
		stepUpAABB := groundedAABB.ExpandTowardsXYZ(movement.X, e.MaxUpStep, movement.Z)
		if !onGroundAfterCollision {
			stepUpAABB = stepUpAABB.ExpandTowardsXYZ(0.0, -1.0e-5, 0.0)
		}
		colliders := collectCollidersIgnoringWorldBorder(e, e.level, entityColliders, stepUpAABB)
		stepHeightToSkip := float32(movementStep.Y)
		for _, candidateStepUpHeight := range CollectCandidateStepUpHeights(groundedAABB, colliders, float32(e.MaxUpStep), stepHeightToSkip) {
			stepFromGround := collideWithShapes(phys.NewVec3(movement.X, float64(candidateStepUpHeight), movement.Z), groundedAABB, colliders)
			if stepFromGround.HorizontalDistanceSqr() > movementStep.HorizontalDistanceSqr() {
				distanceToGround := aabb.MinY - groundedAABB.MinY
				return stepFromGround.SubtractXYZ(0.0, distanceToGround, 0.0)
			}
		}
	}
	return movementStep
}

// CollideBoundingBox mirrors Entity.collideBoundingBox(Entity, Vec3, AABB, Level, List<VoxelShape>).
func CollideBoundingBox(source *Entity, movement phys.Vec3, boundingBox phys.AABB, level CollisionGetter, entityColliders []*shapes.VoxelShape) (ret phys.Vec3) {
	colliders := collectCollidersIgnoringWorldBorder(source, level, entityColliders, boundingBox.ExpandTowards(movement))
	return collideWithShapes(movement, boundingBox, colliders)
}

func collectCollidersIgnoringWorldBorder(source *Entity, level CollisionGetter, entityColliders []*shapes.VoxelShape, boundingBox phys.AABB) (ret []*shapes.VoxelShape) {
	if len(entityColliders) > 0 {
		ret = append(ret, entityColliders...)
	}
	ret = append(ret, level.GetBlockCollisions(source, boundingBox)...)
	return
}

func collideWithShapes(movement phys.Vec3, boundingBox phys.AABB, shapeList []*shapes.VoxelShape) (ret phys.Vec3) {
	if len(shapeList) == 0 {
		return movement
	}
	resolvedMovement := phys.Vec3ZERO
	for _, axis := range phys.AxisStepOrder(movement) {
		axisMovement := movement.Get(axis)
		if axisMovement != 0.0 {
			collision := shapes.Collide(axis, boundingBox.MoveVec3(resolvedMovement), shapeList, axisMovement)
			resolvedMovement = resolvedMovement.With(axis, collision)
		}
	}
	return resolvedMovement
}

// CollectCandidateStepUpHeights mirrors Entity.collectCandidateStepUpHeights(...).
func CollectCandidateStepUpHeights(boundingBox phys.AABB, colliders []*shapes.VoxelShape, maxStepHeight float32, stepHeightToSkip float32) (ret []float32) {
	set := make(map[float32]struct{}, 4)
	for _, collider := range colliders {
		for _, coord := range collider.GetCoords(core.AxisY) {
			relativeCoord := float32(coord - boundingBox.MinY)
			if !(relativeCoord < 0.0) && relativeCoord != stepHeightToSkip {
				if relativeCoord > maxStepHeight {
					break
				}
				set[relativeCoord] = struct{}{}
			}
		}
	}
	ret = make([]float32, 0, len(set))
	for candidate := range set {
		ret = append(ret, candidate)
	}
	sort.Slice(ret, func(i int, j int) bool { return ret[i] < ret[j] })
	return
}

// MoveRelative mirrors Entity.moveRelative(float, Vec3).
func (e *Entity) MoveRelative(speed float32, input phys.Vec3) {
	delta := GetInputVector(input, speed, e.yRot)
	e.AddDeltaMovement(delta)
}

// GetInputVector mirrors Entity.getInputVector(Vec3, float, float).
func GetInputVector(input phys.Vec3, speed float32, yRot float32) (ret phys.Vec3) {
	length := input.LengthSqr()
	if length < 1.0e-7 {
		return phys.Vec3ZERO
	}
	movement := input
	if length > 1.0 {
		movement = input.Normalize()
	}
	movement = movement.Scale(float64(speed))
	sin := util.Sin(float64(yRot * util.MthDegToRad))
	cos := util.Cos(float64(yRot * util.MthDegToRad))
	return phys.NewVec3(
		movement.X*float64(cos)-movement.Z*float64(sin),
		movement.Y,
		movement.Z*float64(cos)+movement.X*float64(sin),
	)
}

// Travel mirrors the LivingEntity.travel ground/air path for a player without
// attributes, mob effects, fluids or climbing (all deferred). input is the
// (xxa, yya, zza) movement vector built from the key input.
func (e *Entity) Travel(input phys.Vec3, speed float32, gravity float64) {
	blockFriction := float32(1.0)
	if e.onGround {
		blockFriction = 0.6
	}
	movement := e.handleRelativeFrictionAndCalculateMovement(input, blockFriction, speed)
	movementY := movement.Y - gravity
	friction := blockFriction * 0.91
	verticalFriction := float32(0.98)
	e.SetDeltaMovement(phys.NewVec3(movement.X*float64(friction), movementY*float64(verticalFriction), movement.Z*float64(friction)))
}

func (e *Entity) handleRelativeFrictionAndCalculateMovement(input phys.Vec3, blockFriction float32, speed float32) (ret phys.Vec3) {
	e.MoveRelative(e.getFrictionInfluencedSpeed(blockFriction, speed), input)
	e.Move(MoverTypeSELF, e.GetDeltaMovement())
	return e.GetDeltaMovement()
}

func (e *Entity) getFrictionInfluencedSpeed(blockFriction float32, speed float32) (ret float32) {
	if e.onGround {
		if blockFriction > 0.6 {
			return speed * (0.21600002 / (blockFriction * blockFriction * blockFriction))
		}
		return speed
	}
	return 0.02
}

// JumpFromGround mirrors LivingEntity.jumpFromGround() with the default jump
// power (0.42).
func (e *Entity) JumpFromGround() {
	movement := e.GetDeltaMovement()
	e.SetDeltaMovement(phys.NewVec3(movement.X, 0.42, movement.Z))
}

func mthEqual(a float64, b float64) (ret bool) { return math.Abs(b-a) < float64(float32(1.0e-5)) }
