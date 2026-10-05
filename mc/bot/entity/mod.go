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
	"sync"

	"github.com/google/uuid"

	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
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
type Entity struct {
	Id       int32
	UUID     uuid.UUID
	Type     int32
	Position Vec3
	Yaw      float32
	Pitch    float32
	HeadYaw  float32
	Velocity Vec3
	OnGround bool

	codec VecDeltaCodec
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
}

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

// Start registers the entity handlers on c and overrides the stale 26.4 packet
// decoders.
func Start(c *proto.Conn) (m *Module) {
	m = NewModule(c)
	// The overridden packet ids are resolved by name so this stays correct if
	// the generated id table changes.
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

	m.Register(m.onSpawnEntity)
	m.Register(m.onMoveEntityPos)
	m.Register(m.onMoveEntityPosRot)
	m.Register(m.onMoveEntityRot)
	m.Register(m.onEntityPositionSync)
	m.Register(m.onTeleportEntity)
	m.Register(m.onEntityDestroy)
	m.Register(m.onEntityHeadRotation)
	m.Register(m.onEntityVelocity)
	return
}

func (m *Module) add(entity *Entity) {
	if _, exists := m.entities[entity.Id]; !exists {
		m.order = append(m.order, entity.Id)
	}
	m.entities[entity.Id] = entity
}

func (m *Module) onSpawnEntity(p *v26_4_snapshot_2.PlayToClientPacketSpawnEntity) (err error) {
	entity := &Entity{
		Id:       p.EntityId,
		UUID:     p.ObjectUUID,
		Type:     p.Type,
		Position: Vec3{p.X, p.Y, p.Z},
		Yaw:      unpackDegrees(p.Yaw),
		Pitch:    unpackDegrees(p.Pitch),
		HeadYaw:  unpackDegrees(p.HeadPitch),
		Velocity: Vec3{p.Velocity.X, p.Velocity.Y, p.Velocity.Z},
	}
	entity.codec.SetBase(entity.Position)
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
		entity.HeadYaw = unpackDegrees(p.HeadYaw)
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
		entity.Position = path.EndPosition
		entity.codec.SetBase(path.EndPosition)
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onMoveEntityPosRot(p *PlayToClientPacketMoveEntityPosRot) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		path := entity.codec.decodePath(p.Delta)
		entity.Position = path.EndPosition
		entity.codec.SetBase(path.EndPosition)
		entity.Yaw = unpackDegrees(p.YRot)
		entity.Pitch = unpackDegrees(p.XRot)
		entity.HeadYaw = entity.Yaw
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onMoveEntityRot(p *PlayToClientPacketMoveEntityRot) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		entity.Yaw = unpackDegrees(p.YRot)
		entity.Pitch = unpackDegrees(p.XRot)
		entity.HeadYaw = entity.Yaw
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onEntityPositionSync(p *PlayToClientPacketEntityPositionSync) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		entity.Position = p.Position.EndPosition
		entity.codec.SetBase(p.Position.EndPosition)
		entity.Yaw = p.YRot
		entity.Pitch = p.XRot
		entity.HeadYaw = p.YRot
		entity.OnGround = p.OnGround
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) onTeleportEntity(p *PlayToClientPacketTeleportEntity) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		change := calculateAbsolute(PositionMoveRotation{
			Position:      entity.Position,
			DeltaMovement: entity.Velocity,
			YRot:          entity.Yaw,
			XRot:          entity.Pitch,
		}, p.Change, p.Relatives)
		entity.Position = change.Position
		entity.codec.SetBase(change.Position)
		entity.Velocity = change.DeltaMovement
		entity.Yaw = change.YRot
		entity.Pitch = change.XRot
		entity.HeadYaw = change.YRot
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
