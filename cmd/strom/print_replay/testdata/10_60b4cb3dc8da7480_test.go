package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v1_8"
)

// Error: EOF
func TestFailedPacket60B4CB3DC8DA7480(t *testing.T) {
	p := &generated.PlayToClientPacketEntityEquipment{}
	b := bytes.NewBuffer([]byte{0x8e, 0x1c, 0x0, 0x0, 0x1, 0x5, 0x1, 0x0, 0x0, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatal("expected decode error %!q(MISSING) but got nil", "EOF")
	}
}
