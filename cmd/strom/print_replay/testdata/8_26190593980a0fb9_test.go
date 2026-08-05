package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v26_2"
)

// Error: to do
func TestFailedPacket26190593980A0FB9(t *testing.T) {
	p := &generated.PlayToClientPacketEntityEquipment{}
	b := bytes.NewReader([]byte{0xd3, 0x1, 0x0, 0x1, 0x9a, 0x7, 0x0, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
