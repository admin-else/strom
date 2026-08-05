package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v1_8"
)

// Error: to do
func TestFailedPacket00F88F767DDBAF62(t *testing.T) {
	p := &generated.PlayToClientPacketEntityMetadata{}
	b := bytes.NewBuffer([]byte{0x3d, 0x0, 0x0, 0x21, 0x1, 0x2c, 0x82, 0x0, 0x3, 0x0, 0x4, 0x0, 0xaa, 0x0, 0x25, 0x1, 0x0, 0x0, 0x0, 0x7f})
	err := p.Decode(b)
	if err == nil {
		t.Fatal("expected decode error %!q(MISSING) but got nil", "to do")
	}
}
