// Package entity tracks the client-side entities a bot can see, fed by the
// 26.4 play packets that add, move, rotate and remove them. It is headless and
// reusable by a bot or by the renderer.
//
// The 26.4 entity movement packets are hand-written because the minecraft-data
// schema for 26.4-snapshot-2 (derived from 26.2) does not match the wire format;
// see packets.go. Mod.Start replaces the stale generated types for those packet
// ids before the connection loop starts.
package entity

import (
	"math"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
	"github.com/admin-else/strom/mc/util"
)

// VecDeltaCodec mirrors net.minecraft.network.protocol.game.VecDeltaCodec.
// Entities keep one codec whose base advances with each movement packet, so
// short deltas can be resolved to absolute positions.
type VecDeltaCodec struct {
	base Vec3
}

// truncationSteps mirrors VecDeltaCodec.TRUNCATION_STEPS.
const truncationSteps = 4096.0

func encodeDelta(input float64) (ret int64) {
	return int64(math.Round(input * truncationSteps))
}

func decodeDelta(v int64) (ret float64) {
	return float64(v) / truncationSteps
}

// SetBase mirrors VecDeltaCodec.setBase.
func (c *VecDeltaCodec) SetBase(base Vec3) { c.base = base }

// Base mirrors VecDeltaCodec.getBase.
func (c *VecDeltaCodec) Base() Vec3 { return c.base }

// Decode mirrors VecDeltaCodec.decode(long, long, long).
func (c *VecDeltaCodec) Decode(xa, ya, za int64) (ret Vec3) {
	if xa == 0 && ya == 0 && za == 0 {
		return c.base
	}
	ret.X = c.base.X
	if xa != 0 {
		ret.X = decodeDelta(encodeDelta(c.base.X) + xa)
	}
	ret.Y = c.base.Y
	if ya != 0 {
		ret.Y = decodeDelta(encodeDelta(c.base.Y) + ya)
	}
	ret.Z = c.base.Z
	if za != 0 {
		ret.Z = decodeDelta(encodeDelta(c.base.Z) + za)
	}
	return
}

// decodePath mirrors VecDelta.decode(VecDeltaCodec). It restores the codec base
// after a stepped path, matching the Java implementation; the caller is
// responsible for setting the base to the path end afterwards.
func (c *VecDeltaCodec) decodePath(delta VecDelta) (ret PositionPath) {
	if !delta.Stepped {
		return PositionPath{EndPosition: c.Decode(int64(delta.Xa), int64(delta.Ya), int64(delta.Za))}
	}
	if len(delta.Steps) == 0 {
		return PositionPath{Stepped: true, EndPosition: c.base}
	}
	originalBase := c.base
	steps := make([]PositionStep, 0, len(delta.Steps))
	for _, step := range delta.Steps {
		pos := c.Decode(int64(step.Xa), int64(step.Ya), int64(step.Za))
		steps = append(steps, PositionStep{Position: pos, TickOffset: step.Ticks})
		c.base = pos
	}
	c.base = originalBase
	return PositionPath{Stepped: true, EndPosition: steps[len(steps)-1].Position, Steps: steps}
}

// Entity is the client-side projection of a tracked entity.
//
// position/rotation are the live (interpolated) values, mirroring Entity's
// x/y/z, yRot, xRot and yHeadRot. oldPosition/oldRotation are snapshotted at the
// start of every client tick, mirroring xo/yo/zo, yRotO, xRotO and yHeadRotO;
// the renderer lerps between the snapshot and the live value by the frame's
// partial tick (EntityRenderer.extractRenderState). tickCount mirrors
// Entity.tickCount.
type Entity struct {
	Id       int32
	UUID     uuid.UUID
	Type     int32
	Velocity Vec3
	OnGround bool

	position    Vec3
	yaw         float32
	pitch       float32
	headYaw     float32
	oldPosition Vec3
	oldYaw      float32
	oldPitch    float32
	oldHeadYaw  float32
	tickCount   int

	metadata     map[uint8]any
	equipment    [equipmentSlotCount]ItemStack
	equipmentSet [equipmentSlotCount]bool

	codec         VecDeltaCodec
	interpolation InterpolationHandler
}

// defaultInterpolationSteps mirrors EntityType.Builder's updateInterval default
// (3), the value SteppedInterpolationHandler.create reads from
// entity.getType().updateInterval().
const defaultInterpolationSteps = 3

