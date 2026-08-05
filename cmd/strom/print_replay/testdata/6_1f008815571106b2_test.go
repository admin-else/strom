package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v1_8"
)

// Error: to do
func TestFailedPacket1F008815571106B2(t *testing.T) {
	p := &generated.PlayToClientPacketEntityMetadata{}
	b := bytes.NewBuffer([]byte{0xd3, 0x14, 0x21, 0x1, 0x2c, 0x7f})
	err := p.Decode(b)
	if err == nil {
		t.Fatal("expected decode error %!q(MISSING) but got nil", "to do")
	}
}
