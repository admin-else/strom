package failed_packets_test

import (
	"bytes"
	"testing"

	generated "git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v26_2"
)

// Error: unexpected EOF
func TestFailedPacket6C9AD835672C5C51(t *testing.T) {
	p := &generated.PlayToClientPacketUpdateTime{}
	b := bytes.NewReader([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0xde, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