// newEntity builds a tracked entity with its interpolation handler wired,
// mirroring Entity's constructor chain. The reduced port models every tracked
// entity as a LivingEntity: it always uses the Stepped handler with the default
// updateInterval (3), so the per-type handler selection (Display/minecart/boat
// use the Linear handler, some types are NO_OP) and per-type updateInterval are
// not reproduced.
func newEntity() (e *Entity) {
	e = &Entity{}
	e.interpolation = NewSteppedInterpolationHandler(e, defaultInterpolationSteps)
	return
}

// Position mirrors Entity.getPosition(float): the previous-tick position lerped
// toward the live position by partialTick.
func (e *Entity) Position(partialTick float32) (ret Vec3) {
	a := float64(partialTick)
	ret.X = util.LerpDouble(a, e.oldPosition.X, e.position.X)
	ret.Y = util.LerpDouble(a, e.oldPosition.Y, e.position.Y)
	ret.Z = util.LerpDouble(a, e.oldPosition.Z, e.position.Z)
	return
}

// Yaw mirrors Entity.getYRot(float).
func (e *Entity) Yaw(partialTick float32) (ret float32) {
	if partialTick == 1.0 {
		return e.yaw
	}
	return util.RotLerpFloat(partialTick, e.oldYaw, e.yaw)
}

// Pitch mirrors Entity.getXRot(float).
func (e *Entity) Pitch(partialTick float32) (ret float32) {
	if partialTick == 1.0 {
		return e.pitch
	}
	return util.LerpFloat(partialTick, e.oldPitch, e.pitch)
}

// HeadYaw mirrors LivingEntity.getYHeadRot(float).
func (e *Entity) HeadYaw(partialTick float32) (ret float32) {
	if partialTick == 1.0 {
		return e.headYaw
	}
	return util.RotLerpFloat(partialTick, e.oldHeadYaw, e.headYaw)
}

// AgeInTicks mirrors EntityRenderState.ageInTicks = tickCount + partialTicks.
func (e *Entity) AgeInTicks(partialTick float32) (ret float32) {
	return float32(e.tickCount) + partialTick
}

// StorePositionAndRotation mirrors Entity.storePositionAndRotation().
func (e *Entity) StorePositionAndRotation() (ret PositionAndRotation) {
	return PositionAndRotation{Position: e.position, YRot: e.yaw, XRot: e.pitch}
}

// SetPositionAndRotation mirrors the setPos + setRot pair the interpolation
// handlers apply each tick. setRot wraps yRot and clamps xRot, mirroring
// Entity.setRot -> setYRot(yRot % 360) / setXRot(clamp(xRot % 360, -90, 90)).
func (e *Entity) SetPositionAndRotation(position Vec3, yRot, xRot float32) {
	e.position = position
	e.yaw = float32(math.Mod(float64(yRot), 360.0))
	e.pitch = util.ClampFloat(float32(math.Mod(float64(xRot), 360.0)), -90.0, 90.0)
}

// SnapTo mirrors Entity.snapTo(double, double, double, float, float): set the
// live pose and collapse the previous-tick snapshot onto it.
func (e *Entity) SnapTo(position Vec3, yRot, xRot float32) {
	e.SetPositionAndRotation(position, yRot, xRot)
	e.oldPosition = position
	e.oldYaw = e.yaw
	e.oldPitch = e.pitch
	e.oldHeadYaw = e.headYaw
}

// onInterpolationStart mirrors Entity.onInterpolationStart; subclasses that
// react to a new interpolation target override it (none does in the port).
func (e *Entity) onInterpolationStart() {}

// Tick mirrors Entity.commonTick's client side: snapshot the old pose, advance
// the interpolation one step, then bump tickCount.
func (e *Entity) Tick() {
	e.oldPosition = e.position
	e.oldYaw = e.yaw
	e.oldPitch = e.pitch
	e.oldHeadYaw = e.headYaw
	if e.interpolation != nil {
		e.interpolation.Interpolate()
	}
	e.tickCount++
}

// MoveOrInterpolateTo mirrors Entity.moveOrInterpolateTo(PositionPath, float, float).
func (e *Entity) MoveOrInterpolateTo(position *PositionPath, yRot, xRot float32) {
	if e.interpolation == nil {
		if position != nil {
			e.position = position.EndPosition
		}
		e.yaw = yRot
		e.pitch = xRot
		return
	}
	if !e.interpolation.InterpolateTo(position, yRot, xRot, true) {
		if position != nil {
			e.position = position.EndPosition
		}
		e.yaw = yRot
		e.pitch = xRot
	}
}

