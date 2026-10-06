package control

import (
	"math"
	"sync"
	"time"

	"github.com/admin-else/strom/mc/entity"
	"github.com/admin-else/strom/mc/phys"
)

// Player constants.
const (
	tickInterval = 50 * time.Millisecond
	moveSpeed    = 0.1
	gravity      = 0.08
	eyeHeight    = 1.62
	playerWidth  = 0.6
	playerHeight = 1.8
	// resyncDistance is how far the local prediction may drift from the server
	// position before the controller snaps to the server (teleport/respawn).
	resyncDistance = 2.0
)

type lookAtAction struct {
	x, y, z    float64
	delayTicks int
}

// PlayerTarget is the server-bound movement sink the controller drives:
// net.minecraft.client.player.LocalPlayer's position/rotation accessors.
type PlayerTarget interface {
	Position() (x float64, y float64, z float64)
	Rotation() (yaw float32, pitch float32)
	SetPosition(x float64, y float64, z float64)
	SetRotation(yaw float32, pitch float32)
}

// Controller runs a local player entity with the ported physics and pushes the
// resulting position and rotation to the server-bound player every tick.
type Controller struct {
	mu     sync.Mutex
	level  entity.CollisionGetter
	player PlayerTarget
	ent    *entity.Entity

	yaw   float32
	pitch float32

	lookAt        *lookAtAction
	moveInput     entity.Input
	moveTicksLeft int
	jumpTicksLeft int
	sneak         bool
	sprint        bool
}

// New creates a controller for the given collision level and player module.
func New(level entity.CollisionGetter, playerMod PlayerTarget) (ret *Controller) {
	x, y, z := playerMod.Position()
	yaw, pitch := playerMod.Rotation()
	ent := entity.NewEntity(level, x, y, z, playerWidth, playerHeight)
	ent.SetRotation(yaw, pitch)
	return &Controller{level: level, player: playerMod, ent: ent, yaw: yaw, pitch: pitch}
}

// Level returns the collision source backing the controller.
func (c *Controller) Level() (ret entity.CollisionGetter) { return c.level }

// Entity returns the local simulated entity.
func (c *Controller) Entity() (ret *entity.Entity) { return c.ent }

// Position returns the local simulation position.
func (c *Controller) Position() (x, y, z float64) { return c.ent.X, c.ent.Y, c.ent.Z }

// Rotation returns the local yaw and pitch.
func (c *Controller) Rotation() (yaw, pitch float32) { return c.yaw, c.pitch }

// LookAt schedules a rotation toward the target after delay. This is the
// "look here in X time" action.
func (c *Controller) LookAt(x float64, y float64, z float64, delay time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lookAt = &lookAtAction{x: x, y: y, z: z, delayTicks: durationTicks(delay)}
}

// Look sets the rotation immediately.
func (c *Controller) Look(yaw float32, pitch float32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.yaw, c.pitch = yaw, pitch
}

// Move holds the given keys (a subset of wasd, plus 'j' to jump) for duration.
// This is the "press wasd for X time" action.
func (c *Controller) Move(keys string, duration time.Duration) (unknown []rune) {
	c.mu.Lock()
	defer c.mu.Unlock()
	in := entity.Input{}
	for _, key := range keys {
		switch key {
		case 'w', 'W':
			in.Forward = true
		case 's', 'S':
			in.Backward = true
		case 'a', 'A':
			in.Left = true
		case 'd', 'D':
			in.Right = true
		case 'j', 'J':
			in.Jump = true
		default:
			unknown = append(unknown, key)
		}
	}
	c.moveInput = in
	c.moveTicksLeft = durationTicks(duration)
	return
}

// Jump keeps the jump key pressed for duration. This is the "jump X seconds"
// action.
func (c *Controller) Jump(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.jumpTicksLeft = durationTicks(duration)
}

// SetSneak toggles crouch movement.
func (c *Controller) SetSneak(sneak bool) { c.mu.Lock(); c.sneak = sneak; c.mu.Unlock() }

// SetSprint toggles sprinting.
func (c *Controller) SetSprint(sprint bool) { c.mu.Lock(); c.sprint = sprint; c.mu.Unlock() }

// Stop clears every active action.
func (c *Controller) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lookAt = nil
	c.moveTicksLeft = 0
	c.jumpTicksLeft = 0
	c.moveInput = entity.Input{}
}

// Tick advances one 20 TPS step: resync to a server teleport, apply the active
// actions, run the entity physics and push the result to the player module.
func (c *Controller) Tick() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	sx, sy, sz := c.player.Position()
	if math.Abs(c.ent.X-sx) > resyncDistance || math.Abs(c.ent.Y-sy) > resyncDistance || math.Abs(c.ent.Z-sz) > resyncDistance {
		c.ent.SetPos(sx, sy, sz)
	}

	if c.lookAt != nil {
		if c.lookAt.delayTicks > 0 {
			c.lookAt.delayTicks--
		} else {
			eye := phys.NewVec3(c.ent.X, c.ent.Y+eyeHeight, c.ent.Z)
			target := phys.NewVec3(c.lookAt.x, c.lookAt.y, c.lookAt.z)
			rotation := target.Subtract(eye).Normalize().Rotation()
			c.yaw = rotation.Y
			c.pitch = rotation.X
			c.lookAt = nil
		}
	}

	input := entity.Input{}
	if c.moveTicksLeft > 0 {
		input = c.moveInput
		c.moveTicksLeft--
	}
	if c.jumpTicksLeft > 0 {
		input.Jump = true
		c.jumpTicksLeft--
	}
	input.Shift = c.sneak
	input.Sprint = c.sprint

	c.ent.SetRotation(c.yaw, c.pitch)
	c.ent.TravelWithInput(input, input.Shift)

	c.player.SetPosition(c.ent.X, c.ent.Y, c.ent.Z)
	c.player.SetRotation(c.yaw, c.pitch)
	return nil
}

func durationTicks(duration time.Duration) (ret int) {
	if duration <= 0 {
		return 0
	}
	return int(duration / tickInterval)
}
