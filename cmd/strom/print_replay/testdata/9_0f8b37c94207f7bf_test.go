package failed_packets_test

import (
	"bytes"
	"testing"

	generated "git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v26_2"
)

// Error: unexpected EOF
func TestFailedPacket0F8B37C94207F7BF(t *testing.T) {
	p := &generated.PlayToClientPacketUpdateTime{}
	b := bytes.NewReader([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x42, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