// MoveOrInterpolateToPath mirrors Entity.moveOrInterpolateTo(PositionPath).
func (e *Entity) MoveOrInterpolateToPath(position *PositionPath) {
	if e.interpolation == nil {
		if position != nil {
			e.position = position.EndPosition
		}
		return
	}
	if !e.interpolation.InterpolateTo(position, 0, 0, false) {
		if position != nil {
			e.position = position.EndPosition
		}
	}
}

// MoveOrInterpolateToLinear mirrors Entity.moveOrInterpolateTo(Vec3, float, float).
func (e *Entity) MoveOrInterpolateToLinear(position Vec3, yRot, xRot float32) {
	e.MoveOrInterpolateTo(&PositionPath{EndPosition: position}, yRot, xRot)
}

// MoveOrInterpolateToRot mirrors Entity.moveOrInterpolateTo(float, float).
func (e *Entity) MoveOrInterpolateToRot(yRot, xRot float32) {
	e.MoveOrInterpolateTo(nil, yRot, xRot)
}

// unpackDegrees mirrors Mth.unpackDegrees(byte).
func unpackDegrees(value int8) (ret float32) {
	return float32(value) * (360.0 / 256.0)
}

// Module tracks entities on a connection.
type Module struct {
	*proto.Conn

	mu       sync.RWMutex
	entities map[int32]*Entity
	order    []int32

	ticker *event.Timer
}

// entityTickInterval is the vanilla client tick cadence (20 TPS), matching
// mc/bot/player's ticker.
const entityTickInterval = 50 * time.Millisecond

// NewModule creates an entity Module around an existing connection without
// registering handlers.
func NewModule(c *proto.Conn) (m *Module) {
	return &Module{
		Conn:     c,
		entities: make(map[int32]*Entity),
	}
}

// Entity returns the tracked entity with the given id.
func (m *Module) Entity(id int32) (entity *Entity, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entity, ok = m.entities[id]
	return
}

// Entities returns a snapshot of the tracked entities in spawn order.
func (m *Module) Entities() (ret []*Entity) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret = make([]*Entity, 0, len(m.order))
	for _, id := range m.order {
		if entity, ok := m.entities[id]; ok {
			ret = append(ret, entity)
		}
	}
	return
}

// Count returns the number of tracked entities.
func (m *Module) Count() (n int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entities)
}

// steppedEntityMovementFirstVersion is the first version whose entity movement
// packets use the stepped VecDelta wire form (a properties varint encoding the
// on-ground bit and step count). Earlier versions still send short deltas.
const steppedEntityMovementFirstVersion = "26.4-snapshot-2"

// usesSteppedEntityMovement reports whether the version's entity movement
// packets use the stepped VecDelta form.
func usesSteppedEntityMovement(version string) (ret bool) {
	i := slices.Index(proto_generated.SupportedVersions, version)
	j := slices.Index(proto_generated.SupportedVersions, steppedEntityMovementFirstVersion)
	if i == -1 || j == -1 {
		return version == steppedEntityMovementFirstVersion
	}
	return i >= j
}

// Start registers the entity handlers on c. For 26.4-snapshot-2 the generated
// movement decoders are stale, so the packet ids are replaced with the
// hand-written stepped form. Older versions keep their generated decoders
// (short deltas) and use the adapter handlers below.
func Start(c *proto.Conn) (m *Module) {
	m = NewModule(c)
	if usesSteppedEntityMovement(c.Version) {
		// The overridden packet ids are resolved by name so this stays correct
		// if the generated id table changes.
		override := func(name string, packetType proto_base.EncodeDecodeAble) {
			info, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(name, c.ProtocolVersion, proto_base.Play, proto_base.ToClient)
			if !ok {
				return
			}
			proto.OverridePacketType(proto_base.ToClient, proto_base.Play, info.PacketId, c.ProtocolVersion, name, packetType)
		}
		override("rel_entity_move", &PlayToClientPacketMoveEntityPos{})
		override("entity_move_look", &PlayToClientPacketMoveEntityPosRot{})
		override("entity_look", &PlayToClientPacketMoveEntityRot{})
		override("sync_entity_position", &PlayToClientPacketEntityPositionSync{})
		override("entity_teleport", &PlayToClientPacketTeleportEntity{})

		m.Register(m.onMoveEntityPos)
		m.Register(m.onMoveEntityPosRot)
		m.Register(m.onMoveEntityRot)
		m.Register(m.onEntityPositionSync)
		m.Register(m.onTeleportEntity)
	} else {
		m.RegisterUntil("26.2",
			m.onRelEntityMoveOld,
			m.onEntityMoveLookOld,
			m.onEntityLookOld,
			m.onSyncEntityPositionOld,
			m.onTeleportEntityOld,
		)
	}

	m.Register(m.onSpawnEntity)
	m.Register(m.onEntityDestroy)
	m.Register(m.onEntityHeadRotation)
	m.Register(m.onEntityVelocity)
	m.Register(m.onEntityMetadata)
	m.Register(m.onEntityEquipment)
	m.startTick()
	return
}

