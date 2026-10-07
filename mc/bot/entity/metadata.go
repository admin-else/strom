package entity

import "github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"

// Entity metadata indices, reconstructed from the SynchedEntityData.defineId
// call order of the class chain Entity -> LivingEntity -> Mob -> AgeableMob ->
// Animal. The index is stable per class, so a sheep's wool is always the index
// after AgeableMob's two fields.
const (
	MetadataSharedFlags       uint8 = 0
	MetadataAirSupply         uint8 = 1
	MetadataCustomName        uint8 = 2
	MetadataCustomNameVisible uint8 = 3
	MetadataSilent            uint8 = 4
	MetadataNoGravity         uint8 = 5
	MetadataPose              uint8 = 6
	MetadataTicksFrozen       uint8 = 7

	MetadataLivingFlags     uint8 = 8
	MetadataHealth          uint8 = 9
	MetadataEffectParticles uint8 = 10
	MetadataEffectAmbience  uint8 = 11
	MetadataArrowCount      uint8 = 12
	MetadataStingerCount    uint8 = 13
	MetadataSleepingPos     uint8 = 14

	MetadataMobFlags uint8 = 15

	MetadataBaby      uint8 = 16
	MetadataAgeLocked uint8 = 17

	MetadataSheepWool uint8 = 18
)

// Shared flag bit positions mirror Entity.FLAG_*.
const (
	SharedFlagOnFire     = 1 << 0
	SharedFlagSneaking   = 1 << 1
	SharedFlagSprinting  = 1 << 3
	SharedFlagSwimming   = 1 << 4
	SharedFlagInvisible  = 1 << 5
	SharedFlagGlowing    = 1 << 6
	SharedFlagFallFlying = 1 << 7
)

// LivingEntity flag bits mirror LivingEntity.LIVING_ENTITY_FLAG_*.
const (
	LivingFlagUsingItem  = 1 << 0
	LivingFlagOffHand    = 1 << 1
	LivingFlagSpinAttack = 1 << 2
)

// Mob flag bits mirror Mob.MOB_FLAG_*.
const (
	MobFlagNoAI       = 1 << 0
	MobFlagLeftHanded = 1 << 1
	MobFlagAggressive = 1 << 2
)

// Metadata returns the raw tracked metadata value at index.
func (e *Entity) Metadata(index uint8) (value any, ok bool) {
	if e.metadata == nil {
		return
	}
	value, ok = e.metadata[index]
	return
}

// MetadataByte returns the tracked byte metadata at index.
func (e *Entity) MetadataByte(index uint8) (value int8, ok bool) {
	raw, ok := e.Metadata(index)
	if !ok {
		return
	}
	value, ok = raw.(int8)
	return
}

// MetadataBool returns the tracked boolean metadata at index.
func (e *Entity) MetadataBool(index uint8) (value bool, ok bool) {
	raw, ok := e.Metadata(index)
	if !ok {
		return
	}
	value, ok = raw.(bool)
	return
}

// MetadataInt returns the tracked int metadata at index.
func (e *Entity) MetadataInt(index uint8) (value int32, ok bool) {
	raw, ok := e.Metadata(index)
	if !ok {
		return
	}
	value, ok = raw.(int32)
	return
}

// MetadataFloat returns the tracked float metadata at index.
func (e *Entity) MetadataFloat(index uint8) (value float32, ok bool) {
	raw, ok := e.Metadata(index)
	if !ok {
		return
	}
	value, ok = raw.(float32)
	return
}

// SharedFlags returns the shared entity flag byte (on fire, sneaking, ...).
func (e *Entity) SharedFlags() (flags int8, ok bool) {
	return e.MetadataByte(MetadataSharedFlags)
}

// LivingFlags returns the LivingEntity flag byte (using item, off-hand, ...).
func (e *Entity) LivingFlags() (flags int8, ok bool) {
	return e.MetadataByte(MetadataLivingFlags)
}

// Baby reports the AgeableMob baby flag.
func (e *Entity) Baby() (baby bool, ok bool) {
	return e.MetadataBool(MetadataBaby)
}

// SheepWool returns the sheep wool colour id (0..15) and the sheared flag.
func (e *Entity) SheepWool() (color uint8, sheared bool, ok bool) {
	raw, ok := e.MetadataByte(MetadataSheepWool)
	if !ok {
		return
	}
	color = uint8(raw) & 15
	sheared = raw&16 != 0
	return
}

func (m *Module) onEntityMetadata(p *v26_4_snapshot_2.PlayToClientPacketEntityMetadata) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		if entity.metadata == nil {
			entity.metadata = make(map[uint8]any)
		}
		for index, entry := range p.Metadata.Val {
			entity.metadata[index] = entry.Value
		}
	}
	m.mu.Unlock()
	return nil
}
