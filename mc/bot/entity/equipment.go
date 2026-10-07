package entity

import "github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"

// EquipmentSlot mirrors net.minecraft.world.entity.EquipmentSlot in the order
// ClientboundSetEquipmentPacket writes it (EquipmentSlot.ordinal()). The wire
// byte packs the slot in the low 7 bits and the continuation flag in bit 7.
type EquipmentSlot int8

const (
	EquipmentMainHand EquipmentSlot = 0
	EquipmentOffHand  EquipmentSlot = 1
	EquipmentFeet     EquipmentSlot = 2
	EquipmentLegs     EquipmentSlot = 3
	EquipmentChest    EquipmentSlot = 4
	EquipmentHead     EquipmentSlot = 5
	EquipmentBody     EquipmentSlot = 6
	EquipmentSaddle   EquipmentSlot = 7

	equipmentSlotCount = 8
)

// ItemStack is the tracked equipment slot contents: the raw item id, the stack
// count and the added data components. Removed components are not tracked.
type ItemStack struct {
	ItemId     int32
	Count      int32
	Components []v26_4_snapshot_2.SlotComponent
}

// slotItemData is the anonymous decoded Slot payload for a non-empty stack; it
// is a type alias so the generated Slot.Anon value type-asserts to it.
type slotItemData = struct {
	ItemId                int32
	AddedComponentCount   int32
	RemovedComponentCount int32
	Components            []v26_4_snapshot_2.SlotComponent
	RemoveComponents      []struct {
		Type v26_4_snapshot_2.SlotComponentType
	}
}

// itemStackFromSlot converts a decoded slot into an ItemStack. It reports false
// for an empty slot.
func itemStackFromSlot(slot v26_4_snapshot_2.Slot) (ret ItemStack, ok bool) {
	if slot.ItemCount <= 0 {
		return
	}
	data, isItem := slot.Anon.(slotItemData)
	if !isItem {
		return
	}
	ret = ItemStack{ItemId: data.ItemId, Count: slot.ItemCount, Components: data.Components}
	ok = true
	return
}

// Equipment returns the tracked item in slot. It reports false when the slot
// was never received and true when it was (even when the slot is empty).
func (e *Entity) Equipment(slot EquipmentSlot) (stack ItemStack, ok bool) {
	if slot < 0 || int(slot) >= equipmentSlotCount {
		return
	}
	return e.equipment[slot], e.equipmentSet[slot]
}

func (m *Module) onEntityEquipment(p *v26_4_snapshot_2.PlayToClientPacketEntityEquipment) (err error) {
	m.mu.Lock()
	if entity, ok := m.entities[p.EntityId]; ok {
		for _, equipment := range p.Equipments {
			slot := EquipmentSlot(equipment.Slot & 0x7f)
			if slot < 0 || int(slot) >= equipmentSlotCount {
				continue
			}
			stack, _ := itemStackFromSlot(equipment.Item)
			entity.equipment[slot] = stack
			entity.equipmentSet[slot] = true
		}
	}
	m.mu.Unlock()
	return nil
}