// startTick advances every tracked entity's interpolation at the vanilla client
// tick cadence, mirroring Minecraft's ClientLevel.tickEntities -> entity.tick()
// -> commonTick (setOldPosAndRot + interpolate + tickCount++). player.Module
// registers its ticker the same way.
func (m *Module) startTick() {
	m.ticker = &event.Timer{}
	m.ticker.Every(entityTickInterval, m.tick)
	m.ticker.Start(m.Conn.Loop)
}

func (m *Module) tick() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range m.order {
		if entity, ok := m.entities[id]; ok {
			entity.Tick()
		}
	}
	return nil
}

func (m *Module) add(entity *Entity) {
	if _, exists := m.entities[entity.Id]; !exists {
		m.order = append(m.order, entity.Id)
	}
	m.entities[entity.Id] = entity
}

func (m *Module) onSpawnEntity(p *v26_4_snapshot_2.PlayToClientPacketSpawnEntity) (err error) {
	entity := newEntity()
	entity.Id = p.EntityId
	entity.UUID = p.ObjectUUID
	entity.Type = p.Type
	entity.position = Vec3{p.X, p.Y, p.Z}
	entity.yaw = unpackDegrees(p.Yaw)
	entity.pitch = unpackDegrees(p.Pitch)
	entity.headYaw = unpackDegrees(p.HeadPitch)
	entity.Velocity = Vec3{p.Velocity.X, p.Velocity.Y, p.Velocity.Z}
	entity.oldPosition = entity.position
	entity.oldYaw = entity.yaw
	entity.oldPitch = entity.pitch
	entity.oldHeadYaw = entity.headYaw
	entity.codec.SetBase(entity.position)
	m.mu.Lock()
	m.add(entity)
	m.mu.Unlock()
	return nil
}

