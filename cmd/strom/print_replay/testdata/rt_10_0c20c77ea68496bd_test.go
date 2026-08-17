package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v1_8"
)

// Decode regression: entity_equipment (1.8) - nbt TAG_End fix
func TestRegressionPacket0C20C77EA68496BD(t *testing.T) {
	p := &generated.PlayToClientPacketEntityEquipment{}
	b := bytes.NewReader([]byte{0x89, 0x16, 0x0, 0x0, 0x1, 0x5, 0x1, 0x0, 0x0, 0x0})
	err := p.Decode(b)
	if err != nil {
		t.Fatalf("regression: decode failed: %v", err)
	}
}
