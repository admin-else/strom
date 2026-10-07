package entity

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

func putMetadataByte(b *bytes.Buffer, index uint8, value int8) {
	b.WriteByte(index)
	_ = proto_base.EncodeVarInt(b, 0) // EntityDataSerializers.BYTE
	b.WriteByte(byte(value))
}

func putMetadataBool(b *bytes.Buffer, index uint8, value bool) {
	b.WriteByte(index)
	_ = proto_base.EncodeVarInt(b, 8) // EntityDataSerializers.BOOLEAN
	if value {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func TestEntityMetadataTracked(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 42)
	putMetadataByte(&b, MetadataSharedFlags, 64) // FLAG_GLOWING
	putMetadataBool(&b, MetadataBaby, true)
	putMetadataByte(&b, MetadataSheepWool, 5|16) // lime, sheared
	b.WriteByte(255)

	var packet v26_4_snapshot_2.PlayToClientPacketEntityMetadata
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	m := NewModule(nil)
	m.entities[42] = &Entity{Id: 42}
	if err := m.onEntityMetadata(&packet); err != nil {
		t.Fatal(err)
	}
	entity := m.entities[42]

	if flags, ok := entity.SharedFlags(); !ok || flags != 64 {
		t.Fatalf("shared flags = %v (ok %v), want 64", flags, ok)
	}
	if baby, ok := entity.Baby(); !ok || !baby {
		t.Fatalf("baby = %v (ok %v), want true", baby, ok)
	}
	color, sheared, ok := entity.SheepWool()
	if !ok || color != 5 || !sheared {
		t.Fatalf("wool = color %d sheared %v (ok %v), want 5 true", color, sheared, ok)
	}
	if _, ok := entity.Metadata(MetadataHealth); ok {
		t.Fatalf("unexpected health metadata at index %d", MetadataHealth)
	}
}

func TestEntityMetadataReplacesEntries(t *testing.T) {
	m := NewModule(nil)
	entity := &Entity{Id: 1}
	m.entities[1] = entity

	first := &v26_4_snapshot_2.PlayToClientPacketEntityMetadata{
		EntityId: 1,
		Metadata: v26_4_snapshot_2.EntityMetadata{Val: map[uint8]v26_4_snapshot_2.EntityMetadataEntry{
			MetadataSharedFlags: {Key: MetadataSharedFlags, Type: "byte", Value: int8(1)},
			MetadataBaby:        {Key: MetadataBaby, Type: "boolean", Value: true},
		}},
	}
	if err := m.onEntityMetadata(first); err != nil {
		t.Fatal(err)
	}

	second := &v26_4_snapshot_2.PlayToClientPacketEntityMetadata{
		EntityId: 1,
		Metadata: v26_4_snapshot_2.EntityMetadata{Val: map[uint8]v26_4_snapshot_2.EntityMetadataEntry{
			MetadataSharedFlags: {Key: MetadataSharedFlags, Type: "byte", Value: int8(2)},
		}},
	}
	if err := m.onEntityMetadata(second); err != nil {
		t.Fatal(err)
	}

	if flags, _ := entity.SharedFlags(); flags != 2 {
		t.Fatalf("shared flags = %d, want the replaced value 2", flags)
	}
	if baby, ok := entity.Baby(); !ok || !baby {
		t.Fatalf("baby entry was lost after an unrelated update")
	}
}

func TestEntityMetadataUnknownEntityIgnored(t *testing.T) {
	m := NewModule(nil)
	packet := &v26_4_snapshot_2.PlayToClientPacketEntityMetadata{
		EntityId: 99,
		Metadata: v26_4_snapshot_2.EntityMetadata{Val: map[uint8]v26_4_snapshot_2.EntityMetadataEntry{
			MetadataSharedFlags: {Key: MetadataSharedFlags, Type: "byte", Value: int8(1)},
		}},
	}
	if err := m.onEntityMetadata(packet); err != nil {
		t.Fatal(err)
	}
}

// putEquipmentSlot appends one equipment entry. The continuation bit (0x80) is
// set when more entries follow, matching ClientboundSetEquipmentPacket.write.
func putEquipmentSlot(t *testing.T, b *bytes.Buffer, slot int8, itemCount, itemId int32, more bool) {
	t.Helper()
	if more {
		b.WriteByte(byte(slot) | 128)
	} else {
		b.WriteByte(byte(slot))
	}
	putVarInt(t, b, itemCount)
	if itemCount > 0 {
		putVarInt(t, b, itemId)
		putVarInt(t, b, 0) // added component count
		putVarInt(t, b, 0) // removed component count
	}
}

func TestEntityEquipmentTracked(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 7)
	putEquipmentSlot(t, &b, int8(EquipmentMainHand), 1, 100, true)
	putEquipmentSlot(t, &b, int8(EquipmentHead), 1, 200, true)
	putEquipmentSlot(t, &b, int8(EquipmentFeet), 0, 0, false)

	var packet v26_4_snapshot_2.PlayToClientPacketEntityEquipment
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	m := NewModule(nil)
	m.entities[7] = &Entity{Id: 7}
	if err := m.onEntityEquipment(&packet); err != nil {
		t.Fatal(err)
	}
	entity := m.entities[7]

	hand, ok := entity.Equipment(EquipmentMainHand)
	if !ok || hand.ItemId != 100 || hand.Count != 1 {
		t.Fatalf("main hand = %+v (ok %v), want item 100 count 1", hand, ok)
	}
	head, ok := entity.Equipment(EquipmentHead)
	if !ok || head.ItemId != 200 || head.Count != 1 {
		t.Fatalf("head = %+v (ok %v), want item 200 count 1", head, ok)
	}
	feet, ok := entity.Equipment(EquipmentFeet)
	if !ok || feet.ItemId != 0 || feet.Count != 0 {
		t.Fatalf("feet = %+v (ok %v), want an empty tracked slot", feet, ok)
	}
	if _, ok := entity.Equipment(EquipmentChest); ok {
		t.Fatalf("chest slot should not be tracked")
	}
}