func (m *Module) onEntityDestroy(p *v26_4_snapshot_2.PlayToClientPacketEntityDestroy) (err error) {
	m.mu.Lock()
	for _, id := range p.EntityIds {
		if _, ok := m.entities[id]; !ok {
			continue
		}
		delete(m.entities, id)
		for i, other := range m.order {
			if other == id {
				m.order = append(m.order[:i], m.order[i+1:]...)
				break
			}
		}
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onEntityHeadRotation(p *v26_4_snapshot_2.PlayToClientPacketEntityHeadRotation) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		entity.headYaw = unpackDegrees(p.HeadYaw)
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onEntityVelocity(p *v26_4_snapshot_2.PlayToClientPacketEntityVelocity) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		entity.Velocity = Vec3{p.Velocity.X, p.Velocity.Y, p.Velocity.Z}
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onMoveEntityPos(p *PlayToClientPacketMoveEntityPos) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		path := entity.codec.decodePath(p.Delta)
		entity.codec.SetBase(path.EndPosition)
		entity.MoveOrInterpolateToPath(&path)
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onMoveEntityPosRot(p *PlayToClientPacketMoveEntityPosRot) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		yRot := unpackDegrees(p.YRot)
		xRot := unpackDegrees(p.XRot)
		path := entity.codec.decodePath(p.Delta)
		entity.codec.SetBase(path.EndPosition)
		entity.MoveOrInterpolateTo(&path, yRot, xRot)
		// The reduced port has no LivingEntity.aiStep head-turn solver, so the
		// head follows the body yaw on a look packet; explicit
		// ClientboundRotateHead packets still override it via headYaw.
		entity.headYaw = yRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onMoveEntityRot(p *PlayToClientPacketMoveEntityRot) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		yRot := unpackDegrees(p.YRot)
		entity.MoveOrInterpolateToRot(yRot, unpackDegrees(p.XRot))
		entity.headYaw = yRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

// The *Old handlers adapt the pre-26.4 generated movement packets (short deltas
// and a boolean on-ground flag) to the same entity update logic.

func (m *Module) onRelEntityMoveOld(p *v1_21_11.PlayToClientPacketRelEntityMove) (err error) {
	return m.onMoveEntityPos(&PlayToClientPacketMoveEntityPos{
		EntityId: p.EntityId,
		Delta:    VecDelta{Xa: p.DX, Ya: p.DY, Za: p.DZ},
		OnGround: p.OnGround,
	})
}

func (m *Module) onEntityMoveLookOld(p *v1_21_11.PlayToClientPacketEntityMoveLook) (err error) {
	return m.onMoveEntityPosRot(&PlayToClientPacketMoveEntityPosRot{
		EntityId: p.EntityId,
		Delta:    VecDelta{Xa: p.DX, Ya: p.DY, Za: p.DZ},
		YRot:     p.Yaw,
		XRot:     p.Pitch,
		OnGround: p.OnGround,
	})
}

func (m *Module) onEntityLookOld(p *v1_21_11.PlayToClientPacketEntityLook) (err error) {
	return m.onMoveEntityRot(&PlayToClientPacketMoveEntityRot{
		EntityId: p.EntityId,
		YRot:     p.Yaw,
		XRot:     p.Pitch,
		OnGround: p.OnGround,
	})
}

func (m *Module) onSyncEntityPositionOld(p *v1_21_11.PlayToClientPacketSyncEntityPosition) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		position := Vec3{p.X, p.Y, p.Z}
		entity.codec.SetBase(position)
		entity.Velocity = Vec3{p.Dx, p.Dy, p.Dz}
		entity.MoveOrInterpolateToLinear(position, p.Yaw, p.Pitch)
		entity.headYaw = p.Yaw
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onTeleportEntityOld(p *v1_21_11.PlayToClientPacketEntityTeleport) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		position := Vec3{p.X, p.Y, p.Z}
		entity.codec.SetBase(position)
		yRot := unpackDegrees(p.Yaw)
		entity.MoveOrInterpolateToLinear(position, yRot, unpackDegrees(p.Pitch))
		entity.headYaw = yRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onEntityPositionSync(p *PlayToClientPacketEntityPositionSync) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		entity.codec.SetBase(p.Position.EndPosition)
		entity.MoveOrInterpolateTo(&p.Position, p.YRot, p.XRot)
		entity.headYaw = p.YRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onTeleportEntity(p *PlayToClientPacketTeleportEntity) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		change := calculateAbsolute(PositionMoveRotation{
			Position:      entity.position,
			DeltaMovement: entity.Velocity,
			YRot:          entity.yaw,
			XRot:          entity.pitch,
		}, p.Change, p.Relatives)
		entity.codec.SetBase(change.Position)
		entity.Velocity = change.DeltaMovement
		entity.MoveOrInterpolateToLinear(change.Position, change.YRot, change.XRot)
		entity.headYaw = change.YRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

// calculateAbsolute mirrors net.minecraft.world.entity.PositionMoveRotation.calculateAbsolute.
func calculateAbsolute(source, change PositionMoveRotation, relatives int32) (ret PositionMoveRotation) {
	ret.Position.X = change.Position.X
	ret.Position.Y = change.Position.Y
	ret.Position.Z = change.Position.Z
	if relatives&RelativeX != 0 {
		ret.Position.X += source.Position.X
	}
	if relatives&RelativeY != 0 {
		ret.Position.Y += source.Position.Y
	}
	if relatives&RelativeZ != 0 {
		ret.Position.Z += source.Position.Z
	}
	ret.YRot = change.YRot
	if relatives&RelativeYRot != 0 {
		ret.YRot += source.YRot
	}
	ret.XRot = change.XRot
	if relatives&RelativeXRot != 0 {
		ret.XRot += source.XRot
	}
	ret.XRot = clamp(ret.XRot, -90, 90)
	ret.DeltaMovement = calculateDeltaVector(source.DeltaMovement, change.DeltaMovement, relatives)
	return
}

func calculateDeltaVector(current, change Vec3, relatives int32) (ret Vec3) {
	ret.X = calculateDelta(current.X, change.X, relatives, RelativeDeltaX)
	ret.Y = calculateDelta(current.Y, change.Y, relatives, RelativeDeltaY)
	ret.Z = calculateDelta(current.Z, change.Z, relatives, RelativeDeltaZ)
	return
}

func calculateDelta(currentDelta, deltaChange float64, relatives int32, relative int32) (ret float64) {
	if relatives&relative != 0 {
		return currentDelta + deltaChange
	}
	return deltaChange
}

func clamp(value, min, max float32) (ret float32) {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