func equipmentPacket(entityId int32, slot EquipmentSlot, itemId int32) (ret *v26_4_snapshot_2.PlayToClientPacketEntityEquipment) {
	return &v26_4_snapshot_2.PlayToClientPacketEntityEquipment{
		EntityId: entityId,
		Equipments: []struct {
			Slot int8
			Item v26_4_snapshot_2.Slot
		}{
			{Slot: int8(slot), Item: v26_4_snapshot_2.Slot{ItemCount: 1, Anon: slotItemData{ItemId: itemId}}},
		},
	}
}

func TestEntityEquipmentUpdatesSlot(t *testing.T) {
	m := NewModule(nil)
	entity := &Entity{Id: 3}
	m.entities[3] = entity
	if err := m.onEntityEquipment(equipmentPacket(3, EquipmentHead, 10)); err != nil {
		t.Fatal(err)
	}
	if head, _ := entity.Equipment(EquipmentHead); head.ItemId != 10 {
		t.Fatalf("head item = %d, want 10", head.ItemId)
	}
	if err := m.onEntityEquipment(equipmentPacket(3, EquipmentHead, 20)); err != nil {
		t.Fatal(err)
	}
	if head, _ := entity.Equipment(EquipmentHead); head.ItemId != 20 {
		t.Fatalf("head item = %d, want the replaced value 20", head.ItemId)
	}
}

func TestEntityEquipmentOutOfRangeSlotIgnored(t *testing.T) {
	m := NewModule(nil)
	entity := &Entity{Id: 4}
	m.entities[4] = entity
	if err := m.onEntityEquipment(equipmentPacket(4, EquipmentSlot(100), 5)); err != nil {
		t.Fatal(err)
	}
	if _, ok := entity.Equipment(EquipmentMainHand); ok {
		t.Fatalf("no slot should be tracked for an out-of-range slot id")
	}
}
